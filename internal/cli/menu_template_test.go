package cli

import (
    "context"
    "strings"
    "testing"
    "time"

    "audio-command-center/internal/ui"
    "audio-command-center/internal/cache"
)

func TestMenuTemplate_Render(t *testing.T) {
    style := NewInteractiveStyle(ui.CyberDarkTheme)
    cm := cache.NewCacheManager(time.Minute)
    mt := NewMenuTemplate(style, cm)

    tests := []struct {
        name     string
        items    []MenuItem[int]
        selected int
        config   MenuConfig
        contains []string // substring checks
    }{
        {
            name:     "renders basic menu",
            items:    MainMenuItems,
            selected: 0,
            config:   MenuConfig{Title: "TEST MENU", ShowHelp: true},
            contains: []string{"TEST MENU", "🔥", "navigate"},
        },
        {
            name:     "handles invalid selection",
            items:    MainMenuItems,
            selected: 99,
            config:   MenuConfig{Title: "TEST", WrapAround: true},
            contains: []string{"Most Wanted"},
        },
        {
            name: "respects disabled items",
            items: []MenuItem[int]{
                {ID: 1, Label: "Enabled", Icon: "✅"},
                {ID: 2, Label: "Disabled", Icon: "❌", Disabled: true},
            },
            selected: 1,
            config:   MenuConfig{Title: "TEST"},
            contains: []string{"⛔", "Disabled"},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := mt.Render(context.Background(), tt.items, tt.selected, tt.config)
            
            for _, substr := range tt.contains {
                if !strings.Contains(result, substr) {
                    t.Errorf("expected %q in rendered output", substr)
                }
            }
        })
    }
}

func BenchmarkMenuTemplate_Render(b *testing.B) {
    style := NewInteractiveStyle(ui.CyberDarkTheme)
    cm := cache.NewCacheManager(time.Minute)
    mt := NewMenuTemplate(style, cm)
    items := make([]MenuItem[int], 100)
    for i := range items {
        items[i] = MenuItem[int]{
            ID:    i,
            Label: "Item " + string(rune('A'+i%26)),
            Icon:  "📦",
        }
    }
    config := MenuConfig{Title: "BENCHMARK", Paginate: true, PageSize: 10}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = mt.Render(context.Background(), items, i%100, config)
    }
}
