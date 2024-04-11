package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/jwtutils"
	"github.com/yigithankarabulut/distributed-mail-queue-service/pkg/response"
	"os"
	"strings"
)

// AuthMiddleware is the middleware for checking the token.
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := response.New()
		tokenStr := c.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			return c.JSON(res.BasicError("token is required", fiber.StatusUnauthorized))
		}
		token, err := jwt.ParseWithClaims(tokenStr, &jwtutils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.JSON(res.BasicError(fmt.Sprintf("error parsing token: %s", err.Error()), fiber.StatusUnauthorized))
		}
		if !token.Valid {
			return c.JSON(res.BasicError("invalid token", fiber.StatusUnauthorized))
		}
		claims, ok := token.Claims.(jwtutils.CustomClaims)
		if !ok {
			return c.JSON(res.BasicError("error parsing token claims", fiber.StatusUnauthorized))
		}
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
