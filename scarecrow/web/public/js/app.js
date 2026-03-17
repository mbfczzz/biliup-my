const { createApp, ref, reactive, computed, onMounted, onUnmounted, watch, nextTick } = Vue;

// ── API helpers ──────────────────────────────────────────────────────────────
const BASE = '';

async function apiFetch(method, path, body) {
  const opts = { method, headers: { 'Content-Type': 'application/json' } };
  if (body !== undefined) opts.body = JSON.stringify(body);
  const res = await fetch(BASE + path, opts);
  if (!res.ok) {
    const text = await res.text().catch(() => res.statusText);
    throw new Error(text || `HTTP ${res.status}`);
  }
  const ct = res.headers.get('content-type') || '';
  if (ct.includes('application/json')) return res.json();
  return null;
}

const api = {
  getStreamers: () => apiFetch('GET', '/v1/streamers'),
  addStreamer: (data) => apiFetch('POST', '/v1/streamers', data),
  updateStreamer: (id, data) => apiFetch('PUT', `/v1/streamers/${id}`, data),
  deleteStreamer: (id) => apiFetch('DELETE', `/v1/streamers/${id}`),
  pauseStreamer: (id) => apiFetch('PUT', `/v1/streamers/${id}/pause`),
  getConfig: () => apiFetch('GET', '/v1/configuration'),
  saveConfig: (data) => apiFetch('PUT', '/v1/configuration', data),
  getStatus: () => apiFetch('GET', '/v1/status'),
  getHistory: () => apiFetch('GET', '/v1/streamer-info'),
  getHistoryFiles: (id) => apiFetch('GET', `/v1/streamer-info/files/${id}`),
};

// ── Status helpers ────────────────────────────────────────────────────────────
function getStatusBadge(status) {
  const map = {
    'Streaming': { cls: 'badge-live', icon: '🔴', label: '直播中' },
    'Idle': { cls: 'badge-idle', icon: '⚪', label: '空闲' },
    'Checking': { cls: 'badge-checking', icon: '🔵', label: '检测中' },
    'TimeRange': { cls: 'badge-timerange', icon: '🟣', label: '非录播时间' },
    'Paused': { cls: 'badge-paused', icon: '🟠', label: '已暂停' },
  };
  return map[status] || { cls: 'badge-idle', icon: '⚪', label: status || '未知' };
}

function getInitials(name) {
  if (!name) return '?';
  const words = name.trim().split(/\s+/);
  if (words.length >= 2) return (words[0][0] + words[1][0]).toUpperCase();
  return name.slice(0, 2).toUpperCase();
}

// ── Toast ────────────────────────────────────────────────────────────────────
const toasts = ref([]);
let toastId = 0;

function showToast(message, type = 'success', duration = 3000) {
  const id = ++toastId;
  toasts.value.push({ id, message, type });
  setTimeout(() => {
    const el = document.querySelector(`[data-toast="${id}"]`);
    if (el) el.classList.add('toast-leave');
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id);
    }, 260);
  }, duration);
}

