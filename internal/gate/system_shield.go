package gate

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	ErrCorruptedStream = errors.New("pipeline integrity error: stream content mismatch or truncation")
	ErrForbiddenContent = errors.New("pipeline security error: ingress stream returned blocked or HTML payload")
)

type SystemShield struct {
	maxMemorySniff int64
}

func NewSystemShield() *SystemShield {
	return &SystemShield{
		maxMemorySniff: 512, // Standard sniff window
	}
}

// InspectStream checks the initial bytes of a stream to verify it's valid media
// and not an accidental HTTP error page, JSON payload, or zero-byte corruption.
func (ss *SystemShield) InspectStream(r io.Reader, safeName string) (io.Reader, error) {
	// Read the first 512 bytes for MIME type sniffing
	buf := make([]byte, ss.maxMemorySniff)
	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return nil, fmt.Errorf("%w: failed to read stream header for %s: %v", ErrCorruptedStream, safeName, err)
	}

	if n == 0 {
		return nil, fmt.Errorf("%w: stream for %s is completely empty (0 bytes)", ErrCorruptedStream, safeName)
	}

	detectedType := http.DetectContentType(buf[:n])

	// Guard against common scraper traps (getting an HTML error page instead of audio)
	if detectedType == "text/html" || detectedType == "application/json" {
		return nil, fmt.Errorf("%w: received %s instead of audio stream for %s", ErrForbiddenContent, detectedType, safeName)
	}

	// Rechain the sniffed bytes back with the rest of the stream using io.MultiReader
	restoredStream := io.MultiReader(BytesReader(buf[:n]), r)

	return restoredStream, nil
}

// Helper structure to convert bytes back to an io.Reader
type byteReader struct {
	data []byte
	off  int
}

func BytesReader(b []byte) io.Reader {
	return &byteReader{data: b}
}

func (br *byteReader) Read(p []byte) (n int, err error) {
	if br.off >= len(br.data) {
		return 0, io.EOF
	}
	n = copy(p, br.data[br.off:])
	br.off += n
	return n, nil
}
