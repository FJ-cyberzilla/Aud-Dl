package cli

import (
    "context"
    "fmt"
    "io"
    "strings"
    "sync"
    "time"
    
    "audio-command-center/internal/ui"
    "audio-command-center/internal/cache"
)
// MenuEvent represents async menu updates
type MenuEvent struct {
    Type    EventType
    Payload any
    Timestamp time.Time
}

type MenuItem[T comparable] struct {
    ID          T
    Label       string
    Icon        string
    Description string
    Disabled    bool
}

type MenuConfig struct {
    Title      string
    Subtitle   string
    ShowHelp   bool
    Paginate   bool
    PageSize   int
    WrapAround bool
}

type EventType int

const (
    EventSelection EventType = iota
    EventFilter
    EventRefresh
    EventError
    EventLoading
)

// MenuTemplate is a high-performance, thread-safe menu renderer
type MenuTemplate struct {
    style     *InteractiveStyle
    cache     *cache.CacheManager
    mu        sync.RWMutex
    metrics   *MenuMetrics
    eventCh   chan MenuEvent
    ctx       context.Context
    cancel    context.CancelFunc
}

// MenuMetrics tracks performance and usage
type MenuMetrics struct {
    RenderCount    int64
    AvgRenderTime  time.Duration
    CacheHitRate   float64
    ActiveSessions int32
    mu             sync.RWMutex
}

func NewMenuTemplate(style *InteractiveStyle, cache *cache.CacheManager) *MenuTemplate {
    ctx, cancel := context.WithCancel(context.Background())
    mt := &MenuTemplate{
        style:   style,
        cache:   cache,
        eventCh: make(chan MenuEvent, 100),
        ctx:     ctx,
        cancel:  cancel,
        metrics: &MenuMetrics{},
    }
    
    // Start event processor
    go mt.processEvents()
    
    return mt
}

func (mt *MenuTemplate) recordMetrics(start time.Time) {
    mt.metrics.mu.Lock()
    defer mt.metrics.mu.Unlock()
    mt.metrics.RenderCount++
    mt.metrics.AvgRenderTime = (mt.metrics.AvgRenderTime + time.Since(start)) / 2
}

func (mt *MenuTemplate) generateCacheKey[T comparable](items []MenuItem[T], selected int, config MenuConfig) string {
    return fmt.Sprintf("%v-%d-%v", items, selected, config)
}

func (mt *MenuTemplate) estimateSize[T comparable](items []MenuItem[T], config MenuConfig) int {
    return len(items) * 100
}

func (mt *MenuTemplate) renderHeaderAsync(ctx context.Context, config MenuConfig, out chan<- string) {
    out <- mt.style.Title(config.Title)
}

func (mt *MenuTemplate) renderFooterAsync(ctx context.Context, config MenuConfig, out chan<- string) {
    out <- mt.style.HintText("Press arrow keys to navigate")
}

func (mt *MenuTemplate) isVisible(index int, selected int, config MenuConfig) bool {
    return true
}

func (mt *MenuTemplate) emitMetrics() {
    mt.metrics.mu.Lock()
    defer mt.metrics.mu.Unlock()
}

// Render renders menu with intelligent caching and streaming
func (mt *MenuTemplate) Render[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    opts ...RenderOption,
) string {
    start := time.Now()
    defer mt.recordMetrics(start)

    // Check cache for identical render
    cacheKey := mt.generateCacheKey(items, selected, config)
    if cached, ok := mt.cache.Get(cacheKey); ok {
        mt.metrics.mu.Lock()
        mt.metrics.CacheHitRate = (mt.metrics.CacheHitRate*0.9 + 0.1)
        mt.metrics.mu.Unlock()
        return cached.(string)
    }

    // Build with streaming approach
    var b strings.Builder
    b.Grow(mt.estimateSize(items, config))

    // Async header rendering
    headerCh := make(chan string, 1)
    go mt.renderHeaderAsync(ctx, config, headerCh)

    // Render items in parallel for large menus
    itemsCh := make(chan string, len(items))
    go mt.renderItemsParallel(ctx, items, selected, config, itemsCh)

    // Footer with progress
    footerCh := make(chan string, 1)
    go mt.renderFooterAsync(ctx, config, footerCh)

    // Assemble with timeout protection
    timeout := time.After(100 * time.Millisecond)
    
    select {
    case header := <-headerCh:
        b.WriteString(header)
        b.WriteString("\n\n")
    case <-timeout:
        b.WriteString(mt.style.WarningMessage("⏳ Loading header..."))
    }

    // Stream items as they arrive
    itemCount := 0
    for itemHTML := range itemsCh {
        b.WriteString(itemHTML)
        b.WriteString("\n")
        itemCount++
        
        // Yield to scheduler for large menus
        if itemCount%50 == 0 {
            time.Sleep(time.Microsecond)
        }
    }

    select {
    case footer := <-footerCh:
        b.WriteString("\n")
        b.WriteString(footer)
    case <-timeout:
        b.WriteString("\n" + mt.style.HintText("Press any key to continue..."))
    }

    result := b.String()
    
    // Cache the result
    mt.cache.Set(cacheKey, result, 5*time.Minute)
    
    return result
}

