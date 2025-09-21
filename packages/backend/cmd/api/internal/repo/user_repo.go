package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

// Helper function to define the column order for user returns
func (r *PostgresUserRepo) userColumns() string {
	return "id, name, email, role, created_at, updated_at"
}

// Helper function to scan a row into a User struct
func (r *PostgresUserRepo) scanUser(scanner interface {
	Scan(dest ...any) error
}) (domain.User, error) {
	var user domain.User
	err := scanner.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	return user, err
}

func (r *PostgresUserRepo) Create(name string, email string, role domain.UserRole) (domain.User, error) {
	user := domain.User{
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO users (name, email, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING ` + r.userColumns()
	row := r.db.QueryRow(query, user.Name, user.Email, user.Role, user.CreatedAt, user.UpdatedAt)
	createdUser, err := r.scanUser(row)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (r *PostgresUserRepo) List(limit, offset int) ([]domain.User, error) {
	query := `SELECT ` + r.userColumns() + ` FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}

func (r *PostgresUserRepo) Get(id string) (domain.User, error) {
	query := `SELECT ` + r.userColumns() + ` FROM users WHERE id = $1`

	row := r.db.QueryRow(query, id)
	user, err := r.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepo) GetByEmail(email string) (domain.User, error) {
	query := `SELECT ` + r.userColumns() + ` FROM users WHERE email = $1`

	row := r.db.QueryRow(query, email)
	user, err := r.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepo) Update(id string, name string, email string, role domain.UserRole) (domain.User, error) {
	query := `UPDATE users SET name = $1, email = $2, role = $3, updated_at = $4 WHERE id = $5 RETURNING ` + r.userColumns()

	row := r.db.QueryRow(query, name, email, role, time.Now(), id)
	user, err := r.scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (r *PostgresUserRepo) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepo) ListByRole(role domain.UserRole, limit, offset int) ([]domain.User, error) {
	query := `SELECT ` + r.userColumns() + ` FROM users WHERE role = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, role, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users by role: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}

func (r *PostgresUserRepo) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
