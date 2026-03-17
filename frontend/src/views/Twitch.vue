<template>
  <div class="page">
    <div class="header">
      <h1>录播管理</h1>
      <button class="btn-primary" @click="showAdd = true">+ 添加频道</button>
    </div>

    <div class="card" v-if="channels.length">
      <table>
        <thead>
          <tr>
            <th>频道名</th>
            <th>URL</th>
            <th>状态</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="c in channels" :key="c.id">
            <td><strong>{{ c.name }}</strong></td>
            <td class="url">{{ c.url }}</td>
            <td><span :class="'status ' + c.status">{{ c.status }}</span></td>
            <td>
              <button class="btn-sm" @click="pause(c.id)">{{ c.paused ? '恢复' : '暂停' }}</button>
              <button class="btn-sm btn-danger" @click="del(c.id)">删除</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="showAdd" class="modal" @click.self="showAdd = false">
      <div class="modal-content">
        <h2>添加频道</h2>
        <div class="form-group">
          <label>频道名</label>
          <input v-model="form.name" placeholder="输入频道名称">
        </div>
        <div class="form-group">
          <label>Twitch URL</label>
          <input v-model="form.url" placeholder="https://www.twitch.tv/...">
        </div>
        <div class="modal-actions">
          <button class="btn-primary" @click="add">确定</button>
          <button class="btn-secondary" @click="showAdd = false">取消</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const channels = ref([])
const showAdd = ref(false)
const form = ref({ name: '', url: '' })

const load = async () => {
  const { data } = await api.getStreamers()
  channels.value = data
}

const add = async () => {
  await api.addStreamer(form.value)
  showAdd.value = false
  form.value = { name: '', url: '' }
  load()
}

const pause = async (id) => {
  await api.pauseStreamer(id)
  load()
}

const del = async (id) => {
  if (confirm('确定删除?')) {
    await api.deleteStreamer(id)
    load()
  }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 32px; }
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}
h1 { font-size: 28px; font-weight: 600; color: #fff; }
.card {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  overflow: hidden;
}
table { width: 100%; border-collapse: collapse; }
th {
  background: #252525;
  padding: 16px;
  text-align: left;
  font-weight: 600;
  color: #a0a0a0;
  font-size: 13px;
  text-transform: uppercase;
}
td {
  padding: 16px;
  border-top: 1px solid #2a2a2a;
  color: #e0e0e0;
}
.url { color: #6366f1; font-size: 14px; }
.status {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}
.status.live { background: #10b981; color: #fff; }
.status.offline { background: #374151; color: #9ca3af; }
.btn-primary {
  background: #6366f1;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: background 0.2s;
}
.btn-primary:hover { background: #4f46e5; }
.btn-secondary {
  background: #374151;
  color: #fff;
  border: none;
  padding: 10px 20px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
}
.btn-sm {
  background: #2a2a2a;
  color: #e0e0e0;
  border: none;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  margin-right: 8px;
}
.btn-sm:hover { background: #374151; }
.btn-danger { background: #dc2626; color: #fff; }
.btn-danger:hover { background: #b91c1c; }
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal-content {
  background: #1a1a1a;
  padding: 32px;
  border-radius: 16px;
  width: 480px;
  border: 1px solid #2a2a2a;
}
.modal-content h2 {
  margin-bottom: 24px;
  color: #fff;
  font-size: 20px;
}
.form-group {
  margin-bottom: 20px;
}
.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #a0a0a0;
  font-size: 14px;
  font-weight: 500;
}
.form-group input {
  width: 100%;
  padding: 12px;
  background: #0f0f0f;
  border: 1px solid #2a2a2a;
  border-radius: 8px;
  color: #e0e0e0;
  font-size: 14px;
}
.form-group input:focus {
  outline: none;
  border-color: #6366f1;
}
.modal-actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}
.modal-actions button { flex: 1; }
</style>
