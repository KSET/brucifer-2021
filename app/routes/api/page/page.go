package page

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/providers/markdown"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/response"
)

func RenderPage(name string) func(ctx *fiber.Ctx) error {
	page := models.Page{}

	err := database.DatabaseProvider().Client().Where(
		"name = ?",
		name,
	).First(&page).Error

	if err != nil {
		panic(err)
	}

	return func(ctx *fiber.Ctx) error {
		return ctx.Render(
			"shell",
			fiber.Map{
				"content": page.Rendered,
			},
			"layouts/main",
		)
	}
}

func CreatePage(ctx *fiber.Ctx) error {
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

	page, err := repo.Page().Create(
		repo.PageCreateOptions{
			Name:     input.Name,
			Markdown: input.Markdown,
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

	db := database.DatabaseProvider().Client()

	var model = models.Page{}
	err = db.Model(model).Where(
		"uuid = ?",
		ctx.Params("id"),
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

	err = repo.Page().Update(
		repo.PageUpdateOptions{
			Model:    &model,
			Name:     input.Name,
			Markdown: input.Markdown,
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
			model,
		),
	)
}

func DeletePage(ctx *fiber.Ctx) (err error) {
	db := database.DatabaseProvider().Client()

	var model = models.Page{}
	err = db.Where(
		"uuid = ?",
		ctx.Params("id"),
	).Delete(model).Error

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
	pages, err := repo.Page().List()

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
			page,
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
