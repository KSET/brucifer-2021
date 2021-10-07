package auth

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/response"
	"brucosijada.kset.org/app/util/auth"
	"brucosijada.kset.org/app/util/cookie"
)

func Login(ctx *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid login request",
				err,
			),
		)
	}

	if input.Password == "" || input.Identity == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Identity and password required",
				nil,
			),
		)
	}

	user := auth.AuthProvider().ValidateUser(input.Identity, input.Password)

	if user != nil {
		cookie.Cookie().SetAuthCookie(user, ctx)

		return ctx.JSON(
			response.Success(
				*user,
			),
		)
	}

	return ctx.JSON(
		response.Error(
			"Incorrect identity or password",
			nil,
		),
	)
}

func Logout(ctx *fiber.Ctx) error {
	cookie.Cookie().RemoveAuthCookie(ctx)

	return ctx.JSON(
		response.Success(
			nil,
		),
	)
}

func GetUser(c *fiber.Ctx) error {
	user := c.Locals("user")

	return c.JSON(
		response.Success(
			user,
		),
	)
}
