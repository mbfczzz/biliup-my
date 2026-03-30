<template>
  <div class="page">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <p class="page-subtitle">配置录制参数、上传线路和全局行为</p>
    </div>

    <div v-if="loading" class="state-box"><div class="spinner"></div><p class="state-text">加载配置中...</p></div>

    <template v-else>
      <div class="config-sections">
        <div v-for="section in sections" :key="section.key" :class="['config-section card', { 'is-open': expanded.has(section.key) }]">
          <div class="section-head" @click="toggle(section.key)">
            <div class="section-left">
              <div class="section-icon" v-html="section.icon"></div>
              <div class="section-text">
                <span class="section-title">{{ section.title }}</span>
                <span class="section-desc">{{ section.desc }}</span>
              </div>
            </div>
            <svg :class="['chevron', { open: expanded.has(section.key) }]" width="16" height="16" viewBox="0 0 16 16" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M6 4l4 4-4 4"/></svg>
          </div>
          <div class="section-body">
            <div class="section-body-inner">
              <div class="section-body-content">
              <!-- 录制配置 -->
              <template v-if="section.key === 'record'">
                <div class="form-group">
                  <label class="form-label">下载器</label>
                  <select class="form-select" v-model="config.downloader">
                    <option :value="null">默认</option>
                    <option value="streamlink">Streamlink</option>
                    <option value="ffmpeg">FFmpeg</option>
                    <option value="stream-gears">Stream-gears</option>
                  </select>
                  <p class="form-hint">录制使用的下载器类型</p>
                </div>
                <div class="form-group">
                  <label class="form-label">文件大小限制 (字节)</label>
                  <input class="form-input" type="number" v-model.number="config.file_size" placeholder="8388608">
                  <p class="form-hint">单个录制文件的最大大小</p>
                </div>
                <div class="form-group">
                  <label class="form-label">分段时间</label>
                  <input class="form-input" v-model="config.segment_time" placeholder="00:00:00">
                  <p class="form-hint">录制分段时长，格式 HH:MM:SS</p>
                </div>
                <div class="form-group">
                  <label class="form-label">文件名前缀</label>
                  <input class="form-input" v-model="config.filename_prefix" placeholder="可选">
                </div>
              </template>
              <!-- 上传配置 -->
              <template v-if="section.key === 'upload'">
                <div class="form-group">
                  <label class="form-label">上传线路</label>
                  <select class="form-select" v-model="config.lines">
                    <option value="AUTO">自动</option>
                    <option value="alia">阿里云 (alia)</option>
                    <option value="bda2">百度 (bda2)</option>
                    <option value="bldsa">B站 (bldsa)</option>
                    <option value="qn">七牛 (qn)</option>
                    <option value="tx">腾讯 (tx)</option>
                    <option value="txa">腾讯加速 (txa)</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="form-label">上传线程数</label>
                  <input class="form-input" type="number" v-model.number="config.threads" min="1" max="16">
                </div>
                <div class="form-group">
                  <label class="form-label">上传器</label>
                  <select class="form-select" v-model="config.uploader">
                    <option :value="null">默认</option>
                    <option value="bili_web">bili_web</option>
                    <option value="Noop">Noop (不上传)</option>
                  </select>
                </div>
              </template>
              <!-- 高级配置 -->
              <template v-if="section.key === 'advanced'">
                <div class="form-group">
                  <label class="form-label">延迟时间 (秒)</label>
                  <input class="form-input" type="number" v-model.number="config.delay" placeholder="10">
                  <p class="form-hint">录制完成后等待上传的延迟</p>
                </div>
                <div class="form-group">
                  <label class="form-label">事件循环间隔 (秒)</label>
                  <input class="form-input" type="number" v-model.number="config.event_loop_interval" placeholder="40">
                </div>
                <div class="form-group">
                  <label class="form-label">检测间隔 (秒)</label>
                  <input class="form-input" type="number" v-model.number="config.checker_sleep" placeholder="120">
                  <p class="form-hint">检查直播状态的时间间隔</p>
                </div>
                <div class="form-group">
                  <label class="form-label">过滤阈值 (MB)</label>
                  <input class="form-input" type="number" v-model.number="config.filtering_threshold" placeholder="2">
                  <p class="form-hint">小于该大小的文件将被过滤</p>
                </div>
                <div class="form-group">
                  <label class="form-label">日志级别</label>
                  <select class="form-select" v-model="config.loggers_level">
                    <option :value="null">默认</option>
                    <option value="DEBUG">DEBUG</option>
                    <option value="INFO">INFO</option>
                    <option value="WARNING">WARNING</option>
                    <option value="ERROR">ERROR</option>
                  </select>
                </div>
                <div class="form-group">
                  <label class="toggle-wrap">
                    <span class="toggle-text">
                      <span class="toggle-label">使用直播封面</span>
                      <span class="toggle-hint">自动下载直播封面作为视频封面</span>
                    </span>
                    <label class="toggle">
                      <input type="checkbox" v-model="config.use_live_cover">
                      <span class="toggle-track"></span>
                    </label>
                  </label>
                </div>
              </template>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="action-bar">
        <button class="btn btn-ghost" @click="load">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><path d="M1 4v6h6"/><path d="M3.51 15a9 9 0 102.13-9.36L1 10"/></svg>
          重置
        </button>
        <button class="btn btn-primary" @click="save">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"><path d="M20 6L9 17l-5-5"/></svg>
          保存配置
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const loading = ref(true)
const config = ref({})
const expanded = ref(new Set())

