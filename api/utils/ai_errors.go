package utils

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
)

var aiBudgetMarkers = []string{"budget", "quota", "usage limit", "limit exceeded", "exceeded your current"}

// ClassifyAIError turns a failed provider call into a message the caller can act
// on. The body is read only for the markers that tell an exhausted budget from a
// rejected key; nothing from it reaches the response.
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

// providerReason lifts the human readable message out of an OpenAI shaped error
// envelope. Only that one field is passed on: the rest of the body is the
// caller's own data and can carry their key back in an auth error.
func providerReason(body string) string {
	var envelope struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal([]byte(body), &envelope); err != nil {
		return ""
	}

	message := strings.Join(strings.Fields(envelope.Error.Message), " ")
	if message == "" || strings.Contains(strings.ToLower(message), "key") {
		return ""
	}

	return TruncateRunes(message, constants.AIMaxProviderReasonLength)
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
