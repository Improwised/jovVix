package config

type QuizConfig struct {
	QuestionTimeLimit     string `envconfig:"QUESTION_TIME_LIMIT"`
	DefaultQuestionPoints int16  `envconfig:"DEFAULT_QUESTION_POINTS"`
	ScoreboardMaxDuration string `envconfig:"SCOREBOARD_MAX_DURATION"`
	FileSize              int64  `envconfig:"MAX_QUIZ_FILE_SIZE"`
}
