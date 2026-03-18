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
.page { padding: 48px; max-width: 1200px; margin: 0 auto; min-height: 100vh; }
.header { margin-bottom: 40px; }
h1 { font-size: 36px; font-weight: 700; color: #fff; margin-bottom: 8px; letter-spacing: -0.5px; }
.subtitle { font-size: 15px; color: #9ca3af; font-weight: 400; }

.streamers-list { display: flex; flex-direction: column; gap: 20px; }
.streamer-section { background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.06); border-radius: 20px; overflow: hidden; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.streamer-section:hover { box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5); border-color: rgba(99, 102, 241, 0.15); }
.streamer-header { padding: 24px 28px; cursor: pointer; display: flex; justify-content: space-between; align-items: center; transition: background 0.2s; }
.streamer-header:hover { background: rgba(99, 102, 241, 0.05); }
.streamer-header h3 { font-size: 18px; color: #fff; margin: 0; font-weight: 600; letter-spacing: -0.3px; }
.toggle { color: #9ca3af; font-size: 14px; transition: transform 0.2s; }

.files-container { padding: 0 28px 28px; border-top: 1px solid rgba(255, 255, 255, 0.04); }
.loading-small { display: flex; align-items: center; gap: 12px; padding: 24px; color: #9ca3af; }
.spinner-sm { width: 24px; height: 24px; border: 2px solid rgba(99, 102, 241, 0.1); border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-small { padding: 24px; color: #6b7280; text-align: center; font-size: 14px; }
.files-list { display: flex; flex-direction: column; gap: 10px; padding-top: 20px; }
.file-item { display: flex; align-items: center; gap: 14px; padding: 16px; background: rgba(15, 15, 15, 0.6); border: 1px solid rgba(255, 255, 255, 0.04); border-radius: 12px; transition: all 0.2s; }
.file-item:hover { background: rgba(15, 15, 15, 0.8); border-color: rgba(99, 102, 241, 0.2); transform: translateX(4px); }
.file-icon { font-size: 22px; filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3)); }
.file-name { flex: 1; color: #e5e7eb; font-size: 14px; font-weight: 500; }
.file-size { color: #9ca3af; font-size: 13px; font-family: 'SF Mono', 'Monaco', 'Courier New', monospace; }
</style>
