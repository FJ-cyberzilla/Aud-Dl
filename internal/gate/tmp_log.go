package gate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type TmpLogManager struct {
	LogPath    string
	ResumePath string
	mu         sync.RWMutex
}

func NewTmpLogManager(configDir string) *TmpLogManager {
	return &TmpLogManager{
		LogPath:    filepath.Join(configDir, "runtime_logs.tmp"),
		ResumePath: filepath.Join(configDir, "resume_state.json"),
	}
}

// RecordHealthQuestCheck logs the timestamp of a successful startup health check
func (tlm *TmpLogManager) RecordHealthQuestCheck() error {
	timestamp := time.Now().Format(time.RFC3339)
	entry := fmt.Sprintf("Health Check: %s\n", timestamp)
	
	f, err := os.OpenFile(tlm.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

func (tlm *TmpLogManager) SaveResumeState(streamID string, state interface{}) error {
	tlm.mu.Lock()
	defer tlm.mu.Unlock()

	states := make(map[string]interface{})
	data, err := os.ReadFile(tlm.ResumePath)
	if err == nil {
		_ = json.Unmarshal(data, &states)
	}

	states[streamID] = state
	updated, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tlm.ResumePath, updated, 0644)
}

func (tlm *TmpLogManager) GetResumeState(streamID string) (json.RawMessage, bool) {
	tlm.mu.RLock()
	defer tlm.mu.RUnlock()

	data, err := os.ReadFile(tlm.ResumePath)
	if err != nil {
		return nil, false
	}

	var states map[string]json.RawMessage
	if err := json.Unmarshal(data, &states); err != nil {
		return nil, false
	}

	raw, exists := states[streamID]
	return raw, exists
}

func (tlm *TmpLogManager) ClearResumeState(streamID string) error {
	tlm.mu.Lock()
	defer tlm.mu.Unlock()

	data, err := os.ReadFile(tlm.ResumePath)
	if err != nil {
		return nil
	}

	var states map[string]interface{}
	if err := json.Unmarshal(data, &states); err != nil {
		return err
	}

	delete(states, streamID)
	updated, err := json.MarshalIndent(states, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tlm.ResumePath, updated, 0644)
}
