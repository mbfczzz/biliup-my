<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">录播管理</h1>
        <p class="page-subtitle">管理 Twitch 直播录制频道</p>
      </div>
      <button class="btn btn-primary" @click="showAdd = true">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="7" y1="1" x2="7" y2="13"/><line x1="1" y1="7" x2="13" y2="7"/></svg>
        添加频道
      </button>
    </div>

    <div v-if="loading" class="state-box">
      <div class="spinner"></div>
      <p class="state-text">加载中...</p>
    </div>

    <div v-else-if="!channels.length" class="empty-box">
      <div class="empty-icon">📺</div>
      <h3 class="empty-title">还没有添加频道</h3>
      <p class="empty-desc">点击"添加频道"按钮开始录制 Twitch 直播</p>
      <button class="btn btn-primary" @click="showAdd = true">添加第一个频道</button>
    </div>

    <div v-else class="card">
      <table>
        <thead>
          <tr>
            <th>频道名称</th>
            <th>直播地址</th>
            <th>状态</th>
            <th style="width:180px">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in channels" :key="c.id">
            <td><span class="cell-name">{{ c.name }}</span></td>
            <td><span class="cell-url">{{ c.url }}</span></td>
            <td>
              <span :class="['badge', c.status === 'Working' ? 'badge-live' : 'badge-offline']">
                {{ c.status === 'Working' ? '直播中' : '离线' }}
              </span>
            </td>
            <td>
              <div class="cell-actions">
                <button class="btn btn-ghost btn-sm" @click="pause(c.id)">{{ c.paused ? '恢复' : '暂停' }}</button>
                <button class="btn btn-danger btn-sm" @click="del(c.id)">删除</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

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
            <p class="form-hint">输入 Twitch 频道地址</p>
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
.cell-name { font-weight:600; color:var(--c-text-1); }
.cell-url { color:var(--c-accent-hover); font-size:13px; font-family:var(--font-mono); }
.cell-actions { display:flex; gap:6px; }
</style>
