package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketManager WebSocket管理器
type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan interface{}
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

// NewWebSocketManager 创建新的WebSocket管理器
func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan interface{}, 100),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Start 启动WebSocket管理器
func (wm *WebSocketManager) Start() {
	for {
		select {
		case client := <-wm.register:
			wm.mu.Lock()
			wm.clients[client] = true
			wm.mu.Unlock()
			log.Printf("WebSocket客户端已连接，当前连接数: %d", len(wm.clients))

		case client := <-wm.unregister:
			wm.mu.Lock()
			delete(wm.clients, client)
			wm.mu.Unlock()
			client.Close()
			log.Printf("WebSocket客户端已断开，当前连接数: %d", len(wm.clients))

		case message := <-wm.broadcast:
			wm.mu.RLock()
			for client := range wm.clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Printf("发送WebSocket消息失败: %v", err)
					client.Close()
					delete(wm.clients, client)
				}
			}
			wm.mu.RUnlock()
		}
	}
}

// Broadcast 广播消息
func (wm *WebSocketManager) Broadcast(message interface{}) {
	wm.broadcast <- message
}

// HandleWebSocket WebSocket处理器
func (wm *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	wm.register <- conn

	// 启动goroutine处理客户端消息
	go func() {
		defer func() {
			wm.unregister <- conn
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				break
			}

			// 处理客户端消息（如果需要）
			log.Printf("收到WebSocket消息: %s", string(message))
		}
	}()
}

// BroadcastHardwareInfo 广播硬件信息
func (wm *WebSocketManager) BroadcastHardwareInfo(data interface{}) {
	message := map[string]interface{}{
		"type": "hardware_info",
		"data": data,
		"time": time.Now().Unix(),
	}
	wm.Broadcast(message)
}

// BroadcastHostList 广播主机列表
func (wm *WebSocketManager) BroadcastHostList(hosts map[string]*HostData) {
	message := map[string]interface{}{
		"type": "host_list",
		"data": hosts,
		"time": time.Now().Unix(),
	}
	wm.Broadcast(message)
}
