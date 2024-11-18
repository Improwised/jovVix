package structs

import (
	"github.com/google/uuid"
)

// All request sturcts
// Request struct have Req prefix

type ReqRegisterUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type ReqUpdateUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
}

type ReqLoginUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ReqAnswerSubmit struct {
	QuestionId   uuid.UUID `json:"id" validate:"required"`
	AnswerKeys   []int     `json:"keys" validate:"required"`
	ResponseTime int       `json:"response_time" validate:"required"`
}

type ReqUpdateQuestion struct {
	Question          string            `json:"question" validate:"required"`
	Type              int               `json:"type" validate:"required"`
	Options           map[string]string `json:"options" validate:"required"`
	Answers           []int             `json:"answers" validate:"required"`
	Points            int16             `json:"points"`
	DurationInSeconds int               `json:"duration_in_seconds" validate:"required"`
	QuestionMedia     string            `json:"question_media" validate:"required"`
	OptionsMedia      string            `json:"options_media" validate:"required"`
	Resource          string            `json:"resource"`
}

type ReqShareQuiz struct {
	Email      string `json:"email" validate:"required,email"`
	Permission string `json:"permission" validate:"required"`
}
