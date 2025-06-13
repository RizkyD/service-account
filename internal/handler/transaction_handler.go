package handler

import (
	"account-service/internal/model"
	"account-service/internal/service"
	"account-service/internal/util"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service  service.TransactionService
	validate *validator.Validate
}

func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: s, validate: validator.New()}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var dto model.CreateTransactionDTO
	if err := c.BodyParser(&dto); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	if err := h.validate.Struct(dto); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	trx, err := h.service.Create(c.Context(), &dto)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan")
	}
	return c.Status(fiber.StatusOK).JSON(trx)
}

func (h *TransactionHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	trx, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		c.Locals("err", err)
		if err.Error() == "transaction not found" {
			return util.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan")
	}
	return c.Status(fiber.StatusOK).JSON(trx)
}

func (h *TransactionHandler) GetAll(c *fiber.Ctx) error {
	limitParam := c.Query("limit", "100")
	offsetParam := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	trxs, err := h.service.GetAll(c.Context(), limit, offset)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan")
	}
	return c.Status(fiber.StatusOK).JSON(trxs)
}

func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}

	var dto model.UpdateTransactionDTO
	if err = c.BodyParser(&dto); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}
	if err = h.validate.Struct(dto); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	err = h.service.Update(c.Context(), id, &dto)
	if err != nil {
		c.Locals("err", err)
		if err.Error() == "transaction not found" {
			return util.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan")
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "ID tidak valid")
	}
	err = h.service.Delete(c.Context(), id)
	if err != nil {
		c.Locals("err", err)
		if err.Error() == "transaction not found" {
			return util.ErrorResponse(c, fiber.StatusNotFound, err.Error())
		}
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan")
	}
	return c.SendStatus(fiber.StatusOK)
}
