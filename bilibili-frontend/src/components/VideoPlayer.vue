<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import DPlayer from 'dplayer'

interface Props {
  videoUrl: string
  coverUrl?: string
  videoId: number
  autoplay?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  autoplay: false,
})

const emit = defineEmits<{
  'timeupdate': [time: number]
}>()

const containerRef = ref<HTMLDivElement | null>(null)
let dp: DPlayer | null = null

// localStorage key
function storageKey(id: number) {
  return `video_progress_${id}`
}

function initPlayer() {
  if (!containerRef.value || !props.videoUrl) return

  // 销毁旧实例
  if (dp) {
    dp.destroy()
    dp = null
  }

  const savedTime = localStorage.getItem(storageKey(props.videoId))
  const startTime = savedTime ? parseFloat(savedTime) : 0

  dp = new DPlayer({
    container: containerRef.value,
    video: {
      url: props.videoUrl,
      pic: props.coverUrl || '',
      type: 'auto',
    },
    theme: '#00f0ff',
    autoplay: props.autoplay,
    lang: 'zh-cn',
    hotkey: true,
    preload: 'metadata',
    screenshot: false,
    mutex: true,
    // 自定义控制条颜色覆盖
    contextmenu: [],
  })

  // 恢复播放进度
  if (startTime > 0) {
    dp.on('loadedmetadata', () => {
      if (dp && startTime < dp.video.duration - 5) {
        dp.seek(startTime)
      }
    })
  }

  // 每 5 秒保存进度
  let saveTimer: ReturnType<typeof setInterval> | null = null
  dp.on('play', () => {
    saveTimer = setInterval(() => {
      if (dp) {
        const current = dp.video.currentTime
        const duration = dp.video.duration
        // 播放到 90% 以上视为看完，清除记录
        if (duration && current / duration > 0.9) {
          localStorage.removeItem(storageKey(props.videoId))
        } else {
          localStorage.setItem(storageKey(props.videoId), String(current))
        }
      }
    }, 5000)
  })

  // 定期触发 timeupdate 事件
  let timeTimer: ReturnType<typeof setInterval> | null = null
  dp.on('play', () => {
    timeTimer = setInterval(() => {
      if (dp) {
        emit('timeupdate', dp.video.currentTime)
      }
    }, 500)
  })

  dp.on('pause', () => {
    if (saveTimer) {
      clearInterval(saveTimer)
      saveTimer = null
    }
    if (timeTimer) {
      clearInterval(timeTimer)
      timeTimer = null
    }
    // 暂停时立即保存
    if (dp) {
      localStorage.setItem(storageKey(props.videoId), String(dp.video.currentTime))
    }
  })

  dp.on('ended', () => {
    if (saveTimer) {
      clearInterval(saveTimer)
      saveTimer = null
    }
    if (timeTimer) {
      clearInterval(timeTimer)
      timeTimer = null
    }
    localStorage.removeItem(storageKey(props.videoId))
  })
}

onMounted(() => {
  initPlayer()
})

onUnmounted(() => {
  if (dp) {
    dp.destroy()
    dp = null
  }
})

// videoUrl 变化时重新初始化
watch(() => props.videoUrl, () => {
  initPlayer()
})
</script>

<template>
  <div ref="containerRef" class="dplayer-wrap" />
</template>

<style scoped>
.dplayer-wrap {
  width: 100%;
  height: 100%;
  background: #0a0a0f;
}

.dplayer-wrap :deep(.dplayer-controller) {
  background: linear-gradient(to top, rgba(0,0,0,0.8), transparent) !important;
}

.dplayer-wrap :deep(.dplayer-bar-wrap) {
  background: rgba(255,255,255,0.15) !important;
}

.dplayer-wrap :deep(.dplayer-bar) {
  background: #00f0ff !important;
  box-shadow: 0 0 6px rgba(0, 240, 255, 0.5) !important;
}

.dplayer-wrap :deep(.dplayer-thumb) {
  background: #00f0ff !important;
  box-shadow: 0 0 8px rgba(0, 240, 255, 0.8) !important;
}

.dplayer-wrap :deep(.dplayer-play-icon),
.dplayer-wrap :deep(.dplayer-volume-icon),
.dplayer-wrap :deep(.dplayer-full-icon),
.dplayer-wrap :deep(.dplayer-setting-icon) {
  color: #e4e5eb !important;
}

.dplayer-wrap :deep(.dplayer-play-icon:hover),
.dplayer-wrap :deep(.dplayer-volume-icon:hover),
.dplayer-wrap :deep(.dplayer-full-icon:hover),
.dplayer-wrap :deep(.dplayer-setting-icon:hover) {
  color: #00f0ff !important;
}

.dplayer-wrap :deep(.dplayer-time) {
  color: #8b8fa3 !important;
  font-family: 'JetBrains Mono', Consolas, monospace !important;
}

.dplayer-wrap :deep(.dplayer-menu) {
  background: #0f1117 !important;
  border: 1px solid rgba(0, 240, 255, 0.2) !important;
}

.dplayer-wrap :deep(.dplayer-menu-item) {
  color: #e4e5eb !important;
}

.dplayer-wrap :deep(.dplayer-menu-item:hover) {
  background: rgba(0, 240, 255, 0.1) !important;
}

.dplayer-wrap :deep(.dplayer-setting-box) {
  background: #0f1117 !important;
  border: 1px solid rgba(0, 240, 255, 0.2) !important;
}

.dplayer-wrap :deep(.dplayer-setting-item) {
  color: #e4e5eb !important;
}

.dplayer-wrap :deep(.dplayer-setting-item:hover) {
  background: rgba(0, 240, 255, 0.1) !important;
}

.dplayer-wrap :deep(.dplayer-setting-ditem) {
  color: #8b8fa3 !important;
}

.dplayer-wrap :deep(.dplayer-setting-ditem:hover) {
  color: #00f0ff !important;
  background: rgba(0, 240, 255, 0.1) !important;
}

.dplayer-wrap :deep(.dplayer-mask) {
  background: rgba(0, 0, 0, 0.4) !important;
}

.dplayer-wrap :deep(.dplayer-loading-icon .di) {
  color: #00f0ff !important;
}

.dplayer-wrap :deep(.dplayer-notice) {
  background: rgba(15, 17, 23, 0.9) !important;
  border: 1px solid rgba(0, 240, 255, 0.3) !important;
  color: #00f0ff !important;
  font-family: 'JetBrains Mono', Consolas, monospace !important;
}
</style>
