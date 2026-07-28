# First 12 Production Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复用户反馈的 1-12 项，让已有管理端、用户端和官网功能由真实数据与真实配置驱动。

**Architecture:** 先修复 Go API 的查询字段、业务校验和事务边界，再同步前端类型与页面状态。任务积分在数据库事务内按有效模型数量扣除，统计与分析只聚合真实记录。

**Tech Stack:** Go、Chi、PostgreSQL、Next.js、React、TypeScript、Astro、Tailwind CSS、pnpm。

## Global Constraints

- 所有接口遵循 `/api/v1` 和 `{success,data,message}` 返回格式。
- 所有用户可见错误使用简体中文。
- 保留已有功能；功能已存在时修复其真实数据链路并跳过重复实现。
- 遵循 `docs/UI-UX设计规范.md` 的 Enterprise SaaS 视觉规范。
- 不读取、输出或写入 Agent 环境中的大模型密钥。

---

### Task 1: 修复任务、订单与模型统计数据契约

**Files:**
- Modify: `services/task/model.go`
- Modify: `services/task/repository.go`
- Modify: `services/order/repository.go`
- Modify: `services/api/internal/handler/admin_task_handler.go`
- Modify: `services/api/internal/handler/aimodel_handler.go`
- Test: `services/task/repository_test.go`
- Test: `services/order/repository_test.go`

**Interfaces:**
- `task.Repository.Delete(ctx, id)` SHALL remove dependent `ai_answers` rows before removing `ai_tasks`.
- Admin task list SHALL expose `project_name`, `username`, `progress`, `questions_count`, and `completed_count`.
- Model stats SHALL expose numeric `avg_latency_ms` and use configured model names where available.

- [ ] 写失败测试，覆盖任务进度边界、订单字段扫描映射和统计空数据。
- [ ] 修正 `taskCols`、任务列表 JOIN、订单列表 SELECT 与 Scan 字段数量。
- [ ] 删除任务时在事务中删除答案并删除任务。
- [ ] 将统计查询增加平均延迟并记录真实任务模型标识。
- [ ] 运行 `go test ./services/task ./services/order ./services/api/...`，确认通过。

### Task 2: 实现多模型计费与真实进度

**Files:**
- Modify: `services/user/repository.go`
- Modify: `services/task/model.go`
- Modify: `services/task/repository.go`
- Modify: `services/api/internal/handler/task_handler.go`
- Modify: `services/api/internal/handler/project_handler.go`
- Test: `services/api/internal/handler/task_handler_test.go`

**Interfaces:**
- Add `NormalizeModels(model string) []string` for comma-separated model IDs, trimming blanks and removing duplicates.
- Add `CalculateCreditCost(models []string) int`, returning model count.
- Task creation SHALL call one transactional repository operation that validates credit balance, deducts `len(models)`, and creates the task.

- [ ] 写失败测试，验证单模型扣 1、多模型按数量扣除、余额不足不扣费。
- [ ] 实现模型解析与积分计算，拒绝空模型列表。
- [ ] 将余额不足消息改为“积分不足，请先充值”。
- [ ] 让创建失败自动回滚积分，避免先扣费后写入失败。
- [ ] 用已完成数/问题总数计算 0-100 进度并返回前端。
- [ ] 运行对应 Go 测试与 `go vet ./...`。

### Task 3: 修复 CDK 与支付配置

**Files:**
- Modify: `services/cdk/model.go`
- Modify: `services/cdk/repository.go`
- Modify: `services/api/internal/handler/cdk_handler.go`
- Modify: `services/payment/model.go`
- Modify: `services/payment/repository.go`
- Add: `services/database/migrations/000015_add_payment_channel.sql`
- Test: `services/cdk/repository_test.go`
- Test: `services/payment/repository_test.go`

**Interfaces:**
- CDK create accepts optional `code`; empty code generates a unique uppercase code.
- CDK credits SHALL be positive; max uses SHALL be zero or positive; duplicate code returns HTTP 409 with Chinese message.
- Payment config adds `channel`, defaulting to `alipay`; Epay request uses the configured channel.

- [ ] 写失败测试，覆盖 CDK 校验、生成代码、重复代码和渠道参数。
- [ ] 增加支付渠道迁移与读写字段。
- [ ] 在 handler 层返回中文校验错误和冲突状态。
- [ ] 将易支付 URL 生成从固定 `alipay` 改为配置渠道。
- [ ] 运行 CDK、payment 测试和数据库迁移静态检查。

