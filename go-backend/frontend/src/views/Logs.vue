<template>
  <div class="page">
    <h2 style="margin-bottom: 20px;">运行日志</h2>

    <el-card class="log-card">
      <template #header>
        <div class="log-header">
          <span>实时日志</span>
          <div>
            <el-tag :type="wsConnected ? 'success' : 'danger'" size="small" effect="dark">
              {{ wsConnected ? '已连接' : '未连接' }}
            </el-tag>
            <el-button size="small" @click="clearLogs" style="margin-left: 8px">清空</el-button>
          </div>
        </div>
      </template>

      <div class="log-container" ref="logContainer">
        <div v-for="(log, index) in logs" :key="index" class="log-line" :class="'log-' + log.level">
          <span class="log-time">{{ formatTime(log.time) }}</span>
          <span class="log-level">[{{ log.level.toUpperCase() }}]</span>
          <span class="log-msg">{{ log.message }}</span>
        </div>
        <div v-if="logs.length === 0" class="log-empty">暂无日志</div>
      </div>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

export default {
  name: 'Logs',
  setup() {
    const logs = ref([])
    const logContainer = ref(null)
    const wsConnected = ref(false)
    let ws = null

    const connectWS = () => {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = window.location.host
      ws = new WebSocket(`${protocol}//${host}/v1/logs`)

      ws.onopen = () => {
        wsConnected.value = true
      }

      ws.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data)
          logs.value.push(msg)

          // 限制最多显示500条
          if (logs.value.length > 500) {
            logs.value = logs.value.slice(-500)
          }

          // 自动滚动到底部
          nextTick(() => {
            if (logContainer.value) {
              logContainer.value.scrollTop = logContainer.value.scrollHeight
            }
          })
        } catch {}
      }

      ws.onclose = () => {
        wsConnected.value = false
        // 5秒后重连
        setTimeout(connectWS, 5000)
      }

      ws.onerror = () => {
        wsConnected.value = false
      }
    }

    const formatTime = (time) => {
      if (!time) return ''
      const d = new Date(time)
      return d.toLocaleTimeString('zh-CN', { hour12: false })
    }

    const clearLogs = () => {
      logs.value = []
    }

    onMounted(connectWS)
    onUnmounted(() => {
      if (ws) ws.close()
    })

    return { logs, logContainer, wsConnected, formatTime, clearLogs }
  },
}
</script>

<style scoped>
.page { max-width: 1200px; }
.log-card { height: calc(100vh - 160px); }
.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.log-container {
  height: calc(100vh - 260px);
  overflow-y: auto;
  background: #1e1e1e;
  border-radius: 4px;
  padding: 12px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
}
.log-line { white-space: pre-wrap; word-break: break-all; }
.log-time { color: #6a9955; margin-right: 8px; }
.log-level { margin-right: 8px; font-weight: bold; }
.log-msg { color: #d4d4d4; }
.log-info .log-level { color: #4fc1ff; }
.log-error .log-level { color: #f44747; }
.log-error .log-msg { color: #f44747; }
.log-warn .log-level { color: #cca700; }
.log-empty { color: #666; text-align: center; padding: 40px; }
</style>
