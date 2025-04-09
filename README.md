# 票务预订系统后端

这是一个基于Go语言开发的现代化票务预订系统后端服务。系统采用微服务架构设计，提供完整的票务管理、用户认证、订单处理等功能。通过RESTful API接口，支持多平台接入，确保系统的高可用性和可扩展性。

## 🚀 功能特性

- **用户管理**
  - 用户注册与登录
  - 角色权限管理
  - JWT + Redis 混合认证授权
  - 用户会话管理
  - 主动登出功能
  - 个人信息管理

- **票务管理**
  - 票务信息CRUD
  - 票务库存管理
  - 票务状态追踪
  - 动态二维码生成
  - 基于活动时间的二维码过期机制
  - 票务验证系统

- **活动管理**
  - 活动信息CRUD
  - 活动时间管理（开始时间、结束时间）
  - 活动状态追踪
  - 活动统计功能（已购票数、已入场数）

- **统计分析**
  - 票务销售统计
  - 活动参与度分析
  - 用户行为分析

- **系统特性**
  - Swagger API文档
  - Docker容器化部署
  - 数据库迁移支持
  - Redis缓存支持
  - 日志系统
  - 错误处理
  - 数据验证
  - 并发控制
  - 数据一致性保证

## 🛠️ 技术栈

<div align="center">
  
| 类别 | 技术 | 版本 | 说明 |
|------|------|------|------|
| **后端** | ![Go](https://img.shields.io/badge/Go-1.24-blue?style=flat&logo=go) | 1.24 | 高性能编程语言 |
| **Web框架** | ![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=flat&logo=go) | v2 | 高性能Web框架 |
| **数据库** | ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat&logo=postgresql) | 15+ | 关系型数据库 |
| **缓存** | ![Redis](https://img.shields.io/badge/Redis-7.0-DC382D?style=flat&logo=redis) | 7.0 | 内存数据库 |
| **ORM** | ![GORM](https://img.shields.io/badge/GORM-v1.25-00ADD8?style=flat&logo=go) | v1.25 | Go ORM库 |
| **认证** | ![JWT](https://img.shields.io/badge/JWT-v5-000000?style=flat&logo=jsonwebtokens) | v5 | 身份认证 |
| **配置** | ![godotenv](https://img.shields.io/badge/godotenv-v1.5-ECD53F?style=flat&logo=dotenv) | v1.5 | 环境变量管理 |
| **文档** | ![Swagger](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger) | 2.0 | API文档生成 |
| **容器** | ![Docker](https://img.shields.io/badge/Docker-24.0-2496ED?style=flat&logo=docker) | 24.0 | 容器化部署 |
| **开发工具** | ![Air](https://img.shields.io/badge/Air-1.52-00ADD8?style=flat&logo=go) | 1.52 | 热重载工具 |
| **二维码** | ![QRCode](https://img.shields.io/badge/QRCode-1.0-000000?style=flat&logo=qrcode) | 1.0 | 二维码生成 |

</div>

## 📋 系统要求

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7.0+
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

4. 更新Swagger文档：
```bash
sh update_swagger.sh
```

## 📊 项目结构

```
.
├── cmd/               # 应用程序入口点
│   └── api/           # API服务入口
├── config/            # 配置管理
├── db/                # 数据库连接和迁移
├── docs/              # Swagger文档
├── handlers/          # HTTP请求处理器
├── middlewares/       # 中间件
├── models/            # 数据模型
├── repositories/      # 数据访问层
├── services/          # 业务逻辑层
├── utils/             # 工具函数
├── .air.toml          # Air热重载配置
├── .env               # 环境变量
├── .env.example       # 环境变量示例
├── docker-compose.yaml # Docker Compose配置
├── Dockerfile         # Docker构建文件
├── go.mod             # Go模块定义
├── go.sum             # Go依赖锁定
├── Makefile           # 项目管理命令
├── README.md          # 项目文档
└── update_swagger.sh  # Swagger更新脚本
```

## 📝 更新日志

### v1.0.2 (2024-04-09)
- 改进票务验证机制
- 优化活动时间管理功能
- 完善用户会话控制
- 提升API响应性能
- 优化数据库索引

### v1.0.1 (2024-04-08)
- 新增动态二维码生成功能
- 实现基于活动时间的二维码过期机制
- 优化票务验证系统
- 完善活动时间管理
- 改进错误处理机制
- 优化代码结构和可维护性
- 增强数据一致性保证
- 添加并发控制机制

## 📄 许可证

MIT License

## 📝 TODO 列表

### 功能增强
- [✓] 实现用户登出功能
- [ ] 集成支付网关（支付宝、微信支付）
- [✓] 添加票务二维码生成与验证
- [ ] 实现票务转赠功能
- [ ] 添加票务收藏功能
- [ ] 实现票务搜索和筛选功能
- [ ] 添加票务推荐系统

### 性能优化
- [✓] 实现Redis缓存层
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
- [✓] 添加项目结构说明
- [ ] 添加详细的API使用示例
- [ ] 编写部署文档
- [ ] 添加性能测试报告
- [ ] 完善开发文档
- [ ] 添加故障处理指南 