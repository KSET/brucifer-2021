package user

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func ListUsers(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(
		response.Success(
			repo.User().List(),
		),
	)
}
