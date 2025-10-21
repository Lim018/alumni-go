package service

import (
    "math"
    "strconv"
    
    "github.com/gofiber/fiber/v2"
    "go-fiber/app/model"
    "go-fiber/app/repository"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

func GetPekerjaanByAlumniIDService(c *fiber.Ctx, db *mongo.Database) error {
    alumniIDStr := c.Params("alumni_id")
    alumniID, err := primitive.ObjectIDFromHex(alumniIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "ID alumni tidak valid",
            "success": false,
        })
    }

    repo := repository.NewPekerjaanRepository(db)
    pekerjaanList, err := repo.FindPekerjaanByAlumniID(alumniID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal mendapatkan data pekerjaan: " + err.Error(),
            "success": false,
        })
    }

    responses := make([]model.PekerjaanResponse, len(pekerjaanList))
    for i, pekerjaan := range pekerjaanList {
        responses[i] = pekerjaan.ToPekerjaanResponse()
    }

    return c.JSON(fiber.Map{
        "message": "Berhasil mendapatkan data pekerjaan untuk alumni",
        "success": true,
        "data":    responses,
    })
}

func CreatePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    var req model.CreatePekerjaanRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Input tidak valid: " + err.Error(),
            "success": false,
        })
    }

    alumniID, err := primitive.ObjectIDFromHex(req.AlumniID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
			"message": "Alumni ID tidak valid",
            "success": false,
        })
    }

    pekerjaan := model.Pekerjaan{
        AlumniID:            alumniID,
        NamaPerusahaan:      req.NamaPerusahaan,
        PosisiJabatan:       req.PosisiJabatan,
        BidangIndustri:      req.BidangIndustri,
        LokasiKerja:         req.LokasiKerja,
        GajiRange:           req.GajiRange,
        TanggalMulaiKerja:   req.TanggalMulaiKerja,
        TanggalSelesaiKerja: req.TanggalSelesaiKerja,
        StatusPekerjaan:     req.StatusPekerjaan,
        DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
    }

    repo := repository.NewPekerjaanRepository(db)
    newPekerjaan, err := repo.CreatePekerjaan(pekerjaan)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal menambahkan pekerjaan: " + err.Error(),
            "success": false,
        })
    }

    return c.Status(201).JSON(fiber.Map{
        "message": "Pekerjaan berhasil ditambahkan",
        "success": true,
        "data":    newPekerjaan.ToPekerjaanResponse(),
    })
}

func UpdatePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "ID tidak valid",
            "success": false,
        })
    }

    var req model.UpdatePekerjaanRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Input tidak valid: " + err.Error(),
            "success": false,
        })
    }

    alumniID, err := primitive.ObjectIDFromHex(req.AlumniID)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Alumni ID tidak valid",
            "success": false,
        })
    }

    pekerjaan := model.Pekerjaan{
        AlumniID:            alumniID,
        NamaPerusahaan:      req.NamaPerusahaan,
        PosisiJabatan:       req.PosisiJabatan,
        BidangIndustri:      req.BidangIndustri,
        LokasiKerja:         req.LokasiKerja,
        GajiRange:           req.GajiRange,
        TanggalMulaiKerja:   req.TanggalMulaiKerja,
        TanggalSelesaiKerja: req.TanggalSelesaiKerja,
        StatusPekerjaan:     req.StatusPekerjaan,
        DeskripsiPekerjaan:  req.DeskripsiPekerjaan,
    }

    repo := repository.NewPekerjaanRepository(db)
    updatedPekerjaan, err := repo.UpdatePekerjaan(id, pekerjaan)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal update pekerjaan: " + err.Error(),
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil diupdate",
        "success": true,
        "data":    updatedPekerjaan.ToPekerjaanResponse(),
    })
}

