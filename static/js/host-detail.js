// 主机详情页JavaScript
class HostDetail {
    constructor() {
        this.ws = null;
        this.historyChart = null;
        this.hostname = this.getHostnameFromUrl();
        this.hostData = null;
        this.init();
    }

    getHostnameFromUrl() {
        const pathParts = window.location.pathname.split('/');
        return pathParts[pathParts.length - 1];
    }

    init() {
        this.connectWebSocket();
        this.initHistoryChart();
        this.loadHostData();
        this.startPeriodicUpdate();
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            console.log('WebSocket连接已建立');
            this.updateConnectionStatus(true);
        };
        
        this.ws.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleWebSocketMessage(data);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket连接已关闭');
            this.updateConnectionStatus(false);
            setTimeout(() => this.connectWebSocket(), 5000);
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket错误:', error);
            this.updateConnectionStatus(false);
        };
    }

    updateConnectionStatus(connected) {
        const statusElement = document.getElementById('connection-status');
        if (connected) {
            statusElement.textContent = '已连接';
            statusElement.className = 'badge bg-success';
        } else {
            statusElement.textContent = '连接断开';
            statusElement.className = 'badge bg-danger';
        }
    }

    handleWebSocketMessage(data) {
        if (data.type === 'hardware_info' && data.data[this.hostname]) {
            this.updateHostData(data.data[this.hostname]);
        }
    }

    async loadHostData() {
        try {
            const response = await fetch(`/api/hosts/${this.hostname}?limit=100`);
            if (!response.ok) {
                throw new Error('主机数据加载失败');
            }
            
            const data = await response.json();
            this.hostData = data;
            this.updateDisplay();
            this.updateHistoryChart();
        } catch (error) {
            console.error('加载主机数据失败:', error);
            this.showError('加载主机数据失败');
        }
    }

    updateHostData(newData) {
        if (this.hostData && this.hostData.hardware_info) {
            this.hostData.hardware_info.push(newData);
            // 保持最近100条记录
            if (this.hostData.hardware_info.length > 100) {
                this.hostData.hardware_info = this.hostData.hardware_info.slice(-100);
            }
        }
        this.updateDisplay();
        this.updateHistoryChart();
    }

    updateDisplay() {
        if (!this.hostData || !this.hostData.hardware_info || this.hostData.hardware_info.length === 0) {
            return;
        }

        const latest = this.hostData.hardware_info[this.hostData.hardware_info.length - 1];
        
        // 更新基本指标
        document.getElementById('cpu-usage').textContent = (latest.cpu.usage || 0).toFixed(1) + '%';
        document.getElementById('memory-usage').textContent = (latest.memory.usage || 0).toFixed(1) + '%';
        
        // 计算平均磁盘使用率
        let avgDiskUsage = 0;
        if (latest.disk && latest.disk.partitions) {
            avgDiskUsage = latest.disk.partitions.reduce((sum, partition) => {
                return sum + (partition.usage || 0);
            }, 0) / latest.disk.partitions.length;
        }
        document.getElementById('disk-usage').textContent = avgDiskUsage.toFixed(1) + '%';
        
        document.getElementById('last-update').textContent = this.formatTime(latest.timestamp);
        
        // 更新详细信息
        this.updateCPUInfo(latest.cpu);
        this.updateMemoryInfo(latest.memory);
        this.updateDiskInfo(latest.disk);
        this.updateNetworkInfo(latest.network);
    }

    updateCPUInfo(cpu) {
        const container = document.getElementById('cpu-info');
        container.innerHTML = `
            <div class="row">
                <div class="col-6">
                    <strong>型号:</strong><br>
                    <span class="text-muted">${cpu.model_name || '未知'}</span>
                </div>
                <div class="col-6">
                    <strong>核心数:</strong><br>
                    <span class="text-muted">${cpu.cores || 0}</span>
                </div>
            </div>
            <div class="row mt-3">
                <div class="col-6">
                    <strong>频率:</strong><br>
                    <span class="text-muted">${this.formatFrequency(cpu.frequency)}</span>
                </div>
                <div class="col-6">
                    <strong>使用率:</strong><br>
                    <div class="progress mt-1">
                        <div class="progress-bar" role="progressbar" 
                             style="width: ${cpu.usage || 0}%" 
                             aria-valuenow="${cpu.usage || 0}" 
                             aria-valuemin="0" aria-valuemax="100">
                            ${(cpu.usage || 0).toFixed(1)}%
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    updateMemoryInfo(memory) {
        const container = document.getElementById('memory-info');
        const totalGB = (memory.total / (1024 * 1024 * 1024)).toFixed(1);
        const usedGB = (memory.used / (1024 * 1024 * 1024)).toFixed(1);
        const freeGB = (memory.free / (1024 * 1024 * 1024)).toFixed(1);
        
        container.innerHTML = `
            <div class="row">
                <div class="col-6">
                    <strong>总内存:</strong><br>
                    <span class="text-muted">${totalGB} GB</span>
                </div>
                <div class="col-6">
                    <strong>已使用:</strong><br>
                    <span class="text-muted">${usedGB} GB</span>
                </div>
            </div>
            <div class="row mt-3">
                <div class="col-6">
                    <strong>可用内存:</strong><br>
                    <span class="text-muted">${freeGB} GB</span>
                </div>
                <div class="col-6">
                    <strong>使用率:</strong><br>
                    <div class="progress mt-1">
                        <div class="progress-bar bg-info" role="progressbar" 
                             style="width: ${memory.usage || 0}%" 
                             aria-valuenow="${memory.usage || 0}" 
                             aria-valuemin="0" aria-valuemax="100">
                            ${(memory.usage || 0).toFixed(1)}%
                        </div>
                    </div>
                </div>
            </div>
        `;
    }

    updateDiskInfo(disk) {
        const container = document.getElementById('disk-info');
        if (!disk || !disk.partitions) {
            container.innerHTML = '<p class="text-muted">暂无磁盘信息</p>';
            return;
        }
        
        let html = '';
        disk.partitions.forEach(partition => {
            const totalGB = (partition.total / (1024 * 1024 * 1024)).toFixed(1);
            const usedGB = (partition.used / (1024 * 1024 * 1024)).toFixed(1);
            const freeGB = (partition.free / (1024 * 1024 * 1024)).toFixed(1);
            
            html += `
                <div class="card mb-3">
                    <div class="card-body">
                        <h6 class="card-title">${partition.device}</h6>
                        <p class="card-text text-muted">挂载点: ${partition.mount_point}</p>
                        <div class="row">
                            <div class="col-4">
                                <small class="text-muted">总容量</small><br>
                                <strong>${totalGB} GB</strong>
                            </div>
                            <div class="col-4">
                                <small class="text-muted">已使用</small><br>
                                <strong>${usedGB} GB</strong>
                            </div>
                            <div class="col-4">
                                <small class="text-muted">可用</small><br>
                                <strong>${freeGB} GB</strong>
                            </div>
                        </div>
                        <div class="mt-2">
                            <small class="text-muted">使用率</small>
                            <div class="progress">
                                <div class="progress-bar bg-warning" role="progressbar" 
                                     style="width: ${partition.usage || 0}%" 
                                     aria-valuenow="${partition.usage || 0}" 
                                     aria-valuemin="0" aria-valuemax="100">
                                    ${(partition.usage || 0).toFixed(1)}%
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            `;
        });
        
        container.innerHTML = html;
    }

    updateNetworkInfo(network) {
        const container = document.getElementById('network-info');
        if (!network || !network.interfaces) {
            container.innerHTML = '<p class="text-muted">暂无网络信息</p>';
            return;
        }
        
        let html = '';
        network.interfaces.forEach(iface => {
            const sentMB = (iface.bytes_sent / (1024 * 1024)).toFixed(2);
            const recvMB = (iface.bytes_recv / (1024 * 1024)).toFixed(2);
            
            html += `
                <div class="card mb-3">
                    <div class="card-body">
                        <h6 class="card-title">${iface.name}</h6>
                        <div class="row">
                            <div class="col-6">
                                <small class="text-muted">发送</small><br>
                                <strong>${sentMB} MB</strong>
                            </div>
                            <div class="col-6">
                                <small class="text-muted">接收</small><br>
                                <strong>${recvMB} MB</strong>
                            </div>
                        </div>
                        <div class="row mt-2">
                            <div class="col-6">
                                <small class="text-muted">发送包数</small><br>
                                <strong>${iface.packets_sent}</strong>
                            </div>
                            <div class="col-6">
                                <small class="text-muted">接收包数</small><br>
                                <strong>${iface.packets_recv}</strong>
                            </div>
                        </div>
                        ${iface.addresses && iface.addresses.length > 0 ? `
                            <div class="mt-2">
                                <small class="text-muted">IP地址</small><br>
                                <code>${iface.addresses.join(', ')}</code>
                            </div>
                        ` : ''}
                    </div>
                </div>
            `;
        });
        
        container.innerHTML = html;
    }

    initHistoryChart() {
        const ctx = document.getElementById('historyChart').getContext('2d');
        this.historyChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: [],
                datasets: [{
                    label: 'CPU使用率',
                    data: [],
                    borderColor: 'rgb(75, 192, 192)',
                    backgroundColor: 'rgba(75, 192, 192, 0.2)',
                    tension: 0.1
                }, {
                    label: '内存使用率',
                    data: [],
                    borderColor: 'rgb(255, 99, 132)',
                    backgroundColor: 'rgba(255, 99, 132, 0.2)',
                    tension: 0.1
                }, {
                    label: '磁盘使用率',
                    data: [],
                    borderColor: 'rgb(255, 205, 86)',
                    backgroundColor: 'rgba(255, 205, 86, 0.2)',
                    tension: 0.1
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100
                    }
                },
                plugins: {
                    legend: {
                        position: 'top',
                    }
                }
            }
        });
    }

    updateHistoryChart() {
        if (!this.hostData || !this.hostData.hardware_info) return;
        
        const data = this.hostData.hardware_info;
        const labels = [];
        const cpuData = [];
        const memoryData = [];
        const diskData = [];
        
        data.forEach(item => {
            labels.push(this.formatTime(item.timestamp));
            cpuData.push(item.cpu.usage || 0);
            memoryData.push(item.memory.usage || 0);
            
            // 计算平均磁盘使用率
            let avgDiskUsage = 0;
            if (item.disk && item.disk.partitions) {
                avgDiskUsage = item.disk.partitions.reduce((sum, partition) => {
                    return sum + (partition.usage || 0);
                }, 0) / item.disk.partitions.length;
            }
            diskData.push(avgDiskUsage);
        });
        
        this.historyChart.data.labels = labels;
        this.historyChart.data.datasets[0].data = cpuData;
        this.historyChart.data.datasets[1].data = memoryData;
        this.historyChart.data.datasets[2].data = diskData;
        
        this.historyChart.update();
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString('zh-CN');
    }

    formatFrequency(freq) {
        if (!freq) return '未知';
        return (freq / 1000).toFixed(1) + ' GHz';
    }

    showError(message) {
        // 可以添加错误提示UI
        console.error(message);
    }

    startPeriodicUpdate() {
        // 每30秒刷新一次数据
        setInterval(() => {
            this.loadHostData();
        }, 30000);
    }
}

// 页面加载完成后初始化主机详情页
document.addEventListener('DOMContentLoaded', () => {
    new HostDetail();
});
