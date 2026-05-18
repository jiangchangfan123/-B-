<script setup lang="ts">
import type { VideoItem, RankItem } from '../types/home'
import VideoCard from './VideoCard.vue'
import { formatNumber } from '../mock/homeData'

interface Props {
  title: string
  videos: VideoItem[]
  ranks: RankItem[]
}

defineProps<Props>()
</script>

<template>
  <section class="section-block">
    <!-- 标题 -->
    <div class="section-block__header">
      <div class="section-block__title-wrap">
        <span class="section-block__prefix">//</span>
        <h2 class="section-block__title">{{ title }}</h2>
      </div>
      <span class="section-block__more">> VIEW_ALL</span>
    </div>

    <div class="section-block__divider" />

    <div class="section-block__body">
      <!-- 左侧视频 -->
      <div class="section-block__videos">
        <VideoCard
          v-for="video in videos.slice(0, 5)"
          :key="video.id"
          :video="video"
          size="small"
        />
      </div>

      <!-- 右侧排行榜 -->
      <div class="section-block__rank">
        <div class="section-block__rank-header">
          <span class="section-block__rank-icon">▲</span>
          <span>热门排行</span>
        </div>
        <div class="section-block__rank-list">
          <div
            v-for="item in ranks"
            :key="item.id"
            class="section-block__rank-item"
          >
            <span
              class="section-block__rank-num"
              :class="{
                'section-block__rank-num--1': item.rank === 1,
                'section-block__rank-num--2': item.rank === 2,
                'section-block__rank-num--3': item.rank === 3,
              }"
            >
              [0{{ item.rank }}]
            </span>
            <span class="section-block__rank-title">{{ item.title }}</span>
            <span class="section-block__rank-views">{{ formatNumber(item.views) }}</span>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.section-block {
  margin-bottom: 48px;
}

.section-block__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.section-block__title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.section-block__prefix {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  color: #b829dd;
  font-weight: 600;
}

.section-block__title {
  font-size: 16px;
  font-weight: 600;
  color: #e4e5eb;
  margin: 0;
  letter-spacing: 0.02em;
}

.section-block__more {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  cursor: pointer;
  letter-spacing: 0.05em;
  transition: text-shadow 0.2s;
}

.section-block__more:hover {
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
}

.section-block__divider {
  height: 1px;
  background: linear-gradient(90deg, rgba(184, 41, 221, 0.4), transparent 80%);
  margin-bottom: 16px;
}

/* 主体 */
.section-block__body {
  display: flex;
  gap: 24px;
}

.section-block__videos {
  flex: 1;
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
}

@media (max-width: 1599px) {
  .section-block__videos {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 1199px) {
  .section-block__videos {
    grid-template-columns: repeat(3, 1fr);
  }
}

/* 排行榜 */
.section-block__rank {
  width: 240px;
  flex-shrink: 0;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  padding: 12px 16px;
}

.section-block__rank-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.05em;
  margin-bottom: 10px;
  padding-bottom: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.section-block__rank-icon {
  color: #b829dd;
}

.section-block__rank-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 0;
  cursor: pointer;
  transition: background 0.15s;
  border-radius: 2px;
  position: relative;
  overflow: hidden;
}

.section-block__rank-item:hover {
  background: rgba(0, 240, 255, 0.04);
}

.section-block__rank-item:hover::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.4);
}

.section-block__rank-num {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
  width: 32px;
}

.section-block__rank-num--1 {
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
}

.section-block__rank-num--2 {
  color: #e4e5eb;
}

.section-block__rank-num--3 {
  color: #5a5d6e;
}

.section-block__rank-title {
  flex: 1;
  font-size: 12px;
  color: #8b8fa3;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.section-block__rank-views {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  flex-shrink: 0;
}
</style>