const sections = [
  { key: 'record', title: '录制配置', desc: '下载器、分段和文件参数', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>' },
  { key: 'upload', title: '上传配置', desc: '上传线路、线程数和上传器', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M12 4v12M8 8l4-4 4 4"/></svg>' },
  { key: 'advanced', title: '高级设置', desc: '延迟、检测间隔和日志级别', icon: '<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>' },
]

const toggle = (key) => {
  const s = new Set(expanded.value)
  if (s.has(key)) s.delete(key); else s.add(key)
  expanded.value = s
}

const load = async () => {
  try { loading.value = true; const { data } = await api.getConfig(); if (data) config.value = data }
  catch (e) { console.error('加载配置失败:', e) }
  finally { loading.value = false }
}

const save = async () => {
  try { await api.updateConfig(config.value); alert('配置保存成功') }
  catch (e) { alert('保存失败: ' + (e.response?.data?.message || e.message)) }
}

onMounted(load)
</script>

<style scoped>
.config-sections { display:flex; flex-direction:column; gap:10px; margin-bottom:24px; }
.config-section { padding:0; overflow:hidden; }

.section-head {
  padding:20px 24px; cursor:pointer;
  display:flex; justify-content:space-between; align-items:center;
  transition:background .15s;
}
.section-head:hover { background:rgba(0,0,0,.015); }

.section-left { display:flex; align-items:center; gap:14px; }
.section-icon {
  width:40px; height:40px; border-radius:12px;
  background:linear-gradient(135deg, var(--c-accent-soft), rgba(196,138,156,.15));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
  flex-shrink:0;
}
.section-text { display:flex; flex-direction:column; gap:2px; }
.section-title { font-size:14px; font-weight:700; color:var(--c-text-1); }
.section-desc { font-size:12px; color:var(--c-text-4); }

.chevron { color:var(--c-text-4); transition:transform .2s var(--ease); flex-shrink:0; }
.chevron.open { transform:rotate(90deg); }

.section-body .form-group { margin-bottom:0; }

/* Toggle switch */
.toggle-wrap {
  display:flex; align-items:center; justify-content:space-between;
  cursor:pointer; gap:16px;
}
.toggle-text { display:flex; flex-direction:column; gap:2px; }
.toggle-label { font-size:14px; font-weight:600; color:var(--c-text-1); }
.toggle-hint { font-size:12px; color:var(--c-text-4); }

.toggle { position:relative; display:inline-block; width:44px; height:24px; flex-shrink:0; }
.toggle input { opacity:0; width:0; height:0; }
.toggle-track {
  position:absolute; inset:0; border-radius:12px;
  background:rgba(0,0,0,.12); cursor:pointer;
  transition:all .25s var(--ease);
}
.toggle-track::before {
  content:''; position:absolute; left:2px; top:2px;
  width:20px; height:20px; border-radius:50%;
  background:#fff; box-shadow:0 1px 3px rgba(0,0,0,.15);
  transition:all .25s var(--ease);
}
.toggle input:checked + .toggle-track {
  background:var(--c-accent);
}
.toggle input:checked + .toggle-track::before {
  transform:translateX(20px);
}

/* Smooth expand with grid */
.section-body {
  display:grid; grid-template-rows:0fr;
  transition:grid-template-rows .3s var(--ease);
}
.is-open .section-body { grid-template-rows:1fr; }
.section-body-inner { overflow:hidden; }
.section-body-content {
  padding:20px 24px 24px;
  display:flex; flex-direction:column; gap:18px;
  border-top:1px solid var(--c-border);
  margin:0 16px;
}

.action-bar { display:flex; gap:10px; justify-content:flex-end; }

.form-select {
  appearance:none;
  background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12' fill='%23999'%3E%3Cpath d='M3 5l3 3 3-3'/%3E%3C/svg%3E");
  background-repeat:no-repeat;
  background-position:right 14px center;
  padding-right:36px;
}
</style>
