package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NoierBB/englishSchool/internal/models"
)

type GroupRepository interface {
	CreateGroup(ctx context.Context, g models.Group) (int, error)
	GetGroup(ctx context.Context) ([]models.Group, error)
	GetGroupById(ctx context.Context, id int) (*models.Group, error)
	AddStudent(ctx context.Context, groupId, studentId int) error
}

type GroupRepo struct {
	db *sql.DB
}

func NewGropRepo(db *sql.DB) *GroupRepo {
	return &GroupRepo{db: db}
}

func (r *GroupRepo) CreateGroup(ctx context.Context, g models.Group) (int, error) {
	const query = `INSERT INTO groups (name, level, teacher_id)
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		g.Name, g.Level, g.TeacherId).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create group: %w", err)
	}
	return id, nil
}

func (r *GroupRepo) GetGroup(ctx context.Context) ([]models.Group, error) {
	var groups []models.Group

	const query = `SELECT id, name, level, teacher_id FROM groups ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select group: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var g models.Group
		if err := rows.Scan(&g.Id, &g.Name, &g.Level, &g.TeacherId); err != nil {
			return nil, fmt.Errorf("scan groups: %w", err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return groups, nil
}

func (r *GroupRepo) GetGroupById(ctx context.Context, id int) (*models.Group, error) {
	const query = `SELECT id, name FROM groups WHERE id = $1`

	var group models.Group
	err := r.db.QueryRowContext(ctx, query, id).Scan(&group.Id, &group.Name)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *GroupRepo) AddStudent(ctx context.Context, groupId, studentId int) error {
	const query = `INSERT INTO group_students (group_id, student_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, groupId, studentId)
	return err
}
