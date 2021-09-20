package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"brucosijada.kset.org/src/template/handlebars"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"

	flag "github.com/spf13/pflag"
)

//go:embed assets/*
var assets embed.FS

//go:embed views/*
var views embed.FS

const defaultPort = 3000
const defaultHost = "0.0.0.0"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
	}

	envPort, err := strconv.ParseInt(os.Getenv("PORT"), 0, 32)
	if err != nil || envPort == 0 {
		envPort = defaultPort
	}

	envHost := os.Getenv("HOST")
	if envHost == "" {
		envHost = defaultHost
	}

	var port int
	var host string

	flag.IntVarP(&port, "port", "p", int(envPort), "Set the port on which the server will run")
	flag.StringVarP(&host, "host", "h", envHost, "Set the host to which the server will bind")

	flag.Parse()

	// db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	//
	// if err != nil {
	// 	panic(err)
	// }

	viewsFolder, _ := fs.Sub(views, "views")
	engine := handlebars.NewFileSystem(
		http.FS(viewsFolder),
		".hbs",
	)

	app := fiber.New(fiber.Config{
		GETOnly:          true,
		DisableKeepalive: true,
		ReadTimeout:      10 * time.Second,
		ServerHeader:     "Microsoft-IIS/7.0",
		AppName:          "Brucifer 2021.",
		Views:            engine,
	})

	// println(db)

	app.Use(recover.New())
	app.Use(etag.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path}\t| ${status} ${latency}\t| ${ua}\t| ${ips} \n",
	}))
	app.Use(func(c *fiber.Ctx) error {
		// start timer
		start := time.Now()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now()
		// Do something with response
		c.Append("Server-Timing", fmt.Sprintf("app;dur=%v", stop.Sub(start).String()))
		// return stack error if exist
		return err
	})

	app.Use(helmet.New(helmet.Config{
		ReferrerPolicy:        "no-referrer-when-downgrade",
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' fonts.googleapis.com 'sha256-Ti7r9yBH+4QYBw011dOX3wsrskik/E0NNN3tNsdFTDw=' 'sha256-9sr6tDEbWwUiCy6TgwEAHgGGz2scGplIKDjkpkDanb0=' 'sha256-FVLanvfJQIpUweRHJtois7ZhhD1+NPuI/KD91VG0pcc='; font-src fonts.gstatic.com",
	}))

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

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", host, port)))
}
