package service

import (
	"alumni-go/app/model"
	"alumni-go/app/repository"
	"alumni-go/helper"
	"database/sql"
	"errors"
)

// Tambahkan field untuk menyimpan secret key
type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

// Ubah constructor untuk menerima db dan jwtSecret
func NewAuthService(db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  repository.NewUserRepository(db),
		jwtSecret: jwtSecret, // Simpan secret key
	}
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Validasi dasar
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	// Ambil user dari repository
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		// Log error internal di sini jika perlu
		return nil, errors.New("invalid username or password")
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Verifikasi password
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Gunakan secret key yang sudah disimpan untuk membuat token
	token, err := helper.GenerateToken(user.ID, user.Username, user.Role, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Hapus password dari respons
	user.Password = ""

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("failed to get user data")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	user.Password = ""
	return user, nil
}
