<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import NavBar from '../components/NavBar.vue'
import SideBar from '../components/SideBar.vue'
import { uploadVideo, getTranscodeStatus, getVideoDetail, updateVideo } from '../api/video'
import { videoCategories } from '../types/video'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// URL 拼接
const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
const SERVER_BASE = API_BASE.replace(/\/api\/v1\/?$/, '')
function getFullUrl(url: string): string {
  if (!url) return ''
  if (url.startsWith('http')) return url
  return SERVER_BASE + url
}

// 编辑模式
const editVideoId = computed(() => {
  const id = route.params.id
  return id ? Number(id) : 0
})
const isEditMode = computed(() => editVideoId.value > 0)

// ========== 状态 ==========
const selectedFile = ref<File | null>(null)
const selectedCover = ref<File | null>(null)
const coverPreview = ref('')
const isDragging = ref(false)
const uploadProgress = ref(0)
const uploadSpeed = ref('')
const uploadState = ref<'idle' | 'uploading' | 'transcoding' | 'done' | 'error'>('idle')
const transcodeLogs = ref<string[]>([])
const currentLogIndex = ref(0)
const videoId = ref(0)

const form = reactive({
  title: '',
  description: '',
  category: 'cine',
})

const titleCount = computed(() => form.title.length)
const descCount = computed(() => form.description.length)

// 格式色块
const formatColor = computed(() => {
  if (!selectedFile.value) return '#5a5d6e'
  const ext = selectedFile.value.name.split('.').pop()?.toLowerCase()
  const colors: Record<string, string> = {
    mp4: '#00f0ff',
    mov: '#b829dd',
    webm: '#00a8ff',
    mkv: '#f0a500',
  }
  return colors[ext || ''] || '#5a5d6e'
})

// 文件大小格式化
function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(1) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

// ========== 文件处理 ==========
const fileInputRef = ref<HTMLInputElement | null>(null)
const coverInputRef = ref<HTMLInputElement | null>(null)

function onSelectClick() {
  fileInputRef.value?.click()
}

function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) handleFile(file)
}

function handleFile(file: File) {
  const allowed = ['.mp4', '.mov', '.webm', '.mkv']
  const ext = '.' + file.name.split('.').pop()?.toLowerCase()
  if (!allowed.includes(ext)) {
    alert('> ERROR: 不支持的格式，请上传 MP4 / MOV / WEBM / MKV')
    return
  }
  if (file.size > 2 * 1024 * 1024 * 1024) {
    alert('> ERROR: 文件超过 2GB 限制')
    return
  }
  selectedFile.value = file
  uploadState.value = 'idle'
  uploadProgress.value = 0
}

function removeFile() {
  selectedFile.value = null
  uploadState.value = 'idle'
  uploadProgress.value = 0
  if (fileInputRef.value) fileInputRef.value.value = ''
}

// ========== 拖拽上传 ==========
function onDragEnter(e: DragEvent) {
  e.preventDefault()
  isDragging.value = true
}
function onDragOver(e: DragEvent) {
  e.preventDefault()
}
function onDragLeave(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false
}
function onDrop(e: DragEvent) {
  e.preventDefault()
  isDragging.value = false
  const file = e.dataTransfer?.files[0]
  if (file) handleFile(file)
}

// ========== 封面处理 ==========
function onCoverClick() {
  coverInputRef.value?.click()
}

function onCoverChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  const allowed = ['image/jpeg', 'image/png', 'image/webp']
  if (!allowed.includes(file.type)) {
    alert('> ERROR: 封面仅支持 JPG / PNG / WebP')
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    alert('> ERROR: 封面超过 5MB')
    return
  }
  selectedCover.value = file
  const reader = new FileReader()
  reader.onload = () => {
    coverPreview.value = reader.result as string
  }
  reader.readAsDataURL(file)
}

function resetCover() {
  selectedCover.value = null
  coverPreview.value = ''
  if (coverInputRef.value) coverInputRef.value.value = ''
}

// ========== 上传 ==========
let progressTimer: ReturnType<typeof setInterval> | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null

