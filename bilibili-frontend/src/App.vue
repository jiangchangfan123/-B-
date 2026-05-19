<script setup lang="ts">
import { watch } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from './stores/user'

const route = useRoute()
const userStore = useUserStore()

// 页面刷新时恢复登录态
userStore.restore()
if (userStore.token && !userStore.userInfo) {
  userStore.fetchUserInfo().catch(() => {})
}

// 路由变化时自动检查登录态
watch(() => route.path, () => {
  if (userStore.token && !userStore.userInfo) {
    userStore.fetchUserInfo().catch(() => {})
  }
})
</script>

<template>
  <router-view />
</template>
