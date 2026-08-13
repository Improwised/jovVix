package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/Improwised/jovvix/api/constants"
	"github.com/Improwised/jovvix/api/pkg/structs"
	"github.com/Improwised/jovvix/api/utils"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type AIQuizService struct {
	client *resty.Client
	logger *zap.Logger
	config *config.AIConfig
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponseFormat struct {
	Type string `json:"type"`
}

type chatCompletionRequest struct {
	Model          string              `json:"model"`
	Messages       []chatMessage       `json:"messages"`
	Temperature    float64             `json:"temperature"`
	MaxTokens      int                 `json:"max_tokens,omitempty"`
	ResponseFormat *chatResponseFormat `json:"response_format,omitempty"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type modelListEntry struct {
	ID string `json:"id"`
}

type modelListResponse struct {
	Data []modelListEntry `json:"data"`
}

var aiLanguageNames = map[string]string{
	constants.AIDefaultLanguage: "English",
	"hindi":                     "Hindi",
	"gujarati":                  "Gujarati",
	"marathi":                   "Marathi",
	"bengali":                   "Bengali",
	"tamil":                     "Tamil",
	"telugu":                    "Telugu",
	"spanish":                   "Spanish",
	"french":                    "French",
	"german":                    "German",
	"portuguese":                "Portuguese",
	"arabic":                    "Arabic",
	"chinese":                   "Simplified Chinese",
	"japanese":                  "Japanese",
}

var aiDifficultyGuidance = map[string]string{
	constants.AIDifficultyEasy:   "common knowledge that a beginner would recognise; single-step recall",
	constants.AIDifficultyMedium: "requires real familiarity with the topic and one step of reasoning",
	constants.AIDifficultyHard:   "requires specialist knowledge, precise detail, or multi-step reasoning; distractors should be near-misses",
}

const aiSystemPrompt = `You are a quiz author. You write factual, self-contained multiple-choice questions for a live quiz game.
You reply with a single JSON object and nothing else: no markdown, no code fences, no commentary.
Code, when you use it, goes in a dedicated JSON string field with newlines escaped as \n, never in markdown fences.`

const aiAvoidPromptHeader = `

This quiz already has the questions below. Do not repeat any of them, and do not ask the same fact in different words:`

const aiUserPromptTemplate = `Create %d multiple-choice quiz questions about "%s" at %s difficulty.

Rules:
- Each question must have exactly %d options.
- Exactly one option is correct.
- "correct_answer" is the 1-based position of the correct option in the "options" array.
- Options must be plausible, mutually exclusive, and similar in length. Never use "All of the above" or "None of the above".
- Do not number the questions. Do not prefix options with letters or numbers.
- No two questions may test the same fact, and no two options within a question may be identical.
- Each question must be answerable in about 30 seconds with no external material.
- Write every question, option, explanation, and the quiz title and description in %s. Use plain text only: no markdown, no LaTeX, no HTML, no images.
- The JSON keys, and the values of "question_type", "question_media" and "options_media", always stay in English exactly as specified below, whatever language the questions are written in.

Code (use only when the topic is about programming, otherwise set both media fields to "text"):
- "question_media" is "text" or "code". No other value is ever allowed.
- Use "code" when the question is about a specific snippet: put the raw snippet in "resource" and keep "question" a plain-text prompt such as "What does this function return?". "question" itself is never code.
- When "question_media" is "text", "resource" must be the empty string "".
- "options_media" is "text" or "code". Use "code" only when every option is itself a snippet, and then each option string is the raw code. One value covers all options of a question.
- Snippets are raw source: no markdown fences, no line numbers, first line flush-left, newlines escaped as \n so the reply stays valid JSON.
- Keep "resource" under 800 characters and each code option under 300 characters. Keep text questions under 200 characters and text options under 100 characters.
- Mix code and text questions naturally rather than making every question a snippet.

Question type:
- "question_type" is "single" or "survey". No other value is ever allowed.
- "single" is the default: exactly one option is factually correct.
- "survey" is for a genuine opinion or preference question where no option is right or wrong, such as which approach the player prefers. Every option counts as correct and everyone scores.
- Set "correct_answer" to 0 for a survey question, and write the "explanation" as one sentence about why the question is interesting rather than why an option is right.
- Use survey sparingly: at most one in every five questions, and only when the topic genuinely has an opinion to canvass. Never dress up a factual question as a survey.

Difficulty guidance: %s

Also invent a quiz "title" of at most 50 characters and a one-sentence "description" of at most 150 characters, both plain text.

Reply with exactly this JSON shape and no other keys:
{"title":"...","description":"...","questions":[{"question":"...","question_type":"single","question_media":"text","resource":"","options":["...","...","...","..."],"options_media":"text","correct_answer":1,"explanation":"one short sentence saying why the correct option is right"}]}

Every question object must contain all eight keys, even when the value is "". The "questions" array must contain exactly %d objects.`

func NewAIQuizService(logger *zap.Logger, aiConfig *config.AIConfig) *AIQuizService {
	client := resty.New().
		SetTimeout(aiConfig.Timeout()).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetTransport(utils.NewSSRFGuardedTransport()).
		SetRedirectPolicy(resty.NoRedirectPolicy())

	return &AIQuizService{
		client: client,
		logger: logger,
		config: aiConfig,
	}
}

// avoid carries the questions the target quiz already has, so a generation for an
// existing quiz does not repeat what is in it. Empty for a brand new quiz.
func (svc *AIQuizService) GenerateQuestions(ctx context.Context, cred AICredentials, req structs.ReqGenerateAIQuestions, avoid []string) (utils.AIGeneration, error) {
	if cred.BaseURL == "" || cred.Model == "" {
		return utils.AIGeneration{}, errors.New(constants.ErrAINotConfigured)
	}

	systemPrompt, userPrompt := buildAIPrompts(req, avoid)

	body := chatCompletionRequest{
		Model:       cred.Model,
		Temperature: svc.config.Temp(),
		MaxTokens:   responseTokenBudget(req.NumberOfQuestions),
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}
	if svc.config.JSONMode {
		body.ResponseFormat = &chatResponseFormat{Type: "json_object"}
	}

	completion, err := svc.postCompletion(ctx, cred, body)
	if err != nil {
		return utils.AIGeneration{}, err
	}

	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	generation, err := utils.ExtractAIGeneration(content)
	if err != nil {
		svc.logger.Error("could not parse ai completion content",
			zap.String("content", truncateForLog(content)),
			zap.String("finish_reason", completion.Choices[0].FinishReason),
			zap.Error(err))
		return utils.AIGeneration{}, errors.New(constants.ErrAIInvalidResponse)
	}

	return generation, nil
}

func (svc *AIQuizService) TestConnection(ctx context.Context, cred AICredentials) (time.Duration, error) {
	if cred.BaseURL == "" || cred.Model == "" {
		return 0, errors.New(constants.ErrAIIncompleteCredentials)
	}

	body := chatCompletionRequest{
		Model:       cred.Model,
		Temperature: 0,
		MaxTokens:   1,
		Messages:    []chatMessage{{Role: "user", Content: "ping"}},
	}

	var out chatCompletionResponse
	started := time.Now()
	res, err := svc.doCompletion(ctx, cred, body, &out)
	elapsed := time.Since(started)

	if err != nil {
		svc.logger.Warn("ai connection test could not reach the provider",
			zap.Error(err))
		return elapsed, errors.New(utils.ClassifyAIError(0, "", err))
	}

	if res.StatusCode() != http.StatusOK {
		svc.logger.Warn("ai connection test returned an unexpected status",
			zap.Int("status", res.StatusCode()),
			zap.String("body", constants.AIProviderBodyRedacted))
		return elapsed, errors.New(utils.ClassifyAIError(res.StatusCode(), res.String(), nil))
	}

	if out.Error != nil {
		return elapsed, errors.New(utils.ClassifyAIError(http.StatusBadRequest, res.String(), nil))
	}

	return elapsed, nil
}

func (svc *AIQuizService) ListModels(ctx context.Context, cred AICredentials) ([]string, error) {
	if cred.BaseURL == "" {
		return nil, errors.New(constants.ErrAIIncompleteCredentials)
	}

	var out modelListResponse
	request := svc.client.R().
		SetContext(ctx).
		SetResult(&out)

	if cred.APIKey != "" {
		request.SetAuthToken(cred.APIKey)
	}

	res, err := request.Get(cred.ModelsURL())
	if err != nil {
		svc.logger.Warn("could not list models from the provider",
			zap.Error(err))
		return nil, errors.New(utils.ClassifyAIError(0, "", err))
	}

	if res.StatusCode() != http.StatusOK {
		svc.logger.Warn("model listing returned an unexpected status",
			zap.Int("status", res.StatusCode()),
			zap.String("body", constants.AIProviderBodyRedacted))
		return nil, errors.New(utils.ClassifyAIError(res.StatusCode(), res.String(), nil))
	}

	models := collectModelIDs(out.Data)
	if len(models) == 0 {
		return nil, errors.New(constants.ErrAIModelsUnavailable)
	}

	return models, nil
}

func collectModelIDs(entries []modelListEntry) []string {
	seen := make(map[string]struct{}, len(entries))
	models := make([]string, 0, len(entries))

	for _, entry := range entries {
		id, err := utils.ValidateAIModel(entry.ID)
		if err != nil {
			continue
		}
		if _, duplicate := seen[id]; duplicate {
			continue
		}
		seen[id] = struct{}{}
		models = append(models, id)
	}

	sort.Strings(models)
	if len(models) > constants.AIMaxListedModels {
		models = models[:constants.AIMaxListedModels]
	}
	return models
}

func responseTokenBudget(numberOfQuestions int) int {
	budget := numberOfQuestions * constants.AIResponseTokenBudget
	if budget < constants.AIMinResponseTokens {
		return constants.AIMinResponseTokens
	}
	return budget
}

func (svc *AIQuizService) doCompletion(ctx context.Context, cred AICredentials, body chatCompletionRequest, out *chatCompletionResponse) (*resty.Response, error) {
	request := svc.client.R().
		SetContext(ctx).
		SetBody(body).
		SetResult(out)

	if cred.APIKey != "" {
		request.SetAuthToken(cred.APIKey)
	}

	return request.Post(cred.CompletionsURL())
}

func (svc *AIQuizService) postCompletion(ctx context.Context, cred AICredentials, body chatCompletionRequest) (*chatCompletionResponse, error) {
	var out chatCompletionResponse

	res, err := svc.doCompletion(ctx, cred, body, &out)

	if err != nil {
		svc.logger.Error("ai completion request failed",
			zap.String("model", body.Model),
			zap.Bool("timeout", errors.Is(err, context.DeadlineExceeded)),
			zap.Error(err))
		return nil, errors.New(utils.ClassifyAIError(0, "", err))
	}

	if res.StatusCode() != http.StatusOK {
		svc.logger.Error("ai completion returned an unexpected status",
			zap.Int("status", res.StatusCode()),
			zap.String("model", body.Model),
			zap.String("body", constants.AIProviderBodyRedacted))

		if res.StatusCode() == http.StatusBadRequest &&
			body.ResponseFormat != nil &&
			strings.Contains(strings.ToLower(res.String()), "response_format") {
			body.ResponseFormat = nil
			return svc.postCompletion(ctx, cred, body)
		}

		return nil, errors.New(utils.ClassifyAIError(res.StatusCode(), res.String(), nil))
	}

	if out.Error != nil {
		svc.logger.Error("ai completion returned an error envelope",
			zap.String("model", body.Model),
			zap.String("type", out.Error.Type),
			zap.String("message", constants.AIProviderBodyRedacted))
		return nil, errors.New(utils.ClassifyAIError(http.StatusBadRequest, res.String(), nil))
	}

	if len(out.Choices) == 0 || strings.TrimSpace(out.Choices[0].Message.Content) == "" {
		svc.logger.Error("ai completion returned no content",
			zap.String("model", body.Model),
			zap.String("body", constants.AIProviderBodyRedacted))
		return nil, errors.New(constants.ErrAIEmptyResponse)
	}

	if out.Choices[0].FinishReason == "length" {
		svc.logger.Warn("ai completion hit the token cap and is probably truncated",
			zap.String("model", body.Model),
			zap.Int("max_tokens", body.MaxTokens),
			zap.Int("completion_tokens", out.Usage.CompletionTokens))
	}

	svc.logger.Info("ai completion succeeded",
		zap.String("model", body.Model),
		zap.String("finish_reason", out.Choices[0].FinishReason),
		zap.Int("prompt_tokens", out.Usage.PromptTokens),
		zap.Int("completion_tokens", out.Usage.CompletionTokens),
		zap.Int("total_tokens", out.Usage.TotalTokens))

	return &out, nil
}

func buildAIPrompts(req structs.ReqGenerateAIQuestions, avoid []string) (string, string) {
	guidance, ok := aiDifficultyGuidance[req.Difficulty]
	if !ok {
		guidance = aiDifficultyGuidance[constants.AIDifficultyMedium]
	}

	language, ok := aiLanguageNames[req.Language]
	if !ok {
		language = aiLanguageNames[constants.AIDefaultLanguage]
	}

	userPrompt := fmt.Sprintf(
		aiUserPromptTemplate,
		req.NumberOfQuestions,
		req.Topic,
		req.Difficulty,
		constants.AIDefaultOptionsPerQuestion,
		language,
		guidance,
		req.NumberOfQuestions,
	)

	if block := buildAvoidBlock(avoid); block != "" {
		userPrompt += block
	}

	return aiSystemPrompt, userPrompt
}

func buildAvoidBlock(avoid []string) string {
	var builder strings.Builder
	listed := 0

	for _, question := range avoid {
		if listed >= constants.AIMaxAvoidQuestions {
			break
		}

		trimmed := strings.Join(strings.Fields(question), " ")
		if trimmed == "" {
			continue
		}

		builder.WriteString("\n- ")
		builder.WriteString(utils.TruncateRunes(trimmed, constants.AIMaxAvoidQuestionLength))
		listed++
	}

	if listed == 0 {
		return ""
	}

	return aiAvoidPromptHeader + builder.String()
}

func truncateForLog(value string) string {
	if len(value) <= constants.AIProviderLogBodyLimit {
		return value
	}
	return value[:constants.AIProviderLogBodyLimit] + "..."
}
