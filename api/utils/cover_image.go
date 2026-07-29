package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/Improwised/jovvix/api/constants"
)

var (
	ErrNotDataURI          = errors.New("cover image is not a data URI")
	ErrMalformedDataURI    = errors.New("cover image data URI is malformed")
	ErrUnsupportedEncoding = errors.New("cover image data URI is not base64 encoded")
	ErrUnsupportedMime     = errors.New("cover image type is not supported")
)

// ValidateCoverImage checks a quiz cover image, which is stored as a base64
// data URI rather than in object storage. It returns a client-facing failure
// message, or an empty string when the image is acceptable. An empty cover
// image is valid — it means the quiz falls back to a static image.
func ValidateCoverImage(coverImage string) string {
	if coverImage == "" {
		return ""
	}
	if len(coverImage) > constants.MaxCoverImageBytes {
		return constants.ErrCoverImageTooLarge
	}
	if _, _, err := DecodeCoverImage(coverImage); err != nil {
		return constants.ErrInvalidCoverImage
	}
	return ""
}

func DecodeCoverImage(dataURI string) (string, []byte, error) {
	const prefix = "data:"

	if !strings.HasPrefix(dataURI, prefix) {
		return "", nil, ErrNotDataURI
	}

	comma := strings.IndexByte(dataURI, ',')
	if comma < 0 {
		return "", nil, ErrMalformedDataURI
	}

	meta := dataURI[len(prefix):comma]
	if !strings.HasSuffix(strings.ToLower(meta), ";base64") {
		return "", nil, ErrUnsupportedEncoding
	}

	mime := strings.ToLower(strings.TrimSpace(strings.SplitN(meta, ";", 2)[0]))
	if !constants.AllowedCoverImageTypes[mime] {
		return "", nil, ErrUnsupportedMime
	}

	raw, err := base64.StdEncoding.DecodeString(dataURI[comma+1:])
	if err != nil {
		return "", nil, err
	}

	return mime, raw, nil
}
