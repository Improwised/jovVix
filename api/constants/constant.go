package constants

// variables
const (
	CookieUser   = "user"
	KratosCookie = "ory_kratos_session"
)

// fiber contexts
const (
	ContextUid            = "userId"
	ContextUser           = "userContext"
	ContextQuizPermission = "quiz_permission"
)

// kratos
const (
	KratosID          = "kratosId"
	KratosUserDetails = "kratosUserDetails"
)

// params
const (
	ParamUid         = "userId"
	UserPlayedQuizId = "user_played_quiz_id"
	Username         = "username"
	MediaQuery       = "media"
	ParamTitle       = "title"
)

// Permissions
const (
	ReadPermission  = "read"
	WritePermission = "write"
	SharePermission = "share"
)

// Email templetes
const (
	QuizEmailSubject = "You have been invited to a quiz"
)

// Success messages
// ...

// Fail messages
// ...
const (
	UsernameRequired = "username is required"
)

// Error messages
const (
	ErrGetUser                     = "error while getting user"
	ErrLoginUser                   = "error while logging in the user"
	ErrInsertUser                  = "error while creating user, please try after sometime"
	ErrDeleteUser                  = "error while deleting user, please try after sometime"
	ErrConvertTypeUser             = "Unable to convert to user-model type from locals"
	ErrHealthCheckDb               = "error while checking health of database"
	ErrUnauthenticated             = "error verifying user identity"
	ErrUnauthorized                = "access denied. You do not have the necessary permissions."
	ErrShareNotAllowed             = "you do not have permission to share this quiz"
	ErrCannotChangeOwnAccess       = "you cannot change your own access to this quiz"
	ErrSharedQuizNotFound          = "shared quiz permission not found"
	ErrGetSharedQuiz               = "error while getting shared quiz permission"
	ErrKratosAuth                  = "error while fetching user from kratos"
	ErrKratosDataInsertion         = "error while inserting user data came from kratos"
	ErrKratosIDEmpty               = "error no session_id found in kratos cookie"
	ErrKratosCookieTime            = "error while parsing the expiration time of the cookie"
	ErrRegisterQuiz                = "error while creating quiz"
	ErrCreatingDemoQuiz            = "error while creating demo quiz"
	ErrQuizNotPublic               = "this quiz is not public"
	ErrGetTotalJoinUser            = "error while getting count of user joined in quiz"
	ErrShareQuiz                   = "error while sharing quiz"
	ErrListShareQuiz               = "error while getting list of shared quizzes"
	ErrFetchAuthorizedUsersError   = "Error fetching authorized users for the selected quiz."
	ErrCheckQuizCreatorExists      = "error while checking quiz creator exists or not"
	ErrGetQuizPermission           = "error while getting quiz pemrmission for user"
	ErrUpdateUserPermissionForQuiz = "error while updating user permission for particular quiz"
	ErrDeleteUserPermissionForQuiz = "error while deleting user permission for particular quiz"
	ErrGetStreakCount              = "error while getting streaks count"
	ErrListCategories              = "error while listing quiz categories"
	ErrCreateCategory              = "error while creating quiz category"
	ErrUpdateCategory              = "error while updating quiz category"
	ErrDeleteCategory              = "error while deleting quiz category"
	ErrCategoryNotFound            = "quiz category not found"
	ErrCategoryAlreadyExists       = "a category with this name already exists"
	ErrInvalidCoverImage           = "cover image must be an image"
	ErrCoverImageTooLarge          = "cover image is too large"
)

// Bad Request Message
const (
	BadRequestSharedQuizIdNotFound = "no shared_quiz_id found"
)

// default Events
const (
	EventUserRegistered = "event:userRegistered"
)

