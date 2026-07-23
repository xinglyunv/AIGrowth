# AI-Growth-Engine UI-UX设计规范

版本：v1.0

更新时间：2026-07

---

# 1. 文档目的

本文档用于规范 AI-Growth-Engine 项目的前端 UI 开发标准。

目标：

统一：

- 官网 Website
- 用户端 Console
- 管理端 Admin
- 数据分析 Dashboard

的设计语言、模板来源、组件规范。

所有开发人员和 AI Agent 在开发前必须阅读本文档。

---

# 2. 产品 UI 定位


## 产品类型

AI 企业级 SaaS 平台。


核心功能：

- AI品牌曝光分析
- AI搜索可见度监控
- 竞争品牌分析
- AI增长建议
- 数据报告生成


---

# 3. 整体设计方向


## 设计关键词


```
Enterprise SaaS

Minimal

Professional

Data Driven

Modern

Trustworthy
```


---

## 参考产品


设计参考：

- Linear
- Vercel Dashboard
- Stripe Dashboard
- OpenAI Dashboard
- Ahrefs
- Semrush


---

# 4. 前端系统划分


项目包含三个前端系统。


```
AI-Growth-Engine

├── Website
│   官网

├── Console
│   用户端 SaaS

└── Admin
    管理后台
```


---

# 5. Website 官网规范


## 5.1 模板来源


使用：

syntro-astro


GitHub:

https://github.com/bekturaslan/syntro-astro


在线：

https://aisaasui.vercel.app


---

## 5.2 技术要求


框架：

Astro


样式：

Tailwind CSS


---

## 5.3 使用范围


仅用于：

- 首页
- 产品介绍
- 功能展示
- 定价页面
- 博客
- 联系页面


---

## 5.4 页面结构


首页：


```
Hero

↓

产品价值

↓

核心功能

↓

使用场景

↓

数据展示

↓

客户案例

↓

价格方案

↓

CTA
```


---

## 5.5 官网设计要求


必须：

- 强调产品价值
- 突出商业转化
- 保持高级感


参考：

Vercel 官网

Stripe 官网


禁止：

- 做成后台页面
- 添加复杂业务组件
- 添加大量表格


---

# 6. 用户端 Console 规范


访问路径：

```
/dashboard
```


---

# 6.1 UI模板


主要使用：

Shadcn UI


官方：

https://ui.shadcn.com/


---

辅助：

Shadcn UI Blocks


地址：

https://ui.shadcn.com/blocks


---

# 6.2 使用场景


用于：

用户登录后的工作空间。


包含：


```
Dashboard

品牌项目

AI检测任务

AI分析结果

竞争分析

报告中心

订阅管理

账户设置
```


---

# 6.3 页面布局


统一使用：

SaaS Dashboard Layout


结构：


```
--------------------------------

Sidebar

Navigation


--------------------------------

Header


--------------------------------

Main Content


--------------------------------
```


---

# 6.4 Sidebar规范


宽度：

```
240px
```


包含：


```
Logo

Dashboard

Projects

Analytics

Reports

Settings
```


---

# 6.5 Console组件规范


优先使用：

Shadcn组件。


包括：


## Card

用于：

指标展示。


## Table

用于：

数据列表。


## Form

用于：

配置。


## Dialog

用于：

操作确认。


## Tabs

用于：

模块切换。


---

# 7. Dashboard 数据分析规范


## UI模板


使用：

Tremor


官网：

https://tremor.so/


---

# 7.1 使用范围


用于：


```
AI Visibility Score

品牌趋势

竞争分析

曝光统计

报告分析
```


---

# 7.2 图表规范


统一使用：


- Metric Card
- Line Chart
- Area Chart
- Bar Chart
- Donut Chart


---

禁止：

- 3D图表
- 炫酷动画
- 复杂视觉效果


---

# 8. 管理端 Admin 规范


访问路径：

```
/admin
```


---

# 8.1 UI模板


推荐：

Shadcn Admin


来源：

https://github.com/satnaing/shadcn-admin


---

备用：

Ant Design Pro


来源：

https://pro.ant.design/


---

# 8.2 使用范围


管理：


```
用户管理

权限管理

套餐管理

订单管理

AI模型管理

Prompt管理

任务管理

日志管理

系统设置
```


---

# 8.3 Admin设计要求


特点：

- 信息密度高
- 操作效率优先
- 数据清晰


参考：

企业后台系统。


---

# 9. 登录注册页面规范


UI来源：

Shadcn Authentication Blocks


用于：


```
/login

/register

/reset-password

/verify
```


---

# 10. 图标规范


统一：

Lucide Icons


官网：

https://lucide.dev/


---

禁止：

- Emoji图标
- 随机SVG
- 风格不一致图标


---

# 11. 字体规范


默认：

```
Inter

-apple-system

Segoe UI

Roboto
```


---

# 12. 色彩规范


## 主色


Primary：

```
#2563EB
```


用途：

- 主按钮
- 重点数据
- 品牌颜色


---

## Success


```
#16A34A
```


用途：

增长。


---

## Warning


```
#F59E0B
```


用途：

风险。


---

## Danger


```
#DC2626
```


用途：

错误。


---

# 13. 圆角规范


按钮：

```
8px
```


卡片：

```
12px
```


弹窗：

```
16px
```


---

# 14. AI Agent UI开发规则


AI生成页面时必须遵守：


```
Use existing UI system.

Website:
syntro-astro


Console:
Shadcn UI


Dashboard:
Tremor


Admin:
Shadcn Admin


Style:

Enterprise SaaS

Minimal

Professional

Data focused


Do not:

- create random components
- use emoji
- use excessive gradients
- change design language
```


---

# 15. 模板引用关系


|模块|模板|
|-|-|
|官网|syntro-astro|
|用户端|Shadcn UI|
|Dashboard|Tremor|
|管理端|Shadcn Admin|
|登录注册|Shadcn Auth|
|图标|Lucide Icons|


---

# 16. 开发禁止事项


禁止：

- 引入新的 UI 框架
- 创建新的设计体系
- 页面风格不统一
- 使用 Emoji
- 使用低质量组件


---

# 17. 最终目标


AI-Growth-Engine UI 应达到：

```
官网：
Vercel / Stripe 风格


用户端：
Linear / OpenAI Dashboard 风格


管理端：
企业级 SaaS 后台风格
```


最终打造：

一个专业、可信、可商业化销售的 AI SaaS 产品。