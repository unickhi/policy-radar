# MIX2 Builder Skill

创建符合 MIX2 平台规范的应用程序。

## 功能

- 自动生成 MIX2 标准项目结构
- 配置 mix2.yaml 部署文件
- 生成数据库迁移脚本
- 遵循技术栈最佳实践

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Tailwind CSS |
| 后端 | Go + Gin + GORM |
| 数据库 | MySQL 8.0 |

## 安装

### Claude Code

```bash
# 方式1: 符号链接（推荐，便于更新）
ln -s $(pwd)/adapters/claude-code ~/.claude/skills/mix2-builder

# 方式2: 复制文件
mkdir -p ~/.claude/skills/mix2-builder
cp adapters/claude-code/SKILL.md ~/.claude/skills/mix2-builder/
```

---

## 使用指南

### 1. 创建应用

在 Claude Code 中直接描述需求：

```
帮我创建一个 MIX2 应用：任务管理系统
```

```
开发一个博客系统，需要文章管理和评论功能
```

```
创建一个库存管理系统，前后端分离，使用 Go 和 Vue
```

Skill 会自动触发，生成：
- 标准项目目录结构
- mix2.yaml 配置文件
- 数据库迁移脚本
- 基础代码模板

### 2. 应用类型

支持两种应用类型：

#### 全栈应用（前后端分离）

```
app-name/
├── frontend/      # Vue 3 前端
├── backend/       # Go 后端
├── mix2.yaml
└── README.md
```

#### 纯前端应用

```
app-name/
├── src/           # Vue 3 源码
├── mix2.yaml
└── README.md
```

### 3. 应用打包

应用开发完成后，打包源码上传到 MIX2 平台：

#### 创建 .mix2ignore

```bash
# 在应用根目录创建
cat > .mix2ignore << 'EOF'
node_modules/
dist/
build/
.git/
*.log
.env
.env.local
.idea/
.vscode/
coverage/
*.test.js
*.spec.js
EOF
```

#### 打包命令

```bash
# 在应用上级目录执行
cd /path/to/apps

# 打包（排除不必要文件）
tar -czf my-app.tar.gz \
  --exclude='my-app/node_modules' \
  --exclude='my-app/dist' \
  --exclude='my-app/.git' \
  --exclude='my-app/*.log' \
  my-app/

# 或使用 .mix2ignore
tar -czf my-app.tar.gz -X my-app/.mix2ignore my-app/
```

### 4. 发布到 MIX2

#### 方式1: Dashboard 上传

1. 访问 MIX2 Dashboard (如 http://localhost:5500)
2. 点击「创建应用」
3. 上传 `.tar.gz` 源码包
4. 平台自动构建和部署

#### 方式2: CLI 部署（计划中）

```bash
cd my-app
mix2 deploy
```

#### 方式3: API 发布

```bash
curl -X POST http://localhost:5500/api/v1/apps \
  -H "Authorization: Bearer $TOKEN" \
  -F "package=@my-app.tar.gz"
```

### 5. MIX2 自动构建

上传源码后，MIX2 自动处理：

| 检测文件 | 运行时 | 构建命令 |
|----------|--------|----------|
| `Dockerfile` | Docker | `docker build` |
| `go.mod` | Go | `go build -o app ./cmd/server` |
| `package.json` | Node.js | `npm install && npm run build` |
| 其他 | 静态站点 | 直接服务 |

---

## 常见场景示例

### 场景1: 内容管理系统

```
创建一个内容管理系统，包含：
- 文章管理（CRUD）
- 分类管理
- 标签管理
- 用户认证（集成 MIX2 SSO）
```

### 场景2: 数据分析仪表盘

```
创建一个数据分析仪表盘应用：
- 纯前端（Vue 3）
- 使用 ECharts 展示图表
- 支持数据筛选和导出
```

### 场景3: API 服务

```
创建一个 RESTful API 服务：
- 纯后端（Go）
- 提供 JSON API
- 需要 MySQL 数据库
```

---

## 检查清单

发布前确认：

- [ ] `mix2.yaml` 配置完整
- [ ] 数据库迁移脚本 `.up.sql` 和 `.down.sql` 成对
- [ ] 共享数据库使用 `${prefix}` 表前缀
- [ ] 无硬编码配置（使用环境变量）
- [ ] 健康检查端点 `/api/health` 可用
- [ ] `.mix2ignore` 文件存在
- [ ] `node_modules/` 和 `dist/` 未打包

---

## 适配器

| 工具 | 状态 | 文件 |
|------|------|------|
| Claude Code | ✅ 可用 | `adapters/claude-code/SKILL.md` |
| OpenCode | 🚧 计划中 | `adapters/opencode/skill.md` |

---

## 相关链接

- [MIX2 平台文档](../docs/)
- [技术规范](../docs/APP_SPEC.md)
- [API 文档](../docs/API.md)
