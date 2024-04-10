package utils

import (
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/response"
	"github.com/Improwised/quizz-app/api/pkg/structs"
)

// swagger:parameters RequestCreateUser
type RequestCreateUser struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqRegisterUser
	}
}

// swagger:response ResponseCreateUser
type ResponseCreateUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestGetUser
type RequestGetUser struct {
	// in:path
	UserId string `json:"userId"`
}

// swagger:response ResponseGetUser
type ResponseGetUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestAuthnUser
type RequestAuthnUser struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqLoginUser
	}
}

// swagger:response ResponseAuthnUser
type ResponseAuthnUser struct {
	// in:body
	Body struct {
		// enum: success
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestFinalScore
type RequestFinalScore struct {
	// in:query
	UserPlayedQuiz string `json:"userPlayedQuiz"`
}

// swagger:response ResponseFinalScore
type ResponseFinalScore struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			response.ResponseFinalScore
		} `json:"data"`
	} `json:"body"`
}

////////////////////
// --- GENERIC ---//
////////////////////

// Response is okay
// swagger:response GenericResOk
type ResOK struct {
	// in:body
	Body struct {
		// enum:success
		Status string `json:"status"`
	}
}

// Fail due to user invalid input
// swagger:response GenericResFailBadRequest
type ResFailBadRequest struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to user invalid input
// swagger:response ResForbiddenRequest
type ResForbiddenRequest struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Server understand request but refuse to authorize it
// swagger:response GenericResFailConflict
type ResFailConflict struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to server understand request but unable to process
// swagger:response GenericResFailUnprocessableEntity
type ResFailUnprocessableEntity struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Fail due to resource not exists
// swagger:response GenericResFailNotFound
type ResFailNotFound struct {
	// in: body
	Body struct {
		// enum: fail
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	} `json:"body"`
}

// Unexpected error occurred
// swagger:response GenericResError
type ResError struct {
	// in: body
	Body struct {
		// enum: error
		Status  string      `json:"status"`
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
	} `json:"body"`
}