async function onSubmit() {
  if (!isEditMode.value && !selectedFile.value) {
    alert('> ERROR: 请选择视频文件')
    return
  }
  if (!form.title.trim()) {
    alert('> ERROR: 请输入标题')
    return
  }

  // 编辑模式：直接提交更新
  if (isEditMode.value) {
    const formData = new FormData()
    formData.append('title', form.title)
    formData.append('description', form.description)
    formData.append('category', form.category)
    if (selectedFile.value) {
      formData.append('video', selectedFile.value)
    }
    if (selectedCover.value) {
      formData.append('cover', selectedCover.value)
    }
    try {
      uploadState.value = 'uploading'
      await updateVideo(editVideoId.value, formData)
      uploadState.value = 'done'
      setTimeout(() => {
        router.push('/personal')
      }, 1500)
    } catch (err: any) {
      uploadState.value = 'error'
      alert(err.message || '> ERROR: 更新失败')
    }
    return
  }

  uploadState.value = 'uploading'
  uploadProgress.value = 0

  const formData = new FormData()
  formData.append('file', selectedFile.value)
  formData.append('title', form.title)
  formData.append('description', form.description)
  formData.append('category', form.category)
  if (selectedCover.value) {
    formData.append('cover', selectedCover.value)
  }

  try {
    const res = await uploadVideo(formData, (progress: number) => {
      uploadProgress.value = progress
      const speedMBps = selectedFile.value
        ? (selectedFile.value.size / 1024 / 1024 / (Date.now() / 1000)).toFixed(1)
        : '0'
      uploadSpeed.value = `> ${speedMBps} MB/s`
    })
    videoId.value = res.id

    // 上传完成，进入转码阶段
    uploadState.value = 'transcoding'
    startTranscodePolling(res.id)
  } catch (err: any) {
    uploadState.value = 'error'
    alert(err.message || '> ERROR: 上传失败')
  }
}

function startTranscodePolling(id: number) {
  transcodeLogs.value = [
    '[INFO] INIT_TRANSCODE_PIPELINE...',
    '[INFO] DECODING_VIDEO_STREAM...',
    '[INFO] RESCALING_TO_480P...',
    '[INFO] ENCODING_H264...',
    '[INFO] MUXING_MP4_CONTAINER...',
    '[INFO] UPLOADING_TO_MINIO...',
  ]
  currentLogIndex.value = 0

  const logTimer = setInterval(() => {
    if (currentLogIndex.value < transcodeLogs.value.length - 1) {
      currentLogIndex.value++
    }
  }, 800)

  pollTimer = setInterval(async () => {
    try {
      const status = await getTranscodeStatus(id)
      if (status.transcode_status === 2) {
        // 转码完成
        if (pollTimer) clearInterval(pollTimer)
        clearInterval(logTimer)
        currentLogIndex.value = transcodeLogs.value.length - 1
        uploadState.value = 'done'
        setTimeout(() => {
          router.push(`/video/${id}`)
        }, 1500)
      } else if (status.transcode_status === 3) {
        if (pollTimer) clearInterval(pollTimer)
        clearInterval(logTimer)
        uploadState.value = 'error'
      }
    } catch {
      // 轮询失败继续
    }
  }, 3000)
}

onMounted(async () => {
  if (!userStore.isLoggedIn) {
    router.push('/login?redirect=/upload')
    return
  }
  // 编辑模式：加载原视频数据
  if (isEditMode.value) {
    try {
      const video = await getVideoDetail(editVideoId.value)
      form.title = video.title
      form.description = video.description
      form.category = video.category
      if (video.cover_url) {
        coverPreview.value = getFullUrl(video.cover_url)
      }
    } catch {
      alert('加载视频信息失败')
      router.push('/personal')
    }
  }
})

onUnmounted(() => {
  if (progressTimer) clearInterval(progressTimer)
  if (pollTimer) clearInterval(pollTimer)
})
</script>

