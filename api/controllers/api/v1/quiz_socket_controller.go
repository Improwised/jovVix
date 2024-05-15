package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
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

func (qc *quizSocketController) Join(c *websocket.Conn) {
	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	userName := c.Query(constants.UserName, "")
	var quizResponse QuizReceiveResponse

	// check for middleware error
	if c.Locals(constants.MiddlewareError) != nil {
		handleMiddleWareError(c, qc, response)
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

	// check user web socket connection is close or not
	go func() {
		defer func() {
			removeUser(qc, session.ID.String(), userName)
			c.Close()
		}()

		for {
			_, p, err := c.ReadMessage()

			if err != nil {
				break
			}
			err = json.Unmarshal([]byte(p), &quizResponse)
			if err != nil {
				break
			}

			if quizResponse.Event == "websocket_close" {
				break
			}
		}
	}()

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

	// when user join at that time publish userName to admin
	publishUserOnJoin(c, qc, response, userName, session.ID.String())

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

func publishUserOnJoin(c *websocket.Conn, qc *quizSocketController, quizResponse QuizSendResponse, userName string, sessionId string) {
	// store data to redis in form of slice
	usersData := []string{}
	var jsonData []byte

	response := quizResponse

	// check weather current session id has an any user to show if not then set. if present then get and add new user to it
	exists, err := qc.helpers.PubSubModel.Client.Exists(qc.helpers.PubSubModel.Ctx, sessionId).Result()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}
	if exists == 0 {
		usersData = append(usersData, userName)
		// Serialize slice to JSON
		jsonData, err = json.Marshal(usersData)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
		}

	} else {
		// get data from redis
		users, err := qc.helpers.PubSubModel.Client.Get(qc.helpers.PubSubModel.Ctx, sessionId).Result()
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
		}
		err = json.Unmarshal([]byte(users), &usersData)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
		}
		usersData = lo.Reverse(usersData)
		usersData = append(usersData, userName)
		usersData = lo.Reverse(usersData)
		jsonData, err = json.Marshal(usersData)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))

		}

	}

	// if quiz is still not start then publish join user data to admin and refresh the page
	err = qc.helpers.PubSubModel.Client.Set(qc.helpers.PubSubModel.Ctx, sessionId, jsonData, time.Minute*100).Err()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}

	err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.EventUserJoined, jsonData).Err()

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}
}

func handleMiddleWareError(c *websocket.Conn, qc *quizSocketController, response QuizSendResponse) {
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

	defer func() {
		isConnected = false
		time.Sleep(1 * time.Second)
		c.Close()
		onRefreshHandleConnectedUser(qc, session.ID.String())
	}()

	// handle code sharing with admin
	handleCodeGeneration(c, qc, session, &isConnected, &response)

	// if connection lost during waiting of start event
	if !(isConnected) {
		response.Component = constants.Loading
		response.Data = constants.AdminDisconnected
		shareEvenWithUser(c, qc, &response, constants.AdminDisconnected, sessionId, int(session.InvitationCode.Int32), constants.ToUser)

		qc.logger.Error("admin disconnected")
		return
	}

	// question and score handler
	questionAndScoreHandler(c, qc, &response, session, &isConnected)
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
				go handleConnectedUser(c, qc)

				isBreak := handleStartQuiz(c, qc.logger, isConnected, response.Action)
				subscriberCount := qc.helpers.PubSubModel.Client.PubSubNumSub(qc.helpers.PubSubModel.Ctx, session.ID.String()).Val()[session.ID.String()]
				if subscriberCount != 0 && isBreak {
					// quiz is start publish for admin to stop looking for user
					err := qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.EventStartQuizByAdmin, constants.EventStartQuizByAdmin).Err()
					if err != nil {
						qc.logger.Error("errro while start quiz", zap.Error(err))
					}
					break

				} else {
					// quiz is start publish for admin to stop looking for user becuse no player found
					err := qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.StartQuizByAdminNoPlayerFound, constants.StartQuizByAdminNoPlayerFound).Err()
					if err != nil {
						qc.logger.Error("errro while start quiz but no player found", zap.Error(err))
					}
					response.Data = constants.NoPlayerFound

					err = utils.JSONFailWs(c, constants.EventStartQuiz, response)
					if err != nil {
						qc.logger.Error(fmt.Sprintf("socket error middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
					}
				}

			}
		}
	}
}

