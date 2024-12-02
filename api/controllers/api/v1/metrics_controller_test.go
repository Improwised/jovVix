package v1_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics(t *testing.T) {

	t.Run("get metrics with valid input", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/metrics")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)
	})
}
