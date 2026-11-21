package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteKratosIdentityUsesAdminURL(t *testing.T) {
	var gotMethod, gotPath string
	admin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod, gotPath = r.Method, r.URL.Path
		w.WriteHeader(http.StatusNoContent)
	}))
	defer admin.Close()

	// Trailing slash as in api/.env.example must not produce "//identities".
	err := deleteKratosIdentity(admin.URL+"/admin/", "8a6c1f6e-0000-4000-8000-000000000001")

	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, gotMethod)
	assert.Equal(t, "/admin/identities/8a6c1f6e-0000-4000-8000-000000000001", gotPath)
}

func TestDeleteKratosIdentityRefusesRedirectAndMissingURL(t *testing.T) {
	elsewhereHit := false
	elsewhere := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		elsewhereHit = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer elsewhere.Close()
	// Kratos's public port answers /admin/* like this.
	redirecting := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Location", elsewhere.URL+"/admin/identities/some-id")
		w.WriteHeader(http.StatusTemporaryRedirect)
	}))
	defer redirecting.Close()

	assert.Error(t, deleteKratosIdentity(redirecting.URL+"/admin", "some-id"))
	assert.False(t, elsewhereHit, "the delete request must not follow a redirect")
	assert.Error(t, deleteKratosIdentity("  ", "some-id"))
}
