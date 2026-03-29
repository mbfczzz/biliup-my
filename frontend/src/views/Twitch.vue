<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">直播录制</h1>
        <p class="page-subtitle">添加和管理直播录制任务，监控录制状态</p>
      </div>
      <button class="btn btn-primary" @click="showAdd = true">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="7" y1="1" x2="7" y2="13"/><line x1="1" y1="7" x2="13" y2="7"/></svg>
        添加频道
      </button>
    </div>

    <!-- Stats bar -->
    <div v-if="channels.length" class="stats-bar">
      <div class="stat-item">
        <span class="stat-value">{{ channels.length }}</span>
        <span class="stat-label">频道总数</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-value stat-live">{{ channels.filter(c => c.status === 'Working').length }}</span>
        <span class="stat-label">录制中</span>
      </div>
      <div class="stat-divider"></div>
      <div class="stat-item">
        <span class="stat-value">{{ channels.filter(c => c.status !== 'Working').length }}</span>
        <span class="stat-label">待机中</span>
      </div>
    </div>

    <div v-if="loading" class="state-box">
      <div class="spinner"></div>
      <p class="state-text">加载中...</p>
    </div>

    <div v-else-if="!channels.length" class="empty-box">
      <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/><path d="M10 10l4-2v4l-4-2z" fill="currentColor" stroke="none"/></svg></div>
      <h3 class="empty-title">还没有添加频道</h3>
      <p class="empty-desc">添加直播频道后即可自动监控并录制</p>
      <button class="btn btn-primary" @click="showAdd = true">添加第一个频道</button>
    </div>

    <template v-else>
      <div class="channel-grid">
        <div v-for="c in channels" :key="c.id" :class="['channel-card card', { 'is-live': c.status === 'Working' }]">
          <div class="channel-header">
            <div class="channel-avatar">
              <span>{{ (c.name || '?')[0].toUpperCase() }}</span>
            </div>
            <div class="channel-info">
              <div class="channel-name">{{ c.name }}</div>
              <div class="channel-url">{{ c.url }}</div>
            </div>
            <span :class="['badge', c.status === 'Working' ? 'badge-live' : 'badge-offline']">
              {{ c.status === 'Working' ? '录制中' : '离线' }}
            </span>
          </div>
          <div v-if="c.status === 'Working'" class="live-indicator">
            <div class="live-bar"></div>
          </div>
          <div class="channel-actions">
            <button class="btn btn-ghost btn-sm" @click="pause(c.id)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path v-if="c.paused" d="M5 3l14 9-14 9V3z"/><template v-else><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></template></svg>
              {{ c.paused ? '恢复' : '暂停' }}
            </button>
            <button class="btn btn-danger btn-sm" @click="del(c.id)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/></svg>
              删除
            </button>
          </div>
        </div>
      </div>
    </template>

    <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
      <div class="modal-box">
        <div class="modal-head">
          <span class="modal-title">添加录制频道</span>
          <button class="modal-close" @click="showAdd = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">频道名称</label>
            <input class="form-input" v-model="form.name" placeholder="例如：某某主播" @keyup.enter="add">
          </div>
          <div class="form-group">
            <label class="form-label">直播地址</label>
            <input class="form-input" v-model="form.url" placeholder="https://www.twitch.tv/..." @keyup.enter="add">
            <p class="form-hint">输入直播频道地址</p>
          </div>
        </div>
        <div class="modal-foot">
          <button class="btn btn-ghost" @click="showAdd = false">取消</button>
          <button class="btn btn-primary" @click="add" :disabled="!form.name || !form.url">确定</button>
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
  try { await api.pauseStreamer(id); await load() }
  catch (error) { alert('操作失败: ' + (error.response?.data?.message || error.message)) }
}

const del = async (id) => {
  if (!confirm('确定要删除这个频道吗？')) return
  try { await api.deleteStreamer(id); await load() }
  catch (error) { alert('删除失败: ' + (error.response?.data?.message || error.message)) }
}

onMounted(load)
</script>

<style scoped>
.stats-bar {
  display:flex; align-items:center; gap:24px;
  padding:16px 24px; margin-bottom:24px;
  background:#fff; border:1px solid var(--c-border);
  border-radius:var(--radius); box-shadow:var(--shadow-sm);
}
.stat-item { display:flex; flex-direction:column; gap:2px; }
.stat-value { font-size:20px; font-weight:800; color:var(--c-text-1); letter-spacing:-.5px; }
.stat-value.stat-live { color:var(--c-success); }
.stat-label { font-size:11px; color:var(--c-text-4); font-weight:500; text-transform:uppercase; letter-spacing:.5px; }
.stat-divider { width:1px; height:32px; background:var(--c-border); }

.channel-grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(min(380px, 100%), 1fr)); gap:16px; }

.channel-card { padding:0; overflow:hidden; }
.channel-card.is-live { border-color:rgba(61,154,95,.2); }

.channel-header { display:flex; align-items:center; gap:14px; padding:20px 22px 16px; }

.channel-avatar {
  width:40px; height:40px; border-radius:10px;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
  font-size:16px; font-weight:700; flex-shrink:0;
}
.is-live .channel-avatar {
  background:linear-gradient(135deg, var(--c-success-soft), rgba(61,154,95,.15));
  color:var(--c-success);
}

.channel-info { flex:1; min-width:0; }
.channel-name { font-size:15px; font-weight:700; color:var(--c-text-1); margin-bottom:3px; }
.channel-url {
  font-size:12px; color:var(--c-text-4); font-family:var(--font-mono);
  overflow:hidden; text-overflow:ellipsis; white-space:nowrap;
}

.live-indicator { padding:0 22px 12px; }
.live-bar {
  height:3px; border-radius:2px;
  background:linear-gradient(90deg, var(--c-success), rgba(61,154,95,.3));
  animation:liveProgress 2s ease-in-out infinite;
}
@keyframes liveProgress {
  0%,100% { opacity:.6; } 50% { opacity:1; }
}

.channel-actions {
  display:flex; gap:8px; padding:12px 22px 18px;
  border-top:1px solid var(--c-border);
}
.channel-actions .btn { flex:1; justify-content:center; }

@media (max-width:768px) {
  .stats-bar { gap:16px; padding:14px 18px; }
  .stat-value { font-size:17px; }
  .channel-grid { grid-template-columns:1fr; }
}
</style>
