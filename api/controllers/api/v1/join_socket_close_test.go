package v1

import (
	"context"
	"database/sql"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
	jvredis "github.com/Improwised/jovvix/api/pkg/redis"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"
	"github.com/fasthttp/websocket"
	fws "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// A host who opens the player socket gets a redirect, and Join returns while its
// reader goroutine still runs. A close frame that arrives at that moment must not
// reach a connection the server has already released (it crashed the API).
func TestJoinSurvivesCloseFrameAfterReturn(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:joinclose?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if _, err = db.Exec(`CREATE TABLE active_quizzes(id text, invitation_code integer, title text, quiz_id text,
		admin_id text, activated_to timestamp, activated_from timestamp, is_active boolean, quiz_analysis text,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP, updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
		current_question text, is_question_active boolean, question_delivery_time timestamp, room_id text,
		is_solo boolean DEFAULT false, quiz_version text DEFAULT '', solo_completed boolean DEFAULT false)`); err != nil {
		t.Fatal(err)
	}
	if _, err = db.Exec(`INSERT INTO active_quizzes(id, invitation_code, title, quiz_id, admin_id, is_active)
		VALUES (?, 123456, 'q', ?, 'host', true)`, uuid.NewString(), uuid.NewString()); err != nil {
		t.Fatal(err)
	}
	gdb := goqu.New("sqlite3", db)
	// Redis is unreachable on purpose: the reader's exit path only logs the error.
	rc := goredis.NewClient(&goredis.Options{Addr: "127.0.0.1:1", MaxRetries: -1, DialTimeout: 50 * time.Millisecond})
	defer rc.Close()
	qc := &quizSocketController{
		db:              gdb,
		activeQuizModel: models.InitActiveQuizModel(gdb, zap.NewNop()),
		logger:          zap.NewNop(),
		redis:           &jvredis.RedisPubSub{PubSubModel: &jvredis.PubSubModel{Ctx: context.Background(), Client: *rc}},
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/join", func(c *fiber.Ctx) error {
		c.Locals(constants.ContextUid, "host")
		c.Locals(constants.ContextUser, models.User{ID: "host"})
		c.Locals(constants.QuizSessionInvitationCode, "123456")
		return c.Next()
	}, fws.New(qc.Join))
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() { _ = app.Listener(ln) }()
	defer func() { _ = app.Shutdown() }()
	url := "ws://" + ln.Addr().String() + "/join"

	closeFrame := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cl, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				return
			}
			defer cl.Close()
			_ = cl.WriteControl(websocket.CloseMessage, closeFrame, time.Now().Add(time.Second))
			_ = cl.SetReadDeadline(time.Now().Add(2 * time.Second))
			for {
				if _, _, err := cl.ReadMessage(); err != nil {
					return
				}
			}
		}()
		if i%50 == 49 {
			time.Sleep(5 * time.Millisecond)
		}
	}
	wg.Wait()
}
