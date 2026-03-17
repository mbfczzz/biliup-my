<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>视频列表</h1>
        <p class="subtitle">查看已上传的视频</p>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="!videos.length" class="empty-state">
      <div class="empty-icon">📹</div>
      <h3>还没有上传视频</h3>
      <p>录制完成后会自动上传到这里</p>
    </div>

    <div v-else class="card">
      <table>
        <thead>
          <tr>
            <th>视频标题</th>
            <th>BV号</th>
            <th>上传时间</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="v in videos" :key="v.bvid">
            <td><strong>{{ v.title }}</strong></td>
            <td class="bvid">{{ v.bvid }}</td>
            <td>{{ formatTime(v.created) }}</td>
            <td><span class="badge badge-success">已上传</span></td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const videos = ref([])
const loading = ref(true)

const load = async () => {
  try {
    loading.value = true
    const { data } = await api.getVideos()
    videos.value = data || []
  } catch (error) {
    console.error('加载失败:', error)
    videos.value = []
  } finally {
    loading.value = false
  }
}

const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

onMounted(load)
</script>

<style scoped>
.page { padding: 40px; max-width: 1400px; margin: 0 auto; }
.header { margin-bottom: 32px; }
h1 { font-size: 32px; font-weight: 700; color: #fff; margin-bottom: 8px; }
.subtitle { font-size: 14px; color: #9ca3af; }

.loading-state { display: flex; flex-direction: column; align-items: center; padding: 80px 20px; color: #9ca3af; }
.spinner { width: 40px; height: 40px; border: 3px solid #2a2a2a; border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 16px; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-state { text-align: center; padding: 80px 20px; background: #1a1a1a; border: 2px dashed #2a2a2a; border-radius: 16px; }
.empty-icon { font-size: 64px; margin-bottom: 16px; opacity: 0.5; }
.empty-state h3 { font-size: 20px; color: #fff; margin-bottom: 8px; }
.empty-state p { color: #9ca3af; }

.card { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 16px; overflow: hidden; }
table { width: 100%; border-collapse: collapse; }
th { background: #252525; padding: 16px 20px; text-align: left; font-weight: 600; color: #9ca3af; font-size: 12px; text-transform: uppercase; }
td { padding: 20px; border-top: 1px solid #2a2a2a; color: #e0e0e0; }
.bvid { color: #6366f1; font-family: monospace; }
.badge { display: inline-block; padding: 6px 12px; border-radius: 20px; font-size: 12px; font-weight: 600; }
.badge-success { background: rgba(16, 185, 129, 0.2); color: #10b981; }
</style>
