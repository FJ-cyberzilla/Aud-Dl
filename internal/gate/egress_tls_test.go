package gate

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestFetchAudioStream_RetryBehavior(t *testing.T) {
	var attempts int32

	// Mock server that fails 2 times then succeeds
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt32(&attempts, 1)
		if current <= 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Println("[Test] Expected failure (503) on attempt", current)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("audio-data"))
		fmt.Println("[Test] Expected success on attempt", current)
	}))
	defer ts.Close()

	client := NewEgressTLSClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("[Test] Initiating fetch with expected retries...")
	stream, err := client.FetchAudioStream(ctx, ts.URL)
	if err != nil {
		t.Fatalf("Fetch failed unexpectedly: %v", err)
	}
	defer stream.Close()

	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("Expected 3 attempts (2 fails + 1 success), got %d", attempts)
	}
	fmt.Println("[Test] Pipeline verification successful.")
}
