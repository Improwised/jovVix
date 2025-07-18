package jwt

import (
	"time"

	"github.com/Improwised/jovvix/api/config"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// ParseToken parse, validate the jwt token
// On valid token it returns the decoded token
func ParseToken(config config.AppConfig, token string) (jwt.Token, error) {
	key, err := jwk.FromRaw([]byte(config.Secret))
	if err != nil {
		return nil, err
	}

	claims, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, key), jwt.WithIssuer(config.JWTIssuer))
	return claims, err
}

func CreateToken(config config.AppConfig, sub string, exp time.Time) (string, error) {
	stringToken := ""
	token, err := jwt.NewBuilder().Subject(sub).Expiration(exp).Issuer(config.JWTIssuer).Build()
	if err != nil {
		return stringToken, err
	}
	key, err := jwk.FromRaw([]byte(config.Secret))
	if err != nil {
		return stringToken, err
	}
	signed, err := jwt.Sign(token, jwt.WithKey(jwa.HS256, key))
	if err != nil {
		return stringToken, err
	}
	stringToken = string(signed)
	return stringToken, nil
}
