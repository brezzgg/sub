package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/brezzgg/sub/internal/models"
	_ "modernc.org/sqlite"
)

// SqliteRepoTest is implementation of [models.Repo]
// just for tests. Not production ready.
type SqliteRepoTest struct {
	db *sql.DB
}

func NewRepo(fileName string) (*SqliteRepoTest, error) {
	db, err := sql.Open("sqlite", fileName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS subs (
			id         TEXT PRIMARY KEY,
			payload    BLOB,
			ttl        INTEGER NOT NULL,
			created_at INTEGER NOT NULL
		) WITHOUT ROWID`)
	if err != nil {
		return nil, fmt.Errorf("migrate error: %w", err)
	}

	return &SqliteRepoTest{db: db}, nil
}

// Get implements [models.Repo].
func (r *SqliteRepoTest) Get(id string) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRow(`SELECT payload, ttl, created_at FROM subs WHERE id = ?`, id).
		Scan(&s.Payload, &s.TTL, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &s, nil
}

// GetAll implements [models.Repo].
func (r *SqliteRepoTest) GetAll() (map[string]*models.Subscription, error) {
	rows, err := r.db.Query(`SELECT id, payload, ttl, created_at FROM subs`)
	if err != nil {
		return nil, models.ErrInternal
	}
	defer rows.Close()

	result := make(map[string]*models.Subscription)
	for rows.Next() {
		var s models.Subscription
		var id string
		if err := rows.Scan(&id, &s.Payload, &s.TTL, &s.CreatedAt); err != nil {
			return nil, models.ErrInternal
		}
		if s.Validate() {
			result[id] = &s
		}
	}
	return result, rows.Err()
}

// GetFunc implements [models.Repo].
func (r *SqliteRepoTest) GetFunc(fn func(s *models.Subscription) bool) ([]*models.Subscription, error) {
	rows, err := r.db.Query(`SELECT payload, ttl, created_at FROM subs`)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	defer rows.Close()

	var result []*models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.Payload, &s.TTL, &s.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		if fn(&s) {
			result = append(result, &s)
		}
	}
	return result, rows.Err()
}

// Remove implements [models.Repo].
func (r *SqliteRepoTest) Remove(id string) error {
	res, err := r.db.Exec(`DELETE FROM subs WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("sub not found")
	}
	return nil
}

// Set implements [models.Repo].
func (r *SqliteRepoTest) Set(id string, s *models.Subscription) error {
	if !s.Validate() {
		return fmt.Errorf("sub not valid")
	}
	_, err := r.db.Exec(
		`INSERT INTO subs(id, payload, ttl, created_at)
		 VALUES(?, ?, ?, ?)
	     ON CONFLICT(id) DO UPDATE SET
         payload=excluded.payload,
         ttl=excluded.ttl,
         created_at=excluded.created_at`,
		id, s.Payload, s.TTL, s.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	return nil
}

var _ models.Repo = (*SqliteRepoTest)(nil)
