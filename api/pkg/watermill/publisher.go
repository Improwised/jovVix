package watermill

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/Improwised/jovvix/api/cli/workers"
	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/database"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill-googlecloud/pkg/googlecloud"
	"github.com/ThreeDotsLabs/watermill-sql/v2/pkg/sql"

	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	redis "github.com/redis/go-redis/v9"
)

type WatermillPublisher struct {
	publisher message.Publisher
}

func InitPublisher(cfg config.AppConfig, isDeadLetterQ bool) (*WatermillPublisher, error) {
	logger = watermill.NewStdLogger(cfg.MQ.Debug, cfg.MQ.Track)
	if isDeadLetterQ {
		return initSqlPub(cfg)
	}
	switch cfg.MQ.Dialect {
	case "amqp":
		return initAmqpPub(cfg)

	case "redis":
		return initRedisPub(cfg)

	case "kafka":
		return initKafkaPub(cfg)

	case "googlecloud":
		return initGoogleCloudPub(cfg)

	case "sql":
		return initSqlPub(cfg)
	default:
		return &WatermillPublisher{}, nil
	}

}

// send message into queue using topic name
//
// struct must from worker package(/cli/workers)
func (wp *WatermillPublisher) Publish(topic string, handle workers.Handler) error {
	// if broker is not set then return nil
	if wp.publisher == nil {
		return nil
	}
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)

	err := enc.Encode(&handle)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), network.Bytes())
	err = wp.publisher.Publish(topic, msg)
	return err
}

func initAmqpPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	amqpConfig := amqp.NewDurableQueueConfig(cfg.MQ.Amqp.AmqbUrl)
	publisher, err := amqp.NewPublisher(amqpConfig, logger)
	return &WatermillPublisher{publisher: publisher}, err
}

func initRedisPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	pubClient := redis.NewClient(&redis.Options{
		Addr:     cfg.MQ.Redis.RedisUrl,
		Username: cfg.MQ.Redis.UserName,
		Password: cfg.MQ.Redis.Password,
	})
	publisher, err := redisstream.NewPublisher(
		redisstream.PublisherConfig{
			Client:     pubClient,
			Marshaller: redisstream.DefaultMarshallerUnmarshaller{},
		},
		logger,
	)
	return &WatermillPublisher{publisher: publisher}, err
}

func initKafkaPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:               cfg.MQ.Kafka.KafkaBroker,
			Marshaler:             kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: kafka.DefaultSaramaSyncPublisherConfig(),
		},
		logger,
	)
	return &WatermillPublisher{publisher: publisher}, err
}

func initGoogleCloudPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	publisher, err := googlecloud.NewPublisher(googlecloud.PublisherConfig{
		ProjectID:      cfg.MQ.GoogleCloud.ProjectID,
		ConnectTimeout: 10 * time.Second,
		PublishTimeout: 10 * time.Second,
		Marshaler:      googlecloud.DefaultMarshalerUnmarshaler{},
	}, logger)

	return &WatermillPublisher{publisher: publisher}, err
}

func initSqlPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	switch cfg.MQ.Sql.Dialect {
	case "postgres":
		return initPostgresPub(cfg)
	case "mysql":
		return initMysqlPub(cfg)
	default:
		return nil, nil
	}

}

func initPostgresPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	db, err := database.PostgresDBConnection(cfg.MQ.Sql)
	if err != nil {
		return nil, err
	}
	publisher, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			// we are customizing schema adapter because default schema adapter has some issue
			SchemaAdapter:        database.PostgreSQLSchema{},
			AutoInitializeSchema: true,
		},

		logger,
	)
	if err != nil {
		return nil, err
	}
	return &WatermillPublisher{publisher: publisher}, nil
}

func initMysqlPub(cfg config.AppConfig) (*WatermillPublisher, error) {
	db, err := database.MysqlDBConnection(cfg.MQ.Sql)
	if err != nil {
		return nil, err
	}
	publisher, err := sql.NewPublisher(
		db,
		sql.PublisherConfig{
			SchemaAdapter:        database.MySQLSchema{},
			AutoInitializeSchema: true,
		},

		logger,
	)
	if err != nil {
		return nil, err
	}

	return &WatermillPublisher{publisher: publisher}, nil
}
