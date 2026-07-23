# 1. 项目介绍

AI-Growth-Engine 是一个 AI 品牌可见度分析与增长优化 SaaS 平台。

核心目标：

帮助企业、机构和个人分析：

- 品牌是否被 AI 认识
- AI 是否推荐该品牌
- 竞争品牌为什么获得曝光
- 如何优化品牌信息，提高 AI 理解能力

产品定位：

AI时代的 SEO / Analytics 平台。

类似：

- Ahrefs（搜索分析）
- Brandwatch（品牌监控）
- Google Search Console（搜索可见性）
- AI Agent（智能分析）


---

# 2. 开发原则

## 产品优先

所有开发必须围绕用户价值。

禁止：

- 为了技术炫技增加功能
- 开发无法产生业务价值的模块


## 简洁优先

代码：

简单 > 复杂

架构：

清晰 > 过度设计


不要提前引入：

- 复杂微服务
- 不必要中间件
- 过度抽象


## AI友好

代码必须：

- 结构清晰
- 命名明确
- 文档完整

方便 AI Agent 理解和继续开发。


---

# 3. 技术栈规范

## 前端

默认：

- Next.js
- React
- TypeScript
- TailwindCSS
- Shadcn UI


UI风格：

参考：

- Linear
- Vercel
- OpenAI Dashboard


禁止：

- 大量动画
- Emoji UI
- 花哨渐变


## 后端

推荐：

- Go
- PostgreSQL
- Redis


API：

RESTful API


## 数据处理

使用：

- PostgreSQL
- Elasticsearch
- Redis


异步任务：

- Queue Worker


---

# 4. 项目架构规范

推荐结构：

```
AI-Growth-Engine

apps/

├── web
│   用户端 SaaS

├── admin
│   管理后台


services/

├── api
│   API服务

├── ai-engine
│   AI分析服务

├── crawler
│   数据采集服务

├── report
│   报告生成服务


packages/

├── ui
├── types
├── utils
```


---

# 5. AI Agent开发规范

Agent必须明确职责。

禁止：

一个 Agent 完成所有事情。


推荐：

```
Question Agent

负责：
生成行业问题


Research Agent

负责：
收集品牌信息


Analysis Agent

负责：
分析AI回答


Report Agent

负责：
生成报告
```


---

# 6. 数据规范

所有核心数据必须保存：

- 创建时间
- 更新时间
- 数据来源
- 操作用户


禁止：

无来源数据。


---

# 7. AI调用规范

所有模型调用必须经过统一接口。

禁止：

业务代码直接调用模型。


正确：

```
Business Service

↓

AI Gateway

↓

LLM Provider
```


支持：

- OpenAI
- Claude
- Gemini
- 国内模型


---

# 8. API规范

API必须使用版本号。

例如：

```
/api/v1/projects
```


返回格式统一：

```json
{
  "success": true,
  "data": {},
  "message": ""
}
```


---

# 9. 数据库规范

表名：

snake_case


例如：

```
brand_projects

ai_tasks

ai_answers

competitor_reports
```


字段：

```
id

created_at

updated_at
```


---

# 10. 前端开发规范

组件必须拆分。

禁止：

一个页面超过500行。


组件命名：

PascalCase


例如：

```
BrandScoreCard

AIReportTable
```


---

# 11. UI规范

所有页面遵循：

Enterprise SaaS Design。


必须：

- 清晰层级
- 数据优先
- 响应式


禁止：

- Emoji
- 游戏化设计
- 过度装饰


---

# 12. Git规范

提交格式：

```
feat:
新增功能


fix:
修复问题


refactor:
重构


docs:
文档
```


例如：

```
feat: add AI visibility dashboard
```


---

# 13. 开发流程

任何功能开发：

必须经过：

1. 阅读需求文档

2. 确认数据库影响

3. 设计 API

4. 实现后端

5. 实现前端

6. 测试

7. 更新文档


---

# 14. 禁止事项

禁止：

- 删除已有功能
- 修改核心架构未经确认
- 写重复代码
- 随意增加依赖
- 泄露用户数据


---

# 15. AI Agent工作方式

开始任务前：

必须：

1. 阅读 README

2. 阅读 docs

3. 理解现有架构

4. 检查已有代码


不要：

直接开始写代码。


---

# 16. 输出要求

AI Agent 输出必须包含：

## 修改内容

修改了什么。


## 文件变化

哪些文件新增/修改。


## 测试结果

是否通过。


## 后续建议

下一步工作。


---

# 最终原则

写代码之前先理解产品。

简单可靠优先。

保持架构清晰。

让未来的 AI Agent 可以继续开发。

# UI规范引用

前端开发必须遵循以下 UI 规范：

## 官网

参考：

docs/UI-UX设计规范.md

技术模板：

syntro-astro

用途：

- Landing Page
- 产品介绍
- Pricing
- Blog


## 用户端 SaaS

参考：

docs/UI-UX设计规范.md

设计方向：

- Enterprise SaaS
- Data Dashboard
- Linear 风格
- OpenAI Dashboard 风格


## 管理端

参考：

docs/UI-UX设计规范.md

设计方向：

- 企业后台
- 数据管理
- 高密度信息展示


## UI开发规则

开发任何页面前：

必须先阅读：

1. docs/UI-UX设计规范.md
2. 当前页面需求文档
3. 已有组件规范


禁止：

- 自定义新的设计风格
- 引入其他 UI 体系
- 随意修改颜色、字体、布局规范