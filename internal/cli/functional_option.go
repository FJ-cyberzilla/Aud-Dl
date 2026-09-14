package cli

// RenderOption functional options pattern
type RenderOption func(*renderContext)

type renderContext struct {
    maxWidth    int
    showNumbers bool
    compact     bool
    colorize    bool
}

// WithMaxWidth sets terminal width for wrapping
func WithMaxWidth(width int) RenderOption {
    return func(c *renderContext) {
        c.maxWidth = width
    }
}

// WithNumbers shows option numbers
func WithNumbers() RenderOption {
    return func(c *renderContext) {
        c.showNumbers = true
    }
}

// CompactMode removes extra spacing
func CompactMode() RenderOption {
    return func(c *renderContext) {
        c.compact = true
    }
}

// NoColor disables ANSI codes
func NoColor() RenderOption {
    return func(c *renderContext) {
        c.colorize = false
    }
}
