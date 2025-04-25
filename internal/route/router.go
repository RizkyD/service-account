package route

import (
	"account-service/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, nasabahHandler *handler.NasabahHandler) {
	api := app.Group("/api/v1")
	SetupNasabahRoutes(api, nasabahHandler)

	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"remark": "Endpoint not found within API group",
		})
	})
}
