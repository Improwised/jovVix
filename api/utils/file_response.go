package utils

import "github.com/gofiber/fiber/v2"

func CsvFileResponse(c *fiber.Ctx, filePath string, fileName string) error {
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename="+fileName)
	return c.SendFile(filePath)
}
