<template>
  <div class="page">
    <div class="header">
      <div>
        <h1>上传模板</h1>
        <p class="subtitle">管理视频上传配置模板</p>
      </div>
      <button class="btn-primary" @click="showAdd = true">+ 添加模板</button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p>加载中...</p>
    </div>

    <div v-else-if="!templates.length" class="empty-state">
      <div class="empty-icon">📤</div>
      <h3>还没有上传模板</h3>
      <p>创建模板来配置视频上传参数</p>
      <button class="btn-primary" @click="showAdd = true">创建第一个模板</button>
    </div>

    <div v-else class="templates-grid">
      <div v-for="t in templates" :key="t.id" class="template-card">
        <div class="template-header">
          <h3>{{ t.name || '未命名模板' }}</h3>
          <button class="btn-danger-sm" @click="del(t.id)">删除</button>
        </div>
        <div class="template-body">
          <div class="info-item">
            <span class="label">标题:</span>
            <span>{{ t.title || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">分区:</span>
            <span>{{ t.tid || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">标签:</span>
            <span>{{ t.tags || '-' }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showAdd" class="modal" @click.self="showAdd = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>添加上传模板</h2>
          <button class="close-btn" @click="showAdd = false">✕</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>模板名称</label>
            <input v-model="form.name" placeholder="例如：游戏区模板">
          </div>
          <div class="form-group">
            <label>视频标题</label>
            <input v-model="form.title" placeholder="视频标题">
          </div>
          <div class="form-group">
            <label>分区ID</label>
            <input v-model="form.tid" type="number" placeholder="171">
          </div>
          <div class="form-group">
            <label>标签</label>
            <input v-model="form.tags" placeholder="标签1,标签2">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn-secondary" @click="showAdd = false">取消</button>
          <button class="btn-primary" @click="add">确定</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const templates = ref([])
const loading = ref(true)
const showAdd = ref(false)
const form = ref({ name: '', title: '', tid: '', tags: '' })

const load = async () => {
  try {
    loading.value = true
    const { data } = await api.getTemplates()
    templates.value = data || []
  } catch (error) {
    console.error('加载失败:', error)
    templates.value = []
  } finally {
    loading.value = false
  }
}

const add = async () => {
  try {
    await api.addTemplate(form.value)
    showAdd.value = false
    form.value = { name: '', title: '', tid: '', tags: '' }
    await load()
  } catch (error) {
    alert('添加失败: ' + (error.response?.data?.message || error.message))
  }
}

const del = async (id) => {
  if (!confirm('确定删除此模板？')) return
  try {
    await api.deleteTemplate(id)
    await load()
  } catch (error) {
    alert('删除失败: ' + (error.response?.data?.message || error.message))
  }
}

onMounted(load)
</script>

<style scoped>
.page { padding: 40px; max-width: 1400px; margin: 0 auto; }
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

.templates-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 20px; }
.template-card { background: #1a1a1a; border: 1px solid #2a2a2a; border-radius: 12px; padding: 20px; }
.template-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; padding-bottom: 16px; border-bottom: 1px solid #2a2a2a; }
.template-header h3 { font-size: 16px; color: #fff; margin: 0; }
.template-body { display: flex; flex-direction: column; gap: 12px; }
.info-item { display: flex; gap: 8px; font-size: 14px; }
.info-item .label { color: #9ca3af; min-width: 60px; }
.info-item span:last-child { color: #e0e0e0; }

.btn-primary { background: #6366f1; color: #fff; border: none; padding: 12px 24px; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 14px; transition: all 0.2s; }
.btn-primary:hover { background: #4f46e5; transform: translateY(-1px); }
.btn-secondary { background: #374151; color: #fff; border: none; padding: 12px 24px; border-radius: 10px; cursor: pointer; font-weight: 600; font-size: 14px; }
.btn-danger-sm { background: #dc2626; color: #fff; border: none; padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 12px; }
.btn-danger-sm:hover { background: #b91c1c; }

.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.8); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: #1a1a1a; border-radius: 20px; width: 520px; max-width: 90vw; border: 1px solid #2a2a2a; }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 24px 32px; border-bottom: 1px solid #2a2a2a; }
.modal-header h2 { font-size: 20px; font-weight: 700; color: #fff; margin: 0; }
.close-btn { background: none; border: none; color: #9ca3af; font-size: 24px; cursor: pointer; width: 32px; height: 32px; border-radius: 8px; }
.close-btn:hover { background: #2a2a2a; color: #fff; }
.modal-body { padding: 32px; }
.form-group { margin-bottom: 24px; }
.form-group:last-child { margin-bottom: 0; }
.form-group label { display: block; margin-bottom: 8px; color: #e0e0e0; font-size: 14px; font-weight: 600; }
.form-group input { width: 100%; padding: 14px 16px; background: #0f0f0f; border: 1px solid #2a2a2a; border-radius: 10px; color: #e0e0e0; font-size: 14px; }
.form-group input:focus { outline: none; border-color: #6366f1; box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1); }
.modal-footer { display: flex; gap: 12px; padding: 24px 32px; border-top: 1px solid #2a2a2a; }
.modal-footer button { flex: 1; }
</style>
