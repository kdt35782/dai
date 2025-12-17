package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sm-medical/pkg/config"
	"time"
)

type Claims struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID int64, username, role string) (string, error) {
	cfg := config.AppConfig.JWT
	
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(cfg.ExpiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.AppConfig.JWT
	
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
