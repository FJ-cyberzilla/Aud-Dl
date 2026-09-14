package ui

import (
	"fmt"
)

type RefreshManager struct{}

func NewRefreshManager() *RefreshManager { return &RefreshManager{} }

func (rm *RefreshManager) ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
