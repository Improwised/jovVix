package cli

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	"github.com/Improwised/quizz-app/api/routes"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"github.com/gofiber/fiber/v2/middleware/cors"
	pMetrics "github.com/Improwised/quizz-app/api/pkg/prometheus"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {

			// Create fiber app
			app := fiber.New(fiber.Config{})

			app.Use(cors.New(cors.Config{
				AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
				AllowOrigins:     "*",
				AllowCredentials: true,
				AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
			}))

			promMetrics := pMetrics.InitPrometheusMetrics()

			// Init eventbus
			events := events.NewEventBus(logger)

			db, err := database.Connect(cfg.DB)
			if err != nil {
				return err
			}

			err = events.SubscribeAll()
			if err != nil {
				return err
			}

			pub, err := watermill.InitPublisher(cfg, false)
			if err != nil {
				return err
			}
			// setup routes
			err = routes.Setup(app, db, logger, cfg, events, promMetrics, pub)
			if err != nil {
				return err
			}

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				if err := app.Listen(cfg.Port); err != nil {
					logger.Panic(err.Error())
				}
			}()

			<-interrupt
			logger.Info("gracefully shutting down...")
			if err := app.Shutdown(); err != nil {
				logger.Panic("error while shutdown server", zap.Error(err))
			}

			logger.Info("server stopped to receive new requests or connection.")
			return nil
		},
	}

	return apiCommand
}
