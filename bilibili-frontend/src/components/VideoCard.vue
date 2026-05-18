<script setup lang="ts">
import type { VideoItem } from '../types/home'
import { formatNumber } from '../mock/homeData'

interface Props {
  video: VideoItem
  size?: 'normal' | 'small'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'normal',
})

function coverStyle(colors: string[]) {
  return {
    background: `linear-gradient(135deg, ${colors[0]} 0%, ${colors[1]} 100%)`,
  }
}
</script>

<template>
  <div class="video-card" :class="{ 'video-card--small': props.size === 'small' }">
    <!-- 封面区 -->
    <div class="video-card__cover">
      <div class="video-card__cover-img" :style="coverStyle(props.video.coverGradient)">
        <div class="video-card__cover-pattern" />
      </div>
      <div class="video-card__duration">{{ props.video.duration }}</div>

      <!-- hover 收藏图标 -->
      <div class="video-card__save">
        <span class="video-card__save-text">[+] SAVE</span>
      </div>

      <!-- 扫描线动画 -->
      <div class="video-card__scan-line" />

      <!-- 故障闪烁层 -->
      <div class="video-card__glitch" />
    </div>

    <!-- 信息区 -->
    <div class="video-card__info">
      <!-- 标题 -->
      <div class="video-card__title-row">
        <span v-if="props.video.tags?.length" class="video-card__tag">
          {{ props.video.tags[0] }}
        </span>
        <h3 class="video-card__title">{{ props.video.title }}</h3>
      </div>

      <!-- UP 主 -->
      <div class="video-card__uploader">
        <div
          class="video-card__avatar"
          :style="{ background: props.video.uploader.avatarColor }"
        />
        <span class="video-card__name">{{ props.video.uploader.name }}</span>
        <span class="video-card__date">{{ props.video.date }}</span>
      </div>

      <!-- 数据行 -->
      <div class="video-card__stats">
        <span class="video-card__stat">
          <span class="video-card__stat-icon">▶</span>
          {{ formatNumber(props.video.views) }}
        </span>
        <span class="video-card__stat">
          <span class="video-card__stat-icon">◆</span>
          {{ formatNumber(props.video.danmaku) }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.video-card {
  cursor: pointer;
}

.video-card__cover {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: border-color 0.2s, box-shadow 0.2s;
}

.video-card:hover .video-card__cover {
  border-color: rgba(0, 240, 255, 0.3);
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.15), inset 0 0 8px rgba(0, 240, 255, 0.05);
}

.video-card__cover-img {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-card__cover-pattern {
  width: 40%;
  height: 40%;
  background: rgba(255, 255, 255, 0.04);
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
}

/* 时长标签 */
.video-card__duration {
  position: absolute;
  bottom: 6px;
  right: 6px;
  padding: 2px 6px;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 2px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #00f0ff;
  letter-spacing: 0.02em;
}

/* 收藏 */
.video-card__save {
  position: absolute;
  top: 6px;
  right: 6px;
  opacity: 0;
  transition: opacity 0.2s;
}

.video-card:hover .video-card__save {
  opacity: 1;
}

.video-card__save-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #00f0ff;
  letter-spacing: 0.05em;
}

/* 扫描线 */
.video-card__scan-line {
  position: absolute;
  top: -2px;
  left: 0;
  width: 100%;
  height: 2px;
  background: rgba(0, 240, 255, 0.6);
  box-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
  opacity: 0;
}

.video-card:hover .video-card__scan-line {
  animation: scan-down 0.6s linear;
}

@keyframes scan-down {
  0% { top: -2px; opacity: 1; }
  100% { top: 100%; opacity: 0; }
}

/* 故障闪烁 */
.video-card__glitch {
  position: absolute;
  inset: 0;
  background: rgba(0, 240, 255, 0);
  pointer-events: none;
  opacity: 0;
}

.video-card:hover .video-card__glitch {
  animation: glitch-flash 0.1s ease;
}

@keyframes glitch-flash {
  0% { opacity: 0; transform: translate(0); }
  25% { opacity: 0.08; transform: translate(2px, 0); background: rgba(0, 240, 255, 0.05); }
  50% { opacity: 0; transform: translate(-2px, 0); }
  75% { opacity: 0.06; transform: translate(1px, 0); background: rgba(184, 41, 221, 0.05); }
  100% { opacity: 0; transform: translate(0); }
}

/* ========== 信息区 ========== */
.video-card__info {
  margin-top: 10px;
}

.video-card__title-row {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-bottom: 6px;
}

.video-card__tag {
  flex-shrink: 0;
  padding: 1px 4px;
  border: 1px solid rgba(184, 41, 221, 0.5);
  border-radius: 2px;
  font-size: 10px;
  color: #b829dd;
  font-family: 'JetBrains Mono', Consolas, monospace;
  letter-spacing: 0.02em;
  margin-top: 2px;
}

.video-card__title {
  font-size: 14px;
  font-weight: 500;
  color: #e4e5eb;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  letter-spacing: 0.02em;
  transition: color 0.15s;
}

.video-card:hover .video-card__title {
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.3);
}

/* UP 主 */
.video-card__uploader {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.video-card__avatar {
  width: 24px;
  height: 24px;
  clip-path: polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%);
  border: 1px solid rgba(0, 240, 255, 0.4);
  flex-shrink: 0;
}

.video-card__name {
  font-size: 12px;
  color: #8b8fa3;
  letter-spacing: 0.02em;
}

.video-card__date {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-left: auto;
}

/* 数据行 */
.video-card__stats {
  display: flex;
  align-items: center;
  gap: 12px;
}

.video-card__stat {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #5a5d6e;
  letter-spacing: 0.02em;
}

.video-card__stat-icon {
  font-size: 10px;
  margin-right: 2px;
}

/* 小尺寸模式 */
.video-card--small .video-card__title {
  font-size: 13px;
  -webkit-line-clamp: 1;
}

.video-card--small .video-card__uploader {
  margin-bottom: 2px;
}
</style>