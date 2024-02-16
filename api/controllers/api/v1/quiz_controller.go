package v1

import (
	"fmt"
	"log"
	"time"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/services"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/gofiber/contrib/websocket"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

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

			copyUserObj.Username = utils.GenerateNewStringHavingSuffixName(userObj.Username, 5, 12)

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

type quizConfigs struct {
	db        *models.QuizModel
	userCtrl  *UserController
	appConfig *config.AppConfig
}

func InitQuizController(db *goqu.Database, userCtrl *UserController, appConfig *config.AppConfig) (*quizConfigs, error) {
	return &quizConfigs{models.InitQuizModel(db), userCtrl, appConfig}, nil
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

func (*quizConfigs) Ping(c *websocket.Conn) {

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

func (qc *quizConfigs) Join(c *websocket.Conn) {

	defer c.Close()

	// check for middleware error
	if !c.Locals(constants.MiddlewarePass).(bool) {
		fmt.Print(utils.WsJSONError(c, "authentication failed", c.Locals(constants.MiddlewareError).(string)))
		time.Sleep(1 * time.Second)
		return
	}

	time.Sleep(1 * time.Second)

	err := utils.WsJSONSuccess(c, "get current code", c.Query("code"))

	if err != nil {
		fmt.Println("Error: ", err)
	}

	time.Sleep(1 * time.Second)

	err = utils.WsJSONSuccess(c, "get user id", map[string]string{constants.ContextUid: c.Locals(constants.ContextUid).(string)})
	if err != nil {
		fmt.Println("Error: ", err)
	}

	time.Sleep(1 * time.Second)

	err = utils.WsJSONSuccess(c, "get user context", map[string]any{constants.ContextUser: c.Locals(constants.ContextUser)})

	if err != nil {
		fmt.Println("Error: ", err)
	}

	c.Locals(constants.ContextUid)

}
