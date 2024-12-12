package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"

	controller "github.com/Improwised/quizz-app/api/controllers/api/v1"

	"github.com/fasthttp/websocket"
)

var invitationCode interface{}

func initializeQuiz(t *testing.T) {
	createQuiz(t, quizTitle)
	generateDemoSession(t)
}

func createWebSocket(t *testing.T, url string, header http.Header) *websocket.Conn {
	cfg := config.LoadTestEnv()
	conn, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s/api/v1/socket/%s", cfg.Port, url), header)
	assert.NoError(t, err, "WebSocket connection failed")
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode, "Unexpected status code")
	assert.Equal(t, "websocket", resp.Header.Get("Upgrade"), "Expected WebSocket upgrade")

	return conn
}

// admin websocket connection
func adminWSConnection(t *testing.T) *websocket.Conn {
	var authenticatedHeader = http.Header{}
	var cookie = &http.Cookie{Name: client.Cookies[0].Name, Value: client.Cookies[0].Value}
	authenticatedHeader.Set("Cookie", cookie.String())

	initializeQuiz(t)
	url := "admin/arrange/" + sessionId
	conn := createWebSocket(t, url, authenticatedHeader)
	return conn
}

func playerWSConnection(t *testing.T, code interface{}) *websocket.Conn {
	var guestUserHeader = http.Header{}
	var cookie = &http.Cookie{Name: userclient.Cookies[0].Name, Value: userclient.Cookies[0].Value}
	guestUserHeader.Set("Cookie", cookie.String())

	url := fmt.Sprintf("join/%v?username=%s", code, guestUserName)
	conn := createWebSocket(t, url, guestUserHeader)
	return conn
}

func validateMapConversion(t *testing.T, input interface{}) map[string]interface{} {
	response, ok := input.(map[string]interface{})
	assert.True(t, ok, "type assertion should true")
	return response
}

func readSocketMessage(t *testing.T, conn *websocket.Conn) map[string]interface{} {
	var msg structs.SocketResponseFormat
	err := conn.ReadJSON(&msg)
	assert.NoError(t, err, "Failed to read WebSocket message")
	response := validateMapConversion(t, msg.Data)
	return response
}

func sendSocketMessage(t *testing.T, conn *websocket.Conn, message controller.QuizReceiveResponse) {
	err := conn.WriteJSON(message)
	assert.NoError(t, err, "Failed to send WebSocket message")
}

func setupUserPlayedQuizzes(t *testing.T, code interface{}) {
	res, err := userclient.
		R().
		EnableTrace().
		Post(fmt.Sprintf("/api/v1/user_played_quizes/%v", code))

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode(), res)

	var responsePlayedQuizValidation utils.ResponsePlayedQuizValidation
	err = json.Unmarshal(res.Body(), &responsePlayedQuizValidation.Body)
	assert.Nil(t, err, "Error in parsing response body")

	assert.Equal(t, sessionId, responsePlayedQuizValidation.Body.Data.SessionId, res)
	userPlayedQuizId = responsePlayedQuizValidation.Body.Data.UserPlayedQuizId
}

func setupQuizWithPlayer(t *testing.T) (*websocket.Conn, *websocket.Conn) {
	adminConn := adminWSConnection(t)
	msg := readSocketMessage(t, adminConn)
	response := validateMapConversion(t, msg["data"])
	code := validateMapConversion(t, response["data"])
	invitationCode = code["code"]
	setupUserPlayedQuizzes(t, invitationCode)
	playerConn := playerWSConnection(t, invitationCode)
	return adminConn, playerConn
}

