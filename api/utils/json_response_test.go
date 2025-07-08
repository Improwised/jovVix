package utils_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/Improwised/jovvix/api/utils"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONSuccess(t *testing.T) {
	app := fiber.New()

	app.Get("/test-success", func(c *fiber.Ctx) error {
		data := map[string]string{"message": "Success"}
		return utils.JSONSuccess(c, fiber.StatusOK, data)
	})

	req := httptest.NewRequest("GET", "/test-success", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	assert.Equal(t, "success", result["status"])
	assert.NotNil(t, result["data"])

	data := result["data"].(map[string]interface{})
	assert.Equal(t, "Success", data["message"])
}

func TestJSONFail(t *testing.T) {
	app := fiber.New()

	app.Get("/test-fail", func(c *fiber.Ctx) error {
		data := map[string]string{"error": "Invalid input"}
		return utils.JSONFail(c, fiber.StatusBadRequest, data)
	})

	req := httptest.NewRequest("GET", "/test-fail", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Decode response into a generic map to inspect the content
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Check if the response is structured as expected
	assert.Equal(t, "fail", result["status"])
	assert.NotNil(t, result["data"])

	data := result["data"].(map[string]interface{})
	assert.Equal(t, "Invalid input", data["error"])
}

func TestJSONError(t *testing.T) {
	app := fiber.New()

	app.Get("/test-error", func(c *fiber.Ctx) error {
		return utils.JSONError(c, fiber.StatusInternalServerError, "Internal Server Error")
	})

	req := httptest.NewRequest("GET", "/test-error", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Decode response into a generic map to inspect the content
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	// Check if the response is structured as expected
	assert.Equal(t, "error", result["status"])
	assert.Equal(t, "Internal Server Error", result["message"])
}
