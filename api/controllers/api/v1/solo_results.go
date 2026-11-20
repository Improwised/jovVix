package v1

import (
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// SoloResults ranks each player's best completed attempt under identical rules.
func (qc *quizSocketController) SoloResults(c *fiber.Ctx) error {
	user, ok := c.Locals(constants.ContextUser).(models.User)
	if !ok {
		return utils.JSONFail(c, 401, "Authentication required")
	}
	id, err := uuid.Parse(c.Params("attempt"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid attempt")
	}
	var attempt struct {
		QuizID         string `db:"quiz_id" json:"quiz_id"`
		Version        string `db:"quiz_version" json:"-"`
		Title          string `db:"title" json:"title"`
		Score          int    `db:"score" json:"score"`
		Correct        int    `db:"correct" json:"correct"`
		Total          int    `db:"total" json:"total"`
		TotalQuestions int    `db:"total_questions" json:"total_questions"`
	}
	found, err := qc.db.ScanStruct(&attempt, `SELECT a.quiz_id::text,a.quiz_version,a.title,COALESCE(sum(r.calculated_score),0)::int score,count(*) FILTER(WHERE r.calculated_score>0)::int correct,count(r.id) FILTER(WHERE r.is_attend=true)::int total,count(r.id)::int total_questions FROM active_quizzes a JOIN user_played_quizzes p ON p.active_quiz_id=a.id LEFT JOIN user_quiz_responses r ON r.user_played_quiz_id=p.id JOIN quizzes q ON q.id=a.quiz_id WHERE p.id=$1 AND p.user_id=$2 AND a.is_solo AND a.solo_completed AND NOT a.is_active AND q.is_public GROUP BY a.id`, id, user.ID)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load results")
	}
	if !found {
		return utils.JSONFail(c, 404, "Completed solo attempt not found")
	}
	type entry struct {
		UserID  string `db:"user_id" json:"-"`
		Name    string `db:"name" json:"name"`
		Score   int    `db:"score" json:"score"`
		Correct int    `db:"correct" json:"correct"`
		Rank    int    `db:"rank" json:"rank"`
		ImgKey  string `db:"img_key" json:"img_key"`
		IsYou   bool   `json:"is_you"`
	}
	var players []entry
	err = qc.db.ScanStructs(&players, `WITH attempts AS (SELECT p.user_id,COALESCE(p.player_name,u.first_name) name,COALESCE(sum(r.calculated_score),0)::int score,count(*) FILTER(WHERE r.calculated_score>0)::int correct,COALESCE(u.img_key,'') img_key,p.created_at,p.id FROM active_quizzes a JOIN user_played_quizzes p ON p.active_quiz_id=a.id JOIN users u ON u.id=p.user_id LEFT JOIN user_quiz_responses r ON r.user_played_quiz_id=p.id WHERE a.quiz_id=$1 AND a.quiz_version=$2 AND a.is_solo AND a.solo_completed AND NOT a.is_active GROUP BY p.id,u.first_name,u.img_key), best AS (SELECT DISTINCT ON(user_id) * FROM attempts ORDER BY user_id,score DESC,created_at,id) SELECT user_id,name,score,correct,img_key,rank() OVER(ORDER BY score DESC)::int rank FROM best ORDER BY score DESC,created_at,id`, attempt.QuizID, attempt.Version)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load leaderboard")
	}
	for i := range players {
		players[i].IsYou = players[i].UserID == user.ID
	}
	return utils.JSONSuccess(c, 200, fiber.Map{"attempt": attempt, "players": players})
}
