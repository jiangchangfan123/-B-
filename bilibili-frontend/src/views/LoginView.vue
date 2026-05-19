<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import type { LoginForm, RegisterForm } from '../types/auth'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const mode = ref<'login' | 'register'>('login')
const loading = ref(false)
const errorMsg = ref('')

const loginForm = reactive<LoginForm>({
  username: '',
  password: '',
})

const registerForm = reactive<RegisterForm>({
  username: '',
  password: '',
  email: '',
})

const activeForm = computed(() => (mode.value === 'login' ? loginForm : registerForm))

function switchMode(m: 'login' | 'register') {
  mode.value = m
  errorMsg.value = ''
}

function validate(): boolean {
  errorMsg.value = ''
  const f = activeForm.value
  if (!f.username) {
    errorMsg.value = '> ERROR: 请输入用户名'
    return false
  }
  if (f.username.length < 3 || f.username.length > 20) {
    errorMsg.value = '> ERROR: 用户名应为 3-20 位'
    return false
  }
  if (!f.password) {
    errorMsg.value = '> ERROR: 请输入密码'
    return false
  }
  if (f.password.length < 6 || f.password.length > 30) {
    errorMsg.value = '> ERROR: 密码应为 6-30 位'
    return false
  }
  if (mode.value === 'register') {
    const r = registerForm
    if (!r.email) {
      errorMsg.value = '> ERROR: 请输入邮箱'
      return false
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(r.email)) {
      errorMsg.value = '> ERROR: 邮箱格式不正确'
      return false
    }
  }
  return true
}

async function onSubmit() {
  if (!validate()) return
  loading.value = true
  try {
    if (mode.value === 'login') {
      await userStore.login(loginForm)
      const redirect = route.query.redirect as string
      router.push(redirect || '/')
    } else {
      await userStore.register(registerForm)
      errorMsg.value = ''
      // 注册成功，切换到登录并回填
      const savedUsername = registerForm.username
      mode.value = 'login'
      loginForm.username = savedUsername
      loginForm.password = ''
      errorMsg.value = '> OK: 注册成功，请登录'
    }
  } catch (e: any) {
    errorMsg.value = e.message || '> ERROR: 请求失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-card">
      <!-- Logo -->
      <div class="auth-logo">
        <span class="auth-logo__text">> NEBULA.TV</span>
        <div class="auth-logo__sub">系统接入 // v2.0</div>
      </div>

      <!-- 切换标签 -->
      <div class="auth-tabs">
        <button
          class="auth-tab"
          :class="{ 'auth-tab--active': mode === 'login' }"
          @click="switchMode('login')"
        >
          // 登录
        </button>
        <button
          class="auth-tab"
          :class="{ 'auth-tab--active': mode === 'register' }"
          @click="switchMode('register')"
        >
          // 注册
        </button>
      </div>

      <!-- 表单 -->
      <form class="auth-form" @submit.prevent="onSubmit">
        <div class="auth-field"
          >
          <label class="auth-field__label">> 用户名</label>
          <input
            v-model="activeForm.username"
            type="text"
            class="auth-field__input"
            placeholder="输入用户名..."
            :disabled="loading"
            maxlength="20"
          />
        </div>

        <div class="auth-field"
          >
          <label class="auth-field__label">> 密码</label>
          <input
            v-model="activeForm.password"
            type="password"
            class="auth-field__input"
            placeholder="6-30位密码..."
            :disabled="loading"
            maxlength="30"
          />
        </div>

        <div v-if="mode === 'register'" class="auth-field"
          >
          <label class="auth-field__label">> 邮箱</label>
          <input
            v-model="registerForm.email"
            type="text"
            class="auth-field__input"
            placeholder="输入邮箱..."
            :disabled="loading"
          />
        </div>

        <!-- 错误提示 -->
        <div v-if="errorMsg" class="auth-error">{{ errorMsg }}</div>

        <!-- 提交按钮 -->
        <button
          type="submit"
          class="auth-submit"
          :disabled="loading"
        >
          <span v-if="loading" class="auth-submit__loading"
            >加载中...</span>
          <span v-else>{{ mode === 'login' ? '[ 执行登录 ]' : '[ 执行注册 ]' }}</span>
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #08090d;
  background-image: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 240, 255, 0.02) 2px,
    rgba(0, 240, 255, 0.02) 4px
  );
  padding: 24px;
}

.auth-card {
  width: 420px;
  background: rgba(15, 17, 23, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 4px;
  padding: 40px 36px;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.auth-card:hover {
  border-color: rgba(0, 240, 255, 0.35);
  box-shadow: 0 0 24px rgba(0, 240, 255, 0.08);
}

/* Logo */
.auth-logo {
  text-align: center;
  margin-bottom: 32px;
}

.auth-logo__text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 22px;
  font-weight: 700;
  color: #00f0ff;
  letter-spacing: 0.08em;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4), 0 0 20px rgba(0, 240, 255, 0.2);
}

.auth-logo__text::after {
  content: '_';
  animation: blink 1s step-end infinite;
  color: #00f0ff;
}

@keyframes blink {
  50% { opacity: 0; }
}

.auth-logo__sub {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-top: 6px;
  letter-spacing: 0.1em;
}

/* Tabs */
.auth-tabs {
  display: flex;
  gap: 24px;
  margin-bottom: 28px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.auth-tab {
  background: none;
  border: none;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  letter-spacing: 0.06em;
  color: #5a5d6e;
  padding: 10px 0;
  cursor: pointer;
  position: relative;
  transition: color 0.1s;
}

.auth-tab:hover {
  color: #8b8fa3;
}

.auth-tab--active {
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
}

.auth-tab--active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 1px;
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.5);
}

/* Form Fields */
.auth-field {
  margin-bottom: 20px;
}

.auth-field__label {
  display: block;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.08em;
  margin-bottom: 8px;
}

.auth-field__input {
  width: 100%;
  height: 44px;
  padding: 0 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid #b829dd;
  border-radius: 4px;
  color: #e4e5eb;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.02em;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
}

.auth-field__input::placeholder {
  color: #5a5d6e;
  font-family: 'JetBrains Mono', Consolas, monospace;
}

.auth-field__input:focus {
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.2);
}

.auth-field__input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Error */
.auth-error {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #ff4d4f;
  letter-spacing: 0.02em;
  margin-bottom: 16px;
  min-height: 16px;
}

/* Submit Button */
.auth-submit {
  width: 100%;
  height: 48px;
  background: #b829dd;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  letter-spacing: 0.04em;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: background 0.15s, color 0.15s;
}

.auth-submit:hover:not(:disabled) {
  background: #00f0ff;
  color: #08090d;
}

.auth-submit:hover:not(:disabled)::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.15),
    transparent
  );
  animation: scanline 0.6s linear;
}

@keyframes scanline {
  to { left: 100%; }
}

.auth-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.auth-submit__loading {
  opacity: 0.7;
}
</style>
