export interface CommentUser {
  id: number
  username: string
  nickname: string
  avatar: string
}

export interface ReplyVO {
  id: number
  content: string
  like_count: number
  is_liked: boolean
  created_at: string
  user: CommentUser
  to_user: CommentUser
}

export interface CommentVO {
  id: number
  content: string
  like_count: number
  is_liked: boolean
  created_at: string
  user: CommentUser
  reply_count: number
  replies: ReplyVO[]
}

export interface CommentListResponse {
  list: CommentVO[]
  total: number
  page: number
  size: number
}
