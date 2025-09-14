package repository

import (
	"alumni-management-system/app/model"
	"alumni-management-system/database"
	"database/sql"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users WHERE username = $1
	`
	
	var u model.User
	err := database.DB.QueryRow(query, username).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &u, nil
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users WHERE id = $1
	`
	
	var u model.User
	err := database.DB.QueryRow(query, id).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &u, nil
}