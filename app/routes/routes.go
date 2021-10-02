package routes

import (
	"brucosijada.kset.org/app/routes/base"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	Base := app.Group("/")
	Base.Get("/", base.Home)
	Base.Get("/kontakt", base.Contact)
	Base.Get("/ulaznice", base.Tickets)
	Base.Get("/pravila", base.Rules)
	Base.Get("/lineup", base.Lineup)
}
