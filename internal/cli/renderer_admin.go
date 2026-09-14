package cli

import (
	"fmt"
)

type AdminDashboardStats struct {
	TotalTracks       int
	StorageUsedMB     float64
	ActiveWorkers     int
	SanctionedDomains int
}

type RendererAdmin struct {
	dispatcher *StyleDispatcher
}

func NewRendererAdmin(dispatcher *StyleDispatcher) *RendererAdmin {
	return &RendererAdmin{
		dispatcher: dispatcher,
	}
}

// RenderAdminDashboard draws the deep system diagnostics and storage statistics panel
func (ra *RendererAdmin) RenderAdminDashboard(stats AdminDashboardStats) {
	fmt.Println("\n========================================")
	fmt.Println(ra.dispatcher.Colorize(" VAULT & ENGINE ADMINISTRATION ", "primary"))
	fmt.Println("========================================")
	fmt.Printf(" Ingestion Workers: %s\n", ra.dispatcher.Colorize(fmt.Sprintf("[%d Active]", stats.ActiveWorkers), "success"))
	fmt.Printf(" Vault Database: %s\n", ra.dispatcher.RenderBadge("SQLITE", true))
	fmt.Printf(" Total Track Count: %s\n", ra.dispatcher.Colorize(fmt.Sprintf("%d tracks", stats.TotalTracks), "primary"))
	fmt.Printf(" Storage Consumed: %s\n", ra.dispatcher.Colorize(fmt.Sprintf("%.2f MB", stats.StorageUsedMB), "warning"))
	fmt.Printf(" Blocked Domains: %s\n", ra.dispatcher.Colorize(fmt.Sprintf("%d endpoints", stats.SanctionedDomains), "danger"))
	fmt.Println("----------------------------------------")
	fmt.Println(" [0] Return to Main Menu")
	fmt.Println("----------------------------------------")
}

// RenderPrompt prints a navigation indicator for the admin panel
func (ra *RendererAdmin) RenderPrompt() {
	fmt.Print(ra.dispatcher.Colorize("Press [0] to go back: ", "bold"))
}
