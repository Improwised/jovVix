package config

type QuizConfig struct {
	QuestionTimeLimit     string `envconfig:"QUESTION_TIME_LIMIT"`
	ScoreboardMaxDuration string `envconfig:"SCOREBOARD_MAX_DURATION"`
}
