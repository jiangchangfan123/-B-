# Nebula.TV — 仿 B 站视频流媒体平台

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.25+-00ADD8?logo=go&logoColor=white" />
  <img src="https://img.shields.io/badge/Gin-1.12-008ECF?logo=gin&logoColor=white" />
  <img src="https://img.shields.io/badge/Vue-3.5-4FC08D?logo=vue.js&logoColor=white" />
  <img src="https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white" />
  <img src="https://img.shields.io/badge/Redis-7.0-DC382D?logo=redis&logoColor=white" />
  <img src="https://img.shields.io/badge/RabbitMQ-3.12-FF6600?logo=rabbitmq&logoColor=white" />
  <img src="https://img.shields.io/badge/MinIO-8.0-C72E49?logo=minio&logoColor=white" />
</p>

<p align="center">
  <b>基于 Go + Vue3 构建的视频分享与流媒体播放平台</b><br/>
  支持视频上传/转码、实时弹幕、楼中楼评论、点赞收藏、消息通知等核心功能
</p>

---


## 📑 目录

- [功能特性](#-功能特性)
- [技术架构](#-技术架构)
- [项目结构](#-项目结构)
- [快速开始](#-快速开始)
  - [环境依赖](#环境依赖)
  - [后端启动](#后端启动)
  - [前端启动](#前端启动)
  - [Docker 部署](#docker-部署)
- [配置说明](#-配置说明)
- [API 概览](#-api-概览)
- [核心实现亮点](#-核心实现亮点)
- [待优化项](#-待优化项)
- [技术栈版本](#-技术栈版本)

---

## ✨ 功能特性

### 🎬 视频系统
- [x] 视频上传：支持 mp4 / mov / webm / mkv 格式，单文件最大 2GB
- [x] FFmpeg 异步转码：上传后自动入队 RabbitMQ，消费者异步转码为 480p H.264
- [x] 智能封面：无封面时自动截取视频第 5 秒作为封面
- [x] MinIO 对象存储：视频原片、转码片、封面、头像统一对象存储管理
- [x] 分类筛选 & 排序：按分区筛选，支持最新/最热排序
- [x] 播放量统计 & 播放历史记录

### 💬 互动系统
- [x] **实时弹幕**：基于 WebSocket 分房间广播，支持滚动/顶部/底部三种类型，按时间轴精准加载
- [x] **楼中楼评论**：`parent_id` + `root_id` 双字段模型，支持一级评论分页 + 二级回复聚合
- [x] **点赞系统**：Redis 缓存用户点赞状态与计数，定时回写 MySQL，联合唯一索引保证幂等
- [x] **收藏系统**：一键收藏/取消，收藏列表 JOIN 视频表展示
- [x] **消息通知**：评论被回复、评论被点赞、视频被点赞三类通知，异步 goroutine 发送不阻塞主流程

### 👤 用户系统
- [x] JWT 认证：Bearer Token，Access Token 有效期 2 小时
- [x] 密码安全：bcrypt 哈希存储
- [x] 个人中心：资料修改、密码修改、头像上传（支持 jpg/png/webp）
- [x] 我的投稿 / 播放历史 / 消息中心

---

## 🏗 技术架构

```
┌─────────────────────────────────────────────────────────────┐
│                        前端层                                │
│   Vue 3 + Vite + TypeScript + Pinia + Vue Router + DPlayer  │
└──────────────────────────┬──────────────────────────────────┘
                           │ HTTP / WebSocket
┌──────────────────────────▼──────────────────────────────────┐
│                        接入层                                │
│              Nginx (反向代理 + 静态资源)                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                        网关层                                │
│   Gin + CORS + JWT Auth(Optional Auth) + Logger + Recovery   │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                        业务层                                │
│  Auth / User / Video / Comment / Danmaku / Like / Favorite  │
│  Notification / Upload / History                            │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                        数据层                                │
│   MySQL (主数据)   Redis (缓存)   MinIO (对象存储)            │
│   RabbitMQ (消息队列)                                        │
└─────────────────────────────────────────────────────────────┘
```

### 请求生命周期

```
HTTP Request → Gin Engine → CORS Middleware → JWT Middleware 
    → Controller (参数校验) → Service (业务编排) 
    → DAO (GORM) → MySQL / Redis / MinIO / RabbitMQ
```

---

## 📁 项目结构

```
.
├── bilibili-backend/          # Go 后端
│   ├── config/                # Viper 配置管理
│   ├── consumer/              # RabbitMQ 消费者（视频转码）
│   ├── controller/            # HTTP Handler（参数校验 & 响应封装）
│   ├── dao/                   # 数据访问层（GORM）
│   ├── middleware/            # Gin 中间件（JWT / CORS / OptionalAuth）
│   ├── model/                 # 数据模型（GORM Struct）
│   ├── router/                # 路由注册
│   ├── service/               # 业务逻辑层
│   ├── utils/                 # 工具包（JWT / MinIO / Redis / RabbitMQ / WebSocket）
│   ├── main.go                # 程序入口
│   └── Dockerfile
│
├── bilibili-frontend/         # Vue3 前端
│   ├── src/
│   ├── package.json
│   └── Dockerfile
│
├── nginx/                     # Nginx 配置
├── scripts/                   # 部署脚本
└── README.md
```

### 后端分层规范

| 层级 | 职责 | 对应目录 |
|------|------|----------|
| Controller | 接收 HTTP 请求，参数校验，调用 Service | `controller/` |
| Service | 业务逻辑编排，事务管理，调用 DAO 与外部服务 | `service/` |
| DAO | 数据库 CRUD，GORM 链式操作 | `dao/` |
| Model | 实体定义，GORM 标签与索引 | `model/` |
| Middleware | 鉴权、跨域、日志、Recovery | `middleware/` |

---

## 🚀 快速开始

### 环境依赖

- **Go** >= 1.25
- **Node.js** >= 20
- **MySQL** >= 8.0
- **Redis** >= 7.0
- **RabbitMQ** >= 3.12
- **MinIO** >= 8.0
- **FFmpeg** >= 5.0（用于视频转码与封面截取）

### 后端启动

```bash
cd bilibili-backend

# 1. 安装依赖
go mod download

# 2. 修改配置文件 config/config.yaml（数据库 / Redis / MinIO / RabbitMQ / JWT）
cp config/config.yaml.example config/config.yaml

# 3. 创建数据库并执行初始化脚本
mysql -u root -p < script/init.sql

# 4. 启动服务
go run main.go

# 或编译后运行
go build -o server main.go
./server
```

服务默认运行在 `http://localhost:8080`。

### 前端启动

```bash
cd bilibili-frontend

# 1. 安装依赖
npm install

# 2. 启动开发服务器
npm run dev

# 或同时启动前后端（需配置后端路径）
npm run dev:all
```

前端默认运行在 `http://localhost:5173`。

### Docker 部署

```bash
# 构建镜像
docker build -t bilibili-backend ./bilibili-backend
docker build -t bilibili-frontend ./bilibili-frontend

# 运行（需提前启动 MySQL / Redis / RabbitMQ / MinIO）
docker run -d -p 8080:8080 --name backend bilibili-backend
docker run -d -p 80:80 --name frontend bilibili-frontend
```

> 生产环境推荐使用 `docker-compose` 编排全部依赖服务。

---

## ⚙ 配置说明

核心配置项位于 `bilibili-backend/config/config.yaml`：

```yaml
server:
  port: "8080"
  mode: "debug"          # debug / release

database:
  driver: "mysql"
  host: "127.0.0.1"
  port: "3306"
  username: "bilibili"
  password: "your_password"
  database: "bilibili"
  charset: "utf8mb4"

jwt:
  secret: "your-secret-key"    # 生产环境请使用强随机字符串
  expires_in: 7200             # Token 有效期（秒）

minio:
  endpoint: "127.0.0.1:9000"
  access_key: "minio"
  secret_key: "minio123456"
  use_ssl: false
  bucket_videos: "videos"
  bucket_covers: "covers"
  bucket_avatars: "avatars"

rabbitmq:
  host: "127.0.0.1"
  port: "5672"
  username: "bilibili"
  password: "your_password"
  queue: "video_transcode"

redis:
  host: "127.0.0.1"
  port: "6379"
  password: "your_password"
  db: 0
```

⚠️ **安全提示**：生产环境请将敏感信息（密码、JWT Secret）通过环境变量注入，避免直接写入配置文件并提交到版本控制。

---

## 📡 API 概览

### 认证
| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/auth/register` | 用户注册 |
| POST | `/api/v1/auth/login` | 用户登录 |

### 用户
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/me` | 当前用户信息 |
| PUT | `/api/v1/users/me` | 更新资料 |
| PUT | `/api/v1/users/me/password` | 修改密码 |
| POST | `/api/v1/users/me/avatar` | 上传头像 |
| GET | `/api/v1/users/me/videos` | 我的投稿 |
| GET | `/api/v1/users/me/history` | 播放历史 |
| GET | `/api/v1/users/me/favorites` | 我的收藏 |

### 视频
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/videos` | 视频列表（支持 category / sort / page） |
| GET | `/api/v1/videos/:id` | 视频详情 |
| POST | `/api/v1/videos` | 上传视频 |
| PUT | `/api/v1/videos/:id` | 更新视频 |
| DELETE | `/api/v1/videos/:id` | 删除视频 |
| GET | `/api/v1/videos/:id/transcode` | 转码状态查询 |
| POST | `/api/v1/videos/:id/like` | 点赞/取消点赞 |
| GET | `/api/v1/videos/:id/like/status` | 点赞状态 |
| POST | `/api/v1/videos/:id/favorite` | 收藏/取消收藏 |
| GET | `/api/v1/videos/:id/favorite/status` | 收藏状态 |

### 评论
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/videos/:id/comments` | 视频评论列表 |
| POST | `/api/v1/videos/:id/comments` | 发表评论（含回复） |
| DELETE | `/api/v1/comments/:id` | 删除评论 |
| POST | `/api/v1/comments/:id/like` | 点赞评论 |

### 弹幕
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ws/danmaku?video_id={id}` | WebSocket 弹幕连接 |
| GET | `/api/v1/videos/:id/danmaku` | 获取时间段弹幕 |
| POST | `/api/v1/videos/:id/danmaku` | 发送弹幕 |

### 通知
| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/users/me/notifications` | 通知列表 |
| GET | `/api/v1/users/me/notifications/unread-count` | 未读数 |
| PUT | `/api/v1/users/me/notifications/:id/read` | 标记已读 |
| PUT | `/api/v1/users/me/notifications/read-all` | 全部已读 |
| DELETE | `/api/v1/users/me/notifications/:id` | 删除通知 |

---

## 🔥 核心实现亮点

### 1. 异步视频转码流水线

视频上传后，原片存入 MinIO，同时向 RabbitMQ 投递转码任务。消费者下载原片、执行 FFmpeg 转码（`scale=-2:480`）、上传转码后文件、更新数据库状态。全程异步，上传接口即时返回，用户可通过轮询 `transcode_status` 查看进度。

### 2. Redis 缓存 + 定时回写点赞数据

用户点赞时，使用 Redis Pipeline 原子更新三处数据：用户点赞 Hash、视频计数器、待同步 Set。每 30 秒定时任务将待同步视频的点赞数从 MySQL `COUNT` 后批量写入 `video.like_count`，显著降低高并发场景下的数据库写压力。

### 3. WebSocket 弹幕分房间广播

`DanmakuHub` 按 `video_id` 管理连接分组，广播时遍历同房间所有连接并自动清理失效连接。配合心跳检测机制，保证长连接稳定性。新弹幕通过独立 goroutine 异步广播，不阻塞 HTTP 响应。

### 4. 楼中楼评论聚合查询

评论表采用 `parent_id`（直接回复目标）+ `root_id`（所属一级评论）设计。查询时先分页获取一级评论，再批量聚合前 3 条二级回复，兼顾查询性能与展示完整性。

---

## 🔧 待优化项

- [ ] 引入 Refresh Token 机制，提升登录态续期体验
- [ ] 首页视频列表增加 Redis 缓存，减轻数据库压力
- [ ] 评论查询 N+1 问题优化为批量 IN 查询
- [ ] 视频搜索功能（MySQL FULLTEXT 或 Elasticsearch）
- [ ] 多码率转码（360p / 480p / 720p）与前端自适应切换
- [ ] 弹幕 Redis 缓冲 + 异步批量落库
- [ ] WebSocket 跨节点广播（Redis Pub/Sub）
- [ ] MinIO 预签名 URL 防盗链
- [ ] 结构化日志（Zap）与分布式追踪

---

## 📦 技术栈版本

| 组件 | 版本 | 用途 |
|------|------|------|
| Go | 1.25 | 后端主语言 |
| Gin | 1.12.0 | Web 框架 |
| GORM | 1.31.1 | ORM |
| MySQL Driver | 1.6.0 | 数据库驱动 |
| JWT | v5.3.1 | Token 认证 |
| bcrypt | v0.51.0 | 密码哈希 |
| MinIO Go SDK | v7.1.0 | 对象存储 |
| RabbitMQ Client | v1.11.0 | 消息队列 |
| go-redis | v9.19.0 | Redis 客户端 |
| Gorilla WebSocket | v1.5.3 | 实时通信 |
| Vue | 3.5.34 | 前端框架 |
| Vite | 8.0.12 | 构建工具 |
| Pinia | 3.0.4 | 状态管理 |
| DPlayer | 1.27.1 | 视频播放器 |

---

## 📄 License

本项目仅供学习交流使用。

---

> 如果这个项目对你有帮助，欢迎点个 ⭐ Star 支持一下！
