# AI-Growth-Engine

AI 品牌可见度分析与增长优化 SaaS 平台。

## 项目结构

```
AI-Growth-Engine/
├── apps/                # 前端应用
│   ├── website/         # 官网 (Astro + Tailwind)
│   ├── dashboard/       # 用户端 (Next.js + Shadcn UI)
│   └── admin/           # 管理后台 (Next.js + Shadcn Admin)
├── services/            # 后端服务 (Go)
├── packages/            # 共享包
├── infrastructure/      # 基础设施配置
├── deploy/              # 部署文件
├── docs/                # 项目文档
├── tests/               # 测试
└── scripts/             # 脚本
```

## 技术栈

- **前端**: Next.js / React / TypeScript / TailwindCSS / Shadcn UI / Tremor
- **官网**: Astro + Tailwind CSS
- **后端**: Go / PostgreSQL / Redis
- **包管理**: pnpm workspace + Turborepo (前端), Go workspace (后端)

## 快速开始

```bash
# 安装依赖
pnpm install

# 启动基础设施
make infra

# 启动开发服务
pnpm dev:dashboard
```

## 项目文档

详细文档见 `AI-Growth-Engine-Knowledg/` 目录。
