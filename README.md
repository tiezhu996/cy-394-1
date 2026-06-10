# 运动数据统计 API 服务

```bash
cp .env.example .env
docker compose up --build
```

为健身应用提供运动记录 CRUD、统计分析、排行榜、目标进度和健康设备 Webhook 的纯后端 API 服务。

## 项目主要功能

- 用户运动记录创建、查询、更新、删除，支持日期范围过滤
- 预置运动类型与 MET 值，用于卡路里估算
- 日、周、月维度统计总时长、距离、卡路里、次数和类型占比
- 查询个人记录 PR
- 本周、本月排行榜，支持好友圈筛选参数
- 每周运动目标设定与完成进度查询
- Webhook 接收外部设备推送运动数据
- JWT 鉴权、统一错误响应、zap 日志

## 本地开发方式

```bash
cd backend
go mod tidy
go run ./cmd
```

API 文档地址：http://localhost:19409/swagger/index.html

## 技术栈

| 分类 | 技术 |
| --- | --- |
| 后端框架 | Go + Gin |
| ORM | GORM |
| 数据库 | PostgreSQL |
| 认证 | JWT |
| API 文档 | Swagger / swaggo |
| 日志 | zap |
| 部署 | Docker Compose |

## 项目目录结构

```text
backend/
├── cmd/
├── docs/
├── internal/
│   ├── router/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── middleware/
│   ├── constants/
│   ├── utils/
│   └── config/
├── Dockerfile
└── go.mod
database/
└── init.sql
```

## 主要 API 列表

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 健康检查 |
| POST | `/api/v1/records` | 创建运动记录 |
| GET | `/api/v1/records` | 按日期范围查询记录 |
| PUT | `/api/v1/records/:id` | 更新运动记录 |
| DELETE | `/api/v1/records/:id` | 删除运动记录 |
| GET | `/api/v1/sports/types` | 查询运动类型字典 |
| GET | `/api/v1/stats` | 查询日/周/月统计 |
| GET | `/api/v1/pr` | 查询个人最佳记录 |
| GET | `/api/v1/rankings` | 查询排行榜 |
| POST | `/api/v1/goals` | 设置每周运动目标 |
| GET | `/api/v1/goals/progress` | 查询目标进度 |
| POST | `/api/v1/webhooks/device` | 接收设备运动数据 |

## 环境变量说明

| 变量 | 说明 |
| --- | --- |
| COMPOSE_PROJECT_NAME | Compose 项目名 |
| POSTGRES_DB | 数据库名 |
| POSTGRES_USER | 数据库用户 |
| POSTGRES_PASSWORD | 数据库密码 |
| DATABASE_DSN | GORM PostgreSQL 连接串 |
| JWT_SECRET | JWT 签名密钥 |
| SERVER_PORT | 后端端口 |

## License

MIT
