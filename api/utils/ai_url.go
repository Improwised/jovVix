package utils

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

// Returned by the dial guard so the classifier need not match on message text.
var ErrAIDisallowedDestination = errors.New(constants.ErrAIBaseUrlPrivate)

// The destination is caller supplied, so the response size must be bounded.
var ErrAIResponseTooLarge = errors.New(constants.ErrAIResponseTooLarge)

var (
	extraBlockedCIDRs = []string{
		"100.64.0.0/10",
		"192.0.0.0/24",
		"198.18.0.0/15",
		"240.0.0.0/4",
	}
	blockedNets []*net.IPNet

	nat64Prefix = net.IPNet{
		IP:   net.ParseIP("64:ff9b::"),
		Mask: net.CIDRMask(96, 128),
	}

	aiModelPattern = regexp.MustCompile(`^[A-Za-z0-9._:/@-]+$`)
)

func init() {
	for _, cidr := range extraBlockedCIDRs {
		if _, network, err := net.ParseCIDR(cidr); err == nil {
			blockedNets = append(blockedNets, network)
		}
	}
}

func IsDisallowedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	if nat64Prefix.Contains(ip) && ip.To4() == nil {
		if embedded := net.IPv4(ip[12], ip[13], ip[14], ip[15]); IsDisallowedIP(embedded) {
			return true
		}
	}

	if mapped := ip.To4(); mapped != nil {
		ip = mapped
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() ||
		ip.IsInterfaceLocalMulticast() || ip.IsMulticast() {
		return true
	}

	for _, network := range blockedNets {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func ValidateAIBaseURL(raw string) (*url.URL, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || len(trimmed) > constants.AIMaxBaseUrlLength {
		return nil, errors.New(constants.ErrAIBaseUrlInvalid)
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return nil, errors.New(constants.ErrAIBaseUrlInvalid)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, errors.New(constants.ErrAIBaseUrlInvalid)
	}

	if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil, errors.New(constants.ErrAIBaseUrlInvalid)
	}

	host := parsed.Hostname()
	if host == "" || len(host) > 253 {
		return nil, errors.New(constants.ErrAIBaseUrlInvalid)
	}

	if port := parsed.Port(); port != "" {
		number, err := net.LookupPort("tcp", port)
		if err != nil || number <= 0 {
			return nil, errors.New(constants.ErrAIBaseUrlInvalid)
		}
	}

	if parsed.Scheme != "https" {
		return nil, errors.New(constants.ErrAIBaseUrlNotHTTPS)
	}

	addresses, err := lookupHost(host)
	if err != nil || len(addresses) == 0 {
		return nil, errors.New(constants.ErrAIBaseUrlUnresolvable)
	}

	for _, address := range addresses {
		if IsDisallowedIP(address) {
			return nil, errors.New(constants.ErrAIBaseUrlPrivate)
		}
	}

	return parsed, nil
}

func lookupHost(host string) ([]net.IP, error) {
	if literal := net.ParseIP(host); literal != nil {
		return []net.IP{literal}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.AIDialTimeoutSeconds*time.Second)
	defer cancel()

	return net.DefaultResolver.LookupIP(ctx, "ip", host)
}

func NewSSRFGuardedTransport() *http.Transport {
	dialer := &net.Dialer{
		Timeout:   constants.AIDialTimeoutSeconds * time.Second,
		KeepAlive: 30 * time.Second,
		Control: func(network, address string, _ syscall.RawConn) error {
			host, _, err := net.SplitHostPort(address)
			if err != nil || IsDisallowedIP(net.ParseIP(host)) {
				return ErrAIDisallowedDestination
			}
			return nil
		},
	}

	return &http.Transport{
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          20,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

// Wraps the SSRF guarded transport with a response size cap.
func NewAIProviderTransport() http.RoundTripper {
	return &limitedTransport{base: NewSSRFGuardedTransport(), limit: constants.AIMaxResponseBytes}
}

type limitedTransport struct {
	base  http.RoundTripper
	limit int64
}

func (t *limitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if res.ContentLength > t.limit {
		res.Body.Close()
		return nil, ErrAIResponseTooLarge
	}

	res.Body = &limitedBody{inner: res.Body, limit: t.limit}
	return res, nil
}

// Fails the read rather than truncating, so an oversized body is reported as
// such instead of surfacing as a json parse error.
type limitedBody struct {
	inner io.ReadCloser
	limit int64
	read  int64
}

func (b *limitedBody) Read(p []byte) (int, error) {
	if b.read > b.limit {
		return 0, ErrAIResponseTooLarge
	}

	n, err := b.inner.Read(p)
	b.read += int64(n)
	if b.read > b.limit {
		return n, ErrAIResponseTooLarge
	}

	return n, err
}

func (b *limitedBody) Close() error {
	return b.inner.Close()
}

func ValidateAIModel(model string) (string, error) {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" || len(trimmed) > constants.AIMaxModelLength || !aiModelPattern.MatchString(trimmed) {
		return "", errors.New(constants.ErrAIModelInvalid)
	}
	return trimmed, nil
}

func ValidateAIApiKey(key string) (string, error) {
	trimmed := strings.TrimSpace(key)
	if trimmed == "" {
		return "", nil
	}
	if len(trimmed) > constants.AIMaxApiKeyLength {
		return "", errors.New(constants.ErrAIApiKeyInvalid)
	}
	for _, character := range trimmed {
		if character < 0x21 || character > 0x7e {
			return "", errors.New(constants.ErrAIApiKeyInvalid)
		}
	}
	return trimmed, nil
}

// Joins a base URL with an endpoint path, accepting a base that already has it.
func AIEndpointURL(baseUrl, path string) string {
	base := strings.TrimSuffix(strings.TrimSpace(baseUrl), "/")
	if strings.HasSuffix(base, path) {
		return base
	}

	// People paste the completions url as their base url.
	base = strings.TrimSuffix(base, constants.AICompletionsPath)

	return base + path
}
