<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">本地文件</h1>
      <p class="page-subtitle">浏览和管理本地录制的视频文件</p>
    </div>

    <div v-if="!streamers.length" class="empty-box">
      <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg></div>
      <h3 class="empty-title">暂无频道</h3>
      <p class="empty-desc">添加录制频道后即可在此查看文件</p>
    </div>

    <div v-else class="sections">
      <div v-for="s in streamers" :key="s.id" class="section card">
        <div class="section-head" @click="toggle(s.id)">
          <div class="section-left">
            <div class="section-avatar">{{ (s.name || '?')[0].toUpperCase() }}</div>
            <div class="section-info">
              <span class="section-name">{{ s.name }}</span>
              <span v-if="files[s.id]" class="section-count">{{ files[s.id].length }} 个文件</span>
            </div>
          </div>
          <svg :class="['chevron', { open: expanded.has(s.id) }]" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M6 4l4 4-4 4"/></svg>
        </div>
        <div v-if="expanded.has(s.id)" class="section-body">
          <div v-if="loadingFiles[s.id]" class="loading-inline"><div class="spinner" style="width:20px;height:20px;border-width:2px;margin-bottom:0"></div><span>加载中...</span></div>
          <div v-else-if="!files[s.id]?.length" class="no-files">暂无录制文件</div>
          <div v-else class="file-list">
            <div v-for="f in files[s.id]" :key="f.file" class="file-row">
              <div class="file-icon-wrap">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><path d="M14 2v6h6"/></svg>
              </div>
              <span class="file-name">{{ f.file }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../api'

const streamers = ref([])
const files = reactive({})
const expanded = ref(new Set())
const loadingFiles = reactive({})

const load = async () => {
  try { const { data } = await api.getStreamers(); streamers.value = data || [] }
  catch { streamers.value = [] }
}

const toggle = async (id) => {
  if (expanded.value.has(id)) { expanded.value.delete(id); expanded.value = new Set(expanded.value) }
  else {
    expanded.value.add(id); expanded.value = new Set(expanded.value)
    if (!files[id]) await loadFiles(id)
  }
}

const loadFiles = async (id) => {
  try { loadingFiles[id] = true; const { data } = await api.getFiles(id); files[id] = data || [] }
  catch { files[id] = [] }
  finally { loadingFiles[id] = false }
}

onMounted(load)
</script>

<style scoped>
.sections { display:flex; flex-direction:column; gap:10px; }

.section-head {
  padding:16px 22px; cursor:pointer;
  display:flex; justify-content:space-between; align-items:center;
  transition:background .15s var(--ease);
}
.section-head:hover { background:rgba(0,0,0,.015); }

.section-left { display:flex; align-items:center; gap:14px; }
.section-avatar {
  width:36px; height:36px; border-radius:10px;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
  font-size:14px; font-weight:700; flex-shrink:0;
}
.section-info { display:flex; flex-direction:column; gap:2px; }
.section-name { font-size:14px; font-weight:600; color:var(--c-text-1); }
.section-count { font-size:11px; color:var(--c-text-4); font-weight:500; }

.chevron { color:var(--c-text-4); transition:transform .2s var(--ease); flex-shrink:0; }
.chevron.open { transform:rotate(90deg); }

.section-body { padding:0 22px 18px; }
.loading-inline { display:flex; align-items:center; gap:10px; padding:16px 0; color:var(--c-text-4); font-size:13px; }
.no-files { padding:20px 0; color:var(--c-text-4); font-size:13px; text-align:center; }

.file-list {
  display:flex; flex-direction:column; gap:2px;
  padding-top:8px; border-top:1px solid var(--c-border);
}
.file-row {
  display:flex; align-items:center; gap:12px;
  padding:10px 12px; border-radius:8px;
  transition:all .15s var(--ease);
}
.file-row:hover { background:var(--c-accent-soft); }
.file-icon-wrap { color:var(--c-text-4); flex-shrink:0; display:flex; }
.file-name {
  font-size:13px; color:var(--c-text-2); font-family:var(--font-mono);
  word-break:break-all;
}
</style>
