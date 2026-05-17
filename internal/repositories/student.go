package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NoierBB/englishSchool/internal/models"
)

type StudentsRepository interface {
	CreateStudent(ctx context.Context, s models.Students) (int, error)
	GetStudents(ctx context.Context) ([]models.Students, error)
	GetStudentById(ctx context.Context, id int) (*models.Students, error)
	UpdateStudent(ctx context.Context, s models.Students) error
	DeleteStudent(ctx context.Context, id int) error
}

type StudentRepository struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

func (r *StudentRepository) CreateStudent(ctx context.Context, s models.Students) (int, error) {
	const query = `INSERT INTO students (user_id, name, age, level)

		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query,
		s.UserId, s.Name, s.Age, s.Level,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create student: %w", err)
	}
	return id, nil
}

func (r *StudentRepository) GetStudents(ctx context.Context) ([]models.Students, error) {
	var students []models.Students

	const query = `SELECT id, user_id, name, age, level FROM students ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
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

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&students.Id, &students.UserId, &students.Name, &students.Age, &students.Level)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("students not found")
		}
		return nil, fmt.Errorf("get students by id: %w", err)
	}
	return &students, nil
}

func (r *StudentRepository) UpdateStudent(ctx context.Context, s models.Students) error {
	const query = `UPDATE students
					SET user_id = $1, name = $2, age = $3, level = $4
					WHERE id = $5`

	res, err := r.db.ExecContext(ctx, query, s.UserId, s.Name, s.Age, s.Level, s.Id)

	if err != nil {
		return fmt.Errorf("update student: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("student not found")
	}
	return nil
}

func (r *StudentRepository) DeleteStudent(ctx context.Context, id int) error {
	const query = `DELETE FROM students WHERE id=$1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete student: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("student not found")
	}
	return nil
}
