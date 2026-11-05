package v1

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HostParticipation stores the host's participant row separately from host authority,
// which remains attached to active_quizzes.admin_id.
func (qc *quizSocketController) HostParticipation(c *fiber.Ctx) error {
	user, ok := c.Locals(constants.ContextUser).(models.User)
	if !ok {
		return utils.JSONFail(c, 401, "Authentication required")
	}
	sessionID, err := uuid.Parse(c.Query("active_quiz_id"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid session")
	}
	tx, err := qc.db.Begin()
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load session")
	}
	defer tx.Rollback()
	var session struct {
		AdminID        string         `db:"admin_id"`
		Current        sql.NullString `db:"current_question"`
		QuestionActive sql.NullBool   `db:"is_question_active"`
		Active         bool           `db:"is_active"`
		End            sql.NullTime   `db:"activated_to"`
	}
	found, err := tx.From("active_quizzes").Select("admin_id", "current_question", "is_question_active", "is_active", "activated_to").Where(goqu.Ex{"id": sessionID}).ForUpdate(exp.Wait).ScanStruct(&session)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load session")
	}
	if !found || session.AdminID != user.ID {
		return utils.JSONFail(c, 403, "Only the host can change participation")
	}
	var participant uuid.UUID
	playerName := user.FirstName
	var savedName sql.NullString
	err = tx.QueryRow(`SELECT id, player_name FROM user_played_quizzes WHERE active_quiz_id=$1 AND user_id=$2 AND is_host=false LIMIT 1`, sessionID, user.ID).Scan(&participant, &savedName)
	found = err == nil
	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}
	if savedName.Valid {
		playerName = savedName.String
	}
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load participation")
	}
	if c.Method() == fiber.MethodPut {
		if session.QuestionActive.Valid || session.Current.Valid || session.End.Valid || !session.Active {
			return utils.JSONFail(c, 409, "Participation can only change in the lobby")
		}
		var body struct {
			Enabled    *bool  `json:"enabled"`
			PlayerName string `json:"player_name"`
		}
		if c.BodyParser(&body) != nil || body.Enabled == nil {
			return utils.JSONFail(c, 400, "Invalid participation option")
		}
		if *body.Enabled && !found {
			var suggestions []string
			playerName, suggestions, err = models.ResolveQuizPlayerName(tx, sessionID, user.FirstName, user.LastName, body.PlayerName)
			if errors.Is(err, models.ErrPlayerNameTaken) {
				return utils.JSONFail(c, http.StatusConflict, fiber.Map{"message": "This name is already taken. Choose another name for this quiz.", "suggestions": suggestions})
			}
			if errors.Is(err, models.ErrInvalidPlayerName) {
				return utils.JSONFail(c, http.StatusBadRequest, err.Error())
			}
			if err != nil {
				return utils.JSONFail(c, http.StatusInternalServerError, "Unable to validate quiz name")
			}
			participant = uuid.New()
			_, err = tx.Insert("user_played_quizzes").Rows(goqu.Record{"id": participant, "user_id": user.ID, "active_quiz_id": sessionID, "is_host": false, "player_name": playerName}).Executor().Exec()
			if err != nil {
				return utils.JSONFail(c, 500, "Unable to add host as player")
			}
			var ids []uuid.UUID
			err = tx.From("active_quiz_questions").Select("question_id").Where(goqu.Ex{"active_quiz_id": sessionID}).ScanVals(&ids)
			if err != nil {
				return utils.JSONFail(c, 500, "Unable to prepare questions")
			}
			for _, id := range ids {
				_, err = tx.Insert("user_quiz_responses").Rows(goqu.Record{"id": uuid.New(), "question_id": id, "user_played_quiz_id": participant, "option_order": nil}).Executor().Exec()
				if err != nil {
					return utils.JSONFail(c, 500, "Unable to prepare answers")
				}
			}
			found = true
		} else if !*body.Enabled && found {
			_, err = tx.Delete("user_quiz_responses").Where(goqu.Ex{"user_played_quiz_id": participant}).Executor().Exec()
			if err == nil {
				_, err = tx.Delete("user_played_quizzes").Where(goqu.Ex{"id": participant}).Executor().Exec()
			}
			if err != nil {
				return utils.JSONFail(c, 500, "Unable to remove host as player")
			}
			found = false
		}
	}
	var answer []int
	if found && session.Current.Valid && session.QuestionActive.Valid && session.QuestionActive.Bool {
		var raw sql.NullString
		_, err = tx.From("user_quiz_responses").Select("answers").Where(goqu.Ex{"user_played_quiz_id": participant, "question_id": session.Current.String}).ScanVal(&raw)
		if err != nil {
			return utils.JSONFail(c, 500, "Unable to load answer")
		}
		if raw.Valid {
			_ = json.Unmarshal([]byte(raw.String), &answer)
		}
	}
	var submitted []struct {
		ID        string `db:"id" json:"id"`
		FirstName string `db:"first_name" json:"first_name"`
		Username  string `db:"username" json:"username"`
		ImageKey  string `db:"img_key" json:"img_key"`
	}
	if session.Current.Valid && session.QuestionActive.Valid && session.QuestionActive.Bool {
		err = tx.From(goqu.T("user_quiz_responses").As("r")).
			Join(goqu.T("user_played_quizzes").As("p"), goqu.On(goqu.I("p.id").Eq(goqu.I("r.user_played_quiz_id")))).
			Join(goqu.T("users").As("u"), goqu.On(goqu.I("u.id").Eq(goqu.I("p.user_id")))).
			Select(goqu.I("u.id"), goqu.L("COALESCE(p.player_name, u.first_name)").As("first_name"), goqu.I("u.username"), goqu.I("u.img_key")).
			Where(goqu.Ex{"p.active_quiz_id": sessionID, "p.is_host": false, "r.question_id": session.Current.String, "r.is_attend": true}).ScanStructs(&submitted)
		if err != nil {
			return utils.JSONFail(c, 500, "Unable to load submitted answers")
		}
	}
	if err = tx.Commit(); err != nil {
		return utils.JSONFail(c, 500, "Unable to save participation")
	}
	result := fiber.Map{"enabled": found, "started": session.QuestionActive.Valid || session.Current.Valid, "question_id": session.Current.String, "name": playerName, "user_id": user.ID, "username": user.Username, "answer": answer, "submitted": submitted}
	if found {
		result["user_played_quiz"] = participant.String()
	}
	return utils.JSONSuccess(c, http.StatusOK, result)
}

