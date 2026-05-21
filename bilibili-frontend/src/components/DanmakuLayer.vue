<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { getDanmaku } from '../api/danmaku'
import type { DanmakuItem } from '../api/danmaku'

const props = defineProps<{
  videoId: number
  currentTime: number
  isPlaying: boolean
}>()

// ========== WebSocket ==========
const ws = ref<WebSocket | null>(null)
const reconnectTimer = ref<ReturnType<typeof setTimeout> | null>(null)
const reconnectDelay = 3000
const maxReconnectAttempts = 5
const reconnectAttempts = ref(0)

const apiBase = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const wsBase = apiBase
  .replace(/^http/, 'ws')
  .replace(/\/api\/v1\/?$/, '')
  .replace(/\/$/, '')
const wsUrl = computed(() => `${wsBase}/ws/danmaku?video_id=${props.videoId}`)

// ========== 弹幕显示 ==========
const containerRef = ref<HTMLDivElement | null>(null)
const maxDanmakuCount = 100

// ========== 弹幕池：存储已加载的所有弹幕 ==========
const danmakuPool = ref<DanmakuItem[]>([])
const shownIds = ref<Set<number>>(new Set())

const currentSecond = computed(() => Math.floor(props.currentTime))
let lastSecond = -1

// ========== WebSocket 连接 ==========
function connectWS() {
  if (reconnectAttempts.value >= maxReconnectAttempts) {
    console.warn('[Danmaku] 已超过最大重连次数')
    return
  }
  if (ws.value) ws.value.close()

  const url = wsUrl.value
  console.log('[Danmaku] 正在连接 WebSocket:', url)

  try {
    ws.value = new WebSocket(url)
  } catch (err) {
    console.error('[Danmaku] 创建 WebSocket 失败:', err)
    scheduleReconnect()
    return
  }

  ws.value.onopen = () => {
    console.log('[Danmaku] WebSocket 已连接')
    reconnectAttempts.value = 0
  }

  ws.value.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (msg.type === 'danmaku' && msg.data) {
        const d = msg.data as DanmakuItem
        // 加入弹幕池
        danmakuPool.value.push(d)
        // 如果当前时间在弹幕时间附近（±1秒），立即显示
        if (Math.abs(d.time_point - currentSecond.value) <= 1 && !shownIds.value.has(d.id)) {
          addDanmakuToScreen(d)
          shownIds.value.add(d.id)
        }
      }
    } catch (err) {
      console.error('[Danmaku] 消息解析失败:', err)
    }
  }

  ws.value.onclose = () => {
    console.log('[Danmaku] WebSocket 已断开')
    scheduleReconnect()
  }

  ws.value.onerror = (err) => {
    console.error('[Danmaku] WebSocket 错误:', err)
    ws.value?.close()
  }
}

function scheduleReconnect() {
  if (reconnectTimer.value) return
  if (reconnectAttempts.value >= maxReconnectAttempts) return
  reconnectAttempts.value++
  reconnectTimer.value = setTimeout(() => {
    reconnectTimer.value = null
    connectWS()
  }, reconnectDelay)
}

// ========== 弹幕渲染（直接 DOM 操作） ==========
function addDanmakuToScreen(danmaku: DanmakuItem) {
  if (!containerRef.value) return
  if (shownIds.value.has(danmaku.id)) return  // 防止重复显示

  // 限制同屏数量
  while (containerRef.value.children.length >= maxDanmakuCount) {
    containerRef.value.firstChild?.remove()
  }

  const el = document.createElement('div')
  el.className = 'danmaku-item'
  el.textContent = danmaku.content
  el.style.color = danmaku.color || '#ffffff'

  // 轨道位置：5%~41%，共 10 条轨道
  const trackIndex = Math.floor(Math.random() * 10)
  const top = 5 + trackIndex * 4
  el.style.top = `${top}%`

  // 随机动画时长 6~10 秒
  const duration = 6 + Math.random() * 4
  el.style.setProperty('--duration', `${duration}s`)

  // 如果当前处于暂停状态，创建后立即暂停动画
  if (!props.isPlaying) {
    el.style.animationPlayState = 'paused'
  }

  containerRef.value.appendChild(el)

  // 动画结束后移除 DOM
  el.addEventListener('animationend', () => {
    el.remove()
  })
}

