package constants

// variables
const (
	CookieUser   = "user"
	KratosCookie = "ory_kratos_session"
)

// fiber contexts
const (
	ContextUid = "userId"
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
	Unauthenticated    = "unauthenticated to access resource"
	Unauthorized       = "unauthorized to access resource"
	InvalidCredentials = "invalid credenticals"
	UserNotExist       = "user does not exists"
)

// Error messages
const (
	ErrGetUser             = "error while get user"
	ErrLoginUser           = "error while login user"
	ErrInsertUser          = "error while creating user, please try after sometime"
	ErrHealthCheckDb       = "error while checking health of database"
	ErrUnauthenticated     = "error verifing user identity"
	ErrKratosAuth          = "error while fetching user from kratos"
	ErrKratosDataInsertion = "error while inserting user data came from kratos"
	ErrKratosIDEmpty       = "error no session_id found in kratos cookie"
	ErrKratosCookieTime    = "error while parsing the expiration time of the cookie"
)

// Events
const (
	EventUserRegistered = "event:userRegistered"
)
