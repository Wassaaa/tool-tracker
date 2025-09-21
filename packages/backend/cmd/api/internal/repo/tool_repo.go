package repo

import (
	"database/sql"
	"fmt"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type PostgresToolRepo struct {
	db *sql.DB
}

func NewPostgresToolRepo(db *sql.DB) *PostgresToolRepo {
	return &PostgresToolRepo{db: db}
}

// Helper function to define the column order for tool returns
func (r *PostgresToolRepo) toolColumns() string {
	return "id, name, status, current_user_id, last_checked_out_at, created_at, updated_at"
}

// Helper function to scan a row into a Tool struct
func (r *PostgresToolRepo) scanTool(scanner interface {
	Scan(dest ...any) error
}) (domain.Tool, error) {
	var tool domain.Tool
	err := scanner.Scan(
		&tool.ID,
		&tool.Name,
		&tool.Status,
		&tool.CurrentUserId,
		&tool.LastCheckedOutAt,
		&tool.CreatedAt,
		&tool.UpdatedAt,
	)
	return tool, err
}

func (r *PostgresToolRepo) Create(name string, status domain.ToolStatus) (domain.Tool, error) {

	tool, err := domain.NewTool(name, status)
	if err != nil {
		return domain.Tool{}, fmt.Errorf("failed to create tool: %w", err)
	}

	query := `INSERT INTO tools (name, status) VALUES ($1, $2) RETURNING ` + r.toolColumns()
	row := r.db.QueryRow(query, tool.Name, tool.Status)
	createdTool, err := r.scanTool(row)
	if err != nil {
		return domain.Tool{}, fmt.Errorf("failed to create tool: %w", err)
	}

	return createdTool, nil
}

func (r *PostgresToolRepo) List(limit, offset int) ([]domain.Tool, error) {
	query := `SELECT ` + r.toolColumns() + ` FROM tools ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools: %w", err)
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		tool, err := r.scanTool(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		tools = append(tools, tool)
	}

	// Check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tools: %w", err)
	}

	return tools, nil
}

func (r *PostgresToolRepo) Get(id string) (domain.Tool, error) {
	query := `SELECT ` + r.toolColumns() + ` FROM tools WHERE id = $1`

	row := r.db.QueryRow(query, id)
	tool, err := r.scanTool(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Tool{}, domain.ErrToolNotFound
		}
		return domain.Tool{}, fmt.Errorf("failed to get tool: %w", err)
	}

	return tool, nil
}

func (r *PostgresToolRepo) Update(t domain.Tool) (domain.Tool, error) {
	query := `UPDATE tools SET name = $1, status = $2, current_user_id = $3 WHERE id = $4 RETURNING ` + r.toolColumns()

	row := r.db.QueryRow(query, t.Name, t.Status, t.CurrentUserId, t.ID)
	tool, err := r.scanTool(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Tool{}, domain.ErrToolNotFound
		}
		return domain.Tool{}, fmt.Errorf("failed to update tool: %w", err)
	}

	return tool, nil
}

func (r *PostgresToolRepo) Delete(id string) error {
	query := `DELETE FROM tools WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrToolNotFound
	}

	return nil
}

func (r *PostgresToolRepo) ListByStatus(status domain.ToolStatus, limit, offset int) ([]domain.Tool, error) {
	query := `SELECT ` + r.toolColumns() + ` FROM tools WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools by status: %w", err)
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		tool, err := r.scanTool(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		tools = append(tools, tool)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tools: %w", err)
	}

	return tools, nil
}

func (r *PostgresToolRepo) ListByUser(userID string, limit, offset int) ([]domain.Tool, error) {
	query := `SELECT ` + r.toolColumns() + ` FROM tools WHERE current_user_id = $1 ORDER BY last_checked_out_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query tools by user: %w", err)
	}
	defer rows.Close()

	var tools []domain.Tool
	for rows.Next() {
		tool, err := r.scanTool(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tool: %w", err)
		}
		tools = append(tools, tool)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tools: %w", err)
	}

	return tools, nil
}

func (r *PostgresToolRepo) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM tools`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tools: %w", err)
	}
	return count, nil
}
