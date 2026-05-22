package services

import (
	"context"
	"errors"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/repositories"
)

type StudentService interface {
	CreateStudent(ctx context.Context, s models.Students) (int, error)
	GetStudents(ctx context.Context) ([]models.Students, error)
	GetStudentById(ctx context.Context, id int) (*models.Students, error)
	UpdateStudent(ctx context.Context, s models.Students) error
	DeleteStudent(ctx context.Context, id int) error
	EmailIsExist(ctx context.Context, email string) (bool, error)
}

type studentService struct {
	repo repositories.StudentRepository
}

func NewStudentRepository(repo repositories.StudentRepository) *studentService {
	return &studentService{repo: repo}
}

func (s *studentService) CreateStudent(ctx context.Context, student models.Students) error {
	if student.Name == "" {
		return errors.New("name is requred")
	}
	return nil
}
