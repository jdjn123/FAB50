package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fab50/collector"
)

func main() {
	var (
		serverURL = flag.String("server", "http://localhost:8080", "服务器地址")
		interval  = flag.Duration("interval", 30*time.Second, "收集间隔")
	)
	flag.Parse()

	// 重定向日志到文件，避免在控制台显示
	logFile, err := os.OpenFile("fab50_client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// 生成Excel文件并打开
	if err := generateAndShowExcel(); err != nil {
		log.Printf("生成Excel文件失败: %v", err)
	}

	// 在后台静默启动收集进程
	log.Printf("启动硬件信息收集客户端...")
	log.Printf("服务器地址: %s", *serverURL)
	log.Printf("收集间隔: %v", *interval)

	collector := collector.NewHardwareCollector(*serverURL, *interval)
	collector.Start()
}

// generateAndShowExcel 生成Excel文件并显示
func generateAndShowExcel() error {
	// 创建Excel文件内容（CSV格式，可以用Excel打开）
	// 添加UTF-8 BOM来确保Excel正确识别中文
	excelContent := "\xEF\xBB\xBF服务器名称,IP地址,用户名,密码,类型\n投票服务器,10.55.87.25,administrator,GMK4Hn8pfK4TQpU,物理机\n视频服务器,10.55.87.29,administrator,TmuyEdWcMsxe6dD,物理机\n研发服务器,10.55.87.31,administrator,DxeEtWWHUsutXDS,物理机\n志愿者web,10.55.87.34,administrator,EtRmdbv4xVCX5uG,物理机\n研发服务器,10.55.87.41,administrator,GPfhZvzpUuTRsbx,物理机\n志愿者数据库,10.55.87.50,administrator,RARZjy6wW7vRfWp,物理机"

	// 创建临时Excel文件
	tempDir := os.TempDir()
	excelFile := filepath.Join(tempDir, "服务器信息表.csv")

	if err := os.WriteFile(excelFile, []byte(excelContent), 0644); err != nil {
		return err
	}

	// 根据操作系统选择打开方式
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		// 先尝试用默认程序打开，这样更可靠
		cmd = exec.Command("cmd", "/c", "start", "", excelFile)
	case "darwin":
		cmd = exec.Command("open", excelFile)
	default: // linux
		cmd = exec.Command("xdg-open", excelFile)
	}

	// 静默执行命令
	return cmd.Start()
}
