package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"

	flag "github.com/spf13/pflag"
)

//go:embed assets/*
var assets embed.FS

//go:embed index.html
var indexHtml string

//go:embed kontakt.html
var kontaktHtml string

func main() {
	var port int
	var host string
	flag.IntVarP(&port, "port", "p", 3000, "Set the port on which the server will run")
	flag.StringVarP(&host, "host", "h", "0.0.0.0", "Set the host to which the server will bind")

	flag.Parse()

	app := fiber.New()

	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		c.Type("ico")

		favicon, _ := assets.ReadFile("assets/images/favicon.ico")

		return c.Send(favicon)
	})

	assetsFs, _ := fs.Sub(assets, "assets")
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root: http.FS(assetsFs),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(indexHtml)
	})

	app.Get("/kontakt", func(c *fiber.Ctx) error {
		c.Type("html")
		return c.SendString(kontaktHtml)
	})

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
