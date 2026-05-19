import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, register as registerApi } from '../api/auth'
import { getMeDetail, updateProfile, uploadAvatar, updatePassword } from '../api/user'
import type { LoginForm, RegisterForm, UserInfo, UserProfile } from '../types/auth'
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
      const res = await getMeDetail()
      userInfo.value = {
        id: res.id,
        username: res.username,
        nickname: res.nickname,
        avatar: res.avatar,
        role: res.role,
        email: res.email,
        sign: res.sign,
        coins: res.coins,
        created_at: res.created_at,
      }
      return res
    } catch {
      logout()
      throw new Error('获取用户信息失败')
    }
  }

  async function updateUserProfile(data: { sign?: string; nickname?: string }) {
    await updateProfile(data)
    // 更新本地状态
    if (userInfo.value) {
      if (data.sign !== undefined) userInfo.value.sign = data.sign
      if (data.nickname !== undefined) userInfo.value.nickname = data.nickname
    }
  }

  async function updateUserAvatar(file: File) {
    const res = await uploadAvatar(file)
    if (userInfo.value) {
      userInfo.value.avatar = res.avatar
    }
    return res.avatar
  }

  async function updateUserPassword(oldPwd: string, newPwd: string) {
    await updatePassword({ oldPassword: oldPwd, newPassword: newPwd, confirmPassword: newPwd })
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
    updateUserProfile,
    updateUserAvatar,
    updateUserPassword,
    logout,
    restore,
  }
})
