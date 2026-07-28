# AI-Growth-Engine Complete Implementation Plan

基于对全部 14 份文档和现有代码库的完整分析生成。

---

## 当前代码库状态 SUMMARY

### 已完成部分

| 模块 | 完成度 | 详情 |
|------|--------|------|
| 数据库迁移 | 3/15+ | `users`, `verification_codes`, `audit_logs`, `brand_projects`, `ai_tasks`, `ai_answers`, `admins` |
| Go API 路由 | 基本完成 | chi router, JWT middleware, 7 handlers |
| API Handlers | 7/20+ | auth, user, project, task, dashboard, admin auth, admin CRUD |
| Dashboard 前端页面 | 14 pages | login, register, forgot-password, dashboard, analytics, profile, projects CRUD, tasks CRUD, reports |
| Admin 前端页面 | 6 pages | login, dashboard, users, orders, models, system |
| AI Engine providers | 4 个目录 | openai/, claude/, gemini/, kimi/ (go.mod 存在, 内容待确认) |
| Website (Astro) | 0 pages | 空目录 |

### 空/未实现模块

| 模块 | 状态 |
|------|------|
| AI Engine analyzer/ | 空目录 |
| AI Engine scorer/ | 空目录 |
| AI Engine generator/ | 空目录 |
| AI Engine prompts/ | 空目录 |
| services/report/ | 仅 go.mod |
| services/billing/ | 仅 go.mod |
| services/scheduler/ | 仅 go.mod |
| services/notification/ | 仅 go.mod |
| services/worker/ | 目录存在但内容待确认 |
| Website (Astro) pages | 完全空 |

---

## Phase 1: Critical Fixes (do first)

### BUG-001: Missing database tables
**严重性**: CRITICAL
**描述**: 文档要求大量表(plans, subscriptions, orders, competitors, reports, notifications, prompts, spaces 等)但有且仅有 3 个迁移文件。
**位置**: `/workspace/services/database/migrations/`
**任务**: 创建至少以下迁移文件:
- `000004_create_plans.sql` - SaaS 套餐表
- `000005_create_subscriptions.sql` - 用户订阅表
- `000006_create_orders.sql` - 订单表
- `000007_create_competitors.sql` - 竞争分析表
- `000008_create_reports.sql` - 报告表
- `000009_create_notifications.sql` - 通知表
- `000010_create_prompts.sql` - Prompt 模板表
- `000011_create_spaces.sql` - 工作空间/租户表
- `000012_create_api_keys.sql` - API Key 表
- `000013_create_brand_infos.sql` - 品牌资料库表

### BUG-002: Missing API endpoints
**严重性**: CRITICAL
**描述**: router.go 只定义了 auth、users、projects、tasks、dashboard、admin 的有限路由。缺少:
- 竞争分析 API (`/api/v1/competitors`)
- 报告生成/查看 API (`/api/v1/reports`)
- 订阅/计费 API (`/api/v1/billing`)
- 通知 API (`/api/v1/notifications`)
- 空间/团队 API (`/api/v1/spaces`)
- AI 引擎管理 API
- Prompt 管理 API (admin)
- 操作日志 API (admin)
**位置**: `/workspace/services/api/internal/router/router.go`

### BUG-003: Dashboard 页面可能是空壳
**严重性**: HIGH
**描述**: 14 个 page.tsx 文件存在，但需要逐个验证内容是否完整实现。需要检查每个 page 是否:
- 调用了真实 API
- 实现了完整的 UI 组件
- 有 loading/error/empty 状态处理

### BUG-004: Admin 页面不完整
**严重性**: HIGH
**描述**: admin 只有 6 个页面，缺少:
- `prompts/` - Prompt 管理 (feature 目录存在但无 page)
- `tasks/` - 任务管理 (feature 目录存在但无 page)
- `logs/` - 日志管理 (feature 目录存在但无 page)

### BUG-005: AI Engine 核心逻辑缺失
**严重性**: CRITICAL
**描述**: `analyzer/`, `scorer/`, `generator/`, `prompts/` 全部为空。这是产品的核心竞争力。
**影响**: 整个 AI 分析管道无法工作。

