package config

type MQConfig struct {
	Dialect     string `envconfig:"MQ_DIALECT"`
	Debug       bool   `envconfig:"MQ_DEBUG"`
	Track       bool   `envconfig:"MQ_TRACK"`
	DeadLetterQ string `envconfig:"DEAD_LETTER_QUEUE"`
	HandlerName string `envconfig:"HANDLER_NAME"`
	Redis       RedisConfig
	Amqp        AmqpConfig
	Kafka       KafkaConfig
	GoogleCloud GoogleCloud
	Sql         Sql
}
type RedisConfig struct {
	RedisUrl      string `envconfig:"REDIS_URI"`
	ConsumerGroup string `envconfig:"CONSUMER_GROUP"`
	UserName      string `envconfig:"REDIS_USERNAME"`
	Password      string `envconfig:"REDIS_PASSWORD"`
}

type AmqpConfig struct {
	AmqbUrl string `envconfig:"AMQB_URI"`
}
type KafkaConfig struct {
	KafkaBroker   []string `envconfig:"KAFKA_BROKER"`
	ConsumerGroup string   `envconfig:"CONSUMER_GROUP"`
}

type GoogleCloud struct {
	ProjectID      string `envconfig:"PROJECT_ID"`
	SubscriptionId string `envconfig:"SUBSCRIPTION_ID"`
}

type Sql struct {
	Dialect     string `envconfig:"MQ_DB_DIALECT"`
	Host        string `envconfig:"MQ_DB_HOST"`
	Port        int    `envconfig:"MQ_DB_PORT"`
	Username    string `envconfig:"MQ_DB_USERNAME"`
	Password    string `envconfig:"MQ_DB_PASSWORD"`
	Db          string `envconfig:"MQ_DB_NAME"`
	QueryString string `envconfig:"DB_QUERYSTRING"`
}
