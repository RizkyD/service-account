package route

import (
	"account-service/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupNasabahRoutes(router fiber.Router, nasabahHandler *handler.NasabahHandler) {
	nasabahGroup := router.Group("/")

	nasabahGroup.Post("/daftar", nasabahHandler.Daftar)     // GET /api/v1/daftar
	nasabahGroup.Put("/tabung", nasabahHandler.Tabung)      // PUT /api/v1/tabung
	nasabahGroup.Put("/tarik", nasabahHandler.Tarik)        // PUT /api/v1/tarik
	nasabahGroup.Get("/saldo/:id", nasabahHandler.GetSaldo) // PUT /api/v1/saldo
}