### BUG-006: 报告服务完全未实现
**严重性**: CRITICAL
**描述**: `services/report/` 仅含 go.mod。无报告生成、PDF 导出、报告模板等。
**影响**: 用户无法获得核心产出 (AI 诊断报告)。

### BUG-007: 商业化系统完全未实现
**严重性**: CRITICAL
**描述**: `services/billing/` 仅含 go.mod。无套餐管理、订阅、支付、订单。
**影响**: 无法收费，无商业模式验证。

### BUG-008: 网站 homepage 完全缺失
**严重性**: HIGH
**描述**: `apps/website/src/pages/` 为空。官网是产品入口，Landing Page 必须建设。

### BUG-009: Worker 和 Scheduler 未实现
**严重性**: HIGH
**描述**: `services/worker/` 有目录结构但任务处理逻辑未确认；`services/scheduler/` 仅含 go.mod。
**影响**: AI 检测无法异步执行，自动监控不可用。

### BUG-010: Notification 服务未实现
**严重性**: MEDIUM
**描述**: `services/notification/` 仅含 go.mod。无站内通知、邮件通知。

### UI-001: UI 组件库状态未知
**严重性**: MEDIUM
**描述**: `packages/ui/` 目录存在但内容未知。需要验证 Shadcn UI 组件是否已正确初始化。

### UI-002: 缺少统一的错误处理
**严重性**: MEDIUM
**描述**: 前端页面缺少 loading/error/empty 的统一处理模式。

---

## Phase 2: Feature Implementation (按功能清单)

引用文档: 16-功能清单.md

### 模块一：账户系统 (AUTH)

#### F-001: AUTH-001~005 — 用户注册登录 (P0)
**状态**: 部分完成
**后端**: 
- [x] 注册 API (`POST /api/v1/auth/register`)
- [x] 登录 API (`POST /api/v1/auth/login`)
- [x] 发送验证码 (`POST /api/v1/auth/send-code`)
- [x] 重置密码 (`POST /api/v1/auth/reset-password`)
- [ ] 第三方登录 (OAuth/GitHub/Google)
- [ ] 登录设备管理 API
- [ ] 双因素认证 API
**前端**:
- [x] 登录页面
- [x] 注册页面
- [x] 忘记密码页面
- [ ] 登录设备管理页面
- [ ] 双因素认证设置页面
**复杂度**: Small (for remaining items)
**依赖**: 无

### 模块二：工作空间系统 (SPACE)

#### F-002: SPACE-001~005 — 工作空间/团队 (P0-P1)
**状态**: 未实现
**后端**:
- [ ] 创建工作空间 API
- [ ] 空间切换 API
- [ ] 成员邀请 API
- [ ] 成员管理 CRUD API
- [ ] 权限角色管理 API
**前端**:
- [ ] 空间创建页面/弹窗
- [ ] 空间切换器 (Sidebar)
- [ ] 成员管理页面
- [ ] 权限设置页面
**复杂度**: Medium
**依赖**: F-001 (用户系统)

### 模块三：品牌项目管理 (PROJECT)

#### F-003: PROJECT-001~005 — 品牌项目 CRUD (P0)
**状态**: 部分完成
**后端**:
- [x] 创建项目 (`POST /api/v1/projects`)
- [x] 项目列表 (`GET /api/v1/projects`)
- [x] 项目详情 (`GET /api/v1/projects/{id}`)
- [x] 更新项目 (`PUT /api/v1/projects/{id}`)
- [x] 删除项目 (`DELETE /api/v1/projects/{id}`)
- [ ] 品牌资料库 API (品牌介绍/FAQ/案例)
**前端**:
- [x] 项目列表页面
- [x] 创建项目页面
- [x] 编辑项目页面
- [x] 项目详情页面
- [ ] 品牌资料库 UI (资料填写/管理)
**复杂度**: Small
**依赖**: F-001

### 模块四：AI 检测任务 (TASK)

