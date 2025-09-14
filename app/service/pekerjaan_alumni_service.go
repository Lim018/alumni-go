package service

import (
	"alumni-management-system/app/model"
	"alumni-management-system/app/repository"
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

	pekerjaan, err := s.repo.Update(id, data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("pekerjaan not found")
		}
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