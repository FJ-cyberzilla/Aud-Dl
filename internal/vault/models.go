package vault

// AcousticFingerprint represents a high-level audio signature from Chromaprint.
type AcousticFingerprint struct {
	FilePath    string
	Duration    float64
	RawSubspace []uint32
}

// HashFingerprint represents a low-level acoustic hash of PCM data.
type HashFingerprint struct {
	Hash   uint64
	Artist string
	Album  string
	Path   string
}