// hasHostedParticipants atomically locks participation when the quiz starts.
// NULL means lobby; false means the countdown or results; true means an open question.
func (qc *quizSocketController) hasHostedParticipants(sessionID uuid.UUID) (bool, error) {
	tx, err := qc.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	var state struct {
		Active         bool           `db:"is_active"`
		QuestionActive sql.NullBool   `db:"is_question_active"`
		Current        sql.NullString `db:"current_question"`
	}
	found, err := tx.From("active_quizzes").Select("is_active", "is_question_active", "current_question").Where(goqu.Ex{"id": sessionID}).ForUpdate(exp.Wait).ScanStruct(&state)
	if err != nil {
		return false, err
	}
	if !found || !state.Active {
		return false, fmt.Errorf("active quiz not found")
	}
	if state.QuestionActive.Valid || state.Current.Valid {
		return false, fmt.Errorf("quiz has already started")
	}
	count, err := tx.From("user_played_quizzes").Where(goqu.Ex{"active_quiz_id": sessionID, "is_host": false}).Count()
	if err != nil || count == 0 {
		return false, err
	}
	_, err = tx.Update("active_quizzes").Set(goqu.Record{"is_question_active": false}).Where(goqu.Ex{"id": sessionID}).Executor().Exec()
	if err != nil {
		return false, err
	}
	if err = tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}
