package utils

import (
	"strings"
	"testing"

	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/pkg/structs"
)

func testOptions() AIQuestionOptions {
	return AIQuestionOptions{Duration: 60, Points: 1, MaxQuestions: 20}
}

func validAIQuestion() structs.AIQuestion {
	return structs.AIQuestion{
		Question:      "What is the capital of France?",
		QuestionType:  constants.AIQuestionTypeSingle,
		QuestionMedia: constants.MediaText,
		Options:       []string{"Paris", "Rome", "Madrid", "Berlin"},
		OptionsMedia:  constants.MediaText,
		CorrectAnswer: 1,
		Explanation:   "Paris has been the capital since 987.",
	}
}

func TestNormalizeAIQuestionsAcceptsAValidQuestion(t *testing.T) {
	result, err := NormalizeAIQuestions([]structs.AIQuestion{validAIQuestion()}, testOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Questions) != 1 || len(result.Issues) != 0 {
		t.Fatalf("questions = %d, issues = %v", len(result.Questions), result.Issues)
	}

	stored := result.Questions[0]
	if stored.Answers[0] != 1 || stored.Type != constants.SingleAnswer {
		t.Fatalf("answers = %v, type = %d", stored.Answers, stored.Type)
	}
	if len(stored.Options) != 4 || stored.Options["1"] != "Paris" {
		t.Fatalf("options = %v", stored.Options)
	}
	if stored.DurationInSeconds != 60 || stored.Points != 1 {
		t.Fatalf("duration = %d, points = %d", stored.DurationInSeconds, stored.Points)
	}
}

func TestNormalizeAIQuestionsRejectsBadQuestions(t *testing.T) {
	cases := []struct {
		name   string
		mutate func(*structs.AIQuestion)
	}{
		{"empty question text", func(q *structs.AIQuestion) { q.Question = "   " }},
		{"question too long", func(q *structs.AIQuestion) {
			q.Question = strings.Repeat("q", constants.AIMaxQuestionLength+1)
		}},
		{"option too long", func(q *structs.AIQuestion) {
			q.Options = []string{strings.Repeat("o", constants.AIMaxOptionLength+1), "b", "c"}
		}},
		{"blank option", func(q *structs.AIQuestion) { q.Options = []string{"Paris", "  ", "Rome"} }},
		{"too few options", func(q *structs.AIQuestion) { q.Options = []string{"only"} }},
		{"too many options", func(q *structs.AIQuestion) {
			q.Options = []string{"a", "b", "c", "d", "e", "f"}
		}},
		{"duplicate options ignoring case", func(q *structs.AIQuestion) {
			q.Options = []string{"Paris", "paris", "Rome"}
		}},
		{"correct answer zero on a single question", func(q *structs.AIQuestion) { q.CorrectAnswer = 0 }},
		{"correct answer negative", func(q *structs.AIQuestion) { q.CorrectAnswer = -1 }},
		{"correct answer past the option count", func(q *structs.AIQuestion) { q.CorrectAnswer = 5 }},
		{"image question media", func(q *structs.AIQuestion) { q.QuestionMedia = constants.MediaImage }},
		{"unknown question media", func(q *structs.AIQuestion) { q.QuestionMedia = "video" }},
		{"image options media", func(q *structs.AIQuestion) { q.OptionsMedia = constants.MediaImage }},
		{"unknown question type", func(q *structs.AIQuestion) { q.QuestionType = "essay" }},
		{"code media without a resource", func(q *structs.AIQuestion) {
			q.QuestionMedia = constants.MediaCode
			q.Resource = ""
		}},
		{"resource too long", func(q *structs.AIQuestion) {
			q.QuestionMedia = constants.MediaCode
			q.Resource = strings.Repeat("x", constants.AIMaxResourceLength+1)
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			question := validAIQuestion()
			tc.mutate(&question)

			result, err := NormalizeAIQuestions([]structs.AIQuestion{question}, testOptions())
			if err == nil {
				t.Fatalf("accepted %d questions, want a rejection", len(result.Questions))
			}
			if len(result.Issues) == 0 {
				t.Fatal("expected an issue describing the rejection")
			}
		})
	}
}

