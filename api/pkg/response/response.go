package response

import "github.com/Improwised/quizz-app/api/models"

type ResponseFinalScore struct {
	FinalScore []models.FinalScoreBoard
}

type ResponseFinalScoreForAdmin struct {
	FinalScore []models.FinalScoreBoardAdmin
}
