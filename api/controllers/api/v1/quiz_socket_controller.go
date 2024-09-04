package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	quizUtilsHelper "github.com/Improwised/quizz-app/api/helpers/utils"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/redis"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
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

type UserInfo struct {
	UserId   string
	UserName string
	IsAlive  bool
}

type quizSocketController struct {
	activeQuizModel         *models.ActiveQuizModel
	quizModel               *models.QuizModel
	userPlayedQuizModel     *models.UserPlayedQuizModel
	questionModel           *models.QuestionModel
	appConfig               *config.AppConfig
	logger                  *zap.Logger
	redis                   *redis.RedisPubSub
	answersSubmittedByUsers chan models.User
	mu                      sync.Mutex
}

func InitQuizConfig(db *goqu.Database, appConfig *config.AppConfig, logger *zap.Logger, redis *redis.RedisPubSub, answersSubmittedByUsers chan models.User) *quizSocketController {

	activeQuizModel := models.InitActiveQuizModel(db, logger)
	quizModel := models.InitQuizModel(db)
	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)

	return &quizSocketController{
		activeQuizModel:         activeQuizModel,
		quizModel:               quizModel,
		userPlayedQuizModel:     userPlayedQuizModel,
		questionModel:           questionModel,
		appConfig:               appConfig,
		logger:                  logger,
		redis:                   redis,
		answersSubmittedByUsers: answersSubmittedByUsers,
	}
}

// for user Join
func (qc *quizSocketController) Join(c *websocket.Conn) {
	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		qc.logger.Error("error while fetching user context from connection")
	}
	var quizResponse QuizReceiveResponse

	invitationCode := quizUtilsHelper.GetString(c.Locals(constants.QuizSessionInvitationCode))

	session, err := qc.activeQuizModel.GetSessionByCode(invitationCode)
	if err != nil {
		if err == sql.ErrNoRows {
			qc.logger.Error(constants.ErrInvitationCodeNotFound, zap.Error(err))
			c.Close()
			return
		}
		qc.logger.Error("error in invitation code", zap.Error(err))
		c.Close()
		return
	}

	userId := quizUtilsHelper.GetString(c.Locals(constants.ContextUid))

	defer func() {
		c.Close()
		qc.logger.Info("connection closed by user")
	}()

	// check user web socket connection is close or not
	go func() {
		for {
			_, p, err := c.ReadMessage()

			if err != nil {
				// if error occurs, change the connection alive status to false
				updateUserData(qc, userId, session.ID.String(), false)
				break
			}
			err = json.Unmarshal([]byte(p), &quizResponse)
			if err != nil {
				updateUserData(qc, userId, session.ID.String(), false)
				break
			}

			if quizResponse.Event == "websocket_close" {
				updateUserData(qc, userId, session.ID.String(), false)
				qc.logger.Info("connection close request is send by the user - " + user.Username)
				break
			}

			if quizResponse.Event == constants.EventPing {
				// if gets ping, change the connection alive status to true
				updateUserData(qc, userId, session.ID.String(), true)
			}
		}
	}()

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
	publishUserOnJoin(qc, response, user.FirstName, userId, session.ID.String())
	response.Action = constants.QuizQuestionStatus

	onConnectHandleUser(c, qc, &response, session)

	// userPlayedQuizId := quizUtilsHelper.GetString(c.Locals(constants.CurrentUserQuiz))
	handleQuestion(c, qc, session, response)
}

