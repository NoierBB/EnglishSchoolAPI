package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/NoierBB/englishSchool/internal/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, u models.User) (int, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, id int) error
	ExistByEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Register(ctx context.Context, email, password string) (int, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	const query = `SELECT id, email, password_hash, role
	FROM users WHERE email = $1`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return &user, err
}

func (r *UserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	var exist bool

	const query = `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`

	err := r.db.QueryRowContext(ctx, query, email).Scan(&exist)

	return exist, err
}

func (r *UserRepository) CreateUserTx(ctx context.Context, tx *sql.Tx, u models.User) (int, error) {
	const query = `INSERT INTO users (email, password_hash, role)
	
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	err := tx.QueryRowContext(ctx, query, u.Email, u.Password, u.Role).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, err
}

func (r *UserRepository) CreateUser(ctx context.Context, u models.User) (int, error) {
	const query = `INSERT INTO users (email, password_hash, role)
	
	VALUES ($1, $2, $3)
	RETURNING id`

	var id int
	err := r.db.QueryRowContext(ctx, query, u.Email, u.Password, u.Role).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, err
}

func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	const query = `SELECT id, email, password_hash, role FROM users ORDER BY id`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Id, &u.Email, &u.Password, &u.Role); err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id int) (*models.User, error) {
	const query = `SELECT id, email, password_hash, role
	FROM users
	WHERE id = $1`

	var users models.User

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&users.Id, &users.Email, &users.Password, &users.Role)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, u models.User) error {
	const query = `UPDATE users
	SET email=$1, password_hash=$2, role=$3
	WHERE id=$4`

	res, err := r.db.ExecContext(ctx, query, u.Email, u.Password, u.Role, u.Id)

	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	const query = `DELETE FROM users WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
