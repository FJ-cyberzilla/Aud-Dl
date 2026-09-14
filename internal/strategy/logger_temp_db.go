package strategy

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "modernc.org/sqlite" // Pure Go CGO-free SQLite driver
)


// LoggerTempDB manages temporary logs in a SQLite database.
type LoggerTempDB struct {
	mu      sync.Mutex
	db      *sql.DB
	maxRows int
	dbPath  string
}

// NewLoggerTempDB initializes a new temporary log database.
func NewLoggerTempDB(tempDir string, maxRows int) (*LoggerTempDB, error) {
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	dbPath := filepath.Join(tempDir, "antibot_tactical_logs.db")

	db, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=synchronous(OFF)")
	if err != nil {
		return nil, fmt.Errorf("failed to open temp log db: %w", err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS tactical_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
			level TEXT,
			component TEXT,
			message TEXT
		);`,
		`CREATE INDEX IF NOT EXISTS idx_logs_ts ON tactical_logs(id);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to initialize temp log schema: %w", err)
		}
	}

	return &LoggerTempDB{
		db:      db,
		maxRows: maxRows,
		dbPath:  dbPath,
	}, nil
}

// PushLog inserts a log entry on disk and prunes old rows if max capacity is exceeded.
func (l *LoggerTempDB) PushLog(ctx context.Context, level, component, message string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	query := `INSERT INTO tactical_logs (timestamp, level, component, message) VALUES (?, ?, ?, ?)`
	_, err := l.db.ExecContext(ctx, query, time.Now().Format(time.RFC3339), level, component, message)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	// Prune old logs to bound disk and memory footprint
	pruneQuery := `DELETE FROM tactical_logs WHERE id <= (SELECT MAX(id) FROM tactical_logs) - ?`
	_, err = l.db.ExecContext(ctx, pruneQuery, l.maxRows)
	if err != nil {
		return fmt.Errorf("failed to prune logs: %w", err)
	}
	return nil
}

// FetchRecent retrieves the N most recent log entries for real-time UI streaming.
func (l *LoggerTempDB) FetchRecent(ctx context.Context, limit int) ([]TacticalEvent, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	rows, err := l.db.QueryContext(ctx, `
		SELECT timestamp, level, component, message 
		FROM tactical_logs 
		ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch logs: %w", err)
	}
	defer rows.Close()

	var events []TacticalEvent
	for rows.Next() {
		var ev TacticalEvent
		var tsStr string
		var lvlStr string
		if err := rows.Scan(&tsStr, &lvlStr, &ev.Component, &ev.Message); err != nil {
			return nil, fmt.Errorf("failed to scan log row: %w", err)
		}
		ev.Level = ParseLogLevel(lvlStr)
		ev.Timestamp, err = time.Parse(time.RFC3339, tsStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse log timestamp: %w", err)
		}
		events = append(events, ev)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return events, nil
}

// Close flushes and cleans up the temporary database file on shutdown.
func (l *LoggerTempDB) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.db != nil {
		if err := l.db.Close(); err != nil {
			return fmt.Errorf("failed to close db: %w", err)
		}
	}
	if err := os.Remove(l.dbPath); err != nil {
		return fmt.Errorf("failed to remove db file: %w", err)
	}
	return nil
}
