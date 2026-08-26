package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
)

func TestClassifyAIErrorStatusMapping(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   string
	}{
		{"unauthorized", http.StatusUnauthorized, `{"error":{"message":"bad key"}}`, constants.ErrAITestKeyRejected},
		{"payment required", http.StatusPaymentRequired, "", constants.ErrAIBudgetExhausted},
		{"forbidden", http.StatusForbidden, `{"error":{"message":"not allowed"}}`, constants.ErrAITestKeyForbidden},
		{"forbidden over budget", http.StatusForbidden, `{"error":{"message":"you exceeded your current quota"}}`, constants.ErrAIBudgetExhausted},
		{"rate limited", http.StatusTooManyRequests, `{"error":{"message":"slow down"}}`, constants.ErrAIRateLimited},
		{"rate limited over budget", http.StatusTooManyRequests, `{"error":{"message":"monthly budget reached"}}`, constants.ErrAIBudgetExhausted},
		{"server error", http.StatusInternalServerError, "", constants.ErrAITestServerError},
		{"bad gateway", http.StatusBadGateway, "", constants.ErrAITestServerError},
		{"unhandled status", http.StatusTeapot, "", constants.ErrAIRequestFailed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifyAIError(tc.status, tc.body, nil); got != tc.want {
				t.Fatalf("ClassifyAIError(%d) = %q, want %q", tc.status, got, tc.want)
			}
		})
	}
}

func TestClassifyAIErrorAppendsProviderReasonOnBadRequestAndNotFound(t *testing.T) {
	body := `{"error":{"message":"model gpt-9 does not exist"}}`

	badRequest := ClassifyAIError(http.StatusBadRequest, body, nil)
	if !strings.HasPrefix(badRequest, constants.ErrAITestBadRequest) {
		t.Fatalf("got %q, want it to start with the bad request message", badRequest)
	}
	if !strings.Contains(badRequest, "model gpt-9 does not exist") {
		t.Fatalf("got %q, want the provider reason appended", badRequest)
	}

	notFound := ClassifyAIError(http.StatusNotFound, body, nil)
	if !strings.HasPrefix(notFound, constants.ErrAITestNotFound) {
		t.Fatalf("got %q, want it to start with the not found message", notFound)
	}
}

// The provider body is the caller's own data coming back, so anything that looks
// like a credential must not be reflected into the response.
func TestProviderReasonDoesNotLeakCredentials(t *testing.T) {
	cases := []struct {
		body   string
		secret string
	}{
		{`{"error":{"message":"Incorrect API key provided: sk-abc123DEF456ghi789"}}`, "sk-abc123DEF456ghi789"},
		{`{"error":{"message":"invalid credential sk-live-0123456789abcdefghij"}}`, "sk-live-0123456789abcdefghij"},
		{`{"error":{"message":"token sk-proj-ZZZZZZZZZZZZZZZZZZZZ not recognized"}}`, "sk-proj-ZZZZZZZZZZZZZZZZZZZZ"},
		{`{"error":{"message":"Bearer sk-abcdefghijklmnopqrstuvwxyz is expired"}}`, "sk-abcdefghijklmnopqrstuvwxyz"},
		{`{"error":{"message":"your api-key AKIAIOSFODNN7EXAMPLE was rejected"}}`, "AKIAIOSFODNN7EXAMPLE"},
		{`{"error":{"message":"github token ghp_0123456789abcdefghijklmnopqrstuvwxyz denied"}}`, "ghp_0123456789abcdefghijklmnopqrstuvwxyz"},
		{`{"error":{"message":"rejected 8Fj2kLm9QpXvRt3wYz7NbAc5"}}`, "8Fj2kLm9QpXvRt3wYz7NbAc5"},
	}

	for _, tc := range cases {
		t.Run(tc.secret, func(t *testing.T) {
			reason := providerReason(tc.body)
			if strings.Contains(reason, tc.secret) {
				t.Fatalf("providerReason leaked the secret in %q", reason)
			}
			if !strings.Contains(reason, aiRedactedToken) {
				t.Fatalf("expected a redaction marker in %q", reason)
			}
		})
	}
}

func TestProviderReasonKeepsHarmlessMessages(t *testing.T) {
	reason := providerReason(`{"error":{"message":"the model is overloaded, retry shortly"}}`)
	if reason != "the model is overloaded, retry shortly" {
		t.Fatalf("got %q, want the message passed through", reason)
	}

	if got := providerReason(`{"error":{"message":"  spaced   out   message "}}`); got != "spaced out message" {
		t.Fatalf("got %q, want whitespace collapsed", got)
	}

	if got := providerReason("not json at all"); got != "" {
		t.Fatalf("got %q, want empty for an unparseable body", got)
	}

	if got := providerReason(`{"error":{"message":""}}`); got != "" {
		t.Fatalf("got %q, want empty for a blank message", got)
	}
}

func TestProviderReasonIsTruncated(t *testing.T) {
	long := strings.Repeat("word ", 200)
	reason := providerReason(`{"error":{"message":"` + long + `"}}`)
	if len([]rune(reason)) > constants.AIMaxProviderReasonLength {
		t.Fatalf("reason is %d runes, want at most %d", len([]rune(reason)), constants.AIMaxProviderReasonLength)
	}
}

func TestClassifyAITransportError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"deadline", context.DeadlineExceeded, constants.ErrAITestTimedOut},
		{"cancelled", context.Canceled, constants.ErrAITestTimedOut},
		{"blocked destination", ErrAIDisallowedDestination, constants.ErrAIBaseUrlPrivate},
		{"oversized response", ErrAIResponseTooLarge, constants.ErrAIResponseTooLarge},
		{"dns failure", &net.DNSError{Err: "no such host", IsNotFound: true}, constants.ErrAIBaseUrlUnresolvable},
		{"connection refused", &net.OpError{Op: "dial", Err: errors.New("connection refused")}, constants.ErrAIProviderUnreachable},
		{"redirect refused", errors.New("auto redirect is disabled"), constants.ErrAIProviderUnreachable},
		{"unknown", errors.New("something else entirely"), constants.ErrAIRequestFailed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ClassifyAIError(0, "", tc.err); got != tc.want {
				t.Fatalf("ClassifyAIError(%v) = %q, want %q", tc.err, got, tc.want)
			}
		})
	}
}

func TestClassifyAIErrorPrefersTransportErrorOverStatus(t *testing.T) {
	got := ClassifyAIError(http.StatusOK, `{"error":{"message":"ignored"}}`, ErrAIDisallowedDestination)
	if got != constants.ErrAIBaseUrlPrivate {
		t.Fatalf("got %q, want the transport error to win", got)
	}
}
