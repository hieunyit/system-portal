package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/lib/pq"
)

// Postgres wraps a sql.DB connection pool.
type Postgres struct {
	DSN string
	DB  *sql.DB
}

// New opens a PostgreSQL connection.
func New(dsn string) (*Postgres, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{DSN: dsn, DB: db}, nil
}

// Close closes the database connection.
func (p *Postgres) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// Migrate executes embedded SQL migration files in order.
func (p *Postgres) Migrate() error {
	dir := filepath.Join("migrations")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if _, err := p.DB.Exec(string(data)); err != nil {
			return err
		}
	}
	return nil
}
