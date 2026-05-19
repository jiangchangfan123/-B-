import request from './request'
import type { LoginForm, RegisterForm, LoginResponse, RegisterResponse, UserInfo } from '../types/auth'

// 基础配置
const BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

// 关闭 Mock，走真实后端
const USE_MOCK = false

// ========== Mock 数据 ==========
const MOCK_TOKEN = 'mock_jwt_token_' + Date.now()
const MOCK_USER: UserInfo = {
  id: 1,
  username: 'void_user',
  avatar: '',
  role: 1,
  email: 'user@nebula.tv',
  sign: '深空漫步者',
  coins: 42,
}

let mockRegistered = false

// ========== API ==========

/**登录 */
export async function login(form: LoginForm): Promise<LoginResponse> {
  if (USE_MOCK) {
    await delay(500)
    if (!form.username || !form.password) {
      throw new Error('> ERROR: 请输入用户名和密码')
    }
    if (form.username.length < 3 || form.username.length > 20) {
      throw new Error('> ERROR: 用户名格式不正确')
    }
    return {
      accessToken: MOCK_TOKEN,
      userInfo: { ...MOCK_USER, username: form.username },
    }
  }
  return request.post('/auth/login', form)
}

/**注册 */
export async function register(form: RegisterForm): Promise<RegisterResponse> {
  if (USE_MOCK) {
    await delay(800)
    if (!form.username || !form.password || !form.email) {
      throw new Error('> ERROR: 请填写所有字段')
    }
    if (form.username.length < 3 || form.username.length > 20) {
      throw new Error('> ERROR: 用户名格式不正确')
    }
    if (form.password.length < 6 || form.password.length > 30) {
      throw new Error('> ERROR: 密码长度不正确')
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(form.email)) {
      throw new Error('> ERROR: 邮箱格式不正确')
    }
    if (mockRegistered && form.username === MOCK_USER.username) {
      throw new Error('> ERROR: 用户名已存在')
    }
    mockRegistered = true
    return { id: Date.now(), username: form.username }
  }
  return request.post('/auth/register', form)
}

/**获取当前用户 */
export async function getMe(): Promise<UserInfo> {
  if (USE_MOCK) {
    await delay(300)
    return { ...MOCK_USER }
  }
  return request.get('/auth/me')
}

// 工具函数
function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms))
}
