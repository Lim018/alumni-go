package service

import (
    "fmt"
    "math"
    "os"
    "path/filepath"
    "strconv"
    
    "go-fiber/app/model"
    "go-fiber/app/repository"
    
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
)

// UploadFotoService - Upload foto dengan validasi JPEG/JPG/PNG max 1MB
func UploadFotoService(c *fiber.Ctx, db *mongo.Database) error {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded",
            "error":   err.Error(),
        })
    }

    maxSize := int64(1 * 1024 * 1024)
    if fileHeader.Size > maxSize {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File size exceeds 1MB",
        })
    }

    allowedTypes := map[string]bool{
        "image/jpeg": true,
        "image/jpg":  true,
        "image/png":  true,
    }

    contentType := fileHeader.Header.Get("Content-Type")
    if !allowedTypes[contentType] {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File type not allowed. Only JPEG, JPG, PNG are accepted",
        })
    }

    userID := c.Locals("user_id").(primitive.ObjectID)
    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    alumniIDStr := c.FormValue("alumni_id")
    var alumniID primitive.ObjectID

    if isAdmin {
        if alumniIDStr == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "alumni_id is required for admin",
            })
        }
        alumniID, err = primitive.ObjectIDFromHex(alumniIDStr)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Invalid alumni_id",
            })
        }
    } else {
        alumniRepo := repository.NewAlumniRepository(db)
        alumni, err := alumniRepo.GetAlumniByUserID(userID)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Alumni record not found for this user",
            })
        }
        alumniID = alumni.ID
    }

    uploadPath := "./uploads/foto"
    if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to create upload directory",
            "error":   err.Error(),
        })
    }

    ext := filepath.Ext(fileHeader.Filename)
    newFileName := uuid.New().String() + ext
    filePath := filepath.Join(uploadPath, newFileName)

    if err := c.SaveFile(fileHeader, filePath); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file",
            "error":   err.Error(),
        })
    }

    fileModel := &model.File{
        FileName:     newFileName,
        OriginalName: fileHeader.Filename,
        FilePath:     filePath,
        FileSize:     fileHeader.Size,
        FileType:     contentType,
        FileCategory: "foto",
        AlumniID:     alumniID,
        UploadedBy:   userID,
    }

    fileRepo := repository.NewFileRepository(db)
    if err := fileRepo.Create(fileModel); err != nil {
        os.Remove(filePath)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file metadata",
            "error":   err.Error(),
        })
    }

    baseURL := c.BaseURL()
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "Foto uploaded successfully",
        "data":    fileModel.ToFileResponse(baseURL),
    })
}

// UploadSertifikatService - Upload sertifikat dengan validasi PDF max 2MB
func UploadSertifikatService(c *fiber.Ctx, db *mongo.Database) error {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded",
            "error":   err.Error(),
        })
    }

    maxSize := int64(2 * 1024 * 1024)
    if fileHeader.Size > maxSize {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File size exceeds 2MB",
        })
    }

    contentType := fileHeader.Header.Get("Content-Type")
    if contentType != "application/pdf" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File type not allowed. Only PDF is accepted",
        })
    }

    userID := c.Locals("user_id").(primitive.ObjectID)
    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    alumniIDStr := c.FormValue("alumni_id")
    var alumniID primitive.ObjectID

    if isAdmin {
        if alumniIDStr == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "alumni_id is required for admin",
            })
        }
        alumniID, err = primitive.ObjectIDFromHex(alumniIDStr)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "success": false,
                "message": "Invalid alumni_id",
            })
        }
    } else {
        alumniRepo := repository.NewAlumniRepository(db)
        alumni, err := alumniRepo.GetAlumniByUserID(userID)
        if err != nil {
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "success": false,
                "message": "Alumni record not found for this user",
            })
        }
        alumniID = alumni.ID
    }

    uploadPath := "./uploads/sertifikat"
    if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to create upload directory",
            "error":   err.Error(),
        })
    }

    ext := filepath.Ext(fileHeader.Filename)
    newFileName := uuid.New().String() + ext
    filePath := filepath.Join(uploadPath, newFileName)

    if err := c.SaveFile(fileHeader, filePath); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file",
            "error":   err.Error(),
        })
    }

    fileModel := &model.File{
        FileName:     newFileName,
        OriginalName: fileHeader.Filename,
        FilePath:     filePath,
        FileSize:     fileHeader.Size,
        FileType:     contentType,
        FileCategory: "sertifikat",
        AlumniID:     alumniID,
        UploadedBy:   userID,
    }

    fileRepo := repository.NewFileRepository(db)
    if err := fileRepo.Create(fileModel); err != nil {
        os.Remove(filePath)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file metadata",
            "error":   err.Error(),
        })
    }

    baseURL := c.BaseURL()
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "uploaded successfully",
        "data":    fileModel.ToFileResponse(baseURL),
    })
}