// Middleware
const (
	// socket
	MiddlewareError     = "middleware_error"
	ErrorTypeConversion = "type conversion failed"

	// http/https
	ErrNotAllowed              = "not allowed to access Resource"
	ErrUserRequiredToCheckRole = "user not logged in"

	// csv
	FileName                    = "file_name"
	MaxRows                     = 500
	FileSize                    = 100000 // TODO: change file size, ~1mb
	MaximumPoints               = 20
	MinimumPoints               = 0
	MaximumDurationInSeconds    = 300
	MinimumDurationInSeconds    = 1
	DefaultDurationInSeconds    = 30
	DefaultCsvPoints            = 1
	SheetName                   = "demo"
	QuizTitle                   = "quiz_title"
	QuizTitleRequired           = "quiz-title is required"
	ErrGettingAttachment        = "error in getting file"
	ErrFileSizeExceed           = "The file is too large to upload. Please select a smaller file."
	ErrFileIsNotInSupportedType = "file has no supported type"
	ErrProblemInUploadFile      = "there was some error in file upload"
	ErrValidatingColumns        = "file columns not in proper format"
	ErrParsingFile              = "error in parsing file"
	ErrRowsReachesToMaxCount    = "rows limit exceed"
	ErrSurveyAnswerLength       = "in survey correct answer should contain all the options as correct"
	ErrSingleAnswerLength       = "in single answer there should be only one correct answer"
	ErrQuestionType             = "please provide a proper question type"
	ErrQuestionId               = "question type id not exists"
	ErrEmptyFile                = "The uploaded file is empty. Please choose a file with content."
	ErrUnsupportedFileType      = "The uploaded file is not a valid CSV. Please check the format and try again."
	ErrEmptyQuestionText        = "question text is required"
	ErrInsufficientOptions      = "at least 2 options are required"
	ErrEmptyCorrectAnswer       = "correct answer is required"
	ErrInvalidCorrectAnswer     = "correct answer must be a number referencing an existing option"
	ErrInvalidPoints            = "points must be a whole number within the allowed range"
	ErrInvalidCSVRows           = "the uploaded CSV has invalid rows, please fix them and try again"
	ErrInvalidQuestionTimeLimit = "question time limit is not configured properly"
	ErrInvalidQuestionMedia     = "question media must be one of: text, image, code"
	ErrInvalidOptionsMedia      = "options media must be one of: text, image, code"

	// ai quiz generation
	ErrAINotConfigured            = "ai quiz generation is not configured"
	ErrAIRequestFailed            = "could not reach the ai service, please try again"
	ErrAIRateLimited              = "the ai service is busy right now, please try again in a moment"
	ErrAIBudgetExhausted          = "the ai service usage limit has been reached, please try again later"
	ErrAIEmptyResponse            = "the ai service returned an empty response, please try again"
	ErrAIInvalidResponse          = "could not read the questions returned by the ai service, please try again"
	ErrAINoValidQuestions         = "the ai service did not return any usable questions, please try again"
	ErrAIInvalidQuestions         = "the generated questions are invalid, please regenerate"
	ErrAITooManyQuestions         = "number of questions must be between 1 and %d"
	ErrAITopicSingleLine          = "topic must be a single line"
	ErrAIDuplicateQuestion        = "duplicate question text"
	ErrAIDuplicateOption          = "two options are identical"
	ErrAITooManyOptions           = "at most %d options are allowed"
	ErrAIBlankOption              = "options must not be blank"
	ErrAICorrectAnswerOutOfRange  = "correct answer does not reference an existing option"
	ErrAIUnsupportedQuestionMedia = "question media must be text or code"
	ErrAIUnsupportedOptionsMedia  = "options media must be text or code"
	ErrAIMissingCodeResource      = "question media is code but no snippet was provided"
	ErrAIResourceTooLong          = "code snippet is too long"
	ErrAIOptionTooLong            = "an option is too long"
	ErrAIQuestionTooLong          = "question text is too long"
	ErrAIMissingQuizId            = "quiz id is required"
	ErrAIAppendFailed             = "could not add the generated questions to this quiz"
	ErrAIBaseUrlInvalid           = "the ai base url is not a valid http(s) url"
	ErrAIBaseUrlNotHTTPS          = "the ai base url must use https"
	ErrAIBaseUrlPrivate           = "this server does not allow ai base urls on local or private networks"
	ErrAIBaseUrlUnresolvable      = "that host name does not resolve"
	ErrAIModelInvalid             = "the ai model name is not valid"
	ErrAIApiKeyInvalid            = "the ai api key is not valid"
	ErrAIIncompleteCredentials    = "ai base url and model are both required"
	ErrAIModelsUnavailable        = "that provider did not return a model list, type the model name instead"
	ErrAIProviderUnreachable      = "could not connect to that base url. the jovvix server makes this call, not your browser"
	ErrAITestKeyRejected          = "the provider rejected that api key"
	ErrAITestKeyForbidden         = "that api key is not allowed to use this model or endpoint"
	ErrAITestNotFound             = "endpoint or model not found. check the base url (many providers need /v1) and the model name"
	ErrAITestBadRequest           = "the provider rejected the request. the model name is probably wrong"
	ErrAITestServerError          = "the provider returned a server error, please try again"
	ErrAITestTimedOut             = "timed out connecting to that base url"
	ErrAITestTLS                  = "tls handshake failed. check http against https"
	ErrAITestNoValidQuestion      = "connected, but this model did not return a usable quiz question. pick another model"
	ErrAIResponseTooLarge         = "the ai service returned too much data, please try again"

	// quiz-id
	QuizId       = "quiz_id"
	QuestionId   = "question_id"
	SharedQuizId = "shared_quiz_id"
	CategoryId   = "category_id"

	// Base64 data-URI cover images inflate ~4/3 over the raw file; 1 MiB of
	// text comfortably covers the 500 KB client-side file limit.
	MaxCoverImageBytes = 1 << 20
)

var AllowedCoverImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/heic": true,
	"image/heif": true,
}

// components
const (
	Waiting  = "Waiting"
	Question = "Question"
	Score    = "Score"
	Loading  = "Loading"
	Running  = "Running"

	ToUser  = 1
	ToAdmin = 2
	ToAll   = 3
)

// constants
const (
	MinInvitationCode = 100000
	MaxInvitationCode = 999999
	Counter           = 5
	Count             = 3
	StreakBaseScore   = 100
)

// Quiz Events
const (
	// Event 1. Authentication <admin side>
	EventAuthentication  = "authentication"
	ActionAuthentication = "authentication check"
	Unauthenticated      = "unauthenticated to access resource"
	InvalidCredentials   = "invalid credentials"

	// Event 2. Authorization <admin/user side>
	EventAuthorization  = "authorization"
	ActionAuthorization = "check for access"
	UserNotExist        = "user does not exists"
	Unauthorized        = "unauthorized to access resource"

	// Event 3. Session Validation <admin>
	EventSessionValidation  = "session_validation"
	ActionSessionValidation = "session validation from server side"
	ErrSessionNotFound      = "session unavailable"

	// Event 4. UserSession Validation <admin/user>
	EventUserSessionValidation   = "user_validation"
	ActionUserSessionValidation  = "user session get or create"
	CurrentUserQuiz              = "user_played_quiz"               // use by web
	ErrUserQuizSessionValidation = "quiz-session-validation-failed" // use by web
	ErrAdminCannotBeUser         = "host cannot be a player in their own quiz"

	EventRedirectToAdmin     = "redirect_to_admin"
	ActionCurrentUserIsAdmin = "current user is admin"

	// Event 4. Active session <admin>
	EventActivateSession          = "session_activation"
	EventSendInvitationCode       = "send_invitation_code" // use by web
	ActionSessionActivation       = "activate demanded session and sent invitation code"
	QuizSessionInvitationCode     = "invitationCode"
	SessionIDParam                = "session_id"
	ActiveQuizObj                 = "current active quiz obj"
	NoPlayerFound                 = "No player found"
	StartQuizByAdminNoPlayerFound = "start quiz by admin but no player found"
	ActionSendUserData            = "send user join data"
	JoinUserOnRunningQuiz         = "join_user_on_running_quiz"

	// Event 5. Join quiz <User>
	EventJoinQuiz                  = "invitation_code_validation"
	ActionJoinQuiz                 = "invitation code validation"
	ErrInvitationCodeInWrongFormat = "invitation code is not in proper format"
	ErrInvitationCodeNotFound      = "invitation code not found" // use by web
	ErrSessionWasCompleted         = "session was completed"     // use by web
	ErrMaxTryToGenerateCode        = "maximum times excide to generate code"

	UserName          = "username"
	UserUkey          = "users_username_ukey"
	Join              = "join access"
	QuizStartsSoon    = "Quiz will start soon"
	ErrUsernameExists = "username already exists"

	// Event 6. Start quiz <admin>
	EventUserJoined       = "user joined"
	EventStartQuiz        = "start_quiz"       // use by web
	EventStartQuizByAdmin = "startQuizByAdmin" // use by admin for start quiz

	// Event 7. Get Questions
	EventSendQuestion              = "send_question"
	ActionSendQuestion             = "send single question to user"
	QuizQuestionStatus             = "quiz question status"
	GetQuestions                   = "get quiz questions"
	NextQuestionWillServeSoon      = "Next question will be coming soon"
	ErrInGettingQuestion           = "error during getting question"
	ErrInGettingTotalQuestionCount = "error during getting total question count"

	EventPublishQuestion       = "publish_question"
	EventStartCount5           = "5_sec_counter" // use by web
	ActionCounter              = "5 second counter"
	EventNextQuestionAsked     = "next_question"         // use by web
	AdminDisconnected          = "admin_is_disconnected" // use by web
	EventAnswerSubmittedByUser = "answer submitted by user"
	ActionAnserSubmittedByUser = "answer submitted by user"

	// Event 8. Submit answer
	ErrQuizNotFound           = "error current quiz not found"
	ErrAnswerSubmit           = "error malfunction in inputs"
	ErrAnswerAlreadySubmitted = "answer already submitted"
	ErrQuestionNotActive      = "question can not receive answers anymore"
	ErrPublishAnswer          = "error while publishing answer in redis"

	// Event skip
	EventSkipAsked  = "ask_skip" // use by web
	WarnSkip        = "Some players haven't submitted their answers yet. Would you like to skip?"
	EventForceSkip  = "ask_force_skip"
	EventSkipTimer  = "skip_timer"
	EventPauseQuiz  = "pause_quiz"
	EventResumeQuiz = "resume_quiz"

	// Event 8. Get score page
	EventShowScore  = "show_score"
	ActionShowScore = "show score page during quiz"

	// Event 9. Terminate quiz
	EventTerminateQuiz  = "terminate_quiz"
	ActionTerminateQuiz = "terminate quiz after completing"

	ErrActiveDeleteQuiz = "Cannot terminate the quiz while it is still active."

	// Event 10. unhandled event
	UnknownError  = "unknown_error"
	ErrJWTExpired = "JWT token expired, Please try again later"

	// Event 11. ping
	EventPing = "ping"
	EventPong = "pong"
)

