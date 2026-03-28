<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">文件管理</h1>
      <p class="page-subtitle">按频道查看录制的文件</p>
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
            <svg :class="['chevron', { open: expanded.has(s.id) }]" width="16" height="16" viewBox="0 0 16 16" fill="currentColor"><path d="M6 4l4 4-4 4"/></svg>
            <span class="section-name">{{ s.name }}</span>
          </div>
          <span v-if="files[s.id]" class="section-count">{{ files[s.id].length }} 个文件</span>
        </div>
        <div v-if="expanded.has(s.id)" class="section-body">
          <div v-if="loadingFiles[s.id]" class="loading-inline"><div class="spinner" style="width:20px;height:20px;border-width:2px;margin-bottom:0"></div><span>加载中...</span></div>
          <div v-else-if="!files[s.id]?.length" class="no-files">暂无文件</div>
          <div v-else class="file-list">
            <div v-for="f in files[s.id]" :key="f.file" class="file-row">
              <svg class="file-icon" width="16" height="16" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd"/></svg>
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
.sections { display:flex; flex-direction:column; gap:8px; }
.section-head { padding:14px 20px; cursor:pointer; display:flex; justify-content:space-between; align-items:center; transition:background .1s var(--ease); }
.section-head:hover { background:rgba(255,255,255,.02); }
.section-left { display:flex; align-items:center; gap:8px; }
.chevron { color:var(--c-text-3); transition:transform .15s var(--ease); }
.chevron.open { transform:rotate(90deg); }
.section-name { font-size:15px; font-weight:600; color:var(--c-text-1); }
.section-count { font-size:12px; color:var(--c-text-3); font-family:var(--font-mono); }
.section-body { padding:0 20px 16px; border-top:1px solid var(--c-border); }
.loading-inline { display:flex; align-items:center; gap:10px; padding:16px 0; color:var(--c-text-3); font-size:13px; }
.no-files { padding:16px 0; color:var(--c-text-3); font-size:13px; text-align:center; }
.file-list { display:flex; flex-direction:column; gap:2px; padding-top:12px; }
.file-row { display:flex; align-items:center; gap:10px; padding:10px 12px; border-radius:8px; transition:background .1s var(--ease); }
.file-row:hover { background:rgba(255,255,255,.03); }
.file-icon { color:var(--c-text-3); flex-shrink:0; }
.file-name { font-size:13px; color:var(--c-text-2); font-family:var(--font-mono); word-break:break-all; }
</style>
