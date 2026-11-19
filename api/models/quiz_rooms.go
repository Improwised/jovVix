package models

import (
	"database/sql"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/doug-martin/goqu/v9"
)

var ErrRoomLeft = errors.New("You have left this room. Join again using its code.")
var ErrRoomClosed = errors.New("This room has ended")

// ResolveRoomMember reserves the display name across all rounds, including away players.
func ResolveRoomMember(db *goqu.Database, roomID, userID, requestedName string) (string, error) {
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var closed bool
	err = tx.QueryRow(`SELECT closed_at IS NOT NULL OR expires_at<=now() FROM quiz_rooms WHERE id=$1 FOR UPDATE`, roomID).Scan(&closed)
	if err != nil {
		return "", err
	}
	if closed {
		return "", ErrRoomClosed
	}
	var name string
	var left sql.NullTime
	err = tx.QueryRow(`SELECT player_name,left_at FROM quiz_room_members WHERE room_id=$1 AND user_id=$2`, roomID, userID).Scan(&name, &left)
	if err == nil {
		if left.Valid {
			return "", ErrRoomLeft
		}
		return name, nil
	}
	if err != sql.ErrNoRows {
		return "", err
	}
	name = strings.TrimSpace(requestedName)
	if name == "" {
		if err = tx.QueryRow(`SELECT first_name FROM users WHERE id=$1`, userID).Scan(&name); err != nil {
			return "", err
		}
		name = strings.TrimSpace(name)
		if len([]rune(name)) > 12 {
			name = string([]rune(name)[:12])
		}
	}
	if utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 12 {
		return "", ErrInvalidPlayerName
	}
	var taken bool
	if err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM quiz_room_members WHERE room_id=$1 AND lower(btrim(player_name))=lower($2))`, roomID, name).Scan(&taken); err != nil {
		return "", err
	}
	if taken {
		return "", ErrPlayerNameTaken
	}
	if _, err = tx.Exec(`INSERT INTO quiz_room_members(room_id,user_id,player_name) VALUES($1,$2,$3)`, roomID, userID, name); err != nil {
		return "", err
	}
	return name, tx.Commit()
}
