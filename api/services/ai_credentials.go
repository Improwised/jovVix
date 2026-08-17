package services

import (
	"errors"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/utils"
)

// One caller's provider for a single request: the server holds no key of its own.
type AICredentials struct {
	BaseURL string
	APIKey  string
	Model   string
}

func (c AICredentials) CompletionsURL() string {
	return utils.AIEndpointURL(c.BaseURL, constants.AICompletionsPath)
}

func (c AICredentials) ModelsURL() string {
	return utils.AIEndpointURL(c.BaseURL, constants.AIModelsPath)
}

func ResolveAICredentials(baseUrl, apiKey, model string) (AICredentials, error) {
	return resolveAICredentials(baseUrl, apiKey, model, true)
}

// For calls that happen before a model is known, such as listing models.
func ResolveAICredentialsForListing(baseUrl, apiKey string) (AICredentials, error) {
	return resolveAICredentials(baseUrl, apiKey, "", false)
}

func resolveAICredentials(baseUrl, apiKey, model string, requireModel bool) (AICredentials, error) {
	baseUrl = strings.TrimSpace(baseUrl)
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)

	if baseUrl == "" && apiKey == "" && model == "" {
		return AICredentials{}, errors.New(constants.ErrAINotConfigured)
	}

	if baseUrl == "" || (requireModel && model == "") {
		return AICredentials{}, errors.New(constants.ErrAIIncompleteCredentials)
	}

	if _, err := utils.ValidateAIBaseURL(baseUrl); err != nil {
		return AICredentials{}, err
	}

	var resolvedModel string
	if model != "" {
		validated, err := utils.ValidateAIModel(model)
		if err != nil {
			return AICredentials{}, err
		}
		resolvedModel = validated
	}

	resolvedKey, err := utils.ValidateAIApiKey(apiKey)
	if err != nil {
		return AICredentials{}, err
	}

	return AICredentials{
		BaseURL: baseUrl,
		APIKey:  resolvedKey,
		Model:   resolvedModel,
	}, nil
}
