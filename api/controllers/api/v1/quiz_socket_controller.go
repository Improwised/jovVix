package v1

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	quizUtilsHelper "github.com/Improwised/jovvix/api/helpers/utils"
	"github.com/Improwised/jovvix/api/models"
	"github.com/Improwised/jovvix/api/pkg/redis"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/services"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	validator "gopkg.in/go-playground/validator.v9"
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
	Avatar   string
	IsAlive  bool
}

type quizSocketController struct {
	activeQuizModel       *models.ActiveQuizModel
	quizModel             *models.QuizModel
	userPlayedQuizModel   *models.UserPlayedQuizModel
	questionModel         *models.QuestionModel
	userQuizResponseModel *models.UserQuizResponseModel
	presignedURLSvc       *services.PresignURLService
	appConfig             *config.AppConfig
	logger                *zap.Logger
	redis                 *redis.RedisPubSub
}

func InitQuizConfig(db *goqu.Database, appConfig *config.AppConfig, logger *zap.Logger, redis *redis.RedisPubSub) (*quizSocketController, error) {

	activeQuizModel := models.InitActiveQuizModel(db, logger)
	quizModel := models.InitQuizModel(db)
	userPlayedQuizModel := models.InitUserPlayedQuizModel(db)
	questionModel := models.InitQuestionModel(db, logger)
	userQuizResponseModel := models.InitUserQuizResponseModel(db)

	presignedURLSvc, err := services.NewFileUploadServices(&appConfig.AWS)
	if err != nil {
		return nil, err
	}

	return &quizSocketController{
		activeQuizModel:       activeQuizModel,
		quizModel:             quizModel,
		userPlayedQuizModel:   userPlayedQuizModel,
		questionModel:         questionModel,
		userQuizResponseModel: userQuizResponseModel,
		presignedURLSvc:       presignedURLSvc,
		appConfig:             appConfig,
		logger:                logger,
		redis:                 redis,
	}, nil
}

// for user Join
func (qc *quizSocketController) Join(c *websocket.Conn) {
	var JoinMu sync.Mutex

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
	isUserConnected := make(chan bool)

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
				qc.logger.Error("error while reading data from websocket", zap.Error(err))
				updateUserData(qc, userId, session.ID.String(), false)
				isUserConnected <- false
				break
			}
			err = json.Unmarshal([]byte(p), &quizResponse)
			if err != nil {
				qc.logger.Error("error while unmarshaling data from websocket", zap.Error(err))
				updateUserData(qc, userId, session.ID.String(), false)
				break
			}

			if quizResponse.Event == "websocket_close" {
				updateUserData(qc, userId, session.ID.String(), false)
				qc.logger.Info("connection close request is send by the user - " + user.Username)
				break
			}

			if quizResponse.Event == constants.EventPing {

				err := func() error {
					JoinMu.Lock()
					defer JoinMu.Unlock()
					return utils.JSONSuccessWs(c, "pong", "")
				}()

				if err != nil {
					qc.logger.Error("error while sending pong message", zap.Error(err))
				}
			}
		}
	}()

	// is user is a host of current quiz
	if userId == session.AdminID {
		response.Action = constants.ActionCurrentUserIsAdmin
		response.Data = map[string]string{"sessionId": session.ID.String()}

		err := func() error {
			JoinMu.Lock()
			defer JoinMu.Unlock()
			return utils.JSONSuccessWs(c, constants.EventRedirectToAdmin, response)
		}()

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket redirect current user is admin: %s event, %s action, %s code", constants.EventRedirectToAdmin, response.Action, invitationCode), zap.Error(err))
		}
		return
	}

	// when user join at that time publish userName to admin
	publishUserOnJoin(qc, response, user.FirstName, userId, user.ImageKey, session.ID.String())
	response.Action = constants.QuizQuestionStatus
	onConnectHandleUser(c, qc, &response, session, &JoinMu)
	// userPlayedQuizId := quizUtilsHelper.GetString(c.Locals(constants.CurrentUserQuiz))
	handleQuestion(c, qc, session, response, isUserConnected, &JoinMu)
}

