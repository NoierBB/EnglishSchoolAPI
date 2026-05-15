package repositories

import (
	"context"
	"database/sql"
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

func (r *StudentRepository) CreateStudent(ctx context.Context, s models.Students) (int, error) {
	const query = `INSERT INTO students (user_id, name, age, level)

		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := r.db.DB.QueryRowContext(ctx, query,
		s.UserId, s.Name, s.Age, s.Level,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create student: %w", err)
	}
	return id, nil
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

func (r *StudentRepository) GetStudentById(ctx context.Context, id int) (*models.Students, error) {
	const query = `SELECT id, user_id, name, age, level
					FROM students
					WHERE id = $1
	`

	var students models.Students

	err := r.db.DB.QueryRowContext(ctx, query, id).
		Scan(&students.Id, &students.UserId, &students.Name, &students.Age, &students.Level)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("students not found")
		}
		return nil, fmt.Errorf("get students by id: %w", err)
	}
	return &students, nil
}
