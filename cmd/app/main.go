package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"audio-command-center/internal/cli"
	"audio-command-center/internal/task"
	"audio-command-center/internal/ui"
)

func main() {
	// 1. Establish root context linked to OS signals (Handles graceful shutdown on Ctrl+C / SIGTERM)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 2. Initialize UI layout and styling components
	refreshManager := ui.NewRefreshManager()
	refreshManager.ClearScreen()

	// 3. Define decoupled action handlers (Avoiding circular imports!)
	// These act as bridges between the CLI presentation layer and your underlying tasks/vaults.
	// For simplicity, using TaskAdmin to access system components.
	ta, err := task.NewTaskAdmin("1.0.0", "./config", "./music_vault", ui.CyberDarkTheme)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize system: %v\n", err)
		os.Exit(1)
	}
	defer ta.Bindings.VaultAdmin.Close()

	searchHandler := func(query string) error {
		fmt.Printf("\n[Dispatcher] Executing search query: '%s'\n", query)
		_, err := ta.Bindings.SearchManager.QueryFallback(ctx, query)
		return err
	}

	statsProvider := func() cli.AdminDashboardStats {
		count, size, err := ta.Bindings.VaultAdmin.GetStats()
		if err != nil {
			return cli.AdminDashboardStats{}
		}
		return cli.AdminDashboardStats{
			TotalTracks:       count,
			StorageUsedMB:     size,
			ActiveWorkers:     ta.GetController().GetDivision().GetActiveWorkers(),
			SanctionedDomains: ta.Bindings.GateHub.GetSanctioner().GetSanctionedCount(),
		}
	}

	// 4. Instantiate the All-in-One CLI Skin Facade with Cyber Dark theme
	skin := cli.NewMenuCrust(cli.ThemeDark)
	skin.SetSearchHandler(searchHandler)
	skin.SetStatsProvider(statsProvider)

	fmt.Println(skin.Dispatcher.Colorize("==> Audio Command Center Initialized Successfully.", "success"))

	// 5. Run the interactive session inside a goroutine listening to context cancellation
	errChan := make(chan error, 1)
	go func() {
		errChan <- skin.BootSkinFacade(ctx)
	}()

	// 6. Select block for graceful exit
	select {
	case <-ctx.Done():
		fmt.Println("\n\033[33m[System] Shutdown signal received. Flushing buffers and exiting gracefully...\033[0m")
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n\033[31m[Error] Session terminated with error: %v\033[0m\n", err)
		}
	}

	refreshManager.ClearScreen()
	fmt.Println("Goodbye!")
}