func publishUserOnJoin(qc *quizSocketController, quizResponse QuizSendResponse, userName string, userId string, avatar string, sessionId string) {
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
		newUser := UserInfo{UserId: userId, UserName: userName, Avatar: avatar, IsAlive: true}
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
		newUser := UserInfo{UserId: userId, UserName: userName, Avatar: avatar, IsAlive: true}
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

	err = qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserJoin, sessionId), publishData).Err()

	if err != nil {
		qc.logger.Error(fmt.Sprintf("socket error publishing event: %s event, %s action", constants.EventUserJoined, response.Action), zap.Error(err))
	}
}

func handleQuestion(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, response QuizSendResponse, isUserConnected chan bool, joinMu *sync.Mutex) {
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
	for {
		select {
		case isConnected := <-isUserConnected:
			if !isConnected {
				return
			}
		case msg := <-ch:
			message := map[string]any{}
			err := json.Unmarshal([]byte(msg.Payload), &message)

			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventJoinQuiz, response.Action), zap.Error(err))
			}

			event := quizUtilsHelper.GetString(message["event"])

			err = func() error {
				joinMu.Lock()
				defer joinMu.Unlock()
				return utils.JSONSuccessWs(c, event, message["response"])
			}()

			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", event, response.Action), zap.Error(err))
			}

			if message["event"] == constants.EventTerminateQuiz {
				return
			}
		}
	}
}

func onConnectHandleUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, joinMu *sync.Mutex) {
	if session.CurrentQuestion.Valid {

		totalQuestion, err := qc.questionModel.GetTotalQuestionCount(session.ID.String())
		if err != nil {
			qc.logger.Error(constants.ErrInGettingTotalQuestionCount, zap.Error(err))
			return
		}

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
			"id":             currentQuestion.ID,
			"no":             currentQuestion.OrderNumber,
			"duration":       duration,
			"question_time":  session.QuestionDeliveryTime.Time,
			"question":       currentQuestion.Question,
			"options":        currentQuestion.Options,
			"totalQuestions": totalQuestion,
			"question_media": currentQuestion.QuestionMedia,
			"options_media":  currentQuestion.OptionsMedia,
			"resource":       currentQuestion.Resource.String,
		}
		response.Data = responseData
		response.Component = constants.Question

		err = func() error {
			joinMu.Lock()
			defer joinMu.Unlock()
			return utils.JSONSuccessWs(c, constants.EventSendQuestion, response)
		}()

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error send current question on connect: %s event, %s action", constants.EventSendQuestion, response.Action), zap.Error(err))
		}
	} else {
		response.Data = constants.QuizStartsSoon

		err := func() error {
			joinMu.Lock()
			defer joinMu.Unlock()
			return utils.JSONSuccessWs(c, constants.EventJoinQuiz, response)
		}()

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
		if err := qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserDisconnect, sessionId), publishData).Err(); err != nil {
			qc.logger.Error("error publishing data in updateUserData", zap.Error(err))
		}
	}
}

// filter user data and publish only alive user to the channel
func filterPublishUsers(qc *quizSocketController, usersData []UserInfo, functionName string) (publishData []byte) {

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
	var arrangeMu sync.Mutex

	isConnected := true
	adminDisconnected := make(chan bool, 1)

	response := QuizSendResponse{
		Component: constants.Waiting,
		Action:    constants.ActionAuthentication,
		Data:      "",
	}

	sessionId := quizUtilsHelper.GetString(c.Locals(constants.SessionIDParam))

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))

	if !ok {
		qc.logger.Error("socket user-model type conversion")

		err := func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONFailWs(c, constants.EventSessionValidation, constants.UnknownError)
		}()

		if err != nil {
			qc.logger.Error("socket user-model type conversion")
		}
		return
	}

	// activate session
	session, err := ActivateAndGetSession(c, qc.activeQuizModel, qc.logger, sessionId, user.ID, &arrangeMu)

	if err != nil {
		qc.logger.Error("get active session", zap.Error(err))

		err := func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONFailWs(c, constants.EventSessionValidation, constants.UnknownError)
		}()

		if err != nil {
			qc.logger.Error("get active session", zap.Error(err))
		}
		return
	}

	defer func() {
		isConnected = false
		adminDisconnected <- true
		time.Sleep(1 * time.Second)
		c.Close()
		qc.logger.Info("connection closed by admin")
	}()

	// handle code sharing with admin
	handleCodeGeneration(c, qc, session, &isConnected, &response, adminDisconnected, &arrangeMu)

	// if connection lost during waiting of start event
	if !(isConnected) {
		response.Component = constants.Loading
		response.Data = constants.AdminDisconnected
		shareEvenWithUser(c, qc, &response, constants.AdminDisconnected, sessionId, int(session.InvitationCode.Int32), constants.ToUser, &arrangeMu)

		qc.logger.Error("admin disconnected")
		return
	}

	// handle when user join during running quiz
	go handleRunningQuizUserJoin(c, qc, adminDisconnected, session.ID.String(), &arrangeMu)

	// question and score handler
	questionAndScoreHandler(c, qc, &response, session, &isConnected, &arrangeMu)
}

