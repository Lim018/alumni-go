package service

import (
    "strconv"
    "math"

    "github.com/gofiber/fiber/v2"
    "go-fiber/app/model"
    "go-fiber/app/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// CreateAlumniService godoc
// @Summary Create new alumni
// @Description Create a new alumni record
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body model.CreateAlumniRequest true "Alumni data"
// @Success 201 {object} map[string]interface{} "Alumni created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni [post]
func CreateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
    var req model.CreateAlumniRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Input tidak valid: " + err.Error(),
            "success": false,
        })
    }

    // Convert UserID string to ObjectID if provided
    var userID primitive.ObjectID
    if req.UserID != "" {
        var err error
        userID, err = primitive.ObjectIDFromHex(req.UserID)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "message": "User ID tidak valid",
                "success": false,
            })
        }
    }

    alumni := model.Alumni{
        NIM:        req.NIM,
        Nama:       req.Nama,
        Jurusan:    req.Jurusan,
        Angkatan:   req.Angkatan,
        TahunLulus: req.TahunLulus,
        Email:      req.Email,
        NoTelepon:  req.NoTelepon,
        Alamat:     req.Alamat,
        UserID:     userID,
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

// UpdateAlumniService godoc
// @Summary Update alumni
// @Description Update an existing alumni record
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Alumni ID"
// @Param body body model.UpdateAlumniRequest true "Alumni data"
// @Success 200 {object} map[string]interface{} "Alumni updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/{id} [put]
func UpdateAlumniService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
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

// DeleteAlumniService godoc
// @Summary Delete alumni
// @Description Delete an alumni record
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Alumni ID"
// @Success 200 {object} map[string]interface{} "Alumni deleted successfully"
// @Failure 400 {object} map[string]interface{} "Invalid ID"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/{id} [delete]
func DeleteAlumniService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
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

// GetAllAlumniServiceDatatable godoc
// @Summary Get all alumni with pagination
// @Description Get list of all alumni with pagination, sorting, and search
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Param sortBy query string false "Sort by field" default(_id)
// @Param order query string false "Sort order (asc/desc)" default(asc)
// @Param search query string false "Search keyword"
// @Success 200 {object} map[string]interface{} "List of alumni"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni [get]
func GetAllAlumniServiceDatatable(c *fiber.Ctx, db *mongo.Database) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "_id")
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

// GetAlumniStatsService godoc
// @Summary Get alumni statistics
// @Description Get alumni statistics grouped by jurusan
// @Tags Alumni
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Alumni statistics"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /alumni/stats/jurusan [get]
func GetAlumniStatsService(c *fiber.Ctx, db *mongo.Database) error {
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