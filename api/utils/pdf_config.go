package utils

import (
	"fmt"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

func getBlackColor() color.Color {
	return color.Color{
		Red:   0,
		Green: 0,
		Blue:  0,
	}
}

func BuildHeading(m pdf.Maroto) {
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(constants.PdfTitleText, props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Color: getBlackColor(),
			})
		})
	})
}

func BuildUsersTables(m pdf.Maroto, orderToUserAndQuestionData map[int][]structs.UserAndQuestionData, orderOfQuestion []int) {
	for ind, orderNo := range orderOfQuestion {
		m.Row(8, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("%s : %d  %v", constants.PdfQueText, ind+1, orderToUserAndQuestionData[orderNo][0].Question), props.Text{
					Style: consts.Bold,
					Size:  10,
					Color: getBlackColor(),
				})
			})
		})
		m.Line(1)

		for _, userData := range orderToUserAndQuestionData[orderNo] {
			m.Row(6, func() {
				// User column
				m.Col(2, func() {
					m.Text(userData.UserName, props.Text{Size: 8, Color: getBlackColor()})
				})

				// Option columns
				for optNo := 1; optNo <= len(userData.Options); optNo++ {
					optKey := fmt.Sprintf("%d", optNo)

					bg := color.Color{Red: 255, Green: 255, Blue: 255} // default white
					txtColor := getBlackColor()
					correctAns := strings.Trim(userData.CorrectAnswer, "[]")
					selectedAns := strings.Trim(userData.SelectedAnswer, "[]")
					if userData.QuestionType == constants.SurveyString && len(selectedAns) > 0 {
						txtColor = color.Color{Green: 255}
					} else if len(selectedAns) == 0 {
						txtColor = getBlackColor()
					} else {
						if optKey == correctAns {
							txtColor = color.Color{Green: 255}
						}
						if optKey == selectedAns && selectedAns != correctAns {
							txtColor = color.Color{Red: 255}
						}
					}

					m.Col(2, func() {
						m.SetBackgroundColor(bg)
						m.Text(userData.Options[optKey], props.Text{Size: 8, Color: txtColor, Align: consts.Center})
						m.SetBackgroundColor(color.Color{Red: 255, Green: 255, Blue: 255})
					})
				}
			})
			m.Line(1)
		}
		m.Row(6, func() {})
	}
}