// handleRunningQuizUserJoin listens for users joining a running quiz and sends the updated count of users to the client.
// It also listens for admin disconnection and gracefully terminates if the admin is disconnected.
func handleRunningQuizUserJoin(c *websocket.Conn, qc *quizSocketController, adminDisconnected chan bool, sessionId string, arrangeMu *sync.Mutex) {
	response := QuizSendResponse{}
	response.Action = constants.JoinUserOnRunningQuiz
	response.Component = constants.Running

	pubsub := qc.redis.PubSubModel.Client.Subscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserJoin, sessionId))
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserJoin, sessionId))
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()

	for {
		select {
		case isDisconnected := <-adminDisconnected:
			if isDisconnected {
				return
			}
		case msg := <-ch:
			response.Data = msg.Payload

			// get total user in active quiz
			totalUserJoin, err := qc.userPlayedQuizModel.GetCountOfTotalJoinUsers(sessionId)
			if err != nil {
				qc.logger.Error(constants.ErrGetTotalJoinUser, zap.Error(err))
				continue
			}
			response.Data = totalUserJoin

			err = func() error {
				arrangeMu.Lock()
				defer arrangeMu.Unlock()
				return utils.JSONSuccessWs(c, constants.JoinUserOnRunningQuiz, response)
			}()

			if err != nil {
				qc.logger.Error("error while sending user data ", zap.Error(err))
			}
		}
	}
}

func ActivateAndGetSession(c *websocket.Conn, activeQuizModel *models.ActiveQuizModel, logger *zap.Logger, sessionId string, userId string, arrangeMu *sync.Mutex) (models.ActiveQuiz, error) {

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

			err = func() error {
				arrangeMu.Lock()
				defer arrangeMu.Unlock()
				return utils.JSONFailWs(c, constants.EventAuthorization, response)
			}()

			if err != nil {
				logger.Error(fmt.Sprintf("socket error authentication host: %s event, %s action", constants.EventAuthorization, response.Action), zap.Error(err))
			}
			return session, err
		} else if err.Error() == constants.ErrSessionWasCompleted {
			response.Action = constants.ActionSessionActivation
			response.Data = constants.ErrSessionWasCompleted

			err = func() error {
				arrangeMu.Lock()
				defer arrangeMu.Unlock()
				return utils.JSONFailWs(c, constants.EventAuthorization, response)
			}()

			if err != nil {
				logger.Error(fmt.Sprintf("socket error authentication host: %s event, %s action", constants.EventAuthorization, response.Action), zap.Error(err))
			}
			return session, err
		}

		response.Action = constants.ActionSessionActivation
		response.Data = constants.UnknownError
		logger.Debug("unknown error was triggered from ActivateAndGetSession")

		err = func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONErrorWs(c, constants.EventActivateSession, response)
		}()

		if err != nil {
			logger.Error(fmt.Sprintf("socket error get or activate session: %s event, %s action", constants.EventActivateSession, response.Action), zap.Error(err))
		}
		return session, err
	}

	c.Locals(constants.ActiveQuizObj, session)

	return session, nil
}

