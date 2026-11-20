package v1

import (
	"database/sql"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type roomState struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	AdminID   string       `db:"admin_id" json:"-"`
	Code      int          `db:"code" json:"code"`
	SessionID uuid.UUID    `db:"current_session_id" json:"session_id"`
	Closed    sql.NullTime `db:"closed_at" json:"-"`
	HostSeen  time.Time    `db:"host_seen_at" json:"-"`
	Expires   time.Time    `db:"expires_at" json:"expires_at"`
	Title     string       `db:"title" json:"title"`
	Active    bool         `db:"is_active" json:"-"`
	Started   bool         `db:"started" json:"-"`
}

func roomBy(tx *goqu.TxDatabase, field string, value any, lock bool) (roomState, error) {
	var r roomState
	// Lock the room in its own statement. A join with FOR UPDATE can lose
	// its row during PostgreSQL's recheck when another selector changes the
	// current session while this transaction waits for the lock.
	if lock {
		var locked uuid.UUID
		if err := tx.QueryRow(`SELECT r.id FROM quiz_rooms r WHERE `+field+`=$1 FOR UPDATE`, value).Scan(&locked); err != nil {
			return r, err
		}
		field, value = "r.id", locked
	}
	query := `SELECT r.id,r.admin_id,r.code,r.current_session_id,r.closed_at,r.expires_at,r.host_seen_at,COALESCE(a.title,'') AS title,a.is_active,(a.current_question IS NOT NULL OR a.is_question_active IS NOT NULL) AS started FROM quiz_rooms r JOIN active_quizzes a ON a.id=r.current_session_id WHERE ` + field + `=$1`
	_, err := tx.ScanStruct(&r, query, value)
	if err != nil {
		return r, err
	}
	if r.ID == uuid.Nil {
		return r, sql.ErrNoRows
	}
	r.AdminID = strings.TrimSpace(r.AdminID)
	return r, nil
}
func (r roomState) phase() string {
	if r.Closed.Valid || time.Now().After(r.Expires) {
		return "closed"
	}
	if !r.Active {
		return "results"
	}
	if r.Started {
		return "playing"
	}
	return "lobby"
}

