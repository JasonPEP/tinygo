# 🚀 Railway 部署指南

本指南将帮助你将 TinyGo 短链接服务部署到 Railway 平台。

## 📋 部署前准备

### 1. 安装 Railway CLI
```bash
# macOS
brew install railway

# 或者使用 npm
npm install -g @railway/cli
```

### 2. 登录 Railway
```bash
railway login
```

## 🚀 部署步骤

### 1. 初始化 Railway 项目
```bash
# 在项目根目录执行
railway init
```

### 2. 设置环境变量
在 Railway Dashboard 中设置以下环境变量：

**必需的环境变量：**
```
TINYGO_AUTH_USERNAME=admin
TINYGO_AUTH_PASSWORD=your_secure_password_here
TINYGO_BASE_URL=https://your-app-name.railway.app
```

**可选的环境变量：**
```
TINYGO_ADDR=:8080
TINYGO_DATABASE_DRIVER=sqlite
TINYGO_DATABASE_DSN=data/tinygo.db
TINYGO_LOG_LEVEL=info
TINYGO_LOG_FORMAT=json
TINYGO_AUTH_SESSION_KEY=your_custom_session_key
TINYGO_AUTH_SESSION_MAX_AGE=3600
```

### 3. 部署应用
```bash
# 部署到 Railway
railway up

# 或者使用 git 推送
git add .
git commit -m "Deploy to Railway"
git push origin main
```

## 🔧 Railway 配置说明

### 环境变量设置
1. 登录 [Railway Dashboard](https://railway.app/dashboard)
2. 选择你的项目
3. 进入 "Variables" 标签页
4. 添加所需的环境变量

### 自动部署
Railway 支持 Git 自动部署：
- 推送到 `main` 分支会自动触发部署
- 每次推送都会重新构建和部署应用

### 健康检查
Railway 会自动检查 `/healthz` 端点来确保应用正常运行。

## 🔐 安全建议

### 1. 强密码
使用强密码作为 `TINYGO_AUTH_PASSWORD`：
```bash
# 生成强密码
openssl rand -base64 32
```

### 2. 自定义会话密钥
设置自定义的会话密钥：
```bash
# 生成随机会话密钥
openssl rand -base64 32
```

### 3. HTTPS
Railway 自动提供 HTTPS 支持，确保 `TINYGO_BASE_URL` 使用 `https://` 协议。

## 📊 监控和日志

### 查看日志
```bash
# 使用 Railway CLI 查看日志
railway logs

# 或者查看实时日志
railway logs --follow
```

### 监控指标
在 Railway Dashboard 中可以查看：
- CPU 使用率
- 内存使用率
- 网络流量
- 请求数量

## 🛠️ 故障排除

### 常见问题

1. **应用启动失败**
   - 检查环境变量是否正确设置
   - 查看 Railway 日志：`railway logs`

2. **认证失败**
   - 确认 `TINYGO_AUTH_USERNAME` 和 `TINYGO_AUTH_PASSWORD` 已设置
   - 检查密码是否包含特殊字符

3. **数据库问题**
   - Railway 使用临时文件系统，重启后数据会丢失
   - 考虑使用 Railway 的 PostgreSQL 插件进行持久化存储

### 获取帮助
- [Railway 文档](https://docs.railway.app/)
- [Railway Discord](https://discord.gg/railway)
- 项目 Issues: 在 GitHub 仓库中创建 Issue

## 🔄 更新部署

### 代码更新
```bash
# 提交更改
git add .
git commit -m "Update application"
git push origin main

# Railway 会自动部署更新
```

### 环境变量更新
在 Railway Dashboard 中更新环境变量后，应用会自动重启。

## 📈 扩展和优化

### 数据库升级
考虑使用 Railway 的 PostgreSQL 插件：
1. 在 Railway Dashboard 中添加 PostgreSQL 插件
2. 更新环境变量：
   ```
   TINYGO_DATABASE_DRIVER=postgres
   TINYGO_DATABASE_DSN=${{Postgres.DATABASE_URL}}
   ```

### 性能优化
- 启用 GORM 连接池
- 配置适当的日志级别
- 监控内存使用情况