package service

import (
	"alumni-go/app/model"
	"alumni-go/app/repository"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanAlumniService struct {
	repo       *repository.PekerjaanAlumniRepository
	alumniRepo *repository.AlumniRepository
}

// Perubahan di sini: NewPekerjaanAlumniService sekarang menerima *sql.DB
func NewPekerjaanAlumniService(db *sql.DB) *PekerjaanAlumniService {
	return &PekerjaanAlumniService{
		// Teruskan 'db' saat membuat kedua repository
		repo:       repository.NewPekerjaanAlumniRepository(db),
		alumniRepo: repository.NewAlumniRepository(db),
	}
}

func (s *PekerjaanAlumniService) GetByID(id int) (*model.PekerjaanAlumniWithAlumni, error) {
	pekerjaan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if pekerjaan == nil {
		return nil, errors.New("pekerjaan not found")
	}
	return pekerjaan, nil
}

func (s *PekerjaanAlumniService) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	// Check if alumni exists
	alumni, err := s.alumniRepo.GetByID(alumniID)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni not found")
	}

	return s.repo.GetByAlumniID(alumniID)
}

func (s *PekerjaanAlumniService) Create(data *model.PekerjaanCreate) (*model.PekerjaanAlumni, error) {
	// Check if alumni exists
	alumni, err := s.alumniRepo.GetByID(data.AlumniID)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni not found")
	}

	// Validate tanggal
	if data.TanggalSelesaiKerja != nil && data.TanggalMulaiKerja.After(*data.TanggalSelesaiKerja) {
		return nil, errors.New("tanggal mulai kerja tidak boleh setelah tanggal selesai kerja")
	}

	pekerjaan, err := s.repo.Create(data)
	if err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (s *PekerjaanAlumniService) Update(id int, data *model.PekerjaanUpdate) (*model.PekerjaanAlumni, error) {
	_, err := s.repo.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pekerjaan not found")
		}
		return nil, err
	}

	if data.TanggalSelesaiKerja != nil && data.TanggalMulaiKerja.After(*data.TanggalSelesaiKerja) {
		return nil, errors.New("tanggal mulai kerja tidak boleh setelah tanggal selesai kerja")
	}

	pekerjaan, err := s.repo.Update(id, data)
	if err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (s *PekerjaanAlumniService) Delete(id int) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("pekerjaan not found")
		}
		return err
	}

	return nil
}

func (s *PekerjaanAlumniService) GetAll(c *fiber.Ctx) (*model.PaginatedResponse, error) {
	// 1. Ekstrak parameter
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")

	// 2. Validasi
	sortByWhitelist := map[string]bool{"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "tanggal_mulai_kerja": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// 3. Panggil Repository
	pekerjaan, err := s.repo.GetAll(page, perPage, search, sortBy, order)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.CountAll(search)
	if err != nil {
		return nil, err
	}

	// 4. Buat Response
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	meta := model.MetaData{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int(totalPages),
	}
	
	response := &model.PaginatedResponse{
		Success: true,
		Message: "Pekerjaan retrieved successfully",
		Data: pekerjaan,
		Meta: meta,
	}

	return response, nil
}