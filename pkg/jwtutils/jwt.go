package jwtutils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type IJwtUtils interface {
	GenerateJwtToken(userID uint, expiration time.Duration) (string, error)
}

type JwtUtils struct{}

func New() *JwtUtils {
	return &JwtUtils{}
}

// CustomClaims is the custom claims for jwt.
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateJwtToken generates a jwt token.
func (j *JwtUtils) GenerateJwtToken(userID uint, expiration time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := &CustomClaims{
		userID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
