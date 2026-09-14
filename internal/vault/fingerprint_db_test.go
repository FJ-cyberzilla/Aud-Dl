package vault

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSerializationHelpers(t *testing.T) {
	input := []uint32{0x12345678, 0x87654321, 0x00000000}
	blob := serializeUint32Slice(input)
	output := deserializeUint32Slice(blob)

	if !reflect.DeepEqual(input, output) {
		t.Errorf("Serialization/Deserialization failed: got %v, want %v", output, input)
	}
}

func TestSQLiteFingerprintDB_Flow(t *testing.T) {
	// Setup temp dir
	tmpDir, err := os.MkdirTemp("", "vault_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	db, err := NewSQLiteFingerprintDB(tmpDir)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer db.Close()

	// Prepare dummy file
	dummyFile := filepath.Join(tmpDir, "test.mp3")
	if err := os.WriteFile(dummyFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	fp := &AcousticFingerprint{
		FilePath:    dummyFile,
		Duration:    123.45,
		RawSubspace: []uint32{0x11111111, 0x22222222},
	}

	ctx := context.Background()
	if err := db.SaveFingerprint(ctx, fp); err != nil {
		t.Fatalf("SaveFingerprint failed: %v", err)
	}

	records, err := db.LoadAllFingerprints(ctx)
	if err != nil {
		t.Fatalf("LoadAllFingerprints failed: %v", err)
	}

	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}

	if records[0].FilePath != fp.FilePath || records[0].Duration != fp.Duration || !reflect.DeepEqual(records[0].RawSubspace, fp.RawSubspace) {
		t.Errorf("Loaded record does not match saved record")
	}
}