func TestArrange(t *testing.T) {
	t.Run("test invalid session id", func(t *testing.T) {

		var authenticatedHeader = http.Header{}
		var cookie = &http.Cookie{Name: client.Cookies[0].Name, Value: client.Cookies[0].Value}
		authenticatedHeader.Set("Cookie", cookie.String())
		url := "admin/arrange/4efdfd74-1451-4aa8-806a-67aa95157067"

		conn := createWebSocket(t, url, authenticatedHeader)
		defer conn.Close()

		msg := readSocketMessage(t, conn)
		response := validateMapConversion(t, msg["data"])
		assert.Equal(t, constants.EventActivateSession, msg["event"])
		assert.Equal(t, constants.UnknownError, response["data"])
	})

	t.Run("test valid session id", func(t *testing.T) {

		adminWebSocket := adminWSConnection(t)
		defer adminWebSocket.Close()
		msg := readSocketMessage(t, adminWebSocket)
		assert.Equal(t, constants.EventSendInvitationCode, msg["event"])

	})
}

func TestHandleCodeGeneration(t *testing.T) {
	t.Run("test pinging", func(t *testing.T) {

		adminConn := adminWSConnection(t)
		defer adminConn.Close()

		var sendMessage = controller.QuizReceiveResponse{
			Event: constants.EventPing,
			Data:  "ping",
		}
		readSocketMessage(t, adminConn) //read invitation code
		sendSocketMessage(t, adminConn, sendMessage)
		msg := readSocketMessage(t, adminConn)

		assert.Equal(t, constants.EventPong, msg["event"]) //expect pong
	})

	t.Run("start quiz without any user join", func(t *testing.T) {
		adminConn := adminWSConnection(t)
		defer adminConn.Close()

		readSocketMessage(t, adminConn) //read invitation code

		var sendMessage = controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
			Data:  "",
		}
		sendSocketMessage(t, adminConn, sendMessage)

		msg := readSocketMessage(t, adminConn)
		response := validateMapConversion(t, msg["data"])
		assert.Equal(t, constants.NoPlayerFound, response["data"]) //expect no player found
	})

	t.Run("start quiz with joined user", func(t *testing.T) {

		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()

		readSocketMessage(t, adminConn) // read user join data
		var sendMessage = controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
			Data:  "",
		}
		sendSocketMessage(t, adminConn, sendMessage)

		msg := readSocketMessage(t, adminConn)
		response := validateMapConversion(t, msg["data"])

		assert.Equal(t, constants.EventStartCount5, msg["event"])
		assert.Equal(t, constants.ActionCounter, response["action"])
		assert.Equal(t, constants.Question, response["component"])
	})
}

func TestQuestionAndScoreHandler(t *testing.T) {
	t.Run("test send question", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()

		readSocketMessage(t, adminConn) // read user join data
		var sendMessage = controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
			Data:  "",
		}
		sendSocketMessage(t, adminConn, sendMessage)

		readSocketMessage(t, adminConn)        // read 5 sec counter
		msg := readSocketMessage(t, adminConn) // read question
		response := validateMapConversion(t, msg["data"])
		question := validateMapConversion(t, response["data"])

		assert.Equal(t, constants.EventSendQuestion, msg["event"])
		assert.Equal(t, "Which city is known as the Eternal City?", question["question"])
	})

	t.Run("test skip question", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()

		readSocketMessage(t, adminConn) // read user join data
		var sendMessage = controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
			Data:  "",
		}
		sendSocketMessage(t, adminConn, sendMessage)

		readSocketMessage(t, adminConn) // read 5 second counter
		readSocketMessage(t, adminConn) // read question

		sendMessage.Event = constants.EventSkipAsked
		sendSocketMessage(t, adminConn, sendMessage)

		msg := readSocketMessage(t, adminConn)
		response := validateMapConversion(t, msg["data"])

		assert.Equal(t, constants.EventSkipAsked, msg["event"])
		assert.Equal(t, constants.WarnSkip, response["data"])

		sendMessage.Event = constants.EventForceSkip
		sendSocketMessage(t, adminConn, sendMessage)

		msg = readSocketMessage(t, adminConn)
		response = validateMapConversion(t, msg["data"])

		assert.Equal(t, constants.EventShowScore, msg["event"])
		assert.Equal(t, constants.ActionShowScore, response["action"])
	})
}

