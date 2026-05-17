package services

import (
	"context"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, u models.User) (int, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserRepository(repo repositories.UserRepository) *userService {
	return &userService{repo: repo}
}