func publishUserOnJoin(qc *quizSocketController, quizResponse QuizSendResponse, userName string, userId string, sessionId string) {
	// store data to redis in form of slice
	var usersData []UserInfo
	var jsonData []byte

	response := quizResponse

	// check weather current session id has an any user to show if not then set. if present then get and add new user to it
	exists, err := qc.redis.PubSubModel.Client.Exists(qc.redis.PubSubModel.Ctx, sessionId).Result()
	if err != nil {
		qc.logger.Error("error while checking if there is any user in redis for the session in publishUserOnJoin", zap.Error(err))
	}
	if exists == 0 {
		newUser := UserInfo{UserId: userId, UserName: userName, IsAlive: true}
		usersData = append(usersData, newUser)
		// Serialize slice to JSON
		jsonData, err = json.Marshal(usersData)
		if err != nil {
			qc.logger.Error("error while marshalling data into json in publishUserOnJoin when there is no data in redis", zap.Error(err))
		}

	} else {
		// get data from redis
		users, err := qc.redis.PubSubModel.Client.Get(qc.redis.PubSubModel.Ctx, sessionId).Result()
		if err != nil {
			qc.logger.Error("error while fetching data from redis in publishUserOnJoin", zap.Error(err))
		}
		err = json.Unmarshal([]byte(users), &usersData)
		if err != nil {
			qc.logger.Error("error while unmarshaling redis in publishUserOnJoin", zap.Error(err))
		}
		for _, data := range usersData {
			if userId == data.UserId {
				qc.logger.Error(fmt.Sprintf("User %s already exist in redis", userName))
				return
			}
		}
		newUser := UserInfo{UserId: userId, UserName: userName, IsAlive: true}
		usersData = append(usersData, newUser)
		jsonData, err = json.Marshal(usersData)
		if err != nil {
			qc.logger.Error("error while marshaling data into json in publishUserOnJoin", zap.Error(err))

		}

	}

	// if quiz is still not start then publish join user data to admin and refresh the page
	err = qc.redis.PubSubModel.Client.Set(qc.redis.PubSubModel.Ctx, sessionId, jsonData, time.Minute*100).Err()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}

	// remove data with isAlive false before publishing
	publishData := filterPublishUsers(qc, usersData, "publishUserOnJoin")

	err = qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, constants.EventUserJoined, publishData).Err()

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}
}

func handleQuestion(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, response QuizSendResponse) {
	pubsub := qc.redis.PubSubModel.Client.Subscribe(qc.redis.PubSubModel.Ctx, session.ID.String())
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.redis.PubSubModel.Ctx, session.ID.String())
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

func onConnectHandleUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz) {
	if session.CurrentQuestion.Valid {

		questionID, err := uuid.Parse(session.CurrentQuestion.String)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("\nquestionID is not being parsed from the current question id of this session and that current question id is - %v\n", session.CurrentQuestion), zap.Error(err))
		}

		currentQuestion, err := qc.questionModel.GetCurrentQuestion(questionID)
		if err != nil {
			qc.logger.Error("unable to get the current question and the question id was "+session.CurrentQuestion.String, zap.Error(err))
		}

		response.Action = constants.ActionSendQuestion
		duration := currentQuestion.DurationInSeconds - int(time.Since(session.QuestionDeliveryTime.Time).Seconds())
		if duration < 0 {
			return
		}
		responseData := map[string]any{
			"id":            currentQuestion.ID,
			"no":            currentQuestion.OrderNumber,
			"duration":      duration,
			"question_time": session.QuestionDeliveryTime.Time,
			"question":      currentQuestion.Question,
			"options":       currentQuestion.Options,
		}
		response.Data = responseData
		response.Component = constants.Question

		err = utils.JSONSuccessWs(c, constants.EventSendQuestion, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error send current question on connect: %s event, %s action", constants.EventSendQuestion, response.Action), zap.Error(err))
		}
	} else {
		response.Data = constants.QuizStartsSoon
		err := utils.JSONSuccessWs(c, constants.EventJoinQuiz, response)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventJoinQuiz, response.Action), zap.Error(err))
		}
	}
}

