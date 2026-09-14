package cli

import (
    "context"
    "encoding/json"
    "net/http"
    "sync"

    "audio-command-center/internal/cache"
    "github.com/gorilla/websocket"
    )

// WebMenu serves menus over WebSocket for remote control
type WebMenu struct {
    *MenuTemplate
    clients sync.Map
    upgrader websocket.Upgrader
}

func NewWebMenu(style *InteractiveStyle, cache *cache.CacheManager) *WebMenu {
    return &WebMenu{
        MenuTemplate: NewMenuTemplate(style, cache),
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool { return true },
            ReadBufferSize: 1024,
            WriteBufferSize: 1024,
        },
    }
}

// ServeWS handles WebSocket connections
func (wm *WebMenu) ServeWS(w http.ResponseWriter, r *http.Request) {
    conn, err := wm.upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
        return
    }
    defer conn.Close()
    
    ctx := context.Background()
    clientID := r.RemoteAddr
    wm.clients.Store(clientID, conn)
    defer wm.clients.Delete(clientID)
    
    for {
        // Read client messages
        var msg struct {
            Type    string          `json:"type"`
            Action  string          `json:"action"`
            Payload json.RawMessage `json:"payload"`
        }
        
        err := conn.ReadJSON(&msg)
        if err != nil {
            break
        }
        
        // Handle menu actions
        switch msg.Action {
        case "render":
            wm.handleRender(ctx, conn, msg.Payload)
        case "select":
            wm.handleSelect(msg.Payload)
        case "search":
            wm.handleSearch(msg.Payload)
        }
    }
}

func (wm *WebMenu) handleRender(ctx context.Context, conn *websocket.Conn, payload json.RawMessage) {
    // Placeholder implementation
}

func (wm *WebMenu) handleSelect(payload json.RawMessage) {
    // Placeholder implementation
}

func (wm *WebMenu) handleSearch(payload json.RawMessage) {
    // Placeholder implementation
}
// Broadcast updates to all connected clients
func (wm *WebMenu) Broadcast(update interface{}) {
    wm.clients.Range(func(key, value interface{}) bool {
        conn := value.(*websocket.Conn)
        if err := conn.WriteJSON(update); err != nil {
            wm.clients.Delete(key)
        }
        return true
    })
}
