package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"bilibili-backend/service"
	"bilibili-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 升级配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 开发环境允许所有 origin，生产环境应限制
		return true
	},
}

// DanmakuController 弹幕控制器
type DanmakuController struct {
	danmakuService *service.DanmakuService
}

// NewDanmakuController 创建弹幕控制器
func NewDanmakuController(danmakuService *service.DanmakuService) *DanmakuController {
	return &DanmakuController{danmakuService: danmakuService}
}

// WebSocket 升级为 WebSocket 连接，按 video_id 分组
func (ctrl *DanmakuController) WebSocket(c *gin.Context) {
	videoIDStr := c.Query("video_id")
	if videoIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_id 参数不能为空"})
		return
	}
	videoID, err := strconv.ParseInt(videoIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "video_id 格式错误"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] 升级失败: %v", err)
		return
	}

	// 注册到 Hub
	utils.DanmakuHubInstance.Register(videoID, conn)

	// 设置心跳处理
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	// 启动读消息的 goroutine（简单弹幕系统不需要处理客户端主动发送的消息，但保持连接活跃）
	go func() {
		defer utils.DanmakuHubInstance.Unregister(videoID, conn)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[WS] 读取消息错误: %v", err)
				}
				break
			}
		}
	}()
}

// GetDanmaku 获取视频某时间段的弹幕
func (ctrl *DanmakuController) GetDanmaku(c *gin.Context) {
	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	startStr := c.DefaultQuery("start", "0")
	endStr := c.DefaultQuery("end", "0")
	start, _ := strconv.Atoi(startStr)
	end, _ := strconv.Atoi(endStr)

	if end == 0 {
		end = start + 1
	}

	list, err := ctrl.danmakuService.GetDanmakuByTimeRange(videoID, start, end)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"list": list,
	})
}

// SendDanmakuRequest 发送弹幕请求体
type SendDanmakuRequest struct {
	Content   string `json:"content" binding:"required,max=100"`
	TimePoint int    `json:"time_point" binding:"gte=0"`
	Color     string `json:"color"`
	Type      int8   `json:"type"`
}

// SendDanmaku 发送弹幕
func (ctrl *DanmakuController) SendDanmaku(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint64)

	videoIDStr := c.Param("id")
	videoID, err := strconv.ParseUint(videoIDStr, 10, 64)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	var req SendDanmakuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	// 默认颜色和类型
	if req.Color == "" {
		req.Color = "#ffffff"
	}
	if req.Type == 0 {
		req.Type = 1
	}

	danmaku, err := ctrl.danmakuService.CreateDanmaku(videoID, uid, req.Content, req.TimePoint, req.Color, req.Type)
	if err != nil {
		utils.Error(c, utils.CodeBadRequest)
		return
	}

	utils.OK(c, gin.H{
		"danmaku": danmaku,
	})
}
