package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
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

		pqErr, ok := quizUtilsHelper.ConvertType[*pq.Error](err)

		if !ok {
			return userObj, fmt.Errorf("SomeError during register admin with new username %s", userObj.Username)
		}

		if pqErr.Code == "23505" {

			if !(retrying && pqErr.Constraint == constants.UserUkey) {
				return userObj, fmt.Errorf("username (%s) already registered", userObj.Username)
			}

			copyUserObj.Password = userObj.Password

			copyUserObj.Username = quizUtilsHelper.GenerateNewStringHavingSuffixName(userObj.Username, 5, 12)

			copyUserObj, err = userSvc.RegisterUser(copyUserObj, events.NewEventBus(logger))

		}

		if err != nil {
			return userObj, fmt.Errorf("SomeError during register admin with new username %s", userObj.Username)
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
	helpers   *quizHelper.HelperGroup
}

func InitQuizConfig(db *goqu.Database, appConfig *config.AppConfig, logger *zap.Logger, helpers *quizHelper.HelperGroup) *quizSocketController {
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

	defer func() {
		c.Close()
	}()

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	// check for middleware error
	if c.Locals(constants.MiddlewareError) != nil {
		// handle error type
		errorInfo, ok := quizUtilsHelper.ConvertType[error](c.Locals(constants.MiddlewareError))
		if ok {
			response.Data = errorInfo.Error()
			err := utils.JSONErrorWs(c, constants.EventAuthentication, response)
			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error in middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
			}
			return
		}

		// handle string type
		errorString := quizUtilsHelper.GetString(c.Locals(constants.MiddlewareError))
		if errorString != "<nil>" {
			response.Data = errorString
			err := utils.JSONFailWs(c, constants.EventAuthentication, response)
			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error in middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
			}
			return
		}

	}

	invitationCode := quizUtilsHelper.GetString(c.Locals(constants.QuizSessionInvitationCode))

	session, ok := quizUtilsHelper.ConvertType[models.ActiveQuiz](c.Locals(constants.ActiveQuizObj))

	if !ok {
		response.Action = constants.ActionSessionValidation
		err := utils.JSONErrorWs(c, constants.UnknownError, response)

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error session type change: %s event, %s action, %s code", constants.EventSessionValidation, response.Action, invitationCode), zap.Error(err))
		}
		return
	}

	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	// is user is a host of current quiz
	if userId == session.AdminID {
		response.Action = constants.ActionCurrentUserIsAdmin
		response.Data = map[string]string{"sessionId": session.ID.String()}
		err := utils.JSONSuccessWs(c, constants.EventRedirectToAdmin, response)

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket redirect current user is admin: %s event, %s action, %s code", constants.EventRedirectToAdmin, response.Action, invitationCode), zap.Error(err))
		}
		return
	}

	response.Action = constants.QuizQuestionStatus
	if session.CurrentQuestion.Valid {
		response.Data = constants.NextQuestionWillServeSoon
	} else {
		response.Data = constants.QuizStartsSoon
	}

	err := utils.JSONSuccessWs(c, constants.EventJoinQuiz, response)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventJoinQuiz, response.Action), zap.Error(err))
	}

	// userPlayedQuizId := quizUtilsHelper.GetString(c.Locals(constants.CurrentUserQuiz))
	handleQuestion(c, qc, session, response)

}

