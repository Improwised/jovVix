package structs

import "database/sql"

// All response structs
// Response struct have Res prefix

type SocketResponseFormat struct {
	EventName string      `json:"event"`
	Data      interface{} `json:"data"`
}

type ResUserPlayedQuiz struct {
	Id             string         `json:"id" db:"id"`
	Title          string         `json:"title" db:"title"`
	Description    sql.NullString `json:"description" db:"description"`
	CreatedAt      string         `json:"created_at" db:"created_at"`
	TotalQuestions string         `json:"total_questions" db:"total_questions"`
}

type ResUserPlayedQuizAnalyticsBoard struct {
	SelectedAnswer   sql.NullString    `db:"selected_answer,omitempty" json:"selected_answer"`
	CorrectAnswer    string            `db:"correct_answer,omitempty" json:"correct_answer"`
	CalculatedScore  int               `db:"calculated_score,omitempty" json:"calculated_score"`
	IsAttend         bool              `db:"is_attend,omitempty" json:"is_attend"`
	ResponseTime     int               `db:"response_time,omitempty" json:"response_time"`
	CalculatedPoints int               `db:"calculated_points,omitempty" json:"calculated_points"`
	Question         string            `db:"question,omitempty" json:"question"`
	RawOptions       []byte            `db:"options,omitempty" json:"raw_options"`
	Options          map[string]string `db:"omitempty" json:"options"`
	QuestionsMedia   string            `db:"question_media" json:"question_media"`
	OptionsMedia     string            `db:"options_media" json:"options_media"`
	Resource         string            `db:"resource" json:"resource"`
	Points           int               `db:"points,omitempty" json:"points"`
	QuestionTypeID   int               `db:"type,omitempty" json:"question_type_id"`
	QuestionType     string            `db:"omitempty" json:"question_type"`
}
