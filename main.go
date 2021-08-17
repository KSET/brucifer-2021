package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

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

	app := fiber.New(fiber.Config{
		GETOnly:          true,
		DisableKeepalive: true,
		ReadTimeout:      10 * time.Second,
		ServerHeader:     "Microsoft-IIS/7.0",
		AppName:          "Brucifer 2021.",
	})

	app.Use(recover.New())
	app.Use(etag.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path}\t| ${status} ${latency}\t| ${ua}\t| ${ips} \n",
	}))

	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")
		// c.Set("Strict-Transport-Security", "max-age=5184000")

		// Go to next middleware:
		return c.Next()
	})

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
		c.Type("html", "utf-8")
		return c.SendString(indexHtml)
	})

	app.Get("/kontakt", func(c *fiber.Ctx) error {
		c.Type("html", "utf-8")
		return c.SendString(kontaktHtml)
	})

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
