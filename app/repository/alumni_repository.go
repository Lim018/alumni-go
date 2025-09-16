package repository

import (
	"alumni-go/app/model"
	"alumni-go/database"
	"database/sql"
	"fmt"
)

type AlumniRepository struct{}

func NewAlumniRepository() *AlumniRepository {
	return &AlumniRepository{}
}

func (r *AlumniRepository) GetAll(page, perPage int, search string) ([]model.Alumni, int64, error) {
	offset := (page - 1) * perPage
	
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	if search != "" {
		whereClause = "WHERE LOWER(nama) LIKE LOWER($" + fmt.Sprintf("%d", argIndex) + ") OR LOWER(nim) LIKE LOWER($" + fmt.Sprintf("%d", argIndex+1) + ")"
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIndex += 2
	}

	// Count total records
	countQuery := "SELECT COUNT(*) FROM alumni " + whereClause
	var total int64
	err := database.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	
	args = append(args, perPage, offset)
	
	rows, err := database.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		alumni = append(alumni, a)
	}

	return alumni, total, nil
}

func (r *AlumniRepository) GetByID(id int) (*model.Alumni, error) {
	query := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni WHERE id = $1
	`
	
	var a model.Alumni
	err := database.DB.QueryRow(query, id).Scan(&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
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
	err := database.DB.QueryRow(query, alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
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
	err := database.DB.QueryRow(query, id, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
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

func (r *AlumniRepository) CheckNIMExists(nim string, excludeID ...int) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM alumni WHERE nim = $1"
	args := []interface{}{nim}
	
	if len(excludeID) > 0 {
		query += " AND id != $2"
		args = append(args, excludeID[0])
	}
	
	query += ")"
	
	var exists bool
	err := database.DB.QueryRow(query, args...).Scan(&exists)
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
	err := database.DB.QueryRow(query, args...).Scan(&exists)
	return exists, err
}