package model

import "github.com/golang-jwt/jwt/v5"

// JWTClaims - JWT token claims
type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}