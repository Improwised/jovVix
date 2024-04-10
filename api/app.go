// Quizz-app
//
//	Schemes: http
//	Host: 127.0.0.1
//	BasePath: /api
//	Version: 0.0.1-alpha
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package main

import (
	"time"

	"github.com/Improwised/quizz-app/api/cli"
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/routinewrapper"
	sentry "github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func main() {

	// Collecting config from env or file or flag
	cfg := config.GetConfig()

	logger, err := logger.NewRootLogger(cfg.Debug, cfg.IsDevelopment)
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	// this function will logged error log in sentry
	sentryLoggedFunc := func() {
		err := recover()

		if err != nil {
			sentry.CurrentHub().Recover(err)
			sentry.Flush(time.Second * 2)
		}
	}

	// routine wrapper will handle go routine error also an log into sentry
	routinewrapper.Init(sentryLoggedFunc)
	defer sentryLoggedFunc()

	err = cli.Init(cfg, logger)
	if err != nil {
		panic(err)
	}

}
