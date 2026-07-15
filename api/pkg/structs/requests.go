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

type ReqCreateQuiz struct {
	Title             string `json:"title" validate:"required"`
	Description       string `json:"description"`
	Points            int16  `json:"points"`
	DurationInSeconds int    `json:"duration_in_seconds"`
	IsPublic          bool   `json:"is_public"`
	CategoryId        string `json:"category_id" validate:"omitempty,uuid"`
	CoverImage        string `json:"cover_image"`
}

type ReqQuizCategory struct {
	Name string `json:"name" validate:"required,max=50"`
}

type ReqUpdateQuizSettings struct {
	Points            int16    `json:"points" validate:"min=0,max=20"`
	DurationInSeconds int      `json:"duration_in_seconds" validate:"required,min=1"`
	QuestionIds       []string `json:"question_ids" validate:"required,min=1,dive,uuid"`
}

type ReqCreateQuestion struct {
	Question          string            `json:"question" validate:"required"`
	Type              int               `json:"type" validate:"required"`
	Options           map[string]string `json:"options" validate:"required"`
	Answers           []int             `json:"answers" validate:"required"`
	Points            int16             `json:"points" validate:"omitempty,min=0,max=20"`
	DurationInSeconds int               `json:"duration_in_seconds" validate:"omitempty,min=1"`
	QuestionMedia     string            `json:"question_media" validate:"required"`
	OptionsMedia      string            `json:"options_media" validate:"required"`
	Resource          string            `json:"resource"`
}

type ReqShareQuiz struct {
	Email      string `json:"email" validate:"required,email"`
	Permission string `json:"permission" validate:"required"`
}
