package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"fab50/types"
)

// Storage 数据存储
type Storage struct {
	mu       sync.RWMutex
	hosts    map[string]*HostData
	maxHosts int
	dataDir  string
}

// HostData 主机数据
type HostData struct {
	Hostname    string                    `json:"hostname"`
	LastUpdate  time.Time                 `json:"last_update"`
	HardwareInfo []types.HardwareInfo     `json:"hardware_info"`
	MaxRecords  int                       `json:"max_records"`
}

// NewStorage 创建新的存储实例
func NewStorage(maxHosts, maxRecords int, dataDir string) *Storage {
	// 确保数据目录存在
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		fmt.Printf("创建数据目录失败: %v\n", err)
	}
	
	return &Storage{
		hosts:     make(map[string]*HostData),
		maxHosts:  maxHosts,
		dataDir:   dataDir,
	}
}

// AddHardwareInfo 添加硬件信息
func (s *Storage) AddHardwareInfo(info *types.HardwareInfo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	hostname := info.Hostname
	
	// 检查主机是否存在
	hostData, exists := s.hosts[hostname]
	if !exists {
		// 如果主机数量超过限制，删除最旧的主机
		if len(s.hosts) >= s.maxHosts {
			s.removeOldestHost()
		}
		
		hostData = &HostData{
			Hostname:    hostname,
			HardwareInfo: make([]types.HardwareInfo, 0),
			MaxRecords:  100, // 每个主机最多保存100条记录
		}
		s.hosts[hostname] = hostData
	}

	// 添加新的硬件信息
	hostData.HardwareInfo = append(hostData.HardwareInfo, *info)
	hostData.LastUpdate = time.Now()

	// 如果记录数超过限制，删除最旧的记录
	if len(hostData.HardwareInfo) > hostData.MaxRecords {
		hostData.HardwareInfo = hostData.HardwareInfo[1:]
	}

	// 保存数据到文件
	s.saveDataToFile(hostname, hostData)
}

// GetHosts 获取所有主机列表
func (s *Storage) GetHosts() map[string]*HostData {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*HostData)
	for hostname, hostData := range s.hosts {
		result[hostname] = hostData
	}
	return result
}

// GetHostData 获取指定主机的数据
func (s *Storage) GetHostData(hostname string) (*HostData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hostData, exists := s.hosts[hostname]
	return hostData, exists
}

// GetLatestHardwareInfo 获取最新的硬件信息
func (s *Storage) GetLatestHardwareInfo() map[string]*types.HardwareInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]*types.HardwareInfo)
	for hostname, hostData := range s.hosts {
		if len(hostData.HardwareInfo) > 0 {
			latest := hostData.HardwareInfo[len(hostData.HardwareInfo)-1]
			result[hostname] = &latest
		}
	}
	return result
}

// removeOldestHost 删除最旧的主机
func (s *Storage) removeOldestHost() {
	var oldestHost string
	var oldestTime time.Time

	for hostname, hostData := range s.hosts {
		if oldestHost == "" || hostData.LastUpdate.Before(oldestTime) {
			oldestHost = hostname
			oldestTime = hostData.LastUpdate
		}
	}

	if oldestHost != "" {
		delete(s.hosts, oldestHost)
	}
}

// Cleanup 清理过期数据
func (s *Storage) Cleanup(maxAge time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for hostname, hostData := range s.hosts {
		if now.Sub(hostData.LastUpdate) > maxAge {
			delete(s.hosts, hostname)
		}
	}
}

// saveDataToFile 保存数据到文件
func (s *Storage) saveDataToFile(hostname string, hostData *HostData) {
	// 创建主机数据文件路径
	filename := filepath.Join(s.dataDir, hostname+".json")
	
	// 将数据序列化为JSON
	data, err := json.MarshalIndent(hostData, "", "  ")
	if err != nil {
		fmt.Printf("序列化主机 %s 数据失败: %v\n", hostname, err)
		return
	}
	
	// 写入文件
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("保存主机 %s 数据到文件失败: %v\n", hostname, err)
		return
	}
	
	fmt.Printf("主机 %s 数据已保存到文件: %s\n", hostname, filename)
}
