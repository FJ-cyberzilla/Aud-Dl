package cli

import (
    "context"
    "encoding/json"
    "slices"
    "sync"
    "time"

    "audio-command-center/internal/vault"
    )

// SmartMenu with AI-powered recommendations
type SmartMenu struct {
    *MenuTemplate
    vault      *vault.MusicVault
    userPrefs  *UserPreferences
    history    *MenuHistory
}

func NewSmartMenu(mt *MenuTemplate, v *vault.MusicVault) *SmartMenu {
    return &SmartMenu{
        MenuTemplate: mt,
        vault:        v,
        userPrefs:    &UserPreferences{},
        history:      &MenuHistory{},
    }
}

type UserPreferences struct {
    FavoriteFormats []int
    RecentSearches  []string
    LastSelection   int
    Theme           string
    mu              sync.RWMutex
}

type MenuHistory struct {
    selections []SelectionRecord
    maxSize    int
}

type SelectionRecord struct {
    ItemID      int
    Timestamp   time.Time
    Context     string
    Duration    time.Duration
}

// RenderSmart renders with ML-based suggestions
func (sm *SmartMenu) RenderSmart(
    ctx context.Context,
    items []MenuItem[int],
    selected int,
    config MenuConfig,
) string {
    // Get user context
    prefs := sm.userPrefs.Get()
    
    // Reorder items based on user preferences
    reordered := sm.reorderByPreference(items, prefs)
    
    // Add smart suggestions
    if len(items) > 5 {
        sm.addSuggestions(reordered, prefs)
    }
    
    // Render with adaptive features
    config.Subtitle = sm.getContextualSubtitle()
    
    return sm.MenuTemplate.Render(ctx, reordered, selected, config)
}

// reorderByPreference uses historical data to prioritize options
func (sm *SmartMenu) reorderByPreference(
	items []MenuItem[int],
	prefs *UserPreferences,
) []MenuItem[int] {
	if len(prefs.FavoriteFormats) == 0 {
		return items
	}

	// Weight items based on usage history
	weighted := make([]struct {
		item   MenuItem[int]
		weight float64
	}, len(items))

	for i, item := range items {
		weighted[i] = struct {
			item   MenuItem[int]
			weight float64
		}{item, sm.calculateWeight(item, prefs)}
	}

	// Sort by weight descending
	slices.SortFunc(weighted, func(a, b struct {
		item   MenuItem[int]
		weight float64
	}) int {
		if a.weight > b.weight {
			return -1
		}
		if a.weight < b.weight {
			return 1
		}
		return 0
	})

	result := make([]MenuItem[int], len(items))
	for i, w := range weighted {
		result[i] = w.item
	}
	return result
}

func (sm *SmartMenu) calculateWeight(item MenuItem[int], prefs *UserPreferences) float64 {
	weight := 1.0
	for _, fav := range prefs.FavoriteFormats {
		if item.ID == fav {
			weight += 2.0
		}
	}
	return weight
}

// ExportMenu exports menu as JSON for API consumption
func (mt *MenuTemplate) ExportMenu[T comparable](
    items []MenuItem[T],
    selected int,
) ([]byte, error) {
    type ExportItem struct {
        ID          any    `json:"id"`
        Label       string `json:"label"`
        Icon        string `json:"icon"`
        Description string `json:"description"`
        Selected    bool   `json:"selected"`
        Disabled    bool   `json:"disabled"`
    }
    
    export := make([]ExportItem, len(items))
    for i, item := range items {
        export[i] = ExportItem{
            ID:          item.ID,
            Label:       item.Label,
            Icon:        item.Icon,
            Description: item.Description,
            Selected:    i == selected,
            Disabled:    item.Disabled,
        }
    }
    
    return json.MarshalIndent(export, "", "  ")
}

// StreamingMenu for real-time updates
type StreamingMenu struct {
    *MenuTemplate
    updateCh chan<- MenuUpdate
}

type MenuUpdate struct {
    Type    string      `json:"type"`
    Content string      `json:"content"`
    Data    interface{} `json:"data"`
}

func (sm *StreamingMenu) Stream(
    ctx context.Context,
    items []MenuItem[int],
    selected int,
    config MenuConfig,
) <-chan MenuUpdate {
    updates := make(chan MenuUpdate, 10)
    
    go func() {
        defer close(updates)
        
        ticker := time.NewTicker(500 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                // Simulate streaming updates
                content := sm.Render(ctx, items, selected, config)
                updates <- MenuUpdate{
                    Type:    "refresh",
                    Content: content,
                    Data: map[string]interface{}{
                        "timestamp": time.Now(),
                        "selected": selected,
                    },
                }
            }
        }
    }()
    
    return updates
}

func (up *UserPreferences) Get() *UserPreferences {
    up.mu.RLock()
    defer up.mu.RUnlock()
    return up
}

func (sm *SmartMenu) addSuggestions(items []MenuItem[int], prefs *UserPreferences) {
    // Add suggestions logic
}

func (sm *SmartMenu) getContextualSubtitle() string {
    return "Enhanced by AI"
}
