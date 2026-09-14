package vault

import (
	"testing"
)

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		fp1      []uint32
		fp2      []uint32
		expected float64
	}{
		{
			name:     "identical arrays",
			fp1:      []uint32{0xFFFFFFFF, 0x00000000},
			fp2:      []uint32{0xFFFFFFFF, 0x00000000},
			expected: 1.0,
		},
		{
			name:     "completely different arrays",
			fp1:      []uint32{0xFFFFFFFF},
			fp2:      []uint32{0x00000000},
			expected: 0.0,
		},
		{
			name:     "partially similar",
			fp1:      []uint32{0b00000000000000000000000011110000}, // 4 bits match
			fp2:      []uint32{0b00000000000000000000000011111111}, // 8 bits total, 32 bits compared
			expected: 0.875, // (32-4)/32 = 28/32 = 0.875
		},
		{
			name:     "different lengths",
			fp1:      []uint32{0xFFFFFFFF, 0xFFFFFFFF},
			fp2:      []uint32{0xFFFFFFFF},
			expected: 1.0, // Should use minLen
		},
		{
			name:     "empty arrays",
			fp1:      []uint32{},
			fp2:      []uint32{},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateSimilarity(tt.fp1, tt.fp2)
			if got != tt.expected {
				t.Errorf("calculateSimilarity() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDuplicateDetector_IsDuplicate(t *testing.T) {
	dd := NewDuplicateDetector()

	// Add some data
	dd.IndexTrack(&AcousticFingerprint{
		FilePath:    "track1.mp3",
		Duration:    100.0,
		RawSubspace: []uint32{0xFFFFFFFF},
	})
	dd.IndexTrack(&AcousticFingerprint{
		FilePath:    "track2.mp3",
		Duration:    200.0,
		RawSubspace: []uint32{0xAAAAAAAA},
	})

	tests := []struct {
		name      string
		target    *AcousticFingerprint
		threshold float64
		want      bool
		wantPath  string
	}{
		{
			name: "exact match",
			target: &AcousticFingerprint{
				FilePath:    "target.mp3",
				Duration:    100.0,
				RawSubspace: []uint32{0xFFFFFFFF},
			},
			threshold: 0.9,
			want:      true,
			wantPath:  "track1.mp3",
		},
		{
			name: "duration mismatch",
			target: &AcousticFingerprint{
				FilePath:    "target.mp3",
				Duration:    105.0, // Outside 3s threshold
				RawSubspace: []uint32{0xFFFFFFFF},
			},
			threshold: 0.9,
			want:      false,
			wantPath:  "",
		},
		{
			name: "similarity too low",
			target: &AcousticFingerprint{
				FilePath:    "target.mp3",
				Duration:    100.0,
				RawSubspace: []uint32{0x00000000},
			},
			threshold: 0.9,
			want:      false,
			wantPath:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, path, _ := dd.IsDuplicate(tt.target, tt.threshold)
			if got != tt.want {
				t.Errorf("IsDuplicate() = %v, want %v", got, tt.want)
			}
			if path != tt.wantPath {
				t.Errorf("IsDuplicate() path = %v, want %v", path, tt.wantPath)
			}
		})
	}
}
