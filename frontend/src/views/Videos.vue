<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">投稿记录</h1>
        <p class="page-subtitle">查看已上传到 B 站的视频及投稿状态</p>
      </div>
      <div v-if="videos.length" class="header-stat">
        <span class="header-stat-value">{{ videos.length }}</span>
        <span class="header-stat-label">个视频</span>
      </div>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载中...</p></div>

    <div v-else-if="!videos.length" class="empty-box">
      <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg></div>
      <h3 class="empty-title">还没有视频</h3>
      <p class="empty-desc">录制完成后的视频会自动显示在这里</p>
    </div>

    <div v-else class="video-list">
      <div v-for="v in videos" :key="v.key" class="video-item card">
        <div class="video-icon">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg>
        </div>
        <div class="video-info">
          <div class="video-name">{{ v.name }}</div>
          <div class="video-meta">
            <span v-if="v.size" class="meta-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M7 10l5 5 5-5M12 15V3"/></svg>
              {{ formatSize(v.size) }}
            </span>
            <span v-if="v.updateTime" class="meta-item">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><path d="M12 6v6l4 2"/></svg>
              {{ formatTime(v.updateTime) }}
            </span>
          </div>
        </div>
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
.header-stat {
  display:flex; align-items:baseline; gap:4px;
  padding:8px 16px; border-radius:10px;
  background:var(--c-accent-soft);
}
.header-stat-value { font-size:18px; font-weight:800; color:var(--c-accent); letter-spacing:-.3px; }
.header-stat-label { font-size:12px; color:var(--c-accent); font-weight:500; }

.video-list { display:flex; flex-direction:column; gap:8px; }

.video-item {
  display:flex; align-items:center; gap:16px;
  padding:16px 20px;
  transition:all .2s var(--ease);
}
.video-item:hover { transform:translateX(4px); }

.video-icon {
  width:42px; height:42px; border-radius:10px; flex-shrink:0;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
}

.video-info { flex:1; min-width:0; }
.video-name {
  font-size:14px; font-weight:600; color:var(--c-text-1); margin-bottom:4px;
  overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
}
.video-meta { display:flex; gap:16px; }
.meta-item {
  display:inline-flex; align-items:center; gap:5px;
  font-size:12px; color:var(--c-text-4); font-family:var(--font-mono);
}
.meta-item svg { opacity:.6; }
</style>
