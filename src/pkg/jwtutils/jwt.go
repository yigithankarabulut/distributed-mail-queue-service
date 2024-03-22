package jwtutils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yigithankarabulut/distributed-mail-queue-service/src/pkg/response"
	"os"
	"strings"
	"time"
)

type IJwtUtils interface {
	AuthMiddleware() func(c *fiber.Ctx) error
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
	claims := CustomClaims{
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

// AuthMiddleware is the middleware for checking the token.
func (j *JwtUtils) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := response.New()
		tokenStr := c.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			return c.JSON(res.BasicError("token is required", fiber.StatusUnauthorized))
		}
		token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.JSON(res.BasicError(fmt.Sprintf("error parsing token: %s", err.Error()), fiber.StatusUnauthorized))
		}
		if !token.Valid {
			return c.JSON(res.BasicError("invalid token", fiber.StatusUnauthorized))
		}
		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			return c.JSON(res.BasicError("error parsing token claims", fiber.StatusUnauthorized))
		}
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
