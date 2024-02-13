package components

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type quizConfigs struct {
	user_id         string
	role            string
	quiz_id         string
	session_id      string
	user_session_id string
}

type QuizSessionCfg struct {
	id               int
	code             string
	session_name     string
	questions        []Question
	current_question int
	quiz_id          Quiz
	session_start_at time.Time
	session_end_at   time.Time
	admin            []UserConfig
	player           []UserConfig
	mutex            sync.Mutex
}

type User struct {
	user_id string
	name    string
	role    string
}

type UserConfig struct {
	conn            *websocket.Conn
	user            User
	user_session_id string
}

// Quiz structure
type Quiz struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Question structure
type Question struct {
	ID       int64    `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answers  []int    `json:"answers"` // Indices of correct options
}

type SessionManager struct {
	sessions map[string]*QuizSessionCfg
	series   int
	mutex    sync.Mutex // Synchronization for thread-safety
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*QuizSessionCfg),
		series:   1,
	}
}

func (sm *SessionManager) GetCode() (int, string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	code := generateRandomString(8)

	for {
		if _, is_available := sm.sessions[code]; !is_available {
			break
		}
		code = generateRandomString(8)
	}

	sm.series += 1

	return sm.series, code
}

func generateID() int64 {
	// Replace with your actual ID generation logic (e.g., using a unique ID generator library)
	return int64(time.Now().UnixNano())
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := big.NewInt(int64(len(charset)))
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			panic(err) // Handle error
		}
		result[i] = charset[randomIndex.Int64()]
	}

	return string(result)
}

func (sm *SessionManager) AddSession(quizID int64, name string, admin []UserConfig, start time.Time, maxDuration time.Duration) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	no, code := sm.GetCode()

	session := &QuizSessionCfg{
		id:               no,
		code:             code,
		session_name:     name,
		questions:        []Question{},
		session_start_at: start,
		session_end_at:   start.Add(maxDuration),
		admin:            admin,
	}
	sm.sessions[code] = session
	return nil
}

func (sm *SessionManager) RemoveSession(code string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	_, ok := sm.sessions[code]
	if !ok {
		return errors.New("session not found")
	}

	delete(sm.sessions, code)
	return nil
}
