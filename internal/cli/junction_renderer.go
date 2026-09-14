package cli

import (
	"fmt"
)

type JunctionRenderer struct {
	dispatcher *StyleDispatcher
	template   *JunctionTemplate
}

func NewJunctionRenderer(dispatcher *StyleDispatcher, template *JunctionTemplate) *JunctionRenderer {
	return &JunctionRenderer{
		dispatcher: dispatcher,
		template:   template,
	}
}

// RenderMainDashboard draws the primary command center interface
func (jr *JunctionRenderer) RenderMainDashboard(activeWorkers int, vaultCount int) {
	fmt.Println("\n========================================")
	fmt.Println(jr.dispatcher.Colorize("     AUDIO COMMAND CENTER v2.6         ", "primary"))
	fmt.Println("========================================")
	
	// Render operational badges
	fmt.Printf(" Status: %s | Vault Tracks: %s\n",
		jr.dispatcher.RenderBadge("ENGINE", true),
		jr.dispatcher.Colorize(fmt.Sprintf("[%d]", vaultCount), "success"),
	)
	fmt.Println("----------------------------------------")
	fmt.Println(" [1] Search & Download Audio")
	fmt.Println(" [2] View Music Vault & Statistics")
	fmt.Println(" [3] Exit")
	fmt.Println("----------------------------------------")
}

// RenderPrompt prints a styled input indicator
func (jr *JunctionRenderer) RenderPrompt() {
	fmt.Print(jr.dispatcher.Colorize("Select an option [1-3]: ", "bold"))
}
