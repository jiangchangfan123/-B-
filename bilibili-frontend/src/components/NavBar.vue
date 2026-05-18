<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { CategoryItem, HotSearchItem } from '../types/home'
import { categoryList, hotSearchList } from '../mock/homeData'

const categories = ref<CategoryItem[]>(categoryList.map((c, i) => ({ ...c, active: i === 0 })))
const hotSearches = ref<HotSearchItem[]>(hotSearchList)

const searchFocused = ref(false)
const searchKeyword = ref('')
const searchPanelRef = ref<HTMLDivElement | null>(null)
const searchInputRef = ref<HTMLInputElement | null>(null)

function onFocus() {
  searchFocused.value = true
}

function onBlurDelay() {
  setTimeout(() => {
    if (!searchPanelRef.value?.contains(document.activeElement)) {
      searchFocused.value = false
    }
  }, 100)
}

function onClickOutside(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (
    searchPanelRef.value &&
    !searchPanelRef.value.contains(target) &&
    searchInputRef.value !== target
  ) {
    searchFocused.value = false
  }
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    searchFocused.value = false
    searchInputRef.value?.blur()
  }
}

function switchCategory(id: string) {
  categories.value = categories.value.map((c) => ({ ...c, active: c.id === id }))
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
  document.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('click', onClickOutside)
  document.removeEventListener('keydown', onKeyDown)
})
</script>

<template>
  <nav class="navbar">
    <!-- 左侧：Logo + 搜索 -->
    <div class="navbar__left">
      <div class="navbar__logo">
        <span class="navbar__logo-text">// VOID</span>
      </div>

      <div ref="searchPanelRef" class="navbar__search-wrapper">
        <div
          class="navbar__search"
          :class="{ 'navbar__search--focused': searchFocused }"
        >
          <input
            ref="searchInputRef"
            v-model="searchKeyword"
            type="text"
            class="navbar__search-input"
            placeholder="> 输入检索指令..."
            @focus="onFocus"
            @blur="onBlurDelay"
          />
          <button class="navbar__search-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/>
            </svg>
          </button>
        </div>

        <!-- 热搜面板 -->
        <Transition name="search-panel">
          <div v-if="searchFocused" class="navbar__search-panel">
            <div class="navbar__search-panel-header">热门检索指令</div>
            <ul class="navbar__search-panel-list">
              <li
                v-for="item in hotSearches"
                :key="item.id"
                class="navbar__search-panel-item"
              >
                <span class="navbar__search-panel-mark">></span>
                <span class="navbar__search-panel-text">{{ item.text }}</span>
                <span
                  v-if="item.heat === 'hot'"
                  class="navbar__search-panel-tag navbar__search-panel-tag--hot"
                >[HOT]</span>
                <span
                  v-else-if="item.heat === 'new'"
                  class="navbar__search-panel-tag navbar__search-panel-tag--new"
                >[NEW]</span>
              </li>
            </ul>
          </div>
        </Transition>
      </div>
    </div>

    <!-- 中间：分类导航 -->
    <div class="navbar__center">
      <div
        v-for="cat in categories"
        :key="cat.id"
        class="navbar__cat-item"
        :class="{ 'navbar__cat-item--active': cat.active }"
        @click="switchCategory(cat.id)"
      >
        <span class="navbar__cat-label">{{ cat.label }}</span>
        <span class="navbar__cat-code">{{ cat.code }}</span>
        <div class="navbar__cat-indicator">
          <span class="navbar__cat-line" />
          <span class="navbar__cat-block" />
        </div>
      </div>
    </div>

    <!-- 右侧：功能区 -->
    <div class="navbar__right">
      <!-- 消息 -->
      <div class="navbar__icon-btn navbar__icon-btn--bell" title="消息">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"/>
          <path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"/>
        </svg>
        <span class="navbar__icon-badge" />
      </div>

      <!-- 历史 -->
      <div class="navbar__icon-btn navbar__icon-btn--history" title="历史记录">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/>
          <polyline points="12 6 12 12 16 14"/>
        </svg>
      </div>

      <!-- 投稿 -->
      <button class="navbar__upload-btn">
        <span class="navbar__upload-text">投稿</span>
      </button>

      <!-- 用户头像 -->
      <div class="navbar__avatar">
        <div class="navbar__avatar-inner" />
      </div>
    </div>

    <!-- 搜索遮罩 -->
    <Transition name="overlay">
      <div v-if="searchFocused" class="navbar__overlay" />
    </Transition>
  </nav>
