# Frontend 使用说明

这是一个基于原生 HTML/CSS/JavaScript 的单页前端，对接 `backend` 的 API。

## 功能

- 用户注册、登录、退出（JWT 存储在 `localStorage`）
- 查询汇率（公开接口）
- 新增汇率（需要登录）
- 汇率换算器：输入金额后按选中汇率实时换算
- 兑换方向一键反转：例如 `USD → CNY` 可切换为 `CNY → USD` 自动重算
- 丝滑反转模式：可勾选“反转时保留结果为新输入金额”，连按反转时更顺滑
- 发布文章（需要登录）
- 查看文章列表和详情
- 点赞文章并查看点赞数

## 启动方式

在 `frontend` 目录下启动静态文件服务，例如：

```bash
cd /home/eltsen/VSCode/CurrencyExchangeApp/frontend
python3 -m http.server 5173
```

浏览器打开：

- http://localhost:5173

## 后端地址

默认请求地址是：

- `http://localhost:3000/api`

如果后端端口不同，可在浏览器控制台执行：

```js
localStorage.setItem('apiBase', 'http://localhost:你的端口/api')
location.reload()
```
