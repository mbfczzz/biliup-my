<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <p class="page-subtitle">配置下载路径、录制参数和上传行为</p>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载配置中...</p></div>

    <template v-else>
      <div class="config-sections">
        <div v-for="section in sections" :key="section.key" :class="['config-section card', { 'is-open': expanded.has(section.key) }]">
          <div class="section-head" @click="toggle(section.key)">
            <div class="section-left">
              <div class="section-icon" v-html="section.icon"></div>
              <div class="section-text">
                <span class="section-title">{{ section.title }}</span>
                <span class="section-desc">{{ section.desc }}</span>
              </div>
            </div>
            <svg :class="['chevron', { open: expanded.has(section.key) }]" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M6 4l4 4-4 4"/></svg>
          </div>
          <div class="section-body">
            <div class="section-body-inner">
              <div class="section-body-content">
              <!-- 路径配置 -->
              <template v-if="section.key === 'path'">
                <div class="form-group">
                  <label class="form-label">下载保存路径</label>
                  <input class="form-input" v-model="config.download_path" placeholder="/path/to/downloads">
                  <p class="form-hint">录制的视频文件保存位置</p>
                </div>
                <div class="form-group">
                  <label class="form-label">上传临时路径</label>
                  <input class="form-input" v-model="config.upload_path" placeholder="/path/to/uploads">
                  <p class="form-hint">上传前的临时文件存储位置</p>
                </div>
              </template>
              <!-- 录制配置 -->
              <template v-if="section.key === 'record'">
                <div class="form-group">
                  <label class="form-label">视频质量</label>
                  <select class="form-select" v-model="config.quality">
                    <option value="best">最佳质量</option>
                    <option value="high">高质量</option>
                    <option value="medium">中等质量</option>
                    <option value="low">低质量</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label">文件格式</label>
                  <select class="form-select" v-model="config.format">
                    <option value="mp4">MP4</option>
                    <option value="flv">FLV</option>
                    <option value="mkv">MKV</option>
                  </select>
                </div>
              </template>
              <!-- 上传配置 -->
              <template v-if="section.key === 'upload'">
                <div class="form-group">
                  <label class="toggle-wrap">
                    <span class="toggle-text">
                      <span class="toggle-label">自动上传</span>
                      <span class="toggle-hint">录制完成后自动上传到 B 站</span>
                    </span>
                    <label class="toggle">
                      <input type="checkbox" v-model="config.auto_upload">
                      <span class="toggle-track"></span>
                    </label>
                  </label>
                </div>
                <div class="form-group">
                  <label class="form-label">并发上传数</label>
                  <input class="form-input" type="number" v-model.number="config.concurrent_uploads" min="1" max="5">
                  <p class="form-hint">同时上传的文件数量 (1-5)</p>
                </div>
              </template>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="action-bar">
        <button class="btn btn-ghost" @click="load">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M1 4v6h6"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
          重置
        </button>
        <button class="btn btn-primary" @click="save">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M20 6L9 17l-5-5"/></svg>
          保存配置
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const loading = ref(true)
const config = ref({ download_path: '', upload_path: '', quality: 'best', format: 'mp4', auto_upload: true, concurrent_uploads: 3 })
const expanded = ref(new Set())

const sections = [
  { key: 'path', title: '路径配置', desc: '设置文件下载和临时存储路径', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>' },
  { key: 'record', title: '录制配置', desc: '视频质量和文件格式偏好', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>' },
  { key: 'upload', title: '上传配置', desc: '自动上传和并发数设置', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M12 4v12M8 8l4-4 4 4"/></svg>' },
]

const toggle = (key) => {
  const s = new Set(expanded.value)
  if (s.has(key)) s.delete(key); else s.add(key)
  expanded.value = s
}

const load = async () => {
  try { loading.value = true; const { data } = await api.getConfig(); if (data) config.value = { ...config.value, ...data } }
  catch (e) { console.error('加载配置失败:', e) }
  finally { loading.value = false }
}

const save = async () => {
  try { await api.updateConfig(config.value); alert('配置保存成功') }
  catch (e) { alert('保存失败: ' + (e.response?.data?.message || e.message)) }
}

onMounted(load)
</script>

<style scoped>
.config-sections { display:flex; flex-direction:column; gap:10px; margin-bottom:24px; }
.config-section { padding:0; overflow:hidden; }

.section-head {
  padding:20px 24px; cursor:pointer;
  display:flex; justify-content:space-between; align-items:center;
  transition:background .15s;
}
.section-head:hover { background:rgba(0,0,0,.015); }

.section-left { display:flex; align-items:center; gap:14px; }
.section-icon {
  width:40px; height:40px; border-radius:12px;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
  flex-shrink:0;
}
.section-text { display:flex; flex-direction:column; gap:2px; }
.section-title { font-size:14px; font-weight:700; color:var(--c-text-1); }
.section-desc { font-size:12px; color:var(--c-text-4); }

.chevron { color:var(--c-text-4); transition:transform .2s var(--ease); flex-shrink:0; }
.chevron.open { transform:rotate(90deg); }

.section-body .form-group { margin-bottom:0; }

/* Toggle switch */
.toggle-wrap {
  display:flex; align-items:center; justify-content:space-between;
  cursor:pointer; gap:16px;
}
.toggle-text { display:flex; flex-direction:column; gap:2px; }
.toggle-label { font-size:14px; font-weight:600; color:var(--c-text-1); }
.toggle-hint { font-size:12px; color:var(--c-text-4); }

.toggle { position:relative; display:inline-block; width:44px; height:24px; flex-shrink:0; }
.toggle input { opacity:0; width:0; height:0; }
.toggle-track {
  position:absolute; inset:0; border-radius:12px;
  background:rgba(0,0,0,.12); cursor:pointer;
  transition:all .25s var(--ease);
}
.toggle-track::before {
  content:''; position:absolute; left:2px; top:2px;
  width:20px; height:20px; border-radius:50%;
  background:#fff; box-shadow:0 1px 3px rgba(0,0,0,.15);
  transition:all .25s var(--ease);
}
.toggle input:checked + .toggle-track {
  background:var(--c-accent);
}
.toggle input:checked + .toggle-track::before {
  transform:translateX(20px);
}

/* Smooth expand with grid */
.section-body {
  display:grid; grid-template-rows:0fr;
  transition:grid-template-rows .3s var(--ease);
}
.is-open .section-body { grid-template-rows:1fr; }
.section-body-inner { overflow:hidden; }
.section-body-content {
  padding:20px 24px 24px;
  display:flex; flex-direction:column; gap:18px;
  border-top:1px solid var(--c-border);
  margin:0 16px;
}

.action-bar { display:flex; gap:10px; justify-content:flex-end; }

.form-select {
  appearance:none;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12' fill='%23999'%3E%3Cpath d='M3 5l3 3 3-3'/%3E%3C/svg%3E");
  background-repeat:no-repeat;
  background-position:right 14px center;
  padding-right:36px;
}
</style>