func handleCodeGeneration(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, isConnected *bool, response *QuizSendResponse, adminDisconnected chan bool, arrangeMu *sync.Mutex) {
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
				handleInvitationCodeSend(c, response, qc.logger, session.InvitationCode.Int32, arrangeMu)
				isInvitationCodeSent = true
				go handleConnectedUser(c, qc, session.ID.String(), adminDisconnected, arrangeMu)

			}

			// once code sent receive start signal
			if isInvitationCodeSent {
				isBreak := handleStartQuiz(c, qc.logger, isConnected, response.Action)

				if isBreak == constants.EventPing {

					err := func() error {
						arrangeMu.Lock()
						defer arrangeMu.Unlock()
						return utils.JSONSuccessWs(c, constants.EventPong, "")
					}()

					if err != nil {
						qc.logger.Error("error while sending pong message", zap.Error(err))
					}
				} else {
					users, err := qc.redis.PubSubModel.Client.Get(qc.redis.PubSubModel.Ctx, session.ID.String()).Result()
					if err != nil {
						qc.logger.Error("error while fetching data from redis inside updateUserData", zap.Error(err))
					}

					var usersData []UserInfo
					err = json.Unmarshal([]byte(users), &usersData)
					if err != nil {
						qc.logger.Error("error while unmarshaling redis inside updateUserData", zap.Error(err))
					}
					if len(usersData) != 0 && isBreak == constants.EventStartQuiz {

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

						err = func() error {
							arrangeMu.Lock()
							defer arrangeMu.Unlock()
							return utils.JSONFailWs(c, constants.EventStartQuiz, response)
						}()

						if err != nil {
							qc.logger.Error(fmt.Sprintf("socket error middleware: %s event, %s action", constants.EventAuthentication, response.Action), zap.Error(err))
						}
					}
				}

			}
		}
	}
}

// handle waiting page
func handleInvitationCodeSend(c *websocket.Conn, response *QuizSendResponse, logger *zap.Logger, invitationCode int32, arrangeMu *sync.Mutex) bool {

	// send code to client
	response.Action = constants.ActionSessionActivation
	response.Data = map[string]int{"code": int(invitationCode)}

	err := func() error {
		arrangeMu.Lock()
		defer arrangeMu.Unlock()
		return utils.JSONSuccessWs(c, constants.EventSendInvitationCode, response)
	}()

	if err != nil {
		logger.Error(fmt.Sprintf("socket error sent code: %s event, %s action", constants.EventSendInvitationCode, response.Action), zap.Error(err))
	}

	return true
}

// when user connect at that time send data to admin
func handleConnectedUser(c *websocket.Conn, qc *quizSocketController, sessionId string, adminDisconnected chan bool, arrangeMu *sync.Mutex) {
	response := QuizSendResponse{}
	response.Action = constants.ActionSendUserData
	response.Component = constants.Waiting

	pubsub := qc.redis.PubSubModel.Client.Subscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserJoin, sessionId), fmt.Sprintf("%s-%s", constants.ChannelUserDisconnect, sessionId), constants.EventTerminateQuiz, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelUserJoin, sessionId), fmt.Sprintf("%s-%s", constants.ChannelUserDisconnect, sessionId), constants.EventTerminateQuiz, constants.EventStartQuizByAdmin, constants.StartQuizByAdminNoPlayerFound)
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()
	usersData := []UserInfo{}

	for {
		select {
		case isDisconnected := <-adminDisconnected:
			if isDisconnected {
				return
			}
		case msg := <-ch:
			response.Data = msg.Payload

			if response.Data == constants.StartQuizByAdminNoPlayerFound {
				continue
			}

			if response.Data == constants.EventStartQuizByAdmin || response.Data == constants.EventTerminateQuiz {
				return
			}

			err := json.Unmarshal([]byte(msg.Payload), &usersData)
			if err != nil {
				qc.logger.Error("error while unmarshaling data inside handleconnectedUser ", zap.Error(err))

				break
			}

			response.Data = usersData

			err = func() error {
				arrangeMu.Lock()
				defer arrangeMu.Unlock()
				return utils.JSONSuccessWs(c, constants.EventSendInvitationCode, response)
			}()
			// sending the user data to the admin
			if err != nil {
				qc.logger.Error("error while sending user data ", zap.Error(err))
			}
		}
	}
}

// start quiz by message event from admin
func handleStartQuiz(c *websocket.Conn, logger *zap.Logger, isConnected *bool, action string) string {
	message := QuizReceiveResponse{}
	err := c.ReadJSON(&message)
	if err != nil {
		logger.Error(fmt.Sprintf("socket error start event handling: %s event, %s action", constants.EventStartQuiz, action), zap.Error(err))
		*isConnected = false
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

func shareEvenWithUser(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, event string, sessionId string, invitationCode int, sentToWhom int, arrangeMu *sync.Mutex) {
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

		err := func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONSuccessWs(c, event, response)
		}()

		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error sending event: %s event, %s action %v code", constants.EventSendQuestion, response.Action, invitationCode), zap.Error(err))
		}
	}
}

