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
.page { padding: 40px; max-width: 1200px; margin: 0 auto; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 32px; }
h1 { font-size: 32px; font-weight: 700; color: #fff; margin-bottom: 8px; }
.subtitle { font-size: 14px; color: #9ca3af; }

.loading-state { display: flex; flex-direction: column; align-items: center; padding: 80px 20px; color: #9ca3af; }
.spinner { width: 40px; height: 40px; border: 3px solid #2a2a2a; border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 16px; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-state { text-align: center; padding: 80px 20px; background: #1a1a1a; border: 2px dashed #2a2a2a; border-radius: 16px; }
.empty-icon { font-size: 64px; margin-bottom: 16px; opacity: 0.5; }
.empty-state h3 { font-size: 20px; color: #fff; margin-bottom: 8px; }
.empty-state p { color: #9ca3af; margin-bottom: 24px; }

.users-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 20px; }
.user-card { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 12px; padding: 20px; display: flex; justify-content: space-between; align-items: center; }
.user-info { flex: 1; }
.user-name { font-size: 16px; font-weight: 600; color: #fff; margin-bottom: 4px; }
.user-id { font-size: 13px; color: #9ca3af; }

.btn-primary { background: #6366f1; color: #fff; border: none; padding: 12px 24px; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 14px; transition: all 0.2s; }
.btn-primary:hover { background: #4f46e5; transform: translateY(-1px); }
.btn-danger-sm { background: #dc2626; color: #fff; border: none; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.btn-danger-sm:hover { background: #b91c1c; }

.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.8); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: #1a1a1a; border-radius: 20px; width: 420px; max-width: 90vw; border: 1px solid #2a2a2a; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 24px 32px; border-bottom: 1px solid #2a2a2a; }
.modal-header h2 { font-size: 20px; font-weight: 700; color: #fff; margin: 0; }
.close-btn { background: none; border: none; color: #9ca3af; font-size: 24px; cursor: pointer; width: 32px; height: 32px; border-radius: 8px; }
.close-btn:hover { background: #2a2a2a; color: #fff; }
.modal-body { padding: 32px; }

.qr-loading { display: flex; flex-direction: column; align-items: center; padding: 40px 20px; }
.qr-container { text-align: center; }
.qr-image { width: 240px; height: 240px; border-radius: 12px; background: #fff; padding: 12px; }
.qr-hint { margin-top: 16px; color: #9ca3af; font-size: 14px; }
</style>
