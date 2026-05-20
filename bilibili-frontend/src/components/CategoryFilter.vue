<script setup lang="ts">
import { ref } from 'vue'

interface CategoryItem {
  id: string
  label: string
}

const categories: CategoryItem[] = [
  { id: 'all', label: '全部' },
  { id: 'cine', label: '影视' },
  { id: 'game', label: '游戏' },
  { id: 'acgn', label: '二次元' },
  { id: 'chef', label: '美食' },
  { id: 'docu', label: '纪录片' },
  { id: 'show', label: '综艺' },
]

const activeId = ref('all')

const emit = defineEmits<{
  (e: 'change', categoryId: string): void
}>()

function onSelect(id: string) {
  if (activeId.value === id) return
  activeId.value = id
  emit('change', id)
}
</script>

<template>
  <div class="category-filter">
    <div class="category-list">
      <button
        v-for="cat in categories"
        :key="cat.id"
        class="category-btn"
        :class="{ active: activeId === cat.id }"
        @click="onSelect(cat.id)"
      >
        {{ cat.label }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.category-filter {
  margin: 16px 0;
  padding: 0 4px;
}

.category-list {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.category-list::-webkit-scrollbar {
  display: none;
}

.category-btn {
  flex-shrink: 0;
  padding: 6px 16px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 4px;
  color: #8b8fa3;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  letter-spacing: 0.02em;
}

.category-btn:hover {
  border-color: rgba(0, 240, 255, 0.3);
  color: #e4e5eb;
}

.category-btn.active {
  background: rgba(0, 240, 255, 0.08);
  border-color: rgba(0, 240, 255, 0.4);
  color: #00f0ff;
  box-shadow: 0 0 8px rgba(0, 240, 255, 0.1);
}
</style>
