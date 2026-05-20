package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"bilibili-backend/dao"
	"bilibili-backend/model"
	"bilibili-backend/utils"
)

type CommentService struct {
	commentDao          *dao.CommentDao
	userDao             *dao.UserDao
	notificationService *NotificationService
}

func NewCommentService(commentDao *dao.CommentDao, userDao *dao.UserDao, notificationService *NotificationService) *CommentService {
	return &CommentService{commentDao: commentDao, userDao: userDao, notificationService: notificationService}
}

// CreateComment 创建评论（一级或二级回复）
func (s *CommentService) CreateComment(videoID, userID uint64, content string, parentID uint64) (*model.Comment, error) {
	comment := &model.Comment{
		VideoID:  videoID,
		UserID:   userID,
		Content:  content,
		ParentID: parentID,
		RootID:   0,
	}

	// 如果是二级回复，需要确定 root_id
	if parentID > 0 {
		parent, err := s.commentDao.GetByID(parentID)
		if err != nil {
			return nil, fmt.Errorf("parent comment not found")
		}
		if parent.RootID > 0 {
			comment.RootID = parent.RootID
		} else {
			comment.RootID = parent.ID
		}
	}

	if err := s.commentDao.Create(comment); err != nil {
		return nil, err
	}

	// 加载用户信息
	user, _ := s.userDao.GetByID(userID)
	comment.User = *user

	// 发送通知：回复评论
	if parentID > 0 && s.notificationService != nil {
		parent, _ := s.commentDao.GetByID(parentID)
		if parent != nil && parent.UserID != userID {
			go s.notificationService.CreateCommentReplyNotification(videoID, comment.ID, userID, parent.UserID, content)
		}
	}

	return comment, nil
}

// CommentVO 一级评论视图
type CommentVO struct {
	ID         uint64    `json:"id"`
	Content    string    `json:"content"`
	LikeCount  int       `json:"like_count"`
	IsLiked    bool      `json:"is_liked"`
	CreatedAt  time.Time `json:"created_at"`
	User       UserVO    `json:"user"`
	ReplyCount int64     `json:"reply_count"`
	Replies    []ReplyVO `json:"replies"`
}

// ReplyVO 二级回复视图
type ReplyVO struct {
	ID        uint64    `json:"id"`
	Content   string    `json:"content"`
	LikeCount int       `json:"like_count"`
	IsLiked   bool      `json:"is_liked"`
	CreatedAt time.Time `json:"created_at"`
	User      UserVO    `json:"user"`
	ToUser    UserVO    `json:"to_user"`
}

// UserVO 用户简信息
type UserVO struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// GetVideoComments 获取视频评论列表
func (s *CommentService) GetVideoComments(videoID uint64, page, size int, sort string, currentUserID uint64) ([]CommentVO, int64, error) {
	comments, total, err := s.commentDao.GetTopLevelByVideoID(videoID, page, size, sort)
	if err != nil {
		return nil, 0, err
	}

	result := make([]CommentVO, 0, len(comments))
	for _, c := range comments {
		vo := CommentVO{
			ID:        c.ID,
			Content:   c.Content,
			LikeCount: c.LikeCount,
			IsLiked:   s.isCommentLiked(currentUserID, c.ID),
			CreatedAt: c.CreatedAt,
			User: UserVO{
				ID:       c.User.ID,
				Username: c.User.Username,
				Nickname: c.User.Nickname,
				Avatar:   c.User.Avatar,
			},
		}

		// 获取回复数量
		replyCount, _ := s.commentDao.GetReplyCountByRootID(c.ID)
		vo.ReplyCount = replyCount

		// 获取前3条回复
		replies, _ := s.commentDao.GetRepliesByRootID(c.ID)
		vo.Replies = make([]ReplyVO, 0, len(replies))
		for _, r := range replies {
			// 找到回复对象的用户信息
			toUser, _ := s.userDao.GetByID(r.UserID)
			// 如果是回复一级评论，to_user 就是一级评论的作者
			if r.ParentID == c.ID {
				parent, _ := s.commentDao.GetByIDWithUser(r.ParentID)
				if parent != nil {
					toUser = &parent.User
				}
			} else {
				// 回复的回复，to_user 是被回复的那个人
				parentReply, _ := s.commentDao.GetByIDWithUser(r.ParentID)
				if parentReply != nil {
					toUser = &parentReply.User
				}
			}

			replyVO := ReplyVO{
				ID:        r.ID,
				Content:   r.Content,
				LikeCount: r.LikeCount,
				IsLiked:   s.isCommentLiked(currentUserID, r.ID),
				CreatedAt: r.CreatedAt,
				User: UserVO{
					ID:       r.User.ID,
					Username: r.User.Username,
					Nickname: r.User.Nickname,
					Avatar:   r.User.Avatar,
				},
			}
			if toUser != nil {
				replyVO.ToUser = UserVO{
					ID:       toUser.ID,
					Username: toUser.Username,
					Nickname: toUser.Nickname,
					Avatar:   toUser.Avatar,
				}
			}
			vo.Replies = append(vo.Replies, replyVO)
		}

		result = append(result, vo)
	}

	return result, total, nil
}

// DeleteComment 删除评论（只能删除自己的）
func (s *CommentService) DeleteComment(commentID, userID uint64) error {
	comment, err := s.commentDao.GetByID(commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return fmt.Errorf("not owner")
	}
	return s.commentDao.Delete(commentID)
}

// ToggleCommentLike 点赞/取消点赞评论
func (s *CommentService) ToggleCommentLike(commentID, userID uint64) (bool, error) {
	ctx := context.Background()
	userKey := fmt.Sprintf("comment:likes:user:%d", userID)
	field := strconv.FormatUint(commentID, 10)

	liked := false

	// 查当前状态
	if utils.RedisClient != nil {
		val, err := utils.RedisClient.HGet(ctx, userKey, field).Result()
		if err == nil && val == "1" {
			liked = true
		}
	}

	if liked {
		// 取消点赞
		if utils.RedisClient != nil {
			utils.RedisClient.HSet(ctx, userKey, field, 0)
			utils.RedisClient.Decr(ctx, fmt.Sprintf("comment:likes:count:%d", commentID))
			utils.RedisClient.Expire(ctx, userKey, 24*time.Hour)
		}
		_ = s.commentDao.IncrementLikeCount(commentID, -1)
		return false, nil
	}

	// 点赞
	if utils.RedisClient != nil {
		utils.RedisClient.HSet(ctx, userKey, field, 1)
		utils.RedisClient.Incr(ctx, fmt.Sprintf("comment:likes:count:%d", commentID))
		utils.RedisClient.Expire(ctx, userKey, 24*time.Hour)
	}
	_ = s.commentDao.IncrementLikeCount(commentID, 1)

	// 发送通知：评论被点赞
	if s.notificationService != nil {
		comment, _ := s.commentDao.GetByID(commentID)
		if comment != nil && comment.UserID != userID {
			go s.notificationService.CreateCommentLikedNotification(commentID, userID, comment.UserID)
		}
	}

	return true, nil
}

// isCommentLiked 检查用户是否点赞了该评论
func (s *CommentService) isCommentLiked(userID, commentID uint64) bool {
	if userID == 0 {
		return false
	}
	ctx := context.Background()
	if utils.RedisClient != nil {
		val, err := utils.RedisClient.HGet(ctx, fmt.Sprintf("comment:likes:user:%d", userID), strconv.FormatUint(commentID, 10)).Result()
		if err == nil {
			return val == "1"
		}
	}
	return false
}
