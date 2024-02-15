package components

type QuizGameManager struct {
	quizSession map[string]QuizSessionCfg
}

func InitQuizGameManager() *QuizGameManager {
	return &QuizGameManager{}
}

func (qgm *QuizGameManager) Join(user *User, code *string) {
}
