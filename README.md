# Policy Radar 项目

政策雷达 - 政策数据管理与展示系统

## 项目结构

```
policy radar/
├── backend/          # Go 后端服务
├── frontend/         # Vue 管理后台前端
├── frontend-h5/      # Vue H5 移动端前端
├── Datas/            # 数据文件
├── docs/             # 文档
├── start.sh          # 一键启动脚本
└── stop.sh           # 一键停止脚本
```

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 后端 API | 8080 | Go 后端服务 |
| 管理后台 | 3000 | Vue 前端管理后台 |
| H5 前端 | 3001 | Vue H5 移动端 |

## 快速启动

### 一键启动所有服务

```bash
./start.sh
```

### 一键停止所有服务

```bash
./stop.sh
```

## 单独启动/停止服务

### 后端服务 (Go)

```bash
# 启动
cd backend
./server

# 后台启动
cd backend
nohup ./server > server.log 2>&1 &

# 停止
kill $(lsof -ti :8080)
```

### 前端管理后台 (Vue)

```bash
# 启动
cd frontend
npm run dev

# 后台启动
cd frontend
nohup npm run dev > frontend.log 2>&1 &

# 停止
kill $(lsof -ti :3000)
```

### H5 前端 (Vue)

```bash
# 启动
cd frontend-h5
npm run dev

# 后台启动
cd frontend-h5
nohup npm run dev > h5.log 2>&1 &

# 停止
kill $(lsof -ti :3001)
```

## 访问地址

启动服务后，可通过以下地址访问：

- **后端 API**: http://localhost:8080
- **管理后台**: http://localhost:3000
- **H5 前端**: http://localhost:3001

## 日志查看

```bash
# 查看后端日志
tail -f backend/server.log

# 查看管理后台日志
tail -f frontend/frontend.log

# 查看 H5 前端日志
tail -f frontend-h5/h5.log
```

## 查看服务状态

```bash
# 查看所有服务端口
lsof -i :8080 -i :3000 -i :3001 | grep LISTEN
```

## 依赖环境

- Go 1.21+
- Node.js 18+
- npm 或 yarn

## 初始化安装

如果是首次运行，需要先安装依赖：

```bash
# 安装前端依赖
cd frontend && npm install
cd ../frontend-h5 && npm install
```