package services_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/services"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// A private quiz shared with another host, who ran it in a room with a player, must not block the
// owner's account deletion (GitHub #46, #70). Needs a migrated database: JOVVIX_TEST_DSN.
func TestDeleteAccountWithSharedQuizPlayedByOthers(t *testing.T) {
	dsn := os.Getenv("JOVVIX_TEST_DSN")
	if dsn == "" {
		t.Skip("requires JOVVIX_TEST_DSN")
	}
	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	n := time.Now().UnixNano() % 1e18
	owner, host, player := fmt.Sprintf("o%019d", n), fmt.Sprintf("h%019d", n), fmt.Sprintf("p%019d", n)
	ownerKratos := uuid.NewString()
	ownerQuiz, hostQuiz, question := uuid.NewString(), uuid.NewString(), uuid.NewString()
	session, earlierRound, room, played := uuid.NewString(), uuid.NewString(), uuid.NewString(), uuid.NewString()

	seed := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO users(id,first_name,last_name,username,email,kratos_id) VALUES
			($1::text,'o','o',left($1::text,1)||right($1::text,11),$4,$5),($2::text,'h','h',left($2::text,1)||right($2::text,11),$6,NULL),($3::text,'p','p',left($3::text,1)||right($3::text,11),NULL,NULL)`,
			[]any{owner, host, player, owner + "@example.com", ownerKratos, host + "@example.com"}},
		{`INSERT INTO quizzes(id,title,creator_id,is_public) VALUES ($1,'owner private',$2,false),($3,'host quiz',$4,false)`,
			[]any{ownerQuiz, owner, hostQuiz, host}},
		{`INSERT INTO questions(id,question,options,answers,duration_in_seconds,type) VALUES ($1,'q','{}','[]',10,1)`,
			[]any{question}},
		{`INSERT INTO quiz_questions(id,quiz_id,question_id) VALUES (gen_random_uuid(),$1,$2)`, []any{ownerQuiz, question}},
		// The host re-shared the quiz, so this row matches neither shared_by nor shared_to of the owner.
		{`INSERT INTO shared_quizzes(id,quiz_id,shared_by,shared_to,permission) VALUES (gen_random_uuid(),$1,$2,'third@example.com','read')`,
			[]any{ownerQuiz, host}},
		{`INSERT INTO active_quizzes(id,quiz_id,admin_id,current_question) VALUES ($1,$2,$3,$4),($5,$6,$3,NULL)`,
			[]any{session, ownerQuiz, host, question, earlierRound, hostQuiz}},
		{`INSERT INTO active_quiz_questions(id,active_quiz_id,question_id,order_no) VALUES (gen_random_uuid(),$1,$2,1)`,
			[]any{session, question}},
		{`INSERT INTO quiz_rooms(id,admin_id,code,current_session_id) VALUES ($1,$2,100000+floor(random()*899999)::int,$3)`,
			[]any{room, host, session}},
		{`UPDATE active_quizzes SET room_id=$1 WHERE id IN ($2,$3)`, []any{room, session, earlierRound}},
		{`INSERT INTO quiz_room_members(room_id,user_id,player_name) VALUES ($1,$2,'player')`, []any{room, player}},
		{`INSERT INTO user_played_quizzes(id,user_id,active_quiz_id) VALUES ($1,$2,$3)`, []any{played, player, session}},
		{`INSERT INTO user_quiz_responses(id,user_played_quiz_id,question_id) VALUES (gen_random_uuid(),$1,$2)`,
			[]any{played, question}},
	}
	// Remove whatever this run left, in foreign-key order, also after a failed seed or deletion.
	t.Cleanup(func() {
		for _, q := range [][]any{
			{`DELETE FROM user_quiz_responses WHERE user_played_quiz_id=$1`, played},
			{`DELETE FROM user_played_quizzes WHERE id=$1`, played},
			{`DELETE FROM quiz_room_members WHERE room_id=$1`, room},
			{`UPDATE active_quizzes SET room_id=NULL WHERE room_id=$1`, room},
			{`DELETE FROM quiz_rooms WHERE id=$1`, room},
			{`DELETE FROM active_quiz_questions WHERE active_quiz_id=$1`, session},
			{`DELETE FROM active_quizzes WHERE id IN ($1,$2)`, session, earlierRound},
			{`DELETE FROM shared_quizzes WHERE quiz_id=$1`, ownerQuiz},
			{`DELETE FROM quiz_questions WHERE quiz_id=$1`, ownerQuiz},
			{`DELETE FROM questions WHERE id=$1`, question},
			{`DELETE FROM quizzes WHERE id IN ($1,$2)`, ownerQuiz, hostQuiz},
			{`DELETE FROM users WHERE id IN ($1,$2,$3)`, owner, host, player},
		} {
			if _, err := sqlDB.Exec(q[0].(string), q[1:]...); err != nil {
				t.Logf("cleanup: %v", err)
			}
		}
	})
	for _, step := range seed {
		_, err := sqlDB.Exec(step.query, step.args...)
		require.NoError(t, err, step.query)
	}

	kratosDeletes := 0
	kratos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && strings.HasSuffix(r.URL.Path, "/identities/"+ownerKratos) {
			kratosDeletes++
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer kratos.Close()

	cfg := config.AppConfig{Kratos: config.KratosConfig{BaseUrl: kratos.URL, AdminUrl: kratos.URL + "/admin"}}
	svc, err := services.NewUserService(goqu.New("postgres", sqlDB), zap.NewNop(), cfg)
	require.NoError(t, err)

	require.NoError(t, svc.DeleteUserDataById(owner, ownerKratos))
	assert.Equal(t, 1, kratosDeletes)

	count := func(query string, a ...any) int {
		var c int
		require.NoError(t, sqlDB.QueryRow(query, a...).Scan(&c), query)
		return c
	}
	assert.Equal(t, 0, count(`SELECT count(*) FROM users WHERE id=$1`, owner))
	assert.Equal(t, 0, count(`SELECT count(*) FROM quizzes WHERE id=$1`, ownerQuiz))
	assert.Equal(t, 0, count(`SELECT count(*) FROM questions WHERE id=$1`, question))
	assert.Equal(t, 0, count(`SELECT count(*) FROM shared_quizzes WHERE quiz_id=$1`, ownerQuiz))
	assert.Equal(t, 0, count(`SELECT count(*) FROM active_quizzes WHERE id=$1`, session))
	assert.Equal(t, 0, count(`SELECT count(*) FROM quiz_rooms WHERE id=$1`, room))
	assert.Equal(t, 0, count(`SELECT count(*) FROM user_played_quizzes WHERE id=$1`, played))
	// Other people's data stays: the host's earlier round is only detached from the closed room.
	assert.Equal(t, 2, count(`SELECT count(*) FROM users WHERE id IN ($1,$2)`, host, player))
	assert.Equal(t, 1, count(`SELECT count(*) FROM active_quizzes WHERE id=$1 AND room_id IS NULL`, earlierRound))
}
