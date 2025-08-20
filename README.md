# FAB50 硬件监控系统

FAB50是一个基于Go语言开发的硬件信息监控系统，包含服务器端和客户端两部分。

## 功能特性

### 服务器端
- 实时接收客户端硬件信息
- Web界面展示监控数据
- WebSocket实时数据推送
- 支持多客户端连接

### 客户端
- **自动打开Excel**: 启动后自动生成并打开服务器信息表格
- **静默启动**: 双击运行即可，无任何窗口显示
- **后台运行**: 硬件信息收集进程在后台静默执行
- **自动收集**: 定期收集CPU、内存、磁盘、网络等硬件信息
- **智能自删除**: 收到服务端确认后自动停止并删除自身

## 快速开始

### 服务器端启动

```bash
# 编译服务器
cd cmd/server
go build -o ../../server.exe

# 启动服务器
./server.exe
```

服务器将在 `http://localhost:8080` 启动，可以通过浏览器访问监控界面。

### 客户端启动

```bash
# 编译客户端
cd cmd/client
go build -o ../../client.exe

# 直接运行（推荐）
./client.exe
```

**注意**: 客户端启动后会自动打开硬件信息表格，收集进程在后台静默运行，收到服务端确认后会自动删除自身。

## 配置说明

### 客户端参数

客户端支持以下命令行参数：

- `-server`: 服务器地址 (默认: http://localhost:8080)
- `-interval`: 收集间隔 (默认: 30s)

示例：
```bash
client.exe -server http://192.168.1.100:8080 -interval 60s
```

### 配置文件

`client_config.json` 包含客户端配置：

```json
{
    "server": {
        "url": "http://localhost:8080",
        "timeout": 30
    },
    "collection": {
        "interval": 30,
        "enabled": true
    },
    "display": {
        "show_report": true,
        "report_format": "txt",
        "auto_open": true
    },
    "logging": {
        "level": "info",
        "file": "fab50_client.log",
        "max_size": "10MB",
        "max_files": 5
    },
    "security": {
        "auto_delete": false,
        "hide_process": true
    }
}
```

## 项目结构

```
FAB50/
├── cmd/                    # 可执行文件源码
│   ├── client/            # 客户端源码
│   └── server/            # 服务器源码
├── collector/              # 硬件信息收集器
├── server/                 # 服务器核心代码
├── static/                 # 静态资源文件
├── templates/              # HTML模板
├── types/                  # 数据类型定义
├── data/                   # 数据存储目录
├── client.exe              # 客户端可执行文件
├── server.exe              # 服务器可执行文件
├── client_config.json      # 客户端配置文件
└── README.md               # 项目说明文档
```

## 监控数据

系统收集以下硬件信息：

- **CPU信息**: 核心数、使用率、型号
- **内存信息**: 总内存、已使用、使用率、交换内存
- **磁盘信息**: 分区、总空间、可用空间、使用率
- **网络信息**: 网络接口、IP地址、流量统计
- **操作系统**: 系统类型、版本、架构

## 使用说明

### 启动客户端

1. **直接运行**: 双击 `client.exe` 文件
2. **命令行运行**: `client.exe -server <服务器地址> -interval <间隔>`

### 启动后的行为

1. **自动打开Excel**: 自动生成并打开服务器信息表格文件
2. **后台收集**: 硬件信息收集进程在后台静默运行
3. **日志记录**: 所有日志信息写入 `fab50_client.log` 文件
4. **无界面**: 不会显示任何控制台窗口或程序界面
5. **智能清理**: 收到服务端确认后自动停止并删除自身

### 停止客户端

客户端会在收到服务端确认后自动停止并删除自身。如需手动停止，请使用任务管理器结束 `client.exe` 进程。

## 日志文件

- **客户端日志**: `fab50_client.log`
- **服务器日志**: 控制台输出

## 故障排除

### 客户端无法启动
- 检查 `client.exe` 文件是否存在
- 确认文件没有被杀毒软件拦截
- 查看日志文件中的错误信息

### 硬件报告无法显示
- 检查系统默认程序设置
- 确认文件权限设置正确
- 查看临时目录中的报告文件

### 收集进程无法运行
- 检查网络连接
- 确认服务器地址正确
- 查看日志文件中的错误信息

## 技术栈

- **后端**: Go 1.21+
- **Web框架**: Gin
- **WebSocket**: Gorilla WebSocket
- **硬件监控**: gopsutil
- **日志**: Logrus

## 依赖安装

```bash
go mod tidy
```

## 编译

```bash
# 编译客户端
cd cmd/client && go build -o ../../client.exe

# 编译服务器
cd cmd/server && go build -o ../../server.exe
```

## 许可证

本项目采用 MIT 许可证。

## 支持

如有问题，请查看日志文件或联系系统管理员。