func TestNormalizeAIQuestionsSurveyMarksEveryOptionCorrect(t *testing.T) {
	question := validAIQuestion()
	question.QuestionType = constants.AIQuestionTypeSurvey
	question.CorrectAnswer = 0

	result, err := NormalizeAIQuestions([]structs.AIQuestion{question}, testOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored := result.Questions[0]
	if stored.Type != constants.Survey {
		t.Fatalf("type = %d, want survey", stored.Type)
	}
	if len(stored.Answers) != 4 {
		t.Fatalf("answers = %v, want one per option", stored.Answers)
	}
	if result.Accepted[0].CorrectAnswer != 0 {
		t.Fatalf("correct answer = %d, want 0 for a survey", result.Accepted[0].CorrectAnswer)
	}
}

func TestNormalizeAIQuestionsDropsDuplicateQuestions(t *testing.T) {
	first := validAIQuestion()
	second := validAIQuestion()
	second.Question = "  WHAT IS   the Capital of France?  "

	result, err := NormalizeAIQuestions([]structs.AIQuestion{first, second}, testOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Questions) != 1 {
		t.Fatalf("kept %d questions, want the duplicate dropped", len(result.Questions))
	}
	if len(result.Issues) != 1 || !strings.Contains(result.Issues[0], constants.ErrAIDuplicateQuestion) {
		t.Fatalf("issues = %v, want a duplicate report", result.Issues)
	}
}

func TestNormalizeAIQuestionsCodeOptionsAreCaseSensitive(t *testing.T) {
	question := validAIQuestion()
	question.OptionsMedia = constants.MediaCode
	question.Options = []string{"Value", "value", "other"}
	question.CorrectAnswer = 1

	result, err := NormalizeAIQuestions([]structs.AIQuestion{question}, testOptions())
	if err != nil {
		t.Fatalf("code options differing only by case should be kept: %v", err)
	}
	if len(result.Questions) != 1 {
		t.Fatalf("questions = %d, want 1", len(result.Questions))
	}
}

func TestNormalizeAIQuestionsTruncatesAtMaxQuestions(t *testing.T) {
	var raw []structs.AIQuestion
	for i := 0; i < 5; i++ {
		question := validAIQuestion()
		question.Question = "Question number " + string(rune('a'+i))
		raw = append(raw, question)
	}

	opts := testOptions()
	opts.MaxQuestions = 3

	result, err := NormalizeAIQuestions(raw, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Questions) != 3 {
		t.Fatalf("questions = %d, want 3", len(result.Questions))
	}
	if len(result.Issues) != 1 || result.Issues[0] != constants.AIQuestionsTruncated {
		t.Fatalf("issues = %v, want a truncation notice", result.Issues)
	}
}

func TestNormalizeAIQuestionsClampsPoints(t *testing.T) {
	opts := testOptions()
	opts.Points = constants.MaximumPoints + 50

	result, err := NormalizeAIQuestions([]structs.AIQuestion{validAIQuestion()}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Questions[0].Points != constants.MaximumPoints {
		t.Fatalf("points = %d, want clamped to %d", result.Questions[0].Points, constants.MaximumPoints)
	}
}

func TestNormalizeAIQuestionsRejectsAnInvalidDuration(t *testing.T) {
	opts := testOptions()
	opts.Duration = 0

	if _, err := NormalizeAIQuestions([]structs.AIQuestion{validAIQuestion()}, opts); err == nil {
		t.Fatal("a zero duration should be rejected")
	}
}

func TestExtractAIGeneration(t *testing.T) {
	object := `{"title":"Geography","description":"A quiz","questions":[{"question":"Capital of France?","options":["Paris","Rome"],"correct_answer":1}]}`

	cases := []struct {
		name      string
		content   string
		wantTitle string
	}{
		{"plain object", object, "Geography"},
		{"fenced with a language tag", "```json\n" + object + "\n```", "Geography"},
		{"fenced without a language tag", "```\n" + object + "\n```", "Geography"},
		{"surrounded by prose", "Here you go:\n" + object + "\nHope that helps!", "Geography"},
		{"bare array falls back", `[{"question":"Capital of France?","options":["Paris","Rome"],"correct_answer":1}]`, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			generation, err := ExtractAIGeneration(tc.content)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if generation.Title != tc.wantTitle {
				t.Fatalf("title = %q, want %q", generation.Title, tc.wantTitle)
			}
			if len(generation.Questions) != 1 {
				t.Fatalf("questions = %d, want 1", len(generation.Questions))
			}
		})
	}
}

func TestExtractAIGenerationRejectsUnusableContent(t *testing.T) {
	for _, content := range []string{"", "   ", "I cannot help with that.", "{}", `{"questions":[]}`, "{not json}"} {
		if _, err := ExtractAIGeneration(content); err == nil {
			t.Fatalf("ExtractAIGeneration(%q) accepted, want error", content)
		}
	}
}

func TestTruncateRunesCountsRunesNotBytes(t *testing.T) {
	if got := TruncateRunes("héllo wörld", 5); got != "héllo" {
		t.Fatalf("got %q, want %q", got, "héllo")
	}

	// A five rune multi byte string is eleven bytes: a byte slice would split it.
	multi := "日本語です"
	if got := TruncateRunes(multi, 5); got != multi {
		t.Fatalf("got %q, want the string unchanged", got)
	}
	if got := TruncateRunes(multi, 2); got != "日本" {
		t.Fatalf("got %q, want %q", got, "日本")
	}
}

func TestSanitizeAITitleFallsBackToTheTopic(t *testing.T) {
	if got := SanitizeAITitle("   ", "world capitals"); got != "world capitals quiz" {
		t.Fatalf("got %q", got)
	}
	if got := SanitizeAITitle("  Spaced   Out  ", "topic"); got != "Spaced Out" {
		t.Fatalf("got %q", got)
	}
	if got := SanitizeAITitle(strings.Repeat("t", 80), "topic"); len([]rune(got)) > constants.AIMaxTitleLength {
		t.Fatalf("title not truncated: %d runes", len([]rune(got)))
	}
}
