package v1

import (
	"reflect"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/google/uuid"
)

type playerOptionOrderProviderFunc func(userPlayedQuizID, questionID uuid.UUID, options map[string]string) (map[string]string, error)

func (f playerOptionOrderProviderFunc) GetOptionsForPlayer(userPlayedQuizID, questionID uuid.UUID, options map[string]string) (map[string]string, error) {
	return f(userPlayedQuizID, questionID, options)
}

func (f playerOptionOrderProviderFunc) GetOptionsAndAnswerKeysForPlayer(userPlayedQuizID, questionID uuid.UUID, options map[string]string, answerKeys []int) (map[string]string, []int, error) {
	orderedOptions, err := f.GetOptionsForPlayer(userPlayedQuizID, questionID, options)
	if err != nil {
		return nil, nil, err
	}
	displayedKeyByOriginal := map[int]int{1: 2, 2: 3, 3: 1}
	orderedAnswerKeys := make([]int, len(answerKeys))
	for index, originalKey := range answerKeys {
		orderedAnswerKeys[index] = displayedKeyByOriginal[originalKey]
	}
	return orderedOptions, orderedAnswerKeys, nil
}

func TestApplyPlayerOptionOrderForEventReordersScoreOptions(t *testing.T) {
	playerID := uuid.New()
	questionID := uuid.New()
	called := false
	scoreData := map[string]any{
		"id": questionID.String(),
		"options": map[string]any{
			"1": "Blue",
			"2": "Red",
			"3": "Green",
		},
		"answers": []any{2.0},
	}
	message := map[string]any{
		"response": map[string]any{"data": scoreData},
	}
	provider := playerOptionOrderProviderFunc(func(gotPlayerID, gotQuestionID uuid.UUID, options map[string]string) (map[string]string, error) {
		called = true
		if gotPlayerID != playerID {
			t.Errorf("got player ID %s, want %s", gotPlayerID, playerID)
		}
		if gotQuestionID != questionID {
			t.Errorf("got question ID %s, want %s", gotQuestionID, questionID)
		}
		if !reflect.DeepEqual(options, map[string]string{"1": "Blue", "2": "Red", "3": "Green"}) {
			t.Errorf("got canonical options %#v", options)
		}
		return map[string]string{"1": options["3"], "2": options["1"], "3": options["2"]}, nil
	})

	if err := applyPlayerOptionOrderForEvent(message, constants.EventShowScore, provider, playerID); err != nil {
		t.Fatalf("applyPlayerOptionOrderForEvent returned an error: %v", err)
	}
	if !called {
		t.Fatal("option order was not requested for the show_score event")
	}

	want := map[string]string{"1": "Green", "2": "Blue", "3": "Red"}
	got, ok := scoreData["options"].(map[string]string)
	if !ok {
		t.Fatalf("score options have type %T, want map[string]string", scoreData["options"])
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got score options %#v, want %#v", got, want)
	}
	if gotAnswers, ok := scoreData["answers"].([]int); !ok || !reflect.DeepEqual(gotAnswers, []int{3}) {
		t.Fatalf("got displayed answer keys %#v, want []int{3}", scoreData["answers"])
	}
}

func TestApplyPlayerOptionOrderForEventSkipsAdminsAndUnrelatedEvents(t *testing.T) {
	called := false
	provider := playerOptionOrderProviderFunc(func(uuid.UUID, uuid.UUID, map[string]string) (map[string]string, error) {
		called = true
		return nil, nil
	})

	for _, tc := range []struct {
		name   string
		event  string
		userID uuid.UUID
	}{
		{name: "admin", event: constants.EventShowScore, userID: uuid.Nil},
		{name: "unrelated event", event: constants.EventTerminateQuiz, userID: uuid.New()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := applyPlayerOptionOrderForEvent(nil, tc.event, provider, tc.userID); err != nil {
				t.Fatalf("applyPlayerOptionOrderForEvent returned an error: %v", err)
			}
		})
	}
	if called {
		t.Fatal("option order provider was called for an admin or unrelated event")
	}
}
