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

type AlumniService struct {
	repo *repository.AlumniRepository
}

func NewAlumniService() *AlumniService {
	return &AlumniService{
		repo: repository.NewAlumniRepository(),
	}
}

// func (s *AlumniService) GetAll(page, perPage int, search string) ([]model.Alumni, model.MetaData, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if perPage < 1 {
// 		perPage = 10
// 	}
// 	if perPage > 100 {
// 		perPage = 100
// 	}

// 	alumni, total, err := s.repo.GetAll(page, perPage, search)
// 	if err != nil {
// 		return nil, model.MetaData{}, err
// 	}

// 	meta := model.MetaData{
// 		Page:       page,
// 		PerPage:    perPage,
// 		Total:      total,
// 		TotalPages: (int(total) + perPage - 1) / perPage,
// 	}

// 	return alumni, meta, nil
// }

func (s *AlumniService) GetByID(id int) (*model.Alumni, error) {
	alumni, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if alumni == nil {
		return nil, errors.New("alumni not found")
	}
	return alumni, nil
}

func (s *AlumniService) Create(data *model.AlumniCreate) (*model.Alumni, error) {
	// Check if NIM already exists
	exists, err := s.repo.CheckNIMExists(data.NIM)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("NIM already exists")
	}

	// Check if email already exists
	exists, err = s.repo.CheckEmailExists(data.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Validate tahun lulus >= angkatan + 3 (minimal 3 tahun kuliah)
	if data.TahunLulus < data.Angkatan+3 {
		return nil, errors.New("tahun lulus harus minimal 3 tahun setelah angkatan")
	}

	return s.repo.Create(data)
}

func (s *AlumniService) Update(id int, data *model.AlumniUpdate) (*model.Alumni, error) {
	// Check if alumni exists
	_, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email already exists (exclude current record)
	exists, err := s.repo.CheckEmailExists(data.Email, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Validate tahun lulus >= angkatan + 3
	if data.TahunLulus < data.Angkatan+3 {
		return nil, errors.New("tahun lulus harus minimal 3 tahun setelah angkatan")
	}

	alumni, err := s.repo.Update(id, data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("alumni not found")
		}
		return nil, err
	}

	return alumni, nil
}

func (s *AlumniService) Delete(id int) error {
	// Check if alumni exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	err = s.repo.Delete(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("alumni not found")
		}
		return err
	}

	return nil
}

func (s *AlumniService) GetAlumniBaruBekerja() ([]model.AlumniPekerjaanSingkat, error) {
	return s.repo.GetAlumniBaruBekerja()
}

func (s *AlumniService) GetAll(c *fiber.Ctx) (*model.PaginatedResponse, error) {
	// 1. Ekstrak parameter dari query URL [cite: 200]
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))
	search := c.Query("search", "")
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")

	// 2. Validasi input untuk keamanan (mencegah SQL Injection) [cite: 158]
	sortByWhitelist := map[string]bool{"id": true, "nama": true, "nim": true, "angkatan": true, "tahun_lulus": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id" // default sort
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}

	// 3. Panggil kedua method repository
	alumni, err := s.repo.GetAll(page, perPage, search, sortBy, order)
	if err != nil {
		return nil, err
	}
	
	total, err := s.repo.CountAll(search)
	if err != nil {
		return nil, err
	}

	// 4. Buat response dengan data dan metadata [cite: 178]
	totalPages := (total + int64(perPage) - 1) / int64(perPage)
	meta := model.MetaData{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: int(totalPages),
	}
	
	response := &model.PaginatedResponse{
		Success: true,
		Message: "Alumni retrieved successfully",
		Data: alumni,
		Meta: meta,
	}

	return response, nil
}