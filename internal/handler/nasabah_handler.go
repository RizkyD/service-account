package handler

import (
	"account-service/internal/model"
	"account-service/internal/service"
	"account-service/internal/util"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type NasabahHandler struct {
	nasabahService service.NasabahService
	validate       *validator.Validate
}

func NewNasabahHandler(s service.NasabahService) *NasabahHandler {
	return &NasabahHandler{
		nasabahService: s,
		validate:       validator.New(),
	}
}

func (h *NasabahHandler) Daftar(c *fiber.Ctx) error {
	var nasabahDTO model.DaftarNasabahDTO

	if err := c.BodyParser(&nasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.validate.Struct(nasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	dataNasabah, err := h.nasabahService.Daftar(c.Context(), &nasabahDTO)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "telah digunakan") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terdapat kesalahan internal pada server")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"no_rekening": dataNasabah.ID,
	})
}

func (h *NasabahHandler) Tabung(c *fiber.Ctx) error {
	var updateSaldoNasabahDTO model.UpdateSaldoNasabahDTO

	if err := c.BodyParser(&updateSaldoNasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.validate.Struct(updateSaldoNasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	err := h.nasabahService.UpdateSaldo(c.Context(), &updateSaldoNasabahDTO, true)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

		if strings.Contains(err.Error(), "tidak mencukupi") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terdapat kesalahan internal pada server")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"saldo": updateSaldoNasabahDTO.Saldo,
	})
}

func (h *NasabahHandler) Tarik(c *fiber.Ctx) error {
	var updateSaldoNasabahDTO model.UpdateSaldoNasabahDTO

	if err := c.BodyParser(&updateSaldoNasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	if err := h.validate.Struct(updateSaldoNasabahDTO); err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Body request tidak valid")
	}

	err := h.nasabahService.UpdateSaldo(c.Context(), &updateSaldoNasabahDTO, false)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

		if strings.Contains(err.Error(), "tidak mencukupi") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}

		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terdapat kesalahan internal pada server")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"saldo": updateSaldoNasabahDTO.Saldo,
	})
}

func (h *NasabahHandler) GetSaldo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusBadRequest, "Nomor rekening tidak valid")
	}

	dataNasabah, err := h.nasabahService.GetSaldo(c.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "tidak ditemukan") {
			c.Locals("err", err)
			return util.ErrorResponse(c, fiber.StatusBadRequest, err.Error())
		}
		c.Locals("err", err)
		return util.ErrorResponse(c, fiber.StatusInternalServerError, "Terdapat kesalahan internal pada server")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"saldo": dataNasabah.Saldo,
	})
}