<template>
  <div class="upload-page">
    <NavBar />
    <SideBar />

    <main class="upload-main">
      <div class="upload-container">
        <!-- 顶部标题 -->
        <div class="upload-header">
          <h1 class="upload-title"
            >// {{ isEditMode ? 'EDIT_PROTOCOL' : 'UPLOAD_PROTOCOL' }} // v1.0</h1>
        </div>

        <!-- 文件选择区 -->
        <div
          v-if="!selectedFile"
          class="drop-zone"
          :class="{ 'drop-zone--active': isDragging }"
          @click="onSelectClick"
          @dragenter="onDragEnter"
          @dragover="onDragOver"
          @dragleave="onDragLeave"
          @drop="onDrop"
        >
          <div class="drop-zone__content">
            <button class="select-btn" @click.stop="onSelectClick"
              >[ SELECT_FILE ]</button>
            <p class="drop-zone__text"
              >{{ isDragging ? '> READY_TO_UPLOAD' : (isEditMode ? '> CLICK_TO_REPLACE_VIDEO (OPTIONAL)' : '> CLICK_TO_UPLOAD') }}</p>
            <p class="drop-zone__hint"
              >SUPPORTED: MP4 / MOV / WEBM / MKV</p>
          </div>
          <input
            ref="fileInputRef"
            type="file"
            accept=".mp4,.mov,.webm,.mkv"
            hidden
            @change="onFileChange"
          />
        </div>

        <!-- 文件信息卡片 -->
        <div v-else class="file-card">
          <div class="file-card__thumb" :style="{ backgroundColor: formatColor }">
            <span class="file-card__ext"
              >{{ selectedFile.name.split('.').pop()?.toUpperCase() }}</span>
          </div>
          <div class="file-card__info">
            <div class="file-card__name"
              >{{ selectedFile.name }}</div>
            <div class="file-card__meta"
              >> {{ formatSize(selectedFile.size) }}</div>
          </div>
          <button class="file-card__remove" @click="removeFile"
            >[X]</button>
        </div>

        <!-- 上传进度 / 转码状态 -->
        <div v-if="uploadState !== 'idle'" class="progress-area">
          <div class="progress-bar">
            <div
              class="progress-fill"
              :class="{
                'progress-fill--pulse': uploadState === 'transcoding',
                'progress-fill--done': uploadState === 'done',
                'progress-fill--error': uploadState === 'error',
              }"
              :style="{ width: uploadProgress + '%' }"
            ></div>
          </div>
          <div class="progress-info">
            <span v-if="uploadState === 'uploading'"
              >> {{ uploadProgress.toFixed(1) }}% // {{ uploadSpeed }}</span>
            <span v-else-if="uploadState === 'transcoding'"
              >> TRANSCODING // 480P...</span>
            <span v-else-if="uploadState === 'done'"
              >> TRANSCODE_COMPLETE // READY_FOR_STREAM</span>
            <span v-else-if="uploadState === 'error'"
              >> TRANSCODE_FAILED // CHECK_LOGS</span>
          </div>

          <!-- 转码日志 -->
          <div v-if="uploadState === 'transcoding'" class="transcode-logs">
            <div
              v-for="(log, i) in transcodeLogs"
              :key="i"
              class="log-line"
              :class="{ 'log-line--active': i === currentLogIndex }"
            >
              {{ log }}
            </div>
          </div>
        </div>

        <!-- 视频信息表单 -->
        <div v-if="selectedFile" class="upload-form">
          <div class="form-group">
            <label class="form-label">> TITLE</label>
            <input
              v-model="form.title"
              type="text"
              class="form-input"
              placeholder="> 输入视频标题..."
              maxlength="200"
            />
            <div class="char-count"
              >{{ titleCount }}/200</div>
          </div>

          <div class="form-group">
            <label class="form-label">> CATEGORY</label>
            <div class="category-list">
              <button
                v-for="cat in videoCategories"
                :key="cat.id"
                class="category-btn"
                :class="{ 'category-btn--active': form.category === cat.id }"
                @click="form.category = cat.id"
              >
                {{ cat.label }}
              </button>
            </div>
          </div>

          <div class="form-group">
            <label class="form-label">> DESCRIPTION</label>
            <textarea
              v-model="form.description"
              class="form-textarea"
              placeholder="> 写点什么吧..."
              maxlength="500"
            ></textarea>
            <div class="char-count"
              >{{ descCount }}/500</div>
          </div>

          <div class="form-group">
            <label class="form-label">> COVER_ART</label>
            <div class="cover-area">
              <div v-if="!coverPreview" class="cover-auto">
                <span class="cover-auto-text"
                  >> AUTO_GENERATE</span>
                <span class="cover-auto-hint"
                  >（后端 FFmpeg 截取第 5 秒）</span>
              </div>
              <img v-else :src="coverPreview" class="cover-preview" />
              <div class="cover-actions">
                <button class="cover-btn" @click="onCoverClick"
                  >> 选择封面</button>
                <button v-if="coverPreview" class="cover-btn cover-btn--reset" @click="resetCover"
                  >> RESET</button>
              </div>
            </div>
            <input
              ref="coverInputRef"
              type="file"
              accept="image/*"
              hidden
              @change="onCoverChange"
            />
          </div>

          <button
            class="submit-btn"
            :disabled="uploadState === 'uploading' || uploadState === 'transcoding'"
            @click="onSubmit"
          >
            <span v-if="uploadState === 'uploading'"
              >> UPLOADING...</span>
            <span v-else-if="uploadState === 'transcoding'"
              >> TRANSCODING...</span>
            <span v-else-if="uploadState === 'done'"
              >> DONE</span>
            <span v-else-if="uploadState === 'error'"
              >> RETRY</span>
            <span v-else>[ {{ isEditMode ? 'EXECUTE_UPDATE' : 'EXECUTE_UPLOAD' }} ]</span>
          </button>
        </div>
      </div>
    </main>
  </div>
