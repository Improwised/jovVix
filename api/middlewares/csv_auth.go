package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (m *Middleware) ValidateCsv(c *fiber.Ctx) error {
	userID := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))
	file, err := c.FormFile("attachment")

	if err != nil {
		m.Logger.Error("error in getting csv file", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrGettingAttachment)
	}

	// size validation
	if file.Size > constants.FileSize {
		m.Logger.Error("error in getting csv file", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrGettingAttachment)
	}

	isMatched := false
	allowedTypes := []string{
		"text/csv",
	}

	for _, types := range allowedTypes {

		if types == file.Header.Get("Content-Type") {
			isMatched = true
			break
		}
	}

	// content type validation
	if !isMatched {
		m.Logger.Error("file type mismatch", zap.Any("file", file))
		return utils.JSONFail(c, fiber.StatusBadRequest, constants.ErrFileIsNotInSupportedType)
	}

	folder := "./uploads"

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		// Create the folder and its parent directories
		if err := os.MkdirAll(folder, 0755); err != nil {
			m.Logger.Error("folder creation failed")
			return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrProblemInUploadFile)
		}
		m.Logger.Debug("folder creation success")
	}

	destination := folder + "/" + strings.TrimSpace(userID) + "_" + file.Filename

	if err := c.SaveFile(file, destination); err != nil {
		m.Logger.Error("error in storing file", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.ErrProblemInUploadFile)
	}

	c.Locals(constants.FileName, destination)

	return c.Next()
}