#### F-004: TASK-001~004 — AI 检测任务 (P0)
**状态**: 部分完成
**后端**:
- [x] 创建任务 (`POST /api/v1/tasks`)
- [x] 任务列表 (`GET /api/v1/tasks`)
- [x] 任务详情 (`GET /api/v1/tasks/{id}`)
- [x] 任务执行 (`POST /api/v1/tasks/{id}/execute`)
- [x] 获取报告 (`GET /api/v1/tasks/{id}/report`)
- [ ] 自动生成问题 (调用 AI Engine)
- [ ] 自定义问题 API
- [ ] 批量检测 API
**前端**:
- [x] 任务列表页面
- [x] 任务详情页面
- [x] 新建任务页面
- [ ] 自定义问题输入 UI
- [ ] 批量检测配置 UI
**复杂度**: Medium
**依赖**: F-005 (AI Engine)

### 模块五：AI 模型接入 (MODEL)

#### F-005: MODEL-001~004 — AI 模型管理与调用 (P0)
**状态**: 部分完成
**后端**:
- [x] Provider 接口定义
- [x] OpenAI provider
- [x] Claude provider
- [x] Gemini provider
- [x] Kimi provider
- [ ] 模型管理 API (CRUD)
- [ ] 模型选择策略
- [ ] AI 调用记录完整实现
**前端** (Admin):
- [x] 模型管理页面 (admin)
- [ ] 调用统计可视化
**复杂度**: Medium
**依赖**: 无

### 模块六：AI Agent 系统 (AGENT)

#### F-006: AGENT-001~006 — AI Agents (P0-P1)
**状态**: 未实现
**后端** (`services/ai-engine/`):
- [ ] AGENT-001: 问题生成 Agent (在 analyzer/ 或 generator/ 中实现)
- [ ] AGENT-002: 查询执行 Agent (复用 providers + worker)
- [ ] AGENT-003: 品牌识别 Agent (在 analyzer/ 中实现)
- [ ] AGENT-004: 竞争分析 Agent (在 analyzer/ 中实现)
- [ ] AGENT-005: 内容分析 Agent (P1)
- [ ] AGENT-006: 优化建议 Agent (P1)
**复杂度**: Large
**依赖**: F-005 (AI providers)

### 模块七：Prompt 管理系统 (PROMPT)

#### F-007: PROMPT-001~003 — Prompt 管理 (P0-P1)
**状态**: 未实现
**后端**:
- [ ] Prompt 模板 CRUD API
- [ ] Prompt 版本控制 API
- [ ] Prompt 测试 API
**前端** (Admin):
- [ ] Prompt 管理页面
- [ ] Prompt 版本对比
- [ ] Prompt 测试台
**复杂度**: Medium
**依赖**: F-005

### 模块八：AI 分析引擎 (ANALYSIS)

#### F-008: ANALYSIS-001~006 — 分析引擎 (P0-P1)
**状态**: 未实现
**后端** (in `services/ai-engine/analyzer/`):
- [ ] ANALYSIS-001: 品牌出现检测 (精确匹配/模糊匹配/同义词)
- [ ] ANALYSIS-002: 品牌曝光统计 (出现次数/比例/模型覆盖)
- [ ] ANALYSIS-003: 推荐排名分析 (第一推荐/列表/仅提及)
- [ ] ANALYSIS-004: 情感分析 (正面/中性/负面) (P1)
- [ ] ANALYSIS-005: 引用来源分析 (P1)
- [ ] ANALYSIS-006: 行业排名分析 (P1)
**复杂度**: Large
**依赖**: F-006 (AI Agents)

### 模块九：评分系统 (SCORE)

#### F-009: SCORE-001~003 — AI Visibility Score (P0-P1)
**状态**: 未实现
**后端** (in `services/ai-engine/scorer/`):
- [ ] SCORE-001: 总评分生成 (0-100)
- [ ] SCORE-002: 评分权重模型 (出现率/排名/覆盖率/内容质量/竞争/趋势)
- [ ] SCORE-003: 评分趋势计算 (日/7天/30天/年)
**前端**:
- [ ] Dashboard 评分展示组件
- [ ] 趋势图表 (使用 Tremor)
**复杂度**: Medium
**依赖**: F-008 (分析引擎)

