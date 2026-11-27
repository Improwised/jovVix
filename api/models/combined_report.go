package models

import (
	"sort"
	"time"

	"github.com/lib/pq"
)

// CombinedSession is one ended session of a quiz that the host can add to a combined report.
type CombinedSession struct {
	ID            string     `json:"id" db:"id"`
	ActivatedFrom time.Time  `json:"activated_from" db:"activated_from"`
	ActivatedTo   *time.Time `json:"activated_to" db:"activated_to"`
	Participants  int        `json:"participants" db:"participants"`
	Selected      bool       `json:"selected" db:"-"`
}

// CombinedPlay is one player's result in one session.
type CombinedPlay struct {
	PlayID       string `db:"play_id"`
	UserID       string `db:"user_id"`
	SessionID    string `db:"session_id"`
	Name         string `db:"name"`
	Username     string `db:"username"`
	ImgKey       string `db:"img_key"`
	Score        int    `db:"score"`
	Correct      int    `db:"correct"`
	Questions    int    `db:"questions"`
	ResponseTime int    `db:"response_time"`
}

// CombinedEntry is one row of the merged leaderboard.
type CombinedEntry struct {
	Rank         int    `json:"rank"`
	Name         string `json:"firstname"`
	Username     string `json:"username"`
	ImgKey       string `json:"img_key"`
	SessionID    string `json:"session_id"`
	Score        int    `json:"score"`
	Correct      int    `json:"correct"`
	Questions    int    `json:"questions"`
	ResponseTime int    `json:"response_time"`
	Plays        int    `json:"plays"`
}

// CombinedReport is the merged leaderboard across the selected sessions of one quiz.
type CombinedReport struct {
	QuizID       string            `json:"quiz_id"`
	Title        string            `json:"title"`
	Sessions     []CombinedSession `json:"sessions"`
	Participants int               `json:"participants"`
	Plays        int               `json:"plays"`
	Leaderboard  []CombinedEntry   `json:"leaderboard"`
}

// better reports whether play a outranks play b: higher score, then faster total response time.
func better(a, b CombinedPlay) bool {
	if a.Score != b.Score {
		return a.Score > b.Score
	}
	if a.ResponseTime != b.ResponseTime {
		return a.ResponseTime < b.ResponseTime
	}
	return a.PlayID < b.PlayID
}

// CombineLeaderboard keeps each player's best play and ranks all players together.
// A player who joined more than one selected session appears once. Rank is a dense
// rank on score, the same rule as the per-session scoreboard; equal scores are listed
// fastest first.
func CombineLeaderboard(plays []CombinedPlay) []CombinedEntry {
	best := map[string]CombinedPlay{}
	count := map[string]int{}
	for _, p := range plays {
		count[p.UserID]++
		if cur, ok := best[p.UserID]; !ok || better(p, cur) {
			best[p.UserID] = p
		}
	}
	kept := make([]CombinedPlay, 0, len(best))
	for _, p := range best {
		kept = append(kept, p)
	}
	sort.Slice(kept, func(i, j int) bool { return better(kept[i], kept[j]) })

	entries := make([]CombinedEntry, len(kept))
	rank := 0
	for i, p := range kept {
		if i == 0 || p.Score != kept[i-1].Score {
			rank++
		}
		entries[i] = CombinedEntry{
			Rank: rank, Name: p.Name, Username: p.Username, ImgKey: p.ImgKey, SessionID: p.SessionID,
			Score: p.Score, Correct: p.Correct, Questions: p.Questions, ResponseTime: p.ResponseTime,
			Plays: count[p.UserID],
		}
	}
	return entries
}

// ErrNoSessions means the host has no ended session with players for the quiz.
type ErrNoSessions struct{}

func (ErrNoSessions) Error() string { return "no ended sessions with players for this quiz" }

// ErrUnknownSession means a requested session is not an ended session of this quiz hosted by the user.
type ErrUnknownSession struct{ ID string }

func (e ErrUnknownSession) Error() string {
	return "session " + e.ID + " is not an ended session of this quiz"
}

// GetCombinedReport merges the selected sessions of quizID hosted by hostID.
// An empty selection selects every session.
func (model *QuizModel) GetCombinedReport(quizID, hostID string, selected []string) (CombinedReport, error) {
	report := CombinedReport{QuizID: quizID, Sessions: []CombinedSession{}, Leaderboard: []CombinedEntry{}}

	// Only the user's own ended sessions are listed, so a host can never select
	// another host's session. Sessions count when a player answered, as on the reports list.
	err := model.db.ScanStructs(&report.Sessions, `SELECT a.id::text AS id, a.activated_from, a.activated_to,
 count(DISTINCT p.id)::int AS participants
 FROM active_quizzes a
 JOIN user_played_quizzes p ON p.active_quiz_id = a.id AND p.is_host = false
 JOIN user_quiz_responses r ON r.user_played_quiz_id = p.id
 WHERE a.quiz_id = $1 AND a.admin_id = $2 AND a.activated_to IS NOT NULL
 GROUP BY a.id ORDER BY a.activated_from, a.id`, quizID, hostID)
	if err != nil {
		return report, err
	}
	if len(report.Sessions) == 0 {
		return report, ErrNoSessions{}
	}
	if _, err = model.db.ScanVal(&report.Title, `SELECT title FROM quizzes WHERE id = $1`, quizID); err != nil {
		return report, err
	}

	known := map[string]bool{}
	for _, s := range report.Sessions {
		known[s.ID] = true
	}
	want := map[string]bool{}
	for _, id := range selected {
		if !known[id] {
			return report, ErrUnknownSession{ID: id}
		}
		want[id] = true
	}
	ids := []string{}
	for i := range report.Sessions {
		s := &report.Sessions[i]
		s.Selected = len(selected) == 0 || want[s.ID]
		if s.Selected {
			ids = append(ids, s.ID)
		}
	}

	plays := []CombinedPlay{}
	err = model.db.ScanStructs(&plays, `SELECT p.id::text AS play_id, p.user_id, a.id::text AS session_id,
 COALESCE(p.player_name, u.first_name, '') AS name, u.username, COALESCE(u.img_key, '') AS img_key,
 COALESCE(sum(r.calculated_score), 0)::int AS score,
 count(r.id) FILTER (WHERE r.calculated_score > 0)::int AS correct,
 (SELECT count(*) FROM active_quiz_questions q WHERE q.active_quiz_id = a.id)::int AS questions,
 COALESCE(sum(r.response_time) FILTER (WHERE r.response_time > 0), 0)::int AS response_time
 FROM active_quizzes a
 JOIN user_played_quizzes p ON p.active_quiz_id = a.id AND p.is_host = false
 JOIN users u ON u.id = p.user_id
 JOIN user_quiz_responses r ON r.user_played_quiz_id = p.id
 WHERE a.id = ANY($1::uuid[])
 GROUP BY p.id, a.id, u.id`, pq.Array(ids))
	if err != nil {
		return report, err
	}
	report.Plays = len(plays)
	report.Leaderboard = CombineLeaderboard(plays)
	report.Participants = len(report.Leaderboard)
	return report, nil
}