// function to update user IsAlive status
func updateUserData(qc *quizSocketController, userId string, sessionId string, isAlive bool) {
	// Fetch data from Redis
	users, err := qc.redis.PubSubModel.Client.Get(qc.redis.PubSubModel.Ctx, sessionId).Result()
	if err != nil {
		qc.logger.Error("error fetching data from redis in updateUserData", zap.Error(err))
		return
	}

	var usersData []UserInfo
	if err := json.Unmarshal([]byte(users), &usersData); err != nil {
		qc.logger.Error("error unmarshaling redis data in updateUserData", zap.Error(err))
		return
	}

	var update bool
	var updatedUserData []UserInfo
	for _, data := range usersData {
		if data.UserId == userId {
			if data.IsAlive == isAlive {
				return
			}
			data.IsAlive = isAlive
			update = true
		}
		if data.IsAlive {
			updatedUserData = append(updatedUserData, data)
		}
	}

	if update {
		jsonData, err := json.Marshal(updatedUserData)
		if err != nil {
			qc.logger.Error("error marshaling updated data in updateUserData", zap.Error(err))
			return
		}

		// Update Redis
		if err := qc.redis.PubSubModel.Client.Set(qc.redis.PubSubModel.Ctx, sessionId, jsonData, time.Minute*100).Err(); err != nil {
			qc.logger.Error("error updating data in redis in updateUserData", zap.Error(err))
			return
		}

		qc.logger.Debug(fmt.Sprintf("IsAlive status updated for user %s (%s)", userId, userId))

		// Publish updated user data
		publishData := filterPublishUsers(qc, updatedUserData, "updateUserData")
		if err := qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, constants.EventUserJoined, publishData).Err(); err != nil {
			qc.logger.Error("error publishing data in updateUserData", zap.Error(err))
		}
	}
}

// filter user data and publish only alive user to the channel
func filterPublishUsers(qc *quizSocketController, usersData []UserInfo, functionName string) (publishData []byte) {

	qc.mu.Lock()
	defer qc.mu.Unlock()

	// Filter out elements where IsAlive is false
	var filteredData []UserInfo
	for _, data := range usersData {
		if data.IsAlive {
			filteredData = append(filteredData, data)
		}
	}

	// store new data into publishData
	publishData, err := json.Marshal(filteredData)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("error while marshaling data into filterPublishUsers, called from %s", functionName), zap.Error(err))
	}

	return publishData
}

// for admin join
func (qc *quizSocketController) Arrange(c *websocket.Conn) {

	isConnected := true

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
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
	session, err := ActivateAndGetSession(c, qc.activeQuizModel, qc.logger, sessionId, user.ID)

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
		qc.logger.Debug("deactivating quiz - " + session.ID.String())
		err := qc.activeQuizModel.Deactivate(session.ID)
		if err != nil {
			qc.logger.Error("error while deactivating quiz", zap.Error(err))
		}
		c.Close()
		qc.logger.Info("connection closed by admin")
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

func ActivateAndGetSession(c *websocket.Conn, activeQuizModel *models.ActiveQuizModel, logger *zap.Logger, sessionId string, userId string) (models.ActiveQuiz, error) {

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	session, err := activeQuizModel.GetOrActivateSession(sessionId, userId)

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
		logger.Debug("unknown error was triggered from ActivateAndGetSession")
		err = utils.JSONErrorWs(c, constants.EventActivateSession, response)
		if err != nil {
			logger.Error(fmt.Sprintf("socket error get or activate session: %s event, %s action", constants.EventActivateSession, response.Action), zap.Error(err))
		}
		return session, err
	}

	c.Locals(constants.ActiveQuizObj, session)

	return session, nil
}