### 模块十：数据中心 (DATA)

#### F-010: DATA-001~003 — 数据采集/清洗/标签 (P0-P1)
**状态**: 未实现
**后端**:
- [ ] DATA-001: AI 回答数据采集管道
- [ ] DATA-002: 数据清洗 (去重/无效回答/错误识别)
- [ ] DATA-003: 数据标签 (行业/地区/模型/时间/类型)
**复杂度**: Medium
**依赖**: F-005, F-008

### 模块十一：Dashboard 数据中心 (DASH)

#### F-011: DASH-001~004 — Dashboard (P0-P1)
**状态**: 部分完成
**后端**:
- [x] Dashboard stats API (`GET /api/v1/dashboard/stats`)
**前端**:
- [x] Dashboard 页面 (需要验证完整性)
- [ ] 品牌趋势图 (Tremor LineChart)
- [ ] 模型表现分析面板 (P1)
- [ ] 竞争分析面板 (P1)
**复杂度**: Medium
**依赖**: F-009 (Score)

### 模块十二：报告系统 (REPORT)

#### F-012: REPORT-001~006 — 报告系统 (P0-P2)
**状态**: 未实现
**后端** (in `services/report/`):
- [ ] REPORT-001: AI 分析报告自动生成
- [ ] REPORT-002: 报告模板管理 (企业/学校/电商/SaaS/个人)
- [ ] REPORT-003: 报告在线查看 API
- [ ] REPORT-004: 报告导出 (PDF/Excel/Markdown/HTML)
- [ ] REPORT-005: 报告分享 (公开链接/密码/有效期)
- [ ] REPORT-006: 白标报告 (P2, 企业版)
**前端**:
- [ ] 报告查看页面 (已存在 route 需验证)
- [ ] 报告列表/筛选
- [ ] 报告导出按钮
- [ ] 报告分享弹窗
**复杂度**: Large
**依赖**: F-006, F-008, F-009

### 模块十三：监控系统 (MONITOR)

#### F-013: MONITOR-001~004 — 自动监控 (P0-P1)
**状态**: 未实现
**后端** (in `services/scheduler/`):
- [ ] MONITOR-001: 自动检测计划 (每天/每周/每月)
- [ ] MONITOR-002: 数据变化检测 (评分/排名/竞争/AI回答)
- [ ] MONITOR-003: 异常检测 (曝光下降/竞争超过)
- [ ] MONITOR-004: 历史数据追踪 (每日快照/评分变化)
**前端**:
- [ ] 监控计划设置页面
- [ ] 变化提醒配置
**复杂度**: Medium
**依赖**: F-008, F-012

### 模块十四：通知系统 (NOTICE)

#### F-014: NOTICE-001~003 — 通知 (P0-P2)
**状态**: 未实现
**后端** (in `services/notification/`):
- [ ] NOTICE-001: 站内通知 (任务完成/报告生成/评分变化)
- [ ] NOTICE-002: 邮件通知 (注册/报告/异常) (P1)
- [ ] NOTICE-003: Webhook 通知 (P2, 企业版)
**前端**:
- [ ] 通知中心页面/下拉菜单
- [ ] 通知列表/标记已读
- [ ] 通知设置页面
**复杂度**: Medium
**依赖**: F-001

### 模块十五：搜索与筛选系统 (SEARCH)

#### F-015: SEARCH-001~002 — 搜索筛选 (P1)
**状态**: 未实现
**后端**:
- [ ] 全局搜索 API (项目/报告/任务/品牌)
- [ ] 高级筛选 API (时间/行业/模型/状态/评分)
**前端**:
- [ ] 全局搜索框 (Header)
- [ ] 高级筛选组件
**复杂度**: Small
**依赖**: F-003, F-004, F-012

### 模块十七~十八：权限与审计 (PERMISSION/AUDIT)

