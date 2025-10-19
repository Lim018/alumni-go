package model

import "time"

// type User struct {
//     ID        int       `json:"id"`
//     Username  string    `json:"username"`
//     Email     string    `json:"email"`
//     Role      string    `json:"role"`
//     CreatedAt time.Time `json:"created_at"`
// }

// type LoginRequest struct {
//     Username string `json:"username"`
//     Password string `json:"password"`
// }

// type LoginResponse struct {
//     User  User   `json:"user"`
//     Token string `json:"token"`
// }

// type JWTClaims struct {
//     UserID   int    `json:"user_id"`
//     Username string `json:"username"`
//     Role     string `json:"role"`
//     jwt.RegisteredClaims
// }

// User - Base model for database representation
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest - Request for POST /login
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse - Response for user data (without sensitive info)
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// LoginResponse - Response for POST /login
type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// UserListResponse - Response for GET /users
type UserListResponse struct {
	Data []UserResponse `json:"data"`
	Meta MetaInfo       `json:"meta"`
}