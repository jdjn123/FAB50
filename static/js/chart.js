// 简化的Chart.js实现
class Chart {
    constructor(ctx, config) {
        this.ctx = ctx;
        this.config = config;
        this.data = config.data;
        this.options = config.options || {};
        this.canvas = ctx.canvas;
        this.width = this.canvas.width;
        this.height = this.canvas.height;
        
        this.init();
    }
    
    init() {
        this.draw();
    }
    
    draw() {
        const ctx = this.ctx;
        ctx.clearRect(0, 0, this.width, this.height);
        
        if (this.config.type === 'line') {
            this.drawLineChart();
        }
    }
    
    drawLineChart() {
        const ctx = this.ctx;
        const datasets = this.data.datasets;
        const labels = this.data.labels;
        
        if (!labels || labels.length === 0) return;
        
        const padding = 40;
        const chartWidth = this.width - padding * 2;
        const chartHeight = this.height - padding * 2;
        const stepX = chartWidth / (labels.length - 1);
        
        // 绘制坐标轴
        ctx.strokeStyle = '#ddd';
        ctx.lineWidth = 1;
        ctx.beginPath();
        ctx.moveTo(padding, padding);
        ctx.lineTo(padding, this.height - padding);
        ctx.lineTo(this.width - padding, this.height - padding);
        ctx.stroke();
        
        // 绘制网格线
        ctx.strokeStyle = '#f0f0f0';
        ctx.lineWidth = 0.5;
        for (let i = 0; i <= 5; i++) {
            const y = padding + (chartHeight / 5) * i;
            ctx.beginPath();
            ctx.moveTo(padding, y);
            ctx.lineTo(this.width - padding, y);
            ctx.stroke();
        }
        
        // 绘制数据线
        datasets.forEach((dataset, datasetIndex) => {
            const data = dataset.data;
            if (!data || data.length === 0) return;
            
            ctx.strokeStyle = dataset.borderColor || '#667eea';
            ctx.lineWidth = 2;
            ctx.beginPath();
            
            data.forEach((value, index) => {
                const x = padding + stepX * index;
                const y = this.height - padding - (value / 100) * chartHeight;
                
                if (index === 0) {
                    ctx.moveTo(x, y);
                } else {
                    ctx.lineTo(x, y);
                }
            });
            
            ctx.stroke();
            
            // 绘制数据点
            ctx.fillStyle = dataset.backgroundColor || '#667eea';
            data.forEach((value, index) => {
                const x = padding + stepX * index;
                const y = this.height - padding - (value / 100) * chartHeight;
                
                ctx.beginPath();
                ctx.arc(x, y, 4, 0, 2 * Math.PI);
                ctx.fill();
            });
        });
        
        // 绘制标签
        ctx.fillStyle = '#666';
        ctx.font = '12px Arial';
        ctx.textAlign = 'center';
        
        labels.forEach((label, index) => {
            const x = padding + stepX * index;
            const y = this.height - padding + 20;
            ctx.fillText(label, x, y);
        });
        
        // 绘制图例
        if (this.options.plugins && this.options.plugins.legend) {
            this.drawLegend();
        }
    }
    
    drawLegend() {
        const ctx = this.ctx;
        const datasets = this.data.datasets;
        const legendY = 20;
        let legendX = this.width - 150;
        
        ctx.font = '12px Arial';
        ctx.textAlign = 'left';
        
        datasets.forEach((dataset, index) => {
            // 绘制颜色块
            ctx.fillStyle = dataset.borderColor || '#667eea';
            ctx.fillRect(legendX, legendY + index * 20, 15, 15);
            
            // 绘制标签
            ctx.fillStyle = '#333';
            ctx.fillText(dataset.label || `Dataset ${index + 1}`, legendX + 20, legendY + index * 20 + 12);
        });
    }
    
    update() {
        this.draw();
    }
}

// 全局Chart对象
window.Chart = Chart;
