package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
)

type CreatePdfService struct {
	pdf       pdf.Maroto
	appConfig *config.AppConfig
}

func NewCreatePdfService(appConfig *config.AppConfig) *CreatePdfService {
	p := pdf.NewMaroto(consts.Portrait, consts.A4)
	return &CreatePdfService{
		pdf:       p,
		appConfig: appConfig,
	}
}

func (cps *CreatePdfService) CreatPdf(orderToUserAndQuestionData map[int][]structs.UserAndQuestionData, orderOfQuestion []int, activeQuizId string) ([]byte, error) {
	cps.pdf.SetPageMargins(20, 10, 20)
	utils.BuildHeading(cps.pdf)

	utils.BuildUsersTables(cps.pdf, orderToUserAndQuestionData, orderOfQuestion)
	pdfPath := filepath.Join(cps.appConfig.PDFS_FILE_PATH, fmt.Sprintf("/%s.pdf", activeQuizId))
	err := cps.pdf.OutputFileAndClose(pdfPath)

	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(pdfPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
