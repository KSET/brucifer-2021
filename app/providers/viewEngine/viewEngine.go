package viewEngine

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/app/version"
	"brucosijada.kset.org/src/template/handlebars"
)

func Init(fs fs.FS) fiber.Views {
	engine := handlebars.NewFileSystem(
		http.FS(fs),
		".hbs",
	)

	buildTime := version.BuildTime()
	buildTimeFormatted := buildTime.UTC().Format(version.TimeFormat)

	engine.AddFunc(
		"_currentYear",
		func() int {
			return time.Now().Year()
		},
	)

	engine.AddFunc(
		"_buildTime",
		func() string {
			return buildTimeFormatted
		},
	)

	return engine
}