func GetAllPekerjaanServiceDatatable(c *fiber.Ctx, db *mongo.Database) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "_id")
    order := c.Query("order", "asc")
    search := c.Query("search", "")

    if page < 1 {
        page = 1
    }
    offset := (page - 1) * limit

    repo := repository.NewPekerjaanRepository(db)
    list, err := repo.GetPekerjaan(search, sortBy, order, limit, offset)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal mendapatkan data pekerjaan alumni: " + err.Error(),
            "success": false,
        })
    }

    total, err := repo.CountPekerjaan(search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal menghitung total pekerjaan alumni: " + err.Error(),
            "success": false,
        })
    }

    responses := make([]model.PekerjaanResponse, len(list))
    for i, pekerjaan := range list {
        responses[i] = pekerjaan.ToPekerjaanResponse()
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
        "message": "Berhasil mendapatkan data pekerjaan alumni",
        "success": true,
        "data":    responses,
        "meta":    meta,
    })
}

func SoftDeletePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error":   "ID tidak valid",
            "success": false,
        })
    }

    userIDInterface := c.Locals("user_id")
    var userID primitive.ObjectID
    
    // Handle different types of user_id from middleware
    switch v := userIDInterface.(type) {
    case primitive.ObjectID:
        userID = v
    case string:
        userID, err = primitive.ObjectIDFromHex(v)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error":   "User ID tidak valid",
                "success": false,
            })
        }
    case int:
        // If your middleware stores as int, you need to query user collection
        return c.Status(400).JSON(fiber.Map{
            "error":   "Format User ID tidak didukung",
            "success": false,
        })
    }

    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    repo := repository.NewPekerjaanRepository(db)
    err = repo.SoftDelete(id, userID, isAdmin)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error":   "Gagal soft delete pekerjaan: " + err.Error(),
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil dihapus (soft delete)",
        "success": true,
    })
}

func GetTrashPekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    sortBy := c.Query("sortBy", "is_delete")
    order := c.Query("order", "desc")
    search := c.Query("search", "")

    if page < 1 {
        page = 1
    }
    offset := (page - 1) * limit

    userIDInterface := c.Locals("user_id")
    var userID primitive.ObjectID
    var err error
    
    switch v := userIDInterface.(type) {
    case primitive.ObjectID:
        userID = v
    case string:
        userID, err = primitive.ObjectIDFromHex(v)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error":   "User ID tidak valid",
                "success": false,
            })
        }
    }

    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    repo := repository.NewPekerjaanRepository(db)
    list, err := repo.GetTrashPekerjaan(userID, isAdmin, search, sortBy, order, limit, offset)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal mendapatkan data trash: " + err.Error(),
            "success": false,
        })
    }

    total, err := repo.CountTrashPekerjaan(userID, isAdmin, search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Gagal menghitung total trash: " + err.Error(),
            "success": false,
        })
    }

    responses := make([]model.PekerjaanTrashResponse, len(list))
    for i, pekerjaan := range list {
        responses[i] = pekerjaan.ToPekerjaanTrashResponse()
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
        "message": "Berhasil mendapatkan data trash",
        "success": true,
        "data":    responses,
        "meta":    meta,
    })
}

func RestorePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "ID tidak valid",
            "success": false,
        })
    }

    userIDInterface := c.Locals("user_id")
    var userID primitive.ObjectID
    
    switch v := userIDInterface.(type) {
    case primitive.ObjectID:
        userID = v
    case string:
        userID, err = primitive.ObjectIDFromHex(v)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error":   "User ID tidak valid",
                "success": false,
            })
        }
    }

    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    repo := repository.NewPekerjaanRepository(db)
    err = repo.RestorePekerjaan(id, userID, isAdmin)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "Data tidak ditemukan atau bukan milik anda",
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil dikembalikan",
        "success": true,
    })
}

func HardDeletePekerjaanService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "ID tidak valid",
            "success": false,
        })
    }

    userIDInterface := c.Locals("user_id")
    var userID primitive.ObjectID
    
    switch v := userIDInterface.(type) {
    case primitive.ObjectID:
        userID = v
    case string:
        userID, err = primitive.ObjectIDFromHex(v)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "error":   "User ID tidak valid",
                "success": false,
            })
        }
    }

    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    repo := repository.NewPekerjaanRepository(db)
    err = repo.HardDeletePekerjaan(id, userID, isAdmin)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "message": "Data tidak ditemukan atau bukan milik anda",
            "success": false,
        })
    }

    return c.JSON(fiber.Map{
        "message": "Pekerjaan berhasil dihapus permanen",
        "success": true,
    })
}