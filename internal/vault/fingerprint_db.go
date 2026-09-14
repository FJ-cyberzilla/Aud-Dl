package vault

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite" // Pure Go CGO-free SQLite driver
)

type SQLiteFingerprintDB struct {
	db *sql.DB
}

func NewSQLiteFingerprintDB(vaultBasePath string) (*SQLiteFingerprintDB, error) {
	dbPath := filepath.Join(vaultBasePath, "vault_fingerprints.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	// Enable WAL mode for high-concurrency read/write operations
	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA synchronous=NORMAL;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set WAL mode: %w", err)
	}

	repo := &SQLiteFingerprintDB{db: db}
	if err := repo.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return repo, nil
}

func (s *SQLiteFingerprintDB) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS fingerprints (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT UNIQUE NOT NULL,
		duration REAL NOT NULL,
		raw_subspace BLOB NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_duration ON fingerprints(duration);
	`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

// SaveFingerprint serializes raw uint32 slices to binary BLOB and stores in SQLite
func (s *SQLiteFingerprintDB) SaveFingerprint(ctx context.Context, fp *AcousticFingerprint) error {
	blob := serializeUint32Slice(fp.RawSubspace)

	query := `
	INSERT INTO fingerprints (file_path, duration, raw_subspace)
	VALUES (?, ?, ?)
	ON CONFLICT(file_path) DO UPDATE SET
		duration = excluded.duration,
		raw_subspace = excluded.raw_subspace;
	`
	_, err := s.db.ExecContext(ctx, query, fp.FilePath, fp.Duration, blob)
	if err != nil {
		return fmt.Errorf("failed to persist fingerprint: %w", err)
	}
	return nil
}

// LoadAllFingerprints retrieves all indexed records into memory on app startup
func (s *SQLiteFingerprintDB) LoadAllFingerprints(ctx context.Context) ([]AcousticFingerprint, error) {
	query := `SELECT file_path, duration, raw_subspace FROM fingerprints`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query fingerprints: %w", err)
	}
	defer rows.Close()

	var records []AcousticFingerprint
	for rows.Next() {
		var fp AcousticFingerprint
		var blob []byte

		if err := rows.Scan(&fp.FilePath, &fp.Duration, &blob); err != nil {
			return nil, err
		}

		// Verify file still exists on disk; if deleted, skip
		if _, err := os.Stat(fp.FilePath); os.IsNotExist(err) {
			continue
		}

		fp.RawSubspace = deserializeUint32Slice(blob)
		records = append(records, fp)
	}

	return records, nil
}

func (s *SQLiteFingerprintDB) Close() error {
	return s.db.Close()
}

// Binary Serialization Helpers
func serializeUint32Slice(slice []uint32) []byte {
	buf := make([]byte, len(slice)*4)
	for i, v := range slice {
		binary.LittleEndian.PutUint32(buf[i*4:], v)
	}
	return buf
}

func deserializeUint32Slice(b []byte) []uint32 {
	slice := make([]uint32, len(b)/4)
	for i := range slice {
		slice[i] = binary.LittleEndian.Uint32(b[i*4:])
	}
	return slice
}
