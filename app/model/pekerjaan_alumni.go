package model

import "time"

type PekerjaanAlumni struct {
	ID                   int       `json:"id" db:"id"`
	AlumniID             int       `json:"alumni_id" db:"alumni_id"`
	NamaPerusahaan       string    `json:"nama_perusahaan" db:"nama_perusahaan"`
	PosisiJabatan        string    `json:"posisi_jabatan" db:"posisi_jabatan"`
	BidangIndustri       string    `json:"bidang_industri" db:"bidang_industri"`
	LokasiKerja          string    `json:"lokasi_kerja" db:"lokasi_kerja"`
	GajiRange            *string   `json:"gaji_range" db:"gaji_range"`
	TanggalMulaiKerja    time.Time `json:"tanggal_mulai_kerja" db:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja  *time.Time `json:"tanggal_selesai_kerja" db:"tanggal_selesai_kerja"`
	StatusPekerjaan      string    `json:"status_pekerjaan" db:"status_pekerjaan"`
	DeskripsiPekerjaan   *string   `json:"deskripsi_pekerjaan" db:"deskripsi_pekerjaan"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type PekerjaanAlumniWithAlumni struct {
	PekerjaanAlumni
	Alumni Alumni `json:"alumni"`
}

type PekerjaanCreate struct {
	AlumniID             int        `json:"alumni_id" validate:"required"`
	NamaPerusahaan       string     `json:"nama_perusahaan" validate:"required,min=2,max=100"`
	PosisiJabatan        string     `json:"posisi_jabatan" validate:"required,min=2,max=100"`
	BidangIndustri       string     `json:"bidang_industri" validate:"required,min=2,max=50"`
	LokasiKerja          string     `json:"lokasi_kerja" validate:"required,min=2,max=100"`
	GajiRange            *string    `json:"gaji_range" validate:"omitempty,max=50"`
	TanggalMulaiKerja    time.Time  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja  *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan      string     `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan   *string    `json:"deskripsi_pekerjaan"`
}

type PekerjaanUpdate struct {
	NamaPerusahaan       string     `json:"nama_perusahaan" validate:"required,min=2,max=100"`
	PosisiJabatan        string     `json:"posisi_jabatan" validate:"required,min=2,max=100"`
	BidangIndustri       string     `json:"bidang_industri" validate:"required,min=2,max=50"`
	LokasiKerja          string     `json:"lokasi_kerja" validate:"required,min=2,max=100"`
	GajiRange            *string    `json:"gaji_range" validate:"omitempty,max=50"`
	TanggalMulaiKerja    time.Time  `json:"tanggal_mulai_kerja" validate:"required"`
	TanggalSelesaiKerja  *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan      string     `json:"status_pekerjaan" validate:"required,oneof=aktif selesai resigned"`
	DeskripsiPekerjaan   *string    `json:"deskripsi_pekerjaan"`
}