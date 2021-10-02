package base

import (
	"github.com/gofiber/fiber/v2"
)

func Home(c *fiber.Ctx) error {
	return c.Render(
		"index",
		fiber.Map{},
		"layouts/main",
	)
}

func Contact(c *fiber.Ctx) error {
	return c.Render(
		"kontakt",
		fiber.Map{},
		"layouts/main",
	)
}

func Tickets(c *fiber.Ctx) error {
	return c.Render(
		"ulaznice",
		fiber.Map{},
		"layouts/main",
	)
}

func Rules(c *fiber.Ctx) error {
	return c.Render(
		"pravila",
		fiber.Map{},
		"layouts/main",
	)
}

func Lineup(c *fiber.Ctx) error {
	return c.Render(
		"lineup",
		fiber.Map{},
		"layouts/main",
	)
}
