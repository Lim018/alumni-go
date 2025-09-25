package repository

import (
	"alumni-go/app/model"
	"database/sql"
	"fmt"
	"strings"
)

type AlumniRepository struct {
	db *sql.DB
}

func NewAlumniRepository(db *sql.DB) *AlumniRepository {
	return &AlumniRepository{db: db}
}

func (r *AlumniRepository) GetByID(id int) (*model.Alumni, error) {
	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni WHERE id = $1
	`
	
	var a model.Alumni
	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
		&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &a, nil
}

func (r *AlumniRepository) Create(alumni *model.AlumniCreate) (*model.Alumni, error) {
	query := `
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
	`
	
	var a model.Alumni
	err := r.db.QueryRow(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
		alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat).Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
		&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return &a, nil
}

func (r *AlumniRepository) Update(id int, alumni *model.AlumniUpdate) (*model.Alumni, error) {
	query := `
		UPDATE alumni 
		SET nama = $2, jurusan = $3, angkatan = $4, tahun_lulus = $5, email = $6, no_telepon = $7, alamat = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
	`
	
	var a model.Alumni
	err := r.db.QueryRow(query, id, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
		alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat).Scan(
		&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
		&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &a, nil
}

func (r *AlumniRepository) Delete(id int) error {
	query := "DELETE FROM alumni WHERE id = $1"
	result, err := r.db.Exec(query, id)
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

func (r *AlumniRepository) CheckNIMExists(nim string, excludeID ...int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM alumni WHERE nim = $1"
	args := []interface{}{nim}
	
	if len(excludeID) > 0 {
		query += " AND id != $2"
		args = append(args, excludeID[0])
	}
	
	query += ")"
	
	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}

func (r *AlumniRepository) CheckEmailExists(email string, excludeID ...int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM alumni WHERE email = $1"
	args := []interface{}{email}
	
	if len(excludeID) > 0 {
		query += " AND id != $2"
		args = append(args, excludeID[0])
	}
	
	query += ")"
	
	var exists bool
	err := r.db.QueryRow(query, args...).Scan(&exists)
	return exists, err
}

func (r *AlumniRepository) GetAlumniBaruBekerja() ([]model.AlumniPekerjaanSingkat, error) {
	query := `
		SELECT
			a.nama,
			a.jurusan,
			a.tahun_lulus,
			p.bidang_industri
		FROM
			alumni a
		JOIN
			pekerjaan_alumni p ON a.id = p.alumni_id
		WHERE
			p.status_pekerjaan = 'aktif' AND AGE(NOW(), p.tanggal_mulai_kerja) < '3 years'
		ORDER BY
			a.nama ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.AlumniPekerjaanSingkat
	for rows.Next() {
		var res model.AlumniPekerjaanSingkat
		err := rows.Scan(&res.Nama, &res.Jurusan, &res.TahunLulus, &res.BidangIndustri)
		if err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, nil
}

func (r *AlumniRepository) GetAll(page, perPage int, search, sortBy, order string) ([]model.Alumni, error) {
	offset := (page - 1) * perPage
	
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE LOWER(nama) LIKE $1 OR LOWER(nim) LIKE $1 OR LOWER(jurusan) LIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)
	
	rows, err := r.db.Query(query, "%"+strings.ToLower(search)+"%", perPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}

	return alumni, nil
}

func (r *AlumniRepository) CountAll(search string) (int64, error) {
	var total int64
	countQuery := `SELECT COUNT(*) FROM alumni WHERE LOWER(nama) LIKE $1 OR LOWER(nim) LIKE $1 OR LOWER(jurusan) LIKE $1`
	
	err := r.db.QueryRow(countQuery, "%"+strings.ToLower(search)+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	
	return total, nil
}