// final scoreboard cookie for user
const UserPlayedQuiz = "user_played_quiz"
const ActiveQuizId = "active_quiz_id"

// database table names
const (
	UserQuizResponsesTable   = "user_quiz_responses"
	UserPlayedQuizzesTable   = "user_played_quizzes"
	QuestionsTable           = "questions"
	UsersTable               = "users"
	ActiveQuizzesTable       = "active_quizzes"
	QuizQuestionsTable       = "quiz_questions"
	ActiveQuizQuestionsTable = "active_quiz_questions"
	QuizzesTable             = "quizzes"
)

// Question Types
const (
	SingleAnswerString = "single answer"
	SurveyString       = "survey"

	SingleAnswer = 1
	Survey       = 2
)

// Media Types
const (
	MediaText  = "text"
	MediaImage = "image"
	MediaCode  = "code"
)

// AI quiz generation
const (
	AIDefaultLanguage = "english"

	AIDifficultyEasy   = "easy"
	AIDifficultyMedium = "medium"
	AIDifficultyHard   = "hard"

	AIQuestionTypeSingle = "single"
	AIQuestionTypeSurvey = "survey"

	AICompletionsPath = "/chat/completions"
	AIModelsPath      = "/models"

	HeaderAIBaseUrl = "X-AI-Base-Url"
	HeaderAIApiKey  = "X-AI-Api-Key"
	HeaderAIModel   = "X-AI-Model"

	AIProviderBodyRedacted = "<omitted: caller supplied provider>"

	AIDefaultMaxQuestions       = 20
	AIMaxQuestionsHardLimit     = 50
	AIDefaultTimeoutSeconds     = 90
	AIDefaultTemperature        = 0.4
	AIDefaultOptionsPerQuestion = 4
	AIMinOptionsPerQuestion     = 2
	AIMaxOptionsPerQuestion     = 5
	AIDefaultQuestionDuration   = 60
	AIDefaultQuestionPoints     = 1
	AIProviderLogBodyLimit      = 800
	AIMaxResourceLength         = 1200
	AIMaxOptionLength           = 400
	AIMaxQuestionLength         = 500
	AIMaxExplanationLength      = 500
	AIMaxTitleLength            = 50
	AIMaxDescriptionLength      = 150
	AIResponseTokenBudget       = 400
	AIMinResponseTokens         = 1500
	AITestTimeoutSeconds        = 20
	AITestGenerateTimeout       = 45
	AITestTopic                 = "general knowledge"
	AIMaxRequestAdaptations     = 3
	AIMaxBaseUrlLength          = 300
	AIMaxModelLength            = 200
	AIMaxApiKeyLength           = 512
	AIDialTimeoutSeconds        = 10
	AIMaxResponseBytes          = 4 << 20
	AIMaxListedModels           = 300
	AIMaxAvoidQuestions         = 40
	AIMaxAvoidQuestionLength    = 120
	AIMaxProviderReasonLength   = 200

	AIQuestionsTruncated     = "the ai returned more questions than requested, extras were dropped"
	NoticeAIFewerQuestions   = "the ai returned %d of the %d questions you asked for"
	NoticeAIDroppedQuestions = "%d generated questions were dropped because they were malformed"

	AIAdaptationDroppedJSONMode     = "dropped response_format"
	AIAdaptationMaxCompletionTokens = "switched to max_completion_tokens"
	AIAdaptationDroppedTemperature  = "dropped temperature"
)

// Pagination and Filters

const (
	PageNumberQueryParam = "page"
	NameQueryParam       = "name"
	OrderQueryParam      = "order"
	OrderByQueryParam    = "orderBy"
	DefaultPageSize      = 10
)

// Channel name for redis pubsub
const (
	ChannelUserJoin       = "user_joined"
	ChannelUserDisconnect = "user_disconnect"
	ChannelSetAnswer      = "set_answer"
)