// RoomState is authoritative; clients can recover after missing any number of events.
func (qc *quizSocketController) RoomState(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	tx, err := qc.db.Begin()
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load room")
	}
	defer tx.Rollback()
	field, value := "r.id", any(c.Params("room"))
	if c.Params("room") == "lookup" {
		if c.Query("code") != "" {
			field, value = "r.code", c.Query("code")
		} else {
			var id sql.NullString
			if c.Query("session_id") != "" {
				sessionID, parseErr := uuid.Parse(c.Query("session_id"))
				if parseErr != nil {
					return utils.JSONFail(c, 400, "Invalid session")
				}
				err = tx.QueryRow(`SELECT room_id FROM active_quizzes WHERE id=$1`, sessionID).Scan(&id)
			} else {
				participantID, parseErr := uuid.Parse(c.Query("participant_id"))
				if parseErr != nil {
					return utils.JSONFail(c, 400, "Invalid participant")
				}
				err = tx.QueryRow(`SELECT a.room_id FROM active_quizzes a JOIN user_played_quizzes p ON p.active_quiz_id=a.id WHERE p.id=$1`, participantID).Scan(&id)
			}
			if err != nil || !id.Valid {
				return utils.JSONFail(c, 404, "Room not found")
			}
			value = id.String
		}
	}
	if field == "r.id" {
		if _, err = uuid.Parse(value.(string)); err != nil {
			return utils.JSONFail(c, 400, "Invalid room")
		}
	}

	r, err := roomBy(tx, field, value, false)
	if err != nil {
		return utils.JSONFail(c, 404, "Room not found")
	}
	host := r.AdminID == user.ID
	if host {
		if _, err = tx.Exec(`UPDATE quiz_rooms SET host_seen_at=now() WHERE id=$1`, r.ID); err != nil {
			return utils.JSONFail(c, 500, "Unable to update host presence")
		}
		r.HostSeen = time.Now()
	}

	var member bool
	if err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM quiz_room_members WHERE room_id=$1 AND user_id=$2 AND left_at IS NULL)`, r.ID, user.ID).Scan(&member); err != nil {
		return utils.JSONFail(c, 500, "Unable to load membership")
	}
	if !host && !member {
		return utils.JSONFail(c, 403, "You have left this room. Join again using its code.")
	}
	if member {
		if _, err = tx.Exec(`UPDATE quiz_room_members SET last_seen_at=now() WHERE room_id=$1 AND user_id=$2`, r.ID, user.ID); err != nil {
			return utils.JSONFail(c, 500, "Unable to update presence")
		}
	}
	type memberInfo struct {
		Name   string `db:"player_name" json:"name"`
		Online bool   `db:"online" json:"online"`
	}
	members := []memberInfo{}
	err = tx.ScanStructs(&members, `SELECT player_name,last_seen_at>now()-interval '90 seconds' AS online FROM quiz_room_members WHERE room_id=$1 AND left_at IS NULL ORDER BY player_name`, r.ID)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load players")
	}
	var participant sql.NullString
	err = tx.QueryRow(`SELECT id::text FROM user_played_quizzes WHERE active_quiz_id=$1 AND user_id=$2 AND is_host=false LIMIT 1`, r.SessionID, user.ID).Scan(&participant)
	if err != nil && err != sql.ErrNoRows {
		return utils.JSONFail(c, 500, "Unable to load participation")
	}
	if err = tx.Commit(); err != nil {
		return utils.JSONFail(c, 500, "Unable to load room")
	}
	return utils.JSONSuccess(c, 200, fiber.Map{"id": r.ID, "code": r.Code, "session_id": r.SessionID, "phase": r.phase(), "title": r.Title, "is_host": host, "host_online": time.Since(r.HostSeen) < 90*time.Second, "members": members, "user_played_quiz": participant.String, "username": user.Username})
}

// JoinRoom serializes membership changes with round selection and closure.
func (qc *quizSocketController) JoinRoom(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	var body struct {
		PlayerName string `json:"player_name"`
	}
	if len(c.Body()) > 0 && c.BodyParser(&body) != nil {
		return utils.JSONFail(c, 400, "Invalid join request")
	}
	tx, err := qc.db.Begin()
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to join room")
	}
	defer tx.Rollback()
	r, err := roomBy(tx, "r.code", c.Params("code"), true)
	if err != nil {
		return utils.JSONFail(c, 404, "invitation code not found")
	}
	if r.phase() == "closed" {
		return utils.JSONFail(c, 410, "This room has ended")
	}
	if r.AdminID == user.ID {
		return utils.JSONSuccess(c, 200, fiber.Map{"room_id": r.ID})
	}
	name := strings.TrimSpace(body.PlayerName)
	var saved string
	err = tx.QueryRow(`SELECT player_name FROM quiz_room_members WHERE room_id=$1 AND user_id=$2`, r.ID, user.ID).Scan(&saved)
	if err == nil {
		name = saved
	} else if err != sql.ErrNoRows {
		return utils.JSONFail(c, 500, "Unable to join room")
	}
	if name == "" {
		name = strings.TrimSpace(user.FirstName)
	}
	if utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 12 {
		return utils.JSONFail(c, 400, "Name must contain 1 to 12 characters")
	}
	var taken bool
	err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM quiz_room_members WHERE room_id=$1 AND user_id<>$2 AND lower(btrim(player_name))=lower($3))`, r.ID, user.ID, name).Scan(&taken)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to join room")
	}
	if !taken {
		err = tx.QueryRow(`SELECT EXISTS(SELECT 1 FROM user_played_quizzes WHERE active_quiz_id=$1 AND user_id=$2 AND is_host=false AND lower(btrim(player_name))=lower($3))`, r.SessionID, r.AdminID, name).Scan(&taken)
		if err != nil {
			return utils.JSONFail(c, 500, "Unable to validate player name")
		}
	}
	if taken {
		return utils.JSONFail(c, 409, fiber.Map{"message": "This name is already taken. Choose another name for this room."})
	}
	_, err = tx.Exec(`INSERT INTO quiz_room_members(room_id,user_id,player_name) VALUES($1,$2,$3) ON CONFLICT(room_id,user_id) DO UPDATE SET left_at=NULL,last_seen_at=now(),player_name=excluded.player_name`, r.ID, user.ID, name)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to join room")
	}
	if err = tx.Commit(); err != nil {
		return utils.JSONFail(c, 500, "Unable to join room")
	}
	return utils.JSONSuccess(c, 200, fiber.Map{"room_id": r.ID})
}

func (qc *quizSocketController) LeaveRoom(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	id, err := uuid.Parse(c.Params("room"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid room")
	}
	_, err = qc.db.Exec(`UPDATE quiz_room_members SET left_at=now() WHERE room_id=$1 AND user_id=$2`, id, user.ID)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to leave room")
	}
	return utils.JSONSuccess(c, 200, "Left room")
}

