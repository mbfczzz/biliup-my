<template>
  <div class="page">
    <!-- ===== List View ===== -->
    <template v-if="!showForm">
      <div class="page-header-row">
        <div>
          <h1 class="page-title">投稿模板</h1>
          <p class="page-subtitle">预设视频标题、分区、标签等投稿参数</p>
        </div>
        <button class="btn btn-primary" @click="openAdd">
          <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5"><line x1="7" y1="1" x2="7" y2="13"/><line x1="1" y1="7" x2="13" y2="7"/></svg>
          添加模板
        </button>
      </div>

      <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载中...</p></div>

      <div v-else-if="!templates.length" class="empty-box">
        <div class="empty-icon"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M12 4v12M8 8l4-4 4 4"/></svg></div>
        <h3 class="empty-title">还没有投稿模板</h3>
        <p class="empty-desc">创建模板来预设投稿参数，一键应用</p>
        <button class="btn btn-primary" @click="openAdd">创建第一个模板</button>
      </div>

      <div v-else class="grid">
        <div v-for="t in templates" :key="t.id" class="tpl-card">
          <div class="tpl-accent"></div>
          <div class="tpl-header">
            <div class="tpl-icon">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6"><path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/><path d="M14 2v6h6M16 13H8M16 17H8M10 9H8"/></svg>
            </div>
            <div class="tpl-title-group">
              <h3 class="tpl-name">{{ t.template_name || '未命名模板' }}</h3>
              <span class="tpl-copyright-badge" :class="t.copyright === 1 ? 'is-original' : ''">
                {{ t.copyright === 1 ? '自制' : '转载' }}
              </span>
            </div>
            <div class="tpl-actions">
              <button class="tpl-action" @click="openEdit(t)" title="编辑">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
              </button>
              <button class="tpl-action tpl-action-del" @click="del(t.id)" title="删除">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M3 6h18M8 6V4a2 2 0 012-2h4a2 2 0 012 2v2M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6"/></svg>
              </button>
            </div>
          </div>
          <div class="tpl-title-preview" v-if="t.title">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg>
            <span>{{ t.title }}</span>
          </div>
          <div class="tpl-body">
            <div class="tpl-info-grid">
              <div class="tpl-info-item" v-if="t.tid"><span class="tpl-info-label">分区</span><span class="tpl-info-value">{{ t.tid }}</span></div>
              <div class="tpl-info-item" v-if="t.user_cookie"><span class="tpl-info-label">账号</span><span class="tpl-info-value">{{ getUserName(t.user_cookie) }}</span></div>
              <div class="tpl-info-item" v-if="t.cover_path"><span class="tpl-info-label">封面</span><span class="tpl-info-value tpl-mono">{{ t.cover_path }}</span></div>
              <div class="tpl-info-item" v-if="t.copyright === 2 && t.copyright_source"><span class="tpl-info-label">来源</span><span class="tpl-info-value tpl-mono">{{ t.copyright_source }}</span></div>
            </div>
            <div class="tpl-tags" v-if="Array.isArray(t.tags) ? t.tags.length : t.tags">
              <span class="tpl-tag" v-for="tag in (Array.isArray(t.tags) ? t.tags : (t.tags || '').split(','))" :key="tag">{{ tag.trim() }}</span>
            </div>
            <div class="tpl-desc-block" v-if="t.description"><p class="tpl-desc">{{ t.description }}</p></div>
            <div class="tpl-features" v-if="t.dolby || t.hires || t.up_close_danmu || t.up_close_reply || t.up_selection_reply || t.no_reprint || t.is_only_self || t.charging_pay">
              <span class="tpl-feat" v-if="t.dolby">杜比音效</span>
              <span class="tpl-feat" v-if="t.hires">Hi-Res</span>
              <span class="tpl-feat warn" v-if="t.up_close_danmu">关闭弹幕</span>
              <span class="tpl-feat warn" v-if="t.up_close_reply">关闭评论</span>
              <span class="tpl-feat" v-if="t.up_selection_reply">精选评论</span>
              <span class="tpl-feat" v-if="t.no_reprint">禁止转载</span>
              <span class="tpl-feat muted" v-if="t.is_only_self">仅自己可见</span>
              <span class="tpl-feat" v-if="t.charging_pay">充电面板</span>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ===== Full-page Form View ===== -->
    <template v-else>
      <div class="form-page-header">
        <button class="back-btn" @click="showForm = false">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
        </button>
        <h1 class="page-title">{{ editing ? '编辑投稿模板' : '新建投稿模板' }}</h1>
        <div class="form-page-actions">
          <button class="btn btn-ghost" @click="showForm = false">取消</button>
          <button class="btn btn-primary" @click="save">
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M20 6L9 17l-5-5"/></svg>
            {{ editing ? '保存' : '创建' }}
          </button>
        </div>
      </div>

      <div class="form-card">
        <!-- 基本信息 -->
        <h3 class="sec-title">基本信息</h3>
        <div class="sec-body">
          <div class="row">
            <label class="row-label">模板名称 <i>*</i></label>
            <div class="row-ctrl"><input class="form-input" v-model="form.template_name" placeholder="例如：jingggxd"></div>
          </div>
          <div class="row">
            <label class="row-label">投稿账号 <i>*</i></label>
            <div class="row-ctrl">
              <select class="form-input" v-model="form.user_cookie">
                <option value="">未选择</option>
                <option v-for="u in users" :key="u.id" :value="u.value">{{ u.name || u.value }}</option>
              </select>
            </div>
          </div>
        </div>

        <!-- 基本设置 -->
        <h3 class="sec-title">基本设置</h3>
        <div class="sec-body">
          <div class="row">
            <label class="row-label">视频标题</label>
            <div class="row-ctrl">
              <input class="form-input" v-model="form.title" placeholder="prx jinggg 直播日常！%Y-%m-%d 直播回放 | {title}">
              <div class="row-hint">
                <code>{streamer}</code>: 录播备注 &nbsp; <code>{title}</code>: 直播标题<br>
                <code>%Y-%m-%d %H_%M_%S</code>: 开始录制时的 年-月-日 时_分_秒
              </div>
            </div>
          </div>
          <div class="row">
            <label class="row-label">类型</label>
            <div class="row-ctrl">
              <div class="type-selector">
                <label class="type-opt" :class="{ active: form.copyright === 2 }"><input type="radio" :value="2" v-model="form.copyright"><span>转载</span></label>
                <label class="type-opt" :class="{ active: form.copyright === 1 }"><input type="radio" :value="1" v-model="form.copyright"><span>自制</span></label>
              </div>
              <input v-if="form.copyright === 2" class="form-input" style="margin-top:10px" v-model="form.copyright_source" placeholder="转载来源 URL">
              <p v-if="form.copyright === 2" class="row-sub-hint">如不填写转载来源默认为直播间地址</p>
            </div>
          </div>
          <div class="row">
            <label class="row-label">分区 <i>*</i></label>
            <div class="row-ctrl">
              <input class="form-input" v-model.number="form.tid" type="number" placeholder="171" style="max-width:200px">
              <p class="row-sub-hint">游戏 / 网络游戏 等分区对应的 ID</p>
            </div>
          </div>
          <div class="row">
            <label class="row-label">标签 <i>*</i></label>
            <div class="row-ctrl">
              <div class="tags-input-wrap">
                <span class="tag-chip" v-for="(tag, i) in tagsList" :key="i">{{ tag }}<button class="tag-x" @click="removeTag(i)">&times;</button></span>
                <input class="tags-input" v-model="tagInput" @keydown.enter.prevent="addTag" @keydown.,.prevent="addTag" placeholder="输入标签后按回车添加">
              </div>
            </div>
          </div>
          <div class="row">
            <label class="row-label">视频封面</label>
            <div class="row-ctrl"><input class="form-input" v-model="form.cover_path" placeholder="/opt/cover/jingggxd/jingggxd.jpg"></div>
          </div>
          <div class="row">
            <label class="row-label">简介</label>
            <div class="row-ctrl">
              <textarea class="form-input form-textarea" v-model="form.description" placeholder="视频简介" rows="5" maxlength="2000"></textarea>
              <span class="char-count">{{ (form.description || '').length }}/2000</span>
            </div>
          </div>
          <div class="row">
            <label class="row-label">额外字段</label>
            <div class="row-ctrl">
              <textarea class="form-input form-textarea form-textarea-sm" v-model="form.extra_fields" placeholder='Json格式，示例：{key: value}' rows="3" maxlength="2000"></textarea>
              <span class="char-count">{{ (form.extra_fields || '').length }}/2000</span>
            </div>
          </div>
          <div class="row">
            <label class="row-label">定时发布</label>
            <div class="row-ctrl">
              <div class="inline-row">
                <label class="toggle"><input type="checkbox" v-model="enableDtime"><span class="toggle-track"><span class="toggle-thumb"></span></span></label>
                <input v-if="enableDtime" class="form-input" type="datetime-local" v-model="dtimeStr" style="max-width:280px">
              </div>
              <p v-if="enableDtime" class="row-sub-hint">(当前+2小时 &lt;= 可选时间 &lt;= 当前+15天，转载稿件请年制发布时间为准，上传速度不佳的机器请调减开启，或设置更大的延迟时间。)</p>
            </div>
          </div>
        </div>

        <!-- 更多设置 -->
        <h3 class="sec-title sec-toggle" @click="showMore = !showMore">
          更多设置
          <svg :class="['chev', { open: showMore }]" width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M4 5.5l3 3 3-3"/></svg>
        </h3>
        <div v-if="showMore" class="sec-body">
          <div class="row">
            <label class="row-label">音效设置</label>
            <div class="row-ctrl">
              <div class="check-group">
                <label class="ck"><input type="checkbox" v-model="form.dolby" true-value="1" false-value="0"><span>杜比音效</span></label>
                <label class="ck"><input type="checkbox" v-model="form.hires" true-value="1" false-value="0"><span>Hi-Res无损音质</span></label>
              </div>
            </div>
          </div>
          <div class="row">
            <label class="row-label">互动设置</label>
            <div class="row-ctrl">
              <div class="check-group">
                <label class="ck"><input type="checkbox" v-model="form.up_close_danmu" :true-value="true" :false-value="false"><span>关闭弹幕</span></label>
                <label class="ck"><input type="checkbox" v-model="form.up_close_reply" :true-value="true" :false-value="false"><span>关闭评论</span></label>
                <label class="ck"><input type="checkbox" v-model="form.up_selection_reply" :true-value="true" :false-value="false"><span>开启精选评论</span></label>
              </div>
            </div>
          </div>
          <div class="row">
            <label class="row-label">粉丝动态</label>
            <div class="row-ctrl"><input class="form-input" v-model="form.dynamic" placeholder="发布时同步发送的动态内容"></div>
          </div>
          <div class="row">
            <label class="row-label">上传插件</label>
            <div class="row-ctrl">
              <select class="form-input" v-model="form.uploader" style="max-width:240px">
                <option value="">默认</option>
                <option value="biliup-rs">biliup-rs</option>
                <option value="biliup-python">biliup-python</option>
              </select>
            </div>
          </div>
          <div class="row">
            <label class="row-label">自制声明</label>
            <div class="row-ctrl"><label class="toggle"><input type="checkbox" v-model="form.no_reprint" true-value="1" false-value="0"><span class="toggle-track"><span class="toggle-thumb"></span></span></label></div>
          </div>
          <div class="row">
            <label class="row-label">仅自己可见</label>
            <div class="row-ctrl"><label class="toggle"><input type="checkbox" v-model="form.is_only_self" true-value="1" false-value="0"><span class="toggle-track"><span class="toggle-thumb"></span></span></label></div>
          </div>
          <div class="row">
            <label class="row-label">开启充电面板</label>
            <div class="row-ctrl"><label class="toggle"><input type="checkbox" v-model="form.charging_pay" true-value="1" false-value="0"><span class="toggle-track"><span class="toggle-thumb"></span></span></label></div>
          </div>
        </div>

        <!-- 简介@替换 -->
        <h3 class="sec-title">简介@替换</h3>
        <div class="sec-body">
          <p class="sec-desc">
            如需在简介中@别人，请使用此选项。示例：<br>
            简介：{streamer}主播直播间地址：{url}【@credit】<br>
            其中的"@credit"会依次替换为下面输入的@
          </p>
          <div v-for="(c, i) in creditsList" :key="i" class="credit-row">
            <input class="form-input" v-model="creditsList[i]" placeholder="@用户名">
            <button class="credit-del" @click="creditsList.splice(i, 1)">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
            </button>
          </div>
          <button class="add-row-btn" @click="creditsList.push('')">
            <svg width="14" height="14" viewBox="0 0 14 14" fill="none" stroke="currentColor" stroke-width="2"><line x1="7" y1="2" x2="7" y2="12"/><line x1="2" y1="7" x2="12" y2="7"/></svg>
            添加行
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const templates = ref([])
const users = ref([])
const loading = ref(true)
const showForm = ref(false)
const showMore = ref(false)
const editing = ref(null)
const tagInput = ref('')
const tagsList = ref([])
const creditsList = ref([])
const enableDtime = ref(false)
const dtimeStr = ref('')

