package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"brucosijada.kset.org/app/models"
	"brucosijada.kset.org/app/providers/database"
	"brucosijada.kset.org/app/providers/hash"
	"brucosijada.kset.org/app/providers/minio"
	"brucosijada.kset.org/app/providers/session"
	"brucosijada.kset.org/app/providers/viewEngine"
	"brucosijada.kset.org/app/repo"
	"brucosijada.kset.org/app/routes"
	"brucosijada.kset.org/app/util/async"
	"brucosijada.kset.org/app/version"

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

type appConfig struct {
	Port int
	Host string
}

func loadConfig() appConfig {
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

	return appConfig{
		Port: port,
		Host: host,
	}
}

func cspConfig() string {
	cspEntries := map[string][]string{
		"default-src": {
			"'self'",
		},
		"img-src": {
			"'self'",
		},
		"media-src": {
			"'self'",
		},
		"style-src": {
			"'self'",
			"fonts.googleapis.com",
			"'unsafe-inline'",
		},
		"font-src": {
			"'self'",
			"fonts.gstatic.com",
		},
		"script-src": {
			"'unsafe-inline'",
		},
	}

	cspItems := make([]string, 0, len(cspEntries))
	for k, v := range cspEntries {
		cspItems = append(
			cspItems,
			fmt.Sprintf(
				"%s %s",
				k,
				strings.Join(v, " "),
			),
		)
	}

	return strings.Join(cspItems, "; ")
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file: ", err)
	}

	conf := loadConfig()

	_, err := async.Async().RunInParallel(
		func() (interface{}, error) {
			return nil, hash.HashProvider().CalibrateMemoryParam(time.Second)
		},
		func() (interface{}, error) {
			return nil, database.DatabaseProvider().Register()
		},
		func() (interface{}, error) {
			return nil, session.SessionProvider().Register()
		},
		func() (interface{}, error) {
			return nil, minio.MinioProvider().Register()
		},
	).All()
	if err != nil {
		panic(err)
	}

	err = func() (err error) {
		db := database.DatabaseProvider().Client()
		var users int64
		err = db.Model(models.User{}).Count(&users).Error
		if err != nil {
			return
		}

		if users == 0 {
			_, err = repo.User().Create(
				repo.UserCreateOptions{
					Email:      os.Getenv("INITIAL_USER_MAIL"),
					Identity:   os.Getenv("INITIAL_USER_USER"),
					Password:   os.Getenv("INITIAL_USER_PASS"),
					Invitation: nil,
				},
			)
		}

		return
	}()
	if err != nil {
		panic(err)
	}

	viewsFolder, _ := fs.Sub(views, "views")

	app := fiber.New(
		fiber.Config{
			DisableKeepalive: true,
			ReadTimeout:      10 * time.Second,
			ServerHeader:     "Microsoft-IIS/7.0",
			AppName:          fmt.Sprintf("Brucifer %d.", version.BuildTime().Year()),
			Views:            viewEngine.Init(viewsFolder),
			BodyLimit:        200 * 1024 * 1024 * 1,
			Prefork:          false,
		},
	)

	// println(db)

	app.Use(recover.New())
	app.Use(etag.New())
	app.Use(
		logger.New(
			logger.Config{
				Format: "[${time}] ${method} ${path}\t| ${status} ${latency}\t| ${ua}\t| ${ips} \n",
			},
		),
	)
	app.Use(
		func(c *fiber.Ctx) error {
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
		},
	)

	app.Use(
		helmet.New(
			helmet.Config{
				ReferrerPolicy:        "no-referrer-when-downgrade",
				ContentSecurityPolicy: cspConfig(),
			},
		),
	)

	app.Get(
		"/favicon.ico",
		func(c *fiber.Ctx) error {
			c.Type("ico")
			c.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int((7*24*time.Hour).Seconds())))

			favicon, _ := assets.ReadFile("assets/images/favicon.ico")

			return c.Send(favicon)
		},
	)

	assetsFs, _ := fs.Sub(assets, "assets")
	app.Use(
		"/assets",
		func(ctx *fiber.Ctx) error {
			ctx.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int((365*24*time.Hour).Seconds())))
			return ctx.Next()
		},
		filesystem.New(
			filesystem.Config{
				Root: http.FS(assetsFs),
			},
		),
	)

	routes.RegisterRoutes(app)

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", conf.Host, conf.Port)))
}
