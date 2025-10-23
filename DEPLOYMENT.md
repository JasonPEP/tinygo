# TinyGo 部署到 Railway 指南

本项目已从 SQLite 迁移到 PostgreSQL，并配置了完整的 Railway 部署方案。

## 项目变更

### 数据库迁移
- ✅ 从 SQLite 迁移到 PostgreSQL
- ✅ 添加了 PostgreSQL 驱动支持
- ✅ 配置了环境变量支持（Railway 自动提供 `DATABASE_URL`）
- ✅ 保持了 SQLite 作为本地开发选项

### 部署配置
- ✅ 创建了 `railway.json` 配置文件
- ✅ 创建了 `nixpacks.toml` 构建配置
- ✅ 创建了 `Dockerfile` 作为备选方案
- ✅ 配置了 GitHub Actions 自动部署

## Railway 部署步骤

### 1. 准备 Railway 账户
1. 访问 [Railway.app](https://railway.app)
2. 使用 GitHub 账户登录
3. 连接你的 GitHub 仓库

### 2. 创建新项目
1. 在 Railway 控制台点击 "New Project"
2. 选择 "Deploy from GitHub repo"
3. 选择你的 `tinygo` 仓库
4. 选择 "Deploy Now"

### 3. 添加 PostgreSQL 数据库
1. 在项目页面点击 "New"
2. 选择 "Database" → "PostgreSQL"
3. Railway 会自动创建数据库并设置 `DATABASE_URL` 环境变量

### 4. 配置环境变量（可选）
Railway 会自动设置以下环境变量：
- `DATABASE_URL` - PostgreSQL 连接字符串
- `PORT` - 应用端口（Railway 自动设置）

你可以手动设置以下环境变量：
- `BASE_URL` - 你的应用域名（如：https://your-app.railway.app）
- `LOG_LEVEL` - 日志级别（默认：info）
- `LOG_FORMAT` - 日志格式（默认：text）

### 5. 部署
Railway 会自动：
1. 检测到 Go 项目
2. 使用 `nixpacks.toml` 配置构建
3. 运行 `go mod download` 下载依赖
4. 运行 `go build -o bin/tinygo ./cmd/server` 构建应用
5. 启动应用

## GitHub Actions 自动部署

### 设置 Secrets
在 GitHub 仓库设置中添加以下 Secrets：
1. `RAILWAY_TOKEN` - Railway API Token
2. `RAILWAY_SERVICE` - Railway 服务名称

### 获取 Railway Token
1. 访问 [Railway Account Settings](https://railway.app/account/tokens)
2. 点击 "Create Token"
3. 复制生成的 token
4. 在 GitHub 仓库设置中添加为 Secret

### 自动部署流程
- 推送到 `main` 分支时自动部署
- 包含测试、构建和部署步骤
- 使用 Railway CLI 进行部署

## 本地开发

### 使用 PostgreSQL（推荐）
```bash
# 安装 PostgreSQL
brew install postgresql  # macOS
# 或使用 Docker
docker run --name postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

# 创建数据库
createdb tinygo

# 设置环境变量
export DATABASE_DRIVER=postgres
export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=tinygo port=5432 sslmode=disable"

# 运行应用
go run ./cmd/server
```

### 使用 SQLite（开发）
```bash
# 设置环境变量
export DATABASE_DRIVER=sqlite
export DATABASE_DSN="data/tinygo.db"

# 运行应用
go run ./cmd/server
```

## 验证部署

### 健康检查
访问 `https://your-app.railway.app/healthz` 应该返回 "ok"

### API 测试
```bash
# 创建短链接
curl -X POST https://your-app.railway.app/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://example.com"}'

# 访问短链接（应该重定向到原始URL）
curl -I https://your-app.railway.app/abc123
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查 `DATABASE_URL` 环境变量
   - 确认 PostgreSQL 服务正在运行

2. **构建失败**
   - 检查 Go 版本（需要 1.25+）
   - 确认所有依赖都已下载

3. **部署失败**
   - 检查 Railway 日志
   - 确认环境变量设置正确

### 查看日志
```bash
# 使用 Railway CLI
railway logs

# 或在 Railway 控制台查看
```

## 项目结构

```
tinygo/
├── cmd/server/main.go          # 应用入口
├── internal/
│   ├── config/                # 配置管理
│   ├── database/              # 数据库连接
│   ├── storage/               # 数据存储层
│   ├── shortener/             # 业务逻辑
│   └── transport/http/        # HTTP 处理器
├── web/                       # 前端资源
├── railway.json               # Railway 配置
├── nixpacks.toml             # 构建配置
├── Dockerfile                # Docker 配置
└── .github/workflows/        # GitHub Actions
```

## 环境变量说明

| 变量名 | 描述 | 默认值 | Railway 自动设置 |
|--------|------|--------|------------------|
| `DATABASE_URL` | PostgreSQL 连接字符串 | - | ✅ |
| `DATABASE_DRIVER` | 数据库驱动 | postgres | - |
| `DATABASE_DSN` | 数据库连接字符串 | - | - |
| `BASE_URL` | 应用基础URL | http://localhost:8080 | - |
| `ADDR` | 监听地址 | :8080 | - |
| `PORT` | 端口号 | 8080 | ✅ |
| `LOG_LEVEL` | 日志级别 | info | - |
| `LOG_FORMAT` | 日志格式 | text | - |

## 下一步

1. 配置自定义域名（可选）
2. 设置 SSL 证书（Railway 自动提供）
3. 配置监控和告警
4. 设置备份策略
