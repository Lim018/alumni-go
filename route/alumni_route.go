package route

import (
	"alumni-go/app/model"
	"alumni-go/app/service"
	"alumni-go/helper"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniHandler struct {
	alumniService *service.AlumniService
}

func NewAlumniHandler() *AlumniHandler {
	return &AlumniHandler{
		alumniService: service.NewAlumniService(),
	}
}

// func (h *AlumniHandler) GetAll(c *fiber.Ctx) error {
// 	page, _ := strconv.Atoi(c.Query("page", "1"))
// 	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
// 	search := c.Query("search", "")

// 	alumni, meta, err := h.alumniService.GetAll(page, perPage, search)
// 	if err != nil {
// 		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
// 	}

// 	return helper.PaginatedSuccessResponse(c, "Alumni retrieved successfully", alumni, meta)
// }

func (h *AlumniHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	alumni, err := h.alumniService.GetByID(id)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusNotFound, err.Error())
	}

	return helper.SuccessResponse(c, "Alumni retrieved successfully", alumni)
}

func (h *AlumniHandler) Create(c *fiber.Ctx) error {
	var data model.AlumniCreate
	if err := c.BodyParser(&data); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if data.NIM == "" || data.Nama == "" || data.Email == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "NIM, nama, and email are required")
	}

	alumni, err := h.alumniService.Create(&data)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "NIM already exists" || err.Error() == "email already exists" {
			status = fiber.StatusConflict
		} else if err.Error() == "tahun lulus harus minimal 3 tahun setelah angkatan" {
			status = fiber.StatusBadRequest
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Alumni created successfully", alumni)
}

func (h *AlumniHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	var data model.AlumniUpdate
	if err := c.BodyParser(&data); err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Basic validation
	if data.Nama == "" || data.Email == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Nama and email are required")
	}

	alumni, err := h.alumniService.Update(id, &data)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "alumni not found" {
			status = fiber.StatusNotFound
		} else if err.Error() == "email already exists" {
			status = fiber.StatusConflict
		} else if err.Error() == "tahun lulus harus minimal 3 tahun setelah angkatan" {
			status = fiber.StatusBadRequest
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Alumni updated successfully", alumni)
}

func (h *AlumniHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ID")
	}

	err = h.alumniService.Delete(id)
	if err != nil {
		status := fiber.StatusInternalServerError
		if err.Error() == "alumni not found" {
			status = fiber.StatusNotFound
		}
		return helper.ErrorResponse(c, status, err.Error())
	}

	return helper.SuccessResponse(c, "Alumni deleted successfully", nil)
}

func (h *AlumniHandler) GetAlumniBaruBekerja(c *fiber.Ctx) error {
	data, err := h.alumniService.GetAlumniBaruBekerja()
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.SuccessResponse(c, "Successfully retrieved alumni working less than 3 years", data)
}

func (h *AlumniHandler) GetAll(c *fiber.Ctx) error {
	response, err := h.alumniService.GetAll(c)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}
	
	// Karena service sudah mengembalikan struct response lengkap, kita tinggal kirim
	return c.Status(fiber.StatusOK).JSON(response)
}