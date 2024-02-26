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
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	ctx := context.Background()

	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return &PubSubModel{
		Ctx:    ctx,
		Client: *client,
	}, nil
}

// err = client.Set(ctx, "name", "Elliot", 0).Err()
// // if there has been an error setting the value
// // handle the error
// if err != nil {
// 	fmt.Println(err)
// }

// go func() {
// 	time.Sleep(5 * time.Second)
// 	err = client.Publish(ctx, "mychannel1", "payload").Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	time.Sleep(2*time.Second)
// 	err = client.Publish(ctx, "mychannel1", "close").Err()
// 	if err != nil {
// 		panic(err)
// 	}
// }()

// time.Sleep(3 * time.Second)
// val, err := client.Get(ctx, "name").Result()
// if err != nil {
// 	fmt.Println(err)
// }

// fmt.Println(val)

// pubsub := client.Subscribe(ctx, "mychannel1")
// defer pubsub.Close()

// ch := pubsub.Channel()

// for msg := range ch {
// 	if msg.Payload == "close"{
// 		break
// 	}
// 	fmt.Println(msg.Channel, msg.Payload)
// }
