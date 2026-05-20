<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import { getMyVideos, getHistory } from '../api/user'
import { deleteVideo } from '../api/video'
import type { MyVideoItem, HistoryItem, UpdatePasswordForm } from '../types/auth'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref<'settings' | 'archive' | 'history' | 'security'>('settings')
const loading = ref(false)
const message = ref('')
const messageType = ref<'ok' | 'error'>('ok')

// 编辑模式开关
const isEditing = ref(false)

// 头像 URL 拼接
const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const SERVER_BASE = API_BASE.replace(/\/api\/v1\/?$/, '')

function getFullAvatarUrl(avatar: string | undefined): string {
  if (!avatar) return ''
  if (avatar.startsWith('http')) return avatar
  return SERVER_BASE + avatar
}

function getFullUrl(url: string | undefined): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return SERVER_BASE + url
}

const avatarUrl = computed(() => getFullAvatarUrl(userStore.userInfo?.avatar))

// 昵称：有昵称显示昵称，没有显示用户名
const displayName = computed(() => {
  return userStore.userInfo?.nickname || userStore.userInfo?.username || ''
})

// ========== 头像上传 ==========
const fileInputRef = ref<HTMLInputElement | null>(null)
const avatarUploading = ref(false)

function onAvatarClick() {
  fileInputRef.value?.click()
}

async function onAvatarChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  if (file.size > 5 * 1024 * 1024) {
    showMessage('> ERROR: 文件超过 5MB 限制', 'error')
    return
  }
  avatarUploading.value = true
  try {
    await userStore.updateUserAvatar(file)
    showMessage('> OK: 头像上传成功', 'ok')
  } catch (err: any) {
    showMessage(err.message || '> ERROR: 上传失败', 'error')
  } finally {
    avatarUploading.value = false
    if (target) target.value = ''
  }
}

// ========== 资料设置 ==========
const profileForm = reactive({
  nickname: '',
  sign: '',
})

const signCount = computed(() => profileForm.sign.length)

// 进入编辑模式
function enterEditMode() {
  isEditing.value = true
  profileForm.nickname = userStore.userInfo?.nickname || ''
  profileForm.sign = userStore.userInfo?.sign || ''
}

// 取消编辑
function cancelEdit() {
  isEditing.value = false
  profileForm.nickname = userStore.userInfo?.nickname || ''
  profileForm.sign = userStore.userInfo?.sign || ''
}

async function saveProfile() {
  if (profileForm.sign.length > 200) {
    showMessage('> ERROR: 签名超过 200 字符', 'error')
    return
  }
  if (profileForm.nickname.length > 32) {
    showMessage('> ERROR: 昵称超过 32 字符', 'error')
    return
  }
  loading.value = true
  try {
    await userStore.updateUserProfile({
      sign: profileForm.sign,
      nickname: profileForm.nickname,
    })
    showMessage('> OK: 保存成功', 'ok')
    isEditing.value = false
  } catch (err: any) {
    showMessage(err.message || '> ERROR: 保存失败', 'error')
  } finally {
    loading.value = false
  }
}

// ========== 我的视频 ==========
const myVideos = ref<MyVideoItem[]>([])
const videoTotal = ref(0)

async function loadVideos() {
  try {
    const res = await getMyVideos()
    myVideos.value = res.list
    videoTotal.value = res.total
  } catch {
    myVideos.value = []
  }
}

async function onDeleteVideo(id: number) {
  if (!confirm('确定删险该视频吗？')) return
  try {
    await deleteVideo(id)
    myVideos.value = myVideos.value.filter(v => v.id !== id)
    videoTotal.value--
  } catch {
    alert('删险失败')
  }
}

function statusLabel(status: number) {
  if (status === 1) return '[已发布]'
  if (status === 2) return '[审核中]'
  if (status === 3) return '[已封禁]'
  return '[未知]'
}

function statusColor(status: number) {
  if (status === 1) return '#00f0ff'
  if (status === 2) return '#f0a500'
  if (status === 3) return '#ff4757'
  return '#5a5d6e'
}

// ========== 播放历史 ==========
const historyList = ref<HistoryItem[]>([])
const historyTotal = ref(0)
const showClearConfirm = ref(false)

async function loadHistory() {
  try {
    const res = await getHistory()
    historyList.value = res.list
    historyTotal.value = res.total
  } catch {
    historyList.value = []
  }
}