#### F-016: PERMISSION-001~003 + AUDIT-001~002 (P0-P1)
**状态**: 未实现
**后端**:
- [ ] RBAC 权限模型 (SuperAdmin/Admin/Owner/Member/Viewer)
- [ ] 数据权限隔离
- [ ] API 权限控制
- [ ] 操作审计日志
- [ ] 安全审计 (登录/API/权限异常)
**前端** (Admin):
- [ ] 权限管理页面
- [ ] 审计日志查看页面
**复杂度**: Medium (RBAC 是基础结构)
**依赖**: F-002 (Space)

### 模块十九：商业化系统 (BILLING)

#### F-017: BILLING-001~006 — 计费订阅 (P0-P2)
**状态**: 未实现
**后端** (in `services/billing/`):
- [ ] BILLING-001: 套餐管理 API (Free/Professional/Business/Enterprise)
- [ ] BILLING-002: 订阅管理 API (升级/降级/取消)
- [ ] BILLING-003: 用量计费 (AI 调用/Token/报告/API)
- [ ] BILLING-004: 订单系统
- [ ] BILLING-005: 支付系统 (微信/支付宝/Stripe) (P1)
- [ ] BILLING-006: 发票系统 (P2)
**前端**:
- [ ] 套餐选择页面 (Plan comparison)
- [ ] 订阅管理页面
- [ ] 用量查看
- [ ] 订单历史
- [ ] 发票申请 (P2)
**复杂度**: Large
**依赖**: F-001

### 模块二十~二十一：开放平台 (API/DEV)

#### F-018: API-001~005 + DEV-001~003 (P1-P3)
**状态**: 未实现
**后端**:
- [ ] API Key 管理 (创建/删除/重置)
- [ ] API 调用统计
- [ ] 品牌查询 API (评分/报告/趋势)
- [ ] AI 检测 API (提交问题返回分析)
- [ ] Webhook 系统 (P2)
- [ ] 开发者资料管理 (P2)
- [ ] SDK (P3)
**前端**:
- [ ] API Key 管理页面
- [ ] API 文档中心 (P2)
**复杂度**: Large
**依赖**: F-001, F-017

---

## Phase 3: Admin Panel Completion

### Admin 当前页面: 6 pages exist
- dashboard (/) - 需要验证完整性
- users/ - 用户管理
- orders/ - 订单管理
- models/ - AI 模型管理
- system/ - 系统设置

### Admin 缺失页面:

#### ADM-001: Admin Dashboard 数据完善 (P0)
**描述**: 展示用户数量、企业数量、任务数量、AI 调用次数、收入、系统状态
**后端**: 需要 admin dashboard stats API
**位置**: `/workspace/apps/admin/src/app/(admin)/page.tsx`

#### ADM-002: Prompt 管理页面 (P0)
**描述**: Prompt 模板管理、版本控制、在线测试
**前端**: 需要创建 page.tsx
**位置**: `/workspace/apps/admin/src/app/(admin)/prompts/page.tsx`
**后端**: 需要 prompt CRUD API 路由

#### ADM-003: 任务管理页面 (P1)
**描述**: 查看运行/失败/等待/历史任务，支持重试和删除
**前端**: 需要创建 page.tsx
**位置**: `/workspace/apps/admin/src/app/(admin)/tasks/page.tsx`
**后端**: 需要 admin task management API 路由

#### ADM-004: 操作日志页面 (P0)
**描述**: 查看管理员操作日志（操作/时间/IP/内容）
**前端**: 需要创建 page.tsx
**位置**: `/workspace/apps/admin/src/app/(admin)/logs/page.tsx`
**后端**: 需要 audit log API 路由

#### ADM-005: 用户详情页面 (P1)
**描述**: 查看用户详细信息（账号/项目/AI 用量/订单/操作日志）
**前端**: 需要创建 `users/[id]/page.tsx`
**位置**: `/workspace/apps/admin/src/app/(admin)/users/[id]/page.tsx`
**后端**: 需要 user detail API

#### ADM-006: 企业管理页面 (P1)
**描述**: 管理企业空间（企业信息/成员/套餐/用量）
**前端**: 需要创建 spaces/ 目录
**后端**: 需要对应的 admin API 路由

#### ADM-007: AI 调用统计页面 (P1)
**描述**: 模型调用次数、Token 消耗、费用、失败次数
**前端**: 需要 stats/ 子页面或整合到 models/
**后端**: 需要 admin AI usage stats API

