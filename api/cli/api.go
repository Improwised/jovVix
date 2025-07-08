package cli

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/database"
	"github.com/Improwised/jovvix/api/pkg/events"
	pMetrics "github.com/Improwised/jovvix/api/pkg/prometheus"
	"github.com/Improwised/jovvix/api/pkg/watermill"
	"github.com/Improwised/jovvix/api/routes"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/cobra"
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
				AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization,Options",
				AllowOrigins:     "http://*:5000, ws://*:3300, wss://*:3300, ws://*:3000, wss://*:3000, http://0.0.0.0:5000, http://127.0.0.1:5000, http://127.0.0.1:3000, "+ cfg.WebUrl,
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
				logger.Error(err.Error())
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

func GetOutboundIP(logger *zap.Logger) net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        logger.Fatal(err.Error())
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}
