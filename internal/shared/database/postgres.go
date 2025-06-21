package database

import (
    "database/sql"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "system-portal/internal/shared/config"
    "system-portal/pkg/logger"

    _ "github.com/lib/pq"
)

// Postgres wraps a sql.DB connection pool and its configuration.
type Postgres struct {
    Config config.DatabaseConfig
    DB     *sql.DB
}

// New creates a new Postgres instance using the provided DatabaseConfig.
// It opens the connection, verifies it via Ping and version check, and returns the Postgres wrapper.
func New(cfg config.DatabaseConfig) (*Postgres, error) {
    // Build DSN from config
    dsn := cfg.DSN()
    logger.Log.WithFields(map[string]interface{}{
        "host":    cfg.Host,
        "port":    cfg.Port,
        "user":    cfg.User,
        "db":      cfg.Name,
        "sslmode": cfg.SSLMode,
    }).Info("connecting to postgres")

    // Open connection using lib/pq driver
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        logger.Log.WithError(err).Error("failed to open postgres connection")
        return nil, err
    }

    // Verify connectivity and credentials via Ping
    if err := db.Ping(); err != nil {
        logger.Log.WithError(err).Error("failed to ping postgres database (check credentials and network)")
        db.Close()
        return nil, err
    }

    // Perform version check to ensure DB responsiveness
    var version string
    if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
        db.Close()
        logger.Log.WithError(err).Error("failed to retrieve postgres version; connection may be invalid")
        return nil, fmt.Errorf("version check failed: %w", err)
    }
    logger.Log.WithField("version", version).Info("postgres version retrieved")

    logger.Log.Info("postgres connection established")
    return &Postgres{Config: cfg, DB: db}, nil
}

// Close closes the database connection pool.
func (p *Postgres) Close() error {
    if p.DB != nil {
        return p.DB.Close()
    }
    return nil
}

// Migrate executes SQL migration files in order, splitting statements and logging details.
func (p *Postgres) Migrate() error {
    // Log working directory for debugging
    cwd, err := os.Getwd()
    if err != nil {
        logger.Log.WithError(err).Error("cannot determine working directory")
    } else {
        logger.Log.WithField("cwd", cwd).Info("working directory")
    }

    dir := filepath.Join("migrations")
    entries, err := os.ReadDir(dir)
    if err != nil {
        logger.Log.WithError(err).WithField("dir", dir).Error("cannot read migrations directory")
        return err
    }
    logger.Log.WithField("count", len(entries)).Info("migration files found")

    sort.Slice(entries, func(i, j int) bool { return entries[i].Name() < entries[j].Name() })
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        filePath := filepath.Join(dir, e.Name())
        logger.Log.WithField("file", e.Name()).Info("applying migration file")

        data, err := os.ReadFile(filePath)
        if err != nil {
            logger.Log.WithError(err).WithField("file", e.Name()).Error("cannot read migration file")
            return err
        }

        // Split SQL into individual statements by semicolon
        statements := strings.Split(string(data), ";")
        for idx, stmt := range statements {
            stmt = strings.TrimSpace(stmt)
            if stmt == "" {
                continue
            }
            logger.Log.WithFields(map[string]interface{}{
                "file":            e.Name(),
                "statement_index": idx,
            }).Info("applying SQL statement")

            res, err := p.DB.Exec(stmt)
            if err != nil {
                logger.Log.WithError(err).WithFields(map[string]interface{}{
                    "file":            e.Name(),
                    "statement_index": idx,
                }).Error("statement execution failed")
                return err
            }

            if rows, err := res.RowsAffected(); err == nil {
                logger.Log.WithFields(map[string]interface{}{
                    "file":            e.Name(),
                    "statement_index": idx,
                    "rows_affected":   rows,
                }).Info("statement executed")
            } else {
                logger.Log.WithError(err).WithField("file", e.Name()).Warn("could not get rows affected for statement")
            }
        }
    }

    logger.Log.Info("database migrations complete")
    return nil
}

