<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import { useUserStore } from '../stores/user'
import {
  getNotifications,
  markAsRead,
  markAllAsRead,
  deleteNotification,
} from '../api/notification'
import type { NotificationItem } from '../api/notification'

const router = useRouter()
const userStore = useUserStore()

const notifications = ref<NotificationItem[]>([])
const loading = ref(false)
const activeTab = ref<'all' | 'unread'>('all')
const page = ref(1)
const hasMore = ref(false)

const typeLabelMap: Record<number, string> = {
  1: '回复',
  2: '赞',
  3: '赞',
}

const typeColorMap: Record<number, string> = {
  1: '#00f0ff',
  2: '#b829dd',
  3: '#b829dd',
}

// 头像 URL 拼接
const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const SERVER_BASE = API_BASE.replace(/\/api\/v1\/?$/, '')

function getFullUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return SERVER_BASE + url
}

async function loadNotifications(reset = false) {
  if (reset) {
    page.value = 1
    notifications.value = []
  }
  loading.value = true
  try {
    const res = await getNotifications(page.value, 20, activeTab.value === 'unread')
    if (reset) {
      notifications.value = res.list
    } else {
      notifications.value.push(...res.list)
    }
    hasMore.value = res.list.length === 20
  } catch (e) {
    console.error('加载消息失败:', e)
  } finally {
    loading.value = false
  }
}

function onTabChange(tab: 'all' | 'unread') {
  activeTab.value = tab
  loadNotifications(true)
}

function onLoadMore() {
  page.value++
  loadNotifications()
}

async function onMarkAsRead(id: number) {
  try {
    await markAsRead(id)
    const item = notifications.value.find((n) => n.id === id)
    if (item) item.is_read = true
  } catch (e) {
    console.error(e)
  }
}

async function onMarkAllAsRead() {
  try {
    await markAllAsRead()
    notifications.value.forEach((n) => (n.is_read = true))
  } catch (e) {
    console.error(e)
  }
}

async function onDelete(id: number) {
  try {
    await deleteNotification(id)
    notifications.value = notifications.value.filter((n) => n.id !== id)
  } catch (e) {
    console.error(e)
  }
}

function onClickItem(item: NotificationItem) {
  if (!item.is_read) {
    onMarkAsRead(item.id)
  }
  if (item.type === 3) {
    router.push(`/video/${item.related_id}`)
  } else {
    // 评论相关，跳到视频详情页并定位到评论
    router.push(`/video/${item.related_id}`)
  }
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  return `${Math.floor(diff / 86400)}天前`
}

onMounted(() => {
  loadNotifications(true)
})
</script>

<template>
  <div class="notifications">
    <NavBar />
    <SideBar />
    <main class="main">
      <div class="main__inner">
        <!-- 标题区 -->
        <div class="header">
          <h1 class="header__title">消息中心</h1>
          <button class="header__action" @click="onMarkAllAsRead">
            全部已读
          </button>
        </div>

        <!-- 标签切换 -->
        <div class="tabs">
          <button
            class="tabs__btn"
            :class="{ active: activeTab === 'all' }"
            @click="onTabChange('all')"
          >
            全部消息
          </button>
          <button
            class="tabs__btn"
            :class="{ active: activeTab === 'unread' }"
            @click="onTabChange('unread')"
          >
            未读消息
          </button>
        </div>

        <!-- 消息列表 -->
        <div class="list">
          <div
            v-for="item in notifications"
            :key="item.id"
            class="item"
            :class="{ unread: !item.is_read }"
            @click="onClickItem(item)"
          >
            <div class="item__avatar">
              <img
                v-if="item.trigger_user?.avatar"
                :src="getFullUrl(item.trigger_user.avatar)"
                alt="avatar"
              />
              <div v-else class="item__avatar-fallback" />
            </div>
            <div class="item__content">
              <div class="item__top">
                <span
                  class="item__type"
                  :style="{ color: typeColorMap[item.type] }"
                >
                  {{ typeLabelMap[item.type] }}
                </span>
                <span class="item__title">{{ item.title }}</span>
              </div>
              <p v-if="item.content" class="item__text">
                {{ item.content }}
              </p>
              <span class="item__time">{{ formatTime(item.created_at) }}</span>
            </div>
            <div class="item__actions">
              <button
                v-if="!item.is_read"
                class="item__btn"
                @click.stop="onMarkAsRead(item.id)"
              >
                标为已读
              </button>
              <button class="item__btn item__btn--danger" @click.stop="onDelete(item.id)">
                删除
              </button>
            </div>
          </div>

          <div v-if="notifications.length === 0 && !loading" class="empty">
            暂无消息
          </div>

          <div v-if="hasMore" class="load-more">
            <button class="load-more__btn" @click="onLoadMore">
              加载更多
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.notifications {
  min-height: 100vh;
  background: #08090d;
}

