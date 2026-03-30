<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">B 站账号</h1>
        <p class="page-subtitle">管理登录凭证，支持多账号切换</p>
      </div>
      <button class="btn btn-primary" @click="showLogin = true">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="7" y1="1" x2="7" y2="13"/><line x1="1" y1="7" x2="13" y2="7"/></svg>
        添加账号
      </button>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载中...</p></div>

    <div v-else-if="!users.length" class="empty-box">
      <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14c-4.418 0-8 1.79-8 4v2h16v-2c0-2.21-3.582-4-8-4z"/></svg></div>
      <h3 class="empty-title">还没有添加账号</h3>
      <p class="empty-desc">扫码登录 B 站账号以上传视频</p>
      <button class="btn btn-primary" @click="showLogin = true">添加第一个账号</button>
    </div>

    <div v-else class="user-list">
      <div v-for="u in users" :key="u.id" class="user-card card">
        <div class="user-left">
          <div class="user-avatar">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14c-4.418 0-8 1.79-8 4v2h16v-2c0-2.21-3.582-4-8-4z"/></svg>
          </div>
          <div class="user-info">
            <div class="user-name">{{ u.name || '未命名账号' }}</div>
            <div class="user-platform">
              <span class="platform-badge">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z"/></svg>
                {{ u.platform || 'bilibili' }}
              </span>
            </div>
          </div>
        </div>
        <button class="btn btn-danger btn-sm" @click="del(u.id)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/></svg>
          删除
        </button>
      </div>
    </div>

    <div v-if="showLogin" class="modal-overlay" @click.self="showLogin = false">
      <div class="modal-box qr-modal">
        <div class="modal-head">
          <span class="modal-title">扫码登录</span>
          <button class="modal-close" @click="showLogin = false">&times;</button>
        </div>
        <div class="modal-body qr-body">
          <div v-if="qrLoading" class="qr-loading">
            <div class="spinner"></div>
            <p class="state-text">正在生成二维码...</p>
          </div>
          <div v-else-if="qrcode" class="qr-wrap">
            <div class="qr-frame">
              <img :src="qrcode" alt="二维码" class="qr-img">
            </div>
            <div class="qr-text">
              <p class="qr-hint">使用 <strong>B 站 APP</strong> 扫码登录</p>
              <p class="qr-sub-hint">扫码后请在手机上确认</p>
            </div>
          </div>
          <div v-else class="qr-error">
            <p class="state-text">二维码获取失败，请重试</p>
            <button class="btn btn-ghost btn-sm" @click="loadQrcode" style="margin-top:12px">重新获取</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import QRCode from 'qrcode'
import api from '../api'

const users = ref([])
const loading = ref(true)
const showLogin = ref(false)
const qrcode = ref('')
const qrLoading = ref(false)

const load = async () => {
  try { loading.value = true; const { data } = await api.getUsers(); users.value = data || [] }
  catch { users.value = [] }
  finally { loading.value = false }
}

const qrData = ref(null)

const loadQrcode = async () => {
  try {
    qrLoading.value = true
    const { data } = await api.getQrcode()
    qrData.value = data
    const url = data.data?.url || ''
    if (url) {
      qrcode.value = await QRCode.toDataURL(url, { width: 280, margin: 2 })
      // Start polling for scan result
      pollLogin(data)
    }
  }
  catch (e) {
    console.error('QR code error:', e, e.response?.status, e.response?.data)
    alert('获取二维码失败: ' + JSON.stringify(e.response?.data || e.message))
  }
  finally { qrLoading.value = false }
}

let pollController = null

const pollLogin = async (qrData) => {
  if (pollController) pollController.abort()
  pollController = new AbortController()
  try {
    const { data } = await api.loginByQrcode(qrData, pollController.signal)
    if (data.code === 0) {
      showLogin.value = false
      await load()
    }
  } catch (e) {
    if (e.name === 'CanceledError' || e.code === 'ERR_CANCELED') return
    console.error('Login error:', e)
    // 二维码过期或其他错误，自动重新获取二维码
    if (showLogin.value) {
      qrcode.value = ''
      await loadQrcode()
    }
  } finally {
    pollController = null
  }
}

const del = async (id) => {
  if (!confirm('确定删除此账号？')) return
  try { await api.deleteUser(id); await load() }
  catch (e) { alert('删除失败: ' + (e.response?.data?.message || e.message)) }
}

watch(showLogin, v => {
  if (v) { loadQrcode() }
  else { qrcode.value = ''; if (pollController) { pollController.abort(); pollController = null } }
})
onBeforeUnmount(() => {
  if (pollController) { pollController.abort(); pollController = null }
})
onMounted(load)
</script>

<style scoped>
.user-list { display:flex; flex-direction:column; gap:10px; }

.user-card {
  padding:18px 22px; display:flex; justify-content:space-between; align-items:center;
  transition:all .2s var(--ease);
}
.user-card:hover { transform:translateX(4px); }

.user-left { display:flex; align-items:center; gap:14px; }
.user-avatar {
  width:42px; height:42px; border-radius:12px;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
  flex-shrink:0;
}
.user-info { display:flex; flex-direction:column; gap:4px; }
.user-name { font-size:14px; font-weight:600; color:var(--c-text-1); }
.user-platform { display:flex; }
.platform-badge {
  display:inline-flex; align-items:center; gap:4px;
  font-size:11px; color:var(--c-text-4); font-weight:500;
  padding:2px 8px 2px 5px; border-radius:6px;
  background:rgba(0,0,0,.03);
}

.qr-modal { width:400px; }
.qr-body { display:flex; justify-content:center; padding-top:16px; padding-bottom:8px; }
.qr-loading { display:flex; flex-direction:column; align-items:center; padding:48px 0; }
.qr-error { display:flex; flex-direction:column; align-items:center; padding:48px 0; }
.qr-wrap { text-align:center; }
.qr-frame {
  display:inline-block; padding:20px;
  background:linear-gradient(145deg, #fafafa, #fff);
  border:1.5px solid rgba(0,0,0,.06);
  border-radius:20px;
  box-shadow:0 4px 20px rgba(0,0,0,.04), 0 0 0 1px rgba(255,255,255,.8) inset;
  transition:transform .2s var(--ease);
}
.qr-frame:hover { transform:scale(1.02); }
.qr-img { width:200px; height:200px; border-radius:12px; display:block; }
.qr-text { margin-top:24px; }
.qr-hint { color:var(--c-text-2); font-size:14px; font-weight:500; }
.qr-hint strong { font-weight:700; color:var(--c-text-1); }
.qr-sub-hint { margin-top:6px; color:var(--c-text-4); font-size:12px; }
</style>
