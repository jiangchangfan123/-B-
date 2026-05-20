# 弹幕系统实施计划

> 📌 核心结论：基于 WebSocket 实现简化弹幕系统，包含发送弹幕、接收弹幕、滚动显示三大核心能力，前后端均从零新建文件实现。

---

## 一、后端实现

### 1.1 数据库表 `danmaku`

```sql
CREATE TABLE danmaku (
    id          BIGINT PRIMARY KEY AUTO_INCREMENT,
    video_id    BIGINT NOT NULL COMMENT '关联视频ID',
    user_id     BIGINT NOT NULL COMMENT '发送者ID',
    content     VARCHAR(100) NOT NULL COMMENT '弹幕内容',
    time_point  INT DEFAULT 0 COMMENT '视频时间点（秒）',
    color       VARCHAR(7) DEFAULT '#ffffff' COMMENT '弹幕颜色',
    type        TINYINT DEFAULT 1 COMMENT '1=滚动 2=顶部 3=底部',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_video_time (video_id, time_point)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

> ⚠️ 注意：`AutoMigrate` 中需追加 `&model.Danmaku{}`。

### 1.2 WebSocket 连接管理器

- 引入 `gorilla/websocket` 依赖
- 新建 `utils/websocket.go`：
  - 管理全局连接池（`map[*websocket.Conn]bool`）
  - 按 `video_id` 分组广播（`map[int64]map[*websocket.Conn]bool`）
  - 心跳机制：每 30 秒服务端发送 `ping`，客户端回复 `pong`
  - 断线清理：连接关闭时从分组中移除

### 1.3 API 设计

| 接口 | 方法 | 说明 |
|---|---|---|
| `WS /ws/danmaku?video_id={id}` | WebSocket | 连接弹幕信道，接收实时弹幕 |
| `GET /api/v1/videos/:id/danmaku?start=&end=` | HTTP | 获取某时间段弹幕，用于回放时加载 |
| `POST /api/v1/videos/:id/danmaku` | HTTP | 发送弹幕，同时通过 WebSocket 广播给同视频所有连接 |

### 1.4 后端文件清单

| 文件 | 职责 |
|---|---|
| `model/danmaku.go` | `Danmaku` 结构体定义 |
| `dao/danmaku.go` | `Create(danmaku)`、`GetByVideoIDTimeRange(videoID, start, end)` |
| `service/danmaku.go` | 保存弹幕 + 调用 WebSocket 广播的业务逻辑 |
| `controller/danmaku.go` | WebSocket Handler（升级连接、心跳、读写 goroutine）、HTTP Handler（获取/发送弹幕） |
| `utils/websocket.go` | 连接管理器、分组广播、心跳、断线清理 |

### 1.5 路由注册

在 `routes/routes.go` 中注册：

- `GET /ws/danmaku` — `controller.DanmakuWebSocket`
- `GET /api/v1/videos/:id/danmaku` — `controller.GetDanmaku`
- `POST /api/v1/videos/:id/danmaku` — `controller.SendDanmaku`（需登录）

---

## 二、前端实现

### 2.1 弹幕组件 `DanmakuLayer.vue`

覆盖在 `VideoPlayer` 上层，核心功能：

- **布局**：绝对定位，覆盖播放器上半部分（`height: 50%`），`pointer-events: none`
- **弹幕样式**：白色文字 + 黑色描边阴影，从右向左水平滚动
- **滚动动画**：CSS `@keyframes danmaku-move` 控制水平位移
- **轨道算法**：简化版 — 弹幕随机分布在 `top: 5%~45%` 的 N 条轨道上，通过 `top` 百分比避免重叠
- **数量限制**：同屏最大 100 条，超出时移除最早的弹幕 DOM 节点
- **输入框**：放在播放器底部，回车发送，未登录时隐藏或提示登录
- **WebSocket 管理**：连接、断线重连、心跳保活、收到消息后创建弹幕元素

### 2.2 播放进度同步

- `VideoDetailView.vue` 监听视频 `timeupdate` 事件
- 当播放到新秒数时，若该秒弹幕未加载过，调用 `GET /api/v1/videos/:id/danmaku?start=当前秒&end=当前秒+1` 加载历史弹幕
- 避免每秒重复请求，用 `Set` 记录已加载的时间区间

### 2.3 API 层 `api/danmaku.ts`

```typescript
// 获取时间段弹幕
getDanmaku(videoID: number, start: number, end: number)

// 发送弹幕
sendDanmaku(videoID: number, data: { content: string; time_point: number; color?: string; type?: number })
```

### 2.4 前端文件清单

| 文件 | 职责 |
|---|---|
| `components/DanmakuLayer.vue` | 弹幕展示层、WebSocket 连接、输入框、弹幕生命周期管理 |
| `api/danmaku.ts` | 弹幕相关 HTTP API 封装 |
| `views/VideoDetailView.vue` | 集成 `<DanmakuLayer>`，传入 `videoId`、同步播放进度 |

---

## 三、执行顺序

| 顺序 | 任务 | 理由 |
|---|---|---|
| 1 | 后端数据库表 + Model + DAO | 底层先就绪 |
| 2 | 后端 `utils/websocket.go` | 核心基础设施 |
| 3 | 后端 Service + Controller + 路由注册 | 完整后端 API |
| 4 | 前端 `api/danmaku.ts` | API 封装层 |
| 5 | 前端 `DanmakuLayer.vue` | 核心弹幕组件 |
| 6 | 前端集成到 `VideoDetailView.vue` | 联调测试 |
| 7 | 前后端联调、测试心跳/断线重连 | 验证稳定性 |

---

## 四、风险与应对

| 风险 | 应对 |
|---|---|
| WebSocket 连接过多导致内存溢出 | 断线时立即清理连接池，限制单视频最大连接数（可暂不实现） |
| 弹幕过多导致浏览器卡顿 | 同屏限制 100 条，CSS 动画优先使用 `transform` 开启 GPU 加速 |
| 弹幕覆盖播放器控制条 | 弹幕层高度限制为 50%，避开底部控制条 |
| 心跳失效导致僵尸连接 | 服务端 30 秒 ping，客户端 pong，超时未响应则主动关闭连接 |
| 跨域问题（WSL 环境） | WebSocket 升级时允许所有 origin，生产环境再限制 |

---

## 行动清单

- [ ] 确认本计划后，从 Step 1（后端数据库表 + Model）开始执行
- [ ] 每完成一个 Step 后简要汇报进度
- [ ] 全部完成后进行前后端联调测试
