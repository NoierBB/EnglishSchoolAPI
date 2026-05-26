// services/user_service.go
package services

import (
	"context"
	"database/sql"
	"errors"
	"log" // добавь для отладки
	"time"

	"github.com/NoierBB/englishSchool/internal/dto"
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
	RegisterStudent(ctx context.Context, req dto.RegisterStudentRequest) (int, int, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type AuthService struct {
	repo        *repositories.UserRepository
	studentRepo *repositories.StudentRepository
	db          *sql.DB
	jwtSecret   string
}

func NewAuthService(repo *repositories.UserRepository, studentRepo *repositories.StudentRepository, db *sql.DB, jwtSecret string) *AuthService {
	return &AuthService{
		repo:        repo,
		studentRepo: studentRepo,
		db:          db,
		jwtSecret:   jwtSecret,
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

	return s.repo.CreateUser(ctx, u)
}

func (s *AuthService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.repo.GetUsers(ctx)
}

func (s *AuthService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetUserById(ctx, id)
}

func (s *AuthService) UpdateUser(ctx context.Context, u models.User) error {
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

// func (s *AuthService) RegisterStudent(ctx context.Context, req dto.RegisterStudentRequest) (int, int, error) {
// 	tx, err := s.db.BeginTx(ctx, nil)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	defer tx.Rollback()

// 	exist, err := s.repo.ExistByEmail(ctx, req.Email)
// 	if err != nil {
// 		return 0, 0, err
// 	}
// 	if exist {
// 		return 0, 0, errors.New("user already exists")
// 	}

// 	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	userId, err := s.repo.CreateUserTx(ctx, tx, models.User{
// 		Email:    req.Email,
// 		Password: string(hash),
// 		Role:     "student",
// 	})
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	studentId, err := s.studentRepo.CreateStudentTx(ctx, tx, models.Students{
// 		UserId: userId,
// 		Name:   req.Name,
// 		Age:    req.Age,
// 		Level:  req.Level,
// 	})
// 	if err != nil {
// 		return 0, 0, err
// 	}

// 	if err := tx.Commit(); err != nil {
// 		return 0, 0, err
// 	}

// 	return userId, studentId, nil
// }

func (s *AuthService) RegisterStudent(ctx context.Context, req dto.RegisterStudentRequest) (int, int, error) {
	log.Println("=== RegisterStudent called ===")
	log.Printf("Request: Email=%s, Name=%s, Age=%d, Level=%s", req.Email, req.Name, req.Age, req.Level)

	// Проверяем, что s не nil
	if s == nil {
		log.Println("ERROR: AuthService is nil")
		return 0, 0, errors.New("auth service is nil")
	}

	// Проверяем, что s.db не nil
	if s.db == nil {
		log.Println("ERROR: db is nil in AuthService")
		return 0, 0, errors.New("database connection is nil")
	}

	// Проверяем, что s.repo не nil
	if s.repo == nil {
		log.Println("ERROR: userRepo is nil in AuthService")
		return 0, 0, errors.New("user repository is nil")
	}

	// Проверяем, что s.studentRepo не nil
	if s.studentRepo == nil {
		log.Println("ERROR: studentRepo is nil in AuthService")
		return 0, 0, errors.New("student repository is nil")
	}

	log.Println("All dependencies are OK, starting transaction")

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return 0, 0, err
	}
	defer tx.Rollback()

	log.Println("Transaction started, checking if email exists")

	exist, err := s.repo.ExistByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("ExistByEmail error: %v", err)
		return 0, 0, err
	}
	if exist {
		log.Println("User already exists")
		return 0, 0, errors.New("user already exists")
	}

	log.Println("Email is unique, hashing password")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		return 0, 0, err
	}

	log.Println("Creating user")

	userId, err := s.repo.CreateUserTx(ctx, tx, models.User{
		Email:    req.Email,
		Password: string(hash),
		Role:     "student",
	})
	if err != nil {
		log.Printf("CreateUserTx error: %v", err)
		return 0, 0, err
	}

	log.Printf("User created with ID: %d, creating student profile", userId)

	studentId, err := s.studentRepo.CreateStudentTx(ctx, tx, models.Students{
		UserId: userId,
		Name:   req.Name,
		Age:    req.Age,
		Level:  req.Level,
	})
	if err != nil {
		log.Printf("CreateStudentTx error: %v", err)
		return 0, 0, err
	}

	log.Println("Committing transaction")

	if err := tx.Commit(); err != nil {
		log.Printf("Commit error: %v", err)
		return 0, 0, err
	}

	log.Printf("Successfully registered student: userID=%d, studentID=%d", userId, studentId)
	return userId, studentId, nil
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
