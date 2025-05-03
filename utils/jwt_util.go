// utils/jwt_util.go
package utils

import (
	"errors"
	"fmt"
	"ticket/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTUtil interface mendefinisikan fungsi-fungsi untuk mengelola token JWT
type JWTUtil interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIDFromToken(token *jwt.Token) (uint, error)
	GetUserRoleFromToken(token *jwt.Token) (string, error)
}

type jwtCustomClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type jwtUtil struct {
	secretKey  string
	expiration time.Duration
}

// NewJWTUtil menciptakan instance baru dari JWTUtil
func NewJWTUtil(secretKey string, expiration time.Duration) JWTUtil {
	return &jwtUtil{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

func (u *jwtUtil) GenerateToken(user *models.User) (string, error) {
	claims := &jwtCustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(u.secretKey))
}

func (u *jwtUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("metode penandatanganan tidak sesuai: %v", token.Header["alg"])
		}
		return []byte(u.secretKey), nil
	})
}

func (u *jwtUtil) GetUserIDFromToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("token tidak valid")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user id dalam token tidak valid")
	}

	return uint(userID), nil
}

func (u *jwtUtil) GetUserRoleFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("token tidak valid")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("role dalam token tidak valid")
	}

	return role, nil
}
