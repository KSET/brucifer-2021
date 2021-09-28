package routes

import (
	"github.com/gofiber/fiber/v2"
)

func registerBaseRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render(
			"index",
			fiber.Map{},
			"layouts/main",
		)
	})

	app.Get("/kontakt", func(c *fiber.Ctx) error {
		return c.Render(
			"kontakt",
			fiber.Map{},
			"layouts/main",
		)
	})

	app.Get("/ulaznice", func(c *fiber.Ctx) error {
		return c.Render(
			"ulaznice",
			fiber.Map{},
			"layouts/main",
		)
	})

	app.Get("/pravila", func(c *fiber.Ctx) error {
		return c.Render(
			"pravila",
			fiber.Map{},
			"layouts/main",
		)
	})

	app.Get("/lineup", func(c *fiber.Ctx) error {
		return c.Render(
			"lineup",
			fiber.Map{},
			"layouts/main",
		)
	})
}
