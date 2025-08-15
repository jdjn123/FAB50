package types

import (
	"time"
)

// HardwareInfo 硬件信息结构体
type HardwareInfo struct {
	Hostname    string    `json:"hostname"`
	Timestamp   time.Time `json:"timestamp"`
	CPU         CPUInfo   `json:"cpu"`
	Memory      MemInfo   `json:"memory"`
	Disk        DiskInfo  `json:"disk"`
	Network     NetInfo   `json:"network"`
	OS          OSInfo    `json:"os"`
}

// CPUInfo CPU信息
type CPUInfo struct {
	ModelName   string  `json:"model_name"`
	Cores       int     `json:"cores"`
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature"`
	Frequency   float64 `json:"frequency"`
}

// MemInfo 内存信息
type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Usage       float64 `json:"usage"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapFree    uint64  `json:"swap_free"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Partitions []PartitionInfo `json:"partitions"`
}

// PartitionInfo 分区信息
type PartitionInfo struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mount_point"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	Usage      float64 `json:"usage"`
}

// NetInfo 网络信息
type NetInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
}

// InterfaceInfo 网络接口信息
type InterfaceInfo struct {
	Name        string   `json:"name"`
	Addresses   []string `json:"addresses"`
	BytesSent   uint64   `json:"bytes_sent"`
	BytesRecv   uint64   `json:"bytes_recv"`
	PacketsSent uint64   `json:"packets_sent"`
	PacketsRecv uint64   `json:"packets_recv"`
}

// OSInfo 操作系统信息
type OSInfo struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Architecture string `json:"architecture"`
	Platform     string `json:"platform"`
}