func handleQuestion(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, response QuizSendResponse) {
	pubsub := qc.helpers.PubSubModel.Client.Subscribe(qc.helpers.PubSubModel.Ctx, session.ID.String())
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.helpers.PubSubModel.Ctx, session.ID.String())
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()

	for msg := range ch {

		message := map[string]any{}
		err := json.Unmarshal([]byte(msg.Payload), &message)

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventJoinQuiz, response.Action), zap.Error(err))
		}

		event := quizUtilsHelper.GetString(message["event"])
		err = utils.JSONSuccessWs(c, event, message["response"])
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", event, response.Action), zap.Error(err))
		}

		if message["event"] == constants.EventTerminateQuiz {
			break
		}
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
	if c.Locals(constants.MiddlewareError) != nil {
		response.Data = quizUtilsHelper.GetString(c.Locals(constants.MiddlewareError))
		err := utils.JSONErrorWs(c, constants.EventAuthentication, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
		}
		time.Sleep(1 * time.Second)
		return
	}

	sessionId := quizUtilsHelper.GetString(c.Locals(constants.SessionIDParam))

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))

	if !ok {
		qc.logger.Error("socket user-model type conversion")
		err := utils.JSONFailWs(c, constants.EventSessionValidation, constants.UnknownError)
		if err != nil {
			qc.logger.Error("socket user-model type conversion")
		}
		return
	}

	// activate session
	session, err := ActivateAndGetSession(c, qc.helpers, qc.logger, sessionId, user.ID)

	if err != nil {
		qc.logger.Error("get active session", zap.Error(err))
		err := utils.JSONFailWs(c, constants.EventSessionValidation, constants.UnknownError)
		if err != nil {
			qc.logger.Error("get active session", zap.Error(err))
		}
		return
	}

	// handle code sharing with admin
	handleCodeGeneration(c, qc, session, &isConnected, &response)

	// if connection lost during waiting of start event
	if !(isConnected) {
		qc.logger.Error("user disconnected:")
		return
	}

	// question component handling
	response.Component = constants.Question
	questions, lastQuestionDeliveryTime, err := qc.helpers.QuizModel.GetSharedQuestions(int(session.InvitationCode.Int32))
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))
		return
	}

	var isFirst bool = lastQuestionDeliveryTime.Valid
	var wg sync.WaitGroup
	response.Component = constants.Question
	for _, question := range questions {
		wg.Add(1)
		if isFirst {
			isFirst = false
			sendQuestion(c, qc, &wg, &response, session, question, lastQuestionDeliveryTime)
		} else {
			sendQuestion(c, qc, &wg, &response, session, question, sql.NullTime{})
		}
		wg.Wait()

		err := utils.JSONSuccessWs(c, constants.EventNextQuestionAsk, response)

		if err != nil {
			qc.logger.Error("socket error during asking for next question", zap.Error(err))
		}

		if !handleNextQuestion(c, qc.logger, &isConnected) {
			response.Component = constants.Loading
			response.Data = constants.AdminDisconnected
			shareEvenWithUser(c, qc, &response, constants.AdminDisconnected, sessionId, int(session.InvitationCode.Int32), false)
			break
		}
	}

	// termination of quiz
	if session.ActivatedFrom.Valid && isConnected {
		terminateQuiz(c, qc, &response, session)
	}
}

func handleCodeGeneration(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, isConnected *bool, response *QuizSendResponse) {
	// is isQuestionActive true -> quiz started
	isInvitationCodeSent := session.CurrentQuestion.Valid

	if !isInvitationCodeSent {
		// handle Waiting page
		for {

			if !(*isConnected) {
				break
			}

			// if code not sent then sent it
			if !isInvitationCodeSent {
				// send code to client
				handleInvitationCodeSend(c, response, qc.logger, session.InvitationCode.Int32)
				isInvitationCodeSent = true
			}

			// once code sent receive start signal
			if isInvitationCodeSent {
				isBreak := handleStartQuiz(c, qc.logger, isConnected, response.Action)
				subscriberCount := qc.helpers.PubSubModel.Client.PubSubNumSub(qc.helpers.PubSubModel.Ctx, session.ID.String()).Val()[session.ID.String()]
				if subscriberCount != 0 && isBreak {
					break
				} else {
					response.Data = constants.NoPlayerFound
					err := utils.JSONFailWs(c, constants.EventStartQuiz, response)
					if err != nil {
						qc.logger.Error(fmt.Sprintf("socket error middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
					}
				}

			}
		}
	}
}

func sendQuestion(c *websocket.Conn, qc *quizSocketController, wg *sync.WaitGroup, response *QuizSendResponse, session models.ActiveQuiz, question models.Question, lastQuestionTimeStamp sql.NullTime) {
	// start counter if not any question running
	if !lastQuestionTimeStamp.Valid {
		response.Component = constants.Question
		response.Action = constants.ActionCounter
		response.Data = map[string]int{"counter": constants.Counter, "count": constants.Count}
		shareEvenWithUser(c, qc, response, constants.EventStartCount5, session.ID.String(), int(session.InvitationCode.Int32), false)
		time.Sleep(time.Duration(constants.Counter) * time.Second)

		// update question status to activate
		err := qc.helpers.QuizModel.UpdateCurrentQuestion(session.ID, question.ID, true)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
		}
	}

	// question sent
	response.Action = constants.ActionSendQuestion
	responseData := map[string]any{
		"id":            question.ID,
		"no":            question.OrderNumber,
		"duration":      question.DurationInSeconds,
		"question_time": lastQuestionTimeStamp.Time,
		"question":      question.Question,
		"options":       question.Options,
	}

	if !lastQuestionTimeStamp.Valid {
		responseData["question_time"] = ""
		response.Data = responseData
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), false)
	} else {
		responseData["duration"] = question.DurationInSeconds - int(time.Since(lastQuestionTimeStamp.Time).Seconds())
		response.Data = responseData
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), true)
	}

	wgForQuestion := &sync.WaitGroup{}
	wgForQuestion.Add(1)
	if !lastQuestionTimeStamp.Valid { // new question
		go handleAnswerSubmission(question.DurationInSeconds, wgForQuestion)
	} else { // current question
		duration := question.DurationInSeconds - int(time.Since(lastQuestionTimeStamp.Time).Seconds())
		go handleAnswerSubmission(duration, wgForQuestion)
	}
	wgForQuestion.Wait()

	// update current status to deactivate
	err := qc.helpers.QuizModel.UpdateCurrentQuestion(session.ID, question.ID, false)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
	}

	// ideal time for score page showing or admin event
	// question sent
	response.Component = constants.Score
	response.Action = constants.ActionShowScore
	response.Data = map[string]any{
		"no": question.OrderNumber,
	}
	shareEvenWithUser(c, qc, response, constants.EventShowScore, session.ID.String(), int(session.InvitationCode.Int32), false)
	time.Sleep(5 * time.Second)
	wg.Done()
}

