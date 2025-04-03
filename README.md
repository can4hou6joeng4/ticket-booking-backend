# 票务预订系统后端

这是一个基于Go语言开发的现代化票务预订系统后端服务。系统采用微服务架构设计，提供完整的票务管理、用户认证、订单处理等功能。通过RESTful API接口，支持多平台接入，确保系统的高可用性和可扩展性。

## 🚀 功能特性

- **用户管理**
  - 用户注册与登录
  - 角色权限管理
  - JWT认证授权
  - 个人信息管理

- **票务管理**
  - 票务信息CRUD
  - 票务分类管理
  - 票务库存管理
  - 票务状态追踪

- **订单系统**
  - 订单创建与支付
  - 订单状态管理
  - 订单历史记录
  - 退款处理

- **系统特性**
  - Swagger API文档
  - Docker容器化部署
  - 数据库迁移支持
  - 日志系统
  - 错误处理
  - 数据验证

## 🛠️ 技术栈

<div align="center">
  
| 类别 | 技术 | 版本 | 说明 |
|------|------|------|------|
| **后端** | ![Go](https://img.shields.io/badge/Go-1.24-blue?style=flat&logo=go) | 1.24 | 高性能编程语言 |
| **Web框架** | ![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=flat&logo=go) | v2 | 高性能Web框架 |
| **数据库** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql) | 15+ | 关系型数据库 |
| **ORM** | ![GORM](https://img.shields.io/badge/GORM-v1.25-00ADD8?style=flat&logo=go) | v1.25 | Go ORM库 |
| **认证** | ![JWT](https://img.shields.io/badge/JWT-v5-000000?style=flat&logo=jsonwebtokens) | v5 | 身份认证 |
| **配置** | ![godotenv](https://img.shields.io/badge/godotenv-v1.5-ECD53F?style=flat&logo=dotenv) | v1.5 | 环境变量管理 |
| **文档** | ![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger) | 2.0 | API文档生成 |
| **容器** | ![Docker](https://img.shields.io/badge/Docker-24.0-2496ED?style=flat&logo=docker) | 24.0 | 容器化部署 |
| **开发工具** | ![Air](https://img.shields.io/badge/Air-1.49-00ADD8?style=flat&logo=go) | 1.49 | 热重载工具 |

</div>

## 📋 系统要求

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 15+
- 内存: 4GB+
- 磁盘空间: 10GB+

## 🔧 安装说明

1. 克隆项目：
```bash
git clone https://github.com/your-username/ticket-booking-project-v1.git
cd ticket-booking-project-v1/backend
```

2. 配置环境变量：
```bash
cp .env.example .env
# 编辑.env文件，设置必要的环境变量
```

3. 使用Docker Compose启动服务：
```bash
docker-compose up -d
```

4. 访问API文档：
```
http://localhost:8081/swagger/index.html
```

## 🚀 开发指南

1. 启动开发环境：
```bash
make dev
```

2. 运行测试：
```bash
make test
```

3. 构建项目：
```bash
make build
```

## 📁 项目结构

```
.
├── cmd/
│   └── main.go                 # 应用程序入口
├── config/
│   ├── config.go              # 配置结构体
│   └── env.go                 # 环境变量配置
├── db/
│   ├── migrations/            # 数据库迁移文件
│   └── seeder/               # 数据库种子数据
├── docs/
│   └── swagger/              # Swagger文档
├── handlers/
│   ├── auth.go               # 认证处理器
│   ├── ticket.go             # 票务处理器
│   └── order.go              # 订单处理器
├── middlewares/
│   ├── auth.go               # 认证中间件
│   ├── logger.go             # 日志中间件
│   └── error.go              # 错误处理中间件
├── models/
│   ├── user.go               # 用户模型
│   ├── ticket.go             # 票务模型
│   └── order.go              # 订单模型
├── repositories/
│   ├── user.go               # 用户数据访问
│   ├── ticket.go             # 票务数据访问
│   └── order.go              # 订单数据访问
├── services/
│   ├── auth.go               # 认证服务
│   ├── ticket.go             # 票务服务
│   └── order.go              # 订单服务
├── utils/
│   ├── jwt.go                # JWT工具
│   ├── validator.go          # 验证工具
│   └── logger.go             # 日志工具
├── .env                      # 环境变量
├── .air.toml                 # Air配置
├── docker-compose.yaml       # Docker Compose配置
├── Dockerfile                # Docker配置
├── go.mod                    # Go模块文件
└── go.sum                    # Go依赖校验
```

## 📝 API文档

项目使用Swagger自动生成API文档，启动服务后可以通过以下地址访问：
```
http://localhost:8081/swagger/index.html
```

## 🔒 环境变量

主要环境变量配置：
- `DB_HOST`: 数据库主机
- `DB_NAME`: 数据库名称
- `DB_USER`: 数据库用户
- `DB_PASSWORD`: 数据库密码
- `JWT_SECRET`: JWT密钥
- `PORT`: 服务端口
- `LOG_LEVEL`: 日志级别
- `ENVIRONMENT`: 运行环境

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件 

## 📝 TODO 列表

### 功能增强
- [ ] 集成支付网关（支付宝、微信支付）
- [ ] 添加票务二维码生成与验证
- [ ] 实现票务转赠功能
- [ ] 添加票务收藏功能
- [ ] 实现票务搜索和筛选功能
- [ ] 添加票务推荐系统

### 性能优化
- [ ] 实现Redis缓存层
- [ ] 添加数据库读写分离
- [ ] 优化数据库查询性能
- [ ] 实现API限流功能
- [ ] 添加请求压缩

### 安全增强
- [ ] 实现双因素认证
- [ ] 添加IP白名单功能
- [ ] 实现敏感数据加密
- [ ] 添加操作日志审计
- [ ] 实现API签名验证

### 监控与运维
- [ ] 集成Prometheus监控
- [ ] 添加Grafana仪表盘
- [ ] 实现日志聚合
- [ ] 添加健康检查接口
- [ ] 实现自动备份功能

### 文档完善
- [ ] 添加详细的API使用示例
- [ ] 编写部署文档
- [ ] 添加性能测试报告
- [ ] 完善开发文档
- [ ] 添加故障处理指南 