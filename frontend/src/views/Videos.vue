<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">视频列表</h1>
      <p class="page-subtitle">查看录制的视频文件</p>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载中...</p></div>

    <div v-else-if="!videos.length" class="empty-box">
      <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg></div>
      <h3 class="empty-title">还没有视频</h3>
      <p class="empty-desc">录制完成后的视频文件会显示在这里</p>
    </div>

    <div v-else class="card">
      <div class="table-wrap">
        <table>
          <thead><tr><th>文件名</th><th>大小</th><th>修改时间</th></tr></thead>
          <tbody>
            <tr v-for="v in videos" :key="v.key">
              <td><span class="cell-name">{{ v.name }}</span></td>
              <td><span class="cell-size">{{ formatSize(v.size) }}</span></td>
              <td><span class="cell-time">{{ formatTime(v.updateTime) }}</span></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const videos = ref([])
const loading = ref(true)

const load = async () => {
  try { loading.value = true; const { data } = await api.getVideos(); videos.value = data || [] }
  catch { videos.value = [] }
  finally { loading.value = false }
}

const formatTime = (ts) => {
  if (!ts) return '-'
  return new Date(ts * 1000).toLocaleString('zh-CN')
}

const formatSize = (bytes) => {
  if (!bytes) return '-'
  const mb = bytes / 1024 / 1024
  return mb < 1024 ? mb.toFixed(1) + ' MB' : (mb / 1024).toFixed(2) + ' GB'
}

onMounted(load)
</script>

<style scoped>
.cell-name { font-weight:600; color:var(--c-text-1); font-size:13px; }
.cell-size { font-family:var(--font-mono); font-size:13px; color:var(--c-text-3); }
.cell-time { font-size:13px; color:var(--c-text-3); }
</style>
