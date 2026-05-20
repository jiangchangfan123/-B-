// =========================
// 视频模块 TypeScript 类型定义
// =========================

/** 视频上传表单 */
export interface UploadForm {
  file: File
  title: string
  description?: string
  category: string
  cover?: File
}

/** 视频上传响应 */
export interface UploadResponse {
  id: number
  title: string
  cover_url: string
  status: number
  transcode_status: number
  created_at: string
}

/** 转码状态响应 */
export interface TranscodeStatusResponse {
  transcode_status: number
  transcoded_url: string
}

/** 视频列表项 */
export interface VideoListItem {
  id: number
  title: string
  cover_url: string
  category: string
  view_count: number
  like_count: number
  created_at: string
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

/** 视频列表响应 */
export interface VideoListResponse {
  list: VideoListItem[]
  total: number
  page: number
  size: number
}

/** 视频详情 */
export interface VideoDetail {
  id: number
  title: string
  description: string
  cover_url: string
  video_url: string
  transcoded_url: string
  category: string
  status: number
  transcode_status: number
  view_count: number
  like_count: number
  comment_count: number
  danmaku_count: number
  created_at: string
  user_info: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

/** 分类选项 */
export interface CategoryOption {
  id: string
  label: string
  code: string
}

export const videoCategories: CategoryOption[] = [
  { id: 'cine', label: '影视', code: '[CINE]' },
  { id: 'game', label: '游戏', code: '[GAME]' },
  { id: 'acgn', label: '二次元', code: '[ACGN]' },
  { id: 'chef', label: '美食', code: '[CHEF]' },
  { id: 'docu', label: '纪录片', code: '[DOCU]' },
  { id: 'show', label: '综艺', code: '[SHOW]' },
]
