package structs

import "github.com/google/uuid"

// All request sturcts
// Request struct have Req prefix

type ReqRegisterUser struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	UserName  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
}

type ReqLoginUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ReqAnswerSubmit struct {
	QuestionId uuid.UUID `json:"id" validate:"required"`
	AnswerKeys []int     `json:"keys" validate:"required"`
}
