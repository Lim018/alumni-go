package repository

import (
	"alumni-management-system/app/model"
	"alumni-management-system/database"
	"database/sql"
	"fmt"
)

type PekerjaanAlumniRepository struct{}

func NewPekerjaanAlumniRepository() *PekerjaanAlumniRepository {
	return &PekerjaanAlumniRepository{}
}

func (r *PekerjaanAlumniRepository) GetAll(page, perPage int, search string) ([]model.PekerjaanAlumniWithAlumni, int64, error) {
	offset := (page - 1) * perPage
	
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = "WHERE LOWER(p.nama_perusahaan) LIKE LOWER($" + fmt.Sprintf("%d", argIndex) + ") OR LOWER(a.nama) LIKE LOWER($" + fmt.Sprintf("%d", argIndex+1) + ")"
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	// Count total records
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni p 
				   JOIN alumni a ON p.alumni_id = a.id ` + whereClause
	var total int64
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri, p.lokasi_kerja,
			   p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, p.status_pekerjaan, p.deskripsi_pekerjaan,
			   p.created_at, p.updated_at,
			   a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.no_telepon, a.alamat,
			   a.created_at, a.updated_at
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	
	args = append(args, perPage, offset)
	
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumniWithAlumni
	for rows.Next() {
		var p model.PekerjaanAlumniWithAlumni
		err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
			&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
			&p.CreatedAt, &p.UpdatedAt,
			&p.Alumni.ID, &p.Alumni.NIM, &p.Alumni.Nama, &p.Alumni.Jurusan, &p.Alumni.Angkatan, &p.Alumni.TahunLulus,
			&p.Alumni.Email, &p.Alumni.NoTelepon, &p.Alumni.Alamat, &p.Alumni.CreatedAt, &p.Alumni.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		pekerjaan = append(pekerjaan, p)
	}

	return pekerjaan, total, nil
}

func (r *PekerjaanAlumniRepository) GetByID(id int) (*model.PekerjaanAlumniWithAlumni, error) {
	query := `
		SELECT p.id, p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri, p.lokasi_kerja,
			   p.gaji_range, p.tanggal_mulai_kerja, p.tanggal_selesai_kerja, p.status_pekerjaan, p.deskripsi_pekerjaan,
			   p.created_at, p.updated_at,
			   a.id, a.nim, a.nama, a.jurusan, a.angkatan, a.tahun_lulus, a.email, a.no_telepon, a.alamat,
			   a.created_at, a.updated_at
		FROM pekerjaan_alumni p
		JOIN alumni a ON p.alumni_id = a.id
		WHERE p.id = $1
	`
	
	var p model.PekerjaanAlumniWithAlumni
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
		&p.CreatedAt, &p.UpdatedAt,
		&p.Alumni.ID, &p.Alumni.NIM, &p.Alumni.Nama, &p.Alumni.Jurusan, &p.Alumni.Angkatan, &p.Alumni.TahunLulus,
		&p.Alumni.Email, &p.Alumni.NoTelepon, &p.Alumni.Alamat, &p.Alumni.CreatedAt, &p.Alumni.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &p, nil
}

func (r *PekerjaanAlumniRepository) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
			   gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
			   created_at, updated_at
		FROM pekerjaan_alumni
		WHERE alumni_id = $1
		ORDER BY tanggal_mulai_kerja DESC
	`
	
	rows, err := database.DB.Query(query, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaan []model.PekerjaanAlumni
	for rows.Next() {
		var p model.PekerjaanAlumni
		err := rows.Scan(&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri,
			&p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja,
			&p.StatusPekerjaan, &p.DeskripsiPekerjaan, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pekerjaan = append(pekerjaan, p)
	}

	return pekerjaan, nil
}

func (r *PekerjaanAlumniRepository) Create(pekerjaan *model.PekerjaanCreate) (*model.PekerjaanAlumni, error) {
	query := `
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
									  gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
				  gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
				  created_at, updated_at
	`
	
	var p model.PekerjaanAlumni
	err := database.DB.QueryRow(query, pekerjaan.AlumniID, pekerjaan.NamaPerusahaan, pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri, pekerjaan.LokasiKerja, pekerjaan.GajiRange, pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja, pekerjaan.StatusPekerjaan, pekerjaan.DeskripsiPekerjaan).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
		&p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &p, nil
}

func (r *PekerjaanAlumniRepository) Update(id int, pekerjaan *model.PekerjaanUpdate) (*model.PekerjaanAlumni, error) {
	query := `
		UPDATE pekerjaan_alumni 
		SET nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, lokasi_kerja = $5,
			gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, status_pekerjaan = $9,
			deskripsi_pekerjaan = $10, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
				  gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan,
				  created_at, updated_at
	`
	
	var p model.PekerjaanAlumni
	err := database.DB.QueryRow(query, id, pekerjaan.NamaPerusahaan, pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri, pekerjaan.LokasiKerja, pekerjaan.GajiRange, pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja, pekerjaan.StatusPekerjaan, pekerjaan.DeskripsiPekerjaan).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan, &p.BidangIndustri, &p.LokasiKerja,
		&p.GajiRange, &p.TanggalMulaiKerja, &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.DeskripsiPekerjaan,
		&p.CreatedAt, &p.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &p, nil
}

func (r *PekerjaanAlumniRepository) Delete(id int) error {
	query := "DELETE FROM pekerjaan_alumni WHERE id = $1"
	result, err := database.DB.Exec(query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}