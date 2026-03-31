# synapse-go

Go 共享基础设施库。为所有 Go 产品提供统一的基础设施层（L1）和平台能力层（L2）。

## 安装

```bash
# 配置私有模块访问（如仓库为 private）
export GOPRIVATE=github.com/techfitmaster/*

go get github.com/techfitmaster/synapse-go@latest
```

## 包列表

### L1 — 基础设施

| 包 | 说明 |
|---|------|
| `config` | 配置结构体（MySQL/Redis/Auth/SMTP）+ `GetEnv()` 环境变量加载 |
| `db` | MySQL 连接池初始化（GORM） |
| `redis` | Redis 连接初始化 |
| `logger` | 结构化日志（zap） |
| `migrate` | 数据库迁移（golang-migrate） |
| `timeutil` | 时间工具函数 |
| `mailer` | 邮件发送（SMTP + NoopMailer） |

### L2 — 平台能力

| 包 | 说明 |
|---|------|
| `resp` | 统一 API 响应格式（Success/Error/SuccessPage + 错误码） |
| `bizerr` | 业务错误类型（BizError + HandleError + 便捷构造函数） |
| `ginutil` | Gin 工具（路由参数解析 + Context Helpers） |
| `middleware` | HTTP 中间件（CORS、RequestID、JWT、角色鉴权、Header Secret） |
| `ratelimit` | 多维限流（RPM/TPM/并发，Redis Lua 脚本） |

## 使用示例

```go
import (
    "github.com/techfitmaster/synapse-go/config"
    "github.com/techfitmaster/synapse-go/db"
    "github.com/techfitmaster/synapse-go/logger"
    "github.com/techfitmaster/synapse-go/resp"
    "github.com/techfitmaster/synapse-go/middleware"
)

func main() {
    cfg := config.Load()
    log := logger.New(cfg.Env)
    gormDB, _ := db.New(cfg.MySQL)

    r := gin.New()
    r.Use(middleware.RequestIDMiddleware())
    r.Use(middleware.CORSMiddleware("http://localhost:3000"))

    r.GET("/health", func(c *gin.Context) {
        resp.Success(c, gin.H{"status": "ok"})
    })
}
```

## API 响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": { ... },
  "trace_id": "abc-123"
}
```

错误码规范：`0` 成功 / `1xxx` 客户端错误 / `2xxx` 鉴权错误 / `5xxx` 服务端错误

## 开发

```bash
make test          # 运行单元测试
make test-verbose  # 详细输出
make lint          # golangci-lint
make coverage      # 生成覆盖率报告
```

## 版本管理

推送 tag 发布新版本：
```bash
git tag v0.3.0
git push origin v0.3.0
```

消费方更新：
```bash
go get github.com/techfitmaster/synapse-go@v0.3.0
```
