package base

import (
	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/repo"
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

func Sponsors(c *fiber.Ctx) error {
	sponsorList, err := repo.Sponsor().ListSimple()

	if err != nil {
		return err
	}

	type sponsor struct {
		Name string
		Src  string
	}

	sponsors := make([]sponsor, len(*sponsorList))
	for i, item := range *sponsorList {
		src := repo.SponsorItemLogo{}
		for _, logo := range item.Logo {
			if src.Width < logo.Width {
				logo := logo
				src = logo
			}
		}

		sponsors[i] = sponsor{
			Name: item.Name,
			Src:  src.Url,
		}
	}

	return c.Render(
		"sponzori",
		fiber.Map{
			"sponsors": sponsors,
		},
		"layouts/main",
	)
}
