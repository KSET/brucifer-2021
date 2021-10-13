package routes

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/middleware"
	"brucosijada.kset.org/app/response"
	"brucosijada.kset.org/app/routes/api/artist"
	"brucosijada.kset.org/app/routes/api/auth"
	"brucosijada.kset.org/app/routes/api/image"
	"brucosijada.kset.org/app/routes/api/page"
	"brucosijada.kset.org/app/routes/api/sponsor"
	"brucosijada.kset.org/app/routes/api/user"
	"brucosijada.kset.org/app/routes/base"
)

func RegisterRoutes(app *fiber.App) {
	Base := app.Group("/")
	Base.Get("/", base.Home)
	Base.Get("/lineup", base.Lineup)
	Base.Get("/sponzori", base.Sponsors)
	Base.Get("/i/:id", image.ShowImage)
	Base.Get("/:pageName", page.RenderPage)

	Api := app.Group(
		"/api", func(c *fiber.Ctx) error {
			c.Type("json")

			return c.Next()
		},
	)

	ApiAuth := Api.Group("/auth")
	ApiAuth.Post("/login", auth.Login)
	ApiAuth.Post("/logout", auth.Logout)
	ApiAuth.Post("/register", auth.Register)
	ApiAuth.Get("/user", middleware.RequireAuth(), auth.GetUser)

	ApiSponsor := Api.Group("/sponsor")
	ApiSponsor.Post("/", middleware.RequireAuth(), sponsor.CreateSponsor)
	ApiSponsor.Patch("/swap/:id", middleware.RequireAuth(), sponsor.SwapSponsors)
	ApiSponsor.Delete("/:id", middleware.RequireAuth(), sponsor.DeleteSponsor)
	ApiSponsor.Get("/", sponsor.ListSponsors)
	ApiSponsor.Get("/:id", middleware.RequireAuth(), sponsor.ShowSponsor)
	ApiSponsor.Patch("/:id", middleware.RequireAuth(), sponsor.UpdateSponsor)

	ApiArtist := Api.Group("/artist")
	ApiArtist.Post("/", middleware.RequireAuth(), artist.CreateArtist)
	ApiArtist.Patch("/swap/:id", middleware.RequireAuth(), artist.SwapArtists)
	ApiArtist.Delete("/:id", middleware.RequireAuth(), artist.DeleteArtist)
	ApiArtist.Get("/", artist.ListArtists)
	ApiArtist.Get("/:id", middleware.RequireAuth(), artist.ShowArtist)
	ApiArtist.Patch("/:id", middleware.RequireAuth(), artist.UpdateArtist)

	ApiUser := Api.Group("/user")
	ApiUser.Get("/", middleware.RequireAuth(), user.ListUsers)
	ApiUser.Patch("/", middleware.RequireAuth(), user.ChangePassword)

	ApiUserInvitation := Api.Group("/user-invitation")
	ApiUserInvitation.Get("/info/:code", user.InvitationInfo)
	ApiUserInvitation.Get("/", middleware.RequireAuth(), user.ListInvitations)
	ApiUserInvitation.Post("/", middleware.RequireAuth(), user.CreateInvitation)
	ApiUserInvitation.Delete("/:id", middleware.RequireAuth(), user.DeleteInvitation)

	ApiPage := Api.Group("/page", middleware.RequireAuth())
	ApiPage.Post("/", page.CreatePage)
	ApiPage.Patch("/:id", page.UpdatePage)
	ApiPage.Delete("/:id", page.DeletePage)
	ApiPage.Get("/", page.ListPages)
	ApiPage.Get("/rendered", page.RenderFromParam)
	ApiPage.Get("/:id", page.ShowPage)

	Api.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(
				response.Error(
					"Not found",
					nil,
				),
			)
		},
	)

	Base.Use(
		func(c *fiber.Ctx) (err error) {
			return c.Status(fiber.StatusNotFound).Render(
				"shell",
				fiber.Map{
					"content": "<h1>Greška 404</h1><h2>Stranica nije pronađena<h2>",
				},
				"layouts/main",
			)
		},
	)
}
