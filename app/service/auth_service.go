package service

import (
	"database/sql"
	"errors"
	"go-fiber/app/model"
	"go-fiber/app/repository"
	"go-fiber/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func LoginService(db *sql.DB, req model.LoginRequest) (*model.LoginResponse, error) {
	user, passwordHash, err := repository.FindUserByUsernameOrEmail(db, req.Username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		return nil, errors.New("username atau password salah")
	}

	token, err := utils.GenerateToken(*user)
	if err != nil {
		return nil, errors.New("gagal generate token")
	}

	return &model.LoginResponse{
		User:  user.ToUserResponse(),
		Token: token,
	}, nil
}

func GetUsersService(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		limit, _ := strconv.Atoi(c.Query("limit", "10"))
		sortBy := c.Query("sortBy", "id")
		order := c.Query("order", "asc")
		search := c.Query("search", "")
		offset := (page - 1) * limit

		sortByWhitelist := map[string]string{
			"id":         "id",
			"username":   "username",
			"email":      "email",
			"role":       "role",
			"created_at": "created_at",
		}
		col, ok := sortByWhitelist[sortBy]
		if !ok {
			col = "id"
		}

		ord := "ASC"
		if strings.ToLower(order) == "desc" {
			ord = "DESC"
		}

		users, err := repository.GetUsersRepo(db, search, col, ord, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
		}

		total, err := repository.CountUsersRepo(db, search)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
		}

		pages := 0
		if total > 0 {
			pages = (total + limit - 1) / limit
		}

		// Convert to response models
		responses := make([]model.UserResponse, len(users))
		for i, user := range users {
			responses[i] = user.ToUserResponse()
		}

		response := model.UserListResponse{
			Data: responses,
			Meta: model.MetaInfo{
				Page:   page,
				Limit:  limit,
				Total:  total,
				Pages:  pages,
				SortBy: col,
				Order:  ord,
				Search: search,
			},
		}
		return c.JSON(response)
	}
}