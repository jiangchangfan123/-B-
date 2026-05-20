# Phase 3 实施计划

## 目标
完成剩余核心功能：简化弹幕系统、消息推送系统、主页分类筛选、左侧导航栏精简。

---

## 模块一：简化弹幕系统

### 设计原则
仅实现核心三个能力：发送弹幕、接收弹幕、滚动显示。不做弹幕屏蔽、关键词过滤、高级弹幕。

### 后端

#### 1.1 数据库表 `danmaku`
```sql
CREATE TABLE danmaku (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id    BIGINT NOT NULL,
    user_id     BIGINT NOT NULL,
    content     VARCHAR(100) NOT NULL,
    time_point  INT DEFAULT 0 COMMENT '视频时间点（秒）',
    color       VARCHAR(7) DEFAULT '#ffffff',
    type        TINYINT DEFAULT 1 COMMENT '1=滚动 2=顶部 3=底部',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_video_time (video_id, time_point)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 1.2 WebSocket 接入
- 引入 `gorilla/websocket`
- 新建 `utils/websocket.go`：管理连接池，按 `video_id` 分组广播
- `ws://localhost:8080/ws/danmaku?video_id=5` 连接
- 心跳机制：每 30 秒 ping/pong

#### 1.3 API 设计
| 接口 | 方法 | 说明 |
|---|---|---|
| `WS /ws/danmaku?video_id={id}` | WebSocket | 连接弹幕信道 |
| `GET /api/v1/videos/:id/danmaku?start=&end=` | HTTP | 获取某时间段弹幕（用于回放时加载） |
| `POST /api/v1/videos/:id/danmaku` | HTTP | 发送弹幕（同时通过 WebSocket 广播） |

#### 1.4 代码文件
- `model/danmaku.go` — 数据模型
- `dao/danmaku.go` — DAO 层（Create、GetByVideoIDTimeRange）
- `service/danmaku.go` — 业务逻辑
- `controller/danmaku.go` — HTTP + WebSocket handler
- `utils/websocket.go` — 连接管理器

### 前端

#### 1.5 弹幕组件 `DanmakuLayer.vue`
- 覆盖在 `VideoPlayer` 上半部分，高度约占播放器 50%
- 绝对定位 `pointer-events: none`，只有输入框区域响应鼠标
- 弹幕样式：白色文字，从右向左横向滚动
- 滚动动画：CSS `@keyframes` 水平位移
- 轨道算法：简化版 — 随机分布在上半部多条轨道，避免重叠用 `top: N%`
- 同屏最大弹幕数限制：100 条，超出则移除最早的
- WebSocket 连接、断线重连、心跳保活
- 输入框：放在播放器底部或右侧，回车发送弹幕，登录后可用

#### 1.6 集成点
- `VideoDetailView.vue` 中在播放器区域加入 `<DanmakuLayer>`
- 播放进度同步：视频播放到 X 秒时，加载该时间段弹幕

---

## 模块二：消息推送系统

### 设计原则
只做三种通知类型，使用 MySQL + HTTP 轮询（不用 WebSocket 或 SSE，降低复杂度）。前端每 30 秒轮询未读数量，页面切换时刷新消息列表。

### 后端

