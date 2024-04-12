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
func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		res := response.New()
		tokenStr := c.Get("Authorization")
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(res.BasicError("missing token", fiber.StatusUnauthorized))
		}
		token, err := jwt.ParseWithClaims(tokenStr, &jwtutils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(res.BasicError(fmt.Sprintf("error parsing token: %s", err.Error()), fiber.StatusUnauthorized))
		}
		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(res.BasicError("invalid token", fiber.StatusUnauthorized))
		}
		claims, ok := token.Claims.(*jwtutils.CustomClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(res.BasicError("invalid token", fiber.StatusUnauthorized))
		}
		id := uint(claims.UserID)
		if id == 0 {
			return c.Status(fiber.StatusUnauthorized).JSON(res.BasicError("invalid token", fiber.StatusUnauthorized))
		}
		c.Locals("userID", id)
		return c.Next()
	}
}
