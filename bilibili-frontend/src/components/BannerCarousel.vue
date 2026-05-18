<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { BannerItem } from '../types/home'
import { bannerList } from '../mock/homeData'

const banners = ref<BannerItem[]>(bannerList)
const currentIndex = ref(0)
const isHovering = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

function startAutoPlay() {
  timer = setInterval(() => {
    if (!isHovering.value) {
      currentIndex.value = (currentIndex.value + 1) % banners.value.length
    }
  }, 5000)
}

function stopAutoPlay() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

function goTo(index: number) {
  currentIndex.value = index
}

function prev() {
  currentIndex.value =
    (currentIndex.value - 1 + banners.value.length) % banners.value.length
}

function next() {
  currentIndex.value = (currentIndex.value + 1) % banners.value.length
}

onMounted(startAutoPlay)
onUnmounted(stopAutoPlay)
</script>

<template>
  <div
    class="banner"
    @mouseenter="isHovering = true"
    @mouseleave="isHovering = false"
  >
    <div class="banner__track">
      <div
        v-for="(item, index) in banners"
        :key="item.id"
        class="banner__slide"
        :class="{ 'banner__slide--active': index === currentIndex }"
        :style="{
          background: `linear-gradient(135deg, ${item.gradientColors[0]}, ${item.gradientColors[1]}, ${item.gradientColors[2]})`
        }"
      >
        <div class="banner__content">
          <h2 class="banner__title">{{ item.title }}</h2>
          <p class="banner__subtitle">{{ item.subtitle }}</p>
        </div>
      </div>
    </div>

    <!-- 切换箭头 -->
    <button class="banner__arrow banner__arrow--prev" @click="prev">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6" />
      </svg>
    </button>
    <button class="banner__arrow banner__arrow--next" @click="next">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6" />
      </svg>
    </button>

    <!-- 指示器 -->
    <div class="banner__dots">
      <button
        v-for="(_, index) in banners"
        :key="index"
        class="banner__dot"
        :class="{ 'banner__dot--active': index === currentIndex }"
        @click="goTo(index)"
      />
    </div>
  </div>
</template>

<style scoped>
.banner {
  position: relative;
  width: 100%;
  height: 280px;
  border-radius: 4px;
  border: 1px solid rgba(0, 240, 255, 0.15);
  overflow: hidden;
  transition: border-color 0.2s;
}

.banner:hover {
  border-color: rgba(0, 240, 255, 0.35);
}

.banner__track {
  position: relative;
  width: 100%;
  height: 100%;
}

.banner__slide {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: flex-end;
  padding: 32px 40px;
  opacity: 0;
  transition: opacity 0.5s ease;
}

.banner__slide--active {
  opacity: 1;
}

.banner__content {
  position: relative;
  z-index: 1;
}

.banner__title {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
  letter-spacing: 0.02em;
  text-shadow: 0 0 20px rgba(184, 41, 221, 0.5);
}

.banner__subtitle {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  color: #00f0ff;
  margin: 0;
  letter-spacing: 0.04em;
}

/* 切换箭头 */
.banner__arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 4px;
  color: #00f0ff;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s, background 0.2s, border-color 0.2s;
  z-index: 2;
}

.banner:hover .banner__arrow {
  opacity: 1;
}

.banner__arrow:hover {
  background: #00f0ff;
  border-color: #00f0ff;
  color: #08090d;
}

.banner__arrow--prev {
  left: 16px;
}

.banner__arrow--next {
  right: 16px;
}

/* 指示器 */
.banner__dots {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 8px;
  z-index: 2;
}

.banner__dot {
  width: 12px;
  height: 2px;
  background: #5a5d6e;
  border: none;
  border-radius: 1px;
  cursor: pointer;
  transition: width 0.3s ease, background 0.2s;
}

.banner__dot--active {
  width: 24px;
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.5);
}
</style>