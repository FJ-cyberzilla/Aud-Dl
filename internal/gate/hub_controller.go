package gate

import (
	"context"
	"fmt"
	"io"
	"time"
)

type GateHubController struct {
	validator  *HubValidator
	shield     *SystemShield
	egress     *EgressTLSClient
	updater    *AutoUpdateController
	tmpLog     *TmpLogManager
	Sanctioner *HubSanctioner
}

func NewGateHubController(appVersion, configDir string) *GateHubController {
	return &GateHubController{
		validator:  NewHubValidator(),
		shield:     NewSystemShield(),
		egress:     NewEgressTLSClient(),
		updater:    NewAutoUpdateController(appVersion, configDir),
		tmpLog:     NewTmpLogManager(configDir),
		Sanctioner: NewHubSanctioner(3, 15*time.Minute),
	}
}

// GetSanctioner provides access to the sanctioning manager.
func (ghc *GateHubController) GetSanctioner() *HubSanctioner {
	return ghc.Sanctioner
}

// IngestAndSanitizeStream validates request parameters first, then ingests stream safely
func (ghc *GateHubController) IngestAndSanitizeStream(ctx context.Context, query, filename string) (io.Reader, error) {
	// 1. Pre-Ingest Input & Path Validation
	if err := ghc.validator.ValidateQuery(query); err != nil {
		return nil, err
	}

	safeName, err := ghc.validator.ValidateSanitizeFilename(filename)
	if err != nil {
		return nil, err
	}

	// 2. Fetch via Anti-Bot Egress TLS
	rawStream, err := ghc.egress.FetchAudioStream(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("egress fetch failed: %w", err)
	}

	// 3. Post-Fetch Magic Byte Payload Shielding
	safeStream, err := ghc.shield.InspectStream(rawStream, safeName)
	if err != nil {
		return nil, fmt.Errorf("security shield rejection: %w", err)
	}

	return safeStream, nil
}

// ExecuteMaintenanceCycle manages startup checks and 7-day auto-update cycles.
func (ghc *GateHubController) ExecuteMaintenanceCycle(ctx context.Context) (bool, string, error) {
	_ = ghc.tmpLog.RecordHealthQuestCheck()
	return ghc.updater.CheckAndApplyWeeklyUpdate(ctx, ghc.tmpLog)
}
