package routes

import (
	"brucosijada.kset.org/routes/base"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	base.RegisterRoutes(app)
}
