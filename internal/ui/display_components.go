package ui

import "context"

type DisplayCompiler struct {
	Theme Theme
}

func NewDisplayCompiler(theme Theme) *DisplayCompiler {
	return &DisplayCompiler{Theme: theme}
}

type DisplayAssist struct {
	Theme Theme
}

func NewDisplayAssist(theme Theme) *DisplayAssist {
	return &DisplayAssist{Theme: theme}
}

type HealthQuest struct {
	Theme Theme
}

func NewHealthQuest(theme Theme) *HealthQuest {
	return &HealthQuest{Theme: theme}
}

func (hq *HealthQuest) RunStartupQuest(ctx context.Context, vaultDir string) bool {
	return true
}
