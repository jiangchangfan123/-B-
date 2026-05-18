<script setup lang="ts">
import NavBar from './components/NavBar.vue'
import SideBar from './components/SideBar.vue'
import BannerCarousel from './components/BannerCarousel.vue'
import VideoGrid from './components/VideoGrid.vue'
import SectionBlock from './components/SectionBlock.vue'
import BackToTop from './components/BackToTop.vue'
import ScrollProgress from './components/ScrollProgress.vue'
import { videoList, allSectionVideos, gameSectionVideos, generateRanks } from './mock/homeData'

const allRanks = generateRanks(allSectionVideos)
const gameRanks = generateRanks(gameSectionVideos)
</script>

<template>
  <div class="app">
    <!-- 顶部导航 -->
    <NavBar />

    <!-- 左侧边栏 -->
    <SideBar />

    <!-- 滚动进度条 -->
    <ScrollProgress />

    <!-- 主体内容 -->
    <main class="main">
      <div class="main__inner">
        <!-- 轮播 Banner -->
        <BannerCarousel />

        <!-- 推荐视频网格 -->
        <section class="main__section">
          <div class="main__section-header">
            <span class="main__section-prefix">//</span>
            <h2 class="main__section-title">推荐视频</h2>
          </div>
          <div class="main__section-divider" />
          <VideoGrid :videos="videoList" />
        </section>

        <!-- 分区推荐一 -->
        <SectionBlock
          title="全部 · 最新归档"
          :videos="allSectionVideos"
          :ranks="allRanks"
        />

        <!-- 分区推荐二 -->
        <SectionBlock
          title="游戏 · 深度协议"
          :videos="gameSectionVideos"
          :ranks="gameRanks"
        />
      </div>
    </main>

    <!-- 回到顶部 -->
    <BackToTop />
  </div>
</template>

<style>
/* ========== 全局样式 ========== */
* {
  box-sizing: border-box;
}

html {
  min-width: 1280px;
  scroll-behavior: smooth;
}

body {
  margin: 0;
  padding: 0;
  background: #08090d;
  color: #e4e5eb;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto,
    'Helvetica Neue', Arial, 'Noto Sans', sans-serif;
  letter-spacing: 0.02em;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow-x: hidden;
}

/* 等宽字体回退 */
@font-face {
  font-family: 'JetBrains Mono';
  src: local('JetBrains Mono'), local('Consolas'), local('Monaco'), local('monospace');
}
</style>

<style scoped>
.app {
  min-height: 100vh;
  background: #08090d;
}

/* 主体内容区 */
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

/* 区块标题 */
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
