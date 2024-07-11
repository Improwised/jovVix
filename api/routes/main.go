package routes

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	controller "github.com/Improwised/quizz-app/api/controllers/api/v1"
	quizHelper "github.com/Improwised/quizz-app/api/helpers/quiz"
	"github.com/Improwised/quizz-app/api/middlewares"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	pMetrics "github.com/Improwised/quizz-app/api/pkg/prometheus"
	"github.com/Improwised/quizz-app/api/pkg/watermill"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	fiber "github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, events *events.Events, pMetrics *pMetrics.PrometheusMetrics, pub *watermill.WatermillPublisher) error {
	mu.Lock()
	defer mu.Unlock()

	// plugins
	app.Use(middlewares.LogHandler(logger, pMetrics))

	swagger_file_path := "./assets/swagger.json"
	swagger_new_file_path := "./assets/new_swagger.json"

	err := newSwagger(swagger_file_path, swagger_new_file_path, config.Port)
	if err != nil {
		return err
	}

	cfg := swagger.Config{
		FilePath: swagger_new_file_path,
		Title:    "Swagger API Docs",
		Path:     "/api/docs",
	}

	app.Use(swagger.New(cfg))

	app.Static("/assets/", "./assets")

	router := app.Group("/api")

	err = healthCheckController(router, goqu, logger)
	if err != nil {
		return err
	}

	err = metricsController(router, goqu, logger, pMetrics)
	if err != nil {
		return err
	}

	helperStructs, err := quizHelper.InitHelper(goqu, config.RedisClient, logger)

	if err != nil {
		return err
	}

	// middleware initialization
	middleware := middlewares.NewMiddleware(config, logger, goqu, helperStructs.UserService)
	playedQuizValidationMiddleware := middlewares.NewPlayedQuizValidationMiddleware(config, logger, goqu, helperStructs)

	v1 := router.Group("/v1")

	v1.Use("/socket", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals(constants.MiddlewareError, nil)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// FinalScoreboard
	err = setUpFinalScoreBoardController(v1, goqu, logger, middleware, events, pub)
	if err != nil {
		return err
	}
	err = setUpAnalyticsBoardController(v1, goqu, logger, middleware, events, pub)
	if err != nil {
		return err
	}

	err = setupAuthController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupUserController(v1, goqu, logger, middleware, events, pub)
	if err != nil {
		return err
	}

	err = quizController(v1, goqu, logger, middleware, playedQuizValidationMiddleware, events, pub, config, helperStructs)
	if err != nil {
		return err
	}

	return nil
}

func newSwagger(file_name, new_file, port string) error {
	// Verify Swagger file exists
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		return fmt.Errorf("%s file does not exist", file_name)
	}

	// Read Swagger Spec into memory
	rawSpec, err := os.ReadFile(file_name)
	if err != nil {
		return fmt.Errorf("failed to read provided Swagger file (%s): %v", file_name, err.Error())
	}

	// Validate we have valid JSON or YAML
	var jsonData map[string]interface{}
	errJSON := json.Unmarshal(rawSpec, &jsonData)
	if errJSON != nil {
		return fmt.Errorf("swagger-json is not in valid format")
	}
	jsonData["host"] = port

	newData, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		return fmt.Errorf("error during host change in swagger")
	}

	file, err := os.Create(new_file)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(newData)

	return err
}

func setupAuthController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, config config.AppConfig) error {
	authController, err := controller.NewAuthController(goqu, logger, config)
	if err != nil {
		return err
	}
	v1.Post("/login", authController.DoAuth)

	if config.Kratos.IsEnabled {
		kratos := v1.Group("/kratos")
		kratos.Get("/auth", middlewares.Authenticated, authController.DoKratosAuth)
	}
	return nil
}

func setupUserController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, events *events.Events, pub *watermill.WatermillPublisher) error {
	userController, err := controller.NewUserController(goqu, logger, events, pub)
	if err != nil {
		return err
	}

	// user route
	userRouter := v1.Group("/user")
	userRouter.Get("/is_admin", middlewares.Authenticated, userController.IsAdmin)
	userRouter.Get("/who", middlewares.Authenticated, userController.GetUserMeta)

	// users route
	usersRouter := v1.Group("/users")
	usersRouter.Get(fmt.Sprintf("/:%s", constants.ParamUid), middlewares.Authenticated, userController.GetUser)
	usersRouter.Post("/", userController.CreateUser)
	return nil
}

