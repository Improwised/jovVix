package models

import (
	"context"

	redis "github.com/redis/go-redis/v9"
)

type PubSubModel struct {
	Ctx    context.Context
	Client redis.Client
}

func InitPubSubModel(addr, password string, db int) (*PubSubModel, error) {
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     addr,
	// 	Password: password,
	// 	DB:       db,
	// })
	// ctx := context.Background()

	// _, err := client.Ping(ctx).Result()

	// if err != nil {
	// 	return nil, err
	// }

	return &PubSubModel{
		// Ctx:    ctx,
		// Client: *client,
	}, nil
}
