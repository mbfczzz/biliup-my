<template>
  <div class="page">
    <div class="page-header">
      <h1>系统配置</h1>
      <p class="subtitle">配置录播和上传参数</p>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载配置中...</p>
    </div>

    <div v-else class="config-grid">
      <div class="card">
        <div class="card-icon">📁</div>
        <h3 class="card-title">路径配置</h3>
        <div class="card-content">
          <div class="form-item">
            <label>下载保存路径</label>
            <input v-model="config.download_path" placeholder="/path/to/downloads">
            <span class="hint">录制的视频文件保存位置</span>
          </div>
          <div class="form-item">
            <label>上传临时路径</label>
            <input v-model="config.upload_path" placeholder="/path/to/uploads">
            <span class="hint">上传前的临时文件存储位置</span>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-icon">⚙️</div>
        <h3 class="card-title">录制配置</h3>
        <div class="card-content">
          <div class="form-item">
            <label>视频质量</label>
            <select v-model="config.quality">
              <option value="best">最佳质量</option>
              <option value="high">高质量</option>
              <option value="medium">中等质量</option>
              <option value="low">低质量</option>
            </select>
          </div>
          <div class="form-item">
            <label>文件格式</label>
            <select v-model="config.format">
              <option value="mp4">MP4</option>
              <option value="flv">FLV</option>
              <option value="mkv">MKV</option>
            </select>
          </div>
        </div>
      </div>

      <div class="card">
        <div class="card-icon">🚀</div>
        <h3 class="card-title">上传配置</h3>
        <div class="card-content">
          <div class="form-item checkbox-item">
            <label class="checkbox-label">
              <input type="checkbox" v-model="config.auto_upload">
              <span>自动上传</span>
            </label>
            <span class="hint">录制完成后自动上传到平台</span>
          </div>
          <div class="form-item">
            <label>并发上传数</label>
            <input type="number" v-model.number="config.concurrent_uploads" min="1" max="5">
            <span class="hint">同时上传的文件数量（1-5）</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="!loading" class="actions">
      <button class="btn-secondary" @click="load">重置</button>
      <button class="btn-primary" @click="save">保存配置</button>
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
  padding: 48px;
  max-width: 1200px;
  margin: 0 auto;
  min-height: 100vh;
}

.page-header {
  margin-bottom: 40px;
}

h1 {
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
  letter-spacing: -0.5px;
}

.subtitle {
  font-size: 15px;
  color: #9ca3af;
  font-weight: 400;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120px 20px;
  color: #6b7280;
}

.spinner {
  width: 48px;
  height: 48px;
  border: 3px solid rgba(99, 102, 241, 0.1);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(360px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.card {
  background: rgba(26, 26, 26, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  padding: 28px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
}

.card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  border-color: rgba(99, 102, 241, 0.2);
}

.card-icon {
  font-size: 32px;
  margin-bottom: 16px;
  filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.3));
}

.card-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 24px;
  letter-spacing: -0.3px;
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-item label {
  font-size: 13px;
  font-weight: 600;
  color: #d1d5db;
  letter-spacing: 0.3px;
}

.form-item input[type="text"],
.form-item input[type="number"],
.form-item select {
  width: 100%;
  padding: 12px 16px;
  background: rgba(15, 15, 15, 0.8);
  border: 1.5px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  color: #e5e7eb;
  font-size: 14px;
  transition: all 0.2s;
  font-family: inherit;
}

.form-item input:focus,
.form-item select:focus {
  outline: none;
  border-color: #6366f1;
  background: rgba(15, 15, 15, 0.95);
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
}

.checkbox-item {
  gap: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-size: 14px;
  color: #e5e7eb;
  font-weight: 500;
}

.checkbox-label input[type="checkbox"] {
  width: 20px;
  height: 20px;
  cursor: pointer;
  accent-color: #6366f1;
}

.hint {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

.actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding-top: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  border: none;
  padding: 14px 32px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
}

.btn-primary:active {
  transform: translateY(0);
}

.btn-secondary {
  background: rgba(55, 65, 81, 0.6);
  color: #d1d5db;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 14px 32px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: rgba(75, 85, 99, 0.8);
  border-color: rgba(255, 255, 255, 0.12);
}
</style>