</template>

<style scoped>
.upload-page {
  min-height: 100vh;
  background: #08090d;
  background-image: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 240, 255, 0.02) 2px,
    rgba(0, 240, 255, 0.02) 4px
  );
}

.upload-main {
  margin-left: 200px;
  margin-top: 64px;
  min-height: calc(100vh - 64px);
  padding: 32px 24px;
}

.upload-container {
  max-width: 720px;
  margin: 0 auto;
}

/* 顶部标题 */
.upload-header {
  margin-bottom: 32px;
}

.upload-title {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 20px;
  font-weight: 600;
  color: #00f0ff;
  text-shadow: 0 0 8px rgba(0, 240, 255, 0.4);
  margin: 0;
  letter-spacing: 0.04em;
}

.upload-title::after {
  content: '_';
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  50% { opacity: 0; }
}

/* 拖拽区域 */
.drop-zone {
  height: 280px;
  border: 2px dashed rgba(0, 240, 255, 0.2);
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: border-color 0.15s, box-shadow 0.15s;
  margin-bottom: 24px;
}

.drop-zone:hover {
  border-color: rgba(0, 240, 255, 0.4);
}

.drop-zone--active {
  border-style: solid;
  border-color: #00f0ff;
  box-shadow: 0 0 20px rgba(0, 240, 255, 0.2);
  background-image: repeating-linear-gradient(
    0deg,
    transparent,
    transparent 2px,
    rgba(0, 240, 255, 0.03) 2px,
    rgba(0, 240, 255, 0.03) 4px
  );
}

.drop-zone__content {
  text-align: center;
}

.select-btn {
  display: inline-block;
  padding: 12px 32px;
  background: transparent;
  border: 1px solid #b829dd;
  border-radius: 4px;
  color: #fff;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  letter-spacing: 0.04em;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: border-color 0.15s, color 0.15s;
  margin-bottom: 16px;
}

.select-btn:hover {
  border-color: #00f0ff;
  color: #00f0ff;
}

.select-btn:hover::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(0, 240, 255, 0.1),
    transparent
  );
  animation: scanline 0.6s linear;
}

@keyframes scanline {
  to { left: 100%; }
}

.drop-zone__text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  letter-spacing: 0.06em;
  margin: 0 0 8px;
}

.drop-zone__hint {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.04em;
  margin: 0;
}

/* 文件卡片 */
.file-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #0f1117;
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 4px;
  margin-bottom: 24px;
}

.file-card__thumb {
  width: 60px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 2px;
}

.file-card__ext {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #fff;
  font-weight: 600;
}

.file-card__info {
  flex: 1;
  min-width: 0;
}

.file-card__name {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  color: #e4e5eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-card__meta {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-top: 4px;
}

.file-card__remove {
  background: none;
  border: none;
  color: #ff4757;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  cursor: pointer;
  opacity: 0.6;
  transition: opacity 0.15s;
}

.file-card__remove:hover {
  opacity: 1;
  animation: glitch 0.3s linear;
}

@keyframes glitch {
  0%, 100% { transform: translate(0); }
  20% { transform: translate(-2px, 1px); }
  40% { transform: translate(2px, -1px); }
  60% { transform: translate(-1px, -1px); }
  80% { transform: translate(1px, 1px); }
}

