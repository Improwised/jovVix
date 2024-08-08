package constants

// variables
const (
	CookieUser   = "user"
	KratosCookie = "ory_kratos_session"
)

// fiber contexts
const (
	ContextUid  = "userId"
	ContextUser = "userContext"
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
	ErrGetUser             = "error while get user"
	ErrLoginUser           = "error while login user"
	ErrInsertUser          = "error while creating user, please try after sometime"
	ErrHealthCheckDb       = "error while checking health of database"
	ErrUnauthenticated     = "error verifying user identity"
	ErrKratosAuth          = "error while fetching user from kratos"
	ErrKratosDataInsertion = "error while inserting user data came from kratos"
	ErrKratosIDEmpty       = "error no session_id found in kratos cookie"
	ErrKratosCookieTime    = "error while parsing the expiration time of the cookie"
	ErrRegisterQuiz        = "error while creating quiz"
	ErrCreatingDemoQuiz    = "error while creating demo quiz"
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
	SheetName                   = "demo"
	QuizTitle                   = "quiz_title"
	QuizTitleRequired           = "quiz-title is required"
	ErrGettingAttachment        = "error in getting file"
	ErrFileSizeExceed           = "file size exceed"
	ErrFileIsNotInSupportedType = "file has no supported type"
	ErrProblemInUploadFile      = "there was some error in file upload"
	ErrValidatingColumns        = "file columns not in proper format"
	ErrParsingFile              = "error in parsing file"
	ErrRowsReachesToMaxCount    = "rows limit exceed"
	ErrSurveyAnswerLength       = "in survey correct answer should contain all the options as correct"
	ErrSingleAnswerLength       = "in single answer there should be only one correct answer"
	ErrQuestionType             = "please provide a proper question type"
	ErrQuestionId               = "question type id not exist"

	// quiz-id
	QuizId = "quiz_id"
)

// components
const (
	Waiting  = "Waiting"
	Question = "Question"
	Score    = "Score"
	Loading  = "Loading"

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
	PageSize          = 10
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
	NoPlayerFound                 = "no player found"
	StartQuizByAdminNoPlayerFound = "start quiz by admin but no player found"

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
	QuizStartsSoon    = "quiz will start soon"
	ErrUsernameExists = "username already exists"

	// Event 6. Start quiz <admin>
	EventUserJoined       = "user joined"
	EventStartQuiz        = "start_quiz"       // use by web
	EventStartQuizByAdmin = "startQuizByAdmin" // use by admin for start quiz

	// Event 7. Get Questions
	EventSendQuestion         = "send_question"
	ActionSendQuestion        = "send single question to user"
	QuizQuestionStatus        = "quiz question status"
	GetQuestions              = "get quiz questions"
	NextQuestionWillServeSoon = "Next question will coming soon"
	ErrInGettingQuestion      = "error during getting question"

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

	// Event skip
	EventSkipAsked = "ask_skip" // use by web
	WarnSkip       = "some player didn't submit their answer yet. would you want to skip?"
	EventForceSkip = "ask_force_skip"
	EventSkipTimer = "skip_timer"

	// Event 8. Get score page
	EventShowScore  = "show_score"
	ActionShowScore = "show score page during quiz"

	// Event 9. Terminate quiz
	EventTerminateQuiz  = "terminate_quiz"
	ActionTerminateQuiz = "terminate quiz after completing"

	// Event . unhandled event
	UnknownError  = "unknown_error"
	ErrJWTExpired = "JWT token expired, Please try again"
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
