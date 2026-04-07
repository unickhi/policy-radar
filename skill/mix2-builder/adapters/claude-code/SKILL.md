---
name: mix2-builder
description: >
  Create applications conforming to MIX2 platform specifications.
  Use when creating MIX2 apps, developing MIX2 platform applications,
  or generating projects deployable to MIX2.
license: MIT
metadata:
  author: iceskysl
  version: "1.0"
  homepage: https://github.com/mix2/infra
---

# MIX2 Application Builder

> 创建符合 MIX2 平台规范的应用程序

## 触发条件

当用户请求以下操作时应用此技能：
- 创建 MIX2 应用
- 开发可部署到 MIX2 的项目
- 生成符合 MIX2 规范的代码

---

## 技术栈规范

### 前端（推荐）

| 组件 | 技术 | 版本 |
|------|------|------|
| 框架 | Vue 3 + TypeScript | ^3.4 |
| 样式 | Tailwind CSS | ^3.4 |
| 构建 | Vite | ^5.0 |
| 状态 | Pinia | ^2.1 |
| 路由 | Vue Router | ^4.2 |

### 后端（推荐）

| 组件 | 技术 | 版本 |
|------|------|------|
| 语言 | Go | ^1.21 |
| 框架 | Gin | ^1.9 |
| ORM | GORM | ^1.25 |

### 数据库

| 类型 | 技术 | 版本 |
|------|------|------|
| 主数据库 | MySQL | 8.0 |

---

## 项目结构

### 全栈应用

```
app-name/
├── frontend/               # Vue 3 前端
│   ├── src/
│   │   ├── api/           # API 封装
│   │   ├── components/    # 组件
│   │   ├── pages/         # 页面
│   │   ├── router/        # 路由
│   │   ├── stores/        # Pinia 状态
│   │   ├── styles/        # 样式
│   │   ├── types/         # TS 类型
│   │   ├── App.vue
│   │   └── main.ts
│   ├── index.html
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── package.json
├── backend/                # Go 后端
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── handler/       # HTTP 处理器
│   │   ├── service/       # 业务逻辑
│   │   ├── repository/    # 数据访问
│   │   ├── model/         # 数据模型
│   │   ├── middleware/    # 中间件
│   │   └── config/        # 配置
│   ├── migrations/
│   ├── go.mod
│   └── go.sum
├── mix2.yaml
└── README.md
```

### 纯前端应用

```
app-name/
├── src/
│   ├── api/
│   ├── components/
│   ├── pages/
│   ├── router/
│   ├── stores/
│   ├── App.vue
│   └── main.ts
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── package.json
└── mix2.yaml
```

---

## mix2.yaml 配置

### 最小配置

```yaml
name: my-app
```

### 完整配置

```yaml
apiVersion: mix2/v1
kind: Application
metadata:
  name: my-app
  version: 1.0.0
  description: 应用描述

spec:
  runtime:
    type: nodejs      # nodejs | golang | static
    version: "20"

  build:
    command: npm run build
    output: dist

  run:
    command: npm start
    port: 3000

  resources:
    cpu: "0.5"
    memory: "512Mi"

  env:
    NODE_ENV: production

  services:
    - name: mysql
      type: mysql
      version: "8.0"
      isolation: shared
      tablePrefix: myapp

  migrations:
    enabled: true
    path: migrations/
    autoRun: true

  ingress:
    enabled: true
    type: path

  auth:
    enabled: false
```

---

## 数据库迁移

### 文件结构

```
migrations/
├── 001_init_schema.up.sql
├── 001_init_schema.down.sql
├── 002_add_users.up.sql
└── 002_add_users.down.sql
```

### 共享数据库模式

使用 `${prefix}` 占位符：

```sql
-- 001_init_schema.up.sql
CREATE TABLE ${prefix}_users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ${prefix}_posts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES ${prefix}_users(id)
);
```

```sql
-- 001_init_schema.down.sql
DROP TABLE IF EXISTS ${prefix}_posts;
DROP TABLE IF EXISTS ${prefix}_users;
```

---

## API 规范

### RESTful 路由

```
GET    /api/v1/resources       # 列表
GET    /api/v1/resources/:id   # 详情
POST   /api/v1/resources       # 创建
PUT    /api/v1/resources/:id   # 更新
DELETE /api/v1/resources/:id   # 删除
```

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误响应

```json
{
  "code": 40001,
  "message": "参数错误",
  "errors": ["email 格式不正确"]
}
```

---

## 代码模板

### Vue 3 组件

```vue
<template>
  <div class="p-4">
    <h1 class="text-2xl font-bold mb-4">{{ title }}</h1>
    <div v-if="loading" class="text-gray-500">加载中...</div>
    <div v-else>
      <!-- 内容 -->
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Props {
  title: string
}
const props = defineProps<Props>()

const loading = ref(false)
const data = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    // API 调用
  } finally {
    loading.value = false
  }
}

onMounted(() => fetchData())
</script>
```

### Go Handler

```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ExampleHandler struct {
    service *service.ExampleService
}

func NewExampleHandler(s *service.ExampleService) *ExampleHandler {
    return &ExampleHandler{service: s}
}

func (h *ExampleHandler) List(c *gin.Context) {
    items, err := h.service.List(c.Request.Context())
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "code":    500,
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "code":    0,
        "message": "success",
        "data": gin.H{"items": items},
    })
}
```

### Go Model

```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Email     string         `gorm:"uniqueIndex;size:255;not null" json:"email"`
    Name      string         `gorm:"size:100" json:"name"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
    return "${prefix}_users"
}
```

---

## 打包发布

### 源码打包

```bash
tar -czf my-app.tar.gz \
  --exclude='my-app/node_modules' \
  --exclude='my-app/dist' \
  --exclude='my-app/.git' \
  my-app/
```

### .mix2ignore

```
node_modules/
dist/
build/
.git/
*.log
.env
.env.local
.idea/
.vscode/
```

---

## 检查清单

创建应用前必须确认：

- [ ] 使用标准技术栈
- [ ] 项目结构符合规范
- [ ] mix2.yaml 配置完整
- [ ] 数据库迁移脚本完整
- [ ] 共享模式使用 `${prefix}` 表前缀
- [ ] API 遵循 RESTful 规范
- [ ] 响应使用统一格式
- [ ] 无硬编码配置（使用环境变量）
- [ ] 健康检查端点 `/api/health`

---

## 工作流程

1. **需求分析** - 确认应用类型和功能需求
2. **创建结构** - 按规范生成目录和配置
3. **实现功能** - 前端 Vue 3 + 后端 Go
4. **配置部署** - 编写 mix2.yaml 和迁移脚本
5. **验证** - 检查结构和配置完整性
