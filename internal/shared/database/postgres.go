package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"

	_ "github.com/jackc/pgx/v5/stdlib"
	"system-portal/internal/shared/config"
	"system-portal/pkg/logger"
)

// Postgres wraps a sql.DB connection pool.
type Postgres struct {
	DSN string
	DB  *sql.DB
}

// New opens a PostgreSQL connection.
func New(cfg config.DatabaseConfig) (*Postgres, error) {
	dsn := cfg.DSN()
	logger.Log.WithFields(map[string]interface{}{
		"host":    cfg.Host,
		"port":    cfg.Port,
		"user":    cfg.User,
		"db":      cfg.Name,
		"sslmode": cfg.SSLMode,
	}).Info("connecting to postgres")

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	logger.Log.Info("postgres connection established")
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
		logger.Log.WithField("file", e.Name()).Info("applying migration")
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if _, err := p.DB.Exec(string(data)); err != nil {
			logger.Log.WithError(err).WithField("file", e.Name()).Error("migration failed")
			return err
		}
	}
	logger.Log.Info("database migrations complete")
	return nil
}
