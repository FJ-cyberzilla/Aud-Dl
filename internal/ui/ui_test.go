package ui

import (
	"strings"
	"testing"
)

func TestDividers_RenderPaginationFooter(t *testing.T) {
	d := NewDividers()
	footer := d.RenderPaginationFooter(1, 5)
	
	if !strings.Contains(footer, "Page [1 / 5]") {
		t.Errorf("Expected footer to contain Page [1 / 5], got: %s", footer)
	}
	if !strings.Contains(footer, "[N]ext") {
		t.Error("Expected footer to contain [N]ext")
	}
}
