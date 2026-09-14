package strategy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLimiter struct{}
func (m *mockLimiter) AcquireSlot(ctx context.Context) error { return nil }
func (m *mockLimiter) ReleaseSlot()                        {}

type mockShifter struct{}
func (m *mockShifter) JitterWait(ctx context.Context) error { return nil }

type mockTactics struct{}
func (m *mockTactics) InjectEvasionHeaders(req *http.Request) {}

type mockDissector struct{}
func (m *mockDissector) InspectResponse(resp *http.Response) DissectionResult {
	return DissectionResult{Type: ChallengeNone, StatusCode: 200}
}

type mockFallback struct{}
func (m *mockFallback) ExecuteWithRetry(ctx context.Context, operation func(context.Context) error) error {
	return operation(ctx)
}

func TestAntibotCommander_PrepareAndExecute(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	commander := NewAntibotCommander(
		&mockTactics{},
		&mockShifter{},
		&mockLimiter{},
		&mockDissector{},
		&mockFallback{},
	)

	req, _ := http.NewRequest("GET", server.URL, nil)
	client := server.Client()
	
	resp, err := commander.PrepareAndExecute(context.Background(), req, client)
	if err != nil {
		t.Fatalf("PrepareAndExecute failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
}
