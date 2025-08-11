package utils

import (
	"fmt"
	"strings"

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
			m.Text("Quiz Report", props.Text{
				Top:   3,
				Style: consts.Bold,
				Align: consts.Center,
				Color: getBlackColor(),
			})
		})
	})
}

func BuildUsersTables(m pdf.Maroto, questionsAndUsersMap map[int][]structs.Users, questions []string, keys []int) {
	for ind, orderNo := range keys {
		question := ""
		if orderNo-1 >= 0 && orderNo-1 < len(questions) {
			question = questions[orderNo-1]
		}

		m.Row(8, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("Question : %d  %v", ind+1, question), props.Text{
					Style: consts.Bold,
					Size:  10,
					Color: getBlackColor(),
				})
			})
		})
		m.Line(1)

		for _, userData := range questionsAndUsersMap[orderNo] {
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

					// Correct option
					if optKey == strings.Split(userData.CorrectAnswer, "")[1] {
						txtColor = color.Color{Green: 255}
					} else if optKey == strings.Split(userData.SelectedAnswer, "")[1] {
						txtColor = color.Color{Red: 255}
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
