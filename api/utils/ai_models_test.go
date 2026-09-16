package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsChatCapableModel(t *testing.T) {
	chatModels := []string{
		"gpt-4o-mini",
		"gpt-4.1",
		"gpt-5-mini",
		"o3-mini",
		"llama-3.3-70b-versatile",
		"llama-3.1-8b-instant",
		"deepseek/deepseek-chat",
		"meta-llama/llama-3.3-70b-instruct",
		"mistral-7b-instruct",
		"gemma2-9b-it",
		"gemini-2.5-flash",
		"gemini-2.5-pro",
		"qwen/qwen3-32b",
		"openai/gpt-oss-120b",
		"anthropic/claude-sonnet-4.5",
	}

	for _, model := range chatModels {
		assert.True(t, IsChatCapableModel(model), "expected %q to be listed", model)
	}

	nonChatModels := []string{
		"",
		"   ",
		"text-embedding-3-small",
		"text-embedding-004",
		"nomic-embed-text",
		"embed-english-v3.0",
		"whisper-1",
		"whisper-large-v3",
		"gpt-4o-transcribe",
		"gpt-4o-mini-tts",
		"playai-tts",
		"tts-1-hd",
		"gpt-4o-audio-preview",
		"gpt-realtime",
		"dall-e-3",
		"imagen-3.0-generate-002",
		"google/gemini-2.5-flash-image-preview",
		"stable-diffusion-xl",
		"veo-3.0-generate-preview",
		"sora-2",
		"omni-moderation-latest",
		"llama-guard-3-8b",
		"meta-llama/llama-prompt-guard-2-86m",
		"rerank-v3.5",
		"gpt-3.5-turbo-instruct",
		"davinci-002",
		"babbage-002",
	}

	for _, model := range nonChatModels {
		assert.False(t, IsChatCapableModel(model), "expected %q to be filtered", model)
	}
}

func TestIsChatCapableModelIgnoresCaseAndPadding(t *testing.T) {
	assert.False(t, IsChatCapableModel("  TEXT-EMBEDDING-3-SMALL  "))
	assert.True(t, IsChatCapableModel("  GPT-4o-Mini  "))
}
