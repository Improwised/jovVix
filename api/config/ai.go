package config

import (
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

type AIConfig struct {
	Temperature    float64 `envconfig:"AI_TEMPERATURE" default:"0.4"`
	TimeoutSeconds int     `envconfig:"AI_TIMEOUT_SECONDS"`
	MaxQuestions   int     `envconfig:"AI_MAX_QUESTIONS"`
	JSONMode       bool    `envconfig:"AI_JSON_MODE" default:"true"`

	GenRateLimit   int `envconfig:"AI_GEN_RATE_LIMIT" default:"10"`
	GenRateWindow  int `envconfig:"AI_GEN_RATE_WINDOW" default:"60"`
	MetaRateLimit  int `envconfig:"AI_META_RATE_LIMIT" default:"30"`
	MetaRateWindow int `envconfig:"AI_META_RATE_WINDOW" default:"60"`
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

// GenerationRateLimit returns the per-user request budget for the outbound
// (generation/test/models) routes, clamped to a sane minimum.
func (a AIConfig) GenerationRateLimit() int {
	if a.GenRateLimit <= 0 {
		return constants.AIDefaultGenRateLimit
	}
	return a.GenRateLimit
}

func (a AIConfig) GenerationRateWindow() int {
	if a.GenRateWindow <= 0 {
		return constants.AIDefaultGenRateWindow
	}
	return a.GenRateWindow
}

func (a AIConfig) MetaRateLimitValue() int {
	if a.MetaRateLimit <= 0 {
		return constants.AIDefaultMetaRateLimit
	}
	return a.MetaRateLimit
}

func (a AIConfig) MetaRateWindowValue() int {
	if a.MetaRateWindow <= 0 {
		return constants.AIDefaultMetaRateWindow
	}
	return a.MetaRateWindow
}