// GetFilesByAlumniService - Get all files for specific alumni
func GetFilesByAlumniService(c *fiber.Ctx, db *mongo.Database) error {
    alumniIDStr := c.Params("alumni_id")
    alumniID, err := primitive.ObjectIDFromHex(alumniIDStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid alumni_id",
        })
    }

    fileCategory := c.Query("category", "") // "foto" or "sertifikat" or empty for all

    fileRepo := repository.NewFileRepository(db)
    files, err := fileRepo.FindByAlumniID(alumniID, fileCategory)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get files",
            "error":   err.Error(),
        })
    }

    baseURL := c.BaseURL()
    var responses []model.FileResponse
    for _, file := range files {
        responses = append(responses, file.ToFileResponse(baseURL))
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Files retrieved successfully",
        "data":    responses,
    })
}

// GetAllFilesService - Get all files with pagination
func GetAllFilesService(c *fiber.Ctx, db *mongo.Database) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    search := c.Query("search", "")
    fileCategory := c.Query("category", "") // "foto" or "sertifikat"

    if page < 1 {
        page = 1
    }
    offset := (page - 1) * limit

    repo := repository.NewFileRepository(db)
    
    files, err := repo.FindAll(search, fileCategory, limit, offset)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get files",
            "error":   err.Error(),
        })
    }

    total, err := repo.Count(search, fileCategory)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to count files",
            "error":   err.Error(),
        })
    }

    baseURL := c.BaseURL()
    var responses []model.FileResponse
    for _, file := range files {
        responses = append(responses, file.ToFileResponse(baseURL))
    }

    meta := model.MetaInfo{
        Page:   page,
        Limit:  limit,
        Total:  total,
        Pages:  int(math.Ceil(float64(total) / float64(limit))),
        Search: search,
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Files retrieved successfully",
        "data":    responses,
        "meta":    meta,
    })
}

// DeleteFileService - Delete file with authorization check
func DeleteFileService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid file ID",
        })
    }

    userID := c.Locals("user_id").(primitive.ObjectID)
    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    fileRepo := repository.NewFileRepository(db)
    
    // Check authorization
    if !isAdmin {
        isOwner, err := fileRepo.IsFileOwnedByUser(id, userID)
        if err != nil || !isOwner {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "success": false,
                "message": "You don't have permission to delete this file",
            })
        }
    }

    file, err := fileRepo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "File not found",
            "error":   err.Error(),
        })
    }

    // Hapus file dari storage
    if err := os.Remove(file.FilePath); err != nil {
        fmt.Println("Warning: Failed to delete file from storage:", err)
    }

    // Hapus dari database
    if err := fileRepo.Delete(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to delete file",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "File deleted successfully",
    })
}

// DownloadFileService - Download file with authorization check
func DownloadFileService(c *fiber.Ctx, db *mongo.Database) error {
    idStr := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idStr)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "Invalid file ID",
        })
    }

    userID := c.Locals("user_id").(primitive.ObjectID)
    role := c.Locals("role").(string)
    isAdmin := role == "admin"

    fileRepo := repository.NewFileRepository(db)
    
    // Check authorization
    if !isAdmin {
        isOwner, err := fileRepo.IsFileOwnedByUser(id, userID)
        if err != nil || !isOwner {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
                "success": false,
                "message": "You don't have permission to download this file",
            })
        }
    }

    file, err := fileRepo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "File not found",
        })
    }

    // Check if file exists
    if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "File not found in storage",
        })
    }

    // Set headers for download
    c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))
    c.Set("Content-Type", file.FileType)

    return c.SendFile(file.FilePath)
}