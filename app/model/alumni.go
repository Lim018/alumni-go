package model

import "time"

type Alumni struct {
	ID         int       `json:"id" db:"id"`
	NIM        string    `json:"nim" db:"nim"`
	Nama       string    `json:"nama" db:"nama"`
	Jurusan    string    `json:"jurusan" db:"jurusan"`
	Angkatan   int       `json:"angkatan" db:"angkatan"`
	TahunLulus int       `json:"tahun_lulus" db:"tahun_lulus"`
	Email      string    `json:"email" db:"email"`
	NoTelepon  *string   `json:"no_telepon" db:"no_telepon"`
	Alamat     *string   `json:"alamat" db:"alamat"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type AlumniCreate struct {
	NIM        string  `json:"nim" validate:"required,min=5,max=20"`
	Nama       string  `json:"nama" validate:"required,min=2,max=100"`
	Jurusan    string  `json:"jurusan" validate:"required,min=2,max=50"`
	Angkatan   int     `json:"angkatan" validate:"required,min=1980,max=2030"`
	TahunLulus int     `json:"tahun_lulus" validate:"required,min=1980,max=2030"`
	Email      string  `json:"email" validate:"required,email,max=100"`
	NoTelepon  *string `json:"no_telepon" validate:"omitempty,max=15"`
	Alamat     *string `json:"alamat"`
}

type AlumniUpdate struct {
	Nama       string  `json:"nama" validate:"required,min=2,max=100"`
	Jurusan    string  `json:"jurusan" validate:"required,min=2,max=50"`
	Angkatan   int     `json:"angkatan" validate:"required,min=1980,max=2030"`
	TahunLulus int     `json:"tahun_lulus" validate:"required,min=1980,max=2030"`
	Email      string  `json:"email" validate:"required,email,max=100"`
	NoTelepon  *string `json:"no_telepon" validate:"omitempty,max=15"`
	Alamat     *string `json:"alamat"`
}

type AlumniPekerjaanSingkat struct {
	Nama           string `json:"nama"`
	Jurusan        string `json:"jurusan"`
	TahunLulus     int    `json:"tahun_lulus"`
	BidangIndustri string `json:"bidang_industri"`
}