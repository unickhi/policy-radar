#!/bin/bash

# Policy Radar 一键停止脚本

echo "========================================="
echo "  Policy Radar 服务停止脚本"
echo "========================================="

# 停止后端服务 (端口 8080)
echo "[1/3] 停止后端服务..."
PID=$(lsof -ti :8080 2>/dev/null)
if [ -n "$PID" ]; then
    kill $PID 2>/dev/null
    echo "  ✓ 后端服务已停止 (PID: $PID)"
else
    echo "  - 后端服务未运行"
fi

# 停止前端管理后台 (端口 3000)
echo "[2/3] 停止前端管理后台..."
PID=$(lsof -ti :3000 2>/dev/null)
if [ -n "$PID" ]; then
    kill $PID 2>/dev/null
    echo "  ✓ 前端管理后台已停止 (PID: $PID)"
else
    echo "  - 前端管理后台未运行"
fi

# 停止H5前端 (端口 3001)
echo "[3/3] 停止H5前端..."
PID=$(lsof -ti :3001 2>/dev/null)
if [ -n "$PID" ]; then
    kill $PID 2>/dev/null
    echo "  ✓ H5前端已停止 (PID: $PID)"
else
    echo "  - H5前端未运行"
fi

echo ""
echo "========================================="
echo "  所有服务已停止!"
echo "========================================="