package response

import "github.com/Improwised/jovvix/api/models"

type ResponseFinalScore struct {
	FinalScore []models.FinalScoreBoard
}

type ResponseFinalScoreForAdmin struct {
	FinalScore []models.FinalScoreBoardAdmin
}
