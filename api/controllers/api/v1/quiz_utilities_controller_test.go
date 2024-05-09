package v1_test

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"testing"

	"github.com/Improwised/quizz-app/api/config"
	v1 "github.com/Improwised/quizz-app/api/controllers/api/v1"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var result struct {
	Status string `json:"status"`
	Data   string `json:"data"`
}

func TestCreateQuizByCSV(t *testing.T) {
	filepath := ".././././app/public/files/demo.csv"

	// login admin
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
		file, err := os.Open(filepath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)

		part, errFile1 := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition": []string{`form-data; name="attachment"; filename=filepath`},
			"Content-Type":        []string{"text/csv"},
		})
		if errFile1 != nil {
			return
		}

		_, errFile1 = io.Copy(part, file)

		if errFile1 != nil {
			return
		}
		err = writer.WriteField("description", "test description")
		if err != nil {
			return
		}
		err = writer.Close()
		if err != nil {
			return
		}

		title := "Demo Quiz"

		res, err = client.R().
			SetBody(payload).
			SetHeader("Content-Type", writer.FormDataContentType()).
			SetCookie(cookies[0]).
			SetResult(&result).
			SetPathParam("title", title).
			Post("api/v1/admin/quizzes/{title}/upload")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode())
		assert.NotNil(t, result.Data)
	})

	t.Run("check if csv is validate or not", func(t *testing.T) {
		err := v1.ValidateCSVFileFormat(filepath)

		assert.Nil(t, err)
	})

	t.Run("extract questions from csv", func(t *testing.T) {
		questions, err := v1.ExtractQuestionsFromCSV(filepath)
		assert.NotEmpty(t, questions)
		assert.Nil(t, err)
	})
}

func TestGenerateDemoSession(t *testing.T) {

	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)
	req := structs.ReqLoginUser{
		Email:    "adminxyz@gmail.com",
		Password: "RZo5(uXD<3#aH0",
	}

	// login admin to get cookie
	res, err := client.
		R().
		EnableTrace().
		SetBody(req).
		Post("/api/v1/login")

	assert.Nil(t, err)
	cookies := res.Cookies()
	assert.NotEqual(t, len(cookies), 0)

	// get quizId
	var quizId *uuid.UUID

	subquery := db.Select("id").From("users").Limit(1)
	ok, err := db.Select("id").From("quizzes").Where(goqu.C("creator_id").In(subquery)).ScanVal(&quizId)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)

	t.Run("Generate demo session", func(t *testing.T) {
		res, err = client.R().
			SetCookie(cookies[0]).
			SetResult(&result).
			SetPathParam("quizId", quizId.String()).
			Post("/api/v1/admin/quizzes/{quizId}/demo_session")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode())
		assert.NotEmpty(t, result.Data)
	})

	t.Cleanup(func() {
		_, err = db.Exec("delete from active_quiz_questions")
		assert.Nil(t, err)

		_, err = db.Exec("delete from active_quizzes")
		assert.Nil(t, err)

		_, err = db.Exec("delete from quiz_questions")
		assert.Nil(t, err)

		_, err = db.Exec("delete from questions")
		assert.Nil(t, err)

		_, err = db.Exec("delete from quizzes")
		assert.Nil(t, err)

		_, err = db.Exec("delete from users")
		assert.Nil(t, err)
	})

}
