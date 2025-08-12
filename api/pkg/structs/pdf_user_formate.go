package structs

type UserAndQuestionData struct {
	Options        map[string]string
	CorrectAnswer  string
	UserName       string
	SelectedAnswer string
	Question       string
	QuestionType   string
}