#### ADM-008: 内容审核管理页面 (P1)
**描述**: 审核用户输入，检测违规/垃圾/恶意内容
**位置**: `/workspace/apps/admin/src/app/(admin)/audit/page.tsx`

#### ADM-009: 套餐/订阅管理页面 (P0)
**描述**: 管理 SaaS 套餐配置，查看所有订阅
**位置**: 可集成到 orders/ 或新目录 plans/

---

## Phase 4: Website Homepage (Astro)

### 参考文档
- UI-UX设计规范.md (syntro-astro template, Vercel/Stripe style)
- 01-项目介绍.md, 02-产品定位与商业模式.md

### 页面结构 (按 UI-UX设计规范.md)

```
/
├── Hero — "AI 品牌可见度分析与增长平台"
│   - 核心标题: "了解 AI 是否认识你的品牌"
│   - 副标题: 价值主张
│   - CTA: "免费检测" / "查看演示"
│
├── 产品价值 — 三大核心价值
│   - AI 可见度分析 (监控品牌曝光)
│   - AI 竞争分析 (竞争对手为什么被推荐)
│   - AI 增长建议 (下一步优化什么)
│
├── 核心功能 — 功能展示卡片
│   - AI Visibility Score
│   - 品牌检测
│   - 竞争分析
│   - 报告中心
│   - 持续监控
│
├── 使用场景 — 行业案例
│   - SaaS 行业
│   - 教育行业
│   - 电商行业
│   - 医疗行业
│   - 个人品牌
│
├── 数据展示 — Demo / 截图 / 数据可视化
│   - Dashboard 预览
│   - 评分卡片
│   - 趋势图
│
├── 客户案例 — Testimonials / 用户证言
│   - 企业案例
│   - 成果数据
│
├── 价格方案 — Pricing
│   - Free (体验)
│   - Pro (中小企业)
│   - Business (营销团队)
│   - Enterprise (定制)
│   - 对比表
│
├── CTA — 最终行动号召
│   - "开始免费使用"
│   - "预约演示"
│
├── Footer
│   - 产品链接
│   - 定价
│   - 博客
│   - 联系方式
```

### 技术要求
- 框架: Astro
- 样式: Tailwind CSS
- 模板: syntro-astro (GitHub: bekturaslan/syntro-astro)
- 设计风格: Vercel / Stripe 风格
- 禁止: Emoji, 大量动画, 花哨渐变

### 额外页面 (Website)
| 页面 | 路径 | 优先级 |
|------|------|--------|
| 功能展示 | /features | P0 |
| 定价 | /pricing | P0 |
| 博客 | /blog | P1 |
| 关于 | /about | P1 |
| 联系我们 | /contact | P1 |

---

## Prioritized Execution Tasklist

### 优先级矩阵

```
Critical Path (必须串行):
  数据库迁移 → API 路由 → Backend Services → Frontend Pages

可并行:
  Website (Astro) 独立于 Dashboard/Admin 开发
  各个 Backend Microservice 可独立开发
  Frontend Pages 可并行开发 (不同 feature)
```

### 并行 Agent 任务分派

#### Agent Group A: Backend Foundation (依赖链优先)
```
A1: 创建缺失的数据库迁移文件 (000004~000013)
    依赖: 无
    输出: 10 个 SQL 迁移文件

A2: 补全 API 路由和 Handler
    依赖: A1
    输出: 扩展 router.go + 新增 handler 文件

A3: 实现 AI 分析引擎 (analyzer/ + scorer/ + prompts/)
    依赖: F-005 (已有 provider)
    输出: 品牌检测/曝光统计/排名分析/评分生成
```

#### Agent Group B: Business Services (可并行)
```
B1: 实现 report service (报告生成/导出/分享)
    依赖: A3
    输出: services/report/ 完整实现

B2: 实现 billing service (套餐/订阅/支付/订单)
    依赖: A1 (数据库表)
    输出: services/billing/ 完整实现

B3: 实现 notification service (站内通知/邮件)
    依赖: 无
    输出: services/notification/ 完整实现

B4: 实现 scheduler service (定时检测/监控)
    依赖: A3
    输出: services/scheduler/ 完整实现

B5: 实现 worker service (异步任务队列)
    依赖: A2, A3
    输出: services/worker/ 完整实现
```

