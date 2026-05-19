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
  avatar: string
  role: number
  email?: string
  sign?: string
  coins?: number
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

/** Token 存储键 */
export const TOKEN_KEY = 'nebula_access_token'
