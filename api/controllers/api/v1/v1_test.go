package v1_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Improwised/quizz-app/api/cli"
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/services"
	goqu "github.com/doug-martin/goqu/v9"
	resty "github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

var client *resty.Client = nil
var userclient *resty.Client = nil
var db *goqu.Database = nil
var sessionCookie string
var userId string
var guestUserName = "testcaseuser"
var guestUserAvatarName = "Chase"

func TestMain(m *testing.M) {
	err := os.Chdir("../../../")
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.LoadTestEnv()
	logger, err := logger.NewRootLogger(true, true)
	if err != nil {
		log.Fatal(err)
	}

	db, err = database.Connect(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("http://%s", cfg.Port)

	client = resty.New().SetBaseURL(url)
	userclient = resty.New().SetBaseURL(url)

	// Setup to follow redirects
	client.SetRedirectPolicy(resty.FlexibleRedirectPolicy(10))

	cmd := cli.GetAPICommandDef(cfg, logger)

	// execute migration in sqlite
	migrationCmd := cli.GetMigrationCommandDef(cfg)
	migrationCmd.SetArgs([]string{"up"})
	err = migrationCmd.Execute()
	if err != nil {
		logger.Fatal("error while execute migration", zap.Error(err))
	}

	go func() {
		err = cmd.Execute()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for the server to be ready
	serverRunning := false
	for count := 0; count < 100; count += 1 {
		client = client.SetTimeout(time.Second * 2)
		res, err := client.R().EnableTrace().Get("/api/healthz")
		if err == nil {
			log.Println("received status code", res.StatusCode())
		}
		if err == nil && res.StatusCode() == http.StatusOK {
			serverRunning = true
			break
		}
		time.Sleep(time.Second * 2)
	}

	if !serverRunning {
		log.Fatal("program exit due to server is not running...")
	}

	client = client.SetTimeout(time.Second * 10)
	log.Println("server is running...")

	// Create and login user to get the session cookie
	identityID, err := registerUser(cfg.Kratos.BaseUrl)
	if err != nil {
		logger.Fatal("setup user and login failed", zap.Error(err))
	}

	userModel, err := models.InitUserModel(db, logger)
	if err != nil {
		logger.Fatal("init user model failed", zap.Error(err))

	}

	user, err := userModel.GetUserByKratosID(identityID)
	if err != nil {
		logger.Fatal("error while get user by kratosId", zap.Error(err))

	}
	userId = user.ID

	err = createGuestUser(guestUserName, guestUserAvatarName)
	if err != nil {
		logger.Fatal("failed to create guest user", zap.Error(err))
	}

	exitCode := m.Run()

	if err := deleteIdentity(logger, cfg, user.ID, identityID); err != nil {
		logger.Fatal("failed to delete identity", zap.Error(err))
	}

	if err := deleteGuestUser(guestUserName); err != nil {
		logger.Fatal("failed to delete identity", zap.Error(err))
	}

	os.Exit(exitCode)
}

func createGuestUser(username, avatarName string) error {
	userRes, err := userclient.
		R().
		EnableTrace().
		Post(fmt.Sprintf("/api/v1/user/%s?avatar_name=%s", username, avatarName))

	if err != nil || userRes.StatusCode() != http.StatusOK {
		return fmt.Errorf("failed to create guest user: %v", err)
	}

	// Extract the `user` cookie
	var userCookie string
	cookies := userRes.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "user" {
			userCookie = cookie.Value
			break
		}
	}

	// set user cookie
	userclient = userclient.SetCookie(&http.Cookie{
		Name:  "user",
		Value: userCookie,
	})

	return nil
}

func deleteGuestUser(username string) error {
	_, err := db.Exec(fmt.Sprintf("delete from users where username='%s'", username))
	return err
}

func registerUser(kratosUrl string) (string, error) {

	// Initiate the registration flow
	regFlowResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0").
		SetResult(&map[string]interface{}{}).
		Get(kratosUrl + "/self-service/registration/browser")

	if err != nil {
		return "", fmt.Errorf("failed to initiate registration flow: %v", err)
	}

	regFlowData := *(regFlowResp.Result().(*map[string]interface{}))
	flowID, ok := regFlowData["id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract registration flow ID from response")
	}

	// Extract CSRF token cookie
	var csrfTokenCookie *http.Cookie
	for _, cookie := range regFlowResp.Cookies() {
		if strings.HasPrefix(cookie.Name, "csrf_token_") {
			csrfTokenCookie = cookie
			break
		}
	}

	if csrfTokenCookie == nil {
		return "", fmt.Errorf("csrf_token cookie not found in response")
	}

	// Extract CSRF token from the response UI nodes
	var csrfToken string
	if ui, ok := regFlowData["ui"].(map[string]interface{}); ok {
		if nodes, ok := ui["nodes"].([]interface{}); ok {
			for _, node := range nodes {
				if nodeMap, ok := node.(map[string]interface{}); ok {
					if attributes, ok := nodeMap["attributes"].(map[string]interface{}); ok {
						if name, ok := attributes["name"].(string); ok && name == "csrf_token" {
							csrfToken, _ = attributes["value"].(string)
							break
						}
					}
				}
			}
		}
	}
	if csrfToken == "" {
		return "", fmt.Errorf("CSRF token not found in response")
	}

	// Submit registration details
	regResp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0").
		SetResult(&map[string]interface{}{}).
		SetBody(map[string]interface{}{
			"traits.email":      "exampletestuser@example.com",
			"password":          "dockerisnotpodman",
			"traits.name.first": "John",
			"traits.name.last":  "Doe",
			"method":            "password",
			"csrf_token":        csrfToken,
		}).
		SetCookie(csrfTokenCookie).
		Post(fmt.Sprintf("%s/self-service/registration?flow=%s", kratosUrl, flowID))

	if err != nil {
		return "", fmt.Errorf("failed to submit registration details: %v", err)
	}

	if regResp.StatusCode() != http.StatusOK && regResp.StatusCode() != http.StatusSeeOther {
		return "", fmt.Errorf("registration failed: %s", regResp.String())
	}

	// Extract the `ory_kratos_session` cookie
	cookies := regResp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "ory_kratos_session" {
			sessionCookie = cookie.Value
			break
		}
	}

	if sessionCookie == "" {
		return "", fmt.Errorf("ory_kratos_session cookie not found in response")
	}

	// Set the session cookie for the client
	client = client.SetCookie(&http.Cookie{
		Name:  "ory_kratos_session",
		Value: sessionCookie,
	})

	// Extract the identity ID from the registration response
	regRespData := *(regResp.Result().(*map[string]interface{}))
	var identityID string
	if identity, ok := regRespData["identity"].(map[string]interface{}); ok {
		if id, ok := identity["id"].(string); ok {
			identityID = id
		}
	}

	// insert data into the user table
	client.
		R().
		EnableTrace().
		Get("/api/v1/kratos/auth")

	fmt.Println("Registration successful:", identityID)
	return identityID, nil
}

func deleteIdentity(logger *zap.Logger, config config.AppConfig, userID, identityID string) error {

	userSvc, err := services.NewUserService(db, logger, config)
	if err != nil {
		return err
	}

	// delete whole user Data (clean up)
	err = userSvc.DeleteUserDataById(userID, identityID)
	if err != nil {
		return err
	}

	return nil
}
