<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>账号管理</h1>
        <p class="subtitle">管理B站登录账号</p>
      </div>
      <button class="btn-primary" @click="showLogin = true">+ 添加账号</button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="!users.length" class="empty-state">
      <div class="empty-icon">👥</div>
      <h3>还没有添加账号</h3>
      <p>扫码登录B站账号以上传视频</p>
      <button class="btn-primary" @click="showLogin = true">添加第一个账号</button>
    </div>

    <div v-else class="users-grid">
      <div v-for="u in users" :key="u.id" class="user-card">
        <div class="user-info">
          <div class="user-name">{{ u.name || '未命名账号' }}</div>
          <div class="user-id">UID: {{ u.uid || '-' }}</div>
        </div>
        <button class="btn-danger-sm" @click="del(u.id)">删除</button>
      </div>
    </div>

    <div v-if="showLogin" class="modal" @click.self="showLogin = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>扫码登录</h2>
          <button class="close-btn" @click="showLogin = false">✕</button>
        </div>
        <div class="modal-body">
          <div v-if="qrLoading" class="qr-loading">
            <div class="spinner"></div>
            <p>生成二维码中...</p>
          </div>
          <div v-else-if="qrcode" class="qr-container">
            <img :src="qrcode" alt="二维码" class="qr-image">
            <p class="qr-hint">请使用B站APP扫码登录</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import api from '../api'

const users = ref([])
const loading = ref(true)
const showLogin = ref(false)
const qrcode = ref('')
const qrLoading = ref(false)

const load = async () => {
  try {
    loading.value = true
    const { data } = await api.getUsers()
    users.value = data || []
  } catch (error) {
    console.error('加载失败:', error)
    users.value = []
  } finally {
    loading.value = false
  }
}

const loadQrcode = async () => {
  try {
    qrLoading.value = true
    const { data } = await api.getQrcode()
    qrcode.value = data.url || ''
  } catch (error) {
    alert('获取二维码失败: ' + (error.response?.data?.message || error.message))
  } finally {
    qrLoading.value = false
  }
}

const del = async (id) => {
  if (!confirm('确定删除此账号？')) return
  try {
    await api.deleteUser(id)
    await load()
  } catch (error) {
    alert('删除失败: ' + (error.response?.data?.message || error.message))
  }
}

watch(showLogin, (val) => {
  if (val) {
    loadQrcode()
  } else {
    qrcode.value = ''
  }
})

onMounted(load)
</script>

<style scoped>
.page { padding: 48px; max-width: 1200px; margin: 0 auto; min-height: 100vh; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
h1 { font-size: 36px; font-weight: 700; color: #fff; margin-bottom: 8px; letter-spacing: -0.5px; }
.subtitle { font-size: 15px; color: #9ca3af; font-weight: 400; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 120px 20px; color: #6b7280; }
.spinner { width: 48px; height: 48px; border: 3px solid rgba(99, 102, 241, 0.1); border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 20px; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-state { text-align: center; padding: 100px 40px; background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 2px dashed rgba(99, 102, 241, 0.2); border-radius: 20px; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); }
.empty-icon { font-size: 72px; margin-bottom: 20px; opacity: 0.6; filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.3)); }
.empty-state h3 { font-size: 22px; color: #fff; margin-bottom: 12px; font-weight: 600; }
.empty-state p { color: #9ca3af; margin-bottom: 32px; font-size: 15px; }

.users-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 24px; }
.user-card { background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.06); border-radius: 20px; padding: 24px; display: flex; justify-content: space-between; align-items: center; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); }
.user-card:hover { transform: translateY(-2px); box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5); border-color: rgba(99, 102, 241, 0.2); }
.user-info { flex: 1; }
.user-name { font-size: 16px; font-weight: 600; color: #fff; margin-bottom: 6px; letter-spacing: -0.2px; }
.user-id { font-size: 13px; color: #9ca3af; font-family: 'SF Mono', 'Monaco', 'Courier New', monospace; }

.btn-primary { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); color: #fff; border: none; padding: 14px 28px; border-radius: 12px; cursor: pointer; font-weight: 600; font-size: 14px; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3); position: relative; overflow: hidden; }
.btn-primary::before { content: ''; position: absolute; top: 0; left: -100%; width: 100%; height: 100%; background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent); transition: left 0.5s; }
.btn-primary:hover::before { left: 100%; }
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4); }
.btn-primary:active { transform: translateY(0); transition: transform 0.1s; }
.btn-primary:focus-visible { outline: 2px solid #6366f1; outline-offset: 2px; }
.btn-danger-sm { background: rgba(220, 38, 38, 0.9); color: #fff; border: 1px solid rgba(220, 38, 38, 0.3); padding: 8px 14px; border-radius: 8px; cursor: pointer; font-size: 12px; font-weight: 600; transition: all 0.2s; }
.btn-danger-sm:hover { background: #dc2626; box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3); }

.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.85); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(8px); }
.modal-content { background: rgba(26, 26, 26, 0.95); backdrop-filter: blur(20px); border-radius: 20px; width: 420px; max-width: 90vw; border: 1px solid rgba(255, 255, 255, 0.08); box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6); }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 28px 32px; border-bottom: 1px solid rgba(255, 255, 255, 0.06); }
.modal-header h2 { font-size: 20px; font-weight: 700; color: #fff; margin: 0; letter-spacing: -0.3px; }
.close-btn { background: none; border: none; color: #9ca3af; font-size: 24px; cursor: pointer; width: 36px; height: 36px; border-radius: 10px; transition: all 0.2s; display: flex; align-items: center; justify-content: center; }
.close-btn:hover { background: rgba(255, 255, 255, 0.08); color: #fff; }
.modal-body { padding: 32px; }

.qr-loading { display: flex; flex-direction: column; align-items: center; padding: 40px 20px; }
.qr-container { text-align: center; }
.qr-image { width: 240px; height: 240px; border-radius: 16px; background: #fff; padding: 16px; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3); }
.qr-hint { margin-top: 20px; color: #9ca3af; font-size: 14px; }
</style>
