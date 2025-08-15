package main

import (
	"flag"
	"log"
	"time"

	"fab50/server"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		port       = flag.String("port", "8080", "服务器端口")
		maxHosts   = flag.Int("max-hosts", 100, "最大主机数量")
		maxRecords = flag.Int("max-records", 1000, "每个主机最大记录数")
		dataDir    = flag.String("data-dir", "./data", "数据保存目录")
	)
	flag.Parse()

	log.Printf("启动硬件监控服务器...")
	log.Printf("端口: %s", *port)
	log.Printf("最大主机数: %d", *maxHosts)
	log.Printf("每个主机最大记录数: %d", *maxRecords)

	// 创建存储
	storage := server.NewStorage(*maxHosts, *maxRecords, *dataDir)

	// 创建WebSocket管理器
	wsManager := server.NewWebSocketManager()
	go wsManager.Start()

	// 创建处理器
	handlers := server.NewHandlers(storage, wsManager)

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 创建路由
	r := gin.Default()

	// 设置路由
	handlers.SetupRoutes(r)

	// 启动数据清理协程
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				storage.Cleanup(24 * time.Hour) // 清理24小时前的数据
			}
		}
	}()

	// 启动服务器
	log.Printf("服务器启动在端口 %s", *port)
	if err := r.Run(":" + *port); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
