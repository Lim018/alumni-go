package helper

import (
	// "alumni-go/config" // Dihapus untuk memutus import cycle
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Perubahan: Tambahkan parameter secretKey string
func GenerateToken(userID int, username, role, secretKey string) (string, error) {
	// secret := config.GetEnv("JWT_SECRET", "your-secret-key") // Baris ini dihapus
	expireHours := 24 // Default 24 hours

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Gunakan secretKey dari parameter
	return token.SignedString([]byte(secretKey))
}

// Perubahan: Tambahkan parameter secretKey string
func ValidateToken(tokenString, secretKey string) (*JWTClaims, error) {
	// secret := config.GetEnv("JWT_SECRET", "your-secret-key") // Baris ini dihapus

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Gunakan secretKey dari parameter
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}