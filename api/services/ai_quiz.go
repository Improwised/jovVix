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

// Pointer and omitempty so a rejected parameter can be dropped and replayed.
type chatCompletionRequest struct {
	Model               string              `json:"model"`
	Messages            []chatMessage       `json:"messages"`
	Temperature         *float64            `json:"temperature,omitempty"`
	MaxTokens           int                 `json:"max_tokens,omitempty"`
	MaxCompletionTokens int                 `json:"max_completion_tokens,omitempty"`
	ResponseFormat      *chatResponseFormat `json:"response_format,omitempty"`
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

// Matched loosely: every provider words a parameter refusal differently.
var (
	aiJSONModeMarkers = []string{
		"response_format",
		"json_object",
		"json mode",
		"json_schema",
		"structured output",
	}
	aiMaxTokensMarkers = []string{
		"max_tokens",
		"max_completion_tokens",
	}
	aiTemperatureMarkers = []string{
		"temperature",
	}
)

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
		SetTransport(utils.NewAIProviderTransport()).
		SetRedirectPolicy(resty.NoRedirectPolicy())

	return &AIQuizService{
		client: client,
		logger: logger,
		config: aiConfig,
	}
}

// avoid holds the quiz's existing questions, empty for a new quiz.
func (svc *AIQuizService) GenerateQuestions(ctx context.Context, cred AICredentials, req structs.ReqGenerateAIQuestions, avoid []string) (utils.AIGeneration, error) {
	if cred.BaseURL == "" || cred.Model == "" {
		return utils.AIGeneration{}, errors.New(constants.ErrAINotConfigured)
	}

	systemPrompt, userPrompt := buildAIPrompts(req, avoid)

	body := svc.completionRequest(cred, systemPrompt, userPrompt, responseTokenBudget(req.NumberOfQuestions))

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

// Asks for a real question: a model can be reachable and still fail the JSON
// contract. Returns the question it produced, shown back as proof.
func (svc *AIQuizService) TestConnection(ctx context.Context, cred AICredentials, opts utils.AIQuestionOptions) (time.Duration, string, error) {
	if cred.BaseURL == "" || cred.Model == "" {
		return 0, "", errors.New(constants.ErrAIIncompleteCredentials)
	}

	systemPrompt, userPrompt := buildAIPrompts(structs.ReqGenerateAIQuestions{
		Topic:             constants.AITestTopic,
		NumberOfQuestions: 1,
		Difficulty:        constants.AIDifficultyEasy,
		Language:          constants.AIDefaultLanguage,
	}, nil)

	body := svc.completionRequest(cred, systemPrompt, userPrompt, constants.AIMinResponseTokens)

	started := time.Now()
	completion, err := svc.postCompletion(ctx, cred, body)
	elapsed := time.Since(started)

	if err != nil {
		return elapsed, "", err
	}

	content := strings.TrimSpace(completion.Choices[0].Message.Content)
	generation, err := utils.ExtractAIGeneration(content)
	if err != nil {
		svc.logger.Warn("ai connection test reached the model but could not parse its reply",
			zap.String("model", cred.Model),
			zap.String("finish_reason", completion.Choices[0].FinishReason),
			zap.String("content", truncateForLog(content)))
		return elapsed, "", errors.New(constants.ErrAITestNoValidQuestion)
	}

	result, err := utils.NormalizeAIQuestions(generation.Questions, opts)
	if err != nil || len(result.Accepted) == 0 {
		svc.logger.Warn("ai connection test produced no usable question",
			zap.String("model", cred.Model),
			zap.Strings("issues", result.Issues))
		return elapsed, "", errors.New(constants.ErrAITestNoValidQuestion)
	}

	return elapsed, result.Accepted[0].Question, nil
}

func (svc *AIQuizService) ListModels(ctx context.Context, cred AICredentials) (structs.ResAIModels, error) {
	var listing structs.ResAIModels

	if cred.BaseURL == "" {
		return listing, errors.New(constants.ErrAIIncompleteCredentials)
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
		return listing, errors.New(utils.ClassifyAIError(0, "", err))
	}

	if res.StatusCode() != http.StatusOK {
		svc.logger.Warn("model listing returned an unexpected status",
			zap.Int("status", res.StatusCode()),
			zap.String("body", constants.AIProviderBodyRedacted))
		return listing, errors.New(utils.ClassifyAIError(res.StatusCode(), res.String(), nil))
	}

	listing = collectModelIDs(out.Data)
	if len(listing.Models) == 0 {
		return listing, errors.New(constants.ErrAIModelsUnavailable)
	}

	svc.logger.Info("listed models from the provider",
		zap.Int("offered", len(out.Data)),
		zap.Int("chat_capable", len(listing.Models)),
		zap.Int("filtered", listing.Filtered),
		zap.Bool("truncated", listing.Truncated))

	return listing, nil
}

// Keeps only what /chat/completions can serve. Filtered counts the drops so a
// shortened list is not presented as the whole inventory.
func collectModelIDs(entries []modelListEntry) structs.ResAIModels {
	seen := make(map[string]struct{}, len(entries))
	models := make([]string, 0, len(entries))
	filtered := 0

	for _, entry := range entries {
		id, err := utils.ValidateAIModel(entry.ID)
		if err != nil {
			continue
		}
		if _, duplicate := seen[id]; duplicate {
			continue
		}
		seen[id] = struct{}{}

		if !utils.IsChatCapableModel(id) {
			filtered++
			continue
		}

		models = append(models, id)
	}

	sort.Strings(models)

	truncated := len(models) > constants.AIMaxListedModels
	if truncated {
		models = models[:constants.AIMaxListedModels]
	}

	return structs.ResAIModels{
		Models:    models,
		Filtered:  filtered,
		Truncated: truncated,
	}
}

func responseTokenBudget(numberOfQuestions int) int {
	budget := numberOfQuestions * constants.AIResponseTokenBudget
	if budget < constants.AIMinResponseTokens {
		return constants.AIMinResponseTokens
	}
	return budget
}

func (svc *AIQuizService) completionRequest(cred AICredentials, systemPrompt, userPrompt string, maxTokens int) chatCompletionRequest {
	temperature := svc.config.Temp()

	body := chatCompletionRequest{
		Model:       cred.Model,
		Temperature: &temperature,
		MaxTokens:   maxTokens,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
	}

	if svc.config.JSONMode {
		body.ResponseFormat = &chatResponseFormat{Type: "json_object"}
	}

	return body
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
	for attempt := 0; attempt <= constants.AIMaxRequestAdaptations; attempt++ {
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

			if isAIParameterRejection(res.StatusCode()) {
				if adaptation := adaptAIRequest(&body, res.String()); adaptation != "" {
					svc.logger.Warn("retrying the ai completion without a parameter the provider rejected",
						zap.String("model", body.Model),
						zap.String("adaptation", adaptation))
					continue
				}
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
				zap.Int("max_completion_tokens", body.MaxCompletionTokens),
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

	return nil, errors.New(constants.ErrAIRequestFailed)
}

func isAIParameterRejection(status int) bool {
	return status == http.StatusBadRequest || status == http.StatusUnprocessableEntity
}

// Drops or swaps the parameter the provider rejected and reports which. Each
// branch needs the parameter still set, so the ladder cannot loop.
func adaptAIRequest(body *chatCompletionRequest, providerBody string) string {
	lowered := strings.ToLower(providerBody)

	if body.ResponseFormat != nil && containsAny(lowered, aiJSONModeMarkers) {
		body.ResponseFormat = nil
		return constants.AIAdaptationDroppedJSONMode
	}

	if body.MaxTokens > 0 && containsAny(lowered, aiMaxTokensMarkers) {
		body.MaxCompletionTokens = body.MaxTokens
		body.MaxTokens = 0
		return constants.AIAdaptationMaxCompletionTokens
	}

	if body.Temperature != nil && containsAny(lowered, aiTemperatureMarkers) {
		body.Temperature = nil
		return constants.AIAdaptationDroppedTemperature
	}

	return ""
}

func containsAny(haystack string, needles []string) bool {
	for _, needle := range needles {
		if strings.Contains(haystack, needle) {
			return true
		}
	}
	return false
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
