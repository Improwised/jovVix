package v1

import (
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

func (qc *quizSocketController) RoomReports(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	rows := []struct {
		ID           string     `db:"id" json:"id"`
		Code         int        `db:"code" json:"code"`
		Ended        bool       `db:"ended" json:"ended"`
		Rounds       int        `db:"rounds" json:"rounds"`
		Participants int        `db:"participants" json:"participants"`
		Started      time.Time  `db:"started" json:"started"`
		Finished     *time.Time `db:"finished" json:"finished"`
	}{}
	err := qc.db.ScanStructs(&rows, `SELECT r.id::text,r.code,r.created_at AT TIME ZONE 'UTC' AS started,LEAST(r.closed_at,CASE WHEN r.expires_at<=now() THEN r.expires_at END) AT TIME ZONE 'UTC' AS finished,(SELECT count(DISTINCT p.user_id) FROM active_quizzes aq JOIN user_played_quizzes p ON p.active_quiz_id=aq.id AND p.is_host=false WHERE aq.room_id=r.id) AS participants,(r.closed_at IS NOT NULL OR r.expires_at<=now()) AS ended,count(a.id) AS rounds FROM quiz_rooms r LEFT JOIN active_quizzes a ON a.room_id=r.id AND a.is_active=false WHERE r.admin_id=$1 GROUP BY r.id ORDER BY r.created_at DESC`, user.ID)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load room reports")
	}
	return utils.JSONSuccess(c, 200, rows)
}
func (qc *quizSocketController) RoomReport(c *fiber.Ctx) error {
	user := c.Locals(constants.ContextUser).(models.User)
	id, err := uuid.Parse(c.Params("room"))
	if err != nil {
		return utils.JSONFail(c, 400, "Invalid room")
	}
	// Leaving revokes ongoing standings access; former members retain final reports.
	var allowed, host, ended bool
	err = qc.db.QueryRow(`SELECT r.admin_id=$2,(r.closed_at IS NOT NULL OR r.expires_at<=now()),(r.admin_id=$2 OR EXISTS(SELECT 1 FROM quiz_room_members m WHERE m.room_id=r.id AND m.user_id=$2 AND (m.left_at IS NULL OR r.closed_at IS NOT NULL OR r.expires_at<=now()))) FROM quiz_rooms r WHERE r.id=$1`, id, user.ID).Scan(&host, &ended, &allowed)
	if err != nil {
		return utils.JSONFail(c, 404, "Room not found")
	}
	if !allowed {
		return utils.JSONFail(c, 403, "You cannot view this room report")
	}
	standings := []struct {
		Name     string `db:"name" json:"name"`
		Username string `db:"username" json:"username"`
		ImageKey string `db:"img_key" json:"img_key"`
		Score    int64  `db:"score" json:"score"`
		Rank     int64  `db:"rank" json:"rank"`
	}{}
	err = qc.db.ScanStructs(&standings, `WITH totals AS (SELECT p.user_id,max(p.player_name) AS name,COALESCE(sum(u.calculated_score),0) AS score FROM active_quizzes a JOIN user_played_quizzes p ON p.active_quiz_id=a.id AND p.is_host=false LEFT JOIN user_quiz_responses u ON u.user_played_quiz_id=p.id WHERE a.room_id=$1 AND a.is_active=false GROUP BY p.user_id) SELECT totals.name,totals.score,users.username,COALESCE(users.img_key,'') AS img_key,dense_rank() OVER(ORDER BY totals.score DESC) AS rank FROM totals JOIN users ON users.id=totals.user_id ORDER BY totals.score DESC,totals.name,totals.user_id`, id)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load standings")
	}
	rounds := []struct {
		ID          string `db:"id" json:"id"`
		Title       string `db:"title" json:"title"`
		Participant string `db:"participant" json:"participant"`
	}{}
	err = qc.db.ScanStructs(&rounds, `SELECT a.id::text,a.title,COALESCE((SELECT p.id::text FROM user_played_quizzes p WHERE p.active_quiz_id=a.id AND p.user_id=$2 AND p.is_host=false LIMIT 1),'') AS participant FROM active_quizzes a WHERE a.room_id=$1 AND a.is_active=false ORDER BY a.activated_from,a.id`, id, user.ID)
	if err != nil {
		return utils.JSONFail(c, 500, "Unable to load played quizzes")
	}
	return utils.JSONSuccess(c, 200, fiber.Map{"standings": standings, "quizzes": rounds, "ended": ended, "is_host": host, "username": user.Username})
}
