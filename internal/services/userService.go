// services/user_service.go
package services

import (
	"context"
	"errors"
	"log" // добавь для отладки
	"time"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, u models.User) (int, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, id int) error
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthService struct {
	repo      *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repositories.UserRepository, secret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: secret,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, u models.User) (int, error) {
	// Хешируем пароль
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return 0, err
		}
		u.Password = string(hash)
	}

	// Устанавливаем роль по умолчанию
	if u.Role == "" {
		u.Role = "student"
	}

	return s.repo.CreateUser(ctx, u)
}

func (s *AuthService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *AuthService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *AuthService) UpdateUser(ctx context.Context, u models.User) error {
	// Если пароль передан, хешируем
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash)
	}
	return s.repo.UpdateUser(ctx, u)
}

func (s *AuthService) DeleteUser(ctx context.Context, id int) error {
	return s.repo.DeleteUser(ctx, id)
}

func (s *AuthService) Register(ctx context.Context, email, password string) (int, error) {
	log.Printf("Register called with email: %s", email)

	if email == "" || password == "" {
		log.Printf("Empty credentials")
		return 0, errors.New("empty credentials")
	}

	exists, err := s.repo.ExistByEmail(ctx, email)
	if err != nil {
		log.Printf("ExistByEmail error: %v", err)
		return 0, err
	}
	if exists {
		log.Printf("User already exists: %s", email)
		return 0, errors.New("user already exist")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Hash error: %v", err)
		return 0, err
	}

	user := models.User{
		Email:    email,
		Password: string(hash),
		Role:     "student", // роль по умолчанию
	}

	log.Printf("Creating user...")
	id, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Printf("CreateUser error: %v", err)
		return 0, err
	}

	log.Printf("User created with id=%d", id)
	return id, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	log.Printf("Login called for email: %s", email)

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Password comparison error: %v", err)
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Id,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Printf("Token signing error: %v", err)
		return "", err
	}

	log.Printf("Login successful for user: %s", email)
	return tokenString, nil
}