.main {
  margin-left: 200px;
  padding-top: 64px;
  min-height: 100vh;
}

.main__inner {
  max-width: 900px;
  margin: 0 auto;
  padding: 24px 32px;
}

/* 标题 */
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.header__title {
  font-size: 20px;
  font-weight: 600;
  color: #e4e5eb;
  margin: 0;
}

.header__action {
  padding: 6px 14px;
  background: transparent;
  border: 1px solid rgba(0, 240, 255, 0.3);
  border-radius: 4px;
  color: #00f0ff;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.header__action:hover {
  background: rgba(0, 240, 255, 0.08);
}

/* 标签 */
.tabs {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  padding-bottom: 12px;
}

.tabs__btn {
  background: transparent;
  border: none;
  color: #5a5d6e;
  font-size: 14px;
  cursor: pointer;
  padding: 4px 0;
  position: relative;
  transition: color 0.2s;
}

.tabs__btn:hover {
  color: #e4e5eb;
}

.tabs__btn.active {
  color: #00f0ff;
}

.tabs__btn.active::after {
  content: '';
  position: absolute;
  bottom: -13px;
  left: 0;
  right: 0;
  height: 2px;
  background: #00f0ff;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.4);
}

/* 列表 */
.list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
}

.item:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(0, 240, 255, 0.15);
}

.item.unread {
  background: rgba(0, 240, 255, 0.03);
  border-left: 2px solid #00f0ff;
}

.item__avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  border: 1px solid rgba(0, 240, 255, 0.2);
}

.item__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.item__avatar-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #b829dd, #4d6bfa);
}

.item__content {
  flex: 1;
  min-width: 0;
}

.item__top {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.item__type {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 6px;
  border: 1px solid currentColor;
  border-radius: 2px;
  opacity: 0.8;
}

.item__title {
  font-size: 14px;
  color: #e4e5eb;
  font-weight: 500;
}

.item__text {
  font-size: 13px;
  color: #8b8fa3;
  margin: 4px 0 6px;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.item__time {
  font-size: 12px;
  color: #5a5d6e;
}

.item__actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s;
}

.item:hover .item__actions {
  opacity: 1;
}

.item__btn {
  padding: 4px 10px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  color: #8b8fa3;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.item__btn:hover {
  border-color: rgba(0, 240, 255, 0.3);
  color: #00f0ff;
}

.item__btn--danger:hover {
  border-color: rgba(255, 80, 80, 0.3);
  color: #ff5050;
}

/* 空状态 */
.empty {
  text-align: center;
  padding: 60px 0;
  color: #5a5d6e;
  font-size: 14px;
}

/* 加载更多 */
.load-more {
  text-align: center;
  padding: 16px 0;
}

.load-more__btn {
  padding: 8px 24px;
  background: transparent;
  border: 1px solid rgba(0, 240, 255, 0.2);
  border-radius: 4px;
  color: #00f0ff;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.load-more__btn:hover {
  background: rgba(0, 240, 255, 0.06);
}
</style>
