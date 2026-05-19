<script setup lang="ts">
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import BannerCarousel from '../components/BannerCarousel.vue'
import VideoGrid from '../components/VideoGrid.vue'
import SectionBlock from '../components/SectionBlock.vue'
import BackToTop from '../components/BackToTop.vue'
import ScrollProgress from '../components/ScrollProgress.vue'
import { videoList, allSectionVideos, gameSectionVideos, generateRanks } from '../mock/homeData'

const allRanks = generateRanks(allSectionVideos)
const gameRanks = generateRanks(gameSectionVideos)
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
          <VideoGrid :videos="videoList" />
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
