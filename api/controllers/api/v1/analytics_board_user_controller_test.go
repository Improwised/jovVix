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

func TestNewAnalyticsBoardUserController(t *testing.T) {
	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)

	events := events.NewEventBus(logger)

	err = events.SubscribeAll()
	assert.Nil(t, err)
	t.Run("check whether controller is being returned or not", func(t *testing.T) {

		analyticsUserController, err := v1.NewAnalyticsBoardUserController(db, logger, events)
		assert.Nil(t, err)

		assert.NotNil(t, analyticsUserController)
	})

	
}
