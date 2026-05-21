<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import VideoPlayer from '../components/VideoPlayer.vue'
import DanmakuLayer from '../components/DanmakuLayer.vue'
import { getVideoDetail, getTranscodeStatus, getVideoList, likeVideo, getLikeStatus, favoriteVideo, getFavoriteStatus } from '../api/video'
import { sendDanmaku } from '../api/danmaku'
import CommentSection from '../components/CommentSection.vue'
import type { VideoDetail, VideoListItem } from '../types/video'

const route = useRoute()
const router = useRouter()
const videoId = Number(route.params.id)

const video = ref<VideoDetail | null>(null)
const loading = ref(true)
const relatedVideos = ref<VideoListItem[]>([])
const descExpanded = ref(false)
const currentTime = ref(0)
const isPlaying = ref(false)

// 弹幕输入
const danmakuInput = ref('')
const isSendingDanmaku = ref(false)
const danmakuLayerRef = ref<InstanceType<typeof DanmakuLayer> | null>(null)

// 返回主页
function goBack() {
  router.push('/')
}

// 互动按钮骨架状态
const actionState = ref({
  liked: false,
  coined: false,
  favorited: false,
})

async function onLike() {
  try {
    const res = await likeVideo(videoId)
    actionState.value.liked = res.liked
    if (video.value) {
      video.value.like_count = res.count
    }
  } catch {
    showToast('点赞失败，请重试')
  }
}
function onCoin() {
  actionState.value.coined = !actionState.value.coined
}
async function onFavorite() {
  try {
    const res = await favoriteVideo(videoId)
    actionState.value.favorited = res.favorited
  } catch {
    showToast('收藏失败，请重试')
  }
}
function onShare() {
  const url = window.location.href
  navigator.clipboard.writeText(url).catch(() => {})
  // 简单 toast 提示
  showToast('链接已复制到剪贴板')
}

async function handleSendDanmaku() {
  const text = danmakuInput.value.trim()
  if (!text) return
  if (isSendingDanmaku.value) return

  // 检查登录状态
  const userStore = useUserStore()
  if (!userStore.isLoggedIn) {
    showToast('请先登录后再发送弹幕')
    setTimeout(() => router.push('/login'), 1500)
    return
  }

  isSendingDanmaku.value = true
  try {
    const timePoint = Math.floor(currentTime.value)
    const res = await sendDanmaku(videoId, {
      content: text,
      time_point: timePoint,
    })

    // 通知弹幕层立即显示
    danmakuLayerRef.value?.addDanmaku(res.danmaku)
    danmakuInput.value = ''
    showToast('弹幕发送成功')
  } catch (err: any) {
    console.error('[Danmaku] 发送失败:', err)
    const msg = err?.message || ''
    if (msg.includes('登录') || msg.includes('Token') || msg.includes('token')) {
      showToast('登录已过期，请重新登录')
      setTimeout(() => router.push('/login'), 1500)
    } else {
      showToast('弹幕发送失败: ' + (msg || '请检查网络'))
    }
  } finally {
    isSendingDanmaku.value = false
  }
}

const toastMsg = ref('')
let toastTimer: ReturnType<typeof setTimeout> | null = null
function showToast(msg: string) {
  toastMsg.value = msg
  if (toastTimer) clearTimeout(toastTimer)
  toastTimer = setTimeout(() => { toastMsg.value = '' }, 2000)
}

// 转码轮询
let pollTimer: ReturnType<typeof setInterval> | null = null
const transcodePolling = ref(false)

onMounted(async () => {
  await loadVideo()
  await loadRelated()
  // 查询点赞状态
  try {
    const status = await getLikeStatus(videoId)
    actionState.value.liked = status.liked
    if (video.value) {
      video.value.like_count = status.count
    }
  } catch {
    // 未登录或其他错误，忽略
  }
  // 查询收藏状态
  try {
    const favStatus = await getFavoriteStatus(videoId)
    actionState.value.favorited = favStatus.favorited
  } catch {
    // 未登录或其他错误，忽略
  }
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
  if (toastTimer) clearTimeout(toastTimer)
})

async function loadVideo() {
  loading.value = true
  try {
    const res = await getVideoDetail(videoId)
    video.value = res
    if (res.transcode_status !== 2) {
      startTranscodePolling()
    }
  } catch {
    video.value = null
  } finally {
    loading.value = false
  }
}

