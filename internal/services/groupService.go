package services

import (
	"context"

	"github.com/NoierBB/englishSchool/internal/models"
	"github.com/NoierBB/englishSchool/internal/repositories"
)

type GroupService interface {
	CreateGroup(ctx context.Context, g models.Group) (int, error)
	GetGroup(ctx context.Context) ([]models.Group, error)
}

type groupService struct {
	repo repositories.GroupRepo
}

func NewGroupService(repo repositories.GroupRepo) *groupService {
	return &groupService{repo: repo}
}