const defaultForm = () => ({
  template_name: '', title: '', tid: '', copyright: 2, copyright_source: '',
  cover_path: '', description: '', dynamic: '', dtime: null,
  dolby: 0, hires: 0, charging_pay: 0, no_reprint: 0, is_only_self: 0,
  uploader: '', user_cookie: '',
  tags: [], credits: null,
  up_selection_reply: false, up_close_reply: false, up_close_danmu: false,
  extra_fields: ''
})

const form = ref(defaultForm())

const addTag = () => {
  const v = tagInput.value.trim()
  if (v && !tagsList.value.includes(v)) tagsList.value.push(v)
  tagInput.value = ''
}
const removeTag = (i) => tagsList.value.splice(i, 1)

const getUserName = (cookie) => {
  const u = users.value.find(u => u.value === cookie)
  return u ? (u.name || u.value) : cookie
}

const openAdd = () => {
  editing.value = null
  form.value = defaultForm()
  tagsList.value = []
  creditsList.value = []
  enableDtime.value = false
  dtimeStr.value = ''
  showMore.value = false
  showForm.value = true
}

const openEdit = (t) => {
  editing.value = t.id
  form.value = { ...defaultForm(), ...JSON.parse(JSON.stringify(t)) }
  tagsList.value = Array.isArray(t.tags) ? [...t.tags] : (t.tags || '').split(',').map(s => s.trim()).filter(Boolean)
  creditsList.value = Array.isArray(t.credits) ? t.credits.map(c => c.name || c) : []
  if (t.dtime) {
    enableDtime.value = true
    const d = new Date(t.dtime * 1000)
    dtimeStr.value = d.toISOString().slice(0, 16)
  } else {
    enableDtime.value = false
    dtimeStr.value = ''
  }
  showMore.value = false
  showForm.value = true
}

