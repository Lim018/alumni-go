package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type File struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    FileName     string             `json:"file_name" bson:"file_name"`
    OriginalName string             `json:"original_name" bson:"original_name"`
    FilePath     string             `json:"file_path" bson:"file_path"`
    FileSize     int64              `json:"file_size" bson:"file_size"`
    FileType     string             `json:"file_type" bson:"file_type"`
    FileCategory string             `json:"file_category" bson:"file_category"`
    AlumniID     primitive.ObjectID `json:"alumni_id" bson:"alumni_id"`
    UploadedBy   primitive.ObjectID `json:"uploaded_by" bson:"uploaded_by"`
    UploadedAt   time.Time          `json:"uploaded_at" bson:"uploaded_at"`
}

type FileUploadRequest struct {
    AlumniID string `json:"alumni_id" form:"alumni_id"`
}

type FileResponse struct {
    ID           string    `json:"id"`
    FileName     string    `json:"file_name"`
    OriginalName string    `json:"original_name"`
    FilePath     string    `json:"file_path"`
    FileSize     int64     `json:"file_size"`
    FileType     string    `json:"file_type"`
    FileCategory string    `json:"file_category"`
    AlumniID     string    `json:"alumni_id"`
    FileURL      string    `json:"file_url"`
    UploadedAt   time.Time `json:"uploaded_at"`
}

func (f *File) ToFileResponse(baseURL string) FileResponse {
    return FileResponse{
        ID:           f.ID.Hex(),
        FileName:     f.FileName,
        OriginalName: f.OriginalName,
        FilePath:     f.FilePath,
        FileSize:     f.FileSize,
        FileType:     f.FileType,
        FileCategory: f.FileCategory,
        AlumniID:     f.AlumniID.Hex(),
        FileURL:      baseURL + "/uploads/" + f.FileCategory + "/" + f.FileName,
        UploadedAt:   f.UploadedAt,
    }
}

type FileListResponse struct {
    Data []FileResponse `json:"data"`
    Meta MetaInfo       `json:"meta"`
}