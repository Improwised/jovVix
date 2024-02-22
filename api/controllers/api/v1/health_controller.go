package v1

import (
	"context"
	"net/http"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/utils"
	goqu "github.com/doug-martin/goqu/v9"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type HealthController struct {
	db     *goqu.Database
	logger *zap.Logger
}

func NewHealthController(db *goqu.Database, logger *zap.Logger) (*HealthController, error) {
	return &HealthController{
		db:     db,
		logger: logger,
	}, nil
}

// Overall check overall health of application as well as dependencies health check
// swagger:route GET /api/v1/healthz Healthcheck overallHealthCheck
//
//	Overall health check
//
//	Overall health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Overall(ctx *fiber.Ctx) error {
	err := healthDb(hc.db)
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

func (hc *HealthController) Self(ctx *fiber.Ctx) error {
	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

// Database health check
// swagger:route GET /api/v1/healthz/db Healthcheck dbHealthCheck
//
//	Database health check
//
//	Database health check
//
// Produces:
// - application/json
//
// Responses:
//
//	200: GenericResOk
//	500: GenericResError
func (hc *HealthController) Db(ctx *fiber.Ctx) error {
	err := healthDb(hc.db)
	if err != nil {
		hc.logger.Error("error while health checking of db", zap.Error(err))
		return utils.JSONError(ctx, http.StatusInternalServerError, constants.ErrHealthCheckDb)
	}

	return utils.JSONSuccess(ctx, http.StatusOK, "ok")
}

///////////////////////
// HealthCheck CORE
//////////////////////

func healthDb(db *goqu.Database) error {
	_, err := db.ExecContext(context.TODO(), "SELECT 1")
	if err != nil {
		return err
	}
	return nil
}