const buildPayload = () => {
  const p = { ...form.value }
  p.tags = tagsList.value
  p.tid = p.tid ? Number(p.tid) : null
  p.copyright = p.copyright ? Number(p.copyright) : null
  p.dolby = Number(p.dolby) || 0
  p.hires = Number(p.hires) || 0
  p.charging_pay = Number(p.charging_pay) || 0
  p.no_reprint = Number(p.no_reprint) || 0
  p.is_only_self = Number(p.is_only_self) || 0
  p.up_selection_reply = p.up_selection_reply ? 1 : 0
  p.up_close_reply = p.up_close_reply ? 1 : 0
  p.up_close_danmu = p.up_close_danmu ? 1 : 0
  if (enableDtime.value && dtimeStr.value) {
    p.dtime = Math.floor(new Date(dtimeStr.value).getTime() / 1000)
  } else { p.dtime = null }
  if (creditsList.value.filter(Boolean).length) { p.credits = creditsList.value.filter(Boolean) }
  else { p.credits = null }
  if (!p.user_cookie) p.user_cookie = null
  if (!p.copyright_source) p.copyright_source = null
  if (!p.cover_path) p.cover_path = null
  if (!p.description) p.description = null
  if (!p.dynamic) p.dynamic = null
  if (!p.extra_fields) p.extra_fields = null
  if (!p.uploader) p.uploader = null
  return p
}

