<script setup lang="ts">
import { ref } from 'vue'
import type { MenuItem } from '../types/home'
import { menuList, systemMenuList } from '../mock/homeData'

const menus = ref<MenuItem[]>(menuList)
const systemMenus = ref<MenuItem[]>(systemMenuList)
const expanded = ref(false)

function setActive(id: string) {
  menus.value = menus.value.map((m) => ({ ...m, active: m.id === id }))
}

function toggleExpand() {
  expanded.value = !expanded.value
}
</script>

<template>
  <aside class="sidebar">
    <!-- 主菜单 -->
    <nav class="sidebar__nav">
      <div
        v-for="item in menus"
        :key="item.id"
        class="sidebar__item"
        :class="{ 'sidebar__item--active': item.active }"
        @click="setActive(item.id)"
      >
        <svg
          class="sidebar__icon"
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path :d="item.icon" />
        </svg>
        <span class="sidebar__label">{{ item.label }}</span>
        <span
          v-if="item.expandable"
          class="sidebar__arrow"
          :class="{ 'sidebar__arrow--expanded': expanded }"
          @click.stop="toggleExpand"
        >
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9" />
          </svg>
        </span>
      </div>
    </nav>

    <!-- 系统面板 -->
    <div class="sidebar__system">
      <div class="sidebar__system-header">
        <span class="sidebar__system-line" />
        <span class="sidebar__system-title">系统</span>
      </div>
      <div
        v-for="item in systemMenus"
        :key="item.id"
        class="sidebar__item sidebar__item--system"
        :class="{ 'sidebar__item--active': item.active }"
        @click="setActive(item.id)"
      >
        <svg
          class="sidebar__icon"
          width="18"
          height="18"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="1.5"
          stroke-linecap="round"
          stroke-linejoin="round"
        >
          <path :d="item.icon" />
        </svg>
        <span class="sidebar__label">{{ item.label }}</span>
      </div>

      <!-- 控制台模式开关 -->
      <div class="sidebar__toggle-row">
        <span class="sidebar__toggle-label">控制台模式</span>
        <div class="sidebar__toggle">
          <div class="sidebar__toggle-track">
            <div class="sidebar__toggle-thumb" />
          </div>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 64px;
  left: 0;
  width: 200px;
  bottom: 0;
  background: #0f1117;
  border-right: 1px solid #b829dd;
  z-index: 900;
  display: flex;
  flex-direction: column;
  padding: 8px 0;
  animation: breath-border 3s ease-in-out infinite;
}

@keyframes breath-border {
  0%, 100% {
    border-color: rgba(184, 41, 221, 0.3);
    box-shadow: 1px 0 4px rgba(184, 41, 221, 0.1);
  }
  50% {
    border-color: rgba(184, 41, 221, 0.8);
    box-shadow: 1px 0 12px rgba(184, 41, 221, 0.25);
  }
}

.sidebar__nav {
  flex: 1;
  padding: 0 8px;
}

.sidebar__item {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  height: 48px;
  padding: 0 16px;
  margin-bottom: 2px;
  border-radius: 4px;
  color: #8b8fa3;
  cursor: pointer;
  overflow: hidden;
  transition: color 0.15s;
}

.sidebar__item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 0;
  background: rgba(0, 240, 255, 0.06);
  transition: width 0.3s ease;
}

.sidebar__item:hover::before {
  width: 100%;
}

.sidebar__item:hover::after {
  content: '';
  position: absolute;
  left: -100%;
  top: 50%;
  width: 100%;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(0, 240, 255, 0.15), transparent);
  animation: scan-sweep 0.3s ease forwards;
}

@keyframes scan-sweep {
  to { left: 100%; }
}

.sidebar__item:hover {
  color: #e4e5eb;
}

.sidebar__item--active {
  color: #fff;
  background: rgba(0, 240, 255, 0.06);
}

.sidebar__item--active::before {
  width: 2px;
  background: #00f0ff;
}

.sidebar__item--active .sidebar__icon {
  color: #00f0ff;
  filter: drop-shadow(0 0 4px #00f0ff);
}

.sidebar__icon {
  flex-shrink: 0;
  color: currentColor;
  transition: color 0.15s, filter 0.15s;
}

.sidebar__label {
  font-size: 14px;
  letter-spacing: 0.02em;
  position: relative;
  z-index: 1;
}

.sidebar__arrow {
  margin-left: auto;
  transition: transform 0.2s;
  opacity: 0.5;
}

.sidebar__arrow--expanded {
  transform: rotate(180deg);
}

/* 系统区域 */
.sidebar__system {
  padding: 0 8px 16px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  margin-top: auto;
}

.sidebar__system-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px 8px;
}

.sidebar__system-line {
  width: 12px;
  height: 1px;
  background: #5a5d6e;
}

.sidebar__system-title {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.1em;
  text-transform: uppercase;
}

.sidebar__item--system {
  margin-bottom: 2px;
}

/* 开关 */
.sidebar__toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  margin-top: 4px;
}

.sidebar__toggle-label {
  font-size: 13px;
  color: #8b8fa3;
  letter-spacing: 0.02em;
}

.sidebar__toggle {
  cursor: pointer;
}

.sidebar__toggle-track {
  width: 36px;
  height: 18px;
  border-radius: 9px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.1);
  position: relative;
  transition: background 0.2s, box-shadow 0.2s;
}

.sidebar__toggle-track:hover {
  background: rgba(184, 41, 221, 0.15);
  box-shadow: 0 0 8px rgba(184, 41, 221, 0.2);
}

.sidebar__toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #b829dd;
  box-shadow: 0 0 4px rgba(184, 41, 221, 0.5);
  transition: left 0.2s;
}
</style>