</template>

<style scoped>
/* ========== 根元素 ========== */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background-color: #0a0a0f;
  background-image: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 240, 255, 0.03) 2px,
    rgba(0, 240, 255, 0.03) 4px
  );
  border-bottom: 1px solid rgba(0, 240, 255, 0.15);
  box-shadow: 0 1px 0 rgba(0, 240, 255, 0.08);
}

/* ========== 左侧 ========== */
.navbar__left {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-shrink: 0;
}

.navbar__logo-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 18px;
  font-weight: 700;
  color: #00f0ff;
  letter-spacing: 0.08em;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4), 0 0 20px rgba(0, 240, 255, 0.2);
}

.navbar__logo-text::after {
  content: '_';
  animation: blink 1s step-end infinite;
  color: #00f0ff;
}

@keyframes blink {
  50% { opacity: 0; }
}

/* ========== 搜索 ========== */
.navbar__search-wrapper {
  position: relative;
}

.navbar__search {
  display: flex;
  align-items: center;
  width: 340px;
  height: 36px;
  border: 1px solid #b829dd;
  border-radius: 4px;
  background: transparent;
  transition: width 0.3s ease, border-color 0.2s ease, box-shadow 0.2s ease;
  overflow: hidden;
}

.navbar__search--focused,
.navbar__search:focus-within {
  width: 480px;
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.3), inset 0 0 4px rgba(0, 240, 255, 0.05);
}

.navbar__search-input {
  flex: 1;
  height: 100%;
  padding: 0 12px;
  border: none;
  outline: none;
  background: transparent;
  color: #e4e5eb;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.02em;
}

.navbar__search-input::placeholder {
  color: #5a5d6e;
  font-family: 'JetBrains Mono', Consolas, monospace;
}

.navbar__search-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  color: #8b8fa3;
  cursor: pointer;
  transition: color 0.15s;
}

.navbar__search-btn:hover {
  color: #00f0ff;
}

/* ========== 热搜面板 ========== */
.navbar__search-panel {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  width: 480px;
  background: #08090d;
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: 4px;
  padding: 12px 0;
  z-index: 1001;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.6), 0 0 16px rgba(0, 240, 255, 0.08);
}

.navbar__search-panel-header {
  padding: 0 16px 8px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  margin-bottom: 4px;
}

.navbar__search-panel-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.navbar__search-panel-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  color: #8b8fa3;
  cursor: pointer;
  transition: background 0.1s, color 0.1s;
}

.navbar__search-panel-item:hover {
  background: rgba(0, 240, 255, 0.06);
  color: #00f0ff;
}

.navbar__search-panel-mark {
  color: #00f0ff;
  font-weight: 600;
}

.navbar__search-panel-text {
  flex: 1;
}

.navbar__search-panel-tag {
  font-size: 10px;
  padding: 1px 4px;
  border-radius: 2px;
  letter-spacing: 0.05em;
}

.navbar__search-panel-tag--hot {
  color: #b829dd;
  border: 1px solid rgba(184, 41, 221, 0.4);
  background: rgba(184, 41, 221, 0.08);
}

.navbar__search-panel-tag--new {
  color: #00f0ff;
  border: 1px solid rgba(0, 240, 255, 0.4);
  background: rgba(0, 240, 255, 0.08);
}

/* 面板动画 */
.search-panel-enter-active,
.search-panel-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.search-panel-enter-from,
.search-panel-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ========== 中间分类 ========== */
.navbar__center {
  display: flex;
  align-items: center;
  gap: 28px;
  flex-shrink: 0;
}

.navbar__cat-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  cursor: pointer;
  padding: 4px 0;
}

.navbar__cat-label {
  font-size: 14px;
  font-weight: 500;
  color: #a0a3b5;
  transition: none;
  letter-spacing: 0.02em;
}

