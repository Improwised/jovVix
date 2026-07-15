package routes

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	controller "github.com/Improwised/jovvix/api/controllers/api/v1"
	"github.com/Improwised/jovvix/api/middlewares"
	pMetrics "github.com/Improwised/jovvix/api/pkg/prometheus"
	"github.com/Improwised/jovvix/api/pkg/redis"
	goqu "github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/contrib/websocket"
	fiber "github.com/gofiber/fiber/v2"
)

var mu sync.Mutex

// Setup func
func Setup(app *fiber.App, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, pMetrics *pMetrics.PrometheusMetrics) error {
	mu.Lock()
	defer mu.Unlock()

	// plugins
	app.Use(middlewares.LogHandler(logger, pMetrics))

	swagger_file_path := "./assets/swagger.json"
	swagger_new_file_path := "./assets/new_swagger.json"

	err := newSwagger(swagger_file_path, swagger_new_file_path, config.WebUrl)
	if err != nil {
		return err
	}

	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/",
		FilePath: swagger_new_file_path,
		Path:     "docs",
		Title:    "Swagger API Docs",
	}))

	router := app.Group("/api")

	err = setupHealthCheckController(router, goqu, logger)
	if err != nil {
		return err
	}

	err = setupMetricsController(router, goqu, logger, pMetrics)
	if err != nil {
		return err
	}

	redis, err := redis.InitRedisPubSub(goqu, config.RedisClient, logger)

	if err != nil {
		return err
	}

	// middleware initialization
	middleware := middlewares.NewMiddleware(config, logger, goqu)

	v1 := router.Group("/v1")

	v1.Use("/socket", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals(constants.MiddlewareError, nil)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// FinalScoreboard
	err = setUpFinalScoreBoardController(v1, goqu, logger, middleware)
	if err != nil {
		return err
	}
	err = setUpAnalyticsBoardController(v1, goqu, logger, config, middleware)
	if err != nil {
		return err
	}

	err = setupAuthController(v1, goqu, logger, config, middleware)
	if err != nil {
		return err
	}

	err = setupUserController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupQuizSocketController(v1, goqu, logger, middleware, config, redis)
	if err != nil {
		return err
	}

	err = setupQuizController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupQuizCategoryController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupQuestionController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupUserPlayedQuizeController(v1, goqu, logger, middleware, config)
	if err != nil {
		return err
	}

	err = setupSharedQuizzesController(v1, goqu, logger, middleware, config)
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

func setupAuthController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, middlewares middlewares.Middleware) error {
	authController, err := controller.NewAuthController(goqu, logger, config)
	if err != nil {
		return err
	}

	if config.Kratos.IsEnabled {
		kratos := v1.Group("/kratos")
		kratos.Get("/auth", authController.DoKratosAuth)
		kratos.Get("/whoami", authController.GetRegisteredUser)
		kratos.Put("/user", middlewares.KratosAuthenticated, authController.UpadateRegisteredUser)
		kratos.Delete("/user", middlewares.KratosAuthenticated, authController.DeleteRegisteredUser)
	}
	return nil
}

func setupUserController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, config config.AppConfig) error {
	userController, err := controller.NewUserController(goqu, logger, config)
	if err != nil {
		return err
	}

	// user route
	userRouter := v1.Group("/user")
	userRouter.Get("/who", middlewares.Authenticated, userController.GetUserMeta)
	userRouter.Post(fmt.Sprintf("/:%s", constants.Username), userController.CreateGuestUser)

	return nil
}

func setupHealthCheckController(api fiber.Router, goqu *goqu.Database, logger *zap.Logger) error {
	healthController, err := controller.NewHealthController(goqu, logger)
	if err != nil {
		return err
	}

	healthz := api.Group("/healthz")
	healthz.Get("/", healthController.Overall)
	healthz.Get("/db", healthController.Db)
	return nil
}

func setupMetricsController(api fiber.Router, db *goqu.Database, logger *zap.Logger, pMetrics *pMetrics.PrometheusMetrics) error {
	metricsController, err := controller.InitMetricsController(db, logger, pMetrics)
	if err != nil {
		return nil
	}

	api.Get("/metrics", metricsController.Metrics)
	return nil
}

