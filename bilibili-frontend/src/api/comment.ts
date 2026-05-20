import request from './request'
import type { CommentListResponse, CommentVO } from '../types/comment'

/** 发表评论 */
export async function createComment(
  videoId: number,
  data: { content: string; parent_id?: number }
): Promise<CommentVO> {
  return request.post(`/videos/${videoId}/comments`, data)
}

/** 获取视频评论列表 */
export async function getVideoComments(
  videoId: number,
  page = 1,
  size = 20,
  sort = 'time'
): Promise<CommentListResponse> {
  return request.get(`/videos/${videoId}/comments?page=${page}&size=${size}&sort=${sort}`)
}

/** 删除评论 */
export async function deleteComment(commentId: number): Promise<void> {
  return request.delete(`/comments/${commentId}`)
}

/** 点赞/取消点赞评论 */
export async function toggleCommentLike(commentId: number): Promise<{ liked: boolean }> {
  return request.post(`/comments/${commentId}/like`)
}
