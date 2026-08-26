package utils

import "strings"

// A model list is the whole inventory, and embedding or image models fail on
// /chat/completions. Matched on markers, not an allowlist, so new models pass.
var nonChatModelMarkers = []string{
	"embedding",
	"embed-",
	"-embed",
	"whisper",
	"transcribe",
	"speech",
	"-tts",
	"tts-",
	"-audio",
	"realtime",
	"dall-e",
	"dalle",
	"imagen",
	"-image",
	"stable-diffusion",
	"veo-",
	"sora",
	"-video",
	"moderation",
	"guard",
	"-safety",
	"rerank",
}

// Legacy text-completion models: "-instruct" cannot be a marker, most
// open-weight chat models are named that way.
var nonChatModelIDs = map[string]bool{
	"gpt-3.5-turbo-instruct":      true,
	"gpt-3.5-turbo-instruct-0914": true,
	"davinci-002":                 true,
	"babbage-002":                 true,
	"code-davinci-002":            true,
}

func IsChatCapableModel(model string) bool {
	lowered := strings.ToLower(strings.TrimSpace(model))
	if lowered == "" || nonChatModelIDs[lowered] {
		return false
	}

	for _, marker := range nonChatModelMarkers {
		if strings.Contains(lowered, marker) {
			return false
		}
	}

	return true
}
