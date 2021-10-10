package response

import (
	"github.com/gofiber/fiber/v2"
)

func Success(data interface{}) fiber.Map {
	return fiber.Map{
		"status": "success",
		"data":   data,
	}
}

func Error(message string, err error) fiber.Map {
	var data *string
	if err != nil {
		str := err.Error()
		data = &str
	}

	return fiber.Map{
		"status":  "error",
		"message": message,
		"data":    data,
	}
}
