package v1_test

import (
	"testing"

	"github.com/Improwised/quizz-app/api/config"
	v1 "github.com/Improwised/quizz-app/api/controllers/api/v1"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyticsBoardAdminController(t *testing.T) {
	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)

	events := events.NewEventBus(logger)

	err = events.SubscribeAll()
	assert.Nil(t, err)
	t.Run("check whether controller is being returned or not", func(t *testing.T) {

		analyticsAdminController, err := v1.NewAnalyticsBoardAdminController(db, logger, events, &cfg)
		assert.Nil(t, err)

		assert.NotNil(t, analyticsAdminController)
	})

}
