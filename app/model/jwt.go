package model

import (
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

// JWTClaims - JWT token claims
type JWTClaims struct {
    UserID   string `json:"user_id"` // Changed from int to string for ObjectID
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

// Helper to convert ObjectID to string for JWT
func NewJWTClaims(userID primitive.ObjectID, username, role string) JWTClaims {
    return JWTClaims{
        UserID:   userID.Hex(),
        Username: username,
        Role:     role,
    }
}