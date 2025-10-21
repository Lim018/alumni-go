package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// Pekerjaan - Base model for MongoDB
type Pekerjaan struct {
    ID                  primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
    AlumniID            primitive.ObjectID  `json:"alumni_id" bson:"alumni_id"`
    NamaPerusahaan      string              `json:"nama_perusahaan" bson:"nama_perusahaan"`
    PosisiJabatan       string              `json:"posisi_jabatan" bson:"posisi_jabatan"`
    BidangIndustri      string              `json:"bidang_industri" bson:"bidang_industri"`
    LokasiKerja         string              `json:"lokasi_kerja" bson:"lokasi_kerja"`
    GajiRange           string              `json:"gaji_range" bson:"gaji_range"`
    TanggalMulaiKerja   time.Time           `json:"tanggal_mulai_kerja" bson:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time          `json:"tanggal_selesai_kerja" bson:"tanggal_selesai_kerja"`
    StatusPekerjaan     string              `json:"status_pekerjaan" bson:"status_pekerjaan"`
    DeskripsiPekerjaan  string              `json:"deskripsi_pekerjaan" bson:"deskripsi_pekerjaan"`
    CreatedAt           time.Time           `json:"created_at" bson:"created_at"`
    UpdatedAt           time.Time           `json:"updated_at" bson:"updated_at"`
    IsDelete            *time.Time          `json:"is_delete,omitempty" bson:"is_delete,omitempty"`
}

// CreatePekerjaanRequest - Request for POST /pekerjaan
type CreatePekerjaanRequest struct {
    AlumniID            string     `json:"alumni_id" validate:"required"`
    NamaPerusahaan      string     `json:"nama_perusahaan" validate:"required"`
    PosisiJabatan       string     `json:"posisi_jabatan" validate:"required"`
    BidangIndustri      string     `json:"bidang_industri" validate:"required"`
    LokasiKerja         string     `json:"lokasi_kerja" validate:"required"`
    GajiRange           string     `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja" validate:"required"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan" validate:"required"`
    DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
}

// UpdatePekerjaanRequest - Request for PUT /pekerjaan/:id
type UpdatePekerjaanRequest struct {
    AlumniID            string     `json:"alumni_id" validate:"required"`
    NamaPerusahaan      string     `json:"nama_perusahaan" validate:"required"`
    PosisiJabatan       string     `json:"posisi_jabatan" validate:"required"`
    BidangIndustri      string     `json:"bidang_industri" validate:"required"`
    LokasiKerja         string     `json:"lokasi_kerja" validate:"required"`
    GajiRange           string     `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja" validate:"required"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan" validate:"required"`
    DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
}

// PekerjaanResponse - Response for single pekerjaan
type PekerjaanResponse struct {
    ID                  string     `json:"id"`
    AlumniID            string     `json:"alumni_id"`
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           string     `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
}

// PekerjaanTrashResponse - Response for trash items
type PekerjaanTrashResponse struct {
    ID                  string     `json:"id"`
    AlumniID            string     `json:"alumni_id"`
    NamaPerusahaan      string     `json:"nama_perusahaan"`
    PosisiJabatan       string     `json:"posisi_jabatan"`
    BidangIndustri      string     `json:"bidang_industri"`
    LokasiKerja         string     `json:"lokasi_kerja"`
    GajiRange           string     `json:"gaji_range"`
    TanggalMulaiKerja   time.Time  `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
    StatusPekerjaan     string     `json:"status_pekerjaan"`
    DeskripsiPekerjaan  string     `json:"deskripsi_pekerjaan"`
    CreatedAt           time.Time  `json:"created_at"`
    UpdatedAt           time.Time  `json:"updated_at"`
    DeletedAt           time.Time  `json:"deleted_at"`
}

// PekerjaanListResponse - Response for GET /pekerjaan (datatable)
type PekerjaanListResponse struct {
    Data []PekerjaanResponse `json:"data"`
    Meta MetaInfo            `json:"meta"`
}

// PekerjaanTrashListResponse - Response for GET /pekerjaan/trash
type PekerjaanTrashListResponse struct {
    Data []PekerjaanTrashResponse `json:"data"`
    Meta MetaInfo                 `json:"meta"`
}

// MetaInfo for pagination
