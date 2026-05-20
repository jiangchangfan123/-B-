package utils

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// DanmakuHub 管理所有弹幕 WebSocket 连接，按 video_id 分组广播
type DanmakuHub struct {
	// video_id -> connections
	rooms map[int64]map[*websocket.Conn]bool
	mu    sync.Mutex
}

// DanmakuHubInstance 全局实例
var DanmakuHubInstance = NewDanmakuHub()

// NewDanmakuHub 创建新的 Hub
func NewDanmakuHub() *DanmakuHub {
	return &DanmakuHub{
		rooms: make(map[int64]map[*websocket.Conn]bool),
	}
}

// Register 将连接注册到指定 video room
func (h *DanmakuHub) Register(videoID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[videoID] == nil {
		h.rooms[videoID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[videoID][conn] = true
	log.Printf("[WS] 连接注册到 video_id=%d，当前连接数: %d", videoID, len(h.rooms[videoID]))
}

// Unregister 从指定 video room 移除连接，并关闭连接
func (h *DanmakuHub) Unregister(videoID int64, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.rooms[videoID]; ok {
		delete(conns, conn)
		log.Printf("[WS] 连接从 video_id=%d 移除，剩余连接数: %d", videoID, len(conns))
		if len(conns) == 0 {
			delete(h.rooms, videoID)
		}
	}
	conn.Close()
}

// Broadcast 向指定 video room 的所有连接广播文本消息
// 如果某个连接写失败，会自动清理该死亡连接
func (h *DanmakuHub) Broadcast(videoID int64, message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	conns, ok := h.rooms[videoID]
	if !ok || len(conns) == 0 {
		return
	}

	var deadConns []*websocket.Conn
	for conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("[WS] 广播写入失败: %v", err)
			deadConns = append(deadConns, conn)
		}
	}

	// 清理失败连接
	for _, conn := range deadConns {
		delete(conns, conn)
		conn.Close()
	}
	if len(conns) == 0 {
		delete(h.rooms, videoID)
	}
}

// StartHeartbeat 启动定时心跳检测
// 每 interval 时间向所有连接发送 Ping，超时未响应则清理
func (h *DanmakuHub) StartHeartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		h.mu.Lock()
		for videoID, conns := range h.rooms {
			var deadConns []*websocket.Conn
			for conn := range conns {
				// WriteControl 是线程安全的，但此处已持有写锁，所以直接用
				if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
					log.Printf("[WS] Ping 失败，清理连接: %v", err)
					deadConns = append(deadConns, conn)
				}
			}
			for _, conn := range deadConns {
				delete(conns, conn)
				conn.Close()
			}
			if len(conns) == 0 {
				delete(h.rooms, videoID)
			}
		}
		h.mu.Unlock()
	}
}
