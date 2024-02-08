package config

// DBConfig type of db config object
type KratosConfig struct {
	IsEnabled            bool   `envconfig:"KRATOS_ENABLED"`
	BaseUrl              string `envconfig:"SERVE_PUBLIC_BASE_URL"`
	UIUrl                string `envconfig:"SELF_SERVICE_DEFAULT_BROWSER_RETURN_URL"`
	AdminUrl             string `envconfig:"SERVE_ADMIN_BASE_URL"`
	CookieExpirationTime string `envconfig:"KRATOS_COOKIE_EXPIRATION_TIME"`
}

type KratosUserDetails struct {
	Identity struct {
		ID     string `json:"id"`
		Traits struct {
			Name struct {
				Last  string `json:"last"`
				First string `json:"first"`
			} `json:"name"`
			Email string `json:"email"`
		} `json:"traits"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"identity"`
}