func questionAndScoreHandler(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, isConnected *bool, arrangeMu *sync.Mutex) {
	// get questions/remaining question
	response.Component = constants.Question
	questions, lastQuestionDeliveryTime, err := qc.quizModel.GetSharedQuestions(int(session.InvitationCode.Int32))
	if err != nil {
		response.Action = constants.ErrInGettingQuestion
		qc.logger.Error(fmt.Sprintf("socket error get remaining questions: %s event, %s action %v code", constants.EventStartQuiz, response.Action, session.InvitationCode), zap.Error(err))

		err := func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONFailWs(c, constants.EventSendQuestion, response)
		}()

		if err != nil {
			qc.logger.Error("error during get remaining question", zap.Error(err))
		}
		return
	}

	totalQuestion, err := qc.questionModel.GetTotalQuestionCount(session.ID.String())
	if err != nil {
		qc.logger.Error(constants.ErrInGettingTotalQuestionCount, zap.Error(err))
		return
	}

	var wg sync.WaitGroup
	chanNextEvent := make(chan bool)
	chanSkipEvent := make(chan bool)
	chanSkipTimer := make(chan bool)
	chanPauseQuiz := make(chan bool)
	var isQuizEnd bool = false

	go listenAllEvents(c, qc, response, session, chanNextEvent, chanSkipEvent, chanSkipTimer, chanPauseQuiz, isQuizEnd, arrangeMu)

	// handle question
	var isFirst bool = lastQuestionDeliveryTime.Valid
	response.Component = constants.Question
	for _, question := range questions {
		wg.Add(1)
		if isFirst { // handle running question
			isFirst = false
			sendSingleQuestion(c, qc, &wg, response, session, question, lastQuestionDeliveryTime, chanSkipEvent, chanSkipTimer, chanPauseQuiz, totalQuestion, arrangeMu)
		} else { // handle new question
			sendSingleQuestion(c, qc, &wg, response, session, question, sql.NullTime{}, chanSkipEvent, chanSkipTimer, chanPauseQuiz, totalQuestion, arrangeMu)
		}

		err := func() error {
			arrangeMu.Lock()
			defer arrangeMu.Unlock()
			return utils.JSONSuccessWs(c, constants.EventNextQuestionAsked, response)
		}()

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
		terminateQuiz(c, qc, response, session, arrangeMu)
		// isQuizEnd = false
	}
}

func listenAllEvents(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, chanNextEvent chan bool, chanSkipEvent chan bool, chanSkipTimer chan bool, chanPauseQuiz chan bool, isQuizEnd bool, arrangeMu *sync.Mutex) {
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
		case constants.EventPauseQuiz:
			isPauseQuiz := message.Data.(bool)
			chanPauseQuiz <- isPauseQuiz
		}
	}

	// handle connection lost during quiz
	if !isQuizEnd {
		response.Component = constants.Loading
		response.Data = constants.AdminDisconnected
		shareEvenWithUser(c, qc, response, constants.AdminDisconnected, session.ID.String(), int(session.InvitationCode.Int32), constants.ToUser, arrangeMu)
	}
}

