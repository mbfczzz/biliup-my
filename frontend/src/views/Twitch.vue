<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>录播管理</h1>
        <p class="subtitle">管理你的直播录制频道</p>
      </div>
      <button class="btn-primary" @click="showAdd = true">
        <span class="icon">+</span>
        添加频道
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 空状态 -->
    <div v-else-if="!channels.length" class="empty-state">
      <div class="empty-icon">📺</div>
      <h3>还没有添加频道</h3>
      <p>点击右上角"添加频道"按钮开始录制直播</p>
      <button class="btn-primary" @click="showAdd = true">添加第一个频道</button>
    </div>

    <!-- 频道列表 -->
    <div v-else class="card">
      <table>
        <thead>
          <tr>
            <th>频道名称</th>
            <th>直播地址</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in channels" :key="c.id">
            <td>
              <div class="channel-name">
                <strong>{{ c.name }}</strong>
              </div>
            </td>
            <td>
              <div class="url">{{ c.url }}</div>
            </td>
            <td>
              <span :class="'badge badge-' + (c.status || 'offline')">
                {{ c.status === 'live' ? '🔴 直播中' : '⚫ 离线' }}
              </span>
            </td>
            <td>
              <div class="actions">
                <button class="btn-sm" @click="pause(c.id)">
                  {{ c.paused ? '▶️ 恢复' : '⏸️ 暂停' }}
                </button>
                <button class="btn-sm btn-danger" @click="del(c.id)">🗑️ 删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 添加频道弹窗 -->
    <div v-if="showAdd" class="modal" @click.self="showAdd = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>添加录制频道</h2>
          <button class="close-btn" @click="showAdd = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>频道名称</label>
            <input
              v-model="form.name"
              placeholder="例如：某某主播"
              @keyup.enter="add"
            >
          </div>
          <div class="form-group">
            <label>直播地址</label>
            <input
              v-model="form.url"
              placeholder="https://www.twitch.tv/..."
              @keyup.enter="add"
            >
            <div class="hint">支持 Twitch、Bilibili、YouTube 等平台</div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="showAdd = false">取消</button>
          <button class="btn-primary" @click="add" :disabled="!form.name || !form.url">
            确定添加
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const channels = ref([])
const loading = ref(true)
const showAdd = ref(false)
const form = ref({ name: '', url: '' })

const load = async () => {
  try {
    loading.value = true
    const { data } = await api.getStreamers()
    channels.value = data || []
  } catch (error) {
    console.error('加载失败:', error)
    channels.value = []
  } finally {
    loading.value = false
  }
}

const add = async () => {
  if (!form.value.name || !form.value.url) return

  try {
    await api.addStreamer(form.value)
    showAdd.value = false
    form.value = { name: '', url: '' }
    await load()
  } catch (error) {
    alert('添加失败: ' + (error.response?.data?.message || error.message))
  }
}

const pause = async (id) => {
  try {
    await api.pauseStreamer(id)
    await load()
  } catch (error) {
    alert('操作失败: ' + (error.response?.data?.message || error.message))
  }
}

const del = async (id) => {
  if (!confirm('确定要删除这个频道吗？')) return

  try {
    await api.deleteStreamer(id)
    await load()
  } catch (error) {
    alert('删除失败: ' + (error.response?.data?.message || error.message))
  }
}

onMounted(load)
</script>

<style scoped>
.page {
  padding: 48px;
  max-width: 1400px;
  margin: 0 auto;
  min-height: 100vh;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.empty-state {
  text-align: center;
  padding: 100px 40px;
  background: rgba(26, 26, 26, 0.6);
  backdrop-filter: blur(20px);
  border: 2px dashed rgba(99, 102, 241, 0.2);
  border-radius: 20px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
}

.empty-icon {
  font-size: 72px;
  margin-bottom: 20px;
  opacity: 0.6;
  filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.3));
}

.empty-state h3 {
  font-size: 22px;
  color: #fff;
  margin-bottom: 12px;
  font-weight: 600;
}

.empty-state p {
  color: #9ca3af;
  margin-bottom: 32px;
  font-size: 15px;
}

.card {
  background: rgba(26, 26, 26, 0.6);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.card:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  border-color: rgba(99, 102, 241, 0.15);
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  background: rgba(15, 15, 15, 0.8);
  padding: 18px 24px;
  text-align: left;
  font-weight: 600;
  color: #9ca3af;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}

td {
  padding: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
  color: #e5e7eb;
}

tbody tr {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

tbody tr:hover {
  background: rgba(99, 102, 241, 0.05);
  transform: scale(1.002);
}

tbody tr:active {
  transform: scale(1);
}

.channel-name strong {
  font-size: 15px;
  color: #fff;
  font-weight: 600;
}

.url {
  color: #818cf8;
  font-size: 13px;
  font-family: 'SF Mono', 'Monaco', 'Courier New', monospace;
}

.badge {
  display: inline-block;
  padding: 6px 14px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
}

.badge-live {
  background: rgba(16, 185, 129, 0.15);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.badge-offline {
  background: rgba(107, 114, 128, 0.15);
  color: #9ca3af;
  border: 1px solid rgba(107, 114, 128, 0.2);
}

.actions {
  display: flex;
  gap: 10px;
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  border: none;
  padding: 14px 28px;
  border-radius: 12px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 8px;
  box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3);
  position: relative;
  overflow: hidden;
}

.btn-primary::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  transition: left 0.5s;
}

.btn-primary:hover::before {
  left: 100%;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4);
}

.btn-primary:active {
  transform: translateY(0);
  transition: transform 0.1s;
}

.btn-primary:focus-visible {
  outline: 2px solid #6366f1;
  outline-offset: 2px;
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background: rgba(55, 65, 81, 0.6);
  color: #d1d5db;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 14px 28px;
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

.btn-sm {
  background: rgba(42, 42, 42, 0.8);
  color: #e5e7eb;
  border: 1px solid rgba(255, 255, 255, 0.06);
  padding: 8px 16px;
  border-radius: 10px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-sm:hover {
  background: rgba(55, 65, 81, 0.9);
  border-color: rgba(255, 255, 255, 0.1);
  transform: translateY(-1px);
}

.btn-danger {
  background: rgba(220, 38, 38, 0.9);
  color: #fff;
  border: 1px solid rgba(220, 38, 38, 0.3);
}

.btn-danger:hover {
  background: #dc2626;
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3);
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(8px);
}

.modal-content {
  background: rgba(26, 26, 26, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  width: 520px;
  max-width: 90vw;
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28px 32px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin: 0;
  letter-spacing: -0.3px;
}

.close-btn {
  background: none;
  border: none;
  color: #9ca3af;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #fff;
}

.modal-body {
  padding: 32px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-group label {
  display: block;
  margin-bottom: 10px;
  color: #d1d5db;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: 0.3px;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  background: rgba(15, 15, 15, 0.8);
  border: 1.5px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  color: #e5e7eb;
  font-size: 14px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  font-family: inherit;
}

.form-group input::placeholder {
  color: #6b7280;
  opacity: 1;
}

.form-group input:hover {
  border-color: rgba(255, 255, 255, 0.12);
}

.form-group input:focus {
  outline: none;
  border-color: #6366f1;
  background: rgba(15, 15, 15, 0.95);
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1);
}

.form-group input:-webkit-autofill {
  -webkit-box-shadow: 0 0 0 1000px rgba(15, 15, 15, 0.95) inset;
  -webkit-text-fill-color: #e5e7eb;
}

.hint {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 24px 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-footer button {
  flex: 1;
}
</style>