func healthCheckController(api fiber.Router, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := api.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}

func metricsController(api fiber.Router, db *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error {
	metricsController, err := controller.InitMetricsController(db, logger, pMetrics)
	if err != nil {
		return nil
	}

	api.Get("/metrics", metricsController.Metrics)
	return nil
}

func quizController(
	v1 fiber.Router,
	db *goqu.Database,
	logger *zap.Logger,
	middleware middlewares.Middleware,
	sessionMiddle middlewares.PlayedQuizValidationMiddleware,
	events *events.Events,
	pub *watermill.WatermillPublisher,
	config config.AppConfig,
	helper *quizHelper.HelperGroup) error {
		
	AnswersSubmittedByUsers := make(chan models.User, 50)

	quizSocketController := controller.InitQuizConfig(db, &config, logger, helper, AnswersSubmittedByUsers)
	quizController := controller.InitQuizController(logger, events, pub, helper, AnswersSubmittedByUsers)

	// middleware format := param-check/pass... , authentication... , authorization..., controller(API/SOCKET)...

	v1.Get(fmt.Sprintf("/socket/join/:%s", constants.QuizSessionInvitationCode), middleware.CheckSessionCode, middleware.CustomAuthenticated, sessionMiddle.PlayedQuizValidation, websocket.New(quizSocketController.Join))

	v1.Post("/quiz/answer", middleware.Authenticated, middleware.CustomAuthenticated, quizController.SetAnswer)

	v1.Get("/quiz/terminate", middleware.Authenticated, quizController.Terminate)

	// admin endpoints
	allowRoles, err := helper.RoleModel.NewAllowedRoles("admin")
	if err != nil {
		return err
	}
	rbObj := middlewares.NewRolePermissionMiddleware(middleware, allowRoles)

	admin := v1.Group("/admin")
	admin.Use(middleware.Authenticated, rbObj.IsAllowed)

	quizzes := admin.Group("/quizzes")

	quizzes.Get("/", quizController.GetQuizByUser)
	quizzes.Post(fmt.Sprintf("/:%s/demo_session", constants.QuizId), quizController.GenerateDemoSession)
	quizzes.Post(fmt.Sprintf("/:%s/upload", constants.QuizTitle), middleware.ValidateCsv, quizController.CreateQuizByCsv)
	quizzes.Get("/list", quizController.GetAdminUploadedQuizzes)

	v1.Get(fmt.Sprintf("/socket/admin/arrange/:%s", constants.SessionIDParam), middleware.CheckSessionId, middleware.CustomAdminAuthenticated, websocket.New(quizSocketController.Arrange))

	return nil
}

// final score board controller setup
func setUpFinalScoreBoardController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, events *events.Events, pub *watermill.WatermillPublisher) error {
	userController, err := controller.NewUserController(goqu, logger, events, pub)
	if err != nil {
		return err
	}

	finalScoreBoardController, err := controller.NewFinalScoreBoardController(goqu, logger, events)
	if err != nil {
		return err
	}

	finalScoreBoardControllerAdmin, err := controller.NewFinalScoreBoardAdminController(goqu, logger, events)
	if err != nil {
		return err
	}

	finalScore := v1.Group("/final_score")
	finalScore.Get("/user", finalScoreBoardController.GetScore)
	finalScore.Get("/admin", middlewares.Authenticated, userController.IsAdmin, finalScoreBoardControllerAdmin.GetScoreForAdmin)

	return nil
}

func setUpAnalyticsBoardController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, events *events.Events, pub *watermill.WatermillPublisher) error {
	userController, err := controller.NewUserController(goqu, logger, events, pub)
	if err != nil {
		return err
	}

	analyticsBoardUserController, err := controller.NewAnalyticsBoardUserController(goqu, logger, events)
	if err != nil {
		return err
	}

	analyticsBoardAdminController, err := controller.NewAnalyticsBoardAdminController(goqu, logger, events)
	if err != nil {
		return err
	}

	analyticsBoard := v1.Group("/analytics_board")
	analyticsBoard.Get("/user", analyticsBoardUserController.GetAnalyticsForUser)
	analyticsBoard.Get("/admin", middlewares.Authenticated, userController.IsAdmin, analyticsBoardAdminController.GetAnalyticsForAdmin)

	return nil
}
