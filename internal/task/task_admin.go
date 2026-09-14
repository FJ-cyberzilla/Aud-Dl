package task

import (
	"context"
	"fmt"
	"os"
	"time"

	"audio-command-center/internal/cache"
	"audio-command-center/internal/ui"
)

type TaskAdmin struct {
	AppVersion string
	ConfigDir  string
	VaultDir   string
	Theme      ui.Theme
	Bindings   *SystemBindings
	Controller *TaskController
}

func NewTaskAdmin(version, configDir, vaultDir string, theme ui.Theme) (*TaskAdmin, error) {
	// 1. Delegate system initialization and wiring to TaskBinder
	binder := NewTaskBinder(version, configDir, vaultDir, theme)
	bindings, err := binder.BindSystemServices()
	if err != nil {
		return nil, fmt.Errorf("system binding failed: %w", err)
	}

	// 2. Wire Division Worker and Interactive Controller
	division := NewTaskDivision(bindings, 100) // 100 is an arbitrary queue size for now

	controller := NewTaskController(bindings, division, theme, cache.NewCacheManager(10*time.Minute))

	return &TaskAdmin{
		AppVersion: version,
		ConfigDir:  configDir,
		VaultDir:   vaultDir,
		Theme:      theme,
		Bindings:   bindings,
		Controller: controller,
	}, nil
}

// BootstrapSystem handles system startup checks and weekly updates
func (ta *TaskAdmin) BootstrapSystem(ctx context.Context) bool {
	// Startup Health Quest
	health := ui.NewHealthQuest(ta.Theme)
	if !health.RunStartupQuest(ctx, ta.VaultDir) {
		fmt.Println(ui.TrueColor("\n⚠️ Health warnings detected. Proceeding with caution...\n", ta.Theme.Warning))
	}

	// Weekly Auto-Update via GateHub
	updated, newVer, err := ta.Bindings.GateHub.ExecuteMaintenanceCycle(ctx)
	if err == nil && updated {
		fmt.Printf(ui.TrueColor("\n🚀 Engine updated to %s! Please restart app.\n", ta.Theme.Success), newVer)
		os.Exit(0)
	}

	return true
}

func (ta *TaskAdmin) LaunchApp(ctx context.Context) {
	ta.Controller.RunInteractiveSession(ctx)
}

func (ta *TaskAdmin) GetController() *TaskController {
	return ta.Controller
}
