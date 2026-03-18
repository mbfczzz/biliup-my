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
.page { padding: 48px; max-width: 1400px; margin: 0 auto; min-height: 100vh; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 40px; }
h1 { font-size: 36px; font-weight: 700; color: #fff; margin-bottom: 8px; letter-spacing: -0.5px; }
.subtitle { font-size: 15px; color: #9ca3af; font-weight: 400; }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 120px 20px; color: #6b7280; }
.spinner { width: 48px; height: 48px; border: 3px solid rgba(99, 102, 241, 0.1); border-top-color: #6366f1; border-radius: 50%; animation: spin 0.8s linear infinite; margin-bottom: 20px; }
@keyframes spin { to { transform: rotate(360deg); } }

.empty-state { text-align: center; padding: 100px 40px; background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 2px dashed rgba(99, 102, 241, 0.2); border-radius: 20px; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); }
.empty-icon { font-size: 72px; margin-bottom: 20px; opacity: 0.6; filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.3)); }
.empty-state h3 { font-size: 22px; color: #fff; margin-bottom: 12px; font-weight: 600; }
.empty-state p { color: #9ca3af; margin-bottom: 32px; font-size: 15px; }

.templates-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 24px; }
.template-card { background: rgba(26, 26, 26, 0.6); backdrop-filter: blur(20px); border: 1px solid rgba(255, 255, 255, 0.06); border-radius: 20px; padding: 28px; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4); }
.template-card:hover { transform: translateY(-2px); box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5); border-color: rgba(99, 102, 241, 0.2); }
.template-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; padding-bottom: 20px; border-bottom: 1px solid rgba(255, 255, 255, 0.06); }
.template-header h3 { font-size: 18px; color: #fff; margin: 0; font-weight: 600; letter-spacing: -0.3px; }
.template-body { display: flex; flex-direction: column; gap: 14px; }
.info-item { display: flex; gap: 10px; font-size: 14px; }
.info-item .label { color: #9ca3af; min-width: 60px; font-weight: 500; }
.info-item span:last-child { color: #e5e7eb; }

.btn-primary { background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%); color: #fff; border: none; padding: 14px 28px; border-radius: 12px; cursor: pointer; font-weight: 600; font-size: 14px; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); box-shadow: 0 4px 16px rgba(99, 102, 241, 0.3); position: relative; overflow: hidden; }
.btn-primary::before { content: ''; position: absolute; top: 0; left: -100%; width: 100%; height: 100%; background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent); transition: left 0.5s; }
.btn-primary:hover::before { left: 100%; }
.btn-primary:hover { transform: translateY(-2px); box-shadow: 0 8px 24px rgba(99, 102, 241, 0.4); }
.btn-primary:active { transform: translateY(0); transition: transform 0.1s; }
.btn-primary:focus-visible { outline: 2px solid #6366f1; outline-offset: 2px; }
.btn-secondary { background: rgba(55, 65, 81, 0.6); color: #d1d5db; border: 1px solid rgba(255, 255, 255, 0.08); padding: 14px 28px; border-radius: 12px; cursor: pointer; font-weight: 600; font-size: 14px; transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); }
.btn-secondary:hover { background: rgba(75, 85, 99, 0.8); border-color: rgba(255, 255, 255, 0.12); transform: translateY(-1px); }
.btn-secondary:active { transform: translateY(0); transition: transform 0.1s; }
.btn-danger-sm { background: rgba(220, 38, 38, 0.9); color: #fff; border: 1px solid rgba(220, 38, 38, 0.3); padding: 8px 14px; border-radius: 8px; cursor: pointer; font-size: 12px; font-weight: 600; transition: all 0.2s; }
.btn-danger-sm:hover { background: #dc2626; box-shadow: 0 4px 12px rgba(220, 38, 38, 0.3); }

.modal { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.85); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(8px); }
.modal-content { background: rgba(26, 26, 26, 0.95); backdrop-filter: blur(20px); border-radius: 20px; width: 520px; max-width: 90vw; border: 1px solid rgba(255, 255, 255, 0.08); box-shadow: 0 24px 64px rgba(0, 0, 0, 0.6); }
.modal-header { display: flex; justify-content: space-between; align-items: center; padding: 28px 32px; border-bottom: 1px solid rgba(255, 255, 255, 0.06); }
.modal-header h2 { font-size: 20px; font-weight: 700; color: #fff; margin: 0; letter-spacing: -0.3px; }
.close-btn { background: none; border: none; color: #9ca3af; font-size: 24px; cursor: pointer; width: 36px; height: 36px; border-radius: 10px; transition: all 0.2s; display: flex; align-items: center; justify-content: center; }
.close-btn:hover { background: rgba(255, 255, 255, 0.08); color: #fff; }
.modal-body { padding: 32px; }
.form-group { margin-bottom: 24px; }
.form-group:last-child { margin-bottom: 0; }
.form-group label { display: block; margin-bottom: 10px; color: #d1d5db; font-size: 13px; font-weight: 600; letter-spacing: 0.3px; }
.form-group input { width: 100%; padding: 14px 16px; background: rgba(15, 15, 15, 0.8); border: 1.5px solid rgba(255, 255, 255, 0.08); border-radius: 12px; color: #e5e7eb; font-size: 14px; transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1); font-family: inherit; }
.form-group input::placeholder { color: #6b7280; opacity: 1; }
.form-group input:hover { border-color: rgba(255, 255, 255, 0.12); }
.form-group input:focus { outline: none; border-color: #6366f1; background: rgba(15, 15, 15, 0.95); box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.1); }
.form-group input:-webkit-autofill { -webkit-box-shadow: 0 0 0 1000px rgba(15, 15, 15, 0.95) inset; -webkit-text-fill-color: #e5e7eb; }
.modal-footer { display: flex; gap: 12px; padding: 24px 32px; border-top: 1px solid rgba(255, 255, 255, 0.06); }
.modal-footer button { flex: 1; }
</style>
