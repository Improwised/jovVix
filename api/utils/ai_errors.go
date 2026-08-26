package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
)

var aiBudgetMarkers = []string{"budget", "quota", "usage limit", "limit exceeded", "exceeded your current"}

// Turns a failed provider call into a message the caller can act on.
func ClassifyAIError(status int, body string, transportErr error) string {
	if transportErr != nil {
		return classifyAITransportError(transportErr)
	}

	budgetRelated := containsAIBudgetMarker(body)
	reason := providerReason(body)

	switch status {
	case http.StatusBadRequest:
		return withReason(constants.ErrAITestBadRequest, reason)
	case http.StatusUnauthorized:
		return constants.ErrAITestKeyRejected
	case http.StatusPaymentRequired:
		return constants.ErrAIBudgetExhausted
	case http.StatusForbidden:
		if budgetRelated {
			return constants.ErrAIBudgetExhausted
		}
		return constants.ErrAITestKeyForbidden
	case http.StatusNotFound:
		return withReason(constants.ErrAITestNotFound, reason)
	case http.StatusTooManyRequests:
		if budgetRelated {
			return constants.ErrAIBudgetExhausted
		}
		return constants.ErrAIRateLimited
	}

	if status >= http.StatusInternalServerError {
		return constants.ErrAITestServerError
	}
	return constants.ErrAIRequestFailed
}

// Lifts the message out of an OpenAI shaped error envelope, redacting anything
// that looks like a credential: an auth error routinely quotes the key.
func providerReason(body string) string {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return ""
	}

	words := strings.Fields(envelope.Error.Message)
	for i, word := range words {
		if looksLikeCredential(word) {
			words[i] = aiRedactedToken
		}
	}

	message := strings.Join(words, " ")
	if strings.TrimSpace(strings.ReplaceAll(message, aiRedactedToken, "")) == "" {
		return ""
	}

	return TruncateRunes(message, constants.AIMaxProviderReasonLength)
}

const aiRedactedToken = "[redacted]"

var (
	credentialPrefixes = []string{"sk-", "sk_", "pk-", "pk_", "api-", "api_", "key-", "key_", "token-", "token_", "akia", "asia", "ghp_", "xox"}
	opaqueTokenPattern = regexp.MustCompile(`^[A-Za-z0-9+/=_-]{20,}$`)
)

// Errs towards redacting: losing a diagnostic word beats leaking a key.
func looksLikeCredential(word string) bool {
	trimmed := strings.Trim(word, `.,;:!?"'()[]{}`)
	if len(trimmed) < 8 {
		return false
	}

	lowered := strings.ToLower(trimmed)
	for _, prefix := range credentialPrefixes {
		if strings.Contains(lowered, prefix) {
			return true
		}
	}

	return opaqueTokenPattern.MatchString(trimmed) && strings.ContainsAny(trimmed, "0123456789")
}

func withReason(message, reason string) string {
	if reason == "" {
		return message
	}
	return message + ". provider said: " + reason
}

func containsAIBudgetMarker(body string) bool {
	lowered := strings.ToLower(body)
	for _, marker := range aiBudgetMarkers {
		if strings.Contains(lowered, marker) {
			return true
		}
	}
	return false
}

func classifyAITransportError(err error) string {
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return constants.ErrAITestTimedOut
	}

	if errors.Is(err, ErrAIDisallowedDestination) {
		return constants.ErrAIBaseUrlPrivate
	}

	if errors.Is(err, ErrAIResponseTooLarge) {
		return constants.ErrAIResponseTooLarge
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return constants.ErrAIBaseUrlUnresolvable
	}

	var certErr *tls.CertificateVerificationError
	var recordErr tls.RecordHeaderError
	if errors.As(err, &certErr) || errors.As(err, &recordErr) {
		return constants.ErrAITestTLS
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return constants.ErrAITestTimedOut
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return constants.ErrAIProviderUnreachable
	}

	if strings.Contains(strings.ToLower(err.Error()), "redirect") {
		return constants.ErrAIProviderUnreachable
	}

	return constants.ErrAIRequestFailed
}
