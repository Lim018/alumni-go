package model

import (
    "time"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// User - Base model for MongoDB
type User struct {
    ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Username     string             `json:"username" bson:"username"`
    Email        string             `json:"email" bson:"email"`
    PasswordHash string             `json:"-" bson:"password_hash"`
    Role         string             `json:"role" bson:"role"`
    CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
}

// LoginRequest - Request for POST /login
type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// UserResponse - Response for user data (without sensitive info)
type UserResponse struct {
    ID       string `json:"id"`
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