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
.page { padding: 48px; max-width: 1400px; margin: 0 auto; min-height: 100vh; }
.header { margin-bottom: 40px; }
h1 { font-size: 36px; font-weight: 700; color: #fff; margin-bottom: 8px; letter-spacing: -0.5px; }
.subtitle { font-size: 15px; color: #9ca3af; font-weight: 400; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 120px 20px; color: #6b7280; }
.spinner { width: 48px; height: 48px; border: 3px solid rgba(99, 102, 241, 0.1); border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 20px; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-state { text-align: center; padding: 100px 40px; background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 2px dashed rgba(99, 102, 241, 0.2); border-radius: 20px; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); }
.empty-icon { font-size: 72px; margin-bottom: 20px; opacity: 0.6; filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.3)); }
.empty-state h3 { font-size: 22px; color: #fff; margin-bottom: 12px; font-weight: 600; }
.empty-state p { color: #9ca3af; font-size: 15px; }

.card { background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.06); border-radius: 20px; overflow: hidden; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
.card:hover { box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5); border-color: rgba(99, 102, 241, 0.15); }
table { width: 100%; border-collapse: collapse; }
th { background: rgba(15, 15, 15, 0.8); padding: 18px 24px; text-align: left; font-weight: 600; color: #9ca3af; font-size: 12px; text-transform: uppercase; letter-spacing: 0.8px; }
td { padding: 24px; border-top: 1px solid rgba(255, 255, 255, 0.04); color: #e5e7eb; }
tbody tr { transition: background 0.2s; }
tbody tr:hover { background: rgba(99, 102, 241, 0.05); }
.bvid { color: #818cf8; font-family: 'SF Mono', 'Monaco', 'Courier New', monospace; }
.badge { display: inline-block; padding: 6px 14px; border-radius: 12px; font-size: 12px; font-weight: 600; }
.badge-success { background: rgba(16, 185, 129, 0.15); color: #34d399; border: 1px solid rgba(16, 185, 129, 0.2); }
</style>
