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
  padding: 40px;
  max-width: 1400px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

/* 空状态 */
.empty-state {
  text-align: center;
  padding: 80px 20px;
  background: #1a1a1a;
  border: 2px dashed #2a2a2a;
  border-radius: 16px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  color: #fff;
  margin-bottom: 8px;
}

.empty-state p {
  color: #9ca3af;
  margin-bottom: 24px;
}

/* 卡片 */
.card {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
}

/* 表格 */
table {
  width: 100%;
  border-collapse: collapse;
}

th {
  background: #252525;
  padding: 16px 20px;
  text-align: left;
  font-weight: 600;
  color: #9ca3af;
  font-size: 12px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

td {
  padding: 20px;
  border-top: 1px solid #2a2a2a;
  color: #e0e0e0;
}

.channel-name strong {
  font-size: 15px;
  color: #fff;
}

.url {
  color: #6366f1;
  font-size: 13px;
  font-family: 'Courier New', monospace;
}

/* 状态徽章 */
.badge {
  display: inline-block;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.badge-live {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.badge-offline {
  background: rgba(107, 114, 128, 0.2);
  color: #9ca3af;
}

/* 操作按钮 */
.actions {
  display: flex;
  gap: 8px;
}

.btn-primary {
  background: #6366f1;
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary:hover {
  background: #4f46e5;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background: #374151;
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: 600;
  font-size: 14px;
  transition: all 0.2s;
}

.btn-secondary:hover {
  background: #4b5563;
}

.btn-sm {
  background: #2a2a2a;
  color: #e0e0e0;
  border: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-sm:hover {
  background: #374151;
  transform: translateY(-1px);
}

.btn-danger {
  background: #dc2626;
  color: #fff;
}

.btn-danger:hover {
  background: #b91c1c;
}

/* 弹窗 */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: #1a1a1a;
  border-radius: 20px;
  width: 520px;
  max-width: 90vw;
  border: 1px solid #2a2a2a;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px;
  border-bottom: 1px solid #2a2a2a;
}

.modal-header h2 {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: #9ca3af;
  font-size: 24px;
  cursor: pointer;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: #2a2a2a;
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
  margin-bottom: 8px;
  color: #e0e0e0;
  font-size: 14px;
  font-weight: 600;
}

.form-group input {
  width: 100%;
  padding: 14px 16px;
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  border-radius: 10px;
  color: #e0e0e0;
  font-size: 14px;
  transition: all 0.2s;
}

.form-group input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.hint {
  margin-top: 8px;
  font-size: 12px;
  color: #6b7280;
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 24px 32px;
  border-top: 1px solid #2a2a2a;
}

.modal-footer button {
  flex: 1;
}
</style>