#### Agent Group C: Frontend Dashboard (可并行)
```
C1: Dashboard 页面完善 (验证+补全)
    依赖: A2
    文件: /workspace/apps/dashboard/src/app/(protected)/dashboard/

C2: Analytics 页面 (趋势图表 + Tremor)
    依赖: A3, A2
    文件: /workspace/apps/dashboard/src/app/(protected)/analytics/

C3: 竞争分析页面
    依赖: A3
    文件: /workspace/apps/dashboard/src/features/competitors/

C4: 报告查看页面完善
    依赖: B1
    文件: /workspace/apps/dashboard/src/app/(protected)/reports/

C5: 订阅管理页面
    依赖: B2
    文件: /workspace/apps/dashboard/src/features/billing/

C6: 品牌资料库 UI
    依赖: A2
    文件: /workspace/apps/dashboard/src/app/(protected)/projects/[id]/

C7: 通知中心 UI
    依赖: B3
    文件: /workspace/apps/dashboard/src/features/

C8: 空间/团队成员管理
    依赖: A2
    文件: /workspace/apps/dashboard/src/features/

C9: API Key 管理页面
    依赖: B2
    文件: /workspace/apps/dashboard/src/features/
```

#### Agent Group D: Frontend Admin (可并行)
```
D1: Admin Dashboard 完善
    依赖: A2
    文件: /workspace/apps/admin/src/app/(admin)/page.tsx

D2: Prompt 管理页面 (新增)
    依赖: A2
    文件: /workspace/apps/admin/src/app/(admin)/prompts/

D3: 任务管理页面 (新增)
    依赖: A2, B5
    文件: /workspace/apps/admin/src/app/(admin)/tasks/

D4: 操作日志页面 (新增)
    依赖: A2
    文件: /workspace/apps/admin/src/app/(admin)/logs/

D5: 用户详情页面 (新增)
    依赖: A2
    文件: /workspace/apps/admin/src/app/(admin)/users/[id]/

D6: AI 调用统计 (新增/完善)
    依赖: A2
    文件: /workspace/apps/admin/src/app/(admin)/models/stats/

D7: 套餐/订阅管理 (新增)
    依赖: B2
    文件: /workspace/apps/admin/src/app/(admin)/plans/
```

#### Agent Group E: Website (独立)
```
E1: 初始化 Astro 项目 (syntro-astro 模板)
    依赖: 无
    文件: /workspace/apps/website/

E2: Landpage Page (Hero → 产品价值 → 核心功能)
    依赖: E1
    文件: /workspace/apps/website/src/pages/index.astro

E3: 使用场景 + 数据展示 Sections
    依赖: E1
    文件: /workspace/apps/website/src/pages/index.astro

E4: Pricing + CTA + Footer
    依赖: E1
    文件: /workspace/apps/website/src/pages/index.astro

E5: /features, /pricing, /blog, /about, /contact pages
    依赖: E1
    文件: /workspace/apps/website/src/pages/
```

### Recommended Execution Order (Wave 1-4)

**Wave 1 — Foundation (all parallel):**
- A1: Database migrations
- B3: Notification service (no deps)
- B5: Worker service core
- E1: Astro init
- F-005: Provider interface complete

**Wave 2 — Backbone (after Wave 1):**
- A2: API routes and handlers (depends on A1)
- B1: Report service (depends on A2)
- B4: Scheduler service (depends on A2)
- A3: AI analysis engine (depends on F-005)
- E2-E3: Website hero + features sections

**Wave 3 — Features (after Wave 2):**
- B2: Billing service
- C1-C4: Dashboard core pages
- D1-D4: Admin management pages
- E4-E5: Website pricing + extra pages

**Wave 4 — Polish (after Wave 3):**
- C5-C9: Billing, notifications, spaces, API Key UI
- D5-D7: Admin detail pages
- All remaining page polish and error handling