// ── Main App ─────────────────────────────────────────────────────────────────
createApp({
  setup() {
    // Navigation
    const currentPage = ref('streamers');

    // Streamer list
    const streamers = ref([]);
    const streamersLoading = ref(false);
    const streamersError = ref('');

    // Streamer modal
    const showStreamerModal = ref(false);
    const editingStreamer = ref(null); // null = add mode
    const streamerForm = reactive({
      url: '',
      remark: '',
      segment_time: '01:00:00',
      file_size: '',
      excluded_keywords: '',
      time_range: '',
      segment_processor: '',
    });
    const streamerSaving = ref(false);

    // Settings modal
    const showSettingsModal = ref(false);
    const settingsForm = reactive({
      downloader: 'stream-gears',
      segment_time: '01:00:00',
      file_size: '',
      twitch_cookie: '',
      twitch_disable_ads: true,
      loggers_level: 'INFO',
    });
    const settingsSaving = ref(false);

    // History
    const history = ref([]);
    const historyLoading = ref(false);
    const expandedHistory = ref({});

    // Logs
    const logLines = ref([]);
    const logsAutoScroll = ref(true);
    const logsConnected = ref(false);
    let logWs = null;

    // Status
    const statusData = ref(null);
    const statusLoading = ref(false);

    // Polling
    let pollTimer = null;

    // ── Streamers ─────────────────────────────────────────────────────────────
    async function loadStreamers() {
      streamersLoading.value = true;
      streamersError.value = '';
      try {
        const data = await api.getStreamers();
        streamers.value = Array.isArray(data) ? data : (data?.data || []);
      } catch (e) {
        streamersError.value = e.message;
      } finally {
        streamersLoading.value = false;
      }
    }

    function openAddStreamer() {
      editingStreamer.value = null;
      Object.assign(streamerForm, {
        url: '', remark: '', segment_time: '01:00:00',
        file_size: '', excluded_keywords: '', time_range: '', segment_processor: '',
      });
      showStreamerModal.value = true;
    }

    function openEditStreamer(s) {
      editingStreamer.value = s;
      Object.assign(streamerForm, {
        url: s.url || '',
        remark: s.remark || s.name || '',
        segment_time: s.segment_time || '01:00:00',
        file_size: s.file_size || '',
        excluded_keywords: Array.isArray(s.excluded_keywords)
          ? s.excluded_keywords.join(', ')
          : (s.excluded_keywords || ''),
        time_range: s.time_range || '',
        segment_processor: s.segment_processor || '',
      });
      showStreamerModal.value = true;
    }

    async function saveStreamer() {
      if (!streamerForm.url.trim()) {
        showToast('请输入 Twitch 频道 URL', 'error');
        return;
      }
      streamerSaving.value = true;
      try {
        const payload = {
          url: streamerForm.url.trim(),
          remark: streamerForm.remark.trim() || undefined,
          segment_time: streamerForm.segment_time || undefined,
          file_size: streamerForm.file_size ? Number(streamerForm.file_size) : undefined,
          excluded_keywords: streamerForm.excluded_keywords
            ? streamerForm.excluded_keywords.split(',').map(s => s.trim()).filter(Boolean)
            : undefined,
          time_range: streamerForm.time_range.trim() || undefined,
          segment_processor: streamerForm.segment_processor.trim() || undefined,
        };
        if (editingStreamer.value?.id !== undefined) {
          await api.updateStreamer(editingStreamer.value.id, payload);
          showToast('频道已更新');
        } else {
          await api.addStreamer(payload);
          showToast('频道已添加');
        }
        showStreamerModal.value = false;
        await loadStreamers();
      } catch (e) {
        showToast(e.message, 'error');
      } finally {
        streamerSaving.value = false;
      }
    }

    async function deleteStreamer(s) {
      if (!confirm(`确认删除 "${s.remark || s.name || s.url}"？`)) return;
      try {
        await api.deleteStreamer(s.id);
        showToast('已删除');
        await loadStreamers();
      } catch (e) {
        showToast(e.message, 'error');
      }
    }

    async function togglePause(s) {
      try {
        await api.pauseStreamer(s.id);
        showToast(s.status === 'Paused' ? '已恢复' : '已暂停');
        await loadStreamers();
      } catch (e) {
        showToast(e.message, 'error');
      }
    }

    // ── Settings ──────────────────────────────────────────────────────────────
    async function openSettings() {
      try {
        const cfg = await api.getConfig();
        if (cfg) {
          settingsForm.downloader = cfg.downloader || 'stream-gears';
          settingsForm.segment_time = cfg.segment_time || '01:00:00';
          settingsForm.file_size = cfg.file_size || '';
          settingsForm.twitch_cookie = cfg.user?.twitch_cookie || '';
          settingsForm.twitch_disable_ads = cfg.twitch_disable_ads !== false;
          settingsForm.loggers_level = cfg.loggers_level || 'INFO';
        }
      } catch (e) {
        showToast('加载配置失败: ' + e.message, 'error');
      }
      showSettingsModal.value = true;
    }

    async function saveSettings() {
      settingsSaving.value = true;
      try {
        const payload = {
          downloader: settingsForm.downloader,
          segment_time: settingsForm.segment_time,
          file_size: settingsForm.file_size ? Number(settingsForm.file_size) : undefined,
          twitch_disable_ads: settingsForm.twitch_disable_ads,
          loggers_level: settingsForm.loggers_level,
          user: settingsForm.twitch_cookie
            ? { twitch_cookie: settingsForm.twitch_cookie }
            : undefined,
        };
        await api.saveConfig(payload);
        showToast('设置已保存');
        showSettingsModal.value = false;
      } catch (e) {
        showToast(e.message, 'error');
      } finally {
        settingsSaving.value = false;
      }
    }

    // ── History ───────────────────────────────────────────────────────────────
    async function loadHistory() {
      historyLoading.value = true;
      try {
        const data = await api.getHistory();
        history.value = Array.isArray(data) ? data : (data?.data || []);
      } catch (e) {
        showToast('加载历史失败: ' + e.message, 'error');
      } finally {
        historyLoading.value = false;
      }
    }

    async function toggleHistoryRow(item) {
      const id = item.id;
      if (expandedHistory.value[id]) {
        expandedHistory.value[id] = null;
        return;
      }
      try {
        const files = await api.getHistoryFiles(id);
        expandedHistory.value[id] = files;
      } catch (e) {
        showToast('加载文件列表失败: ' + e.message, 'error');
      }
    }

    // ── Logs ──────────────────────────────────────────────────────────────────
    function connectLogs() {
      if (logWs && logWs.readyState < 2) return;
      const wsUrl = (location.protocol === 'https:' ? 'wss:' : 'ws:')
        + '//' + location.host + '/v1/ws/logs?file=download.log';
      logWs = new WebSocket(wsUrl);
      logWs.onopen = () => { logsConnected.value = true; };
      logWs.onclose = () => { logsConnected.value = false; };
      logWs.onerror = () => { logsConnected.value = false; };
      logWs.onmessage = (e) => {
        const line = e.data || '';
        let cls = 'default';
        const low = line.toLowerCase();
        if (low.includes('[error]') || low.includes('error:')) cls = 'error';
        else if (low.includes('[warn]') || low.includes('warning:')) cls = 'warn';
        else if (low.includes('[info]') || low.includes('info:')) cls = 'info';
        else if (low.includes('[debug]')) cls = 'debug';
        logLines.value.push({ text: line, cls });
        if (logLines.value.length > 2000) logLines.value.shift();
        if (logsAutoScroll.value) {
          nextTick(() => {
            const el = document.getElementById('log-container');
            if (el) el.scrollTop = el.scrollHeight;
          });
        }
      };
    }

    function disconnectLogs() {
      if (logWs) { logWs.close(); logWs = null; }
    }

    function clearLogs() { logLines.value = []; }

    function downloadLogs() {
      const text = logLines.value.map(l => l.text).join('\n');
      const a = document.createElement('a');
      a.href = 'data:text/plain;charset=utf-8,' + encodeURIComponent(text);
      a.download = 'scarecrow-download.log';
      a.click();
    }

    // ── Status ────────────────────────────────────────────────────────────────
    async function loadStatus() {
      statusLoading.value = true;
      try {
        statusData.value = await api.getStatus();
      } catch (e) {
        statusData.value = null;
      } finally {
        statusLoading.value = false;
      }
    }

    // ── Page switching ────────────────────────────────────────────────────────
    function navigate(page) {
      if (currentPage.value === 'logs' && page !== 'logs') disconnectLogs();
      currentPage.value = page;
      if (page === 'streamers') loadStreamers();
      if (page === 'history') loadHistory();
      if (page === 'logs') connectLogs();
      if (page === 'status') loadStatus();
    }

    // ── Polling ───────────────────────────────────────────────────────────────
    function startPolling() {
      pollTimer = setInterval(async () => {
        if (currentPage.value === 'streamers') await loadStreamers();
        if (currentPage.value === 'status') await loadStatus();
      }, 5000);
    }

    // ── Lifecycle ─────────────────────────────────────────────────────────────
    onMounted(() => {
      loadStreamers();
      startPolling();
    });

    onUnmounted(() => {
      clearInterval(pollTimer);
      disconnectLogs();
    });

    return {
      currentPage, navigate,
      streamers, streamersLoading, streamersError, loadStreamers,
      showStreamerModal, editingStreamer, streamerForm, streamerSaving,
      openAddStreamer, openEditStreamer, saveStreamer, deleteStreamer, togglePause,
      showSettingsModal, settingsForm, settingsSaving, openSettings, saveSettings,
      history, historyLoading, expandedHistory, toggleHistoryRow,
      logLines, logsAutoScroll, logsConnected, clearLogs, downloadLogs,
      statusData, statusLoading,
      toasts, showToast,
      getStatusBadge, getInitials,
    };
  },

  template: `
<div id="app-root">
  <!-- Navigation -->
  <nav class="nav">
    <a class="nav-brand" href="#">
      <div class="brand-icon">📡</div>
      稻草人的SCARECROW
    </a>
    <div class="nav-tabs">
      <button class="nav-tab" :class="{ active: currentPage === 'streamers' }" @click="navigate('streamers')">
        <span>📺 频道</span>
      </button>
      <button class="nav-tab" :class="{ active: currentPage === 'history' }" @click="navigate('history')">
        <span>🗂 历史</span>
      </button>
      <button class="nav-tab" :class="{ active: currentPage === 'logs' }" @click="navigate('logs')">
        <span>📋 日志</span>
      </button>
      <button class="nav-tab" :class="{ active: currentPage === 'status' }" @click="navigate('status')">
        <span>📊 状态</span>
      </button>
    </div>
    <div class="nav-actions">
      <button class="btn btn-secondary btn-sm" @click="openSettings">⚙️ 设置</button>
    </div>
  </nav>

  <!-- Main -->
  <main class="main">

    <!-- ── Streamers Page ── -->
    <div v-if="currentPage === 'streamers'">
      <div class="page-header">
        <div>
          <div class="page-title">录制频道</div>
          <div class="page-subtitle">管理 Twitch 直播录制</div>
        </div>
        <button class="btn btn-primary" @click="openAddStreamer">+ 添加频道</button>
      </div>

      <div v-if="streamersLoading && !streamers.length" style="text-align:center;padding:40px">
        <div class="spinner"></div>
      </div>
      <div v-else-if="streamersError" class="empty-state">
        <div class="empty-icon">⚠️</div>
        <div class="empty-title">加载失败</div>
        <div class="empty-desc">{{ streamersError }}</div>
        <button class="btn btn-primary" @click="loadStreamers">重试</button>
      </div>
      <div v-else-if="!streamers.length" class="empty-state">
        <div class="empty-icon">📡</div>
        <div class="empty-title">还没有录制频道</div>
        <div class="empty-desc">添加 Twitch 频道开始录制</div>
        <button class="btn btn-primary" @click="openAddStreamer">+ 添加频道</button>
      </div>
      <div v-else class="streamers-grid">
        <div v-for="s in streamers" :key="s.id" class="card streamer-card">
          <div class="streamer-card-header">
            <div class="streamer-avatar">
              <img v-if="s.cover" :src="s.cover" :alt="s.name" @error="$event.target.style.display='none'">
              <span>{{ getInitials(s.remark || s.name) }}</span>
            </div>
            <div class="streamer-info">
              <div class="streamer-name">{{ s.remark || s.name || '未命名' }}</div>
              <div class="streamer-url">{{ s.url }}</div>
            </div>
            <div>
              <span class="badge" :class="getStatusBadge(s.status).cls">
                <span class="badge-dot"></span>
                {{ getStatusBadge(s.status).label }}
              </span>
            </div>
          </div>
          <div class="streamer-card-body" v-if="s.room_title">
            <div style="font-size:13px;color:var(--secondary-label)">{{ s.room_title }}</div>
          </div>
          <div class="streamer-card-footer">
            <button class="btn btn-ghost btn-sm" @click="togglePause(s)">
              {{ s.status === 'Paused' ? '▶ 恢复' : '⏸ 暂停' }}
            </button>
            <button class="btn btn-ghost btn-sm" @click="openEditStreamer(s)">✏️ 编辑</button>
            <button class="btn btn-ghost btn-sm" style="color:var(--red);margin-left:auto" @click="deleteStreamer(s)">🗑 删除</button>
          </div>
        </div>
      </div>
    </div>

    <!-- ── History Page ── -->
    <div v-if="currentPage === 'history'">
      <div class="page-header">
        <div>
          <div class="page-title">录制历史</div>
          <div class="page-subtitle">查看历史录制记录</div>
        </div>
        <button class="btn btn-secondary btn-sm" @click="loadHistory">🔄 刷新</button>
      </div>
      <div v-if="historyLoading" style="text-align:center;padding:40px"><div class="spinner"></div></div>
      <div v-else-if="!history.length" class="empty-state">
        <div class="empty-icon">🗂</div>
        <div class="empty-title">暂无历史记录</div>
      </div>
      <div v-else class="card">
        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>主播</th>
                <th>标题</th>
                <th>时间</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              <template v-for="item in history" :key="item.id">
                <tr>
                  <td><strong>{{ item.name || item.streamer }}</strong></td>
                  <td style="max-width:300px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis">
                    {{ item.title || '—' }}
                  </td>
                  <td style="white-space:nowrap;color:var(--secondary-label);font-size:13px">
                    {{ item.time || item.created_at || '—' }}
                  </td>
                  <td>
                    <button class="btn btn-ghost btn-sm" @click="toggleHistoryRow(item)">
                      {{ expandedHistory[item.id] ? '▲ 收起' : '▼ 文件' }}
                    </button>
                  </td>
                </tr>
                <tr v-if="expandedHistory[item.id]">
                  <td colspan="4" style="padding:0">
                    <div style="padding:12px 16px;background:var(--gray6)">
                      <div v-if="!expandedHistory[item.id].length" style="color:var(--secondary-label);font-size:13px">无文件记录</div>
                      <div v-for="f in expandedHistory[item.id]" :key="f.id || f.file_name"
                           style="font-size:13px;padding:4px 0;border-bottom:1px solid var(--separator)">
                        📄 {{ f.file_name || f.name }}
                        <span style="color:var(--secondary-label);margin-left:8px">{{ f.size ? (f.size / 1024 / 1024).toFixed(1) + ' MB' : '' }}</span>
                      </div>
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- ── Logs Page ── -->
    <div v-if="currentPage === 'logs'">
      <div class="page-header">
        <div>
          <div class="page-title">实时日志</div>
          <div class="page-subtitle">
            <span v-if="logsConnected" style="color:var(--green)">● 已连接</span>
            <span v-else style="color:var(--red)">● 未连接</span>
          </div>
        </div>
        <div style="display:flex;gap:8px;align-items:center">
          <label class="toggle-group" style="font-size:13px;gap:8px">
            <span>自动滚动</span>
            <span class="toggle">
              <input type="checkbox" v-model="logsAutoScroll">
              <span class="toggle-slider"></span>
            </span>
          </label>
          <button class="btn btn-secondary btn-sm" @click="clearLogs">🗑 清空</button>
          <button class="btn btn-secondary btn-sm" @click="downloadLogs">⬇ 下载</button>
        </div>
      </div>
      <div class="log-container" id="log-container">
        <div v-if="!logLines.length" style="color:#555;font-style:italic">等待日志...</div>
        <div v-for="(line, i) in logLines" :key="i" class="log-line" :class="line.cls">{{ line.text }}</div>
      </div>
    </div>

    <!-- ── Status Page ── -->
    <div v-if="currentPage === 'status'">
      <div class="page-header">
        <div>
          <div class="page-title">系统状态</div>
        </div>
        <button class="btn btn-secondary btn-sm" @click="loadStatus">🔄 刷新</button>
      </div>
      <div v-if="statusLoading" style="text-align:center;padding:40px"><div class="spinner"></div></div>
      <div v-else-if="!statusData" class="empty-state">
        <div class="empty-icon">📊</div>
        <div class="empty-title">无法获取状态</div>
      </div>
      <div v-else>
        <div class="status-grid">
          <div class="card status-card">
            <div class="status-value" style="color:var(--blue)">{{ statusData.version || '—' }}</div>
            <div class="status-label">版本</div>
          </div>
          <div class="card status-card">
            <div class="status-value" style="color:var(--green)">{{ statusData.recording || statusData.active_recordings || 0 }}</div>
            <div class="status-label">录制中</div>
          </div>
          <div class="card status-card">
            <div class="status-value" style="color:var(--orange)">{{ statusData.total || statusData.total_streamers || 0 }}</div>
            <div class="status-label">频道总数</div>
          </div>
          <div class="card status-card">
            <div class="status-value" style="color:var(--purple)">{{ statusData.downloader || '—' }}</div>
            <div class="status-label">下载器</div>
          </div>
        </div>
        <div class="card" style="margin-top:16px" v-if="statusData.streamers?.length">
          <div class="card-body">
            <div class="section-title">当前录制中的频道</div>
            <div v-for="s in statusData.streamers" :key="s.id || s.name" style="padding:8px 0;border-bottom:1px solid var(--separator);font-size:14px">
              🔴 {{ s.name || s.url }}
              <span style="color:var(--secondary-label);font-size:12px;margin-left:8px">{{ s.title || '' }}</span>
            </div>
          </div>
        </div>
        <div class="card" style="margin-top:16px" v-if="statusData">
          <div class="card-body">
            <div class="section-title">原始数据</div>
            <pre style="font-size:12px;color:var(--secondary-label);overflow-x:auto;white-space:pre-wrap">{{ JSON.stringify(statusData, null, 2) }}</pre>
          </div>
        </div>
      </div>
    </div>

  </main>

  <!-- ── Add/Edit Streamer Modal ── -->
  <div class="modal-overlay" v-if="showStreamerModal" @click.self="showStreamerModal=false">
    <div class="modal">
      <div class="modal-header">
        <div class="modal-title">{{ editingStreamer ? '编辑频道' : '添加频道' }}</div>
        <button class="modal-close" @click="showStreamerModal=false">×</button>
      </div>
      <div class="modal-body">
        <div class="form-group">
          <label class="form-label">Twitch 频道 URL *</label>
          <input class="form-control" v-model="streamerForm.url" placeholder="https://www.twitch.tv/channelname" />
        </div>
        <div class="form-group">
          <label class="form-label">备注名称</label>
          <input class="form-control" v-model="streamerForm.remark" placeholder="可选，显示用名称" />
        </div>
        <div class="form-group">
          <label class="form-label">分段时长</label>
          <input class="form-control" v-model="streamerForm.segment_time" placeholder="01:00:00" />
        </div>
        <div class="form-group">
          <label class="form-label">最大文件大小 (字节，可选)</label>
          <input class="form-control" type="number" v-model="streamerForm.file_size" placeholder="留空表示不限" />
        </div>
        <div class="form-group">
          <label class="form-label">排除关键词（英文逗号分隔）</label>
          <input class="form-control" v-model="streamerForm.excluded_keywords" placeholder="广告, 测试" />
        </div>
        <div class="form-group">
          <label class="form-label">录制时段（JSON格式，可选）</label>
          <input class="form-control" v-model="streamerForm.time_range" placeholder='["08:00","22:00"]' />
        </div>
        <div class="form-group">
          <label class="form-label">分段后处理器（可选）</label>
          <input class="form-control" v-model="streamerForm.segment_processor" placeholder="处理脚本路径" />
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="showStreamerModal=false">取消</button>
        <button class="btn btn-primary" @click="saveStreamer" :disabled="streamerSaving">
          <span v-if="streamerSaving" class="spinner" style="width:14px;height:14px;border-width:2px"></span>
          {{ editingStreamer ? '保存' : '添加' }}
        </button>
      </div>
    </div>
  </div>

  <!-- ── Settings Modal ── -->
  <div class="modal-overlay" v-if="showSettingsModal" @click.self="showSettingsModal=false">
    <div class="modal modal-wide">
      <div class="modal-header">
        <div class="modal-title">⚙️ 录制设置</div>
        <button class="modal-close" @click="showSettingsModal=false">×</button>
      </div>
      <div class="modal-body">
        <div class="section-title">下载器</div>
        <div class="form-group">
          <label class="form-label">下载器选择</label>
          <select class="form-control" v-model="settingsForm.downloader">
            <option value="stream-gears">stream-gears（推荐）</option>
            <option value="ffmpeg">ffmpeg</option>
            <option value="streamlink">streamlink</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">默认分段时长</label>
          <input class="form-control" v-model="settingsForm.segment_time" placeholder="01:00:00" />
        </div>
        <div class="form-group">
          <label class="form-label">默认最大文件大小 (字节，可选)</label>
          <input class="form-control" type="number" v-model="settingsForm.file_size" placeholder="留空表示不限" />
        </div>
        <div class="divider"></div>
        <div class="section-title">Twitch</div>
        <div class="form-group">
          <label class="form-label">Twitch Cookie (auth-token)</label>
          <input class="form-control" v-model="settingsForm.twitch_cookie" placeholder="可选，用于订阅内容" type="password" />
        </div>
        <div class="form-group">
          <div class="toggle-group">
            <span class="toggle-label">禁用广告 (twitch-disable-ads)</span>
            <label class="toggle">
              <input type="checkbox" v-model="settingsForm.twitch_disable_ads">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
        <div class="divider"></div>
        <div class="section-title">日志</div>
        <div class="form-group">
          <label class="form-label">日志级别</label>
          <select class="form-control" v-model="settingsForm.loggers_level">
            <option value="DEBUG">DEBUG</option>
            <option value="INFO">INFO</option>
            <option value="WARNING">WARNING</option>
            <option value="ERROR">ERROR</option>
          </select>
        </div>
      </div>
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="showSettingsModal=false">取消</button>
        <button class="btn btn-primary" @click="saveSettings" :disabled="settingsSaving">
          <span v-if="settingsSaving" class="spinner" style="width:14px;height:14px;border-width:2px"></span>
          保存设置
        </button>
      </div>
    </div>
  </div>

  <!-- ── Toast Container ── -->
  <div class="toast-container">
    <div
      v-for="t in toasts" :key="t.id"
      :data-toast="t.id"
      class="toast"
      :class="'toast-' + t.type"
    >
      <span v-if="t.type==='success'">✓</span>
      <span v-else-if="t.type==='error'">✕</span>
      <span v-else>ℹ</span>
      {{ t.message }}
    </div>
  </div>
</div>
  `,
}).mount('#app');
