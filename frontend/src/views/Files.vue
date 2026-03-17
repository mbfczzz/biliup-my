<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>文件管理</h1>
        <p class="subtitle">查看录制的视频文件</p>
      </div>
    </div>

    <div class="streamers-list">
      <div v-for="s in streamers" :key="s.id" class="streamer-section">
        <div class="streamer-header" @click="toggleFiles(s.id)">
          <h3>{{ s.name }}</h3>
          <span class="toggle">{{ expandedIds.includes(s.id) ? '▼' : '▶' }}</span>
        </div>

        <div v-if="expandedIds.includes(s.id)" class="files-container">
          <div v-if="loadingFiles[s.id]" class="loading-small">
            <div class="spinner-sm"></div>
            <span>加载中...</span>
          </div>

          <div v-else-if="!files[s.id]?.length" class="empty-small">
            暂无文件
          </div>

          <div v-else class="files-list">
            <div v-for="f in files[s.id]" :key="f.name" class="file-item">
              <span class="file-icon">📄</span>
              <span class="file-name">{{ f.name }}</span>
              <span class="file-size">{{ formatSize(f.size) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const streamers = ref([])
const files = ref({})
const expandedIds = ref([])
const loadingFiles = ref({})

const load = async () => {
  try {
    const { data } = await api.getStreamers()
    streamers.value = data || []
  } catch (error) {
    console.error('加载失败:', error)
  }
}

const toggleFiles = async (id) => {
  const index = expandedIds.value.indexOf(id)
  if (index > -1) {
    expandedIds.value.splice(index, 1)
  } else {
    expandedIds.value.push(id)
    if (!files.value[id]) {
      await loadFiles(id)
    }
  }
}

const loadFiles = async (id) => {
  try {
    loadingFiles.value[id] = true
    const { data } = await api.getFiles(id)
    files.value[id] = data || []
  } catch (error) {
    console.error('加载文件失败:', error)
    files.value[id] = []
  } finally {
    loadingFiles.value[id] = false
  }
}

const formatSize = (bytes) => {
  if (!bytes) return '-'
  const mb = bytes / 1024 / 1024
  if (mb < 1024) return mb.toFixed(2) + ' MB'
  return (mb / 1024).toFixed(2) + ' GB'
}

onMounted(load)
</script>

<style scoped>
.page { padding: 40px; max-width: 1200px; margin: 0 auto; }
.header { margin-bottom: 32px; }
h1 { font-size: 32px; font-weight: 700; color: #fff; margin-bottom: 8px; }
.subtitle { font-size: 14px; color: #9ca3af; }

.streamers-list { display: flex; flex-direction: column; gap: 16px; }
.streamer-section { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 12px; overflow: hidden; }
.streamer-header { padding: 20px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; transition: background 0.2s; }
.streamer-header:hover { background: #252525; }
.streamer-header h3 { font-size: 16px; color: #fff; margin: 0; }
.toggle { color: #9ca3af; font-size: 12px; }

.files-container { padding: 0 20px 20px; border-top: 1px solid #2a2a2a; }
.loading-small { display: flex; align-items: center; gap: 12px; padding: 20px; color: #9ca3af; }
.spinner-sm { width: 20px; height: 20px; border: 2px solid #2a2a2a; border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-small { padding: 20px; color: #6b7280; text-align: center; }
.files-list { display: flex; flex-direction: column; gap: 8px; padding-top: 16px; }
.file-item { display: flex; align-items: center; gap: 12px; padding: 12px; background: #0f0f0f; border-radius: 8px; }
.file-icon { font-size: 20px; }
.file-name { flex: 1; color: #e0e0e0; font-size: 14px; }
.file-size { color: #9ca3af; font-size: 13px; }
</style>
