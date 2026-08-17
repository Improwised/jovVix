package services

import (
	"net/http"
	"strconv"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/stretchr/testify/assert"
)

func testCompletionRequest() chatCompletionRequest {
	temperature := 0.4
	return chatCompletionRequest{
		Model:          "some-model",
		Temperature:    &temperature,
		MaxTokens:      1500,
		ResponseFormat: &chatResponseFormat{Type: "json_object"},
	}
}

func TestAdaptAIRequestDropsJSONMode(t *testing.T) {
	bodies := []string{
		`{"error":{"message":"response_format is not supported"}}`,
		`{"error":{"message":"Invalid value: 'json_object'"}}`,
		`{"error":{"message":"This model does not support JSON mode"}}`,
		`{"error":{"message":"json_schema is unavailable for this model"}}`,
	}

	for _, providerBody := range bodies {
		request := testCompletionRequest()

		adaptation := adaptAIRequest(&request, providerBody)

		assert.Equal(t, constants.AIAdaptationDroppedJSONMode, adaptation, providerBody)
		assert.Nil(t, request.ResponseFormat)
		assert.Equal(t, 1500, request.MaxTokens)
		assert.NotNil(t, request.Temperature)
	}
}

func TestAdaptAIRequestSwitchesToMaxCompletionTokens(t *testing.T) {
	request := testCompletionRequest()
	request.ResponseFormat = nil

	adaptation := adaptAIRequest(&request,
		`{"error":{"message":"Unsupported parameter: 'max_tokens' is not supported with this model. Use 'max_completion_tokens' instead."}}`)

	assert.Equal(t, constants.AIAdaptationMaxCompletionTokens, adaptation)
	assert.Equal(t, 0, request.MaxTokens)
	assert.Equal(t, 1500, request.MaxCompletionTokens)
}

func TestAdaptAIRequestDropsTemperature(t *testing.T) {
	request := testCompletionRequest()
	request.ResponseFormat = nil
	request.MaxTokens = 0
	request.MaxCompletionTokens = 1500

	adaptation := adaptAIRequest(&request,
		`{"error":{"message":"Unsupported value: 'temperature' does not support 0.4 with this model"}}`)

	assert.Equal(t, constants.AIAdaptationDroppedTemperature, adaptation)
	assert.Nil(t, request.Temperature)
}

// Each branch requires its parameter to still be set, so a provider that keeps
// complaining runs out of adaptations instead of looping on the same one.
func TestAdaptAIRequestStopsWhenNothingIsLeftToChange(t *testing.T) {
	request := testCompletionRequest()
	providerBody := `{"error":{"message":"response_format, max_tokens and temperature are all wrong"}}`

	assert.Equal(t, constants.AIAdaptationDroppedJSONMode, adaptAIRequest(&request, providerBody))
	assert.Equal(t, constants.AIAdaptationMaxCompletionTokens, adaptAIRequest(&request, providerBody))
	assert.Equal(t, constants.AIAdaptationDroppedTemperature, adaptAIRequest(&request, providerBody))
	assert.Equal(t, "", adaptAIRequest(&request, providerBody))
}

func TestAdaptAIRequestIgnoresUnrelatedFailures(t *testing.T) {
	request := testCompletionRequest()

	adaptation := adaptAIRequest(&request, `{"error":{"message":"model not found"}}`)

	assert.Equal(t, "", adaptation)
	assert.NotNil(t, request.ResponseFormat)
	assert.Equal(t, 1500, request.MaxTokens)
}

func TestIsAIParameterRejection(t *testing.T) {
	assert.True(t, isAIParameterRejection(http.StatusBadRequest))
	assert.True(t, isAIParameterRejection(http.StatusUnprocessableEntity))
	assert.False(t, isAIParameterRejection(http.StatusUnauthorized))
	assert.False(t, isAIParameterRejection(http.StatusTooManyRequests))
	assert.False(t, isAIParameterRejection(http.StatusOK))
}

func TestCollectModelIDsFiltersAndCounts(t *testing.T) {
	listing := collectModelIDs([]modelListEntry{
		{ID: "llama-3.3-70b-versatile"},
		{ID: "whisper-large-v3"},
		{ID: "llama-3.1-8b-instant"},
		{ID: "playai-tts"},
		{ID: "llama-guard-3-8b"},
		{ID: "llama-3.3-70b-versatile"},
		{ID: "not a valid model id!"},
		{ID: ""},
	})

	assert.Equal(t, []string{"llama-3.1-8b-instant", "llama-3.3-70b-versatile"}, listing.Models)
	assert.Equal(t, 3, listing.Filtered)
	assert.False(t, listing.Truncated)
}

func TestCollectModelIDsReportsTruncation(t *testing.T) {
	entries := make([]modelListEntry, 0, constants.AIMaxListedModels+5)
	for index := 0; index < constants.AIMaxListedModels+5; index++ {
		entries = append(entries, modelListEntry{ID: "model-" + strconv.Itoa(index)})
	}

	listing := collectModelIDs(entries)

	assert.Len(t, listing.Models, constants.AIMaxListedModels)
	assert.True(t, listing.Truncated)
}
