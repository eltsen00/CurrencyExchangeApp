# CurrencyExchangeApp

一个基于 Go + 原生前端的货币兑换项目，包含两大能力：

- 汇率管理：查询汇率、新增汇率
- 新闻功能：发布文章、查看详情、文章点赞
- 实时换算：输入金额按选中汇率实时计算，并支持兑换方向一键反转
- 丝滑反转模式：反转时可自动把上次换算结果作为新输入金额，便于连续反转

## 项目结构

```text
CurrencyExchangeApp/
├── backend/    # Go + Gin + Gorm + MySQL + Redis
└── frontend/   # 原生 HTML/CSS/JavaScript 单页应用
```

## 技术栈

- 后端：Go、Gin、Gorm、Viper、JWT、Redis
- 前端：HTML、CSS、JavaScript（无构建工具）
- 存储：MySQL（业务数据）+ Redis（缓存与点赞）

## 快速启动

### 1) 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

后端默认地址：`http://localhost:3000`

> 启动前请确保 MySQL 与 Redis 已运行，且 `backend/config/config.yml` 中数据库配置正确。

### 2) 启动前端

```bash
cd frontend
python3 -m http.server 5173
```

浏览器访问：`http://localhost:5173`

前端默认请求地址：`http://localhost:3000/api`

## 为什么前端可以用 Python 启动？

`frontend` 是静态资源（`index.html`、`styles.css`、`app.js`），不需要 Node 打包。

`python3 -m http.server` 只是起一个本地静态文件服务器，让浏览器通过 HTTP 访问这些文件；任何能提供静态文件服务的工具都可以替代它（如 Nginx、Caddy、VS Code Live Server）。

## 核心功能说明

### 认证

- 注册：`POST /api/auth/register`
- 登录：`POST /api/auth/login`
- 认证方式：JWT，放在请求头 `Authorization`

### 汇率

- 查询汇率（公开）：`GET /api/exchangeRates`
- 新增汇率（需登录）：`POST /api/exchangeRates`

### 新闻（文章）

- 发布文章（需登录）：`POST /api/articles`
- 查询列表（需登录）：`GET /api/articles`
- 查询详情（需登录）：`GET /api/articles/:id`
- 点赞（需登录）：`POST /api/articles/:id/like`
- 点赞数（需登录）：`GET /api/articles/:id/like`

## 常见问题

### 1) 前端请求失败

- 检查后端是否启动在 `3000` 端口
- 检查前端 `localStorage.apiBase` 是否被改错
- 检查请求头是否带了 `Authorization`（需要登录的接口）

### 2) 文章列表返回 404

后端在没有文章时会返回 `404 No articles found`，这是当前后端逻辑的设计。

## 后续建议

- 为汇率增加按币种筛选与排序
- 为新闻增加分页
- 生产环境改为 `config.yml` + 环境变量组合配置