func sendSingleQuestion(c *websocket.Conn, qc *quizSocketController, wg *sync.WaitGroup, response *QuizSendResponse, session models.ActiveQuiz, question models.Question, lastQuestionTimeStamp sql.NullTime, chanSkipEvent chan bool, chanSkipTimer chan bool, chanPauseQuiz chan bool, totalQuestions int64, arrangeMu *sync.Mutex) {

	defer wg.Done()

	if question.QuestionMedia == "image" {
		presignedURL, err := qc.presignedURLSvc.GetPresignedURL(question.Resource.String, 5*time.Minute)
		if err != nil {
			qc.logger.Error("error while generating presign url")
		}
		question.Resource = sql.NullString{String: presignedURL, Valid: true}
	}

	if question.OptionsMedia == "image" {
		for i, v := range question.Options {
			presignedURL, err := qc.presignedURLSvc.GetPresignedURL(v, 1*time.Minute)
			if err != nil {
				qc.logger.Error("error while generating presign url")
			}
			question.Options[i] = presignedURL
		}
	}

	totalUserJoin, err := qc.userPlayedQuizModel.GetCountOfTotalJoinUsers(session.ID.String())
	if err != nil {
		qc.logger.Error(constants.ErrGetTotalJoinUser, zap.Error(err))
		return
	}

	var questionStartTime time.Time
	
	// start counter if not any question running
	if !lastQuestionTimeStamp.Valid {
		response.Component = constants.Question
		response.Action = constants.ActionCounter
		response.Data = map[string]int{"counter": constants.Counter, "count": constants.Count}
		shareEvenWithUser(c, qc, response, constants.EventStartCount5, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll, arrangeMu)
		time.Sleep(time.Duration(constants.Counter) * time.Second)

		// Set the question start time to NOW (after counter finishes)
		questionStartTime = time.Now()

		// update question status to activate
		err := qc.quizModel.UpdateCurrentQuestion(session.ID, question.ID, true)
		if err != nil {
			qc.logger.Error(fmt.Sprintf("socket error update current question: %s event, %s action %v code", constants.EventSendQuestion, response.Action, session.InvitationCode), zap.Error(err))
			return
		}
	} else {
		// For running questions, use the existing timestamp
		questionStartTime = lastQuestionTimeStamp.Time
	}

	// question sent
	response.Action = constants.ActionSendQuestion
	responseData := map[string]any{
		"id":             question.ID,
		"quiz_id":        question.QuizId,
		"no":             question.OrderNumber,
		"duration":       question.DurationInSeconds,
		"start_time":     questionStartTime.Format(time.RFC3339), 
		"question":       question.Question,
		"options":        question.Options,
		"question_media": question.QuestionMedia,
		"options_media":  question.OptionsMedia,
		"resource":       question.Resource.String,
		"totalQuestions": totalQuestions,
		"totalJoinUser":  totalUserJoin,
	}
	response.Data = responseData
	if !lastQuestionTimeStamp.Valid { // handling new question
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll, arrangeMu)
	} else { // handling running question
		shareEvenWithUser(c, qc, response, constants.EventSendQuestion, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin, arrangeMu)
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
	go handleAnswerSubmission(c, qc, session, question.ID, duration, wgForQuestion, chanSkipEvent, response, arrangeMu)
	wgForQuestion.Wait()

	// update current status to deactivate
	err = qc.quizModel.UpdateCurrentQuestion(session.ID, question.ID, false)
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
	userResponses, err := qc.userQuizResponseModel.GetUsersResponses(session.ID, question.ID)
	if err != nil {
		qc.logger.Error("error during get userResponses", zap.Error(err))
	}

	scoreboardMaxDurationEnv := qc.appConfig.Quiz.ScoreboardMaxDuration
	scoreboardMaxDuration := 20

	if scoreboardMaxDurationEnv != "" {
		if parsedDuration, err := strconv.Atoi(scoreboardMaxDurationEnv); err == nil {
			scoreboardMaxDuration = parsedDuration
		}
	}

	response.Data = map[string]any{
		"question_no":    question.OrderNumber,
		"quiz_id":        question.QuizId,
		"rankList":       userRankBoard,
		"question":       question.Question,
		"answers":        question.Answers,
		"options":        question.Options,
		"question_media": question.QuestionMedia,
		"options_media":  question.OptionsMedia,
		"resource":       question.Resource.String,
		"duration":       scoreboardMaxDuration,
		"totalQuestions": totalQuestions,
		"userResponses":  userResponses,
	}
	shareEvenWithUser(c, qc, response, constants.EventShowScore, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin, arrangeMu)

	response.Data = map[string]any{
		"question_no":    question.OrderNumber,
		"quiz_id":        question.QuizId,
		"rankList":       userRankBoard,
		"question":       question.Question,
		"answers":        question.Answers,
		"options":        question.Options,
		"question_media": question.QuestionMedia,
		"options_media":  question.OptionsMedia,
		"resource":       question.Resource.String,
		"duration":       scoreboardMaxDuration,
		"totalQuestions": totalQuestions,
	}
	shareEvenWithUser(c, qc, response, constants.EventShowScore, session.ID.String(), int(session.InvitationCode.Int32), constants.ToUser, arrangeMu)

	wgForSkipTimer := &sync.WaitGroup{}
	wgForSkipTimer.Add(1)

	// skip timer
	go handleSkipTimer(c, qc, wgForSkipTimer, response, session, chanSkipTimer, chanPauseQuiz, scoreboardMaxDuration, arrangeMu)
	wgForSkipTimer.Wait()
}