#### 2.1 数据库表 `notifications`
```sql
CREATE TABLE notifications (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id     BIGINT NOT NULL COMMENT '接收者',
    type        TINYINT NOT NULL COMMENT '1=评论回复 2=评论被点赞 3=视频被点赞',
    title       VARCHAR(100) NOT NULL,
    content     VARCHAR(255),
    related_id  BIGINT COMMENT '关联ID（视频ID或评论ID）',
    is_read     TINYINT DEFAULT 0,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_read (user_id, is_read),
    INDEX idx_user_created (user_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 2.2 触发点
| 业务场景 | 触发位置 | 接收者 |
|---|---|---|
| A 回复了 B 的评论 | `CommentService.CreateComment` | B（被回复者） |
| A 点赞了 B 的评论 | `CommentService.ToggleCommentLike` | B（评论作者） |
| A 点赞了 B 的视频 | `LikeService.LikeVideo` | B（视频作者） |

#### 2.3 API 设计
| 接口 | 方法 | 说明 |
|---|---|---|
| `GET /api/v1/users/me/notifications` | GET | 获取消息列表（支持未读筛选） |
| `GET /api/v1/users/me/notifications/unread-count` | GET | 未读数量（NavBar 红点用） |
| `PUT /api/v1/users/me/notifications/:id/read` | PUT | 标记单条已读 |
| `PUT /api/v1/users/me/notifications/read-all` | PUT | 全部已读 |
| `DELETE /api/v1/users/me/notifications/:id` | DELETE | 删除消息 |

#### 2.4 代码文件
- `model/notification.go` — 数据模型
- `dao/notification.go` — DAO 层
- `service/notification.go` — 业务逻辑 + 触发包装
- `controller/notification.go` — API handler
- `middleware/notification_trigger.go` — 可选，用于在 Service 层注入触发

### 前端

#### 2.5 API 层 `api/notification.ts`
- 封装上述 5 个接口

#### 2.6 NavBar 红点
- `NavBar.vue` 增加消息图标 + 未读数量红点
- 页面加载时调用 `GET unread-count`
- 每 30 秒轮询更新

#### 2.7 消息中心页 `NotificationsView.vue`
- 路径 `/notifications`
- 消息列表：类型标签 + 内容 + 时间
- 未读消息高亮显示
- 点击跳转：视频被点赞 → 视频详情页，评论被回复 → 视频详情页置底评论
- 顶部按钮：全部已读、删除已读

#### 2.8 路由注册
- `/notifications` 路由，需登录

---

## 模块三：主页分类筛选

### 现状问题
- `mock/homeData.ts` 中 `categoryList` 只有 6 个分类
- `HomeView.vue` 的 `categoryPaletteMap` 有 10 个分类
- 后端 `ListPublished` 已支持 `category` 参数筛选
- 但主页没有分类标签栏，无法点击分类筛选

### 实施方案

#### 3.1 统一分类列表
前后端统一使用 6 个分类（与原 `mock/homeData.ts` 保持一致）：
```
all, cine, game, acgn, chef, docu, show
```
对应中文：全部、影视、游戏、二次元、美食、纪录片、综艺

#### 3.2 新建分类标签栏组件 `CategoryFilter.vue`
- 横向滚动条，显示所有分类标签
- 点击切换 `activeCategory`，发出 `filter-change` 事件
- 当前选中标签底部有青色下划线
- 放在 `HomeView.vue` 中 `BannerCarousel` 下方

#### 3.3 HomeView 修改
- 增加 `activeCategory` ref，默认 `all`
- `loadVideos` 根据 `activeCategory` 传参给 `getVideoList`
- 分类变化时重新加载视频列表
- `SectionBlock` 的标题动态显示当前分类名

#### 3.4 API 层修改
- `api/video.ts` 中 `getVideoList` 已支持 `category` 参数，确认前端正确传参即可

---

## 模块四：左侧导航栏精简

### 现状分析
- 当前主菜单：首页、收藏、动态、热门、频道、分区
- 系统菜单：设置、主题切换
- 路由映射：部分页面不存在（feed、hot、channel、partition、theme）

### 精简方案

#### 4.1 保留的菜单项
| 位置 | 项目 | 跳转 |
|---|---|---|
| 主菜单 | 首页 | `/` |
| 主菜单 | 收藏 | `/favorites` |
| 主菜单 | 消息 | `/notifications` （新增） |
| 系统区 | 个人中心 | `/personal` |

#### 4.2 移除的菜单项
- 动态、热门、频道、分区、主题切换
- 原“设置”改为“个人中心”，更直观

#### 4.3 修改文件
- `mock/homeData.ts` — 重写 `menuList` 和 `systemMenuList`
- `SideBar.vue` — 修改 `routeMap`，移除无效映射
- 如果某个菜单页面不存在且不需要，可以去掉路由

---

## 执行顺序建议

| 顶序 | 模块 | 理由 |
|---|---|---|
| 1 | 分类筛选 + 导航栏精简 | 工作量小，纯前端修改，快速出效果 |
| 2 | 消息推送 | 后端增加表 + API，前端增加页面，与现有功能联动 |
| 3 | 弹幕系统 | 工作量最大，涉及 WebSocket，放最后 |

---

## 需要新增的后端文件

```
model/
  danmaku.go
  notification.go
dao/
  danmaku.go
  notification.go
service/
  danmaku.go
  notification.go
controller/
  danmaku.go
  notification.go
utils/
  websocket.go  (新建)
```

## 需要新增/修改的前端文件

```
components/
  CategoryFilter.vue    (新建)
  DanmakuLayer.vue      (新建)
views/
  NotificationsView.vue (新建)
api/
  danmaku.ts            (新建)
  notification.ts       (新建)
mock/homeData.ts        (修改分类列表 + 菜单)
views/HomeView.vue      (加入分类筛选)
components/SideBar.vue  (精简菜单)
components/NavBar.vue   (增加消息红点)
router/index.ts         (注册新页面)
```

---

## 后端数据库迁移

`main.go` 中 `AutoMigrate` 增加：
```go
&model.Danmaku{}, &model.Notification{}
```

---

## 风险与应对

| 风险 | 应对 |
|---|---|
| WebSocket 在生产环境下需要配置 Nginx 代理 | 计划中注明为开发环境配置，部署时补充 |
| 弹幕过多导致浏览器卡顿 | 设置同屏最大弹幕数（100 条），超出则移除 |
| 轮询消息对服务器造成压力 | 未读数量可缓存在 localStorage，30 秒轮询一次压力可控 |
| 弹幕覆盖播放器控制条 | 弹幕层只覆盖上半部分（50% 高度），避开底部控制条区域 |
