<script setup lang="ts">
import type { VideoItem } from '../types/home'
import VideoCard from './VideoCard.vue'

interface Props {
  videos: VideoItem[]
  loading?: boolean
}

withDefaults(defineProps<Props>(), {
  loading: false,
})
</script>

<template>
  <div class="video-grid">
    <!-- 加载中骨架屏 -->
    <template v-if="loading">
      <div v-for="i in 10" :key="i" class="video-grid__skeleton">
        <div class="video-grid__skeleton-cover">
          <div class="video-grid__skeleton-shimmer" />
        </div>
        <div class="video-grid__skeleton-line video-grid__skeleton-line--title" />
        <div class="video-grid__skeleton-line video-grid__skeleton-line--short" />
      </div>
    </template>

    <!-- 实际卡片 -->
    <VideoCard
      v-for="video in videos"
      v-else
      :key="video.id"
      :video="video"
    />
  </div>
</template>

<style scoped>
.video-grid {
  display: grid;
  gap: 20px 20px;
}

/* 响应式列数 */
@media (min-width: 1600px) {
  .video-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}

@media (min-width: 1200px) and (max-width: 1599px) {
  .video-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 1199px) {
  .video-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* ========== 骨架屏 ========== */
.video-grid__skeleton {
  display: flex;
  flex-direction: column;
}

.video-grid__skeleton-cover {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #14161f;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.03);
}

.video-grid__skeleton-shimmer {
  position: absolute;
  top: 0;
  left: -150%;
  width: 150%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 240, 255, 0.06),
    transparent
  );
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  100% { left: 150%; }
}

.video-grid__skeleton-line {
  height: 12px;
  background: #14161f;
  border-radius: 2px;
  margin-top: 10px;
  overflow: hidden;
  position: relative;
}

.video-grid__skeleton-line::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 240, 255, 0.04),
    transparent
  );
  animation: shimmer 1.5s infinite;
}

.video-grid__skeleton-line--title {
  width: 90%;
}

.video-grid__skeleton-line--short {
  width: 60%;
  margin-top: 8px;
}
</style>
