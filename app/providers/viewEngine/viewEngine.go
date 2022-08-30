package viewEngine

import (
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"brucosijada.kset.org/src/template/handlebars"
)

func Init(fs fs.FS) fiber.Views {
	engine := handlebars.NewFileSystem(
		http.FS(fs),
		".hbs",
	)

	return engine
}
