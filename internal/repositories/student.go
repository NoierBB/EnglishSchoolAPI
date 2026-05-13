package repositories

import (
	"context"
	"fmt"

	"github.com/NoierBB/englishSchool/internal/models"
	db "github.com/NoierBB/englishSchool/pkg/db"
)

type StudentRepository struct {
	db *db.PostgresDB
}

func NewStudentRepository(db *db.PostgresDB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) GetStudents(ctx context.Context) ([]models.Students, error) {
	var students []models.Students

	const query = `SELECT id, user_id, name, age, level FROM students`
	rows, err := r.db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select student: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var doc models.Students
		if err := rows.Scan(&doc.Id, &doc.UserId, &doc.Name, &doc.Age, &doc.Level); err != nil {
			return nil, fmt.Errorf("scan student: %w", err)
		}
		students = append(students, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return students, nil
}
