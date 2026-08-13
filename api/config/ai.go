package config

import (
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

type AIConfig struct {
	Temperature    float64 `envconfig:"AI_TEMPERATURE" default:"0.4"`
	TimeoutSeconds int     `envconfig:"AI_TIMEOUT_SECONDS"`
	MaxQuestions   int     `envconfig:"AI_MAX_QUESTIONS"`
	JSONMode       bool    `envconfig:"AI_JSON_MODE"`
}

func (a AIConfig) Timeout() time.Duration {
	if a.TimeoutSeconds > 0 {
		return time.Duration(a.TimeoutSeconds) * time.Second
	}
	return constants.AIDefaultTimeoutSeconds * time.Second
}

func (a AIConfig) QuestionLimit() int {
	if a.MaxQuestions <= 0 {
		return constants.AIDefaultMaxQuestions
	}
	if a.MaxQuestions > constants.AIMaxQuestionsHardLimit {
		return constants.AIMaxQuestionsHardLimit
	}
	return a.MaxQuestions
}

func (a AIConfig) Temp() float64 {
	if a.Temperature < 0 || a.Temperature > 2 {
		return constants.AIDefaultTemperature
	}
	return a.Temperature
}