func questionAndScoreHandler(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, isConnected *bool) {
	// get questions/remaining question
	response.Component = constants.Question
	questions, lastQuestionDeliveryTime, err := qc.helpers.QuizModel.GetSharedQuestions(int(session.InvitationCode.Int32))
	if err != nil {
		response.Action = constants.ErrInGettingQuestion
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))
		err := utils.JSONFailWs(c, constants.EventSendQuestion, response)
		if err != nil {
			qc.logger.Error("error during get remaining question", zap.Error(err))
		}
		return
	}

	var wg sync.WaitGroup
	chanNextEvent := make(chan bool)
	chanSkipEvent := make(chan bool)
	chanSkipTimer := make(chan bool)
	var isQuizEnd bool = false

	go listenAllEvents(c, qc, response, session, chanNextEvent, chanSkipEvent, chanSkipTimer, isQuizEnd)

	// handle question
	var isFirst bool = lastQuestionDeliveryTime.Valid
	response.Component = constants.Question
	for _, question := range questions {
		wg.Add(1)
		if isFirst { // handle running question
			isFirst = false
			sendSingleQuestion(c, qc, &wg, response, session, question, lastQuestionDeliveryTime, chanSkipEvent, chanSkipTimer)
		} else { // handle new question
			sendSingleQuestion(c, qc, &wg, response, session, question, sql.NullTime{}, chanSkipEvent, chanSkipTimer)
		}
		err := utils.JSONSuccessWs(c, constants.EventNextQuestionAsked, response)

		if err != nil {
			qc.logger.Error("socket error during asking for next question", zap.Error(err))
		}

		// handle next question
		if <-chanNextEvent {
			continue
		}

		wg.Wait()
	}

	// termination of quiz
	if session.ActivatedFrom.Valid && *isConnected {
		terminateQuiz(c, qc, response, session)
		// isQuizEnd = false
	}
}

func listenAllEvents(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, chanNextEvent chan bool, chanSkipEvent chan bool, chanSkipTimer chan bool, isQuizEnd bool) {
	for {
		message := QuizReceiveResponse{}
		err := c.ReadJSON(&message)

		if err != nil {
			qc.logger.Error("error in receiving message from question", zap.Error(err))
			// isConnected = false
			break
		}

		switch message.Event {
		case constants.EventSkipAsked:
			chanSkipEvent <- false
		case constants.EventForceSkip:
			chanSkipEvent <- true
		case constants.EventNextQuestionAsked:
			chanNextEvent <- true
		case constants.EventSkipTimer:
			chanSkipTimer <- true
		}
	}

	// handle connection lost during quiz
	if !isQuizEnd {
		response.Component = constants.Loading
		response.Data = constants.AdminDisconnected
		shareEvenWithUser(c, qc, response, constants.AdminDisconnected, session.ID.String(), int(session.InvitationCode.Int32), constants.ToUser)
	}
}

