<template>
  <div class="page">
    <div class="header">
      <h1>系统配置</h1>
    </div>

    <div class="card" v-if="config">
      <div class="form-group">
        <label>下载路径</label>
        <input v-model="config.download_path" placeholder="输入下载路径">
      </div>
      <div class="form-group">
        <label>上传路径</label>
        <input v-model="config.upload_path" placeholder="输入上传路径">
      </div>
      <button class="btn-primary" @click="save">保存配置</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const config = ref(null)

const load = async () => {
  const { data } = await api.getConfig()
  config.value = data
}

const save = async () => {
  await api.updateConfig(config.value)
  alert('保存成功')
}

onMounted(load)
</script>

<style scoped>
.page { padding: 32px; }
.header { margin-bottom: 24px; }
h1 { font-size: 28px; font-weight: 600; color: #fff; }
.card {
  background: #1a1a1a;
  border: 1px solid #2a2a2a;
  border-radius: 12px;
  padding: 32px;
  max-width: 600px;
}
.form-group {
  margin-bottom: 24px;
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
.btn-primary {
  background: #6366f1;
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: background 0.2s;
}
.btn-primary:hover { background: #4f46e5; }
</style>
