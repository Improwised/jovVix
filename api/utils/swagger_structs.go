package utils

import (
	"mime/multipart"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/structs"
)

// swagger:parameters RequestAnalyticsBoardForAdmin
type RequestAnalyticsBoardForAdmin struct {
	// in:query
	// required: true
	ActiveQuizId string `json:"active_quiz_id"`
}

// swagger:response ResponseAnalyticsBoardForAdmin
type ResponseAnalyticsBoardForAdmin struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.AnalyticsBoardAdmin
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestAnalyticsBoardForUser
type RequestAnalyticsBoardForUser struct {
	// in:query
	// required: true
	UserPlayedQuiz string `json:"user_played_quiz"`
}

// swagger:response ResponseAnalyticsBoardForUser
type ResponseAnalyticsBoardForUser struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.AnalyticsBoardUser
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestCreateQuickUser
type RequestCreateQuickUser struct {
	// in:path
	Username string `json:"username"`
	// in:query
	// required: true
	Avatar string `json:"avatar_name"`
}

// swagger:response ResponseUserDetails
type ResponseUserDetails struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			models.User
		} `json:"data"`
	} `json:"body"`
}

// swagger:response ResponseGetRegisteredUser
type ResponseGetRegisteredUser struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			config.KratosUserDetails
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestUpadateRegisteredUser
type RequestUpadateRegisteredUser struct {
	// in:body
	// required: true
	Body struct {
		structs.ReqUpdateUser
	}
}

// swagger:parameters RequestFinalScoreForAdmin
type RequestFinalScoreForAdmin struct {
	// in:query
	// required: true
	ActiveQuizId string `json:"active_quiz_id"`
}

// swagger:response ResponseFinalScoreForAdmin
type ResponseFinalScoreForAdmin struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.FinalScoreBoardAdmin
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestFinalScoreForUser
type RequestFinalScoreForUser struct {
	// in:query
	// required: true
	UserPlayedQuiz string `json:"user_played_quiz"`
}

// swagger:response ResponseFinalScoreForUser
type ResponseFinalScoreForUser struct {
	//in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.FinalScoreBoard
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestInsertImage
type RequestInsertImage struct {
	// in:query
	// required:true
	QuizId string `json:"quiz_id"`

	// in: formData
	// required: true
	// type: file
	// swagger:file
	// name: image
	File *multipart.FileHeader `json:"image-attachment"`
}

// swagger:response ResponseInsertImage
type ResponseInsertImage struct {
	// in:body
	Body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}

// swagger:response ResponseAdminUploadedQuiz
type ResponseAdminUploadedQuiz struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.QuizWithQuestions
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestGetQuizAnalysis
type RequestGetQuizAnalysis struct {
	// in:path
	ActiveQuizId string `json:"active_quiz_id"`
}

// swagger:response ResponseGetQuizAnalysis
type ResponseGetQuizAnalysis struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.QuizAnalysis
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestListQuizzesAnalysis
type RequestListQuizzesAnalysis struct {
	// in:query
	// required: true
	OrderBy string `json:"orderBy"`
}

// swagger:response RsponseListQuizzesAnalysis
type RsponseListQuizzesAnalysis struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			Data  []models.QuizzesAnalysis `json:"data"`
			Count int64                    `json:"count"`
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestQuizCreated
type RequestQuizCreated struct {
	// in:path
	// required:true
	// description: Title of the quiz
	QuizTitle string `json:"quiz_title"`

	// in: formData
	// required: true
	// description: The CSV file containing quiz questions
	// type: file
	// swagger:file
	// name: attachment
	File *multipart.FileHeader `json:"attachment"`

	// in: formData
	// required: false
	// description: A description of the quiz
	// required: true
	Description string `json:"description"`
}

// swagger:response ResponseQuizCreated
type ResponseQuizCreated struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		QuizId string `json:"quizId"`
	} `json:"body"`
}

