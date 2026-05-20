<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useUserStore } from '../stores/user'
import {
  createComment,
  getVideoComments,
  deleteComment,
  toggleCommentLike,
} from '../api/comment'
import type { CommentVO, ReplyVO } from '../types/comment'

const props = defineProps<{ videoId: number }>()
const emit = defineEmits<{
  (e: 'updateCount', count: number): void
}>()

const userStore = useUserStore()
const comments = ref<CommentVO[]>([])
const total = ref(0)
const page = ref(1)
const size = 20
const sort = ref<'time' | 'hot'>('time')
const loading = ref(false)
const hasMore = computed(() => comments.value.length < total.value)

// 输入框
const inputContent = ref('')
const replyTarget = ref<{ commentId: number; toUserName: string } | null>(null)
const submitting = ref(false)

// 展开的回复
const expandedReplies = ref<Set<number>>(new Set())

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const SERVER_BASE = API_BASE.replace(/\/api\/v1\/?$/, '')

function getFullUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return SERVER_BASE + url
}

function formatTime(iso: string): string {
  const d = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now.getTime() - d.getTime()) / 1000)
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}天前`
  return d.toLocaleDateString('zh-CN')
}

async function loadComments(append = false) {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getVideoComments(props.videoId, page.value, size, sort.value)
    if (append) {
      comments.value.push(...res.list)
    } else {
      comments.value = res.list
    }
    total.value = res.total
    emit('updateCount', res.total)
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function switchSort(s: 'time' | 'hot') {
  if (sort.value === s) return
  sort.value = s
  page.value = 1
  loadComments(false)
}

function loadMore() {
  page.value++
  loadComments(true)
}

async function submitComment() {
  const content = inputContent.value.trim()
  if (!content) return
  if (!userStore.isLoggedIn) return

  submitting.value = true
  try {
    const parentId = replyTarget.value?.commentId || 0
    await createComment(props.videoId, { content, parent_id: parentId })
    inputContent.value = ''
    replyTarget.value = null
    page.value = 1
    await loadComments(false)
  } catch {
    alert('评论失败')
  } finally {
    submitting.value = false
  }
}

function startReply(commentId: number, toUserName: string) {
  replyTarget.value = { commentId, toUserName }
  document.getElementById('comment-input')?.focus()
}

function cancelReply() {
  replyTarget.value = null
}

async function onDelete(commentId: number) {
  if (!confirm('确定删除这条评论吗？')) return
  try {
    await deleteComment(commentId)
    comments.value = comments.value.filter(c => c.id !== commentId)
    // 同时删除回复
    comments.value.forEach(c => {
      c.replies = c.replies.filter(r => r.id !== commentId)
    })
    total.value = Math.max(0, total.value - 1)
    emit('updateCount', total.value)
  } catch {
    alert('删除失败')
  }
}

async function onToggleLike(commentId: number, isReply: boolean, parentId?: number) {
  if (!userStore.isLoggedIn) {
    alert('请先登录')
    return
  }
  try {
    const res = await toggleCommentLike(commentId)
    if (isReply && parentId) {
      const parent = comments.value.find(c => c.id === parentId)
      if (parent) {
        const reply = parent.replies.find(r => r.id === commentId)
        if (reply) {
          reply.is_liked = res.liked
          reply.like_count += res.liked ? 1 : -1
        }
      }
    } else {
      const comment = comments.value.find(c => c.id === commentId)
      if (comment) {
        comment.is_liked = res.liked
        comment.like_count += res.liked ? 1 : -1
      }
    }
  } catch {
    // ignore
  }
}

function toggleReplies(commentId: number) {
  if (expandedReplies.value.has(commentId)) {
    expandedReplies.value.delete(commentId)
  } else {
    expandedReplies.value.add(commentId)
  }
}

watch(() => props.videoId, () => {
  page.value = 1
  comments.value = []
  loadComments(false)
}, { immediate: true })
</script>

<template>
  <div class="comment-section">
    <!-- 头部 -->
    <div class="comment-header">
      <span class="comment-count">评论 ({{ total }})</span>
      <div class="sort-tabs">
        <button
          class="sort-btn"
          :class="{ active: sort === 'hot' }"
          @click="switchSort('hot')"
        >最热</button>
        <button
          class="sort-btn"
          :class="{ active: sort === 'time' }"
          @click="switchSort('time')"
        >最新</button>
      </div>
    </div>

    <!-- 输入框 -->
    <div v-if="userStore.isLoggedIn" class="comment-input-area">
      <div class="input-avatar">
        {{ (userStore.userInfo?.nickname || userStore.userInfo?.username || '?').charAt(0).toUpperCase() }}
      </div>
      <div class="input-box-wrap">
        <div v-if="replyTarget" class="reply-hint">
          回复 <span class="reply-name">@{{ replyTarget.toUserName }}</span>
          <button class="cancel-reply" @click="cancelReply">取消</button>
        </div>
        <textarea
          id="comment-input"
          v-model="inputContent"
          class="comment-textarea"
          :placeholder="replyTarget ? '写下你的回复...' : '发一条友善的评论...'"
          rows="3"
        />
        <div class="input-actions">
          <button
            class="submit-btn"
            :disabled="!inputContent.trim() || submitting"
            @click="submitComment"
          >
            {{ submitting ? '发送中...' : '发送' }}
          </button>
        </div>
      </div>
    </div>
    <div v-else class="comment-login-hint">
      <router-link to="/login" class="login-link">登录</router-link> 后发表评论
    </div>

    <!-- 评论列表 -->
    <div class="comment-list">
      <div v-for="comment in comments" :key="comment.id" class="comment-item">
        <!-- 一级评论 -->
        <div class="comment-main">
          <img
            v-if="comment.user.avatar"
            :src="getFullUrl(comment.user.avatar)"
            class="comment-avatar"
            @error="$event.target.style.display='none'"
          />
          <div v-else class="comment-avatar comment-avatar--text">
            {{ (comment.user.nickname || comment.user.username || '?').charAt(0).toUpperCase() }}
          </div>
          <div class="comment-body">
            <div class="comment-user">{{ comment.user.nickname || comment.user.username }}</div>
            <div class="comment-time">{{ formatTime(comment.created_at) }}</div>
            <div class="comment-content">{{ comment.content }}</div>
            <div class="comment-actions">
              <button
                class="action-link"
                :class="{ liked: comment.is_liked }"
                @click="onToggleLike(comment.id, false)"
              >
                ▲ {{ comment.like_count || '点赞' }}
              </button>
              <button class="action-link" @click="startReply(comment.id, comment.user.nickname || comment.user.username)">
                回复
              </button>
              <button
                v-if="userStore.userInfo?.id === comment.user.id"
                class="action-link action-delete"
                @click="onDelete(comment.id)"
              >
                删除
              </button>
            </div>
          </div>
        </div>

        <!-- 二级回复 -->
        <div v-if="comment.replies.length > 0" class="reply-list">
          <div
            v-for="reply in comment.replies"
            v-show="expandedReplies.has(comment.id) || comment.replies.indexOf(reply) < 3"
            :key="reply.id"
            class="reply-item"
          >
            <img
              v-if="reply.user.avatar"
              :src="getFullUrl(reply.user.avatar)"
              class="reply-avatar"
              @error="$event.target.style.display='none'"
            />
            <div v-else class="reply-avatar reply-avatar--text">
              {{ (reply.user.nickname || reply.user.username || '?').charAt(0).toUpperCase() }}
            </div>
            <div class="reply-body">
              <span class="reply-user">{{ reply.user.nickname || reply.user.username }}</span>
              <span class="reply-to">
                回复 <span class="reply-to-name">@{{ reply.to_user.nickname || reply.to_user.username }}</span>
              </span>
              <div class="reply-content">{{ reply.content }}</div>
              <div class="reply-actions">
                <button
                  class="action-link"
                  :class="{ liked: reply.is_liked }"
                  @click="onToggleLike(reply.id, true, comment.id)"
                >
                  ▲ {{ reply.like_count || '点赞' }}
                </button>
                <button class="action-link" @click="startReply(reply.id, reply.user.nickname || reply.user.username)">
                  回复
                </button>
                <button
                  v-if="userStore.userInfo?.id === reply.user.id"
                  class="action-link action-delete"
                  @click="onDelete(reply.id)"
                >
                  删除
                </button>
              </div>
            </div>
          </div>
          <button
            v-if="comment.reply_count > 3"
            class="expand-replies"
            @click="toggleReplies(comment.id)"
          >
            {{ expandedReplies.has(comment.id) ? '收起回复' : `展开更多回复 (${comment.reply_count - 3})` }}
          </button>
        </div>
      </div>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore" class="load-more">
      <button class="load-more-btn" :disabled="loading" @click="loadMore">
        {{ loading ? '加载中...' : '加载更多评论' }}
      </button>
    </div>
    <div v-else-if="comments.length > 0" class="no-more">
      没有更多评论了
    </div>
    <div v-else class="empty-comments">
      暂无评论，来抢沙发吧～
    </div>
  </div>
</template>

<style scoped>
.comment-section {
  margin-top: 24px;
  padding: 20px;
  background: #0f1117;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.comment-count {
  font-size: 16px;
  font-weight: 600;
  color: #e4e5eb;
}

.sort-tabs {
  display: flex;
  gap: 16px;
}

.sort-btn {
  background: none;
  border: none;
  color: #5a5d6e;
  font-size: 13px;
  cursor: pointer;
  padding: 4px 0;
  transition: color 0.2s;
}

.sort-btn:hover,
.sort-btn.active {
  color: #00f0ff;
}

/* 输入框 */
.comment-input-area {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
}

.input-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  flex-shrink: 0;
}

.input-box-wrap {
  flex: 1;
}

.reply-hint {
  font-size: 12px;
  color: #5a5d6e;
  margin-bottom: 6px;
}

.reply-name {
  color: #00f0ff;
}

.cancel-reply {
  background: none;
  border: none;
  color: #888;
  font-size: 12px;
  cursor: pointer;
  margin-left: 8px;
}

.cancel-reply:hover {
  color: #e4e5eb;
}

.comment-textarea {
  width: 100%;
  padding: 10px 12px;
  background: #08090d;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: #e4e5eb;
  font-size: 14px;
  resize: vertical;
  outline: none;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.comment-textarea:focus {
  border-color: rgba(0, 240, 255, 0.4);
}

.comment-textarea::placeholder {
  color: #5a5d6e;
}

.input-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.submit-btn {
  padding: 6px 20px;
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background: rgba(0, 240, 255, 0.2);
}

.submit-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.comment-login-hint {
  text-align: center;
  padding: 20px;
  color: #5a5d6e;
  font-size: 14px;
  margin-bottom: 16px;
  background: #08090d;
  border-radius: 4px;
  border: 1px dashed rgba(255, 255, 255, 0.06);
}

.login-link {
  color: #00f0ff;
  text-decoration: none;
}

.login-link:hover {
  text-decoration: underline;
}

/* 评论列表 */
.comment-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  padding-bottom: 16px;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-main {
  display: flex;
  gap: 12px;
}

.comment-avatar,
.reply-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-avatar--text,
.reply-avatar--text {
  background: rgba(0, 240, 255, 0.1);
  border: 1px solid rgba(0, 240, 255, 0.3);
  color: #00f0ff;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.reply-avatar {
  width: 32px;
  height: 32px;
}

.reply-avatar--text {
  font-size: 12px;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-user,
.reply-user {
  font-size: 14px;
  font-weight: 500;
  color: #00f0ff;
}

.comment-time {
  font-size: 12px;
  color: #5a5d6e;
  margin-top: 2px;
}

.comment-content {
  font-size: 14px;
  color: #c8c9d0;
  margin-top: 6px;
  line-height: 1.6;
  word-break: break-all;
}

.comment-actions,
.reply-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.action-link {
  background: none;
  border: none;
  color: #5a5d6e;
  font-size: 12px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
}

.action-link:hover {
  color: #e4e5eb;
}

.action-link.liked {
  color: #00f0ff;
}

.action-delete:hover {
  color: #ff4d4d;
}

/* 回复列表 */
.reply-list {
  margin-left: 52px;
  margin-top: 12px;
  padding-left: 12px;
  border-left: 2px solid rgba(255, 255, 255, 0.04);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.reply-item {
  display: flex;
  gap: 10px;
}

.reply-body {
  flex: 1;
}

.reply-to {
  font-size: 12px;
  color: #5a5d6e;
  margin-left: 4px;
}

.reply-to-name {
  color: #00f0ff;
}

.reply-content {
  font-size: 13px;
  color: #c8c9d0;
  margin-top: 4px;
  line-height: 1.5;
  word-break: break-all;
}

.expand-replies {
  background: none;
  border: none;
  color: #5a5d6e;
  font-size: 12px;
  cursor: pointer;
  padding: 4px 0;
  text-align: left;
}

.expand-replies:hover {
  color: #00f0ff;
}

/* 加载更多 */
.load-more {
  text-align: center;
  margin-top: 20px;
}

.load-more-btn {
  padding: 8px 32px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #5a5d6e;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
}

.load-more-btn:hover:not(:disabled) {
  border-color: rgba(0, 240, 255, 0.3);
  color: #00f0ff;
}

.no-more,
.empty-comments {
  text-align: center;
  padding: 40px 0;
  color: #5a5d6e;
  font-size: 13px;
}
</style>
