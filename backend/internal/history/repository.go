package history

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

const maxHistoryRows = 100

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new history row and prunes rows beyond the most recent
// maxHistoryRows in the same transaction, so the table never grows past the
// retention window.
func (r *Repository) Create(ctx context.Context, expression, result string) (History, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return History{}, err
	}
	defer tx.Rollback(ctx)

	var h History
	err = tx.QueryRow(ctx,
		`INSERT INTO history (expression, result) VALUES ($1, $2)
		 RETURNING id, expression, result, created_at`,
		expression, result,
	).Scan(&h.ID, &h.Expression, &h.Result, &h.CreatedAt)
	if err != nil {
		return History{}, err
	}

	_, err = tx.Exec(ctx,
		`DELETE FROM history
		 WHERE id NOT IN (
		     SELECT id FROM history ORDER BY created_at DESC, id DESC LIMIT $1
		 )`,
		maxHistoryRows,
	)
	if err != nil {
		return History{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return History{}, err
	}
	return h, nil
}

// List returns up to the most recent maxHistoryRows entries, newest first.
func (r *Repository) List(ctx context.Context) ([]History, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, expression, result, created_at
		 FROM history
		 ORDER BY created_at DESC, id DESC
		 LIMIT $1`,
		maxHistoryRows,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]History, 0)
	for rows.Next() {
		var h History
		if err := rows.Scan(&h.ID, &h.Expression, &h.Result, &h.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, h)
	}
	return items, rows.Err()
}
