package service

import (
	"database/sql"
	"strconv"
	"math"

	"github.com/gofiber/fiber/v2"
	"go-fiber/app/model"
	"go-fiber/app/repository"
)

func CreateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	var req model.CreateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Input tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Convert request to database model
	alumni := model.Alumni{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	newAlumni, err := repository.CreateAlumni(db, alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menambahkan alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Alumni berhasil ditambahkan",
		"success": true,
		"data":    newAlumni.ToAlumniResponse(),
	})
}

func UpdateAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	var req model.UpdateAlumniRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Input tidak valid: " + err.Error(),
			"success": false,
		})
	}

	// Convert request to database model
	alumni := model.Alumni{
		NIM:        req.NIM,
		Nama:       req.Nama,
		Jurusan:    req.Jurusan,
		Angkatan:   req.Angkatan,
		TahunLulus: req.TahunLulus,
		Email:      req.Email,
		NoTelepon:  req.NoTelepon,
		Alamat:     req.Alamat,
	}

	updatedAlumni, err := repository.UpdateAlumni(db, id, alumni)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal update alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Alumni berhasil diupdate",
		"success": true,
		"data":    updatedAlumni.ToAlumniResponse(),
	})
}

func DeleteAlumniService(c *fiber.Ctx, db *sql.DB) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID tidak valid",
			"success": false,
		})
	}

	if err := repository.DeleteAlumni(db, id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menghapus alumni: " + err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Alumni berhasil dihapus",
		"success": true,
	})
}

func GetAllAlumniServiceDatatable(c *fiber.Ctx, db *sql.DB) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	alumniList, err := repository.GetAlumniRepo(db, search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mendapatkan data alumni: " + err.Error(),
			"success": false,
		})
	}

	total, err := repository.CountAlumniRepo(db, search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menghitung total alumni: " + err.Error(),
			"success": false,
		})
	}

	// Convert to response models
	responses := make([]model.AlumniResponse, len(alumniList))
	for i, alumni := range alumniList {
		responses[i] = alumni.ToAlumniResponse()
	}

	meta := model.MetaInfo{
		Page:   page,
		Limit:  limit,
		Total:  total,
		Pages:  int(math.Ceil(float64(total) / float64(limit))),
		SortBy: sortBy,
		Order:  order,
		Search: search,
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mendapatkan data alumni",
		"success": true,
		"data":    responses,
		"meta":    meta,
	})
}

func GetAlumniStatsService(c *fiber.Ctx, db *sql.DB) error {
	stats, err := repository.GetAlumniStatsByJurusan(db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal mendapatkan statistik: " + err.Error(),
			"success": false,
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mendapatkan statistik alumni",
		"success": true,
		"data":    stats,
	})
}