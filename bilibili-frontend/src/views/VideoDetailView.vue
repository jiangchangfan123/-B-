<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import { getVideoDetail, getTranscodeStatus, getVideoList } from '../api/video'
import type { VideoDetail, VideoListItem } from '../types/video'

const route = useRoute()
const router = useRouter()
const videoId = Number(route.params.id)

const video = ref<VideoDetail | null>(null)
const loading = ref(true)
const relatedVideos = ref<VideoListItem[]>([])
const isPlaying = ref(false)
const descExpanded = ref(false)

// 转码轮询
let pollTimer: ReturnType<typeof setInterval> | null = null
const transcodePolling = ref(false)

onMounted(async () => {
  await loadVideo()
  await loadRelated()
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

async function loadVideo() {
  loading.value = true
  try {
    const res = await getVideoDetail(videoId)
    video.value = res
    // 如果转码未完成，开始轮询
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

function formatTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

function onPlay() {
  isPlaying.value = true
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

            <div v-else class="player-video">
              <video
                :src="getFullUrl(video.transcoded_url)"
                controls
                class="video-element"
                @play="onPlay"
                :poster="getFullUrl(video.cover_url)"
              ></video>
            </div>
          </div>

          <!-- 视频信息 -->
          <div class="video-info">
            <h1 class="video-title"
              >{{ video.title }}</h1>

            <div class="video-meta-row">
              <span class="video-category"
                >[{{ video.category.toUpperCase() }}]</span>
              <div class="video-stats">
                <span class="stat-item"
                  >\u25b6 {{ formatViews(video.view_count) }}</span>
                <span class="stat-item"
                  >\u25c6 {{ video.like_count }}</span>
                <span class="stat-item"
                  >\u25a0 {{ video.comment_count }}</span>
              </div>
            </div>

            <div class="video-up">
              <div class="up-avatar">
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
</style>
