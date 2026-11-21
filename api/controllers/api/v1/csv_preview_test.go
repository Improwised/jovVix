package v1

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// The import modal reads fields, headers, mapping, unmapped_required and rows from this response.
func TestPreviewQuestionsCsvShape(t *testing.T) {
	upload := func(t *testing.T, csv, mapping string) (int, map[string]any, string) {
		path := filepath.Join(t.TempDir(), "upload.csv")
		if err := os.WriteFile(path, []byte(csv), 0o600); err != nil {
			t.Fatal(err)
		}
		app := fiber.New()
		ctrl := &QuestionController{logger: zap.NewNop()}
		app.Post("/preview", func(c *fiber.Ctx) error {
			c.Locals(constants.FileName, path)
			return c.Next()
		}, ctrl.PreviewQuestionsCsv)

		body := "mapping=" + mapping
		req := httptest.NewRequest("POST", "/preview", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}
		raw, _ := io.ReadAll(res.Body)
		var parsed struct {
			Data map[string]any `json:"data"`
		}
		_ = json.Unmarshal(raw, &parsed)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Fatal("preview must delete the uploaded file")
		}
		return res.StatusCode, parsed.Data, string(raw)
	}

	csv := "Prompt,Kind,A,B,Right\nQ1,single answer,x,y,1\nQ2,single answer,x,y,9\n"

	status, data, raw := upload(t, csv, "")
	if status != 200 || len(data["unmapped_required"].([]any)) != 5 || len(data["rows"].([]any)) != 0 {
		t.Fatalf("unmapped upload: %d %s", status, raw)
	}

	status, data, raw = upload(t, csv, `{"Question Text":"Prompt","Question Type":"Kind","Option 1":"A","Option 2":"B","Correct Answer":"Right"}`)
	rows, _ := data["rows"].([]any)
	if status != 200 || len(rows) != 2 || len(data["fields"].([]any)) != 12 {
		t.Fatalf("mapped upload: %d %s", status, raw)
	}
	second := rows[1].(map[string]any)
	if second["row"].(float64) != 3 || len(second["errors"].([]any)) == 0 || second["values"].(map[string]any)["Question Text"] != "Q2" {
		t.Fatalf("second row: %v", second)
	}

	status, _, raw = upload(t, csv, `{"Question Text":"Missing"}`)
	if status != 400 || !strings.Contains(raw, constants.ErrInvalidCSVMapping) {
		t.Fatalf("bad mapping: %d %s", status, raw)
	}
}
