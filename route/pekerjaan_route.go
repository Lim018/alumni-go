package route

import (
	"alumni-go/app/model"
	"alumni-go/app/service"
	"alumni-go/helper"
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanHandler struct {
	pekerjaanService *service.PekerjaanAlumniService
}

// NewPekerjaanHandler sekarang menerima *sql.DB
func NewPekerjaanHandler(db *sql.DB) *PekerjaanHandler {
	return &PekerjaanHandler{
		pekerjaanService: service.NewPekerjaanAlumniService(db),
	}
}

func (h *PekerjaanHandler) GetAll(c *fiber.Ctx) error {
	response, err := h.pekerjaanService.GetAll(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *PekerjaanHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	pekerjaan, err := h.pekerjaanService.GetByID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan retrieved successfully", pekerjaan)
}

func (h *PekerjaanHandler) GetByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid Alumni ID")
	}

	pekerjaan, err := h.pekerjaanService.GetByAlumniID(alumniID)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "alumni not found" {
			status = fiber.StatusNotFound
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan by alumni retrieved successfully", pekerjaan)
}

func (h *PekerjaanHandler) Create(c *fiber.Ctx) error {
	var data model.PekerjaanCreate
	if err := c.BodyParser(&data); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	pekerjaan, err := h.pekerjaanService.Create(&data)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "alumni not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "tanggal mulai kerja tidak boleh setelah tanggal selesai kerja" {
			status = fiber.StatusBadRequest
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan created successfully", pekerjaan)
}

func (h *PekerjaanHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	var data model.PekerjaanUpdate
	if err := c.BodyParser(&data); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	pekerjaan, err := h.pekerjaanService.Update(id, &data)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "pekerjaan not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "tanggal mulai kerja tidak boleh setelah tanggal selesai kerja" {
			status = fiber.StatusBadRequest
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan updated successfully", pekerjaan)
}

func (h *PekerjaanHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	err = h.pekerjaanService.Delete(id)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "pekerjaan not found" {
			status = fiber.StatusNotFound
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Pekerjaan deleted successfully", nil)
}
