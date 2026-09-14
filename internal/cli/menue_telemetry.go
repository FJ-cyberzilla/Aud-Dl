package cli

import (
    "context"
    "log/slog"
    "runtime"
    "sync"
    "time"

    "audio-command-center/internal/cache"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
    )

// TelemetryMenu wraps menu with observability
type TelemetryMenu struct {
    *MenuTemplate
    logger    *slog.Logger
    tracer    trace.Tracer
    metrics   *TelemetryMetrics
}

func NewTelemetryMenu(
    style *InteractiveStyle,
    cache *cache.CacheManager,
    logger *slog.Logger,
) *TelemetryMenu {
    tracer := otel.Tracer("github.com/yourproject/cli")
    
    return &TelemetryMenu{
        MenuTemplate: NewMenuTemplate(style, cache),
        logger:       logger,
        tracer:       tracer,
        metrics:      &TelemetryMetrics{},
    }
}

// RenderWithTelemetry instruments menu rendering
func (tm *TelemetryMenu) RenderWithTelemetry[T comparable](
    ctx context.Context,
    items []MenuItem[T],
    selected int,
    config MenuConfig,
    opts ...RenderOption,
) string {
    ctx, span := tm.tracer.Start(ctx, "menu.render")
    defer span.End()
    
    span.SetAttributes(
        attribute.Int("items.count", len(items)),
        attribute.Int("selected.index", selected),
        attribute.String("menu.title", config.Title),
    )
    
    start := time.Now()
    result := tm.Render(ctx, items, selected, config, opts...)
    duration := time.Since(start)
    
    // Record metrics
    tm.metrics.record(duration, len(items))
    
    // Log if slow
    if duration > 100*time.Millisecond {
        tm.logger.Warn("Slow menu render",
            slog.Duration("duration", duration),
            slog.Int("items", len(items)),
            slog.String("title", config.Title),
        )
    }
    
    span.SetAttributes(
        attribute.Float64("render.duration_ms", float64(duration.Milliseconds())),
        attribute.Int("result.length", len(result)),
    )
    
    return result
}

// Expose metrics for monitoring
func (tm *TelemetryMenu) GetMetrics() map[string]interface{} {
    tm.metrics.mu.RLock()
    defer tm.metrics.mu.RUnlock()
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return map[string]interface{}{
        "render_count": tm.metrics.renderCount,
        "avg_render_ms": tm.metrics.avgRenderMs,
        "cache_hit_rate": tm.metrics.cacheHitRate,
        "memory_alloc": m.Alloc,
        "goroutines": runtime.NumGoroutine(),
    }
}

type TelemetryMetrics struct {
    mu            sync.RWMutex
    renderCount   int64
    totalDuration time.Duration
    avgRenderMs   float64
    cacheHitRate  float64
}

func (m *TelemetryMetrics) record(duration time.Duration, itemCount int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    m.renderCount++
    m.totalDuration += duration
    m.avgRenderMs = float64(m.totalDuration.Milliseconds()) / float64(m.renderCount)
}
