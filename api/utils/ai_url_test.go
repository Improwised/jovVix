package utils

import (
	"net"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
)

func TestIsDisallowedIP(t *testing.T) {
	cases := []struct {
		name       string
		ip         string
		disallowed bool
	}{
		{"nil is refused", "", true},
		{"loopback v4", "127.0.0.1", true},
		{"loopback v6", "::1", true},
		{"rfc1918 ten", "10.0.0.1", true},
		{"rfc1918 172", "172.16.5.4", true},
		{"rfc1918 192", "192.168.1.1", true},
		{"link local, cloud metadata", "169.254.169.254", true},
		{"unspecified", "0.0.0.0", true},
		{"carrier grade nat", "100.64.0.1", true},
		{"benchmarking range", "198.18.0.1", true},
		{"reserved class e", "240.0.0.1", true},
		{"multicast", "224.0.0.1", true},
		{"unique local v6", "fd00::1", true},
		{"link local v6", "fe80::1", true},
		{"ipv4 mapped loopback", "::ffff:127.0.0.1", true},
		{"ipv4 mapped metadata", "::ffff:169.254.169.254", true},
		{"nat64 embedded private", "64:ff9b::10.0.0.1", true},
		{"nat64 embedded metadata", "64:ff9b::169.254.169.254", true},

		{"public v4", "8.8.8.8", false},
		{"public v4 other", "1.1.1.1", false},
		{"public v6", "2606:4700:4700::1111", false},
		{"nat64 embedded public", "64:ff9b::8.8.8.8", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var ip net.IP
			if tc.ip != "" {
				ip = net.ParseIP(tc.ip)
				if ip == nil {
					t.Fatalf("could not parse %q as an ip", tc.ip)
				}
			}

			if got := IsDisallowedIP(ip); got != tc.disallowed {
				t.Fatalf("IsDisallowedIP(%q) = %v, want %v", tc.ip, got, tc.disallowed)
			}
		})
	}
}

// IP literals keep these cases off the network: lookupHost short circuits on a
// literal, so no resolver is involved.
func TestValidateAIBaseURLRejects(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr string
	}{
		{"empty", "", constants.ErrAIBaseUrlInvalid},
		{"blank", "   ", constants.ErrAIBaseUrlInvalid},
		{"not a url", "not a url", constants.ErrAIBaseUrlInvalid},
		{"no scheme", "api.example.com/v1", constants.ErrAIBaseUrlInvalid},
		{"unsupported scheme", "ftp://8.8.8.8/v1", constants.ErrAIBaseUrlInvalid},
		{"file scheme", "file:///etc/passwd", constants.ErrAIBaseUrlInvalid},
		{"userinfo", "https://user:pass@8.8.8.8/v1", constants.ErrAIBaseUrlInvalid},
		{"query string", "https://8.8.8.8/v1?key=secret", constants.ErrAIBaseUrlInvalid},
		{"fragment", "https://8.8.8.8/v1#part", constants.ErrAIBaseUrlInvalid},
		{"over length", "https://8.8.8.8/" + strings.Repeat("a", constants.AIMaxBaseUrlLength), constants.ErrAIBaseUrlInvalid},
		{"plain http", "http://8.8.8.8/v1", constants.ErrAIBaseUrlNotHTTPS},
		{"loopback", "https://127.0.0.1/v1", constants.ErrAIBaseUrlPrivate},
		{"private", "https://10.1.2.3/v1", constants.ErrAIBaseUrlPrivate},
		{"cloud metadata", "https://169.254.169.254/v1", constants.ErrAIBaseUrlPrivate},
		{"ipv6 loopback", "https://[::1]/v1", constants.ErrAIBaseUrlPrivate},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := ValidateAIBaseURL(tc.raw)
			if err == nil {
				t.Fatalf("ValidateAIBaseURL(%q) accepted %v, want error %q", tc.raw, parsed, tc.wantErr)
			}
			if err.Error() != tc.wantErr {
				t.Fatalf("ValidateAIBaseURL(%q) error = %q, want %q", tc.raw, err.Error(), tc.wantErr)
			}
		})
	}
}

func TestValidateAIBaseURLAccepts(t *testing.T) {
	cases := []string{
		"https://8.8.8.8/v1",
		"https://8.8.8.8:8443/v1",
		"  https://1.1.1.1/v1  ",
		"https://[2606:4700:4700::1111]/v1",
	}

	for _, raw := range cases {
		t.Run(raw, func(t *testing.T) {
			parsed, err := ValidateAIBaseURL(raw)
			if err != nil {
				t.Fatalf("ValidateAIBaseURL(%q) = %v, want no error", raw, err)
			}
			if parsed.Scheme != "https" {
				t.Fatalf("scheme = %q, want https", parsed.Scheme)
			}
		})
	}
}

func TestValidateAIModel(t *testing.T) {
	valid := []string{"gpt-4o-mini", "meta/llama-3.1-8b", "gemini-2.5-flash", "org@model_v1.2"}
	for _, model := range valid {
		if _, err := ValidateAIModel(model); err != nil {
			t.Fatalf("ValidateAIModel(%q) = %v, want no error", model, err)
		}
	}

	invalid := []string{"", "   ", "model name with spaces", "model\nname", "model;drop", strings.Repeat("m", constants.AIMaxModelLength+1)}
	for _, model := range invalid {
		if _, err := ValidateAIModel(model); err == nil {
			t.Fatalf("ValidateAIModel(%q) accepted, want error", model)
		}
	}
}

func TestValidateAIApiKey(t *testing.T) {
	key, err := ValidateAIApiKey("  sk-abc123  ")
	if err != nil || key != "sk-abc123" {
		t.Fatalf("ValidateAIApiKey trimmed = %q, %v", key, err)
	}

	if key, err := ValidateAIApiKey("   "); err != nil || key != "" {
		t.Fatalf("blank key = %q, %v, want empty and no error", key, err)
	}

	invalid := []string{"key with space", "key\nnewline", "key\ttab", strings.Repeat("k", constants.AIMaxApiKeyLength+1)}
	for _, candidate := range invalid {
		if _, err := ValidateAIApiKey(candidate); err == nil {
			t.Fatalf("ValidateAIApiKey(%q) accepted, want error", candidate)
		}
	}
}

func TestAIEndpointURL(t *testing.T) {
	cases := []struct {
		name string
		base string
		path string
		want string
	}{
		{"plain base", "https://api.example.com/v1", constants.AICompletionsPath, "https://api.example.com/v1/chat/completions"},
		{"trailing slash", "https://api.example.com/v1/", constants.AICompletionsPath, "https://api.example.com/v1/chat/completions"},
		{"padded", "  https://api.example.com/v1  ", constants.AICompletionsPath, "https://api.example.com/v1/chat/completions"},
		{"base already is the completions url", "https://api.example.com/v1/chat/completions", constants.AICompletionsPath, "https://api.example.com/v1/chat/completions"},
		{"completions url asked for models", "https://api.example.com/v1/chat/completions", constants.AIModelsPath, "https://api.example.com/v1/models"},
		{"models path", "https://api.example.com/v1", constants.AIModelsPath, "https://api.example.com/v1/models"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := AIEndpointURL(tc.base, tc.path); got != tc.want {
				t.Fatalf("AIEndpointURL(%q, %q) = %q, want %q", tc.base, tc.path, got, tc.want)
			}
		})
	}
}
