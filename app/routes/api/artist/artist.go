package artist

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

func CreateArtist(ctx *fiber.Ctx) (err error) {
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
	if logo, err = ctx.FormFile("image"); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request: Missing file",
				err,
			),
		)
	}

	artist, err := repo.Artist().Create(
		repo.ArtistCreateOptions{
			Name:     input.Name,
			Logo:     logo,
			Uploader: ctx.Locals("user").(models.User),
		},
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
			repo.Artist().ToArtistItem(artist),
		),
	)
}

func ListArtists(ctx *fiber.Ctx) error {
	items, err := repo.Artist().ListSimple()
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

func SwapArtists(ctx *fiber.Ctx) (err error) {
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
				"Invalid swap artist ID",
				err,
			),
		)
	}
	if with, err = uuid.Parse(input.With); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid with artist ID",
				err,
			),
		)
	}

	err = repo.Artist().Swap(swap, with)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	newArtists, err := repo.Artist().ListSimple()
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
			newArtists,
		),
	)
}

func DeleteArtist(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	var artist models.Artist
	db.Model(&models.Artist{}).Preload("Image").Where("uuid = ?", ctx.Params("id")).First(&artist)

	if !artist.Exists() {
		return ctx.JSON(
			response.Success(
				nil,
			),
		)
	}

	if err = db.Delete(&artist).Error; err != nil {
		return ctx.JSON(
			response.Error(
				"Something went horribly wrong",
				err,
			),
		)
	}

	return ctx.JSON(
		response.Success(
			repo.Artist().ToArtistItem(&artist),
		),
	)
}

func UpdateArtist(ctx *fiber.Ctx) (err error) {
	type form struct {
		Name string
	}

	db := database.DatabaseProvider().Client()

	var artist models.Artist
	db.Model(&models.Artist{}).Preload("Image").Preload("Image.Variations").Where(
		"uuid = ?",
		ctx.Params("id"),
	).First(&artist)

	if !artist.Exists() {
		return ctx.Status(fiber.StatusNotFound).JSON(
			response.Error(
				"Artist not found",
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

	err = repo.Artist().Update(
		repo.ArtistUpdateOptions{
			Model:    &artist,
			Name:     input.Name,
			Logo:     logo,
			Uploader: ctx.Locals("user").(models.User),
		},
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
			repo.Artist().ToArtistItem(&artist),
		),
	)
}

func ShowArtist(ctx *fiber.Ctx) (err error) {
	var id uuid.UUID

	if id, err = uuid.Parse(ctx.Params("id")); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid artist ID",
				err,
			),
		)
	}

	artist, err := repo.Artist().GetByUuidSimple(id)
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
			artist,
		),
	)
}