// swagger:parameters RequestGenerateDemoSession
type RequestGenerateDemoSession struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
}

// swagger:response ResponseGenerateDemoSession
type ResponseGenerateDemoSession struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		QuizId string `json:"quizId"`
	} `json:"body"`
}

// swagger:parameters RequestListQuestionByQuizId
type RequestListQuestionByQuizId struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
}

// swagger:response ResponseListQuestionByQuizId
type ResponseListQuestionByQuizId struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			structs.ResQuestionAnalytics
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestListUserPlayedQuizes
type RequestListUserPlayedQuizes struct {
	// in:query
	Page string `json:"page"`
	// in:query
	Title string `json:"title"`
}

// swagger:response ResponseListUserPlayedQuizes
type ResponseListUserPlayedQuizes struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			structs.ResUserPlayedQuizWithCount
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestListUserPlayedQuizesWithQuestionById
type RequestListUserPlayedQuizesWithQuestionById struct {
	// in:path
	// required: true
	UserPlayedQuizId string `json:"user_played_quiz_id"`
}

// swagger:response ResponseListUserPlayedQuizesWithQuestionById
type ResponseListUserPlayedQuizesWithQuestionById struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			structs.ResUserPlayedQuizAnalyticsBoard
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestPlayedQuizValidation
type RequestPlayedQuizValidation struct {
	// in:path
	// required: true
	InvitationCode string `json:"invitationCode"`
}

// swagger:response ResponsePlayedQuizValidation
type ResponsePlayedQuizValidation struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			UserPlayedQuizId string `json:"user_played_quiz"`
			SessionId        string `json:"session_id"`
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestGetQuestionById
type RequestGetQuestionById struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
	// in:path
	// required: true
	QuestionId string `json:"question_id"`
}

// swagger:response ResponseGetQuestionById
type ResponseGetQuestionById struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   struct {
			structs.QuestionAnalytics
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestUpdateQuestionById
type RequestUpdateQuestionById struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
	// in:path
	// required: true
	QuestionId string `json:"question_id"`
}

// swagger:parameters RequestDeleteQuestionById
type RequestDeleteQuestionById struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
	// in:path
	// required: true
	QuestionId string `json:"question_id"`
}

// swagger:parameters RequestShareQuiz
type RequestShareQuiz struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
	// in:body
	// required: true
	Body struct {
		structs.ReqShareQuiz
	}
}

// swagger:parameters RequestListQuizAuthorizedUsers
type RequestListQuizAuthorizedUsers struct {
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
}

// swagger:response ResponseListQuizAuthorizedUsers
type ResponseListQuizAuthorizedUsers struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			structs.ResUserWithQuizPermission
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestListSharedQuizzes
type RequestListSharedQuizzes struct {
	// in:query
	// required: true
	Type string `json:"type"`
}

// swagger:response ResponseListSharedQuizzes
type ResponseListSharedQuizzes struct {
	// in:body
	Body struct {
		Status string `json:"status"`
		Data   []struct {
			models.QuizWithQuestions
		} `json:"data"`
	} `json:"body"`
}

// swagger:parameters RequestUpdateUserPermissionOfQuiz
type RequestUpdateUserPermissionOfQuiz struct {
	// in:query
	// required: true
	SharedQuizId string `json:"shared_quiz_id"`
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
	// in:body
	// required: true
	Body struct {
		structs.ReqShareQuiz
	}
}

// swagger:parameters RequestDeleteUserPermissionOfQuiz
type RequestDeleteUserPermissionOfQuiz struct {
	// in:query
	// required: true
	SharedQuizId string `json:"shared_quiz_id"`
	// in:path
	// required: true
	QuizId string `json:"quiz_id"`
}

////////////////////
// --- GENERIC ---//
////////////////////

// swagger:response ResponseOkWithMessage
type ResponseOkWithMessage struct {
	// in:body
	Body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}

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
