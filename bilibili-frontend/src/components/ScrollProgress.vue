<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const progress = ref(0)

function updateProgress() {
  const scrollTop = window.scrollY
  const docHeight = document.documentElement.scrollHeight - window.innerHeight
  progress.value = docHeight > 0 ? (scrollTop / docHeight) * 100 : 0
}

onMounted(() => {
  window.addEventListener('scroll', updateProgress)
})

onUnmounted(() => {
  window.removeEventListener('scroll', updateProgress)
})
</script>

<template>
  <div class="scroll-progress">
    <div class="scroll-progress__bar" :style="{ height: progress + '%' }" />
  </div>
</template>

<style scoped>
.scroll-progress {
  position: fixed;
  right: 0;
  top: 0;
  width: 2px;
  height: 100vh;
  z-index: 9999;
  pointer-events: none;
}

.scroll-progress__bar {
  width: 100%;
  background: #00f0ff;
  box-shadow: 0 0 8px rgba(0, 240, 255, 0.5);
  transition: height 0.1s linear;
}
</style>
