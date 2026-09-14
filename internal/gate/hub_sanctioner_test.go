package gate

import (
	"testing"
	"time"
)

func TestHubSanctioner_RecordViolationAndInspect(t *testing.T) {
	maxOffenses := 2
	banDuration := 100 * time.Millisecond
	hs := NewHubSanctioner(maxOffenses, banDuration)
	domain := "example.com"

	// 1. First violation: Should be Warning
	level := hs.RecordViolation(domain, "First strike")
	if level != SanctionWarning {
		t.Errorf("Expected Warning, got %v", level)
	}

	// 2. Second violation: Should escalate to TempBan
	level = hs.RecordViolation(domain, "Second strike")
	if level != SanctionTempBan {
		t.Errorf("Expected TempBan, got %v", level)
	}

	// 3. InspectSanction should return error
	err := hs.InspectSanction(domain)
	if err == nil {
		t.Error("Expected sanction error, got nil")
	}

	// 4. Wait for ban to expire
	time.Sleep(150 * time.Millisecond)

	// 5. InspectSanction should now be clear
	err = hs.InspectSanction(domain)
	if err != nil {
		t.Errorf("Expected sanction to be cleared, got error: %v", err)
	}
}

func TestHubSanctioner_ClearSanction(t *testing.T) {
	hs := NewHubSanctioner(1, 1*time.Minute)
	domain := "clearme.com"

	hs.RecordViolation(domain, "strike")
	hs.ClearSanction(domain)

	if hs.GetSanctionedCount() != 0 {
		t.Error("Expected sanction count to be 0 after clear")
	}
}
