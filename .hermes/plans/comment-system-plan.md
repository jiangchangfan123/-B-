# 评论系统实施计划

## 目标
为 B 站克隆项目实现完整的视频评论系统，支持一级评论 + 二级回复、评论点赞、按热度/时间排序。

## Phase 1: 数据库模型与迁移

### 1.1 创建评论表模型 `model/comment.go`
```go
type Comment struct {
    ID         uint64    `gorm:"primaryKey" json:"id"`
    VideoID    uint64    `gorm:"not null;index" json:"video_id"`
    UserID     uint64    `gorm:"not null;index" json:"user_id"`
    Content    string    `gorm:"type:text;not null" json:"content"`
    ParentID   uint64    `gorm:"default:0;index" json:"parent_id"` // 回复的评论ID，0=一级评论
    RootID     uint64    `gorm:"default:0;index" json:"root_id"`   // 一级评论ID
    LikeCount  int       `gorm:"default:0" json:"like_count"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    // 关联
    User       User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
```

### 1.2 在 `bootstrap/db.go` 中自动迁移
添加 `&model.Comment{}` 到 AutoMigrate 列表。

## Phase 2: 后端 DAO 层

### 2.1 `dao/comment.go`
- `Create(comment *model.Comment) error`
- `GetByVideoID(videoID uint64, page, size int, sort string) ([]model.Comment, int64, error)`
  - sort="hot" 按 like_count desc, created_at desc
  - sort="time" 按 created_at desc
  - 只查一级评论 (parent_id = 0)
- `GetRepliesByRootID(rootID uint64) ([]model.Comment, error)`
- `GetByID(id uint64) (*model.Comment, error)`
- `Delete(id uint64) error`
- `IncrementLikeCount(id uint64, delta int) error`

## Phase 3: 后端 Service 层

### 3.1 `service/comment.go`
- `CreateComment(videoID, userID uint64, content string, parentID uint64) (*model.Comment, error)`
  - parentID=0 → 一级评论，root_id=0
  - parentID>0 → 二级回复，需要查 parent 的 root_id
- `GetVideoComments(videoID uint64, page, size int, sort string) (*CommentListResult, error)`
  - 返回一级评论列表 + 每个一级评论的前3条回复 + 回复总数
- `DeleteComment(commentID, userID uint64) error`
  - 只能删除自己的评论
- `ToggleCommentLike(commentID, userID uint64) (bool, error)`
  - 使用 Redis 记录用户点赞状态（SET: comment_likes:{comment_id}）
  - 返回是否点赞 + 更新 like_count

### 3.2 返回结构
```go
type CommentVO struct {
    ID        uint64      `json:"id"`
    Content   string      `json:"content"`
    LikeCount int         `json:"like_count"`
    IsLiked   bool        `json:"is_liked"`
    CreatedAt time.Time   `json:"created_at"`
    User      UserInfoVO  `json:"user"`
    ReplyCount int        `json:"reply_count"`
    Replies   []ReplyVO   `json:"replies"` // 前3条回复
}

type ReplyVO struct {
    ID        uint64      `json:"id"`
    Content   string      `json:"content"`
    LikeCount int         `json:"like_count"`
    IsLiked   bool        `json:"is_liked"`
    CreatedAt time.Time   `json:"created_at"`
    User      UserInfoVO  `json:"user"`
    ToUser    UserInfoVO  `json:"to_user"` // 回复给谁
}
```

## Phase 4: 后端 Controller 与路由

### 4.1 `controller/comment.go`
- `Create` POST `/api/v1/videos/:id/comments`
- `List` GET `/api/v1/videos/:id/comments?page=&size=&sort=`
- `Delete` DELETE `/api/v1/comments/:id`
- `ToggleLike` POST `/api/v1/comments/:id/like`

### 4.2 `router.go` 注册路由

## Phase 5: 前端 API 层

### 5.1 `api/comment.ts`
```ts
export function createComment(videoId: number, data: { content: string; parent_id?: number })
export function getVideoComments(videoId: number, params: { page: number; size: number; sort: string })
export function deleteComment(commentId: number)
export function toggleCommentLike(commentId: number)
```

## Phase 6: 前端评论组件

### 6.1 `components/CommentSection.vue`
- 顶部：评论数 + 排序切换（最热 / 最新）
- 输入框：登录后显示，支持发表一级评论
- 一级评论列表：
  - 头像、用户名、时间、内容
  - [点赞] [回复] [删除（自己的）]
  - 展开回复：显示二级回复列表 + "展开更多回复"
  - 二级回复支持再回复（@某人）
- 分页加载："加载更多评论"

### 6.2 类型定义 `types/comment.ts`

## Phase 7: 视频详情页集成

### 7.1 `VideoView.vue`
- 底部接入 `<CommentSection :video-id="videoId" />`
- 评论数同步显示在视频信息区

## Phase 8: 测试验证

1. 发表一级评论
2. 回复一级评论（二级）
3. 点赞评论
4. 删除自己的评论
5. 切换排序（最热/最新）
6. 分页加载更多
