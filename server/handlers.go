package server

import (
	"net/http"
	"strconv"
	"time"

	"fab50/types"

	"github.com/gin-gonic/gin"
)

// Handlers HTTP处理器
type Handlers struct {
	storage *Storage
	ws      *WebSocketManager
}

// NewHandlers 创建新的处理器
func NewHandlers(storage *Storage, ws *WebSocketManager) *Handlers {
	return &Handlers{
		storage: storage,
		ws:      ws,
	}
}

// SetupRoutes 设置路由
func (h *Handlers) SetupRoutes(r *gin.Engine) {
	// API路由
	api := r.Group("/api")
	{
		api.POST("/hardware", h.handleHardwareInfo)
		api.GET("/hosts", h.handleGetHosts)
		api.GET("/hosts/:hostname", h.handleGetHostData)
		api.GET("/latest", h.handleGetLatest)
	}

	// WebSocket路由
	r.GET("/ws", func(c *gin.Context) {
		h.ws.HandleWebSocket(c.Writer, c.Request)
	})

	// 静态文件路由
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// 页面路由
	r.GET("/", h.handleIndex)
	r.GET("/host/:hostname", h.handleHostDetail)
}

// handleHardwareInfo 处理硬件信息提交
func (h *Handlers) handleHardwareInfo(c *gin.Context) {
	var info types.HardwareInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的JSON数据"})
		return
	}

	// 添加时间戳
	if info.Timestamp.IsZero() {
		info.Timestamp = time.Now()
	}

	// 存储数据
	h.storage.AddHardwareInfo(&info)

	// 广播最新数据
	latest := h.storage.GetLatestHardwareInfo()
	h.ws.BroadcastHardwareInfo(latest)

	// 通知客户端停止并自删除
	c.JSON(http.StatusOK, gin.H{
		"message":   "数据接收成功",
		"action":    "stop_and_delete",
		"timestamp": time.Now().Unix(),
	})
}

// handleGetHosts 获取主机列表
func (h *Handlers) handleGetHosts(c *gin.Context) {
	hosts := h.storage.GetHosts()
	c.JSON(http.StatusOK, hosts)
}

// handleGetHostData 获取指定主机的数据
func (h *Handlers) handleGetHostData(c *gin.Context) {
	hostname := c.Param("hostname")

	// 获取查询参数
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 100
	}

	hostData, exists := h.storage.GetHostData(hostname)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "主机不存在"})
		return
	}

	// 限制返回的记录数
	if len(hostData.HardwareInfo) > limit {
		hostData.HardwareInfo = hostData.HardwareInfo[len(hostData.HardwareInfo)-limit:]
	}

	c.JSON(http.StatusOK, hostData)
}

// handleGetLatest 获取最新的硬件信息
func (h *Handlers) handleGetLatest(c *gin.Context) {
	latest := h.storage.GetLatestHardwareInfo()
	c.JSON(http.StatusOK, latest)
}

// handleIndex 处理首页
func (h *Handlers) handleIndex(c *gin.Context) {
	hosts := h.storage.GetHosts()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "硬件监控系统",
		"hosts": hosts,
	})
}

// handleHostDetail 处理主机详情页
func (h *Handlers) handleHostDetail(c *gin.Context) {
	hostname := c.Param("hostname")
	hostData, exists := h.storage.GetHostData(hostname)
	if !exists {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"title": "主机不存在",
		})
		return
	}

	c.HTML(http.StatusOK, "host.html", gin.H{
		"title":    "主机详情 - " + hostname,
		"hostname": hostname,
		"hostData": hostData,
	})
}
