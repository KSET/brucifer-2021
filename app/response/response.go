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
	return fiber.Map{
		"status":  "error",
		"message": message,
		"data":    err.Error(),
	}
}
