<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import VideoGrid from '../components/VideoGrid.vue'
import { getFavorites } from '../api/video'
import type { VideoItem } from '../types/home'

const router = useRouter()
const videos = ref<VideoItem[]>([])
const loading = ref(true)
const error = ref('')

onMounted(async () => {
  try {
    const res = await getFavorites(1, 50)
    videos.value = res.list.map((v: any) => ({
      id: v.video_id || v.id,
      title: v.title,
      cover_url: v.cover_url,
      coverGradient: ['#0a0a0f', '#1a1c24'],
      duration: '00:00',
      uploader: {
        name: '',
        avatarColor: '#00f0ff',
      },
      views: v.view_count || 0,
      danmaku: 0,
      date: '',
      category: v.category || '',
    }))
  } catch (e: any) {
    error.value = e?.message || '加载失败'
  } finally {
    loading.value = false
  }
})

function handleCardClick(id: number) {
  router.push(`/video/${id}`)
}
</script>

<template>
  <div class="favorites-page">
    <NavBar />
    <SideBar />

    <main class="favorites-main">
      <div class="favorites-header">
        <h1 class="favorites-title">// MY_FAVORITES</h1>
        <span class="favorites-count">TOTAL: {{ videos.length }}</span>
      </div>

      <div v-if="error" class="favorites-error">{{ error }}</div>

      <VideoGrid
        :videos="videos"
        :loading="loading"
        @card-click="handleCardClick"
      />

      <div v-if="!loading && videos.length === 0" class="favorites-empty">
        > 暂无收藏视频
      </div>
    </main>
  </div>
</template>

<style scoped>
.favorites-page {
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

.favorites-main {
  margin-left: 200px;
  margin-top: 64px;
  min-height: calc(100vh - 64px);
  padding: 24px;
}

.favorites-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.favorites-title {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 18px;
  font-weight: 600;
  color: #e4e5eb;
  letter-spacing: 0.04em;
  margin: 0;
}

.favorites-count {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  border: 1px solid rgba(0, 240, 255, 0.3);
  padding: 2px 10px;
  border-radius: 2px;
}

.favorites-error,
.favorites-empty {
  padding: 80px 0;
  text-align: center;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
}
</style>
