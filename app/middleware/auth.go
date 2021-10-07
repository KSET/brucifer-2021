package middleware

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
	"brucosijada.kset.org/app/util/auth"
)

func RequireAuth() fiber.Handler {
	auth := auth.AuthProvider()

	return func(c *fiber.Ctx) error {
		if userId, valid := auth.GetAuthUserId(c); valid {
			if user := repo.User().GetById(userId); user.Exists() {
				c.Locals("user", user)

				return c.Next()
			}
		}

		return c.Status(fiber.StatusUnauthorized).JSON(
			response.Error(
				"Authentication required",
				nil,
			),
		)
	}
}