function startTranscodePolling() {
  transcodePolling.value = true
  pollTimer = setInterval(async () => {
    try {
      const status = await getTranscodeStatus(videoId)
      if (status.transcode_status === 2 && video.value) {
        video.value.transcode_status = 2
        video.value.transcoded_url = status.transcoded_url
        video.value.status = 1
        if (pollTimer) clearInterval(pollTimer)
        transcodePolling.value = false
      } else if (status.transcode_status === 3) {
        if (pollTimer) clearInterval(pollTimer)
        transcodePolling.value = false
      }
    } catch {
      // 继续轮询
    }
  }, 3000)
}

async function loadRelated() {
  try {
    const res = await getVideoList(1, 8)
    relatedVideos.value = res.list.filter(v => v.id !== videoId).slice(0, 6)
  } catch {
    relatedVideos.value = []
  }
}

function formatViews(n: number): string {
  if (n >= 10000) return (n / 10000).toFixed(1) + 'W'
  if (n >= 1000) return (n / 1000).toFixed(1) + 'K'
  return String(n)
}

function formatDuration(seconds?: number): string {
  if (!seconds || seconds <= 0) return '00:00'
  const m = Math.floor(seconds / 60).toString().padStart(2, '0')
  const s = (seconds % 60).toString().padStart(2, '0')
  return `${m}:${s}`
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const SERVER_BASE = API_BASE.replace(/\/api\/v1\/?$/, '')

function getFullUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return SERVER_BASE + url
}
</script>

<template>
  <div class="video-page">
    <NavBar />
    <SideBar />

    <main class="video-main">
      <!-- 返回主页按钮 -->
      <div class="back-bar">
        <button class="back-btn" @click="goBack">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="15 18 9 12 15 6" />
          </svg>
          <span>返回主页</span>
        </button>
      </div>

      <!-- Toast 提示 -->
      <div v-if="toastMsg" class="video-toast">{{ toastMsg }}</div>

      <div v-if="loading" class="video-loading"
        >> LOADING_DATA...</div>

      <div v-else-if="!video" class="video-empty"
        >> VIDEO_NOT_FOUND</div>

      <div v-else class="video-layout">
        <!-- 左侧：播放器 + 信息 -->
        <div class="video-left">
          <!-- 播放器 -->
          <div class="player-wrap">
            <div v-if="video.transcode_status !== 2" class="player-placeholder">
              <img v-if="video.cover_url" :src="getFullUrl(video.cover_url)" class="player-cover" />
              <div class="player-overlay">
                <div class="transcode-status"
                  >> TRANSCODING_IN_PROGRESS...</div>
                <div class="transcode-detail"
                  >> 480P_ENCODING // ETA: UNKNOWN</div>
                <div v-if="transcodePolling" class="transcode-pulse"
                  ></div>
              </div>
            </div>

            <template v-else>
              <DanmakuLayer
                ref="danmakuLayerRef"
                :video-id="video.id"
                :current-time="currentTime"
                :is-playing="isPlaying"
              />
              <VideoPlayer
                :video-url="getFullUrl(video.transcoded_url || video.video_url)"
                :cover-url="getFullUrl(video.cover_url)"
                :video-id="video.id"
                @timeupdate="currentTime = $event"
                @playstate="isPlaying = $event"
              />
            </template>
          </div>

          <!-- 弹幕输入框 -->
          <div class="danmaku-send-bar">
            <input
              v-model="danmakuInput"
              type="text"
              class="danmaku-send-input"
              placeholder="发一条友好的弹幕吧~"
              maxlength="100"
              @keyup.enter="handleSendDanmaku"
            />
            <button
              class="danmaku-send-action"
              :disabled="isSendingDanmaku || !danmakuInput.trim()"
              @click="handleSendDanmaku"
            >
              {{ isSendingDanmaku ? '发送中...' : '发射弹幕' }}
            </button>
          </div>

          <!-- 视频信息 -->
          <div class="video-info">
            <!-- 互动按钮区 -->
            <div class="action-bar">
              <button
                class="action-btn"
                :class="{ 'action-btn--active': actionState.liked }"
                @click="onLike"
              >
                <span class="action-icon">▲</span>
                <span class="action-label">点赞</span>
                <span class="action-count">{{ formatViews(video.like_count) }}</span>
              </button>
              <button
                class="action-btn"
                :class="{ 'action-btn--active': actionState.coined }"
                @click="onCoin"
              >
                <span class="action-icon">◆</span>
                <span class="action-label">投币</span>
              </button>
              <button
                class="action-btn"
                :class="{ 'action-btn--active': actionState.favorited }"
                @click="onFavorite"
              >
                <span class="action-icon">★</span>
                <span class="action-label">收藏</span>
              </button>
              <button class="action-btn" @click="onShare">
                <span class="action-icon">⇧</span>
                <span class="action-label">分享</span>
              </button>
            </div>

            <h1 class="video-title"
              >{{ video.title }}</h1>

            <div class="video-meta-row">
              <span class="video-category"
                >[{{ video.category.toUpperCase() }}]</span>
              <div class="video-stats">
                <span class="stat-item"
                  >\u25b6 {{ formatViews(video.view_count) }}</span>
                <span class="stat-item"
                  >\u23f1 {{ formatDuration(video.duration) }}</span>
                <span class="stat-item"
                  >\u25c6 {{ video.like_count }}</span>
                <span class="stat-item"
                  >\u25a0 {{ video.comment_count }}</span>
              </div>
            </div>

            <div class="video-up">
              <img
                v-if="video.user_info?.avatar"
                :src="getFullUrl(video.user_info.avatar)"
                class="up-avatar up-avatar--img"
                alt="avatar"
                @error="$event.target.style.display='none'"
              />
              <div v-else class="up-avatar">
                {{ (video.user_info?.nickname || video.user_info?.username || '?').charAt(0).toUpperCase() }}
              </div>
              <div class="up-info">
                <div class="up-name"
                  >{{ video.user_info?.nickname || video.user_info?.username }}</div>
                <div class="up-time"
                  >> {{ formatTime(video.created_at) }}</div>
              </div>
            </div>

            <div class="video-desc">
              <p
                class="desc-text"
                :class="{ 'desc-text--expanded': descExpanded }"
              >
                {{ video.description || '> 暂无简介...' }}
              </p>
              <button
                v-if="video.description && video.description.length > 80"
                class="desc-toggle"
                @click="descExpanded = !descExpanded"
              >
                {{ descExpanded ? '> COLLAPSE' : '> EXPAND' }}
              </button>
            </div>
          </div>

          <!-- 评论区 -->
          <CommentSection
            :video-id="videoId"
            @update-count="(n) => { if (video) video.comment_count = n }"
          />
        </div>

        <!-- 右侧：推荐列表 -->
        <aside class="video-right">
          <div class="related-header"
            >// RELATED_ARCHIVE</div>
          <div class="related-list">
            <div
              v-for="item in relatedVideos"
              :key="item.id"
              class="related-card"
              @click="router.push(`/video/${item.id}`)"
            >
              <div class="related-thumb">
                <div class="related-placeholder"
                  >{{ item.category.toUpperCase() }}</div>
              </div>
              <div class="related-info">
                <div class="related-title"
                  >{{ item.title }}</div>
                <div class="related-meta"
                  >\u25b6 {{ formatViews(item.view_count) }} \u00b7 {{ item.user?.nickname || item.user?.username }}</div>
              </div>
            </div>
          </div>
        </aside>
      </div>
    </main>
  </div>
