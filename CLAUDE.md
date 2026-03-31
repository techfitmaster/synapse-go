# CLAUDE.md

## 项目概述
synapse-go — Go 共享基础设施库，为所有 Go 产品提供统一的基础设施和平台能力。

## 技术栈
- Go 1.24
- Gin（HTTP 框架）
- GORM + MySQL（数据库）
- go-redis/v8（缓存）
- zap（日志）
- golang-migrate（迁移）
- golang-jwt/v5（JWT 鉴权）

## 架构
纯 Go library（无 main 入口），按功能分包：

```
config/       → 配置结构体 + 环境变量加载
db/           → MySQL 连接池（GORM）
redis/        → Redis 连接
logger/       → 结构化日志（zap）
migrate/      → 数据库迁移
resp/         → 统一 API 响应格式 + 错误码
bizerr/       → 业务错误类型 + Handler 映射
ginutil/      → Gin 工具（参数解析 + Context Helpers）
middleware/   → HTTP 中间件（CORS、RequestID、JWT、角色鉴权、Header Secret）
ratelimit/    → 多维限流（Redis Lua）
mailer/       → 邮件发送
timeutil/     → 时间工具
```

## 常用命令
```bash
make test          # go test ./... -race -cover（跳过集成测试）
make test-all      # 包含集成测试（需要 MySQL + Redis）
make lint          # golangci-lint run ./...
make coverage      # 生成覆盖率 HTML 报告
```

## 消费方
- 818-cargo（`go get github.com/techfitmaster/synapse-go`）
- 未来所有 Go 产品

## 规范
- 每个包独立、无循环依赖
- 所有导出函数必须有测试
- 版本通过 git tag 管理（语义版本）
- 合并 main 需 Albert 审批
