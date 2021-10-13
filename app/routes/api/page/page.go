package page

import (
	"mime/multipart"
	"sort"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/providers/markdown"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func RenderPage(ctx *fiber.Ctx) error {
	page := models.Page{}

	err := database.DatabaseProvider().Client().Where(
		"name = ?",
		ctx.Params("pageName"),
	).Preload("Background.Variations").Preload("BackgroundMobile.Variations").First(&page).Error

	if err != nil {
		return ctx.Next()
	}

	var background *repo.ImageItemVariation = nil
	if page.Background != nil {
		bg := repo.Image().ToImageItem(page.Background)
		sort.Slice(
			bg.Variations, func(i, j int) bool {
				return bg.Variations[i].Width > bg.Variations[j].Width
			},
		)
		background = &bg.Variations[0]
	}

	var backgroundMobile *repo.ImageItemVariation = nil
	if page.BackgroundMobile != nil {
		bg := repo.Image().ToImageItem(page.BackgroundMobile)
		sort.Slice(
			bg.Variations, func(i, j int) bool {
				return bg.Variations[i].Width > bg.Variations[j].Width
			},
		)
		backgroundMobile = &bg.Variations[0]
	}

	return ctx.Render(
		"shell",
		fiber.Map{
			"content":          page.Rendered,
			"background":       background,
			"backgroundMobile": backgroundMobile,
		},
		"layouts/main",
	)
}

func CreatePage(ctx *fiber.Ctx) (err error) {
	var input struct {
		Name     string
		Markdown string
	}
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	var background *multipart.FileHeader
	if background, err = ctx.FormFile("background"); err != nil && err != fasthttp.ErrMissingFile {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request for background",
				err,
			),
		)
	}

	var backgroundMobile *multipart.FileHeader
	if backgroundMobile, err = ctx.FormFile("backgroundMobile"); err != nil && err != fasthttp.ErrMissingFile {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request for mobile background",
				err,
			),
		)
	}

	page, err := repo.Page().Create(
		repo.PageCreateOptions{
			Name:             input.Name,
			Markdown:         input.Markdown,
			Uploader:         ctx.Locals("user").(models.User),
			Background:       background,
			BackgroundMobile: backgroundMobile,
		},
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
			page,
		),
	)
}

func UpdatePage(ctx *fiber.Ctx) (err error) {
	var input struct {
		Name               string
		Markdown           string
		BackgroundId       string
		BackgroundMobileId string
	}
	if err = ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request",
				err,
			),
		)
	}

	var background *multipart.FileHeader
	if background, err = ctx.FormFile("background"); err != nil && err != fasthttp.ErrMissingFile {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request for background",
				err,
			),
		)
	}

	var backgroundMobile *multipart.FileHeader
	if backgroundMobile, err = ctx.FormFile("backgroundMobile"); err != nil && err != fasthttp.ErrMissingFile {
		return ctx.Status(fiber.StatusBadRequest).JSON(
			response.Error(
				"Invalid request for mobile background",
				err,
			),
		)
	}

	db := database.DatabaseProvider().Client()

	var model = models.Page{}
	err = db.Model(model).Where(
		"uuid = ?",
		ctx.Params("id"),
	).Preload(
		"Background",
	).Preload(
		"BackgroundMobile",
	).Preload(
		"Background.Variations",
	).Preload(
		"BackgroundMobile.Variations",
	).Find(&model).Error

	if err != nil {
		return ctx.JSON(
			response.Error(
				"Something went wrong",
				err,
			),
		)
	}

	if !model.Exists() {
		return ctx.Status(fiber.StatusNotFound).JSON(
			response.Error(
				"Page not found",
				err,
			),
		)
	}

	var bg *models.Image = nil
	if model.Background != nil && model.Background.UUID.String() == input.BackgroundId {
		bg = model.Background
	}

	var bgMobile *models.Image = nil
	if model.BackgroundMobile != nil && model.BackgroundMobile.UUID.String() == input.BackgroundMobileId {
		bgMobile = model.BackgroundMobile
	}

	data := repo.PageUpdateOptions{
		Model:             &model,
		Name:              input.Name,
		Markdown:          input.Markdown,
		Uploader:          ctx.Locals("user").(models.User),
		Background:        background,
		BackgroundM:       bg,
		BackgroundMobile:  backgroundMobile,
		BackgroundMobileM: bgMobile,
	}

	err = repo.Page().Update(data)

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
			repo.Page().ToPageItem(&model),
		),
	)
}

func DeletePage(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	err = db.Where(
		"uuid = ?",
		ctx.Params("id"),
	).Delete(&models.Page{}).Error

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
			nil,
		),
	)
}

func ListPages(ctx *fiber.Ctx) error {
	pages, err := repo.Page().ListSimple()

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
			pages,
		),
	)
}

func ShowPage(ctx *fiber.Ctx) (err error) {
	var id uuid.UUID

	if id, err = uuid.Parse(ctx.Params("id")); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(
			response.Error(
				"Invalid page ID",
				err,
			),
		)
	}

	page, err := repo.Page().GetByUuid(id)
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
			repo.Page().ToPageItem(page),
		),
	)
}

func RenderFromParam(c *fiber.Ctx) error {
	md := c.Query("p")

	content := markdown.MarkdownProvider().Render(md)

	return c.Render(
		"shell",
		fiber.Map{
			"content": content,
		},
		"layouts/main",
	)
}
