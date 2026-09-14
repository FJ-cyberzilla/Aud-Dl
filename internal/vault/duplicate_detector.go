package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"math/bits"
	"os/exec"
)

type FpCalcOutput struct {
	Duration    float64  `json:"duration"`
	Fingerprint string   `json:"fingerprint"`
	RawArray    []uint32 `json:"raw"`
}

type DuplicateDetector struct {
	index []AcousticFingerprint
}

func NewDuplicateDetector() *DuplicateDetector {
	return &DuplicateDetector{
		index: make([]AcousticFingerprint, 0),
	}
}

// ExtractFingerprint calls the Chromaprint 'fpcalc' binary to get the acoustic signature
func (dd *DuplicateDetector) ExtractFingerprint(ctx context.Context, filePath string) (*AcousticFingerprint, error) {
	cmd := exec.CommandContext(ctx, "fpcalc", "-raw", "-json", filePath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fpcalc execution failed (ensure Chromaprint is installed): %w", err)
	}

	var parsed FpCalcOutput
	if err := json.Unmarshal(output, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse fpcalc JSON response: %w", err)
	}

	return &AcousticFingerprint{
		FilePath:    filePath,
		Duration:    parsed.Duration,
		RawSubspace: parsed.RawArray,
	}, nil
}

// IsDuplicate compares a new track fingerprint against all existing indexed tracks in the Vault
// Similarity threshold ranges from 0.0 (unrelated) to 1.0 (exact audio match).
func (dd *DuplicateDetector) IsDuplicate(target *AcousticFingerprint, threshold float64) (bool, string, float64) {
	for _, indexed := range dd.index {
		// Quick pre-check: Duration must be within 3 seconds of target
		if target.Duration < indexed.Duration-3.0 || target.Duration > indexed.Duration+3.0 {
			continue
		}

		score := calculateSimilarity(target.RawSubspace, indexed.RawSubspace)
		if score >= threshold {
			return true, indexed.FilePath, score
		}
	}

	return false, "", 0.0
}

// IndexTrack adds a newly stored track fingerprint into the Vault index
func (dd *DuplicateDetector) IndexTrack(fp *AcousticFingerprint) {
	dd.index = append(dd.index, *fp)
}

// Bitwise Hamming Distance similarity calculation between acoustic arrays
func calculateSimilarity(fp1, fp2 []uint32) float64 {
	minLen := len(fp1)
	if len(fp2) < minLen {
		minLen = len(fp2)
	}

	if minLen == 0 {
		return 0.0
	}

	var totalBits int
	var matchingBits int

	for i := 0; i < minLen; i++ {
		// XOR highlights differing bits
		diff := fp1[i] ^ fp2[i]
		differingBits := bits.OnesCount32(diff)

		matchingBits += (32 - differingBits)
		totalBits += 32
	}

	return float64(matchingBits) / float64(totalBits)
}
