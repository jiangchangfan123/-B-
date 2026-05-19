import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi, getMe } from '../api/auth'
import type { LoginForm, RegisterForm, UserInfo } from '../types/auth'
import { TOKEN_KEY } from '../types/auth'

export const useUserStore = defineStore('user', () => {
  // ========== State ==========
  const token = ref<string>(localStorage.getItem(TOKEN_KEY) || '')
  const userInfo = ref<UserInfo | null>(null)

  // ========== Getters ==========
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)

  // ========== Actions ==========
  async function login(form: LoginForm) {
    const res = await loginApi(form)
    token.value = res.accessToken
    userInfo.value = res.userInfo
    localStorage.setItem(TOKEN_KEY, res.accessToken)
    return res
  }

  async function register(form: RegisterForm) {
    return registerApi(form)
  }

  async function fetchUserInfo() {
    if (!token.value) {
      logout()
      throw new Error('未登录')
    }
    try {
      const res = await getMe()
      userInfo.value = res
      return res
    } catch {
      logout()
      throw new Error('获取用户信息失败')
    }
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  function restore() {
    const saved = localStorage.getItem(TOKEN_KEY)
    if (saved) {
      token.value = saved
    }
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    login,
    register,
    fetchUserInfo,
    logout,
    restore,
  }
})
