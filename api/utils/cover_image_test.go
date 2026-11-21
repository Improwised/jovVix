package utils

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/stretchr/testify/assert"
)

// A one-pixel GIF, small enough to inline and a real decodable image.
const gifPayload = "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"

func dataURI(mime, payload string) string {
	return "data:" + mime + ";base64," + payload
}

func TestDecodeCoverImageAllowedTypes(t *testing.T) {
	for _, mime := range []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/heic", "image/heif"} {
		gotMime, raw, err := DecodeCoverImage(dataURI(mime, gifPayload))

		assert.NoError(t, err, mime)
		assert.Equal(t, mime, gotMime)
		assert.NotEmpty(t, raw)
	}
}

func TestDecodeCoverImageRejectsSVG(t *testing.T) {
	svg := base64.StdEncoding.EncodeToString([]byte(`<svg xmlns="http://www.w3.org/2000/svg"><script>alert(1)</script></svg>`))

	_, _, err := DecodeCoverImage(dataURI("image/svg+xml", svg))

	// The MIME is echoed back as Content-Type, so an SVG here would be a
	// script-execution vector on the API origin.
	assert.ErrorIs(t, err, ErrUnsupportedMime)
}

func TestDecodeCoverImageRejectsMalformedInput(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected error
	}{
		{"not a data uri", "https://example.com/a.png", ErrNotDataURI},
		{"no comma", "data:image/png;base64", ErrMalformedDataURI},
		{"not base64 encoded", "data:image/png,rawbytes", ErrUnsupportedEncoding},
		{"disallowed type", dataURI("application/pdf", gifPayload), ErrUnsupportedMime},
		{"empty string", "", ErrNotDataURI},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := DecodeCoverImage(tc.input)
			assert.ErrorIs(t, err, tc.expected)
		})
	}
}

func TestDecodeCoverImageRejectsCorruptBase64(t *testing.T) {
	_, _, err := DecodeCoverImage(dataURI("image/png", "not!valid!base64"))

	assert.Error(t, err)
}

func TestDecodeCoverImageNormalisesMimeCasing(t *testing.T) {
	mime, _, err := DecodeCoverImage("data:IMAGE/PNG;BASE64," + gifPayload)

	assert.NoError(t, err)
	assert.Equal(t, "image/png", mime)
}

func TestValidateCoverImage(t *testing.T) {
	// Empty is valid — the quiz falls back to static art.
	assert.Equal(t, "", ValidateCoverImage(""))

	assert.Equal(t, "", ValidateCoverImage(dataURI("image/png", gifPayload)))

	svg := base64.StdEncoding.EncodeToString([]byte("<svg/>"))
	assert.Equal(t, constants.ErrInvalidCoverImage, ValidateCoverImage(dataURI("image/svg+xml", svg)))

	assert.Equal(t, constants.ErrInvalidCoverImage, ValidateCoverImage("https://example.com/a.png"))

	oversized := dataURI("image/png", strings.Repeat("A", constants.MaxCoverImageBytes))
	assert.Equal(t, constants.ErrCoverImageTooLarge, ValidateCoverImage(oversized))
}

func TestDecodeCoverImageMimeParameters(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"spaces around mime", "data: image/png ;base64," + gifPayload},
		{"extra parameter", "data:image/png;charset=utf-8;base64," + gifPayload},
		{"uppercase with parameter", "data:IMAGE/PNG;CHARSET=UTF-8;BASE64," + gifPayload},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mime, raw, err := DecodeCoverImage(tc.input)

			assert.NoError(t, err)
			assert.Equal(t, "image/png", mime)
			assert.NotEmpty(t, raw)
		})
	}
}

func TestDecodeCoverImageReturnsExactBytes(t *testing.T) {
	want, err := base64.StdEncoding.DecodeString(gifPayload)
	assert.NoError(t, err)

	_, got, err := DecodeCoverImage(dataURI("image/gif", gifPayload))

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestValidateCoverImageSizeBoundary(t *testing.T) {
	// "data:image/png;foo=barrr;base64," is 32 bytes, so a payload of
	// MaxCoverImageBytes-32 'A's makes the whole URI exactly MaxCoverImageBytes
	// long. The limit is `len > MaxCoverImageBytes`, so exactly at the limit
	// is accepted and one byte over is rejected.
	prefix := "data:image/png;foo=barrr;base64,"
	payload := strings.Repeat("A", constants.MaxCoverImageBytes-len(prefix))

	assert.Equal(t, "", ValidateCoverImage(prefix+payload))
	assert.Equal(t, constants.ErrCoverImageTooLarge, ValidateCoverImage(prefix+payload+"A"))
}

func TestValidateCoverImageRejectsCorruptBase64(t *testing.T) {
	assert.Equal(t, constants.ErrInvalidCoverImage, ValidateCoverImage("data:image/png;base64,not!valid!base64"))
}
