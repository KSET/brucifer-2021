package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func CreateInvitation(ctx *fiber.Ctx) (err error) {
	user := ctx.Locals("user").(models.User)
	input := struct {
		Comment string `json:"comment"`
	}{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	invitation, err := repo.UserInvitation().Create(
		user.ID,
		input.Comment,
	)

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
			invitation,
		),
	)
}

func DeleteInvitation(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	invitation := &models.UserInvitation{}
	err = db.Where("uuid = ?", ctx.Params("id")).First(&invitation).Error
	if err != nil {
		return ctx.JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	if invitation.UsedByID != 0 {
		return ctx.JSON(
			response.Error(
				"Code already used",
				nil,
			),
		)
	}

	err = db.Delete(&invitation).Error
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
			nil,
		),
	)
}

func ListInvitations(ctx *fiber.Ctx) (err error) {
	invitations, err := repo.UserInvitation().List()

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
			*invitations,
		),
	)
}

func InvitationInfo(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	invitation := struct {
		Code      string      `json:"code"`
		CreatorId uint        `json:"-"`
		Creator   models.User `json:"-"`
		By        string      `json:"by" gorm:"-"`
	}{}
	tx := db.Model(
		&models.UserInvitation{},
	).Preload(
		"Creator",
	).Where(
		"code = ? and used_by_id = 0",
		ctx.Params("code"),
	).First(&invitation)

	err = tx.Error

	if err == gorm.ErrRecordNotFound {
		return ctx.Status(fiber.StatusNotFound).JSON(
			response.Error(
				"Invitation not found",
				nil,
			),
		)
	}

	if err != nil {
		return ctx.JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	invitation.By = invitation.Creator.Username

	return ctx.JSON(
		response.Success(
			invitation,
		),
	)
}
