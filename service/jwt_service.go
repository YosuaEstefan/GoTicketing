// service/jwt_service.go
package service

import (
	"errors"
	"fmt"
	"ticket/models"
	"time"

	"github.com/golang-jwt/jwt/v5"

)

type JWTService interface {
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

type jwtService struct {
	secretKey  string
	expiration time.Duration
}

func NewJWTService(secretKey string, expiration time.Duration) JWTService {
	return &jwtService{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

func (s *jwtService) GenerateToken(user *models.User) (string, error) {
	claims := &jwtCustomClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.secretKey))
}

func (s *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
}

func (s *jwtService) GetUserIDFromToken(token *jwt.Token) (uint, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("invalid user id in token")
	}

	return uint(userID), nil
}

func (s *jwtService) GetUserRoleFromToken(token *jwt.Token) (string, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return "", errors.New("invalid role in token")
	}

	return role, nil
}
