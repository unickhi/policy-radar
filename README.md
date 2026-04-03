# Policy Radar 项目

政策雷达 - 电力行业标准政策管理与展示系统

## 项目结构

```
policy-radar/
├── frontend/               # Vue 3 前端（PC管理后台 + H5移动端）
│   ├── src/
│   │   ├── api/           # API 封装
│   │   ├── components/    # 公共组件
│   │   ├── pages/         # 页面
│   │   │   ├── admin/     # PC管理后台页面
│   │   │   └── mobile/    # H5移动端页面
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # Pinia 状态管理
│   │   └── styles/        # 样式文件
│   └── package.json
├── backend/                # Go 后端服务
│   ├── cmd/server/        # 服务入口
│   ├── internal/          # 内部代码
│   │   ├── handler/       # HTTP 处理器
│   │   ├── service/       # 业务逻辑
│   │   ├── repository/    # 数据访问
│   │   ├── model/         # 数据模型
│   │   └── config/        # 配置
│   ├── migrations/        # 数据库迁移脚本
│   └── go.mod
├── mix2.yaml              # MIX2 平台配置
├── .mix2ignore            # 部署忽略文件
├── start.sh               # 一键启动脚本
└── stop.sh                # 一键停止脚本
```

## 技术栈

### 前端
- Vue 3 + TypeScript
- Element Plus (PC端UI)
- Tailwind CSS
- Pinia 状态管理
- Vue Router

### 后端
- Go 1.21+
- Gin 框架
- GORM
- SQLite / MySQL

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 后端 API | 8080 | Go 后端服务 |
| 前端服务 | 3000 | PC管理后台 + H5移动端 |

## 快速启动

### 一键启动所有服务

```bash
./start.sh
```

### 一键停止所有服务

```bash
./stop.sh
```

## 访问地址

启动服务后，可通过以下地址访问：

- **后端 API**: http://localhost:8080
- **PC 管理后台**: http://localhost:3000/admin/national
- **H5 移动端**: http://localhost:3000/m/

### 路由说明

| 路由 | 说明 |
|------|------|
| `/admin/national` | 国标政策管理 |
| `/admin/industry` | 行标政策管理 |
| `/admin/local` | 地标政策管理 |
| `/admin/category` | 政策分类管理 |
| `/admin/recommend` | 推荐政策管理 |
| `/admin/hot-update` | 政策热更新 |
| `/admin/dashboard` | 数据看板 |
| `/m/` | H5移动端首页 |
| `/m/detail/:type/:id` | H5政策详情 |

## 日志查看

```bash
# 查看后端日志
tail -f backend/server.log

# 查看前端日志
tail -f frontend/frontend.log
```

## 查看服务状态

```bash
# 查看所有服务端口
lsof -i :8080 -i :3000 | grep LISTEN
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

# 编译后端（如果需要）
cd backend && go build -o server ./cmd/server
```

## MIX2 部署

本项目符合 MIX2 平台规范，可直接部署到 MIX2 平台。

配置文件：`mix2.yaml`