// 暴露方法供父组件调用（发送弹幕后本地立即显示）
defineExpose({
  addDanmaku: (d: DanmakuItem) => {
    danmakuPool.value.push(d)
    if (!shownIds.value.has(d.id)) {
      addDanmakuToScreen(d)
      shownIds.value.add(d.id)
    }
  }
})

// ========== 加载弹幕 ==========
async function loadAllDanmaku() {
  try {
    // 加载该视频前 2 小时的所有弹幕（覆盖绝大多数视频）
    const res = await getDanmaku(props.videoId, 0, 7200)
    danmakuPool.value = res.list
    console.log(`[Danmaku] 已加载 ${res.list.length} 条弹幕`)
    // 立即显示当前时间附近的弹幕
    showDanmakuAtTime(currentSecond.value)
  } catch (err) {
    console.error('[Danmaku] 加载弹幕失败:', err)
  }
}

// 显示指定秒数对应的弹幕（time_point 在 [second, second+1) 范围内）
function showDanmakuAtTime(second: number) {
  const candidates = danmakuPool.value.filter(d =>
    d.time_point >= second &&
    d.time_point < second + 1 &&
    !shownIds.value.has(d.id)
  )
  candidates.forEach(d => {
    addDanmakuToScreen(d)
    shownIds.value.add(d.id)
  })
}

// 清空当前屏幕所有弹幕
function clearScreenDanmaku() {
  if (!containerRef.value) return
  // 移除所有子元素
  while (containerRef.value.firstChild) {
    containerRef.value.firstChild.remove()
  }
}

// 播放进度变化时，显示对应时间点的弹幕
watch(currentSecond, (second, oldSecond) => {
  if (second < 0) return

  // 检测进度条跳（用户拖动进度条或重新播放）
  // 跳跃超过 2 秒认为是人为操作，清空屏幕并重新显示
  const jump = Math.abs(second - lastSecond)
  if (lastSecond >= 0 && jump > 2) {
    console.log(`[Danmaku] 检测到进度跳转: ${lastSecond}s -> ${second}s, 清空屏幕重新加载`)
    clearScreenDanmaku()
    shownIds.value.clear()
    // 显示新时间点的弹幕
    showDanmakuAtTime(second)
  } else {
    // 正常播放，显示当前秒的弹幕
    showDanmakuAtTime(second)
  }
  lastSecond = second
})

// 监听播放/暂停状态：控制所有弹幕动画的播放/暂停
watch(() => props.isPlaying, (playing) => {
  if (!containerRef.value) return
  const items = containerRef.value.querySelectorAll('.danmaku-item')
  items.forEach((el) => {
    ;(el as HTMLElement).style.animationPlayState = playing ? 'running' : 'paused'
  })
})

// ========== 生命周期 ==========
onMounted(() => {
  connectWS()
  loadAllDanmaku()
})

onUnmounted(() => {
  if (reconnectTimer.value) clearTimeout(reconnectTimer.value)
  if (ws.value) ws.value.close()
})
</script>

<template>
  <div class="danmaku-layer">
    <!-- 弹幕滚动区域 -->
    <div class="danmaku-container" ref="containerRef"></div>
  </div>
</template>

<style scoped>
.danmaku-layer {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
  z-index: 10;
}

.danmaku-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 55%;
  overflow: hidden;
}

/*
 * 使用 :deep() 确保动态创建的弹幕元素能获得样式。
 * Vue scoped CSS 不会自动给 document.createElement 创建的元素添加 data-v 属性，
 * :deep() 会编译为 [data-v-xxxxx] .danmaku-item ，通过父元素匹配后代。
 */
:deep(.danmaku-item) {
  position: absolute;
  left: 100%;
  white-space: nowrap;
  font-size: 18px;
  font-weight: 500;
  color: #ffffff;
  text-shadow:
    0 0 2px rgba(0, 0, 0, 0.9),
    0 0 4px rgba(0, 0, 0, 0.7),
    1px 1px 2px rgba(0, 0, 0, 0.8);
  pointer-events: none;
  will-change: left;
  animation-name: danmaku-move;
  animation-duration: var(--duration, 8s);
  animation-timing-function: linear;
  animation-fill-mode: forwards;
}

@keyframes danmaku-move {
  from {
    left: 100%;
  }
  to {
    left: -100%;
  }
}
</style>
