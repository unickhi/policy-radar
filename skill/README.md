# MIX2 Skills

本目录包含 MIX2 平台相关的 Claude Code / OpenCode 技能定义。

## 目录结构

```
skill/
├── README.md                           # 本文件
├── mix2-builder/                       # MIX2 应用构建技能
│   ├── README.md                       # 技能说明
│   └── adapters/                       # 适配不同 AI 工具
│       ├── claude-code/                # Claude Code 格式
│       │   └── SKILL.md                # 技能定义文件
│       └── opencode/                   # OpenCode 格式（预留）
│           └── skill.md
└── ...                                 # 其他技能
```

## 安装技能

### Claude Code

```bash
# 方式1: 复制到 Claude Code skills 目录
cp -r mix2-builder/adapters/claude-code/SKILL.md ~/.claude/skills/mix2-builder/

# 方式2: 创建符号链接（推荐，便于更新）
ln -s $(pwd)/mix2-builder/adapters/claude-code ~/.claude/skills/mix2-builder
```

### OpenCode

```bash
# 待定
cp -r mix2-builder/adapters/opencode/skill.md ~/.opencode/skills/mix2-builder/
```

## 可用技能

| 技能 | 描述 | Claude Code | OpenCode |
|------|------|-------------|----------|
| mix2-builder | 创建符合 MIX2 平台规范的应用 | ✅ | 🚧 |

## 技能规范

### Claude Code 技能格式

```markdown
---
name: skill-name
description: >
  技能描述，用于自动触发。
  支持多行描述。
license: MIT
metadata:
  author: author-name
  version: "1.0"
---

# 技能标题

技能内容...
```

### 目录约定

- `adapters/` - 存放不同 AI 工具的适配器
- `claude-code/` - Claude Code 专用格式
- `opencode/` - OpenCode 专用格式

## 贡献

欢迎贡献新的技能或改进现有技能。
