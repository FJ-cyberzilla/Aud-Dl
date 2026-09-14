package cli

import (
	"context"
	"fmt"
	"audio-command-center/internal/cache"
	"audio-command-center/internal/ui"
)

type SearchExecutionHandler func(query string) error
type AdminStatsProvider func() AdminDashboardStats

type MenuJunction struct {
    onSearch SearchExecutionHandler
    getStats AdminStatsProvider
    theme    ui.Theme
    cache    *cache.CacheManager
}

func NewMenuJunction(theme ui.Theme, cache *cache.CacheManager) *MenuJunction {
    return &MenuJunction{theme: theme, cache: cache}
}

// StartInteractiveSession runs the main interactive loop
func (mj *MenuJunction) StartInteractiveSession(ctx context.Context) error {
    menu := NewInteractiveMenu(mj.theme, mj.cache)

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            choice, err := menu.PromptMainMenu(ctx)
            if err != nil {
                return err
            }

            switch choice {
            case 0: // Search
                query, err := menu.PromptSearchQuery(ctx)
                if err == nil && mj.onSearch != nil {
                    mj.onSearch(query)
                }
            case 1: // Stats
                if mj.getStats != nil {
                    stats := mj.getStats()
                    // This is where we'd render stats, maybe via MenuCrust's admin renderer
                    fmt.Printf("\n--- System Stats ---\nTracks: %d\nStorage: %.2f MB\nWorkers: %d\nSanctions: %d\n",
                        stats.TotalTracks, stats.StorageUsedMB, stats.ActiveWorkers, stats.SanctionedDomains)
                }
            case 2: // Settings
                fmt.Println("Settings not implemented yet.")
            case 3: // Exit
                return nil
            }
        }
    }
}
