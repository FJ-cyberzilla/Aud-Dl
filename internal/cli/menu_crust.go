package cli

import (
	"context"
	"time"

	"audio-command-center/internal/cache"
)

type MenuCrust struct {
	Dispatcher *StyleDispatcher
	renderer   *JunctionRenderer
	admin      *RendererAdmin
	junction   *MenuJunction
}

func NewMenuCrust(theme ThemeMode) *MenuCrust {
	dispatcher := NewStyleDispatcher(theme)
	template := NewJunctionTemplate()
	cacheInstance := cache.NewCacheManager(10 * time.Minute)
	
	return &MenuCrust{
		Dispatcher: dispatcher,
		renderer:   NewJunctionRenderer(dispatcher, template),
		admin:      NewRendererAdmin(dispatcher),
		junction:   NewMenuJunction(dispatcher.theme, cacheInstance),
	}
}

// BootSkinFacade initializes and starts the unified CLI presentation skin
func (mc *MenuCrust) BootSkinFacade(ctx context.Context) error {
	// Delegates directly to the underlying event junction loop
	return mc.junction.StartInteractiveSession(ctx)
}

func (mc *MenuCrust) SetSearchHandler(h SearchExecutionHandler) {
	mc.junction.onSearch = h
}

func (mc *MenuCrust) SetStatsProvider(p AdminStatsProvider) {
	mc.junction.getStats = p
}

// DisplayAdminView renders the administrative stats panel instantly through the facade
func (mc *MenuCrust) DisplayAdminView(stats AdminDashboardStats) {
	mc.admin.RenderAdminDashboard(stats)
	mc.admin.RenderPrompt()
}
