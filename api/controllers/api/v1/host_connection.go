package v1

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"time"
)

// A lease prevents two host tabs from advancing the same quiz concurrently.
func (qc *quizSocketController) hostConnection(c *websocket.Conn, sessionID string) (func(), bool) {
	client := qc.redis.PubSubModel.Client
	ctx := qc.redis.PubSubModel.Ctx
	key := "host-connection:" + sessionID
	token := uuid.NewString()
	ok, err := client.SetNX(ctx, key, token, 15*time.Second).Result()
	if err != nil || !ok {
		return nil, false
	}
	stop := make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				renewed, err := client.Eval(ctx, `if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('expire',KEYS[1],15) else return 0 end`, []string{key}, token).Int()
				if err != nil || renewed != 1 {
					c.Close()
					return
				}
			}
		}
	}()
	return func() {
		close(stop)
		client.Eval(ctx, `if redis.call('get',KEYS[1])==ARGV[1] then return redis.call('del',KEYS[1]) else return 0 end`, []string{key}, token)
	}, true
}
