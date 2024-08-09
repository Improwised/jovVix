package middlewares

import (
	"strings"

	pMetrics "github.com/Improwised/quizz-app/api/pkg/prometheus"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	zapCoreField []zapcore.Field
	// Add path that needs to excluded from logging
	ignorePathList = []string{
		"/docs",
		"/assets/redoc.css",
		"/assets/redoc.standalone.js",
		"/assets/swagger.json",
		"/favicon.ico",
	}
)

// Handler will log each request
func LogHandler(logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()
		if err != nil {
			return err
		}

		exits := lo.Contains(ignorePathList, ctx.Path()) || strings.HasPrefix(string(ctx.Response().Header.ContentType()), "image/") || strings.HasPrefix(string(ctx.Response().Header.ContentType()), "text/")
		if !exits {
			zapCoreField = []zapcore.Field{
				zap.String("host", ctx.Hostname()),
				zap.String("method", string(ctx.Request().Header.Method())),
				zap.String("uri", ctx.BaseURL()),
				zap.String("protocol", ctx.Protocol()),
				zap.String("username", string(ctx.Request().URI().Username())),
				zap.String("requestHeaders", string(ctx.Request().Header.Header())),
				zap.String("responseHeaders", string(ctx.Response().Header.Header())),
				zap.String("request", string(ctx.Request().Body())),
				zap.String("response", ctx.Response().String()),
				zap.Int("status", ctx.Response().Header.StatusCode()),
				zap.Int("size", ctx.Response().Header.ContentLength()),
			}
			if ctx.Response().Header.StatusCode() >= 100 && ctx.Response().Header.StatusCode() <= 399 {
				logger.Debug("Handled successful request", zapCoreField...)
			} else {
				logger.Debug("handled error request", zapCoreField...)
			}
		}

		// For /metrics endpoint count in next request
		// Because /metrics endpoint response is send first and
		// Respected status code counter increase next
		switch statusCode := ctx.Response().StatusCode(); {
		case statusCode >= 200 && statusCode < 300:
			pMetrics.RequestsMetrics.WithLabelValues("2xx").Inc()
		case statusCode >= 300 && statusCode < 400:
			pMetrics.RequestsMetrics.WithLabelValues("3xx").Inc()
		case statusCode >= 400 && statusCode < 500:
			pMetrics.RequestsMetrics.WithLabelValues("4xx").Inc()
		case statusCode >= 500:
			pMetrics.RequestsMetrics.WithLabelValues("5xx").Inc()
		}
		return nil

	}
}
