package route

import (
	"alumni-go/app/model"
	"alumni-go/app/service"
	"alumni-go/helper"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanHandler struct {
	pekerjaanService *service.PekerjaanAlumniService
}

func NewPekerjaanHandler() *PekerjaanHandler {
	return &PekerjaanHandler{
		pekerjaanService: service.NewPekerjaanAlumniService(),
	}
}

func (h *PekerjaanHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")

	pekerjaan, meta, err := h.pekerjaanService.GetAll(page, perPage, search)
	if err != nil {
		return helper.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return helper.PaginatedSuccessResponse(c, "Pekerjaan retrieved successfully", pekerjaan, meta)
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

	// Parse tanggal from string if needed
	if c.FormValue("tanggal_mulai_kerja") != "" {
		if t, err := time.Parse("2006-01-02", c.FormValue("tanggal_mulai_kerja")); err == nil {
			data.TanggalMulaiKerja = t
		}
	}
	if c.FormValue("tanggal_selesai_kerja") != "" {
		if t, err := time.Parse("2006-01-02", c.FormValue("tanggal_selesai_kerja")); err == nil {
			data.TanggalSelesaiKerja = &t
		}
	}

	// Basic validation
	if data.AlumniID == 0 || data.NamaPerusahaan == "" || data.PosisiJabatan == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Alumni ID, nama perusahaan, and posisi jabatan are required")
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

	// Parse tanggal from string if needed
	if c.FormValue("tanggal_mulai_kerja") != "" {
		if t, err := time.Parse("2006-01-02", c.FormValue("tanggal_mulai_kerja")); err == nil {
			data.TanggalMulaiKerja = t
		}
	}
	if c.FormValue("tanggal_selesai_kerja") != "" {
		if t, err := time.Parse("2006-01-02", c.FormValue("tanggal_selesai_kerja")); err == nil {
			data.TanggalSelesaiKerja = &t
		}
	}

	// Basic validation
	if data.NamaPerusahaan == "" || data.PosisiJabatan == "" {
		return helper.ErrorResponse(c, fiber.StatusBadRequest, "Nama perusahaan and posisi jabatan are required")
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
