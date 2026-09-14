package cli

import (
	"fmt"
	"audio-command-center/internal/ui"
)

type InteractiveStyle struct {
	Theme ui.Theme
}

func NewInteractiveStyle(theme ui.Theme) *InteractiveStyle {
	return &InteractiveStyle{Theme: theme}
}

func (s *InteractiveStyle) ActiveOption(text string) string {
	return fmt.Sprintf("%s ❯ %s%s", ui.TrueColor("", s.Theme.Accent), ui.TrueColor(text, s.Theme.Primary), "\033[0m")
}

func (s *InteractiveStyle) InactiveOption(text string) string {
	return fmt.Sprintf("   %s", ui.TrueColor(text, s.Theme.Subtext))
}

func (s *InteractiveStyle) Title(text string) string {
	return ui.TrueColor(text, s.Theme.Highlight)
}

func (s *InteractiveStyle) HintText(text string) string {
	return ui.TrueColor(text, s.Theme.Subtext)
}

func (s *InteractiveStyle) WarningMessage(text string) string {
	return ui.TrueColor(text, s.Theme.Warning)
}

func (s *InteractiveStyle) ErrorMessage(text string) string {
	return ui.TrueColor(text, s.Theme.Error)
}

func (s *InteractiveStyle) LoadingAnimation(text string) {
	fmt.Printf("%s\n", ui.TrueColor(text, s.Theme.Accent))
}

func (s *InteractiveStyle) HeaderBanner(title string) string {
	border := "────────────────────────────────────────────────────────────"
	return fmt.Sprintf("%s\n  %s\n%s", 
		ui.TrueColor(border, s.Theme.Secondary), 
		ui.TrueColor(title, s.Theme.Highlight), 
		ui.TrueColor(border, s.Theme.Secondary),
	)
}

func (s *InteractiveStyle) PromptLabel(label string) string {
	return fmt.Sprintf("%s %s: ", ui.TrueColor("?", s.Theme.Warning), ui.TrueColor(label, s.Theme.Text))
}
