# TinyGo 短链接服务

一个简洁、高效的短链接服务，使用 Go 语言开发，支持 Web UI 和 REST API。

## ✨ 特性

- 🔗 **短链接生成** - 支持自定义短码或自动生成
- 📊 **统计分析** - 点击次数统计和访问记录
- 🌐 **Web UI** - 现代化的管理界面
- 🔌 **REST API** - 完整的 API 接口
- 💾 **数据持久化** - SQLite 数据库存储
- ⚙️ **配置管理** - YAML 配置文件和环境变量支持
- 📝 **结构化日志** - 使用 logrus 进行日志记录

## 🚀 快速开始

### 1. 克隆项目
```bash
git clone <repository-url>
cd tinygo
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 运行服务
```bash
go run ./cmd/server
```

或者编译后运行：
```bash
go build -o bin/urlshort ./cmd/server
./bin/urlshort
```

### 4. 访问服务
- **Web UI**: http://localhost:8080/
- **API 文档**: 见下方 API 接口说明

## 📁 项目结构

```
tinygo/
├── cmd/server/           # 应用程序入口
├── internal/             # 内部包
│   ├── config/          # 配置管理
│   ├── database/        # 数据库连接
│   ├── logger/          # 日志系统
│   ├── shortener/       # 核心业务逻辑
│   ├── storage/         # 数据存储层
│   └── transport/http/  # HTTP 传输层
├── web/                 # Web 资源
│   ├── static/         # 静态文件
│   └── templates/      # HTML 模板
├── configs/             # 配置文件
├── data/                # 数据文件（自动生成）
└── pkg/                 # 公共包
```

## ⚙️ 配置说明

配置文件位于 `configs/config.yaml`：

```yaml
# 服务器配置
addr: ":8080"
base_url: "http://localhost:8080"

# 数据库配置
database:
  driver: "sqlite"
  dsn: "data/tinygo.db"

# 日志配置
log_level: "info"
log_format: "text"
```

### 环境变量支持

所有配置都支持环境变量覆盖，使用 `TINYGO_` 前缀：

```bash
export TINYGO_ADDR=":9090"
export TINYGO_DATABASE_DSN="data/prod.db"
export TINYGO_LOG_LEVEL="debug"
```

## 🔌 API 接口

### 创建短链接
```bash
POST /api/shorten
Content-Type: application/json

{
  "long_url": "https://example.com",
  "custom_code": "mycode"  # 可选
}
```

### 获取链接信息
```bash
GET /api/links/{code}
```

### 删除链接
```bash
DELETE /api/links/{code}
```

### 获取统计信息
```bash
GET /admin/stats
```

### 短链接重定向
```bash
GET /{code}
```

## 🛠️ 开发说明

### 数据库自动创建
- 首次运行时，程序会自动创建 `data/` 目录
- SQLite 数据库文件会自动生成
- 数据库表结构会自动迁移

### 日志级别
- `debug`: 详细调试信息
- `info`: 一般信息（默认）
- `warn`: 警告信息
- `error`: 错误信息

### 日志格式
- `text`: 人类可读格式（默认）
- `json`: JSON 格式，适合日志收集系统

## 📦 依赖库

- **gorilla/mux**: HTTP 路由
- **gorm**: ORM 数据库操作
- **logrus**: 结构化日志
- **viper**: 配置管理
- **sqlite**: 数据库驱动

## 🎯 学习目标

这个项目展示了以下 Go 语言特性：

1. **项目结构** - 遵循 golang-standards/project-layout
2. **依赖注入** - 手动依赖注入模式
3. **接口设计** - 清晰的接口抽象
4. **错误处理** - Go 风格的错误处理
5. **并发安全** - 数据库操作的并发安全
6. **配置管理** - 多环境配置支持
7. **日志系统** - 结构化日志记录
8. **Web 开发** - HTTP 服务和静态文件服务

## 📄 许可证

MIT License