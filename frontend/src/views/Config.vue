<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>系统配置</h1>
        <p class="subtitle">配置录播和上传参数</p>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载配置中...</p>
    </div>

    <!-- 配置表单 -->
    <div v-else class="config-container">
      <div class="card">
        <div class="card-header">
          <h3>📁 路径配置</h3>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label>下载保存路径</label>
            <input
              v-model="config.download_path"
              placeholder="/path/to/downloads"
            >
            <div class="hint">录制的视频文件保存位置</div>
          </div>
          <div class="form-group">
            <label>上传临时路径</label>
            <input
              v-model="config.upload_path"
              placeholder="/path/to/uploads"
            >
            <div class="hint">上传前的临时文件存储位置</div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3>⚙️ 录制配置</h3>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label>视频质量</label>
            <select v-model="config.quality">
              <option value="best">最佳质量</option>
              <option value="high">高质量</option>
              <option value="medium">中等质量</option>
              <option value="low">低质量</option>
            </select>
            <div class="hint">录制视频的画质设置</div>
          </div>
          <div class="form-group">
            <label>文件格式</label>
            <select v-model="config.format">
              <option value="mp4">MP4</option>
              <option value="flv">FLV</option>
              <option value="mkv">MKV</option>
            </select>
            <div class="hint">录制视频的文件格式</div>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <h3>🚀 上传配置</h3>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label>
              <input type="checkbox" v-model="config.auto_upload">
              自动上传
            </label>
            <div class="hint">录制完成后自动上传到平台</div>
          </div>
          <div class="form-group">
            <label>并发上传数</label>
            <input
              type="number"
              v-model.number="config.concurrent_uploads"
              min="1"
              max="5"
            >
            <div class="hint">同时上传的文件数量（1-5）</div>
          </div>
        </div>
      </div>

      <div class="actions">
        <button class="btn-secondary" @click="load">重置</button>
        <button class="btn-primary" @click="save">
          💾 保存配置
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const loading = ref(true)
const config = ref({
  download_path: '',
  upload_path: '',
  quality: 'best',
  format: 'mp4',
  auto_upload: true,
  concurrent_uploads: 3
})

const load = async () => {
  try {
    loading.value = true
    const { data } = await api.getConfig()
    if (data) {
      config.value = { ...config.value, ...data }
    }
  } catch (error) {
    console.error('加载配置失败:', error)
  } finally {
    loading.value = false
  }
}

const save = async () => {
  try {
    await api.updateConfig(config.value)
    alert('✅ 配置保存成功')
  } catch (error) {
    alert('❌ 保存失败: ' + (error.response?.data?.message || error.message))
  }
}

onMounted(load)
</script>

<style scoped>
.page {
  padding: 40px;
  max-width: 900px;
  margin: 0 auto;
}

.header {
  margin-bottom: 32px;
}

h1 {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
}

.subtitle {
  font-size: 14px;
  color: #9ca3af;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  color: #9ca3af;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #2a2a2a;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 配置容器 */
.config-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 卡片 */
.card {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

.card-header {
  padding: 20px 24px;
  border-bottom: 1px solid #2a2a2a;
  background: #252525;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.card-body {
  padding: 24px;
}

/* 表单 */
.form-group {
  margin-bottom: 24px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #e0e0e0;
  font-size: 14px;
  font-weight: 600;
}

.form-group input[type="text"],
.form-group input[type="number"],
.form-group select {
  width: 100%;
  padding: 14px 16px;
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  color: #e0e0e0;
  font-size: 14px;
  transition: all 0.2s;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.form-group input[type="checkbox"] {
  width: auto;
  margin-right: 8px;
  cursor: pointer;
}

.form-group label:has(input[type="checkbox"]) {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.hint {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
}

/* 操作按钮 */
.actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 32px;
}

.btn-primary {
  background: #6366f1;
  color: #fff;
  border: none;
  padding: 14px 32px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-primary:hover {
  background: #4f46e5;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.btn-secondary {
  background: #374151;
  color: #fff;
  border: none;
  padding: 14px 32px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #4b5563;
}
</style>
