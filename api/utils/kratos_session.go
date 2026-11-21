package utils

import (
	"net/http"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/go-resty/resty/v2"
)

// KratosSessionFailure maps a failed Kratos whoami call to the status and message the client gets.
// Kratos answers 401 for a missing, expired or revoked session: the client gets 401 and a message
// a person understands (GitHub #77). Other Kratos statuses pass through. A transport error has no
// status (resty reports 0), so it becomes 502.
func KratosSessionFailure(res *resty.Response, err error) (int, string) {
	if err != nil || res == nil || res.StatusCode() == 0 {
		return http.StatusBadGateway, constants.ErrKratosAuth
	}
	if res.StatusCode() == http.StatusUnauthorized {
		return http.StatusUnauthorized, constants.ErrSessionExpired
	}
	return res.StatusCode(), constants.ErrKratosAuth
}