func (qc *quizSocketController) RoomQuizzes(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	id, err := uuid.Parse(c.Params("room"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid room")
	}
	var host bool
	err = qc.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM quiz_rooms WHERE id=$1 AND admin_id=$2 AND closed_at IS NULL AND expires_at>now())`, id, user.ID).Scan(&host)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load quizzes")
	}
	if !host {
		return utils.JSONFail(c, 403, "Only the host can choose a quiz")
	}
	type choice struct {
		ID            uuid.UUID `db:"id" json:"id"`
		Title         string    `db:"title" json:"title"`
		IsPublic      bool      `db:"is_public" json:"is_public"`
		QuestionCount int       `db:"question_count" json:"question_count"`
	}
	quizzes := []choice{}
	err = qc.db.ScanStructs(&quizzes, `SELECT q.id,q.title,COALESCE(q.is_public,false) AS is_public,(SELECT count(*) FROM quiz_questions qq WHERE qq.quiz_id=q.id) AS question_count FROM quizzes q WHERE (q.is_public=true OR q.creator_id=$1 OR ($2<>'' AND EXISTS(SELECT 1 FROM shared_quizzes s WHERE s.quiz_id=q.id AND s.shared_to=$2))) AND EXISTS(SELECT 1 FROM quiz_questions qq WHERE qq.quiz_id=q.id) ORDER BY q.title,q.id`, user.ID, user.Email)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load quizzes")
	}
	return utils.JSONSuccess(c, 200, quizzes)
}

func (qc *quizSocketController) NextRoomQuiz(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	id, err := uuid.Parse(c.Params("room"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid room")
	}
	var body struct {
		QuizID          uuid.UUID `json:"quiz_id"`
		ExpectedSession uuid.UUID `json:"expected_session_id"`
	}
	if c.BodyParser(&body) != nil || body.QuizID == uuid.Nil || body.ExpectedSession == uuid.Nil {
		return utils.JSONFail(c, 400, "Invalid quiz selection")
	}
	tx, err := qc.db.Begin()
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to select quiz")
	}
	defer tx.Rollback()
	r, err := roomBy(tx, "r.id", id, true)
	if err != nil {
		return utils.JSONFail(c, 404, "Room not found")
	}
	if r.AdminID != user.ID {
		return utils.JSONFail(c, 403, "Only the host can choose a quiz")
	}
	if r.phase() == "closed" {
		return utils.JSONFail(c, 410, "This room has ended")
	}
	if r.SessionID != body.ExpectedSession || r.phase() != "results" {
		return utils.JSONFail(c, 409, "The room has changed. Refresh before choosing a quiz.")
	}
	var title string
	err = tx.QueryRow(`SELECT title FROM quizzes q WHERE id=$1 AND (is_public=true OR creator_id=$2 OR ($3<>'' AND EXISTS(SELECT 1 FROM shared_quizzes s WHERE s.quiz_id=q.id AND s.shared_to=$3))) FOR SHARE`, body.QuizID, user.ID, user.Email).Scan(&title)
	if err != nil {
		return utils.JSONFail(c, 403, "This quiz is no longer available to host")
	}
	session := uuid.New()
	_, err = tx.Exec(`INSERT INTO active_quizzes(id,title,quiz_id,admin_id,room_id,invitation_code,is_active,activated_from) VALUES($1,$2,$3,$4,$5,$6,true,now())`, session, title, body.QuizID, user.ID, r.ID, r.Code)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to create next quiz")
	}
	if err = models.CopyRoomQuestions(tx, session, body.QuizID.String()); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.JSONFail(c, 400, "Quiz must contain playable questions")
		}
		if qc.logger != nil {
			qc.logger.Error("unable to copy next room quiz questions", zap.Error(err), zap.String("quiz_id", body.QuizID.String()))
		}
		return utils.JSONFail(c, 500, "Unable to prepare the next quiz. Please try again.")
	}
	// Carry a playing host's preference into the new round with fresh responses.
	var hostName string
	hostErr := tx.QueryRow(`SELECT player_name FROM user_played_quizzes WHERE active_quiz_id=$1 AND user_id=$2 AND is_host=false AND EXISTS(SELECT 1 FROM quizzes WHERE id=$3 AND is_public=true) LIMIT 1`, r.SessionID, user.ID, body.QuizID).Scan(&hostName)
	if hostErr == nil {
		participant := uuid.New()
		if _, err = tx.Exec(`INSERT INTO user_played_quizzes(id,user_id,active_quiz_id,is_host,player_name) VALUES($1,$2,$3,false,$4)`, participant, user.ID, session, hostName); err != nil {
			return utils.JSONFail(c, 500, "Unable to restore host participation")
		}
		if err = models.CopyRoomResponses(tx, participant, session); err != nil {
			return utils.JSONFail(c, 500, "Unable to prepare host answers")
		}
	} else if hostErr != sql.ErrNoRows {
		return utils.JSONFail(c, 500, "Unable to load host participation")
	}
	_, err = tx.Exec(`UPDATE quiz_rooms SET current_session_id=$1,expires_at=now()+interval '24 hours' WHERE id=$2`, session, id)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to select quiz")
	}
	if err = tx.Commit(); err != nil {
		return utils.JSONFail(c, 500, "Unable to select quiz")
	}
	return utils.JSONSuccess(c, 201, fiber.Map{"session_id": session})
}

func (qc *quizSocketController) CloseRoom(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	id, err := uuid.Parse(c.Params("room"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid room")
	}
	tx, err := qc.db.Begin()
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to end room")
	}
	defer tx.Rollback()
	r, err := roomBy(tx, "r.id", id, true)
	if err != nil {
		return utils.JSONFail(c, 404, "Room not found")
	}
	if r.AdminID != user.ID {
		return utils.JSONFail(c, 403, "Only the host can end the room")
	}
	if r.Active {
		return utils.JSONFail(c, 409, "End the current quiz before ending the room")
	}
	if _, err = tx.Exec(`UPDATE quiz_rooms SET closed_at=COALESCE(closed_at,now()) WHERE id=$1`, id); err != nil {
		return utils.JSONFail(c, 500, "Unable to end room")
	}
	if err = tx.Commit(); err != nil {
		return utils.JSONFail(c, 500, "Unable to end room")
	}
	return utils.JSONSuccess(c, 200, "Room ended")
}
