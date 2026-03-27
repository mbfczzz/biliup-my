<template>
  <div class="page">
    <div class="page-header-row">
      <div>
        <h1 class="page-title">上传模板</h1>
        <p class="page-subtitle">管理 B 站视频上传配置模板</p>
      </div>
      <button class="btn btn-primary" @click="showAdd = true">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="7" y1="1" x2="7" y2="13"/><line x1="1" y1="7" x2="13" y2="7"/></svg>
        添加模板
      </button>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载中...</p></div>

    <div v-else-if="!templates.length" class="empty-box">
      <div class="empty-icon">📤</div>
      <h3 class="empty-title">还没有上传模板</h3>
      <p class="empty-desc">创建模板来配置视频上传参数</p>
      <button class="btn btn-primary" @click="showAdd = true">创建第一个模板</button>
    </div>

    <div v-else class="grid">
      <div v-for="t in templates" :key="t.id" class="tpl-card card">
        <div class="tpl-top">
          <h3 class="tpl-name">{{ t.template_name || t.name || '未命名模板' }}</h3>
          <button class="btn btn-danger btn-sm" @click="del(t.id)">删除</button>
        </div>
        <div class="tpl-fields">
          <div class="tpl-field"><span class="tpl-label">标题</span><span class="tpl-value">{{ t.title || '-' }}</span></div>
          <div class="tpl-field"><span class="tpl-label">分区</span><span class="tpl-value">{{ t.tid || '-' }}</span></div>
          <div class="tpl-field"><span class="tpl-label">标签</span><span class="tpl-value">{{ Array.isArray(t.tags) ? t.tags.join(', ') : (t.tags || '-') }}</span></div>
        </div>
      </div>
    </div>

    <div v-if="showAdd" class="modal-overlay" @click.self="showAdd = false">
      <div class="modal-box">
        <div class="modal-head">
          <span class="modal-title">添加上传模板</span>
          <button class="modal-close" @click="showAdd = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group"><label class="form-label">模板名称</label><input class="form-input" v-model="form.template_name" placeholder="例如：游戏区模板"></div>
          <div class="form-group"><label class="form-label">视频标题</label><input class="form-input" v-model="form.title" placeholder="视频标题"></div>
          <div class="form-group"><label class="form-label">分区 ID</label><input class="form-input" v-model="form.tid" type="number" placeholder="171"></div>
          <div class="form-group"><label class="form-label">标签</label><input class="form-input" v-model="form.tags" placeholder="标签1,标签2"></div>
        </div>
        <div class="modal-foot">
          <button class="btn btn-ghost" @click="showAdd = false">取消</button>
          <button class="btn btn-primary" @click="add">确定</button>
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
const form = ref({ template_name: '', title: '', tid: '', tags: '' })

const load = async () => {
  try { loading.value = true; const { data } = await api.getTemplates(); templates.value = data || [] }
  catch { templates.value = [] }
  finally { loading.value = false }
}

const add = async () => {
  try {
    const payload = { ...form.value, tags: form.value.tags ? form.value.tags.split(',').map(s => s.trim()) : [] }
    await api.addTemplate(payload)
    showAdd.value = false; form.value = { template_name: '', title: '', tid: '', tags: '' }; await load()
  } catch (e) { alert('添加失败: ' + (e.response?.data?.message || e.message)) }
}

const del = async (id) => {
  if (!confirm('确定删除此模板？')) return
  try { await api.deleteTemplate(id); await load() }
  catch (e) { alert('删除失败: ' + (e.response?.data?.message || e.message)) }
}

onMounted(load)
</script>

<style scoped>
.grid { display:grid; grid-template-columns:repeat(auto-fill, minmax(min(320px, 100%), 1fr)); gap:16px; }
.tpl-card { padding:24px; }
.tpl-top { display:flex; justify-content:space-between; align-items:center; margin-bottom:18px; padding-bottom:14px; border-bottom:1px solid var(--c-border); }
.tpl-name { font-size:16px; font-weight:700; color:var(--c-text-1); margin:0; }
.tpl-fields { display:flex; flex-direction:column; gap:10px; }
.tpl-field { display:flex; gap:8px; font-size:13px; }
.tpl-label { color:var(--c-text-3); min-width:48px; font-weight:500; }
.tpl-value { color:var(--c-text-2); }
</style>
