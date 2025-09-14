package service

import (
	"alumni-management-system/app/model"
	"alumni-management-system/app/repository"
	"alumni-management-system/helper"
	"errors"
)

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
	return &AuthService{
		userRepo: repository.NewUserRepository(),
	}
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// Basic validation
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Get user by username
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("failed to authenticate user")
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if !helper.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT token
	token, err := helper.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Remove password from response
	user.Password = ""

	return &model.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) GetUserByID(id int) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user ID")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("failed to get user data")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Remove password from response
	user.Password = ""
	
	return user, nil
}

func (s *AuthService) ValidateUserRole(userID int, allowedRoles []string) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	// Check if user role is in allowed roles
	for _, role := range allowedRoles {
		if user.Role == role {
			return true, nil
		}
	}

	return false, nil
}