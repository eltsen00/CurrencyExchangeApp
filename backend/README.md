# Backend 说明

后端基于 Go + Gin，提供货币兑换与新闻相关 API。

## 运行环境

- Go 1.25+
- MySQL 8+
- Redis 6+

## 配置文件

配置位于：`config/config.yml`

```yaml
app:
  name: CurrencyExchangeApp
  port: 3000

database:
  host: localhost
  port: 3306
  user: root
  password: 123456
  name: test
  MaxOpenConns: 100
  MaxIdleConns: 50
```

Redis 连接当前在代码中固定为：`localhost:6379`（见 `config/redis.go`）。

## 启动步骤

```bash
go mod tidy
go run main.go
```

服务默认监听：`http://localhost:3000`

## 接口总览

Base URL: `http://localhost:3000/api`

### 公开接口

- `GET /exchangeRates`：获取所有汇率
- `POST /auth/register`：用户注册
- `POST /auth/login`：用户登录

### 需认证接口

以下接口都需要在 Header 里携带：

- `Authorization: <JWT_TOKEN>`

接口列表：

- `POST /exchangeRates`：新增汇率
- `POST /articles`：发布文章
- `GET /articles`：获取文章列表
- `GET /articles/:id`：获取文章详情
- `POST /articles/:id/like`：点赞文章
- `GET /articles/:id/like`：获取文章点赞数

## 数据模型

### User

- `username`（唯一）
- `password`（bcrypt 哈希）

### ExchangeRate

- `fromCurrency`
- `toCurrency`
- `rate`
- `date`（服务端写入当前时间）

### Article

- `title`
- `preview`
- `content`

## 缓存与点赞设计

- 文章列表缓存 Key：`articles_cache`
- 文章详情缓存 Key：`article_id_cache{id}`
- 点赞计数 Key：`article:{id}:likes`

## 调试示例

### 注册

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"123456"}'
```

### 登录

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"123456"}'
```

### 新增汇率（带 Token）

```bash
curl -X POST http://localhost:3000/api/exchangeRates \
  -H "Content-Type: application/json" \
  -H "Authorization: <TOKEN>" \
  -d '{"fromCurrency":"USD","toCurrency":"CNY","rate":7.23}'
```

## 注意事项

- 首次写入时，部分表会在对应接口中触发 `AutoMigrate`
- `GET /articles` 在无数据时返回 `404`
- CORS 当前允许 `*`，生产环境建议收敛到明确域名
