package v1

import (
	"fmt"
	"log"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizHelper "github.com/Improwised/quizz-app/api/helpers/quiz"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type QuizSendResponse struct {
	Component string `json:"component"` // simulates a page
	Action    string `json:"action"`    // action is a description of an event
	Data      any    `json:"data"`      // optional data
}

type QuizReceiveResponse struct {
	Component string `json:"component"` // simulates a page
	Event     string `json:"event"`     // event
	Data      any    `json:"data"`      // optional data
}

func CreateQuickUser(db *goqu.Database, logger *zap.Logger, userObj models.User, retrying bool, emailValidation bool) (models.User, error) {
	userModel, err := models.InitUserModel(db)

	if err != nil {
		return userObj, err
	}

	if emailValidation {
		isUnique, err := userModel.IsUniqueEmail(userObj.Email)

		if err != nil {
			return userObj, fmt.Errorf("someError occurred during register user %v", err)
		}

		if !isUnique {
			return userObj, fmt.Errorf("email is already registered")
		}
	}

	userSvc := services.NewUserService(&userModel)

	copyUserObj, err := userSvc.RegisterUser(userObj, events.NewEventBus(logger))

	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {

			if !(retrying && pqErr.Constraint == constants.UserUkey) {
				return userObj, fmt.Errorf("username (%s) already registered", userObj.Username)
			}

			copyUserObj.Password = userObj.Password

			copyUserObj.Username = quizUtilsHelper.GenerateNewStringHavingSuffixName(userObj.Username, 5, 12)

			copyUserObj, err = userSvc.RegisterUser(copyUserObj, events.NewEventBus(logger))

			if err != nil {
				return userObj, fmt.Errorf("SomeError during register admin with new username %s", userObj.Username)
			}

		}

	}

	userObj.ID = copyUserObj.ID
	userObj.Username = copyUserObj.Username

	return userObj, err
}

type quizSocketController struct {
	db        *models.QuizModel
	appConfig *config.AppConfig
	logger    *zap.Logger
	helpers   *quizHelper.HelperStructs
}

func InitQuizConfig(db *goqu.Database, appConfig *config.AppConfig, logger *zap.Logger, helpers *quizHelper.HelperStructs) *quizSocketController {
	return &quizSocketController{models.InitQuizModel(db), appConfig, logger, helpers}
}

type ping struct {
	payloadType int
	payload     []byte
	is_close    bool
}

func pingResponse(pingRequest ping, c *websocket.Conn) {
	log.Println("write:", pingRequest.payloadType, string(pingRequest.payload))
	var response []byte = pingRequest.payload

	err := c.WriteMessage(pingRequest.payloadType, response)

	if pingRequest.is_close || err != nil {
		time.Sleep(1 * time.Second)
		log.Println("write:", pingRequest.payloadType, string(pingRequest.payload))
		c.Close()
	}

}

func (*quizSocketController) Ping(c *websocket.Conn) {

	defer func() {
		log.Println("user disconnected")
	}()

	var (
		initial_time   time.Time = time.Now()
		pingObj        ping      = ping{websocket.TextMessage, []byte("You can send \"hi\", \"bye\" message in socket, other messages are send back"), false}
		defaultChannel chan ping = make(chan ping)
	)

	pingResponse(pingObj, c)

	go func(pingObj ping, defaultChannel chan ping) {

		var (
			payload     []byte
			payloadType int
			err         error
		)

		for {
			if payloadType, payload, err = c.ReadMessage(); err != nil {
				pingObj.payloadType = websocket.CloseAbnormalClosure
				pingObj.payload = []byte(err.Error())
				pingObj.is_close = true
				defaultChannel <- pingObj
				return
			} else if string(payload) == "bye" || payloadType == websocket.CloseMessage {
				pingObj.payloadType = payloadType
				pingObj.payload = []byte("nice to meet you")
				pingObj.is_close = true
				defaultChannel <- pingObj
				return
			} else if string(payload) == "hi" {
				pingObj.payloadType = websocket.TextMessage
				pingObj.payload = []byte("hello: " + time.Now().Format("2006-01-02 15:04:05 MST"))
				pingObj.is_close = false
				defaultChannel <- pingObj
			}
			pingObj.payloadType = payloadType
			pingObj.is_close = false
			pingObj.payload = payload
			defaultChannel <- pingObj
		}
	}(pingObj, defaultChannel)

	tick := time.NewTicker(5 * time.Second)

	for {
		select {
		case pingObj = <-defaultChannel:
			pingResponse(pingObj, c)

			if pingObj.is_close {
				return
			}
			initial_time = time.Now()

		case <-tick.C:
			if initial_time.Add(15 * time.Second).After(time.Now()) {
				pingObj.payloadType = websocket.TextMessage
				pingObj.payload = []byte("Knock knock")
				pingResponse(pingObj, c)
			} else {
				pingObj.payloadType = websocket.CloseMessage
				pingObj.payload = []byte("We will meet soon!!!")
				pingResponse(pingObj, c)
				return
			}
		}
	}
}

