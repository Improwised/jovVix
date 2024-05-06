package v1_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"

	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestCreateQuizByCSV(t *testing.T) {

	req := structs.ReqLoginUser{
		Email:    "adminxyz@gmail.com",
		Password: "RZo5(uXD<3#aH0",
	}

	res, err := client.
		R().
		EnableTrace().
		SetBody(req).
		Post("/api/v1/login")

	assert.Nil(t, err)
	cookies := res.Cookies()
	assert.NotEqual(t, len(cookies), 0)

	t.Run("Upload CSV file for question", func(t *testing.T) {
		url := "http://127.0.0.1:3500/api/v1/admin/quizzes/title/upload"
		file, err := os.Open(".././././app/public/files/demo.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)

		part, errFile1 := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition": []string{`form-data; name="attachment"; filename=".././././app/public/files/demo.csv"`},
			"Content-Type":        []string{"text/csv"},
		})
		if errFile1 != nil {
			fmt.Println(errFile1)
			return
		}

		_, errFile1 = io.Copy(part, file)

		if errFile1 != nil {
			fmt.Println(errFile1)
			return
		}
		_ = writer.WriteField("description", "test description")
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = client.R().
			SetBody(payload).
			SetHeader("Content-Type", writer.FormDataContentType()).
			SetCookie(cookies[0]).
			Post(url)

		if err != nil {
			fmt.Println(err)
			return
		}
	})
}