### Task 4: 修复管理端页面

**Files:**
- Modify: `apps/admin/src/app/(admin)/models/stats/page.tsx`
- Modify: `apps/admin/src/app/(admin)/tasks/page.tsx`
- Modify: `apps/admin/src/app/(admin)/orders/page.tsx`
- Modify: `apps/admin/src/app/(admin)/cdk/page.tsx`
- Modify: `apps/admin/src/app/(admin)/payment/page.tsx`
- Modify: `apps/admin/src/app/(admin)/settings/page.tsx`

**Interfaces:**
- Frontend types match API snake_case fields and derive progress from API values.
- Payment UI exposes `alipay`, `wxpay`, `qqpay`, and `bank` channels.
- SMS provider select exposes `smsbao`; its username/password fields appear when selected.

- [ ] 更新模型统计的延迟字段和 Token 显示，空值显示“暂无数据”。
- [ ] 删除与重试操作显示接口错误，成功后刷新当前列表。
- [ ] 修正订单字段映射，安全处理金额、用户名、套餐和时间。
- [ ] CDK 表单允许自动生成代码并校验额度与次数。
- [ ] 支付表单保存 `channel`，短信宝进入统一服务商选择。
- [ ] 运行 `pnpm --filter @aige/admin typecheck` 与 `pnpm --filter @aige/admin build`。

### Task 5: 修复用户端积分提示、仪表盘与分析

**Files:**
- Modify: `apps/dashboard/src/app/(protected)/projects/[id]/tasks/new/page.tsx`
- Modify: `apps/dashboard/src/app/(protected)/dashboard/page.tsx`
- Modify: `apps/dashboard/src/app/(protected)/analytics/page.tsx`
- Modify: `apps/dashboard/src/lib/api.ts`

**Interfaces:**
- New task page sends normalized selected model IDs and displays `selectedModelIds.length` as required credits.
- API error mapping converts balance errors to “积分不足，请先充值”。
- Dashboard and analytics preserve real empty states and reload on selected date range.

- [ ] 写模型选择和积分提示的组件级测试或静态类型约束。
- [ ] 阻止零模型提交，显示中文校验提示。
- [ ] 修正任务接口返回分页结构的读取，避免把空数据误判为假数据。
- [ ] 让分析页按 7/30/90 天重新请求并按任务报告聚合。
- [ ] 去除无报告时使用统计默认分数填充趋势的行为。
- [ ] 运行 `pnpm --filter @aige/dashboard typecheck` 与 build。

### Task 6: 主题切换与官网 Host 配置

**Files:**
- Modify: `apps/website/astro.config.ts`
- Modify: `apps/website/src/layouts/Layout.astro`
- Modify: `apps/dashboard/src/app/(protected)/layout.tsx`
- Modify: `apps/admin/src/app/(admin)/layout.tsx`
- Add: `apps/website/src/components/ThemeToggle.astro`
- Add: `apps/dashboard/src/components/ThemeToggle.tsx`
- Add: `apps/admin/src/components/ThemeToggle.tsx`

**Interfaces:**
- Theme toggle persists `light` or `dark` in `localStorage` under `aige-theme`.
- Astro server accepts `.monkeycode-ai.online` and the explicit preview host.

- [ ] 写主题值解析测试或共享常量校验。
- [ ] 保留现有设计 token，为三端导航增加主题入口。
- [ ] 增加 `dark` class 初始化脚本，避免刷新闪烁。
- [ ] 在 Astro allowedHosts 中同时保留通配域名和当前预览 Host。
- [ ] 运行 website、dashboard、admin 三端 build。

### Task 7: 全量验证与文档同步

**Files:**
- Modify: `.env.example` only when a required user-facing setting is missing
- Modify: `.monkeycode/specs/first-12-production-fixes/design.md` when implementation details differ

- [ ] 运行 `gofmt -w` 处理 Go 修改文件。
- [ ] 运行 `go test ./...`、`go vet ./...`。
- [ ] 运行 `pnpm typecheck`、`pnpm lint` 和三个 build 命令。
- [ ] 使用 `deploy-website` 启动可预览 Web 应用并检查官网 Host 配置。
- [ ] 检查 git diff，确认只包含本阶段功能与设计文档变更。
