package repository

import (
    "database/sql"
    "go-fiber/app/model"

    "fmt"
	"log"
)

func CreateAlumni(db *sql.DB, alumni model.Alumni) (*model.Alumni, error) {
    query := `INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat) 
              VALUES ($1,$2,$3,$4,$5,$6,$7,$8) 
              RETURNING id, created_at, updated_at`
    err := db.QueryRow(query,
        alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
        alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat).
        Scan(&alumni.ID, &alumni.CreatedAt, &alumni.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &alumni, nil
}

func UpdateAlumni(db *sql.DB, id int, alumni model.Alumni) (*model.Alumni, error) {
    query := `UPDATE alumni 
              SET nim=$1, nama=$2, jurusan=$3, angkatan=$4, tahun_lulus=$5, email=$6, no_telepon=$7, alamat=$8, updated_at=NOW()
              WHERE id=$9 RETURNING id, created_at, updated_at`
    err := db.QueryRow(query,
        alumni.NIM, alumni.Nama, alumni.Jurusan, alumni.Angkatan,
        alumni.TahunLulus, alumni.Email, alumni.NoTelepon, alumni.Alamat, id).
        Scan(&alumni.ID, &alumni.CreatedAt, &alumni.UpdatedAt)
    if err != nil {
        return nil, err
    }
    return &alumni, nil
}

func DeleteAlumni(db *sql.DB, id int) error {
    _, err := db.Exec(`DELETE FROM alumni WHERE id=$1`, id)
    return err
}

//datatable
func GetAlumniRepo(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	// default columns to prevent SQL injection
	allowedSort := map[string]bool{"id": true, "nama": true, "angkatan": true, "tahun_lulus": true}
	if !allowedSort[sortBy] {
		sortBy = "id"
	}
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nama ILIKE $1 OR nim ILIKE $1 OR email ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	log.Println("SQL:", query)
	log.Println("Params:", "%"+search+"%", limit, offset)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var a model.Alumni
		if err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan,
			&a.TahunLulus, &a.Email, &a.NoTelepon, &a.Alamat,
			&a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

func CountAlumniRepo(db *sql.DB, search string) (int, error) {
	var total int
	query := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR email ILIKE $1`
	err := db.QueryRow(query, "%"+search+"%").Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}

func GetAlumniStatsByJurusan(db *sql.DB) ([]map[string]interface{}, error) {
    query := `
        SELECT jurusan, COUNT(*) as total 
        FROM alumni 
        GROUP BY jurusan 
        ORDER BY total DESC
    `
    
    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var stats []map[string]interface{}
    for rows.Next() {
        var jurusan string
        var total int
        if err := rows.Scan(&jurusan, &total); err != nil {
            return nil, err
        }
        stats = append(stats, map[string]interface{}{
            "jurusan": jurusan,
            "total":   total,
        })
    }
    return stats, nil
}