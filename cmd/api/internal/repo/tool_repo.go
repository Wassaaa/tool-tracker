package repo

import (
	"database/sql"
	"time"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type PostgresToolRepo struct {
	db *sql.DB
}

func NewPostgresToolRepo(db *sql.DB) *PostgresToolRepo {
	return &PostgresToolRepo{db: db}
}

func (r *PostgresToolRepo) Create(name string) (domain.Tool, error) {
	tool := domain.Tool{
		Name:      name,
		Status:    "IN_OFFICE",
		CreatedAt: time.Now(),
	}
	query := `INSERT INTO tools (name, status, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.db.QueryRow(query, tool.Name, tool.Status, tool.CreatedAt).Scan(&tool.ID)
	return tool, err
}

func (r *PostgresToolRepo) List(limit, offset int) ([]domain.Tool, error) {
	query := `SELECT id, name, status, created_at FROM tools ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		var tool domain.Tool
		err := rows.Scan(&tool.ID, &tool.Name, &tool.Status, &tool.CreatedAt)
		if err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}

	return tools, rows.Err()
}
