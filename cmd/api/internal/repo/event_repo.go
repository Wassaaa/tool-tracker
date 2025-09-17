package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

type PostgresEventRepo struct {
	db *sql.DB
}

func NewPostgresEventRepo(db *sql.DB) *PostgresEventRepo {
	return &PostgresEventRepo{db: db}
}

// Helper function to define the column order for event returns
func (r *PostgresEventRepo) eventColumns() string {
	return "id, type, tool_id, user_id, actor_id, notes, metadata, created_at"
}

// Helper function to scan a row into an Event struct
func (r *PostgresEventRepo) scanEvent(scanner interface {
	Scan(dest ...any) error
}) (domain.Event, error) {
	var event domain.Event
	err := scanner.Scan(
		&event.ID,
		&event.Type,
		&event.ToolID,
		&event.UserID,
		&event.ActorID,
		&event.Notes,
		&event.Metadata,
		&event.CreatedAt,
	)
	return event, err
}

func (r *PostgresEventRepo) Create(eventType domain.EventType, toolID *string, userID *string, actorID *string, notes string, metadata string) (domain.Event, error) {
	event := domain.Event{
		Type:      eventType,
		ToolID:    toolID,
		UserID:    userID,
		ActorID:   actorID,
		Notes:     notes,
		Metadata:  metadata,
		CreatedAt: time.Now(),
	}

	query := `INSERT INTO events (type, tool_id, user_id, actor_id, notes, metadata, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING ` + r.eventColumns()
	row := r.db.QueryRow(query, event.Type, event.ToolID, event.UserID, event.ActorID, event.Notes, event.Metadata, event.CreatedAt)
	createdEvent, err := r.scanEvent(row)
	if err != nil {
		return domain.Event{}, fmt.Errorf("failed to create event: %w", err)
	}

	return createdEvent, nil
}

func (r *PostgresEventRepo) List(limit, offset int) ([]domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event, err := r.scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

func (r *PostgresEventRepo) Get(id string) (domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events WHERE id = $1`

	row := r.db.QueryRow(query, id)
	event, err := r.scanEvent(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Event{}, domain.ErrEventNotFound
		}
		return domain.Event{}, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

func (r *PostgresEventRepo) ListByType(eventType domain.EventType, limit, offset int) ([]domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events WHERE type = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by type: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event, err := r.scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

func (r *PostgresEventRepo) ListByTool(toolID string, limit, offset int) ([]domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events WHERE tool_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, toolID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by tool: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event, err := r.scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

func (r *PostgresEventRepo) ListByUser(userID string, limit, offset int) ([]domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events WHERE user_id = $1 OR actor_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query events by user: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event, err := r.scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// EventFilter represents filtering options for events
type EventFilter struct {
	Type   *domain.EventType
	ToolID *string
	UserID *string
}

func (r *PostgresEventRepo) ListWithFilter(filter EventFilter, limit, offset int) ([]domain.Event, error) {
	query := `SELECT ` + r.eventColumns() + ` FROM events WHERE 1=1`
	args := []interface{}{}
	argIndex := 1

	if filter.Type != nil {
		query += fmt.Sprintf(` AND type = $%d`, argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.ToolID != nil {
		query += fmt.Sprintf(` AND tool_id = $%d`, argIndex)
		args = append(args, *filter.ToolID)
		argIndex++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(` AND (user_id = $%d OR actor_id = $%d)`, argIndex, argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events with filter: %w", err)
	}
	defer rows.Close()

	var events []domain.Event
	for rows.Next() {
		event, err := r.scanEvent(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

func (r *PostgresEventRepo) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM events`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}
	return count, nil
}