func TestJoin(t *testing.T) {
	t.Run("test pinging from player connection", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()
		sendMessage := controller.QuizReceiveResponse{
			Event: constants.EventPing,
		}
		sendSocketMessage(t, playerConn, sendMessage)

		msg := readSocketMessage(t, playerConn)
		assert.Equal(t, constants.EventPong, msg["event"])
	})

	t.Run("test admin join quiz as a player", func(t *testing.T) {
		adminConn := adminWSConnection(t)
		defer adminConn.Close()
		msg := readSocketMessage(t, adminConn)
		response := validateMapConversion(t, msg["data"])
		code := validateMapConversion(t, response["data"])

		var authenticatedHeader = http.Header{}
		var cookie = &http.Cookie{Name: client.Cookies[0].Name, Value: client.Cookies[0].Value}
		authenticatedHeader.Set("Cookie", cookie.String())
		url := fmt.Sprintf("join/%v?username=%s", code["code"], "John")
		playerConn := createWebSocket(t, url, authenticatedHeader)
		defer playerConn.Close()

		msg = readSocketMessage(t, playerConn)
		assert.Equal(t, constants.EventRedirectToAdmin, msg["event"])

		response = validateMapConversion(t, msg["data"])
		assert.Equal(t, constants.ActionCurrentUserIsAdmin, response["action"])
	})

	t.Run("test user join in running quiz", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()

		readSocketMessage(t, playerConn)
		sendMessage := controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
		}
		sendSocketMessage(t, adminConn, sendMessage)
		readSocketMessage(t, playerConn) // read 5 second counter
		readSocketMessage(t, playerConn) // read question

		// close the current player socket and create new connection in running quiz
		playerConn.Close()
		time.Sleep(2 * time.Second) // wait for 2 seconds before join again for test duration
		playerConn = playerWSConnection(t, invitationCode)
		defer playerConn.Close()

		msg := readSocketMessage(t, playerConn)
		response := validateMapConversion(t, msg["data"])
		question := validateMapConversion(t, response["data"])
		assert.Equal(t, constants.EventSendQuestion, msg["event"])
		assert.LessOrEqual(t, 58.00, question["duration"].(float64))
		assert.Equal(t, constants.ActionSendQuestion, response["action"])
	})
}

func TestSetAnswer(t *testing.T) {

	t.Run("submit answer with all correct details", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()

		readSocketMessage(t, playerConn)
		sendMessage := controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
		}
		sendSocketMessage(t, adminConn, sendMessage)
		readSocketMessage(t, playerConn)        // read 5 second counter
		msg := readSocketMessage(t, playerConn) // read question
		assert.Equal(t, constants.EventSendQuestion, msg["event"])
		response := validateMapConversion(t, msg["data"])

		question := validateMapConversion(t, response["data"])
		req := map[string]interface{}{
			"id":            question["id"].(string),
			"keys":          []int{1},
			"response_time": 1000,
		}
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("user_played_quiz", userPlayedQuizId).
			SetQueryParam("session_id", sessionId).
			SetBody(req).
			Post("/api/v1/quiz/answer")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusAccepted, res.StatusCode())
	})

	t.Run("submint answer with incorrect question id", func(t *testing.T) {
		adminConn, playerConn := setupQuizWithPlayer(t)
		defer adminConn.Close()
		defer playerConn.Close()

		readSocketMessage(t, playerConn)
		sendMessage := controller.QuizReceiveResponse{
			Event: constants.EventStartQuiz,
		}
		sendSocketMessage(t, adminConn, sendMessage)
		readSocketMessage(t, playerConn)        // read 5 second counter
		msg := readSocketMessage(t, playerConn) // read question
		assert.Equal(t, constants.EventSendQuestion, msg["event"])

		req := map[string]interface{}{
			"id":            uuid.New().String(),
			"keys":          []int{1},
			"response_time": 1000,
		}
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam("user_played_quiz", userPlayedQuizId).
			SetQueryParam("session_id", sessionId).
			SetBody(req).
			Post("/api/v1/quiz/answer")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode())
	})
}
