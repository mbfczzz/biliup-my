<template>
  <div class="config">
    <h1>配置</h1>
    <div v-if="config">
      <div class="field">
        <label>下载路径</label>
        <input v-model="config.download_path">
      </div>
      <div class="field">
        <label>上传路径</label>
        <input v-model="config.upload_path">
      </div>
      <button @click="save">保存</button>
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
.config { padding: 20px; }
.field { margin: 15px 0; }
.field label { display: block; margin-bottom: 5px; font-weight: bold; }
.field input { width: 400px; padding: 8px; }
button { margin-top: 20px; padding: 10px 20px; cursor: pointer; }
</style>
