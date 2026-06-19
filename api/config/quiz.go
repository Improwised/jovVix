package config

import "strings"

type QuizConfig struct {
	QuestionTimeLimit      string   `envconfig:"QUESTION_TIME_LIMIT"`
	DefaultQuestionPoints  int16    `envconfig:"DEFAULT_QUESTION_POINTS"`
	ScoreboardMaxDuration  string   `envconfig:"SCOREBOARD_MAX_DURATION"`
	FileSize               int64    `envconfig:"MAX_QUIZ_FILE_SIZE"`
	PublicQuizAdminEmails  []string `envconfig:"PUBLIC_QUIZ_ADMIN_EMAILS"`
	ActiveQuizTTLHours     int      `envconfig:"ACTIVE_QUIZ_TTL_HOURS"`
	ActiveQuizSweepMinutes int      `envconfig:"ACTIVE_QUIZ_SWEEP_MINUTES"`
}

// IsPublicQuizAdmin reports whether the given email is allowed to publish public quizzes.
// Comparison is case-insensitive and trims whitespace around each configured entry.
func (q QuizConfig) IsPublicQuizAdmin(email string) bool {
	if email == "" {
		return false
	}
	target := strings.ToLower(strings.TrimSpace(email))
	for _, admin := range q.PublicQuizAdminEmails {
		if strings.EqualFold(strings.TrimSpace(admin), target) {
			return true
		}
	}
	return false
}
