package v1

import (
	"testing"
	"time"

	"github.com/Improwised/jovvix/api/constants"
)

// The host's socket reader must never block on a control message nobody is waiting for
// (GitHub #39): a double Skip, repeated scoreboard shortcuts, or a malformed pause.
func TestDispatchAdminEventNeverBlocksTheReader(t *testing.T) {
	next, skip, skipTimer, pause := make(chan bool, 1), make(chan bool, 1), make(chan bool, 1), make(chan bool, 1)
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
	if v := <-skip; !v {
		t.Fatal("want the first (force) skip kept")
	}
	if v := <-pause; !v {
		t.Fatal("want the pause kept")
	}

	// A handler that starts a new phase drops what an earlier phase left.
	dispatchAdminEvent(QuizReceiveResponse{Event: constants.EventForceSkip}, false, next, skip, skipTimer, pause)
	drainStaleAdminEvent(skip)
	select {
	case <-skip:
		t.Fatal("a skip from the previous question survived the drain")
	default:
	}

	// Solo sessions ignore host pacing.
	dispatchAdminEvent(QuizReceiveResponse{Event: constants.EventForceSkip}, true, next, skip, skipTimer, pause)
	if len(skip) != 0 {
		t.Fatal("solo session accepted a skip")
	}
}
