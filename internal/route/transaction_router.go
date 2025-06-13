package route

import (
	"account-service/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupTransactionRoutes(router fiber.Router, h *handler.TransactionHandler) {
	trGroup := router.Group("/transactions")
	trGroup.Post("/", h.Create)      // POST /api/v1/transactions
	trGroup.Get("/:id", h.GetByID)   // GET /api/v1/transactions/:id
	trGroup.Get("/", h.GetAll)       // GET /api/v1/transactions
	trGroup.Put("/:id", h.Update)    // PUT /api/v1/transactions/:id
	trGroup.Delete("/:id", h.Delete) // DELETE /api/v1/transactions/:id
}
