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
	ParamUid = "userId"
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
)

// default Events
const (
	EventUserRegistered = "event:userRegistered"
)

// Middleware
const (
	// socket
	MiddlewarePass  = "allowed"
	MiddlewareError = "middleware_error"

	// http/https
	ErrNotAllowed = "Not allowed to access Resource"
)

// components
const (
	Waiting  = "Waiting"
	Question = "Question"
	Score    = "Score"
)

// constants
const (
	MinCode = 100000
	MaxCode = 999999
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
	EventSessionValidation  = "validate session"
	ActionSessionValidation = "session validation from server side"
	ErrSessionNotFound      = "session unavailable"

	// Event 4. Active session <admin>
	EventActivateSession    = "activate session"
	ActionSessionActivation = "activate demanded session and sent code"
	EventSendCode           = "send code to admin"
	QuizSessionCode         = "code"
	SessionIDParam          = "session_id"
	SessionObj              = "current session object"

	// Event 5. Join quiz <User>
	EventJoinQuiz          = "code validation"
	ActionJoinQuiz         = "quiz code validation"
	ErrCodeInWrongFormat   = "code is not in proper format"
	ErrCodeNotFound        = "code not found" // use by web
	ErrCodeNotActive       = "code is not active"
	ErrSessionWasCompleted = "session was completed"

	UserName          = "username"
	UserUkey          = "users_username_ukey"
	Join              = "join access"
	QuizStartsSoon    = "quiz will start soon"
	ErrUsernameExists = "username already exists"

	// Event 6. Start quiz <admin>
	EventStartQuiz = "start quiz" // use by web

	// Event 7. Get Questions
	QuizQuestionStatus        = "quiz question status"
	GetQuestions              = "get quiz questions"
	NextQuestionWillServeSoon = "Next question will coming soon"

	// Event . unhandled event
	UnknownError = "unknown_error"
)
