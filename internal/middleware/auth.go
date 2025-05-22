package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gitlab.com/timkado/api/daisi-rest-postgres/pkg/utils"
)

func AuthenticateBearerToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate Token
		decryptedToken, err := utils.Decrypt(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("companyId", decryptedToken)

		return c.Next()
	}
}