.navbar__cat-code {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #5a5d6e;
  letter-spacing: 0.05em;
}

.navbar__cat-item:hover .navbar__cat-label {
  color: #e4e5eb;
  text-shadow: 0 0 8px rgba(228, 229, 235, 0.3);
}

.navbar__cat-item--active .navbar__cat-label {
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
}

.navbar__cat-item--active .navbar__cat-code {
  color: rgba(0, 240, 255, 0.5);
}

.navbar__cat-indicator {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 2px;
  height: 2px;
}

.navbar__cat-line {
  width: 0;
  height: 1px;
  background: #00f0ff;
  transition: width 0.1s step-end;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.5);
}

.navbar__cat-block {
  width: 4px;
  height: 4px;
  background: transparent;
  transition: background 0.1s step-end;
}

.navbar__cat-item--active .navbar__cat-line {
  width: 16px;
}

.navbar__cat-item--active .navbar__cat-block {
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.6);
  animation: pulse-block 2s ease-in-out infinite;
}

@keyframes pulse-block {
  0%, 100% { opacity: 0.6; box-shadow: 0 0 4px rgba(0, 240, 255, 0.3); }
  50% { opacity: 1; box-shadow: 0 0 10px rgba(0, 240, 255, 0.8); }
}

/* ========== 右侧功能 ========== */
.navbar__right {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-shrink: 0;
}

.navbar__icon-btn {
  position: relative;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8b8fa3;
  cursor: pointer;
  transition: color 0.15s;
}

.navbar__icon-btn:hover {
  color: #00f0ff;
}

.navbar__icon-btn--bell:hover {
  animation: glitch-bell 0.2s cubic-bezier(0.25, 0.46, 0.45, 0.94) both infinite;
  color: #00f0ff;
  text-shadow: 2px 0 #b829dd, -2px 0 #00f0ff;
}

@keyframes glitch-bell {
  0%, 100% { transform: translate(0); }
  20% { transform: translate(-1px, 1px); }
  40% { transform: translate(-1px, -1px); }
  60% { transform: translate(1px, 1px); }
  80% { transform: translate(1px, -1px); }
}

.navbar__icon-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 6px;
  height: 6px;
  background: #b829dd;
  border-radius: 50%;
  box-shadow: 0 0 4px rgba(184, 41, 221, 0.6);
}

/* 历史数据环 */
.navbar__icon-btn--history {
  position: relative;
}

.navbar__icon-btn--history::before {
  content: '';
  position: absolute;
  inset: -4px;
  border: 1px dashed rgba(0, 240, 255, 0);
  border-radius: 50%;
  transition: border-color 0.3s;
}

.navbar__icon-btn--history:hover::before {
  border-color: rgba(0, 240, 255, 0.4);
  animation: rotate-ring 4s linear infinite;
}

@keyframes rotate-ring {
  to { transform: rotate(360deg); }
}

/* 投稿按钮 */
.navbar__upload-btn {
  position: relative;
  overflow: hidden;
  height: 32px;
  padding: 0 18px;
  background: #b829dd;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.navbar__upload-btn:hover {
  background: #00f0ff;
  color: #08090d;
}

.navbar__upload-btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  pointer-events: none;
}

.navbar__upload-btn:hover::after {
  animation: sweep-upload 0.5s linear;
}

@keyframes sweep-upload {
  100% { left: 100%; }
}

/* 头像 */
.navbar__avatar {
  width: 40px;
  height: 40px;
  padding: 1px;
  border: 1px solid #00f0ff;
  transition: border-color 0.2s, box-shadow 0.2s;
  cursor: pointer;
}

.navbar__avatar:hover {
  border-color: #b829dd;
  box-shadow: 0 0 8px rgba(184, 41, 221, 0.5);
}

.navbar__avatar-inner {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #14161f, #2a2d3a);
}

/* ========== 遮罩层 ========== */
.navbar__overlay {
  position: fixed;
  inset: 0;
  background: rgba(8, 9, 13, 0.7);
  z-index: 999;
  backdrop-filter: blur(2px);
}

.overlay-enter-active,
.overlay-leave-active {
  transition: opacity 0.2s ease;
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}
</style>