func (qc *quizSocketController) Join(c *websocket.Conn) {

	defer c.Close()

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	// check for middleware error
	if !c.Locals(constants.MiddlewarePass).(bool) {
		response.Data = c.Locals(constants.MiddlewareError).(error).Error()
		err := utils.JSONFailWs(c, constants.EventAuthentication, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error in middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
		}
		return
	}

	code := c.Locals(constants.QuizSessionCode).(string)

	session, err := qc.helpers.QuizSessionModel.GetSessionByCode(code)

	if err != nil {
		response.Action = constants.ActionJoinQuiz
		response.Data = constants.ErrCodeNotFound
		err = utils.JSONFailWs(c, constants.EventJoinQuiz, response)

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error in session get in code: %s event, %s action, %s code", constants.EventJoinQuiz, response.Action, code), zap.Error(err))
		}
		return
	}

	if !session.IsActive || session.ActivatedTo.Valid {
		response.Data = constants.ErrCodeNotFound
		err = utils.JSONFailWs(c, constants.EventJoinQuiz, response)

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error check session active: %s event, %s action, %s code", constants.EventJoinQuiz, response.Action, code), zap.Error(err))
		}
		return
	}

	response.Data = constants.QuizStartsSoon
	err = utils.JSONSuccessWs(c, constants.EventJoinQuiz, response)

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventJoinQuiz, response.Action), zap.Error(err))
	}

}

func (qc *quizSocketController) Arrange(c *websocket.Conn) {

	isConnected := true
	defer func() {
		isConnected = false
		time.Sleep(1 * time.Second)
		c.Close()
	}()

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	// checks for any middleware errors
	if !c.Locals(constants.MiddlewarePass).(bool) {
		response.Data = c.Locals(constants.MiddlewareError).(string)
		err := utils.JSONErrorWs(c, constants.EventAuthentication, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
		}
		time.Sleep(1 * time.Second)
		return
	}

	sessionId := c.Locals(constants.SessionIDPram).(string)

	user := c.Locals(constants.ContextUser).(models.User)

	// activate session
	session, err := ActivateAndGetSession(c, qc.helpers, qc.logger, sessionId, user.ID)

	if err != nil {
		return
	}

	// is isQuestionActive true -> quiz started
	isCodeSent := session.IsQuestionActive.Valid

	if !isCodeSent {
		// handle Waiting page
		for {

			if !isConnected {
				break
			}

			// if code not sent then sent it
			if !isCodeSent {
				// send code to client
				handleCodeSend(c, response, qc.logger, session.Code.Int32)
				isCodeSent = true
			}

			// once code sent receive start signal
			if isCodeSent {
				isBreak := handleStartQuiz(c, qc.logger, &isConnected, response.Action)
				if isBreak {
					break
				}
			}
		}
	}

	response.Component = constants.Question
	questions, err := qc.helpers.QuizModel.GetSharedQuestions(int(session.Code.Int32))

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.Code), zap.Error(err))
	}

	response.Data = questions
	// handle question page iteration
	err = utils.JSONSuccessWs(c, constants.EventStartQuiz, response)

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error success: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.Code), zap.Error(err))
	}
}

// Activate session

func ActivateAndGetSession(c *websocket.Conn, helpers *quizHelper.HelperStructs, logger *zap.Logger, sessionId string, userId string) (models.QuizSession, error) {

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	session, err := helpers.QuizSessionModel.GetActiveSession(sessionId, userId)

	if err != nil {
		if err.Error() == constants.Unauthenticated {
			response.Action = constants.ActionSessionValidation
			response.Data = constants.Unauthorized
			err = utils.JSONFailWs(c, constants.EventAuthorization, response)
			if err != nil {
				logger.Error(fmt.Sprintf("socket error authentication host: %s event, %s action", constants.EventAuthorization, response.Action), zap.Error(err))
			}
			return session, err
		}

		response.Action = constants.ActionSessionActivation
		response.Data = constants.UnknownError
		err = utils.JSONErrorWs(c, constants.EventActivateSession, response)
		if err != nil {
			logger.Error(fmt.Sprintf("socket error get or activate session: %s event, %s action", constants.EventActivateSession, response.Action), zap.Error(err))
		}
		return session, err
	}

	c.Locals(constants.SessionObj, session)

	return session, nil
}

// handle waiting page

func handleCodeSend(c *websocket.Conn, response QuizSendResponse, logger *zap.Logger, code int32) bool {

	// send code to client
	response.Action = constants.ActionSessionActivation
	response.Data = map[string]int{"code": int(code)}

	err := utils.JSONSuccessWs(c, constants.EventSendCode, response)

	if err != nil {
		logger.Error(fmt.Sprintf("socket error sent code: %s event, %s action", constants.EventSendCode, response.Action), zap.Error(err))
	}

	return true
}

func handleStartQuiz(c *websocket.Conn, logger *zap.Logger, isConnected *bool, action string) bool {
	message := QuizReceiveResponse{}
	err := c.ReadJSON(&message)

	if err != nil {
		logger.Error(fmt.Sprintf("socket error start event handling: %s event, %s action", constants.EventStartQuiz, action), zap.Error(err))
		*isConnected = false
		return true
	}

	if message.Event == constants.EventStartQuiz {
		return true
	}

	return false
}
