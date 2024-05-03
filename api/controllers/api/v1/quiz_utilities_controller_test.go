package v1_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"bytes"
// 	"fmt"
// 	"mime/multipart"
// 	"os"
// 	"testing"

// 	"github.com/Improwised/quizz-app/api/pkg/structs"
// 	"github.com/stretchr/testify/assert"
// )

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

		// attachmentPath := ".././././app/public/files/demo.csv"

		// // Create a multipart form (using bytes.Buffer for efficiency)
		// var formBuffer bytes.Buffer
		// multipartWriter := multipart.NewWriter(&formBuffer)

		// // Add description field
		// description := "Demo description"
		// descriptionWriter, err := multipartWriter.CreateFormField("description")
		// if err != nil {
		// 	fmt.Println("Error creating description field:", err)
		// 	return
		// }
		// _, err = descriptionWriter.Write([]byte(description))
		// if err != nil {
		// 	fmt.Println("Error writing description:", err)
		// 	return
		// }

		// // Add attachment field
		// attachmentWriter, err := multipartWriter.CreateFormFile("attachment", attachmentPath)
		// if err != nil {
		// 	fmt.Println("Error creating attachment field:", err)
		// 	return
		// }

		// // Open the attachment file
		// attachmentBytes, err := os.ReadFile(attachmentPath)
		// if err != nil {
		// 	fmt.Println("Error reading attachment file:", err)
		// 	return
		// }

		// // Write attachment content
		// _, err = attachmentWriter.Write(attachmentBytes)
		// if err != nil {
		// 	fmt.Println("Error writing attachment:", err)
		// 	return
		// }

		// // Close the multipart writer
		// err = multipartWriter.Close()
		// if err != nil {
		// 	fmt.Println("Error closing multipart writer:", err)
		// 	return
		// }

		// fmt.Println("formbuffer/**********************", formBuffer)
		// res, err = client.R().
		// 	SetHeader("Content-Type", "application/json").
		// 	SetHeader("Content-Type", multipartWriter.FormDataContentType()).
		// 	SetHeader("Content-Type", "text/csv").
		// 	SetBody(formBuffer.Bytes()).
		// 	SetCookies(cookies).
		// 	Post("/api/v1/admin/quizzes/demo/upload")

		// fmt.Println("response******************", res)
		// if err != nil {
		// 	t.Errorf("Error sending request: %v", err)
		// }

		// --------------------------------------------------------------------------------------------------------------------------
		url := "http://localhost:3500/api/v1/admin/quizzes/title/upload"
		method := "POST"
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		file, err := os.Open(".././././app/public/files/demo.csv")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		part1, errFile1 := writer.CreateFormFile("attachment", filepath.Base(".././././app/public/files/demo.csv"))
		if errFile1 != nil {
			fmt.Println(errFile1)
			return
		}

		_, errFile1 = io.Copy(part1, file)

		if errFile1 != nil {
			fmt.Println(errFile1)
			return
		}
		_ = writer.WriteField("description", "this is description")
		err = writer.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Cookie", cookies[0].Name+"="+cookies[0].Value)

		req.Header.Set("Content-Type", "text/csv")
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))

		// res, err := client.
		// 	R().
		// 	EnableTrace().
		// 	SetBody(payload).
		// 	SetCookies(cookies).
		// 	SetHeader("Content-Type", writer.FormDataContentType()).
		// 	Post("/api/v1/admin/quizzes/demo/upload")

		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// fmt.Println("response *************", res)

	})

}
