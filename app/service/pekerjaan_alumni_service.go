package service

import (
	"alumni-go/app/model"
	"alumni-go/app/repository"
	"database/sql"
	"errors"
)

type PekerjaanAlumniService struct {
	repo      *repository.PekerjaanAlumniRepository
	alumniRepo *repository.AlumniRepository
}

func NewPekerjaanAlumniService() *PekerjaanAlumniService {
	return &PekerjaanAlumniService{
		repo:      repository.NewPekerjaanAlumniRepository(),
		alumniRepo: repository.NewAlumniRepository(),
	}
}

func (s *PekerjaanAlumniService) GetAll(page, perPage int, search string) ([]model.PekerjaanAlumniWithAlumni, model.MetaData, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	pekerjaan, total, err := s.repo.GetAll(page, perPage, search)
	if err != nil {
		return nil, model.MetaData{}, err
	}

	meta := model.MetaData{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: (int(total) + perPage - 1) / perPage,
	}

	return pekerjaan, meta, nil
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

	// PERBAIKAN: Panggil fungsi Create dari repository
	pekerjaan, err := s.repo.Create(data)
	if err != nil {
		// Jika ada error saat create, langsung kembalikan errornya
		return nil, err
	}

	return pekerjaan, nil
}

func (s *PekerjaanAlumniService) Update(id int, data *model.PekerjaanUpdate) (*model.PekerjaanAlumni, error) {
	// 1. Periksa dulu apakah data pekerjaan yang akan diupdate memang ada
	_, err := s.repo.GetByID(id)
	if err != nil {
		// repo.GetByID sudah mengembalikan error jika tidak ditemukan, 
		// tapi kita bisa bungkus dengan pesan yang lebih spesifik jika mau.
		if err == sql.ErrNoRows {
			return nil, errors.New("pekerjaan not found")
		}
		return nil, err
	}

	// 2. Validasi tanggal
	if data.TanggalSelesaiKerja != nil && data.TanggalMulaiKerja.After(*data.TanggalSelesaiKerja) {
		return nil, errors.New("tanggal mulai kerja tidak boleh setelah tanggal selesai kerja")
	}

	// 3. Panggil repository untuk melakukan update ke database
	pekerjaan, err := s.repo.Update(id, data)
	if err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

func (s *PekerjaanAlumniService) Delete(id int) error {
	// Check if pekerjaan exists
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