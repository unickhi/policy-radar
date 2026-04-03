#!/bin/bash

# Policy Radar 一键启动脚本

PROJECT_DIR="/Users/zhangyuhao/Desktop/极客JIKER/国王/policy radar"

echo "========================================="
echo "  Policy Radar 服务启动脚本"
echo "========================================="

# 检查并释放8080端口（可能被Zookeeper占用）
if lsof -i :8080 | grep -q LISTEN; then
    echo "[准备] 检测到8080端口被占用，正在释放..."
    brew services stop zookeeper 2>/dev/null
    sleep 2
    # 如果仍被占用，强制停止
    PID=$(lsof -ti :8080 2>/dev/null)
    if [ -n "$PID" ]; then
        kill $PID 2>/dev/null
        sleep 1
    fi
    echo "  ✓ 8080端口已释放"
fi

# 启动后端服务
echo "[1/2] 启动后端服务 (端口 8080)..."
cd "$PROJECT_DIR/backend"
if [ -f "server" ]; then
    nohup ./server > server.log 2>&1 &
    echo "  ✓ 后端服务已启动"
else
    echo "  ✗ 错误: server 可执行文件不存在"
fi

# 启动前端服务
echo "[2/2] 启动前端服务 (端口 3000)..."
cd "$PROJECT_DIR/frontend"
nohup npm run dev > frontend.log 2>&1 &
echo "  ✓ 前端服务已启动"

echo ""
echo "========================================="
echo "  所有服务已启动完成!"
echo "========================================="
echo ""
echo "访问地址:"
echo "  - 后端API:    http://localhost:8080"
echo "  - PC管理后台: http://localhost:3000/admin/national"
echo "  - H5移动端:   http://localhost:3000/m/"
echo ""
echo "日志文件位置:"
echo "  - 后端:       backend/server.log"
echo "  - 前端:       frontend/frontend.log"