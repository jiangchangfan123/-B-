<script setup lang="ts">
import { ref, onMounted } from 'vue'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import BannerCarousel from '../components/BannerCarousel.vue'
import VideoGrid from '../components/VideoGrid.vue'
import SectionBlock from '../components/SectionBlock.vue'
import BackToTop from '../components/BackToTop.vue'
import ScrollProgress from '../components/ScrollProgress.vue'
import { getVideoList } from '../api/video'
import { generateRanks } from '../mock/homeData'
import type { VideoItem } from '../types/home'

const videoList = ref<VideoItem[]>([])
const loading = ref(false)

// 随机颜色（用于无封面时的 fallback）
const coverPalettes = [
  ['#14161f', '#b829dd'],
  ['#14161f', '#00f0ff'],
  ['#14161f', '#4d6bfa'],
  ['#0f1117', '#b829dd'],
  ['#0f1117', '#00f0ff'],
  ['#0a0a0f', '#4d6bfa'],
]
const avatarColors = ['#00f0ff', '#b829dd', '#4d6bfa', '#e4e5eb', '#8b8fa3']

function randomPalette(): string[] {
  return coverPalettes[Math.floor(Math.random() * coverPalettes.length)]
}
function randomAvatarColor(): string {
  return avatarColors[Math.floor(Math.random() * avatarColors.length)]
}

// 相对时间格式化
function formatRelativeTime(iso: string): string {
  const date = new Date(iso)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  if (diffSec < 60) return '> just now'
  const diffMin = Math.floor(diffSec / 60)
  if (diffMin < 60) return `> ${diffMin}m ago`
  const diffHour = Math.floor(diffMin / 60)
  if (diffHour < 24) return `> ${diffHour}h ago`
  const diffDay = Math.floor(diffHour / 24)
  if (diffDay < 30) return `> ${diffDay}d ago`
  return `> ${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

// 封面 fallback（当后端没有封面时，用分类确定颜色）
const categoryPaletteMap: Record<string, string[]> = {
  cine: ['#14161f', '#b829dd'],
  game: ['#14161f', '#00f0ff'],
  acgn: ['#14161f', '#4d6bfa'],
  chef: ['#0f1117', '#b829dd'],
  docu: ['#0f1117', '#00f0ff'],
  show: ['#0a0a0f', '#4d6bfa'],
  scifi: ['#14161f', '#b829dd'],
  tech: ['#14161f', '#00f0ff'],
  music: ['#0f1117', '#4d6bfa'],
  life: ['#0a0a0f', '#b829dd'],
}

function adaptVideoList(rawList: any[]): VideoItem[] {
  return rawList.map((v) => {
    const palette = v.cover_url
      ? undefined
      : (categoryPaletteMap[v.category] || randomPalette())
    return {
      id: v.id,
      title: v.title,
      cover_url: v.cover_url || undefined,
      coverGradient: palette || randomPalette(),
      duration: v.duration && v.duration > 0
        ? `${Math.floor(v.duration / 60).toString().padStart(2, '0')}:${(v.duration % 60).toString().padStart(2, '0')}`
        : '00:00',
      uploader: {
        name: v.user?.nickname || v.user?.username || '未知用户',
        avatarColor: randomAvatarColor(),
      },
      views: v.view_count || 0,
      danmaku: v.danmaku_count || 0,
      date: formatRelativeTime(v.created_at),
      category: v.category || 'all',
    }
  })
}

const allSectionVideos = ref<VideoItem[]>([])
const gameSectionVideos = ref<VideoItem[]>([])

const allRanks = ref(generateRanks([]))
const gameRanks = ref(generateRanks([]))

async function loadVideos() {
  loading.value = true
  try {
    const res = await getVideoList(1, 20)
    const adapted = adaptVideoList(res.list)
    videoList.value = adapted
    allSectionVideos.value = adapted
    gameSectionVideos.value = adapted.filter((v) => v.category === 'game')
    allRanks.value = generateRanks(adapted)
    gameRanks.value = generateRanks(gameSectionVideos.value)
  } catch (e) {
    console.error('加载视频列表失败:', e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadVideos()
})
</script>

<template>
  <div class="home">
    <NavBar />
    <SideBar />
    <ScrollProgress />
    <main class="main">
      <div class="main__inner">
        <BannerCarousel />
        <section class="main__section">
          <div class="main__section-header">
            <span class="main__section-prefix">//</span>
            <h2 class="main__section-title">推荐视频</h2>
          </div>
          <div class="main__section-divider" />
          <VideoGrid :videos="videoList" :loading="loading" />
        </section>
        <SectionBlock title="全部 · 最新归档" :videos="allSectionVideos" :ranks="allRanks" />
        <SectionBlock title="游戏 · 深度协议" :videos="gameSectionVideos" :ranks="gameRanks" />
      </div>
    </main>
    <BackToTop />
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh;
  background: #08090d;
}
.main {
  margin-left: 200px;
  margin-top: 64px;
  min-height: calc(100vh - 64px);
  padding: 24px;
}
.main__inner {
  max-width: 1800px;
  margin: 0 auto;
}
.main__section {
  margin-top: 48px;
  margin-bottom: 48px;
}
.main__section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}
.main__section-prefix {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #b829dd;
  font-weight: 600;
}
.main__section-title {
  font-size: 16px;
  font-weight: 600;
  color: #e4e5eb;
  margin: 0;
  letter-spacing: 0.02em;
}
.main__section-divider {
  height: 1px;
  background: linear-gradient(90deg, rgba(184, 41, 221, 0.4), transparent 80%);
  margin-bottom: 16px;
}
</style>
