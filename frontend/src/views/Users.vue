<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">账号管理</h1>
        <p class="page-subtitle">管理 B 站登录账号</p>
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

    <div v-else class="grid">
      <div v-for="u in users" :key="u.id" class="user-card card">
        <div class="user-info">
          <div class="user-name">{{ u.name || '未命名账号' }}</div>
          <div class="user-id">{{ u.platform || 'bilibili' }}</div>
        </div>
        <button class="btn btn-danger btn-sm" @click="del(u.id)">删除</button>
      </div>
    </div>

    <div v-if="showLogin" class="modal-overlay" @click.self="showLogin = false">
      <div class="modal-box">
        <div class="modal-head">
          <span class="modal-title">扫码登录</span>
          <button class="modal-close" @click="showLogin = false">&times;</button>
        </div>
        <div class="modal-body qr-body">
          <div v-if="qrLoading" class="state-box" style="padding:40px 0">
            <div class="spinner"></div>
            <p class="state-text">生成二维码中...</p>
          </div>
          <div v-else-if="qrcode" class="qr-wrap">
            <img :src="qrcode" alt="二维码" class="qr-img">
            <p class="qr-hint">使用 B 站 APP 扫码登录</p>
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
  try { loading.value = true; const { data } = await api.getUsers(); users.value = data || [] }
  catch { users.value = [] }
  finally { loading.value = false }
}

const loadQrcode = async () => {
  try { qrLoading.value = true; const { data } = await api.getQrcode(); qrcode.value = data.url || '' }
  catch (e) { alert('获取二维码失败: ' + (e.response?.data?.message || e.message)) }
  finally { qrLoading.value = false }
}

const del = async (id) => {
  if (!confirm('确定删除此账号？')) return
  try { await api.deleteUser(id); await load() }
  catch (e) { alert('删除失败: ' + (e.response?.data?.message || e.message)) }
}

watch(showLogin, v => { if (v) loadQrcode(); else qrcode.value = '' })
onMounted(load)
</script>

<style scoped>
.grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(min(300px, 100%), 1fr)); gap:14px; }
.user-card { padding:18px 22px; display:flex; justify-content:space-between; align-items:center; }
.user-name { font-size:14px; font-weight:600; color:var(--c-text-1); margin-bottom:3px; }
.user-id { font-size:12px; color:var(--c-text-3); font-family:var(--font-mono); }
.qr-body { display:flex; justify-content:center; }
.qr-wrap { text-align:center; }
.qr-img { width:200px; height:200px; border-radius:12px; background:#fff; padding:12px; }
.qr-hint { margin-top:16px; color:var(--c-text-3); font-size:13px; }
</style>