---

## Data Models — Missing Tables Schema Plan

### 000004_plans
```sql
CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,           -- Free/Professional/Business/Enterprise
    code VARCHAR(50) NOT NULL UNIQUE,     -- free/pro/business/enterprise
    description TEXT,
    price_monthly DECIMAL(10,2),
    price_yearly DECIMAL(10,2),
    max_projects INTEGER DEFAULT 1,
    max_ai_queries INTEGER DEFAULT 10,
    max_reports INTEGER DEFAULT 5,
    max_team_members INTEGER DEFAULT 1,
    api_access BOOLEAN DEFAULT false,
    features JSONB,                       -- 功能特性列表
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000005_subscriptions
```sql
CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active/canceled/expired/trialing
    billing_cycle VARCHAR(10) NOT NULL DEFAULT 'monthly', -- monthly/yearly
    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    trial_end TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000006_orders
```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    order_no VARCHAR(100) NOT NULL UNIQUE,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'CNY',
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- pending/paid/canceled/refunding/refunded
    payment_method VARCHAR(50),
    payment_time TIMESTAMPTZ,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000007_competitors
```sql
CREATE TABLE competitors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES brand_projects(id),
    name VARCHAR(255) NOT NULL,
    website VARCHAR(500),
    mention_count INTEGER DEFAULT 0,
    rank_position INTEGER,
    advantages JSONB,
    analysis JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000008_reports
```sql
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES brand_projects(id),
    task_id UUID REFERENCES ai_tasks(id),
    user_id UUID NOT NULL REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    type VARCHAR(50) DEFAULT 'standard',  -- standard/full/industry
    visibility_score INTEGER,
    content JSONB NOT NULL,               -- 完整报告内容
    summary TEXT,
    status VARCHAR(20) DEFAULT 'draft',   -- draft/published/archived
    share_token VARCHAR(100),
    share_expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000009_notifications
```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,            -- task_complete/report_ready/score_change/system
    title VARCHAR(255) NOT NULL,
    content TEXT,
    is_read BOOLEAN DEFAULT false,
    related_id UUID,                     -- 关联对象 ID
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000010_prompts
```sql
CREATE TABLE prompts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    purpose VARCHAR(100) NOT NULL,        -- question_generation/analysis/score/recommendation
    version INTEGER DEFAULT 1,
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'draft',   -- draft/published/archived
    created_by UUID REFERENCES admins(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE prompt_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    prompt_id UUID NOT NULL REFERENCES prompts(id),
    version INTEGER NOT NULL,
    content TEXT NOT NULL,
    changelog TEXT,
    created_by UUID REFERENCES admins(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000011_spaces
```sql
CREATE TABLE spaces (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(500),
    industry VARCHAR(100),
    region VARCHAR(200),
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE space_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    space_id UUID NOT NULL REFERENCES spaces(id),
    user_id UUID NOT NULL REFERENCES users(id),
    role VARCHAR(20) NOT NULL DEFAULT 'member', -- owner/admin/member/viewer
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000012_api_keys
```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    key_prefix VARCHAR(10) NOT NULL,       -- 显示前缀 (aig_xxx)
    permissions JSONB,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

### 000013_brand_infos
```sql
CREATE TABLE brand_infos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES brand_projects(id),
    brand_intro TEXT,
    product_intro TEXT,
    service_intro TEXT,
    faq JSONB,
    advantages TEXT[],
    cases JSONB,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 总量估算

| 维度 | 数量 |
|------|------|
| 新数据库迁移 | 10 个 |
| 新 API 路由 | ~30-40 条 |
| 新后端 Handler | ~15-20 个 |
| 新后端 Service 逻辑 | ~50-80 个函数 |
| Dashboard 新/完善页面 | ~10-15 个 |
| Admin 新/完善页面 | ~8-12 个 |
| Website 页面 | ~6 个 |
| AI Engine 内部模块 | ~15-20 个 |
| 总估算人天 (单 Agent) | ~30-40 天 |
| 总估算人天 (5 Agent 并行) | ~10-14 天 |