const save = async () => {
  if (!form.value.template_name) { alert('请输入模板名称'); return }
  try {
    const payload = buildPayload()
    if (editing.value) payload.id = editing.value
    await api.addTemplate(payload)
    showForm.value = false
    await load()
  } catch (e) { alert('保存失败: ' + (e.response?.data?.message || e.message)) }
}

const del = async (id) => {
  if (!confirm('确定删除此模板？')) return
  try { await api.deleteTemplate(id); await load() }
  catch (e) { alert('删除失败: ' + (e.response?.data?.message || e.message)) }
}

const load = async () => {
  try { loading.value = true; const { data } = await api.getTemplates(); templates.value = data || [] }
  catch { templates.value = [] }
  finally { loading.value = false }
}

const loadUsers = async () => {
  try { const { data } = await api.getUsers(); users.value = data || [] }
  catch { users.value = [] }
}

onMounted(() => { load(); loadUsers() })
</script>

<style scoped>
/* ===== Card Grid ===== */
.grid {
  display:grid;
  grid-template-columns:repeat(auto-fill, minmax(min(400px, 100%), 1fr));
  gap:18px;
}

/* ===== Template Card ===== */
.tpl-card {
  background:#fff; border:1px solid rgba(0,0,0,.05); border-radius:18px;
  overflow:hidden; transition:all .3s var(--ease);
  box-shadow:0 1px 3px rgba(0,0,0,.03), 0 4px 12px rgba(0,0,0,.02);
}
.tpl-card:hover { box-shadow:0 4px 20px rgba(0,0,0,.06); border-color:rgba(0,0,0,.08); transform:translateY(-2px); }
.tpl-accent { height:3px; background:linear-gradient(90deg, var(--c-accent), rgba(196,138,156,.3), transparent); }
.tpl-header { display:flex; align-items:center; gap:14px; padding:20px 24px 0; }
.tpl-icon {
  width:42px; height:42px; border-radius:12px;
  background:linear-gradient(135deg, rgba(196,138,156,.12), rgba(196,138,156,.04));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center; flex-shrink:0;
}
.tpl-title-group { flex:1; min-width:0; display:flex; align-items:center; gap:8px; }
.tpl-name { font-size:16px; font-weight:700; color:var(--c-text-1); margin:0; overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.tpl-copyright-badge { flex-shrink:0; padding:2px 8px; border-radius:6px; font-size:10px; font-weight:700; letter-spacing:.3px; background:rgba(0,0,0,.04); color:var(--c-text-4); }
.tpl-copyright-badge.is-original { background:var(--c-success-soft); color:var(--c-success); }
.tpl-actions { display:flex; gap:2px; }
.tpl-action { width:32px; height:32px; border-radius:8px; border:none; cursor:pointer; background:transparent; color:var(--c-text-4); display:flex; align-items:center; justify-content:center; transition:all .15s; }
.tpl-action:hover { background:var(--c-accent-soft); color:var(--c-accent); }
.tpl-action-del:hover { background:var(--c-danger-soft); color:var(--c-danger); }
.tpl-title-preview { display:flex; align-items:center; gap:8px; margin:14px 24px 0; padding:10px 14px; background:var(--c-bg-input); border-radius:10px; font-size:13px; color:var(--c-text-2); overflow:hidden; white-space:nowrap; text-overflow:ellipsis; }
.tpl-title-preview svg { color:var(--c-text-4); flex-shrink:0; }
.tpl-body { padding:16px 24px 22px; }
.tpl-info-grid { display:grid; grid-template-columns:repeat(2, 1fr); gap:8px 16px; margin-bottom:12px; }
.tpl-info-item { display:flex; flex-direction:column; gap:2px; }
.tpl-info-label { font-size:10px; font-weight:700; color:var(--c-text-4); text-transform:uppercase; letter-spacing:.8px; }
.tpl-info-value { font-size:13px; color:var(--c-text-2); overflow:hidden; text-overflow:ellipsis; white-space:nowrap; }
.tpl-mono { font-family:var(--font-mono); font-size:11px; }
.tpl-tags { display:flex; flex-wrap:wrap; gap:6px; margin-bottom:12px; }
.tpl-tag { padding:4px 12px; border-radius:8px; font-size:11px; font-weight:600; background:linear-gradient(135deg, rgba(196,138,156,.1), rgba(196,138,156,.05)); color:var(--c-accent); }
.tpl-desc-block { padding:10px 14px; border-radius:10px; background:rgba(0,0,0,.015); border-left:3px solid rgba(196,138,156,.3); margin-bottom:12px; }
.tpl-desc { font-size:12px; color:var(--c-text-3); line-height:1.6; white-space:pre-wrap; max-height:52px; overflow:hidden; display:-webkit-box; -webkit-line-clamp:3; -webkit-box-orient:vertical; }
.tpl-features { display:flex; flex-wrap:wrap; gap:6px; }
.tpl-feat { display:inline-flex; align-items:center; gap:4px; padding:3px 10px; border-radius:6px; font-size:10px; font-weight:600; letter-spacing:.3px; background:var(--c-accent-soft); color:var(--c-accent); }
.tpl-feat.warn { background:rgba(212,148,10,.08); color:var(--c-warning); }
.tpl-feat.muted { background:rgba(0,0,0,.04); color:var(--c-text-4); }

/* ===== Full-page Form ===== */
.form-page-header {
  display:flex; align-items:center; gap:16px; margin-bottom:32px; flex-wrap:wrap;
}
.back-btn {
  width:40px; height:40px; border-radius:12px; border:1.5px solid rgba(0,0,0,.08);
  background:#fff; color:var(--c-text-2); cursor:pointer;
  display:flex; align-items:center; justify-content:center;
  transition:all .15s var(--ease); box-shadow:var(--shadow-sm);
}
.back-btn:hover { border-color:var(--c-accent); color:var(--c-accent); background:var(--c-accent-soft); }
.form-page-header .page-title { flex:1; margin:0; }
.form-page-actions { display:flex; gap:10px; }

/* Form card */
.form-card {
  background:#fff; border:1px solid rgba(0,0,0,.05); border-radius:18px;
  box-shadow:0 1px 3px rgba(0,0,0,.03); padding:0; overflow:hidden;
}

/* Section */
.sec-title {
  font-size:15px; font-weight:700; color:var(--c-text-1);
  padding:24px 36px 12px; margin:0;
  border-bottom:1px solid rgba(0,0,0,.04);
}
.sec-toggle {
  cursor:pointer; display:flex; align-items:center; gap:8px;
  user-select:none; transition:color .15s;
}
.sec-toggle:hover { color:var(--c-accent); }
.chev { transition:transform .25s var(--ease); margin-left:auto; }
.chev.open { transform:rotate(180deg); }

.sec-body { padding:8px 36px 24px; }
.sec-desc { font-size:13px; color:var(--c-text-4); line-height:1.8; margin-bottom:16px; }

/* Horizontal row layout (label left, control right) */
.row {
  display:flex; gap:20px; padding:18px 0;
  border-bottom:1px solid rgba(0,0,0,.03);
  align-items:flex-start;
}
.row:last-child { border-bottom:none; }
.row-label {
  width:120px; min-width:120px; flex-shrink:0;
  font-size:14px; font-weight:600; color:var(--c-text-1);
  padding-top:10px; line-height:1.4;
}
.row-label i { color:var(--c-danger); font-style:normal; }
.row-ctrl { flex:1; min-width:0; }

.row-hint {
  margin-top:8px; padding:10px 14px; border-radius:10px;
  background:rgba(0,0,0,.02); font-size:12px; color:var(--c-text-4); line-height:1.7;
}
.row-hint code {
  padding:1px 6px; border-radius:4px;
  background:rgba(196,138,156,.08); color:var(--c-accent);
  font-family:var(--font-mono); font-size:11px; font-weight:600;
}
.row-sub-hint { margin-top:6px; font-size:12px; color:var(--c-accent); }

/* Type selector */
.type-selector { display:flex; gap:16px; }
.type-opt {
  display:flex; align-items:center; gap:8px; cursor:pointer;
  font-size:14px; color:var(--c-text-3); transition:color .15s;
}
.type-opt input { accent-color:var(--c-accent); width:17px; height:17px; cursor:pointer; }
.type-opt.active { color:var(--c-text-1); font-weight:600; }

/* Tags */
.tags-input-wrap {
  display:flex; flex-wrap:wrap; gap:6px; align-items:center;
  padding:8px 14px; border-radius:12px;
  border:1.5px solid rgba(0,0,0,.05); background:var(--c-bg-input);
  min-height:44px; transition:all .2s var(--ease);
}
.tags-input-wrap:hover { border-color:var(--c-border-hover); background:rgba(255,255,255,.8); }
.tags-input-wrap:focus-within { border-color:var(--c-accent); box-shadow:0 0 0 4px var(--c-accent-glow); background:#fff; }
.tag-chip {
  display:flex; align-items:center; gap:3px;
  padding:4px 10px; border-radius:8px; font-size:12px; font-weight:600;
  background:var(--c-accent-soft); color:var(--c-accent);
}
.tag-x { background:none; border:none; color:var(--c-accent); cursor:pointer; font-size:14px; line-height:1; opacity:.5; transition:opacity .15s; padding:0 1px; }
.tag-x:hover { opacity:1; }
.tags-input { border:none; outline:none; background:transparent; font-size:14px; color:var(--c-text-1); flex:1; min-width:80px; }

/* Checkboxes */
.check-group { display:flex; gap:24px; flex-wrap:wrap; padding:4px 0; }
.ck { display:flex; align-items:center; gap:8px; font-size:14px; color:var(--c-text-2); cursor:pointer; }
.ck input { accent-color:var(--c-accent); width:16px; height:16px; cursor:pointer; }

/* Toggle switch */
.toggle { position:relative; display:inline-flex; cursor:pointer; }
.toggle input { position:absolute; opacity:0; }
.toggle-track {
  width:44px; height:24px; border-radius:12px;
  background:rgba(0,0,0,.1); transition:all .25s var(--ease);
  display:flex; align-items:center; padding:0 3px;
}
.toggle-thumb {
  width:18px; height:18px; border-radius:50%; background:#fff;
  box-shadow:0 1px 4px rgba(0,0,0,.15); transition:all .25s var(--ease);
}
.toggle input:checked + .toggle-track { background:var(--c-accent); }
.toggle input:checked + .toggle-track .toggle-thumb { transform:translateX(20px); }

.inline-row { display:flex; align-items:center; gap:14px; }

/* Textarea */
.form-textarea { resize:vertical; min-height:100px; font-family:inherit; line-height:1.6; }
.form-textarea-sm { min-height:70px; }
.char-count { display:block; text-align:right; margin-top:4px; font-size:11px; color:var(--c-text-4); font-variant-numeric:tabular-nums; }

/* Credits */
.credit-row { display:flex; gap:8px; margin-bottom:8px; align-items:center; }
.credit-row .form-input { flex:1; }
.credit-del {
  width:32px; height:32px; border-radius:8px; border:none;
  background:transparent; color:var(--c-text-4); cursor:pointer;
  display:flex; align-items:center; justify-content:center; transition:all .15s;
}
.credit-del:hover { background:var(--c-danger-soft); color:var(--c-danger); }
.add-row-btn {
  display:inline-flex; align-items:center; gap:6px;
  padding:8px 16px; border-radius:10px; border:1.5px dashed rgba(0,0,0,.1);
  background:transparent; font-size:12px; font-weight:600;
  color:var(--c-text-4); cursor:pointer; transition:all .15s;
}
.add-row-btn:hover { border-color:var(--c-accent); color:var(--c-accent); background:var(--c-accent-soft); }

@media (max-width:768px) {
  .grid { grid-template-columns:1fr; }
  .tpl-header { padding:16px 18px 0; }
  .tpl-title-preview { margin:12px 18px 0; }
  .tpl-body { padding:14px 18px 18px; }
  .tpl-info-grid { grid-template-columns:1fr; }
  .row { flex-direction:column; gap:6px; }
  .row-label { width:auto; min-width:auto; padding-top:0; }
  .sec-title { padding:20px 20px 10px; }
  .sec-body { padding:8px 20px 20px; }
  .form-page-actions { width:100%; }
  .form-page-actions .btn { flex:1; }
}
</style>
