<template>
  <div class="twitch">
    <h1>Twitch 录播</h1>
    <button @click="showAdd = true">添加频道</button>

    <table v-if="channels.length">
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
          <td>{{ c.name }}</td>
          <td>{{ c.url }}</td>
          <td>{{ c.status }}</td>
          <td>
            <button @click="pause(c.id)">{{ c.paused ? '恢复' : '暂停' }}</button>
            <button @click="del(c.id)">删除</button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="showAdd" class="modal">
      <div class="modal-content">
        <h2>添加 Twitch 频道</h2>
        <input v-model="form.name" placeholder="频道名">
        <input v-model="form.url" placeholder="Twitch URL">
        <button @click="add">确定</button>
        <button @click="showAdd = false">取消</button>
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
.twitch { padding: 20px; }
button { margin: 10px 5px; padding: 8px 16px; cursor: pointer; }
table { width: 100%; border-collapse: collapse; margin-top: 20px; }
th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
th { background: #9146ff; color: white; }
.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.5); display: flex; align-items: center; justify-content: center; }
.modal-content { background: white; padding: 30px; border-radius: 8px; }
.modal-content input { display: block; width: 300px; margin: 10px 0; padding: 8px; }
</style>