/* 进度条 */
.progress-area {
  margin-bottom: 24px;
}

.progress-bar {
  width: 100%;
  height: 4px;
  background: #14161f;
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #b829dd, #00f0ff);
  border-radius: 2px;
  transition: width 0.2s linear;
  position: relative;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 40px;
  background: linear-gradient(90deg, transparent, rgba(255,255,255,0.3), transparent);
  animation: shimmer 1s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(0); }
  100% { transform: translateX(100px); }
}

.progress-fill--pulse {
  background: #00f0ff;
  animation: pulse-opacity 1.5s ease-in-out infinite;
}

@keyframes pulse-opacity {
  0%, 100% { opacity: 0.4; }
  50% { opacity: 1; }
}

.progress-fill--done {
  background: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.5);
}

.progress-fill--error {
  background: #ff4757;
}

.progress-info {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #00f0ff;
  margin-top: 12px;
  letter-spacing: 0.04em;
}

/* 转码日志 */
.transcode-logs {
  margin-top: 16px;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 4px;
}

.log-line {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  line-height: 1.8;
  transition: color 0.15s;
}

.log-line--active {
  color: #00f0ff;
}

/* 表单 */
.upload-form {
  display: flex;
  flex-direction: column;
  gap: 28px;
}

.form-group {
  margin-bottom: 4px;
}

.form-label {
  display: block;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  letter-spacing: 0.08em;
  margin-bottom: 10px;
}

.form-input,
.form-textarea {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid #b829dd;
  border-radius: 4px;
  color: #e4e5eb;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 13px;
  letter-spacing: 0.02em;
  outline: none;
  transition: border-color 0.15s, box-shadow 0.15s;
  padding: 0 16px;
}

.form-input {
  height: 48px;
}

.form-textarea {
  height: 120px;
  padding: 12px 16px;
  resize: none;
}

.form-input:focus,
.form-textarea:focus {
  border-color: #00f0ff;
  box-shadow: 0 0 12px rgba(0, 240, 255, 0.2);
}

.form-input::placeholder,
.form-textarea::placeholder {
  color: #5a5d6e;
}

.char-count {
  text-align: right;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #5a5d6e;
  margin-top: 6px;
}

/* 分类选择 */
.category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.category-btn {
  padding: 8px 20px;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: #5a5d6e;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  letter-spacing: 0.04em;
  cursor: pointer;
  transition: all 0.1s;
}

.category-btn:hover {
  border-color: #8b8fa3;
  color: #8b8fa3;
}

.category-btn--active {
  border-color: #00f0ff;
  background: rgba(0, 240, 255, 0.08);
  color: #00f0ff;
  box-shadow: 0 0 8px rgba(0, 240, 255, 0.15);
}

/* 封面 */
.cover-area {
  display: flex;
  align-items: center;
  gap: 16px;
}

.cover-auto {
  width: 120px;
  height: 72px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px dashed rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.cover-auto-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 10px;
  color: #5a5d6e;
}

.cover-auto-hint {
  font-size: 10px;
  color: #3a3d4e;
}

.cover-preview {
  width: 120px;
  height: 72px;
  object-fit: cover;
  border-radius: 4px;
  border: 1px solid rgba(0, 240, 255, 0.2);
}

.cover-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.cover-btn {
  background: none;
  border: none;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 11px;
  color: #8b8fa3;
  cursor: pointer;
  transition: color 0.15s;
  text-align: left;
}

.cover-btn:hover {
  color: #00f0ff;
}

.cover-btn--reset {
  color: #ff4757;
}

.cover-btn--reset:hover {
  color: #ff6b7a;
}

/* 提交按钮 */
.submit-btn {
  width: 100%;
  height: 52px;
  background: #b829dd;
  border: none;
  border-radius: 4px;
  color: #fff;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 14px;
  letter-spacing: 0.04em;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: background 0.15s, color 0.15s;
  margin-top: 8px;
}

.submit-btn:hover:not(:disabled) {
  background: #00f0ff;
  color: #08090d;
}

.submit-btn:hover:not(:disabled)::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.15),
    transparent
  );
  animation: scanline 0.6s linear;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
