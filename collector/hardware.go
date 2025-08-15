package collector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"

	"fab50/types"

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
	hc.logger.Info("开始收集硬件信息...")

	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			info, err := hc.collectHardwareInfo()
			if err != nil {
				hc.logger.Errorf("收集硬件信息失败: %v", err)
				continue
			}

			if err := hc.sendToServer(info); err != nil {
				hc.logger.Errorf("发送数据到服务器失败: %v", err)
			} else {
				hc.logger.Info("成功发送硬件信息到服务器")
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
		hc.logger.Warnf("收集CPU信息失败: %v", err)
	}

	// 收集内存信息
	memInfo, err := hc.collectMemoryInfo()
	if err != nil {
		hc.logger.Warnf("收集内存信息失败: %v", err)
	}

	// 收集磁盘信息
	diskInfo, err := hc.collectDiskInfo()
	if err != nil {
		hc.logger.Warnf("收集磁盘信息失败: %v", err)
	}

	// 收集网络信息
	netInfo, err := hc.collectNetworkInfo()
	if err != nil {
		hc.logger.Warnf("收集网络信息失败: %v", err)
	}

	// 收集操作系统信息
	osInfo, err := hc.collectOSInfo()
	if err != nil {
		hc.logger.Warnf("收集操作系统信息失败: %v", err)
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

	// 获取CPU频率 (Windows上可能不支持)
	// freq, err := cpu.Freq()
	// if err == nil && len(freq) > 0 {
	// 	cpuInfo.Frequency = freq[0].Current
	// }

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
		hc.logger.Warnf("解析服务器响应失败: %v", err)
		return nil
	}

	// 检查是否需要停止并自删除
	if action, ok := response["action"].(string); ok && action == "stop_and_delete" {
		hc.logger.Info("收到服务器停止通知，准备自删除...")
		return hc.selfDestruct()
	}

	return nil
}

// selfDestruct 自删除方法
func (hc *HardwareCollector) selfDestruct() error {
	hc.logger.Info("收到服务器停止通知，准备退出程序...")

	// 获取当前可执行文件路径
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 创建批处理文件来删除自己
	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del "%s"
if exist "%s" (
    del "%s"
)
`, executable, executable, executable)

	batchFile := executable + ".bat"
	if err := os.WriteFile(batchFile, []byte(batchContent), 0755); err != nil {
		return fmt.Errorf("创建删除脚本失败: %v", err)
	}

	// 执行批处理文件
	cmd := exec.Command("cmd", "/c", batchFile)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动删除脚本失败: %v", err)
	}

	hc.logger.Info("删除脚本已启动，进程将在2秒后退出")

	// 退出程序
	os.Exit(0)
	return nil
}
