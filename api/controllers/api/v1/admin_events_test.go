package v1

import (
	"sync"
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/models"
)

// received reads one value without blocking; ok is false when nothing is waiting.
func received(ch chan bool) (value, ok bool) {
	select {
	case value = <-ch:
		return value, true
	default:
		return false, false
	}
}

// The host's socket reader must never block on a control message nobody is waiting for
// (GitHub #39): a double Skip, repeated scoreboard shortcuts, or a malformed pause.
func TestDispatchAdminEventNeverBlocksTheReader(t *testing.T) {
	next, skip, skipTimer, pause := newAdminEventChan(), newAdminEventChan(), newAdminEventChan(), newAdminEventChan()
	messages := []QuizReceiveResponse{
		{Event: constants.EventForceSkip},
		{Event: constants.EventForceSkip}, // the second click of a double click
		{Event: constants.EventSkipAsked},
		{Event: constants.EventSkipTimer},
		{Event: constants.EventSkipTimer},
		{Event: constants.EventPauseQuiz, Data: true},
		{Event: constants.EventPauseQuiz, Data: "yes"}, // used to panic
		{Event: constants.EventNextQuestionAsked},
		{Event: constants.EventNextQuestionAsked},
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for _, m := range messages {
			dispatchAdminEvent(m, false, next, skip, skipTimer, pause)
		}
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("dispatchAdminEvent blocked with no handler receiving")
	}

	// The first event of each kind waits for a handler that was busy.
	if v, ok := received(skip); !ok || !v {
		t.Fatalf("want the first (force) skip kept, got %v, %v", v, ok)
	}
	if v, ok := received(pause); !ok || !v {
		t.Fatalf("want the pause kept, got %v, %v", v, ok)
	}
	if _, ok := received(skipTimer); !ok {
		t.Fatal("want the skip_timer kept")
	}
	if _, ok := received(next); !ok {
		t.Fatal("want the next_question kept")
	}

	// Solo sessions ignore host pacing.
	dispatchAdminEvent(QuizReceiveResponse{Event: constants.EventForceSkip}, true, next, skip, skipTimer, pause)
	if len(skip) != 0 {
		t.Fatal("solo session accepted a skip")
	}
}

// Pause and resume set a state: the newest one wins while the scoreboard handler is busy.
func TestDispatchAdminEventKeepsNewestPause(t *testing.T) {
	next, skip, skipTimer, pause := newAdminEventChan(), newAdminEventChan(), newAdminEventChan(), newAdminEventChan()
	for _, isPause := range []bool{true, false, true, false} {
		dispatchAdminEvent(QuizReceiveResponse{Event: constants.EventPauseQuiz, Data: isPause}, false, next, skip, skipTimer, pause)
	}
	if v, ok := received(pause); !ok || v {
		t.Fatalf("pause, resume, pause, resume: want resume (false) kept, got %v, %v", v, ok)
	}
}

// A skip_timer or pause left from an earlier phase must not end or pause the next scoreboard.
// Without the drain, the stale skip ends the scoreboard at once, or the stale pause reaches
// shareEvenWithUser with a nil controller and panics.
func TestHandleSkipTimerDropsStaleEvents(t *testing.T) {
	skipTimer, pause := newAdminEventChan(), newAdminEventChan()
	for _, ch := range []chan bool{skipTimer, pause} {
		select {
		case ch <- true:
		default:
			t.Fatal("host control channels must hold one event")
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	start := time.Now()
	handleSkipTimer(nil, nil, &wg, nil, models.ActiveQuiz{}, skipTimer, pause, 1, nil, make(chan struct{}))
	if d := time.Since(start); d < 900*time.Millisecond {
		t.Fatalf("scoreboard ended after %v: a stale skip_timer was honoured", d)
	}
}