function formatTime(iso: string) {
  const d = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

function clearHistory() {
  showClearConfirm.value = false
  historyList.value = []
  historyTotal.value = 0
  showMessage('> OK: 历史已清空', 'ok')
}

// ========== 账号安全 ==========
const pwdForm = reactive<UpdatePasswordForm>({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
})

async function updatePwd() {
  if (!pwdForm.oldPassword || !pwdForm.newPassword || !pwdForm.confirmPassword) {
    showMessage('> ERROR: 请填写所有字段', 'error')
    return
  }
  if (pwdForm.newPassword.length < 6 || pwdForm.newPassword.length > 30) {
    showMessage('> ERROR: 新密码应为 6-30 位', 'error')
    return
  }
  if (pwdForm.newPassword !== pwdForm.confirmPassword) {
    showMessage('> ERROR: 两次输入的新密码不一致', 'error')
    return
  }
  loading.value = true
  try {
    await userStore.updateUserPassword(pwdForm.oldPassword, pwdForm.newPassword)
    showMessage('> OK: 密码修改成功，请重新登录', 'ok')
    setTimeout(() => {
      userStore.logout()
      router.push('/login')
    }, 1500)
  } catch (err: any) {
    showMessage(err.message || '> ERROR: 密码修改失败', 'error')
  } finally {
    loading.value = false
  }
}

// ========== 退出登录 ==========
function onLogout() {
  userStore.logout()
  router.push('/')
}

// ========== 通用 ==========
function showMessage(msg: string, type: 'ok' | 'error' = 'ok') {
  message.value = msg
  messageType.value = type
  setTimeout(() => { message.value = '' }, 3000)
}

function switchTab(tab: 'settings' | 'archive' | 'history' | 'security') {
  activeTab.value = tab
  message.value = ''
  isEditing.value = false
  if (tab === 'archive') loadVideos()
  if (tab === 'history') loadHistory()
}

onMounted(() => {
  if (!userStore.isLoggedIn) {
    router.push('/login?redirect=/personal')
    return
  }
})
</script>

<template>
  <div class="personal-page">
    <!-- 顶部导航 -->
    <NavBar />
    <SideBar />

    <main class="personal-main">
      <!-- 左侧个人信息卡片 -->
      <aside class="personal-sidebar">
        <div class="profile-card">
          <!-- 头像 -->
          <div class="avatar-wrap" @click="onAvatarClick">
            <div class="avatar-box" :class="{ 'avatar-uploading': avatarUploading }">
              <img
                v-if="avatarUrl"
                :src="avatarUrl"
                class="avatar-img"
                @error="userStore.userInfo && (userStore.userInfo.avatar = '')"
              />
              <div v-else class="avatar-default">
                {{ displayName.charAt(0).toUpperCase() }}
              </div>
              <div class="avatar-overlay">
                <span class="avatar-overlay-text"
                  >> 上传</span>
              </div>
            </div>
            <input
              ref="fileInputRef"
              type="file"
              accept="image/*"
              hidden
              @change="onAvatarChange"
            />
          </div>

          <!-- 昵称 -->
          <div class="profile-name">{{ displayName }}</div>
          <div class="profile-id">@{{ userStore.userInfo?.username }} \u00b7 ID: {{ userStore.userInfo?.id || 0 }}</div>

          <!-- 签名 -->
          <div class="profile-sign">
            {{ userStore.userInfo?.sign || '> 暂无签名...' }}
          </div>

          <!-- 数据栏 -->
          <div class="profile-stats">
            <div class="stat-item">
              <span class="stat-icon">▶</span>
              <span class="stat-num">12</span>
              <span class="stat-label">投稿</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">◆</span>
              <span class="stat-num">3.4K</span>
              <span class="stat-label">播放</span>
            </div>
            <div class="stat-item">
              <span class="stat-icon">■</span>
              <span class="stat-num">128</span>
              <span class="stat-label">粉丝</span>
            </div>
          </div>

          <!-- 上传头像按钮 -->
          <button class="profile-upload-btn" @click="onAvatarClick">
            [ 上传新头像 ]
          </button>

          <!-- 返回主页 + 退出登录 -->
          <div class="profile-actions">
            <button class="profile-home-btn" @click="router.push('/')">
              [ 返回主页 ]
            </button>
            <button class="profile-logout-btn" @click="onLogout">
              [ 退出登录 ]
            </button>
          </div>
        </div>
      </aside>

      <!-- 右侧内容区 -->
      <div class="personal-content">
        <!-- Tab 导航 -->
        <div class="tab-nav">
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'settings' }"
            @click="switchTab('settings')"
          >
            // 资料设置 [SETTINGS]
          </button>
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'archive' }"
            @click="switchTab('archive')"
          >
            // 我的视频 [ARCHIVE]
          </button>
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'history' }"
            @click="switchTab('history')"
          >
            // 播放历史 [HISTORY]
          </button>
          <button
            class="tab-btn"
            :class="{ 'tab-btn--active': activeTab === 'security' }"
            @click="switchTab('security')"
          >
            // 账号安全 [SECURITY]
          </button>
        </div>

        <!-- 消息提示 -->
        <div
          v-if="message"
          class="global-msg"
          :class="{ 'global-msg--error': messageType === 'error' }"
        >
          {{ message }}
        </div>

        <!-- Tab 1: 资料设置 -->
        <div v-if="activeTab === 'settings'" class="tab-panel">
          <!-- 浏览模式 -->
          <div v-if="!isEditing" class="profile-readonly">
            <div class="ro-row">
              <span class="ro-label">> 昵称</span>
              <span class="ro-value">{{ displayName }}</span>
            </div>
            <div class="ro-row">
              <span class="ro-label">> 用户名</span>
              <span class="ro-value">{{ userStore.userInfo?.username }}</span>
              <span class="ro-lock"
                >> LOCKED</span>
            </div>
            <div class="ro-row">
              <span class="ro-label">> 邮箱</span>
              <span class="ro-value">{{ userStore.userInfo?.email || '> 未设置' }}</span>
              <span class="ro-lock"
                >> LOCKED</span>
            </div>
            <div class="ro-row">
              <span class="ro-label">> 个性签名</span>
              <span class="ro-value ro-value--sign"
                >{{ userStore.userInfo?.sign || '> 暂无签名...' }}</span>
            </div>
            <button class="submit-btn" style="margin-top: 24px;" @click="enterEditMode">
              [ 编辑资料 ]
            </button>
          </div>

          <!-- 编辑模式 -->
          <div v-else class="profile-editing">
            <div class="form-group">
              <label class="form-label">> 昵称</label>
              <input
                v-model="profileForm.nickname"
                type="text"
                class="form-input"
                placeholder="输入昵称..."
                maxlength="32"
              />
            </div>

            <div class="form-group">
              <label class="form-label">> 用户名</label>
              <div class="form-locked">
                <input
                  :value="userStore.userInfo?.username"
                  type="text"
                  class="form-input form-input--locked"
                  disabled
                />
                <span class="lock-tag"
                  >> LOCKED</span>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">> 邮箱</label>
              <div class="form-locked">
                <input
                  :value="userStore.userInfo?.email || ''"
                  type="text"
                  class="form-input form-input--locked"
                  disabled
                />
                <span class="lock-tag"
                  >> LOCKED</span>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">> 个性签名</label>
              <textarea
                v-model="profileForm.sign"
                class="form-textarea"
                placeholder="写点什么吧..."
                maxlength="200"
              ></textarea>
              <div class="char-count" :class="{ 'char-count--over': signCount > 200 }">
                {{ signCount }}/200
              </div>
            </div>

            <div class="edit-actions">
              <button class="submit-btn" :disabled="loading" @click="saveProfile">
                <span v-if="loading"
                  >> SYNCING...</span>
                <span v-else>[ 保存更改 ]</span>
              </button>
              <button class="cancel-btn" @click="cancelEdit">
                [ 取消 ]
              </button>
            </div>
          </div>
        </div>

        <!-- Tab 2: 我的视频 -->
        <div v-if="activeTab === 'archive'" class="tab-panel">
          <div class="archive-stats">
            <span class="archive-stat">TOTAL_UPLOADS: {{ myVideos.length }}</span>
            <span class="archive-stat">TOTAL_VIEWS: 45.2K</span>
            <span class="archive-stat">PENDING: 2</span>
          </div>

          <div v-if="myVideos.length > 0" class="video-list">
            <div
              v-for="video in myVideos"
              :key="video.id"
              class="video-card"
              @click="router.push(`/video/${video.id}`)"
            >
              <div class="video-cover">
                <img
                  v-if="video.cover_url"
                  :src="getFullUrl(video.cover_url)"
                  class="video-cover-img"
                  alt="cover"
                  @error="$event.target.style.display='none'"
                />
                <div v-else class="video-cover-placeholder">{{ video.category }}</div>
              </div>
              <div class="video-info">
                <div class="video-title">{{ video.title }}</div>
                <div class="video-meta">
                  <span
                    class="video-status"
                    :style="{ borderColor: statusColor(video.status), color: statusColor(video.status) }"
                  >
                    {{ statusLabel(video.status) }}
                  </span>
                  <span class="video-views">{{ video.views }} 播放</span>
                </div>
              </div>
              <div class="video-actions">
                <button class="action-btn" @click.stop="router.push(`/upload/${video.id}`)">[ 编辑 ]</button>
                <button class="action-btn" @click.stop="onDeleteVideo(video.id)">[ 删除 ]</button>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <div class="empty-text"
              >> NO_DATA_FOUND</div>
            <button class="submit-btn" style="margin-top: 16px; width: auto; padding: 0 24px;" @click="router.push('/upload')">
              [ 立即上传 ]
            </button>
          </div>
        </div>

        <!-- Tab 3: 播放历史 -->
        <div v-if="activeTab === 'history'" class="tab-panel">
          <div v-if="historyList.length > 0" class="history-list">
            <div
              v-for="item in historyList"
              :key="item.id"
              class="history-item"
              @click="router.push(`/video/${item.id}`)"
            >
              <div class="history-cover">
                <img
                  v-if="item.cover_url"
                  :src="getFullUrl(item.cover_url)"
                  class="history-cover-img"
                  alt="cover"
                  @error="$event.target.style.display='none'"
                />
                <div v-else class="history-cover-placeholder">▶</div>
              </div>
              <div class="history-info">
                <div class="history-title">{{ item.title }}</div>
                <div class="history-time"
                  >> {{ formatTime(item.watched_at) }}</div>
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            <div class="empty-text"
              >> NO_HISTORY_FOUND</div>
          </div>

          <button
            v-if="historyList.length > 0"
            class="clear-btn"
            @click="showClearConfirm = true"
          >
            [ 清空所有历史 ]
          </button>

          <!-- 清空确认 -->
          <div v-if="showClearConfirm" class="confirm-dialog">
            <div class="confirm-box">
              <div class="confirm-text"
                >CONFIRM_DELETE? [Y/N]</div>
              <div class="confirm-actions">
                <button class="confirm-btn confirm-btn--yes" @click="clearHistory"
                  >[Y]</button>
                <button class="confirm-btn confirm-btn--no" @click="showClearConfirm = false"
                  >[N]</button>
              </div>
            </div>
          </div>
        </div>

        <!-- Tab 4: 账号安全 -->
        <div v-if="activeTab === 'security'" class="tab-panel">
          <div class="form-group">
            <label class="form-label">> 旧密码</label>
            <input
              v-model="pwdForm.oldPassword"
              type="password"
              class="form-input"
              placeholder="输入旧密码..."
            />
          </div>

          <div class="form-group">
            <label class="form-label">> 新密码</label>
            <input
              v-model="pwdForm.newPassword"
              type="password"
              class="form-input"
              placeholder="6-30位新密码..."
            />
          </div>

          <div class="form-group">
            <label class="form-label">> 确认新密码</label>
            <input
              v-model="pwdForm.confirmPassword"
              type="password"
              class="form-input"
              placeholder="再次输入新密码..."
            />
          </div>

          <button class="submit-btn" :disabled="loading" @click="updatePwd">
            <span v-if="loading"
              >> SYNCING...</span>
            <span v-else>[ 更新密码 ]</span>
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.personal-page {
  min-height: 100vh;
  background: #08090d;
  background-image: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 240, 255, 0.02) 2px,
    rgba(0, 240, 255, 0.02) 4px
  );
}

