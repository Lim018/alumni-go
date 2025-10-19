package model

import "time"

// type Pekerjaan struct {
// 	ID                 int        `json:"id"`
// 	AlumniID           int        `json:"alumni_id"`
// 	NamaPerusahaan     string     `json:"nama_perusahaan"`
// 	PosisiJabatan      string     `json:"posisi_jabatan"`
// 	BidangIndustri     string     `json:"bidang_industri"`
// 	LokasiKerja        string     `json:"lokasi_kerja"`
// 	GajiRange          string     `json:"gaji_range"`
// 	TanggalMulaiKerja  time.Time  `json:"tanggal_mulai_kerja"`
// 	TanggalSelesaiKerja *time.Time `json:"tanggal_selesai_kerja"`
// 	StatusPekerjaan    string     `json:"status_pekerjaan"`
// 	DeskripsiPekerjaan string     `json:"deskripsi_pekerjaan"`
// 	CreatedAt          time.Time  `json:"created_at"`
// 	UpdatedAt          time.Time  `json:"updated_at"`
// 	IsDelete           *time.Time `json:"is_delete"`
// }

// Pekerjaan - Base model for database representation
type Pekerjaan struct {
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
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
	IsDelete            *time.Time `json:"is_delete"`
}

// CreatePekerjaanRequest - Request for POST /pekerjaan
type CreatePekerjaanRequest struct {
	AlumniID            int        `json:"alumni_id" validate:"required"`
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
	AlumniID            int        `json:"alumni_id" validate:"required"`
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
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
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

// PekerjaanWithAlumniResponse - Response for GET /pekerjaan/alumni/:alumni_id
type PekerjaanWithAlumniResponse struct {
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
	AlumniNama          string     `json:"alumni_nama"`
	AlumniNIM           string     `json:"alumni_nim"`
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

// PekerjaanListResponse - Response for GET /pekerjaan (datatable)
type PekerjaanListResponse struct {
	Data []PekerjaanResponse `json:"data"`
	Meta MetaInfo            `json:"meta"`
}

// PekerjaanTrashResponse - Response for trash items (includes IsDelete)
type PekerjaanTrashResponse struct {
	ID                  int        `json:"id"`
	AlumniID            int        `json:"alumni_id"`
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
	DeletedAt           time.Time  `json:"deleted_at"` // Renamed from IsDelete for clarity
}

// PekerjaanTrashListResponse - Response for GET /pekerjaan/trash
type PekerjaanTrashListResponse struct {
	Data []PekerjaanTrashResponse `json:"data"`
	Meta MetaInfo                 `json:"meta"`
}