# Makefile for FAB50 Hardware Monitoring System

.PHONY: help build-server build-client run-server run-client clean test deps

# 默认目标
help:
	@echo "FAB50 硬件监控系统"
	@echo ""
	@echo "可用命令:"
	@echo "  deps          - 下载依赖"
	@echo "  build-server  - 构建服务器"
	@echo "  build-client  - 构建客户端"
	@echo "  build         - 构建所有组件"
	@echo "  run-server    - 运行服务器"
	@echo "  run-client    - 运行客户端"
	@echo "  clean         - 清理构建文件"
	@echo "  test          - 运行测试"

# 下载依赖
deps:
	go mod download
	go mod tidy

# 构建服务器
build-server:
	@echo "构建服务器..."
	go build -o bin/server cmd/server/main.go

# 构建客户端
build-client:
	@echo "构建客户端..."
	go build -o bin/client cmd/client/main.go

# 构建所有组件
build: build-server build-client
	@echo "构建完成!"

# 运行服务器
run-server:
	@echo "启动服务器..."
	go run cmd/server/main.go

# 运行客户端
run-client:
	@echo "启动客户端..."
	go run cmd/client/main.go

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	go clean

# 运行测试
test:
	go test ./...

# 安装依赖并构建
all: deps build
	@echo "项目准备完成!"

# 开发模式 - 同时运行服务器和客户端
dev: deps
	@echo "开发模式启动..."
	@echo "服务器将在 http://localhost:8080 启动"
	@echo "按 Ctrl+C 停止"
	@trap 'kill %1; kill %2' SIGINT; \
	go run cmd/server/main.go & \
	sleep 2 && \
	go run cmd/client/main.go & \
	wait
