package cli

import (
	"context"
	"audio-command-center/internal/cache"
	"audio-command-center/internal/ui"
	"github.com/manifoldco/promptui"
)

type InteractiveMenu struct {
	template *MenuTemplate
	style    *InteractiveStyle
	theme    ui.Theme
}

func NewInteractiveMenu(theme ui.Theme, cache *cache.CacheManager) *InteractiveMenu {
	style := NewInteractiveStyle(theme)
	template := NewMenuTemplate(style, cache)
	return &InteractiveMenu{
		template: template,
		style:    style,
		theme:    theme,
	}
}

func (im *InteractiveMenu) PromptMainMenu(ctx context.Context) (int, error) {
	prompt := promptui.Select{
		Label: "Main Menu",
		Items: []string{
			"🔍 Search & Download",
			"📊 View Vault Stats",
			"⚙️  Settings",
			"🚪 Exit",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}
	return index, nil
}

func (im *InteractiveMenu) PromptFormatSelection(ctx context.Context) (string, error) {
	prompt := promptui.Select{
		Label: "Select Download Format",
		Items: []string{"MP3 (V0)", "FLAC (Lossless)", "WAV"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (im *InteractiveMenu) PromptSearchQuery(ctx context.Context) (string, error) {
	prompt := promptui.Prompt{
		Label: "Search Query",
	}

	return prompt.Run()
}
