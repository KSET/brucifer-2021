package sponsor

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func CreateSponsor(ctx *fiber.Ctx) (err error) {
	type form struct {
		Name string
	}

	var input form
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	var logo *multipart.FileHeader
	if logo, err = ctx.FormFile("logo"); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request: Missing file",
				err,
			),
		)
	}

	sponsor, err := repo.Sponsor().Create(
		input.Name,
		logo,
		ctx.Locals("user").(models.User).ID,
	)

	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				err.Error(),
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			sponsor,
		),
	)
}

func ListSponsors(ctx *fiber.Ctx) error {
	items, err := repo.Sponsor().ListSimple()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			items,
		),
	)
}

func UpdateSponsor(ctx *fiber.Ctx) (err error) {
	type form struct {
		Name string
	}

	db := database.DatabaseProvider().Client()

	var sponsor = models.Sponsor{}
	db.Model(sponsor).Preload("Image").Preload("Image.Variations").Where(
		"uuid = ?",
		ctx.Params("id"),
	).First(&sponsor)

	if !sponsor.Exists() {
		return ctx.Status(fiber.StatusNotFound).JSON(
			response.Error(
				"Sponsor not found",
				nil,
			),
		)
	}

	var input form
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	var logo *multipart.FileHeader
	if logo, err = ctx.FormFile("image"); err != nil && err != fasthttp.ErrMissingFile {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request: Missing file",
				err,
			),
		)
	}

	err = repo.Sponsor().Update(
		&sponsor,
		input.Name,
		logo,
		ctx.Locals("user").(models.User).ID,
	)

	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				err.Error(),
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			repo.Sponsor().ToSponsorItem(&sponsor),
		),
	)
}

func ShowSponsor(ctx *fiber.Ctx) (err error) {
	var id uuid.UUID

	if id, err = uuid.Parse(ctx.Params("id")); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid sponsor ID",
				err,
			),
		)
	}

	sponsor, err := repo.Sponsor().GetByUuidSimple(id)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				err.Error(),
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			sponsor,
		),
	)
}

func SwapSponsors(ctx *fiber.Ctx) (err error) {
	input := struct {
		With string `json:"with"`
	}{}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	var swap, with uuid.UUID

	if swap, err = uuid.Parse(ctx.Params("id")); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid swap sponsor ID",
				err,
			),
		)
	}
	if with, err = uuid.Parse(input.With); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid with sponsor ID",
				err,
			),
		)
	}

	err = repo.Sponsor().Swap(swap, with)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	newSponsors, err := repo.Sponsor().ListSimple()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			newSponsors,
		),
	)
}

func DeleteSponsor(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	var sponsor models.Sponsor
	db.Model(&models.Sponsor{}).Preload("Image").Where("uuid = ?", ctx.Params("id")).First(&sponsor)

	if !sponsor.Exists() {
		return ctx.JSON(
			response.Success(
				nil,
			),
		)
	}

	if err = db.Delete(&sponsor).Error; err != nil {
		return ctx.JSON(
			response.Error(
				"Something went horribly wrong",
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			sponsor,
		),
	)
}