func handleAnswerSubmission(duration int, wg *sync.WaitGroup) {
	time.Sleep(time.Duration(duration) * time.Second)
	wg.Done()
}

// Activate session

func ActivateAndGetSession(c *websocket.Conn, helpers *quizHelper.HelperGroup, logger *zap.Logger, sessionId string, userId string) (models.ActiveQuiz, error) {

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	session, err := helpers.ActiveQuizModel.GetOrActivateSession(sessionId, userId)

	if err != nil {
		if err.Error() == constants.Unauthenticated {
			response.Action = constants.ActionSessionValidation
			response.Data = constants.Unauthorized
			err = utils.JSONFailWs(c, constants.EventAuthorization, response)
			if err != nil {
				logger.Error(fmt.Sprintf("socket error authentication host: %s event, %s action", constants.EventAuthorization, response.Action), zap.Error(err))
			}
			return session, err
		} else if err.Error() == constants.ErrSessionWasCompleted {
			response.Action = constants.ActionSessionActivation
			response.Data = constants.ErrSessionWasCompleted
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

	c.Locals(constants.ActiveQuizObj, session)

	return session, nil
}

// handle waiting page

func handleInvitationCodeSend(c *websocket.Conn, response *QuizSendResponse, logger *zap.Logger, invitationCode int32) bool {

	// send code to client
	response.Action = constants.ActionSessionActivation
	response.Data = map[string]int{"code": int(invitationCode)}

	err := utils.JSONSuccessWs(c, constants.EventSendInvitationCode, response)

	if err != nil {
		logger.Error(fmt.Sprintf("socket error sent code: %s event, %s action", constants.EventSendInvitationCode, response.Action), zap.Error(err))
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

func handleNextQuestion(c *websocket.Conn, logger *zap.Logger, isConnected *bool) bool {
	message := QuizReceiveResponse{}
	err := c.ReadJSON(&message)

	if err != nil {
		logger.Error(fmt.Sprintf("socket error during asking for next question: %s event", constants.EventNextQuestionAsk), zap.Error(err))
		*isConnected = false
		return false
	}

	if message.Event == constants.EventNextQuestionAsk {
		return true
	}

	return false
}

func shareEvenWithUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, event string, sessionId string, invitationCode int, prevent bool) {
	payload := map[string]any{"event": event, "response": response}
	data, err := json.Marshal(payload)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error marshal redis payload: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
	}

	if !prevent {
		// send event to user
		err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, sessionId, data).Err()

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action %v code", constants.EventPublishQuestion, response.Action, invitationCode), zap.Error(err))
		}
	}

	// send event to admin
	err = utils.JSONSuccessWs(c, event, response)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error sending event: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
	}
}

func terminateQuiz(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz) {
	response.Component = constants.Score
	response.Data = constants.ActionTerminateQuiz
	shareEvenWithUser(c, qc, response, constants.EventTerminateQuiz, session.ID.String(), int(session.InvitationCode.Int32), false)

	err := qc.helpers.ActiveQuizModel.Deactivate(session.ID)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))
		return
	}
}
