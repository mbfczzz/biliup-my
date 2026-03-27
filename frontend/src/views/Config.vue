<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">系统配置</h1>
      <p class="page-subtitle">配置录播和上传参数</p>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载配置中...</p></div>

    <template v-else>
      <div class="config-grid">
        <div class="card config-card">
          <div class="config-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M2 6a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1H8a3 3 0 00-3 3v1.5a1.5 1.5 0 01-3 0V6z" clip-rule="evenodd"/><path d="M6 12a2 2 0 012-2h8a2 2 0 012 2v2a2 2 0 01-2 2H2h2a2 2 0 002-2v-2z"/></svg>
          </div>
          <h3 class="config-title">路径配置</h3>
          <div class="config-fields">
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
          </div>
        </div>

        <div class="card config-card">
          <div class="config-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd"/></svg>
          </div>
          <h3 class="config-title">录制配置</h3>
          <div class="config-fields">
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
          </div>
        </div>

        <div class="card config-card">
          <div class="config-icon">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor"><path d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM6.293 6.707a1 1 0 010-1.414l3-3a1 1 0 011.414 0l3 3a1 1 0 01-1.414 1.414L11 5.414V13a1 1 0 11-2 0V5.414L7.707 6.707a1 1 0 01-1.414 0z"/></svg>
          </div>
          <h3 class="config-title">上传配置</h3>
          <div class="config-fields">
            <div class="form-group">
              <label class="checkbox-wrap">
                <input type="checkbox" v-model="config.auto_upload">
                <span>自动上传</span>
              </label>
              <p class="form-hint">录制完成后自动上传到 B 站</p>
            </div>
            <div class="form-group">
              <label class="form-label">并发上传数</label>
              <input class="form-input" type="number" v-model.number="config.concurrent_uploads" min="1" max="5">
              <p class="form-hint">同时上传的文件数量 (1-5)</p>
            </div>
          </div>
        </div>
      </div>

      <div class="action-bar">
        <button class="btn btn-ghost" @click="load">重置</button>
        <button class="btn btn-primary" @click="save">保存配置</button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const loading = ref(true)
const config = ref({ download_path: '', upload_path: '', quality: 'best', format: 'mp4', auto_upload: true, concurrent_uploads: 3 })

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
.config-grid { display:grid; grid-template-columns:repeat(auto-fit, minmax(320px, 1fr)); gap:16px; margin-bottom:24px; }
.config-card { padding:24px; }
.config-icon { width:36px; height:36px; border-radius:10px; background:rgba(99,102,241,.1); color:var(--c-accent); display:flex; align-items:center; justify-content:center; margin-bottom:14px; }
.config-title { font-size:16px; font-weight:700; color:var(--c-text-1); margin-bottom:20px; }
.config-fields { display:flex; flex-direction:column; gap:16px; }
.config-fields .form-group { margin-bottom:0; }

.checkbox-wrap { display:flex; align-items:center; gap:8px; cursor:pointer; font-size:14px; color:var(--c-text-1); font-weight:500; }
.checkbox-wrap input[type="checkbox"] { width:18px; height:18px; accent-color:var(--c-accent); cursor:pointer; }

.action-bar { display:flex; gap:10px; justify-content:flex-end; }

.form-select {
  appearance:none;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12' fill='%2371717a'%3E%3Cpath d='M3 5l3 3 3-3'/%3E%3C/svg%3E");
  background-repeat:no-repeat;
  background-position:right 12px center;
  padding-right:32px;
}
</style>
