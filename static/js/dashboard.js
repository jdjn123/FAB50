// 仪表板JavaScript
class Dashboard {
    constructor() {
        this.ws = null;
        this.performanceChart = null;
        this.hostsData = {};
        this.init();
    }

    init() {
        this.connectWebSocket();
        this.initPerformanceChart();
        this.loadInitialData();
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
            // 尝试重新连接
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
        switch (data.type) {
            case 'hardware_info':
                this.updateHostsData(data.data);
                this.updateDashboard();
                break;
            case 'host_list':
                this.updateHostsData(data.data);
                this.updateDashboard();
                break;
        }
    }

    async loadInitialData() {
        try {
            const response = await fetch('/api/latest');
            const data = await response.json();
            this.updateHostsData(data);
            this.updateDashboard();
        } catch (error) {
            console.error('加载初始数据失败:', error);
        }
    }

    updateHostsData(data) {
        this.hostsData = data;
    }

    updateDashboard() {
        this.updateMetrics();
        this.updateHostsList();
        this.updatePerformanceChart();
    }

    updateMetrics() {
        const hosts = Object.values(this.hostsData);
        const totalHosts = hosts.length;
        
        let totalCpu = 0;
        let totalMemory = 0;
        let totalDisk = 0;
        
        hosts.forEach(host => {
            totalCpu += host.cpu.usage || 0;
            totalMemory += host.memory.usage || 0;
            
            // 计算平均磁盘使用率
            if (host.disk && host.disk.partitions) {
                const diskUsage = host.disk.partitions.reduce((sum, partition) => {
                    return sum + (partition.usage || 0);
                }, 0) / host.disk.partitions.length;
                totalDisk += diskUsage;
            }
        });
        
        const avgCpu = totalHosts > 0 ? (totalCpu / totalHosts).toFixed(1) : 0;
        const avgMemory = totalHosts > 0 ? (totalMemory / totalHosts).toFixed(1) : 0;
        const avgDisk = totalHosts > 0 ? (totalDisk / totalHosts).toFixed(1) : 0;
        
        document.getElementById('total-hosts').textContent = totalHosts;
        document.getElementById('avg-cpu').textContent = avgCpu + '%';
        document.getElementById('avg-memory').textContent = avgMemory + '%';
        document.getElementById('total-disk').textContent = avgDisk + '%';
    }

    updateHostsList() {
        const container = document.getElementById('hosts-container');
        container.innerHTML = '';
        
        Object.entries(this.hostsData).forEach(([hostname, host]) => {
            const hostCard = this.createHostCard(hostname, host);
            container.appendChild(hostCard);
        });
    }

    createHostCard(hostname, host) {
        const col = document.createElement('div');
        col.className = 'col-md-4 col-lg-3 mb-4';
        
        const isOnline = this.isHostOnline(host);
        const statusClass = isOnline ? 'status-online' : 'status-offline';
        const statusText = isOnline ? '在线' : '离线';
        
        col.innerHTML = `
            <div class="card host-card h-100">
                <div class="card-body">
                    <div class="d-flex justify-content-between align-items-start mb-3">
                        <h6 class="card-title mb-0">${hostname}</h6>
                        <span class="badge ${statusClass}">${statusText}</span>
                    </div>
                    
                    <div class="mb-3">
                        <small class="text-muted">CPU使用率</small>
                        <div class="progress mb-2">
                            <div class="progress-bar" role="progressbar" 
                                 style="width: ${host.cpu.usage || 0}%" 
                                 aria-valuenow="${host.cpu.usage || 0}" 
                                 aria-valuemin="0" aria-valuemax="100">
                                ${(host.cpu.usage || 0).toFixed(1)}%
                            </div>
                        </div>
                    </div>
                    
                    <div class="mb-3">
                        <small class="text-muted">内存使用率</small>
                        <div class="progress mb-2">
                            <div class="progress-bar bg-info" role="progressbar" 
                                 style="width: ${host.memory.usage || 0}%" 
                                 aria-valuenow="${host.memory.usage || 0}" 
                                 aria-valuemin="0" aria-valuemax="100">
                                ${(host.memory.usage || 0).toFixed(1)}%
                            </div>
                        </div>
                    </div>
                    
                    <div class="mb-3">
                        <small class="text-muted">操作系统</small>
                        <div class="text-truncate">${host.os.name || '未知'} ${host.os.version || ''}</div>
                    </div>
                    
                    <div class="text-muted small">
                        最后更新: ${this.formatTime(host.timestamp)}
                    </div>
                </div>
                <div class="card-footer">
                    <a href="/host/${hostname}" class="btn btn-sm btn-outline-primary w-100">
                        查看详情
                    </a>
                </div>
            </div>
        `;
        
        return col;
    }

    isHostOnline(host) {
        const now = new Date();
        const lastUpdate = new Date(host.timestamp);
        const diffMinutes = (now - lastUpdate) / (1000 * 60);
        return diffMinutes < 5; // 5分钟内更新认为在线
    }

    formatTime(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleString('zh-CN');
    }

    initPerformanceChart() {
        const ctx = document.getElementById('performanceChart').getContext('2d');
        this.performanceChart = new Chart(ctx, {
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

    updatePerformanceChart() {
        const hosts = Object.values(this.hostsData);
        if (hosts.length === 0) return;
        
        const now = new Date();
        const timeLabel = now.toLocaleTimeString('zh-CN');
        
        // 计算平均值
        const avgCpu = hosts.reduce((sum, host) => sum + (host.cpu.usage || 0), 0) / hosts.length;
        const avgMemory = hosts.reduce((sum, host) => sum + (host.memory.usage || 0), 0) / hosts.length;
        
        // 更新图表数据
        this.performanceChart.data.labels.push(timeLabel);
        this.performanceChart.data.datasets[0].data.push(avgCpu);
        this.performanceChart.data.datasets[1].data.push(avgMemory);
        
        // 保持最近20个数据点
        if (this.performanceChart.data.labels.length > 20) {
            this.performanceChart.data.labels.shift();
            this.performanceChart.data.datasets[0].data.shift();
            this.performanceChart.data.datasets[1].data.shift();
        }
        
        this.performanceChart.update();
    }

    startPeriodicUpdate() {
        // 每30秒更新一次图表
        setInterval(() => {
            this.updatePerformanceChart();
        }, 30000);
    }
}

// 页面加载完成后初始化仪表板
document.addEventListener('DOMContentLoaded', () => {
    new Dashboard();
});
