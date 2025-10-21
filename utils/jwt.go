package utils

import (
    "os"
    "time"
    
    "go-fiber/app/model"
    
    "github.com/golang-jwt/jwt/v5"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

func GenerateToken(user model.User) (string, error) {
    claims := model.JWTClaims{
        UserID:   user.ID.Hex(), // Convert ObjectID to string
        Username: user.Username,
        Role:     user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := os.Getenv("JWT_SECRET")
    
    return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string) (*model.JWTClaims, error) {
    secret := os.Getenv("JWT_SECRET")
    
    token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, jwt.ErrInvalidKey
}

// Helper function to get ObjectID from JWT claims
func GetUserIDFromClaims(claims *model.JWTClaims) (primitive.ObjectID, error) {
    return primitive.ObjectIDFromHex(claims.UserID)
}