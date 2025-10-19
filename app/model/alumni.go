package model

import "time"

// type Alumni struct {
// 	ID         int       `json:"id"`
// 	NIM        string    `json:"nim"`
// 	Nama       string    `json:"nama"`
// 	Jurusan    string    `json:"jurusan"`
// 	Angkatan   int       `json:"angkatan"`
// 	TahunLulus int       `json:"tahun_lulus"`
// 	Email      string    `json:"email"`
// 	NoTelepon  string    `json:"no_telepon"`
// 	Alamat     string    `json:"alamat"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// 	UserID     int       `json:"user_id"`
// }

// Alumni - Base model for database representation
type Alumni struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  string    `json:"no_telepon"`
	Alamat     string    `json:"alamat"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	UserID     int       `json:"user_id"`
}

// CreateAlumniRequest - Request for POST /alumni
type CreateAlumniRequest struct {
	NIM        string `json:"nim" validate:"required"`
	Nama       string `json:"nama" validate:"required"`
	Jurusan    string `json:"jurusan" validate:"required"`
	Angkatan   int    `json:"angkatan" validate:"required"`
	TahunLulus int    `json:"tahun_lulus" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	NoTelepon  string `json:"no_telepon" validate:"required"`
	Alamat     string `json:"alamat"`
}

// UpdateAlumniRequest - Request for PUT /alumni/:id
type UpdateAlumniRequest struct {
	NIM        string `json:"nim" validate:"required"`
	Nama       string `json:"nama" validate:"required"`
	Jurusan    string `json:"jurusan" validate:"required"`
	Angkatan   int    `json:"angkatan" validate:"required"`
	TahunLulus int    `json:"tahun_lulus" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	NoTelepon  string `json:"no_telepon" validate:"required"`
	Alamat     string `json:"alamat"`
}

// AlumniResponse - Response for single alumni
type AlumniResponse struct {
	ID         int       `json:"id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  string    `json:"no_telepon"`
	Alamat     string    `json:"alamat"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// AlumniListResponse - Response for GET /alumni (datatable)
type AlumniListResponse struct {
	Data []AlumniResponse `json:"data"`
	Meta MetaInfo         `json:"meta"`
}

// AlumniStatsByJurusanResponse - Response for GET /alumni/stats/jurusan
type AlumniStatsByJurusanResponse struct {
	Jurusan string `json:"jurusan"`
	Total   int    `json:"total"`
}