func sendSingleQuestion(c *websocket.Conn, qc *quizSocketController, wg *sync.WaitGroup, response *QuizSendResponse, session models.ActiveQuiz, question models.Question, lastQuestionTimeStamp sql.NullTime, chanSkipEvent chan bool, chanSkipTimer chan bool) {
	// start counter if not any question running
	if !lastQuestionTimeStamp.Valid {
		response.Component = constants.Question
		response.Action = constants.ActionCounter
		response.Data = map[string]int{"counter": constants.Counter, "count": constants.Count}
		shareEvenWithUser(c, qc, response, constants.EventStartCount5, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)
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

	if !lastQuestionTimeStamp.Valid { // handling new question
		responseData["question_time"] = ""
		response.Data = responseData
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)
	} else { // handling running question
		responseData["duration"] = question.DurationInSeconds - int(time.Since(lastQuestionTimeStamp.Time).Seconds())
		response.Data = responseData
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin)
	}

	wgForQuestion := &sync.WaitGroup{}
	wgForQuestion.Add(1)
	var duration int
	if !lastQuestionTimeStamp.Valid { // new question
		duration = question.DurationInSeconds
	} else { // handle running question
		duration = question.DurationInSeconds - int(time.Since(lastQuestionTimeStamp.Time).Seconds())
		if duration < 0 {
			duration = 1
		}
	}
	go handleAnswerSubmission(c, qc, session, question.ID, duration, wgForQuestion, chanSkipEvent, response)
	wgForQuestion.Wait()

	// update current status to deactivate
	err := qc.helpers.QuizModel.UpdateCurrentQuestion(session.ID, question.ID, false)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
	}

	// score-board rendering
	response.Component = constants.Score
	response.Action = constants.ActionShowScore
	userRankBoard, err := qc.helpers.UserPlayedQuizModel.GetRank(session.ID, question.ID)

	if err != nil {
		qc.logger.Error("error during get userRankBoard", zap.Error(err))
	}

	response.Data = map[string]any{
		"rankList": userRankBoard,
		"question": question.Question,
		"answers":  question.Answers,
		"options":  question.Options,
		"duration": 20,
	}

	shareEvenWithUser(c, qc, response, constants.EventShowScore, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)

	wgForSkipTimer := &sync.WaitGroup{}
	wgForSkipTimer.Add(1)

	// skip 20 sec timer
	go handleSkipTimer(wgForSkipTimer, chanSkipTimer)
	wgForSkipTimer.Wait()

	wg.Done()
}

func handleSkipTimer(wg *sync.WaitGroup, chanSkipTimer chan bool) {
	defer wg.Done()

	isTimeout := time.NewTicker(time.Duration(20) * time.Second)

	for {
		select {
		case <-isTimeout.C:
			return
		case isSkip := <-chanSkipTimer: // skip 20 sec if admin clicks on skip button
			if isSkip {
				return
			}
		}
	}
}

func handleAnswerSubmission(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, questionId uuid.UUID, duration int, wg *sync.WaitGroup, chanSkipEvent chan bool, response *QuizSendResponse) {
	defer wg.Done()

	isTimeout := time.NewTicker(time.Duration(duration) * time.Second)

	for {
		select {
		case <-isTimeout.C:
			return
		case isForce := <-chanSkipEvent:
			if isForce {
				return
			} else {
				ok, err := qc.helpers.QuizModel.IsAllAnswerGathered(session.ID, questionId)
				if err != nil {
					qc.logger.Error("error during listening skip event", zap.Error(err))
				}
				if ok {
					return
				} else { // send warning if all participant not given answer
					response.Data = constants.WarnSkip
					shareEvenWithUser(c, qc, response, constants.EventSkipAsked, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin)
				}
			}
		}
	}

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

// start quiz by message event from admin
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

func shareEvenWithUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, event string, sessionId string, invitationCode int, sentToWhom int) {
	payload := map[string]any{"event": event, "response": response}
	data, err := json.Marshal(payload)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error marshal redis payload: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
	}

	if sentToWhom == constants.ToUser || sentToWhom == constants.ToAll {
		// send event to user
		err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, sessionId, data).Err()

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action %v code", constants.EventPublishQuestion, response.Action, invitationCode), zap.Error(err))
		}
	}

	if sentToWhom == constants.ToAdmin || sentToWhom == constants.ToAll {
		// send event to admin
		err = utils.JSONSuccessWs(c, event, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error sending event: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
		}
	}
}

