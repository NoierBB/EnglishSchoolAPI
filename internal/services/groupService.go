package services

import (
	"context"
	"errors"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/repositories"
)

type GroupService interface {
	CreateGroup(ctx context.Context, g models.Group) (int, error)
	GetGroup(ctx context.Context) ([]models.Group, error)
	GetGroupById(ctx context.Context, id int) (*models.Group, error)
	AddStudent(ctx context.Context, groupId, studentId int) error
}

type groupService struct {
	gRepo repositories.GroupRepo
	sRepo repositories.StudentRepository
}

func NewGroupService(repo repositories.GroupRepo) *groupService {
	return &groupService{gRepo: repo}
}

func (s *groupService) AddStudent(ctx context.Context, groupId, studentId int) error {

	_, err := s.gRepo.GetGroupById(ctx, groupId)
	if err != nil {
		return errors.New("group not found")
	}

	_, err = s.sRepo.GetStudentById(ctx, studentId)
	if err != nil {
		return errors.New("student not found")
	}
	return s.gRepo.AddStudent(ctx, groupId, studentId)
}
