package user

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/providers/hash"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
	"brucosijada.kset.org/app/util/async"
)

func ListUsers(ctx *fiber.Ctx) (err error) {
	return ctx.JSON(
		response.Success(
			repo.User().List(),
		),
	)
}

func ChangePassword(ctx *fiber.Ctx) (err error) {
	user := ctx.Locals("user").(models.User)
	input := struct {
		Password struct {
			Old string `json:"old"`
			New string `json:"new"`
		} `json:"password"`
	}{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	data, _ := async.Async().RunInParallel(
		func() (ret interface{}, err error) {
			ret = hash.HashProvider().HashPassword(input.Password.New)
			return
		},
		func() (ret interface{}, err error) {
			ret = hash.HashProvider().VerifyPassword(input.Password.Old, user.Password)
			return
		},
	).All()

	newHash := data[0].(string)
	passwordsMatch := data[1].(bool)

	if !passwordsMatch {
		return ctx.JSON(
			response.Error(
				"Passwords don't match",
				nil,
			),
		)
	}

	user.Password = newHash

	err = database.DatabaseProvider().Client().Save(&user).Error

	if err != nil {
		return ctx.JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			user,
		),
	)
}
