package gate

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type AutoUpdateController struct {
	AppVersion string
	ConfigDir  string
}

func NewAutoUpdateController(appVersion, configDir string) *AutoUpdateController {
	return &AutoUpdateController{
		AppVersion: appVersion,
		ConfigDir:  configDir,
	}
}

type VersionPayload struct {
	LatestVersion string `json:"latest_version"`
}

// CheckAndApplyWeeklyUpdate checks for updates and returns (updated, version, error)
func (auc *AutoUpdateController) CheckAndApplyWeeklyUpdate(ctx context.Context, tmpLog *TmpLogManager) (bool, string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://raw.githubusercontent.com/username/repo/main/version.json", nil)
	if err != nil {
		return false, auc.AppVersion, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, auc.AppVersion, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, auc.AppVersion, nil
	}

	var payload VersionPayload
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return false, auc.AppVersion, err
	}

	if payload.LatestVersion != auc.AppVersion && payload.LatestVersion != "" {
		return true, payload.LatestVersion, nil
	}

	return false, auc.AppVersion, nil
}
