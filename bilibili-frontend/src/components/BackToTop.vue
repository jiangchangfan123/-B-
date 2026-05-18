<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const visible = ref(false)

function checkScroll() {
  visible.value = window.scrollY > 400
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => {
  window.addEventListener('scroll', checkScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', checkScroll)
})
</script>

<template>
  <button
    class="back-to-top"
    :class="{ 'back-to-top--visible': visible }"
    @click="scrollToTop"
  >
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <polyline points="18 15 12 9 6 15" />
    </svg>
  </button>
</template>

<style scoped>
.back-to-top {
  position: fixed;
  bottom: 40px;
  right: 40px;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid #00f0ff;
  border-radius: 4px;
  color: #00f0ff;
  cursor: pointer;
  opacity: 0;
  transform: translateY(10px);
  transition: opacity 0.3s, transform 0.3s, border-color 0.2s, color 0.2s;
  z-index: 800;
}

.back-to-top--visible {
  opacity: 1;
  transform: translateY(0);
}

.back-to-top:hover {
  border-color: #b829dd;
  color: #b829dd;
}

.back-to-top:hover svg {
  animation: arrow-up 0.3s ease forwards;
}

@keyframes arrow-up {
  0% { transform: translateY(0); }
  50% { transform: translateY(-2px); opacity: 0.6; }
  100% { transform: translateY(0); }
}
</style>