// renderItemsParallel renders items using a worker pool
func (mt *MenuTemplate) renderItemsParallel[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    out chan<- string,
) {
    defer close(out)
    
    const workerCount = 4
    type job struct {
        index int
        item  MenuItem[T]
    }
    
    jobs := make(chan job, len(items))
    results := make(chan string, len(items))
    
    // Start workers
    var wg sync.WaitGroup
    for w := 0; w < workerCount; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case <-ctx.Done():
                    return
                default:
                    results <- mt.renderItem(job.item, job.index, selected)
                }
            }
        }()
    }
    
    // Feed jobs
    go func() {
        for i, item := range items {
            if !config.Paginate || mt.isVisible(i, selected, config) {
                jobs <- job{index: i, item: item}
            }
        }
        close(jobs)
    }()
    
    // Wait for completion
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results in order
    rendered := make([]string, len(items))
    for result := range results {
        // Parse index from result (simplified)
        // In production, use a more robust ordering mechanism
        rendered = append(rendered, result)
    }
    
    for _, r := range rendered {
        select {
        case <-ctx.Done():
            return
        case out <- r:
        }
    }
}

// renderItem renders a single menu item
func (mt *MenuTemplate) renderItem[T comparable](
    item MenuItem[T],
    index int,
    selected int,
) string {
    if item.Disabled {
        return mt.style.InactiveOption("⛔ " + item.Label)
    }
    
    if index == selected {
        var parts []string
        parts = append(parts, mt.style.ActiveOption("▸ "+item.Icon+" "+item.Label))
        if item.Description != "" {
            parts = append(parts, "  "+mt.style.HintText(item.Description))
        }
        return strings.Join(parts, "\n")
    }
    
    return mt.style.InactiveOption("  " + item.Icon + " " + item.Label)
}

// processEvents handles async menu events
func (mt *MenuTemplate) processEvents() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case event := <-mt.eventCh:
            mt.handleEvent(event)
        case <-ticker.C:
            mt.emitMetrics()
        case <-mt.ctx.Done():
            return
        }
    }
}

func (mt *MenuTemplate) handleEvent(event MenuEvent) {
    switch event.Type {
    case EventLoading:
        mt.style.LoadingAnimation(event.Payload.(string))
    case EventError:
        mt.style.ErrorMessage(event.Payload.(string))
    case EventRefresh:
        mt.cache.Flush()
    }
}

// Adaptive rendering based on terminal capabilities
func (mt *MenuTemplate) RenderAdaptive[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    w io.Writer,
) error {
    // Detect terminal capabilities
    termInfo := ui.DetectTerminal()
    
    if termInfo.IsTTY && termInfo.Width > 80 {
        return mt.renderRich(ctx, items, selected, config, w)
    }
    
    if termInfo.IsTTY && termInfo.Width < 40 {
        return mt.renderCompact(ctx, items, selected, config, w)
    }
    
    return mt.renderPlain(items, selected, config, w)
}

// renderRich provides full-featured rendering
func (mt *MenuTemplate) renderRich[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    w io.Writer,
) error {
    // Use table-based rendering with borders
    table := ui.NewTable(mt.style.Theme, ui.Column{Title: "", Width: 5}, ui.Column{Title: "Option", Width: 20}, ui.Column{Title: "Description", Width: 30})

    for i, item := range items {
        icon := item.Icon
        if i == selected {
            icon = "▸ " + icon
        }
        table.AddRow(icon, item.Label, item.Description)
    }

    return table.Render(w)
    }
// renderCompact for small terminals
func (mt *MenuTemplate) renderCompact[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    w io.Writer,
) error {
    for i, item := range items {
        if i == selected {
            fmt.Fprintf(w, "> %s\n", item.Label)
        } else {
            fmt.Fprintf(w, "  %s\n", item.Label)
        }
    }
    return nil
}

// renderPlain for non-TTY output
func (mt *MenuTemplate) renderPlain[T comparable](
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    w io.Writer,
) error {
    for i, item := range items {
        fmt.Fprintf(w, "%d. %s\n", i+1, item.Label)
    }
    return nil
}
