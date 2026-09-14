package task

import (
	"net/http"
	"time"

	"audio-command-center/internal/cache"
	"audio-command-center/internal/gate"
	"audio-command-center/internal/processor"
	"audio-command-center/internal/search"
	"audio-command-center/internal/ui"
	"audio-command-center/internal/vault"
)

type SystemBindings struct {
	GateHub       *gate.GateHubController
	Transcoder    *processor.MemoryPipedTranscoder
	Tagger        *processor.EliteTagger
	VaultRouter   *vault.SmartVaultRouter
	Detector      *vault.PersistentDetector
	VaultAdmin    *vault.VaultAdmin
	SearchManager *search.SearchManager
	Compiler      *ui.DisplayCompiler
	Assist        *ui.DisplayAssist
}

type TaskBinder struct {
	AppVersion string
	ConfigDir  string
	VaultDir   string
	Theme      ui.Theme
}

func NewTaskBinder(appVersion, configDir, vaultDir string, theme ui.Theme) *TaskBinder {
	return &TaskBinder{
		AppVersion: appVersion,
		ConfigDir:  configDir,
		VaultDir:   vaultDir,
		Theme:      theme,
	}
}

func (tb *TaskBinder) BindSystemServices() (*SystemBindings, error) {
	// Unified Gate Assistant
	gateHub := gate.NewGateHubController(tb.AppVersion, tb.ConfigDir)

	// Processor & Vault Services
	transcoder := processor.NewMemoryPipedTranscoder()
	tagger := processor.NewEliteTagger()
	router := vault.NewSmartVaultRouter(tb.VaultDir)
	compiler := ui.NewDisplayCompiler(tb.Theme)
	assist := ui.NewDisplayAssist(tb.Theme)
	searchManager := search.NewSearchManager(cache.NewCacheManager(10 * time.Minute))
	searchManager.AddPlatform(&search.YouTubeAudioPlatform{Client: &http.Client{}})
	searchManager.AddPlatform(&search.SoundcloudPlatform{Client: &http.Client{}})
	searchManager.AddPlatform(&search.BandcampPlatform{Client: &http.Client{}})

	detector := vault.NewPersistentDetector()
	vaultAdmin, err := vault.NewVaultAdmin(tb.VaultDir)
	if err != nil {
		return nil, err
	}

	return &SystemBindings{
		GateHub:       gateHub,
		Transcoder:    transcoder,
		Tagger:        tagger,
		VaultRouter:   router,
		Detector:      detector,
		VaultAdmin:    vaultAdmin,
		SearchManager: searchManager,
		Compiler:      compiler,
		Assist:        assist,
	}, nil
}
