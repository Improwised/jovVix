package events

import (
	"fmt"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	evbus "github.com/asaskevich/EventBus"
	"go.uber.org/zap"
)

//go:generate mockery --name IEvents
type IEvents interface {
	SubscribeUserRegistered() error
	Publish(event string, data interface{})
}

type Events struct {
	Bus evbus.Bus `json:"bus"`
}

var logger *zap.Logger = nil

func NewEventBus(lg *zap.Logger) *Events {
	logger = lg
	return &Events{
		Bus: evbus.New(),
	}
}

func userRegistered(data structs.EventUserRegistered) {
	logger.Info(fmt.Sprintf("User Registered email: %s", data.Email))
	// We can send message to SQS, Redis, or any other place
}

func (eb *Events) SubscribeUserRegistered() error {
	return eb.Bus.Subscribe(constants.EventUserRegistered, userRegistered)
}

func (eb *Events) Publish(event string, data interface{}) {
	eb.Bus.Publish(event, data)
}

// SubscribeAll start all subscribers
func (eb *Events) SubscribeAll() error {
	err := eb.Bus.Subscribe(constants.EventUserRegistered, userRegistered)
	if err != nil {
		return err
	}
	return nil
}
