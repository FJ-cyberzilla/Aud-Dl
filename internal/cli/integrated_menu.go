package cli

import "context"

// SimpleMenu provides a fluent API for common menus
type SimpleMenu struct {
    template *MenuTemplate
    items    []MenuItem[int]
    config   MenuConfig
}

func NewSimpleMenu(template *MenuTemplate, items []MenuItem[int]) *SimpleMenu {
    return &SimpleMenu{
        template: template,
        items:    items,
        config: MenuConfig{
            Title:      "Menu",
            ShowHelp:   true,
            Paginate:   false,
            WrapAround: true,
        },
    }
}

func (m *SimpleMenu) WithTitle(title string) *SimpleMenu {
    m.config.Title = title
    return m
}

func (m *SimpleMenu) WithSubtitle(sub string) *SimpleMenu {
    m.config.Subtitle = sub
    return m
}

func (m *SimpleMenu) WithPagination(pageSize int) *SimpleMenu {
    m.config.Paginate = true
    m.config.PageSize = pageSize
    return m
}

func (m *SimpleMenu) WithHelp(show bool) *SimpleMenu {
    m.config.ShowHelp = show
    return m
}

func (m *SimpleMenu) Render(ctx context.Context, selected int) string {
    return m.template.Render(ctx, m.items, selected, m.config)
}

// Usage example:
// menu := NewSimpleMenu(template, MainMenuItems).
//     WithTitle("AUDIO COMMAND CENTER").
//     WithSubtitle("v2.0.0").
//     WithPagination(5)
// output := menu.Render(selectedIndex)