</template>

<style scoped>
.video-page {
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

.video-main {
  margin-left: 200px;
  margin-top: 64px;
  min-height: calc(100vh - 64px);
  padding: 24px;
}

.video-loading,
.video-empty {
  padding: 120px 0;
  text-align: center;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}

/* 布局 */
.video-layout {
  display: flex;
  gap: 24px;
}

.video-left {
  flex: 1;
  min-width: 0;
}

.video-right {
  width: 320px;
  flex-shrink: 0;
}

/* 播放器 */
.player-wrap {
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #0a0a0f;
  border: 1px solid rgba(0, 240, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
  position: relative;
}

.player-placeholder {
  width: 100%;
  height: 100%;
  position: relative;
}

.player-cover {
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0.5;
}

.player-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.transcode-status {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 16px;
  color: #00f0ff;
  text-shadow: 0 0 12px rgba(0, 240, 255, 0.4);
  letter-spacing: 0.04em;
}

.transcode-detail {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}

.transcode-pulse {
  width: 40px;
  height: 4px;
  background: #00f0ff;
  border-radius: 2px;
  margin-top: 8px;
  animation: pulse-bar 1.5s ease-in-out infinite;
}

@keyframes pulse-bar {
  0%, 100% { opacity: 0.3; width: 20px; }
  50% { opacity: 1; width: 60px; }
}

.player-video {
  width: 100%;
  height: 100%;
}

.video-element {
  width: 100%;
  height: 100%;
  object-fit: contain;
  background: #0a0a0f;
}

/* 弹幕发送栏 */
.danmaku-send-bar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
  padding: 8px 12px;
  background: #0f1117;
  border: 1px solid rgba(0, 240, 255, 0.1);
  border-radius: 4px;
}

.danmaku-send-input {
  flex: 1;
  height: 34px;
  padding: 0 12px;
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: 4px;
  background: #0a0a0f;
  color: #e4e5eb;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}

.danmaku-send-input::placeholder {
  color: #5a5d6e;
}

.danmaku-send-input:focus {
  border-color: rgba(0, 240, 255, 0.5);
}

.danmaku-send-action {
  height: 34px;
  padding: 0 18px;
  border: none;
  border-radius: 4px;
  background: rgba(0, 240, 255, 0.15);
  color: #00f0ff;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  white-space: nowrap;
}

.danmaku-send-action:hover:not(:disabled) {
  background: rgba(0, 240, 255, 0.25);
}

.danmaku-send-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 视频信息 */
.video-info {
  margin-top: 20px;
  padding: 20px;
  background: #0f1117;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
}

.video-title {
  font-size: 18px;
  font-weight: 600;
  color: #e4e5eb;
  margin: 0 0 12px;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.video-meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.video-category {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  padding: 2px 8px;
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 2px;
}

.video-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  color: #8b8fa3;
}

/* UP 主 */
.video-up {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 0;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 12px;
}

.up-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #0a0a0f, #1a1c24);
  border: 1px solid rgba(0, 240, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  font-weight: 600;
  color: #00f0ff;
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
}

.up-avatar--img {
  object-fit: cover;
  background: none;
}

.up-info {
  flex: 1;
}

.up-name {
  font-size: 14px;
  font-weight: 500;
  color: #e4e5eb;
}

.up-time {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-top: 2px;
}

/* 简介 */
.video-desc {
  margin-top: 12px;
}

.desc-text {
  font-size: 14px;
  color: #8b8fa3;
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  white-space: pre-line;
}

.desc-text--expanded {
  -webkit-line-clamp: unset;
}

.desc-toggle {
  background: none;
  border: none;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #00f0ff;
  cursor: pointer;
  margin-top: 8px;
  padding: 0;
}

/* 右侧推荐 */
.related-header {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  font-weight: 600;
  color: #b829dd;
  letter-spacing: 0.08em;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(184, 41, 221, 0.2);
  margin-bottom: 16px;
}

.related-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.related-card {
  display: flex;
  gap: 12px;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: background 0.1s;
}

.related-card:hover {
  background: rgba(0, 240, 255, 0.03);
}

.related-card:hover .related-title {
  color: #00f0ff;
}

.related-thumb {
  width: 120px;
  height: 72px;
  flex-shrink: 0;
  background: linear-gradient(135deg, #14161f, #2a2d3a);
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.related-placeholder {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #5a5d6e;
}

.related-info {
  flex: 1;
  min-width: 0;
}

.related-title {
  font-size: 13px;
  color: #e4e5eb;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 4px;
  transition: color 0.15s;
}

.related-meta {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
}

/* ========== 返回按钮 ========== */
.back-bar {
  margin-bottom: 16px;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: 4px;
  padding: 6px 14px;
  color: #8b8fa3;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}

.back-btn:hover {
  color: #00f0ff;
  border-color: rgba(0, 240, 255, 0.5);
  background: rgba(0, 240, 255, 0.05);
}

/* ========== Toast 提示 ========== */
.video-toast {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  background: rgba(15, 17, 23, 0.95);
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 4px;
  padding: 10px 24px;
  color: #00f0ff;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  z-index: 2000;
  animation: toast-in 0.2s ease;
}

@keyframes toast-in {
  from { opacity: 0; transform: translateX(-50%) translateY(-8px); }
  to   { opacity: 1; transform: translateX(-50%) translateY(0); }
}

/* ========== 互动按钮区 ========== */
.action-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  padding: 8px 16px;
  color: #8b8fa3;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s;
}

.action-btn:hover {
  background: rgba(0, 240, 255, 0.06);
  border-color: rgba(0, 240, 255, 0.3);
  color: #e4e5eb;
}

.action-btn--active {
  background: rgba(0, 240, 255, 0.1);
  border-color: rgba(0, 240, 255, 0.5);
  color: #00f0ff;
}

.action-btn--active .action-icon {
  color: #00f0ff;
  filter: drop-shadow(0 0 4px rgba(0, 240, 255, 0.5));
}

.action-icon {
  font-size: 12px;
  color: #5a5d6e;
  transition: color 0.15s;
}

.action-label {
  letter-spacing: 0.02em;
}

.action-count {
  margin-left: 2px;
  color: #5a5d6e;
  font-size: 12px;
}

.action-btn--active .action-count {
  color: #00f0ff;
}
</style>
