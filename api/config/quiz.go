package config

type QuizConfig struct {
	QuestionTimeLimit     string `envconfig:"QUESTION_TIME_LIMIT"`
	ScoreboardMaxDuration string `envconfig:"SCOREBOARD_MAX_DURATION"`
	FileSize              int64  `envconfig:"QUIZ_FILE_SIZE" default:"1048576"`
}