func terminateQuiz(c *websocket.Conn, qc *quizSocketController, response *QuizSendResponse, session models.ActiveQuiz, arrangeMu *sync.Mutex) {

	response.Component = constants.Score
	response.Data = constants.ActionTerminateQuiz
	shareEvenWithUser(c, qc, response, constants.EventTerminateQuiz, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAll, arrangeMu)

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

func handleSkipTimer(c *websocket.Conn, qc *quizSocketController, wg *sync.WaitGroup, response *QuizSendResponse, session models.ActiveQuiz, chanSkipTimer chan bool, chanPauseQuiz chan bool, scoreboardMaxDuration int, arrangeMu *sync.Mutex) {
	defer wg.Done()

	remainingTime := time.Duration(scoreboardMaxDuration) * time.Second
	startTime := time.Now()
	isTimeout := time.NewTimer(remainingTime)
	timerPaused := false

	for {
		select {
		case <-isTimeout.C:
			if !timerPaused {
				return
			}
		case isSkip := <-chanSkipTimer:
			if isSkip {
				return
			}
		case isPause := <-chanPauseQuiz:
			if isPause {
				// Stop the timer and calculate the remaining time
				if !isTimeout.Stop() {
					// drain the channel if the timer has expired
					<-isTimeout.C
				}
				remainingTime -= time.Since(startTime)
				timerPaused = true

				// send event to the user
				response.Component = constants.Score
				response.Data = constants.EventPauseQuiz
				shareEvenWithUser(c, qc, response, constants.EventPauseQuiz, session.ID.String(), int(session.InvitationCode.Int32), constants.ToUser, arrangeMu)
			} else {
				// Resume with the remaining time
				if timerPaused {
					startTime = time.Now()
					isTimeout.Reset(remainingTime)
					timerPaused = false

					// send event to the user
					response.Component = constants.Score
					response.Data = constants.EventResumeQuiz
					shareEvenWithUser(c, qc, response, constants.EventResumeQuiz, session.ID.String(), int(session.InvitationCode.Int32), constants.ToUser, arrangeMu)
				}
			}
		}
	}
}

func handleAnswerSubmission(c *websocket.Conn, qc *quizSocketController, session models.ActiveQuiz, questionId uuid.UUID, duration int, wg *sync.WaitGroup, chanSkipEvent chan bool, response *QuizSendResponse, arrangeMu *sync.Mutex) {
	defer wg.Done()

	isTimeout := time.NewTicker(time.Duration(duration) * time.Second)

	pubsub := qc.redis.PubSubModel.Client.Subscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelSetAnswer, session.ID.String()))
	defer func() {
		if pubsub != nil {
			err := pubsub.Unsubscribe(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelSetAnswer, session.ID.String()))
			if err != nil {
				qc.logger.Error("unsubscribe failed", zap.Error(err))
			}
			pubsub.Close()
		}
	}()

	ch := pubsub.Channel()

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
					shareEvenWithUser(c, qc, response, constants.EventSkipAsked, session.ID.String(), int(session.InvitationCode.Int32), constants.ToAdmin, arrangeMu)
				}
			}
		case msg := <-ch:
			user := models.User{}

			err := json.Unmarshal([]byte(msg.Payload), &user)
			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error send waiting message: %s event, %s action", constants.EventSendQuestion, constants.ActionAnserSubmittedByUser), zap.Error(err))
			}

			response.Data = user
			response.Action = constants.ActionAnserSubmittedByUser

			err = func() error {
				arrangeMu.Lock()
				defer arrangeMu.Unlock()
				return utils.JSONSuccessWs(c, constants.EventAnswerSubmittedByUser, response)
			}()

			if err != nil {
				qc.logger.Error(fmt.Sprintf("socket error sending event: %s event, %s action, %v user", constants.EventSendQuestion, response.Action, user), zap.Error(err))
			}
		}
	}
}

