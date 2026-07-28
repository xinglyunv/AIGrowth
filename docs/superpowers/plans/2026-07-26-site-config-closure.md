# 官网配置闭环优化实施计划

> **For agentic workers:** 当前会话内按任务顺序执行并在每个任务后验证。

**Goal:** 让官网公共页面、六套主题和管理端设置使用统一、可校验、可恢复的网站配置。

**Architecture:** 后端继续使用 `site_settings` key-value 表，将结构化字段编码为 JSON 并通过 `SiteConfig` 返回。官网增加共享配置读取与默认合并，管理端使用同一份嵌套类型编辑配置，保存前校验并在失败时恢复可操作状态。

**Tech Stack:** Go、PostgreSQL、Astro、Next.js 14、React、TypeScript、Tailwind CSS。

## Global Constraints

- 保留六套主题：`syntro`、`saas-landing`、`saas-ui`、`hikari`、`cruip`、`nextly`。
- 官网文案、Logo、标题、导航、按钮、统计、链接、SEO 和主题视觉参数支持后台配置。
- 官网不使用 emoji 或 Unicode 图标。
- 不新增依赖，不修改核心架构。
- API 失败时官网使用默认配置，管理端结束 loading 并显示中文错误。

### Task 1: 完善配置边界

**Files:**
- Modify: `services/setting/model.go`
- Modify: `services/setting/repository.go`
- Test: `services/setting/model_test.go`

- [x] 为 `SiteConfig` 补齐 Footer、SEO 和仪表盘示例类型。
- [x] 为 JSON 字段增加默认值、空数组回退和主题白名单测试。
- [x] 运行 `go test ./...`（工作目录 `services/setting`）。

### Task 2: 建立官网共享配置读取器

**Files:**
- Create: `apps/website/src/lib/site-config.ts`
- Modify: `apps/website/src/layouts/Layout.astro`
- Modify: `apps/website/src/components/ThemeHome.astro`
- Test: `apps/website` production build

- [x] 将默认配置和 API 合并逻辑集中到 `site-config.ts`。
- [x] 让公共导航、CTA、Logo、页脚和 SEO 使用读取器返回的数据。
- [x] 让六套主题读取同一个配置对象，仅保留布局和样式差异。
- [x] 使用 CSS 变量注入可配置颜色并保留主题默认值。
- [x] 运行 `pnpm --filter @aige/website build`。

### Task 3: 统一官网内容页面

**Files:**
- Modify: `apps/website/src/pages/features.astro`
- Modify: `apps/website/src/pages/about.astro`
- Modify: `apps/website/src/pages/pricing.astro`
- Modify: `apps/website/src/pages/contact.astro`

- [x] 替换页面中的品牌标题、统计、CTA、联系方式和页脚重复默认值。
- [x] 保留页面自身的产品说明结构，所有可变品牌内容从共享配置读取。
- [x] 扫描 `apps/website/src`，移除剩余 Unicode 箭头、播放符号和 emoji。
- [x] 运行官网构建并检查生成页面数量。

### Task 4: 重整管理端设置页

**Files:**
- Modify: `apps/admin/src/app/(admin)/settings/page.tsx`
- Modify: `apps/admin/src/lib/api.ts`

- [x] 合并旧首页字段和结构化首页字段，保留一套编辑入口。
- [x] 按品牌、导航、首页、主题、联系和系统分组展示。
- [x] 增加必填项、颜色值、URL 和统计项校验。
- [x] 保存前提交完整嵌套 `SiteConfig`，展示保存中、成功和失败状态。
- [x] API 超时、网络错误和非 JSON 响应都转为中文错误。
- [x] 运行 `pnpm --filter @aige/admin build`。

### Task 5: 端到端验证

**Files:**
- Test: `services/api`
- Test: `services/setting`
- Test: `apps/website`
- Test: `apps/admin`

- [x] 运行 `go test ./...`（工作目录 `services/api`）。
- [x] 运行 `go test ./...`（工作目录 `services/setting`）。
- [x] 运行官网和管理端生产构建。
- [x] 使用正确数据库和开发 JWT 启动 API，验证 `/api/v1/settings` 返回完整配置。
- [x] 验证本地端口 `3000`、`3001`、`3002` 返回 `200`。
- [x] 运行 `git diff --check`。
