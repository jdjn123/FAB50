package main

import (
	"flag"
	"log"
	"time"

	"fab50/collector"
)

func main() {
	var (
		serverURL = flag.String("server", "http://localhost:8080", "服务器地址")
		interval  = flag.Duration("interval", 30*time.Second, "收集间隔")
	)
	flag.Parse()

	log.Printf("启动硬件信息收集客户端...")
	log.Printf("服务器地址: %s", *serverURL)
	log.Printf("收集间隔: %v", *interval)

	collector := collector.NewHardwareCollector(*serverURL, *interval)
	collector.Start()
}