func (qc *quizSocketController) SetAnswer(c *fiber.Ctx) error {
	currentQuiz := c.Query(constants.CurrentUserQuiz)
	sessionId := c.Query(constants.SessionIDParam)

	// validations
	if currentQuiz == "" {
		qc.logger.Error(constants.ErrQuizNotFound)
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuizNotFound)
	}

	currentQuizId, err := uuid.Parse(currentQuiz)
	if err != nil {
		qc.logger.Error("invalid UUID")
		return utils.JSONFail(c, http.StatusBadRequest, "invalid UUID")
	}

	user, ok := quizUtilsHelper.ConvertType[models.User](c.Locals(constants.ContextUser))
	if !ok {
		qc.logger.Error("Unable to convert to user-model type from locals")
		return utils.JSONFail(c, http.StatusInternalServerError, "Unable to convert to user-model type from locals")
	}

	var answer structs.ReqAnswerSubmit

	err = json.Unmarshal(c.Body(), &answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, err.Error())
	}

	validate := validator.New()
	err = validate.Struct(answer)
	if err != nil {
		return utils.JSONFail(c, http.StatusBadRequest, utils.ValidatorErrorString(err))
	}

	// check for question is active or not to receive answers
	currentQuestion, err := qc.userPlayedQuizModel.GetCurrentActiveQuestion(sessionId)
	if err != nil {
		if err == sql.ErrNoRows {
			qc.logger.Error("error during answer submit get current active question", zap.Any("answers", answer), zap.Any("current_quiz_id", currentQuizId))
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerSubmit)
		}
		qc.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusBadRequest, constants.UnknownError)
	}

	if currentQuestion != answer.QuestionId {
		qc.logger.Error(constants.ErrQuestionNotActive)
		return utils.JSONFail(c, http.StatusBadRequest, constants.ErrQuestionNotActive)
	}

	answers, answerPoints, answerDurationInSeconds, questionType, err := qc.questionModel.GetAnswersPointsDurationType(answer.QuestionId.String())
	if err != nil {
		qc.logger.Error("error while get answer, points, duration and type")
		return utils.JSONFail(c, http.StatusBadRequest, "error while get answer, points, duration and type")
	}

	// calculate points
	points, score := utils.CalculatePointsAndScore(answer, answers, answerPoints, answerDurationInSeconds, questionType)

	streakCount, err := qc.userPlayedQuizModel.GetStreakCount(currentQuizId, answer.QuestionId)
	if err != nil {
		if err == sql.ErrNoRows {
			qc.logger.Error(constants.ErrGetStreakCount, zap.Error(err))
			return utils.JSONError(c, http.StatusBadRequest, constants.ErrQuizNotFound)
		}
		qc.logger.Error(constants.ErrGetStreakCount, zap.Error(err))
		return utils.JSONError(c, http.StatusInternalServerError, constants.ErrGetStreakCount)
	}

	// add streak score and update streak also
	finalScore, newSreakCount := utils.CalculateStreakScore(streakCount, score)

	// Submit answer
	if err := qc.userQuizResponseModel.SubmitAnswer(currentQuizId, answer, points, finalScore, newSreakCount); err != nil {
		if err == sql.ErrNoRows {
			return utils.JSONFail(c, http.StatusBadRequest, constants.ErrAnswerAlreadySubmitted)
		}
		qc.logger.Error("error during answer submit", zap.Error(err))
		return utils.JSONFail(c, http.StatusInternalServerError, constants.UnknownError)
	}

	// Publish to Redis in a goroutine
	go func() {
		data, err := json.Marshal(user)
		if err != nil {
			qc.logger.Error("Error marshaling user data", zap.Error(err))
			return
		}

		if err := qc.redis.PubSubModel.Client.Publish(qc.redis.PubSubModel.Ctx, fmt.Sprintf("%s-%s", constants.ChannelSetAnswer, sessionId), data).Err(); err != nil {
			qc.logger.Error("Error publishing answer to Redis", zap.Error(err))
		}
	}()

	return utils.JSONSuccess(c, http.StatusAccepted, nil)
}

func (ctrl *quizSocketController) Terminate(c *fiber.Ctx) error {
	return utils.JSONSuccess(c, http.StatusOK, nil)
}