.personal-main {
  margin-left: 200px;
  margin-top: 64px;
  min-height: calc(100vh - 64px);
  padding: 24px;
  display: flex;
  gap: 24px;
}

/* 左侧边栏 */
.personal-sidebar {
  width: 280px;
  flex-shrink: 0;
}

.profile-card {
  background: #0f1117;
  border: 1px solid rgba(0, 240, 255, 0.1);
  border-radius: 4px;
  padding: 28px 24px;
}

/* 头像 */
.avatar-wrap {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
  cursor: pointer;
}

.avatar-box {
  width: 120px;
  height: 120px;
  border: 2px solid #00f0ff;
  position: relative;
  overflow: hidden;
  transition: box-shadow 0.15s;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.2);
}

.avatar-box:hover {
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.4);
}

.avatar-uploading::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 240, 255, 0.3),
    transparent
  );
  animation: scanline 0.6s linear infinite;
}

@keyframes scanline {
  to { left: 100%; }
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-default {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 48px;
  font-weight: 700;
  color: #00f0ff;
  background: linear-gradient(135deg, #0a0a0f, #1a1c24);
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.1s;
}

.avatar-box:hover .avatar-overlay {
  opacity: 1;
}

.avatar-overlay-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  letter-spacing: 0.06em;
}

/* 用户信息 */
.profile-name {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 20px;
  font-weight: 600;
  color: #e4e5eb;
  text-align: center;
  letter-spacing: 0.02em;
}

.profile-id {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  text-align: center;
  margin-top: 4px;
  letter-spacing: 0.04em;
}

.profile-sign {
  font-size: 14px;
  color: #8b8fa3;
  text-align: center;
  margin-top: 12px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  min-height: 42px;
}

/* 数据栏 */
.profile-stats {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #8b8fa3;
}

.stat-icon {
  color: #00f0ff;
  font-size: 10px;
}

.stat-num {
  color: #e4e5eb;
  font-weight: 600;
}

.stat-label {
  color: #5a5d6e;
}

/* 上传头像按钮 */
.profile-upload-btn {
  width: 100%;
  height: 40px;
  margin-top: 20px;
  background: #b829dd;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.04em;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: background 0.15s, color 0.15s;
}

.profile-upload-btn:hover {
  background: #00f0ff;
  color: #08090d;
}

.profile-upload-btn:hover::after {
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

/* 底部操作区 */
.profile-actions {
  display: flex;
  gap: 10px;
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.profile-home-btn,
.profile-logout-btn {
  flex: 1;
  height: 36px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition: all 0.15s;
}

.profile-home-btn {
  background: transparent;
  color: #8b8fa3;
}

.profile-home-btn:hover {
  border-color: #00f0ff;
  color: #00f0ff;
  background: rgba(0, 240, 255, 0.05);
}

.profile-logout-btn {
  background: transparent;
  color: #ff4757;
  border-color: rgba(255, 71, 87, 0.2);
}

.profile-logout-btn:hover {
  border-color: #ff4757;
  background: rgba(255, 71, 87, 0.08);
}

/* 右侧内容 */
.personal-content {
  flex: 1;
  min-width: 0;
}

/* Tab 导航 */
.tab-nav {
  display: flex;
  gap: 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 24px;
}

.tab-btn {
  background: none;
  border: none;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  letter-spacing: 0.04em;
  color: #5a5d6e;
  padding: 12px 0;
  cursor: pointer;
  position: relative;
  transition: color 0.1s;
  white-space: nowrap;
}

.tab-btn:hover {
  color: #8b8fa3;
}

.tab-btn--active {
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
}

.tab-btn--active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 1px;
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.5);
}

/* 全局消息 */
.global-msg {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  margin-bottom: 16px;
  padding: 8px 12px;
  background: rgba(0, 240, 255, 0.05);
  border: 1px solid rgba(0, 240, 255, 0.15);
  border-radius: 4px;
}

.global-msg--error {
  color: #ff4d4f;
  background: rgba(255, 77, 79, 0.05);
  border-color: rgba(255, 77, 79, 0.15);
}

/* 表单 */
.tab-panel {
  background: #0f1117;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
  padding: 32px;
}

/* 浏览模式 */
.profile-readonly {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.ro-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
}

.ro-label {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.08em;
  width: 80px;
  flex-shrink: 0;
}

.ro-value {
  flex: 1;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #e4e5eb;
}

.ro-value--sign {
  color: #8b8fa3;
  line-height: 1.5;
}

.ro-lock {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}

/* 编辑模式 */
.profile-editing {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.edit-actions {
  display: flex;
  gap: 16px;
}

.edit-actions .submit-btn {
  flex: 1;
}

.form-group {
  margin-bottom: 4px;
}

.form-label {
  display: block;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.08em;
  margin-bottom: 8px;
}

.form-input,
.form-textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid #b829dd;
  border-radius: 4px;
  color: #e4e5eb;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.02em;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  padding: 0 16px;
}

.form-input {
  height: 44px;
}

.form-textarea {
  height: 100px;
  padding: 12px 16px;
  resize: none;
}

.form-input:focus,
.form-textarea:focus {
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.2);
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: #5a5d6e;
}

