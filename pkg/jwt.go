package pkg

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Payload struct {
	Id   string `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

type JWTErr struct {
	Type string
	Err  error
}

func NewPayload(id, role string) *Payload {
	return &Payload{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    os.Getenv("JWT_ISSUER"),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
}

func (c *Payload) GenerateToken() (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("Secret not provided")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString([]byte(jwtSecret))
}

func (c *Payload) VerifyToken(token string) JWTErr {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return JWTErr{
			Type: "System",
			Err:  errors.New("Secret not provided"),
		}
	}
	parsedToken, err := jwt.ParseWithClaims(token, c, func(t *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return JWTErr{
			Type: "Token",
			Err:  err,
		}
	}
	if !parsedToken.Valid {
		return JWTErr{
			Type: "Token",
			Err:  errors.New("Expired Token"),
		}
	}

	return JWTErr{}
}
