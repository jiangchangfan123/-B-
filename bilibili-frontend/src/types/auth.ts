// =========================
// 认证模块 TypeScript 类型定义
// =========================

/** 登录表单 */
export interface LoginForm {
  username: string
  password: string
}

/** 注册表单 */
export interface RegisterForm {
  username: string
  password: string
  email: string
}

/** 用户信息（不含敏感字段） */
export interface UserInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  role: number
  email?: string
  sign?: string
  coins?: number
  created_at?: string
}

/** 完整用户资料（个人中心用） */
export interface UserProfile {
  id: number
  username: string
  nickname: string
  email: string
  avatar: string
  sign: string
  role: number
  coins: number
  created_at: string
}

/** 更新资料请求 */
export interface UpdateProfileForm {
  sign?: string
  nickname?: string
}

/** 修改密码请求 */
export interface UpdatePasswordForm {
  oldPassword: string
  newPassword: string
  confirmPassword: string
}

/** 统一响应格式 */
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

/** 登录响应 */
export interface LoginResponse {
  accessToken: string
  userInfo: UserInfo
}

/** 注册响应 */
export interface RegisterResponse {
  id: number
  username: string
}

/** 视频项（个人中心列表用） */
export interface MyVideoItem {
  id: number
  title: string
  cover_url: string
  views: number
  status: number // 1=已发布 2=审核中 3=已封禁
  category: string
  created_at: string
}

/** 播放历史项 */
export interface HistoryItem {
  id: number
  title: string
  cover_url: string
  views: number
  watched_at: string
}

/** Token 存储键 */
export const TOKEN_KEY = 'nebula_access_token'