func terminateQuiz(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz) {
	response.Component = constants.Score
	response.Data = constants.ActionTerminateQuiz
	shareEvenWithUser(c, qc, response, constants.EventTerminateQuiz, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)

	err := qc.helpers.ActiveQuizModel.Deactivate(session.ID)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))
		return
	}

	// here logic of publishing data of user to admin that terminate quiz so no need to listen for joining users
	err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.EventTerminateQuiz, constants.EventTerminateQuiz).Err()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error while terminationg quiz %s", constants.ActionTerminateQuiz), zap.Error(err))
		return
	}

}

// when user connect at that time send data to admin
func handleConnectedUser(c *websocket.Conn, qc *quizSocketController) {
	response := QuizSendResponse{}
	response.Action = "send user join data"
	response.Component = constants.Waiting
	var mutex sync.Mutex
	defer mutex.Unlock()

	mutex.Lock()

	pubsub := qc.helpers.PubSubModel.Client.Subscribe(qc.helpers.PubSubModel.Ctx, constants.EventUserJoined, constants.EventTerminateQuiz, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.helpers.PubSubModel.Ctx, constants.EventUserJoined, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()
	userData := []string{}

	for msg := range ch {
		response.Data = msg.Payload

		if response.Data == constants.EventStartQuizByAdmin || response.Data == constants.EventTerminateQuiz || response.Data == constants.StartQuizByAdminNoPlayerFound {
			break
		}

		err := json.Unmarshal([]byte(msg.Payload), &userData)
		if err != nil {
			qc.logger.Error("error while unmarshaling data inside handleconnect user ", zap.Error(err))

			break
		}

		response.Data = userData
		err = utils.JSONSuccessWs(c, constants.EventSendInvitationCode, response)
		if err != nil {
			qc.logger.Error("error while sending user data ", zap.Error(err))
		}
	}

}

// remove user from redis after quiz is over and web socket connection close
func removeUser(qc *quizSocketController, topic string, userName string) {
	var mutex sync.Mutex
	defer mutex.Unlock()
	mutex.Lock()
	// get data from redis
	usersData := []string{}
	result, err := qc.helpers.PubSubModel.Client.Get(qc.helpers.PubSubModel.Ctx, topic).Result()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error %s", constants.EventSendInvitationCode), zap.Error(err))
	}

	err = json.Unmarshal([]byte(result), &usersData)
	if err != nil {
		qc.logger.Error("error while unmarshaling data inside removeuser ", zap.Error(err))
	}
	found := lo.IndexOf(usersData, userName)
	if len(usersData) == 1 {
		usersData = []string{}
	} else {
		usersData = append(usersData[:found], usersData[found+1:]...)
	}
	jsonData, err := json.Marshal(usersData)
	if err != nil {
		qc.logger.Error("error while marshling data inside removeuser ", zap.Error(err))
	}

	err = qc.helpers.PubSubModel.Client.Set(qc.helpers.PubSubModel.Ctx, topic, jsonData, time.Minute*100).Err()
	if err != nil {
		qc.logger.Error("error while seting data to redis  inside removeuser ", zap.Error(err))
	}

	err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.EventUserJoined, jsonData).Err()

	if err != nil {
		qc.logger.Error("error while publishig data to redis  inside removeuser ", zap.Error(err))
	}

}

// when admin refresh page at that time send redis conncted user data to admin
func onRefreshHandleConnectedUser(qc *quizSocketController, topic string) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()
	usersData, err := qc.helpers.PubSubModel.Client.Get(qc.helpers.PubSubModel.Ctx, topic).Result()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error %s", constants.EventSendInvitationCode), zap.Error(err))
	}

	err = qc.helpers.PubSubModel.Client.Publish(qc.helpers.PubSubModel.Ctx, constants.EventUserJoined, usersData).Err()

	if err != nil {
		qc.logger.Error("error while publishig data to redis  inside removeuser ", zap.Error(err))
	}
}