func setupQuizSocketController(v1 fiber.Router, db *goqu.Database, logger *zap.Logger, middleware middlewares.Middleware, config config.AppConfig, redis *redis.RedisPubSub) error {
	quizSocketController, err := controller.InitQuizConfig(db, &config, logger, redis)
	if err != nil {
		return err
	}

	// CustomAuthenticated (not kratos-only) so guests can host public quizzes too.
	// GetOrActivateSession still enforces admin_id == userId, so a guest can only
	// arrange a session they themselves created.
	v1.Get(fmt.Sprintf("/socket/admin/arrange/:%s", constants.SessionIDParam), middleware.CheckSessionId, middleware.CustomAuthenticated, websocket.New(quizSocketController.Arrange))
	v1.Get(fmt.Sprintf("/socket/join/:%s", constants.QuizSessionInvitationCode), middleware.CheckSessionCode, middleware.CustomAuthenticated, websocket.New(quizSocketController.Join))
	v1.Post("/quiz/answer", middleware.Authenticated, middleware.CustomAuthenticated, quizSocketController.SetAnswer)
	v1.Get("/quiz/terminate", middleware.Authenticated, quizSocketController.Terminate)

	return nil
}

func setupQuizController(v1 fiber.Router, db *goqu.Database, logger *zap.Logger, middleware middlewares.Middleware, config config.AppConfig) error {
	quizController, err := controller.InitQuizController(db, logger, &config)
	if err != nil {
		return err
	}

	admin := v1.Group("/admin")
	admin.Use(middleware.KratosAuthenticated)

	// Public quizzes listing — must be registered BEFORE the auth-guarded /quizzes group
	// so it is not protected by KratosAuthenticated.
	v1.Get("/quizzes/public", quizController.GetPublicQuizzes)

	// Hosting a public quiz is open to guests as well as registered users, so it uses
	// Authenticated (guest JWT or kratos) instead of the kratos-only /quizzes group.
	// The handler itself enforces that the quiz is actually public.
	v1.Post(fmt.Sprintf("/quizzes/:%s/public_session", constants.QuizId), middleware.Authenticated, quizController.GeneratePublicSession)

	quizzes := v1.Group("/quizzes")
	quizzes.Use(middleware.KratosAuthenticated)

	quizzes.Post(fmt.Sprintf("/:%s/demo_session", constants.QuizId), quizController.GenerateDemoSession)
	quizzes.Post(fmt.Sprintf("/:%s/upload", constants.QuizTitle), middleware.ValidateCsv, middleware.KratosAuthenticated, quizController.CreateQuizByCsv)
	quizzes.Post("/", quizController.CreateQuiz)
	quizzes.Get("/", quizController.GetAdminUploadedQuizzes)
	quizzes.Put(fmt.Sprintf("/:%s/settings", constants.QuizId), middleware.QuizPermission, middleware.VerifyQuizEditAccess, quizController.UpdateQuizSettings)
	quizzes.Delete(fmt.Sprintf("/:%s", constants.QuizId), middleware.QuizPermission, middleware.VerifyQuizEditAccess, quizController.DeleteQuizById)

	report := admin.Group("/reports")
	report.Get("/list", quizController.ListQuizzesAnalysis)
	report.Get(fmt.Sprintf("/:%s/analysis", constants.ActiveQuizId), middleware.KratosAuthenticated, quizController.GetQuizAnalysis)
	return nil
}

func setupQuizCategoryController(v1 fiber.Router, db *goqu.Database, logger *zap.Logger, middleware middlewares.Middleware, config config.AppConfig) error {
	quizCategoryController, err := controller.InitQuizCategoryController(db, logger, &config)
	if err != nil {
		return err
	}

	// Category listing is public — the homepage groups public quizzes by category.
	// Registered BEFORE the auth-guarded /categories group so it stays open.
	v1.Get("/categories", quizCategoryController.ListCategories)

	// Managing categories requires Kratos auth; on top of that, each handler
	// enforces the public-quiz admin email allowlist (same gate as publishing).
	categories := v1.Group("/categories")
	categories.Use(middleware.KratosAuthenticated)

	categories.Post("/", quizCategoryController.CreateCategory)
	categories.Put(fmt.Sprintf("/:%s", constants.CategoryId), quizCategoryController.UpdateCategory)
	categories.Delete(fmt.Sprintf("/:%s", constants.CategoryId), quizCategoryController.DeleteCategory)

	return nil
}

