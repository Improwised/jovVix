package utils

import (
	"strings"

	"github.com/Improwised/jovvix/api/constants"
)

// ValidateCoverImage checks a quiz cover image, which is stored as a base64
// data URI rather than in object storage. It returns a client-facing failure
// message, or an empty string when the image is acceptable. An empty cover
// image is valid — it means the quiz falls back to a static image.
func ValidateCoverImage(coverImage string) string {
	if coverImage == "" {
		return ""
	}
	if !strings.HasPrefix(coverImage, "data:image/") {
		return constants.ErrInvalidCoverImage
	}
	if len(coverImage) > constants.MaxCoverImageBytes {
		return constants.ErrCoverImageTooLarge
	}
	return ""
}
