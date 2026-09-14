package strategy

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"
)

type TempDBResolver struct {
	mu          sync.Mutex
	baseDir     string
	connections map[string]*sql.DB
}

func NewTempDBResolver(baseDir string) (*TempDBResolver, error) {
	if baseDir == "" {
		baseDir = filepath.Join(os.TempDir(), "audio_cc_strategy")
	}
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create temp resolver dir: %w", err)
	}
	return &TempDBResolver{
		baseDir:     baseDir,
		connections: make(map[string]*sql.DB),
	}, nil
}

// GetSharedHandle returns an optimized, single-writer multi-reader SQLite handle with busy-wait timeout
func (r *TempDBResolver) GetSharedHandle(dbName string) (*sql.DB, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if db, exists := r.connections[dbName]; exists {
		return db, nil
	}

	dbPath := filepath.Join(r.baseDir, dbName+".db")
	// Enforce WAL journal mode, 10s busy timeout, and memory mapping for maximum throughput
	dsn := fmt.Sprintf("%s?_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)&_pragma=busy_timeout(10000)&_pragma=mmap_size(268435456)", dbPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open shared temp db [%s]: %w", dbName, err)
	}

	// Optimize connection pool limits for high concurrency
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	r.connections[dbName] = db
	return db, nil
}

func (r *TempDBResolver) CleanupAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for name, db := range r.connections {
		_ = db.Close()
		_ = os.Remove(filepath.Join(r.baseDir, name+".db"))
		_ = os.Remove(filepath.Join(r.baseDir, name+".db-wal"))
		_ = os.Remove(filepath.Join(r.baseDir, name+".db-shm"))
	}
	_ = os.RemoveAll(r.baseDir)
}