.form-input--locked {
  background: rgba(255, 255, 255, 0.03);
  border-color: rgba(255, 255, 255, 0.08);
  color: #5a5d6e;
  cursor: not-allowed;
}

.form-locked {
  position: relative;
}

.lock-tag {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}

.char-count {
  text-align: right;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-top: 6px;
}

.char-count--over {
  color: #ff4d4f;
}

.submit-btn {
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

.submit-btn:hover:not(:disabled) {
  background: #00f0ff;
  color: #08090d;
}

.submit-btn:hover:not(:disabled)::after {
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

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.cancel-btn {
  width: 120px;
  height: 48px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 4px;
  color: #8b8fa3;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.cancel-btn:hover {
  border-color: #00f0ff;
  color: #00f0ff;
}

/* 视频列表 */
.archive-stats {
  display: flex;
  gap: 24px;
  margin-bottom: 20px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}

.video-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.video-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  transition: border-color 0.15s;
}

.video-card:hover {
  border-color: rgba(0, 240, 255, 0.2);
}

.video-cover {
  width: 120px;
  height: 72px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #14161f, #2a2d3a);
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-cover-placeholder {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
}

.video-info {
  flex: 1;
  min-width: 0;
}

.video-title {
  font-size: 14px;
  font-weight: 500;
  color: #e4e5eb;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.video-meta {
  display: flex;
  align-items: center;
  gap: 12px;
}

.video-status {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  padding: 2px 6px;
  border: 1px solid;
  border-radius: 2px;
}

.video-views {
  font-size: 12px;
  color: #5a5d6e;
}

.video-actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.action-btn {
  background: none;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: #8b8fa3;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  padding: 6px 12px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.action-btn:hover {
  border-color: #00f0ff;
  color: #00f0ff;
}

/* 播放历史 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  cursor: pointer;
  transition: border-color 0.15s;
}

.history-item:hover {
  border-color: rgba(0, 240, 255, 0.2);
}

.history-item:hover .history-title {
  color: #00f0ff;
}

.history-cover {
  width: 60px;
  height: 40px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #14161f, #2a2d3a);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.video-cover {
  width: 120px;
  height: 72px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #14161f, #2a2d3a);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}

.video-cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.history-cover-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 2px;
}

.history-cover-placeholder {
  font-size: 14px;
  color: #5a5d6e;
}

.history-info {
  flex: 1;
  min-width: 0;
}

.history-title {
  font-size: 14px;
  color: #e4e5eb;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color 0.15s;
}

.history-time {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
}

.clear-btn {
  margin-top: 20px;
  width: 100%;
  height: 40px;
  background: transparent;
  border: 1px solid #ff4757;
  border-radius: 4px;
  color: #ff4757;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition: background 0.15s;
}

.clear-btn:hover {
  background: rgba(255, 71, 87, 0.1);
}

/* 确认对话框 */
.confirm-dialog {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(8, 9, 13, 0.8);
  z-index: 2000;
}

.confirm-box {
  background: #0f1117;
  border: 1px solid rgba(255, 71, 87, 0.3);
  border-radius: 4px;
  padding: 32px 40px;
  text-align: center;
}

.confirm-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #ff4757;
  margin-bottom: 20px;
  letter-spacing: 0.04em;
}

.confirm-actions {
  display: flex;
  gap: 16px;
  justify-content: center;
}

.confirm-btn {
  width: 60px;
  height: 36px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  cursor: pointer;
  transition: border-color 0.15s, color 0.15s;
}

.confirm-btn--yes {
  color: #ff4757;
  border-color: rgba(255, 71, 87, 0.3);
}

.confirm-btn--yes:hover {
  border-color: #ff4757;
  background: rgba(255, 71, 87, 0.1);
}

.confirm-btn--no {
  color: #8b8fa3;
}

.confirm-btn--no:hover {
  border-color: #00f0ff;
  color: #00f0ff;
}

/* 空状态 */
.empty-state {
  padding: 60px 0;
  text-align: center;
}

.empty-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}
</style>
