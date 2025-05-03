package service

import (
	"errors"
	"ticket/models"
	"ticket/repository"
	"ticket/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(user *models.User) error
	Login(email, password string) (string, error)
	GetUserByID(id uint) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwtUtil  utils.JWTUtil
}

func NewUserService(userRepo repository.UserRepository, jwtUtil utils.JWTUtil) UserService {
	return &userService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
	}
}

func (s *userService) Register(user *models.User) error {
	// Memeriksa apakah pengguna sudah ada
	existingUser, err := s.userRepo.FindByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("pengguna dengan email ini sudah ada")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	// Menetapkan peran default "user" jika tidak ditentukan
	if user.Role == "" {
		user.Role = "user"
	}

	// Menyimpan pengguna
	return s.userRepo.Save(user)
}

func (s *userService) Login(email, password string) (string, error) {
	// Mencari pengguna berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	// Memverifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	// Membuat token JWT menggunakan jwtUtil
	token, err := s.jwtUtil.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
