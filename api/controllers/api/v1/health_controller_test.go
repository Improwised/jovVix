package v1_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverall(t *testing.T) {

	t.Run("check overall health of application", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/healthz")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}

func TestDb(t *testing.T) {

	t.Run("check overall health of application", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/healthz/db")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
