package v1

import (
	"github.com/Improwised/quizz-app/api/models"
	pMetrics "github.com/Improwised/quizz-app/api/pkg/prometheus"
	goqu "github.com/doug-martin/goqu/v9"
	adaptor "github.com/gofiber/adaptor/v2"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type MetricsController struct {
	userModel *models.UserModel
	logger    *zap.Logger
	pMetrics  *pMetrics.PrometheusMetrics
}

func InitMetricsController(db *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) (*MetricsController, error) {
	userModel, err := models.InitUserModel(db, logger)
	if err != nil {
		return nil, err
	}
	return &MetricsController{
		userModel: &userModel,
		logger:    logger,
		pMetrics:  pMetrics,
	}, nil
}

// Prometheus metrics endpoint
// swagger:route GET /metrics Metrics ReqMetrics
//
//	Prometheus metrics endpoint
//
//	Prometheus metrics endpoint
//
// Produces:
// - text/plain
//
// Responses:
func (mc *MetricsController) Metrics(ctx *fiber.Ctx) error {
	users, err := mc.userModel.CountUsers()
	if err != nil {
		mc.logger.Error("error while getting user count", zap.Error(err))
	} else {
		mc.pMetrics.UserMetrics.Set(float64(users))
	}

	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
