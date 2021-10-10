package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
	"brucosijada.kset.org/app/util/auth"
	"brucosijada.kset.org/app/util/cookie"
)

func Login(ctx *fiber.Ctx) error {
	input := struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}{}
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

func Register(ctx *fiber.Ctx) error {
	input := struct {
		Email    string `json:"email"`
		Identity string `json:"identity"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid register request",
				err,
			),
		)
	}

	if input.Code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invite code required",
				nil,
			),
		)
	}

	if input.Password == "" || input.Identity == "" || input.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Identity, email, and password required",
				nil,
			),
		)
	}

	invitation, err := repo.UserInvitation().GetByCode(input.Code)
	if err == gorm.ErrRecordNotFound {
		return ctx.Status(fiber.StatusNotFound).JSON(
			response.Error(
				"Invitation does not exist",
				err,
			),
		)
	}
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}
	if invitation.UsedByID != 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			response.Error(
				"Invitation already used",
				err,
			),
		)
	}

	user, err := repo.User().Create(
		repo.UserCreateOptions{
			Email:      input.Email,
			Identity:   input.Identity,
			Password:   input.Password,
			Invitation: invitation,
		},
	)

	return ctx.JSON(
		response.Success(
			fiber.Map{
				"user":       user,
				"invitation": invitation,
			},
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