func setupQuestionController(v1 fiber.Router, db *goqu.Database, logger *zap.Logger, middleware middlewares.Middleware, config config.AppConfig) error {
	questionController, err := controller.InitQuestionController(db, logger, &config)
	if err != nil {
		return err
	}

	questionRouter := v1.Group(fmt.Sprintf("/quizzes/:%s/questions", constants.QuizId))
	questionRouter.Use(middleware.KratosAuthenticated, middleware.QuizPermission)

	questionRouter.Get("/", questionController.ListQuestionsWithAnswerByQuizId)
	questionRouter.Post("/", middleware.VerifyQuizEditAccess, questionController.CreateQuestion)
	questionRouter.Post("/upload", middleware.VerifyQuizEditAccess, middleware.ValidateCsv, questionController.ImportQuestionsByCsv)
	questionRouter.Get(fmt.Sprintf("/:%s", constants.QuestionId), middleware.VerifyQuizEditAccess, questionController.GetQuestionById)
	questionRouter.Put(fmt.Sprintf("/:%s", constants.QuestionId), middleware.VerifyQuizEditAccess, questionController.UpdateQuestionById)
	questionRouter.Delete(fmt.Sprintf("/:%s", constants.QuestionId), middleware.VerifyQuizEditAccess, questionController.DeleteQuestionById)

	return nil
}

// final score board controller setup
func setUpFinalScoreBoardController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware) error {
	finalScoreBoardController, err := controller.NewFinalScoreBoardController(goqu, logger)
	if err != nil {
		return err
	}

	finalScoreBoardControllerAdmin, err := controller.NewFinalScoreBoardAdminController(goqu, logger)
	if err != nil {
		return err
	}

	finalScore := v1.Group("/final_score")
	finalScore.Get("/user", finalScoreBoardController.GetScore)
	finalScore.Get("/admin", middlewares.KratosAuthenticated, finalScoreBoardControllerAdmin.GetScoreForAdmin)

	return nil
}

func setUpAnalyticsBoardController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, config config.AppConfig, middlewares middlewares.Middleware) error {
	analyticsBoardUserController, err := controller.NewAnalyticsBoardUserController(goqu, logger, &config)
	if err != nil {
		return err
	}

	analyticsBoardAdminController, err := controller.NewAnalyticsBoardAdminController(goqu, logger, &config)
	if err != nil {
		return err
	}

	analyticsBoard := v1.Group("/analytics_board")
	analyticsBoard.Get("/user", analyticsBoardUserController.GetAnalyticsForUser)
	analyticsBoard.Get("/admin", middlewares.KratosAuthenticated, analyticsBoardAdminController.GetAnalyticsForAdmin)

	return nil
}

func setupUserPlayedQuizeController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, config config.AppConfig) error {
	userPlayedQuizeController, err := controller.NewUserPlayedQuizeController(goqu, logger, &config)
	if err != nil {
		return err
	}

	userRouter := v1.Group("/user_played_quizes")
	userRouter.Get("/", middlewares.KratosAuthenticated, userPlayedQuizeController.ListUserPlayedQuizes)
	userRouter.Get(fmt.Sprintf("/:%s", constants.UserPlayedQuizId), userPlayedQuizeController.ListUserPlayedQuizesWithQuestionById)
	userRouter.Post(fmt.Sprintf("/:%s", constants.QuizSessionInvitationCode), middlewares.Authenticated, userPlayedQuizeController.PlayedQuizValidation)
	return nil
}

func setupSharedQuizzesController(v1 fiber.Router, goqu *goqu.Database, logger *zap.Logger, middlewares middlewares.Middleware, config config.AppConfig) error {
	sharedQuizzesController, err := controller.NewSharedQuizzesController(goqu, logger, &config)
	if err != nil {
		return err
	}

	sharedQuizzesRouter := v1.Group("/shared_quizzes")
	sharedQuizzesRouter.Use(middlewares.KratosAuthenticated)

	sharedQuizzesRouter.Get("/", sharedQuizzesController.ListSharedQuizzes)
	sharedQuizzesRouter.Post(fmt.Sprintf("/:%s", constants.QuizId), middlewares.QuizPermission, middlewares.VerifyQuizShareAccess, sharedQuizzesController.ShareQuiz)
	sharedQuizzesRouter.Get(fmt.Sprintf("/:%s", constants.QuizId), middlewares.QuizPermission, middlewares.VerifyQuizShareAccess, sharedQuizzesController.ListQuizAuthorizedUsers)
	sharedQuizzesRouter.Put(fmt.Sprintf("/:%s", constants.QuizId), middlewares.QuizPermission, middlewares.VerifyQuizShareAccess, sharedQuizzesController.UpdateUserPermissionOfQuiz)
	sharedQuizzesRouter.Delete(fmt.Sprintf("/:%s", constants.QuizId), middlewares.QuizPermission, middlewares.VerifyQuizShareAccess, sharedQuizzesController.DeleteUserPermissionOfQuiz)
	return nil
}