func handleCodeGeneration(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, isConnected *bool, response *QuizSendResponse) {
	// is isQuestionActive true -> quiz started
	isInvitationCodeSent := session.CurrentQuestion.Valid

	if !isInvitationCodeSent {
		// handle Waiting page
		for {

			qc.mu.Lock()
			if !(*isConnected) {
				qc.mu.Unlock()
				break
			}

			// if code not sent then sent it
			if !isInvitationCodeSent {
				// send code to client
				handleInvitationCodeSend(c, response, qc.logger, session.InvitationCode.Int32)
				isInvitationCodeSent = true
				go handleConnectedUser(c, qc)

			}
			qc.mu.Unlock()

			// once code sent receive start signal
			if isInvitationCodeSent {
				isBreak := handleStartQuiz(c, qc.logger, isConnected, response.Action, &qc.mu)
				users, err := qc.redis.PubSubModel.Client.Get(qc.redis.PubSubModel.Ctx, session.ID.String()).Result()
				if err != nil {
					qc.logger.Error("error while fetching data from redis inside updateUserData", zap.Error(err))
				}

				var usersData []UserInfo
				err = json.Unmarshal([]byte(users), &usersData)
				if err != nil {
					qc.logger.Error("error while unmarshaling redis inside updateUserData", zap.Error(err))
				}

				if isBreak == constants.EventPing {
					continue
				} else if len(usersData) != 0 && isBreak == constants.EventStartQuiz {

					// quiz is start publish for admin to stop looking for user
					err := qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, constants.EventStartQuizByAdmin, constants.EventStartQuizByAdmin).Err()
					if err != nil {
						qc.logger.Error("error while start quiz", zap.Error(err))
					}
					break
				} else {
					// quiz is start publish for admin to stop looking for user becuse no player found
					err := qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, constants.StartQuizByAdminNoPlayerFound, constants.StartQuizByAdminNoPlayerFound).Err()
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

// when user connect at that time send data to admin
func handleConnectedUser(c *websocket.Conn, qc *quizSocketController) {
	response := QuizSendResponse{}
	response.Action = "send user join data"
	response.Component = constants.Waiting

	pubsub := qc.redis.PubSubModel.Client.Subscribe(qc.redis.PubSubModel.Ctx, constants.EventUserJoined, constants.EventTerminateQuiz, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.redis.PubSubModel.Ctx, constants.EventUserJoined, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()
	usersData := []UserInfo{}

	for msg := range ch {
		response.Data = msg.Payload

		if response.Data == constants.EventStartQuizByAdmin || response.Data == constants.EventTerminateQuiz || response.Data == constants.StartQuizByAdminNoPlayerFound {
			break
		}

		err := json.Unmarshal([]byte(msg.Payload), &usersData)
		if err != nil {
			qc.logger.Error("error while unmarshaling data inside handleconnectedUser ", zap.Error(err))

			break
		}

		usersName := []string{}
		for _, data := range usersData {
			usersName = append(usersName, data.UserName)
		}
		response.Data = usersName
		err = utils.JSONSuccessWs(c, constants.EventSendInvitationCode, response) // sending the user data to the admin
		if err != nil {
			qc.logger.Error("error while sending user data ", zap.Error(err))
		}
	}
}

// start quiz by message event from admin
func handleStartQuiz(c *websocket.Conn, logger *zap.Logger, isConnected *bool, action string, mu *sync.Mutex) string {
	message := QuizReceiveResponse{}
	err := c.ReadJSON(&message)
	if err != nil {
		logger.Error(fmt.Sprintf("socket error start event handling: %s event, %s action", constants.EventStartQuiz, action), zap.Error(err))
		mu.Lock()
		*isConnected = false
		mu.Unlock()
		return constants.UnknownError
	}

	if message.Event == constants.EventStartQuiz {
		return constants.EventStartQuiz
	}

	if message.Event == constants.EventPing {
		return constants.EventPing
	}

	return constants.UnknownError
}

func shareEvenWithUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, event string, sessionId string, invitationCode int, sentToWhom int) {
	payload := map[string]any{"event": event, "response": response}
	data, err := json.Marshal(payload)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error marshal redis payload: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
	}

	if sentToWhom == constants.ToUser || sentToWhom == constants.ToAll {
		// send event to user
		err = qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, sessionId, data).Err()

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

func questionAndScoreHandler(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, isConnected *bool) {
	// get questions/remaining question
	response.Component = constants.Question
	questions, lastQuestionDeliveryTime, err := qc.quizModel.GetSharedQuestions(int(session.InvitationCode.Int32))
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
			sendSingleQuestion(c, qc, &wg, response, session, question, lastQuestionDeliveryTime, chanSkipEvent, chanSkipTimer, len(questions))
		} else { // handle new question
			sendSingleQuestion(c, qc, &wg, response, session, question, sql.NullTime{}, chanSkipEvent, chanSkipTimer, len(questions))
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

func sendSingleQuestion(c *websocket.Conn, qc *quizSocketController, wg *sync.WaitGroup, response *QuizSendResponse, session models.ActiveQuiz, question models.Question, lastQuestionTimeStamp sql.NullTime, chanSkipEvent chan bool, chanSkipTimer chan bool, totalQuestions int) {

	defer wg.Done()

	// start counter if not any question running
	if !lastQuestionTimeStamp.Valid {
		response.Component = constants.Question
		response.Action = constants.ActionCounter
		response.Data = map[string]int{"counter": constants.Counter, "count": constants.Count}
		shareEvenWithUser(c, qc, response, constants.EventStartCount5, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)
		time.Sleep(time.Duration(constants.Counter) * time.Second)

		// update question status to activate
		err := qc.quizModel.UpdateCurrentQuestion(session.ID, question.ID, true)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
			return
		}
	}

	// question sent
	response.Action = constants.ActionSendQuestion
	responseData := map[string]any{
		"id":             question.ID,
		"no":             question.OrderNumber,
		"duration":       question.DurationInSeconds,
		"question_time":  lastQuestionTimeStamp.Time,
		"question":       question.Question,
		"options":        question.Options,
		"totalQuestions": totalQuestions,
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
	err := qc.quizModel.UpdateCurrentQuestion(session.ID, question.ID, false)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
		return
	}

	// score-board rendering
	response.Component = constants.Score
	response.Action = constants.ActionShowScore
	userRankBoard, err := qc.userPlayedQuizModel.GetRank(session.ID, question.ID)

	if err != nil {
		qc.logger.Error("error during get userRankBoard", zap.Error(err))
		return
	}

	response.Data = map[string]any{
		"rankList":       userRankBoard,
		"question":       question.Question,
		"answers":        question.Answers,
		"options":        question.Options,
		"duration":       20,
		"totalQuestions": totalQuestions,
	}

	shareEvenWithUser(c, qc, response, constants.EventShowScore, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)

	wgForSkipTimer := &sync.WaitGroup{}
	wgForSkipTimer.Add(1)

	// skip 20 sec timer
	go handleSkipTimer(wgForSkipTimer, chanSkipTimer)
	wgForSkipTimer.Wait()
}

func terminateQuiz(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz) {

	qc.mu.Lock()
	defer qc.mu.Unlock()

	response.Component = constants.Score
	response.Data = constants.ActionTerminateQuiz
	shareEvenWithUser(c, qc, response, constants.EventTerminateQuiz, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll)

	err := qc.activeQuizModel.Deactivate(session.ID)
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))
		return
	}

	qc.logger.Info("terminateQuiz")
	// here logic of publishing data of user to admin that terminate quiz so no need to listen for joining users
	err = qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, constants.EventTerminateQuiz, constants.EventTerminateQuiz).Err()
	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error while terminationg quiz %s", constants.ActionTerminateQuiz), zap.Error(err))
		return
	}
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
				ok, err := qc.quizModel.IsAllAnswerGathered(session.ID, questionId)
				if err != nil {
					qc.logger.Error("error during listening skip event", zap.Error(err))
					return
				}
				if ok {
					return
				} else { // send warning if all participant not given answer
					response.Data = constants.WarnSkip
					shareEvenWithUser(c, qc, response, constants.EventSkipAsked, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin)
				}
			}
		case user := <-qc.answersSubmittedByUsers:
			response.Data = user
			response.Action = constants.ActionAnserSubmittedByUser
			err := utils.JSONSuccessWs(c, constants.EventAnswerSubmittedByUser, response)
			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error sending event: %s event, %s action, %v user", constants.EventSendQuestion, response.Action, user), zap.Error(err))
				return
			}
		}
	}
}
