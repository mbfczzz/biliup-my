<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-left">
        <h1 class="logo">Biliup</h1>
        <span class="subtitle">Twitch 录制 → B站上传</span>
      </div>
      <div class="header-right">
        <el-tag :type="statusType" effect="dark" size="large">
          {{ statusText }}
        </el-tag>
      </div>
    </el-header>

    <el-container>
      <el-aside width="200px" class="app-aside">
        <el-menu :default-active="$route.path" router>
          <el-menu-item index="/">
            <span>录制管理</span>
          </el-menu-item>
          <el-menu-item index="/templates">
            <span>上传模板</span>
          </el-menu-item>
          <el-menu-item index="/config">
            <span>系统配置</span>
          </el-menu-item>
          <el-menu-item index="/account">
            <span>B站账号</span>
          </el-menu-item>
          <el-menu-item index="/logs">
            <span>运行日志</span>
          </el-menu-item>
        </el-menu>
      </el-aside>

      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import api from './api'

export default {
  name: 'App',
  setup() {
    const status = ref(null)
    const statusText = ref('连接中...')
    const statusType = ref('info')
    let timer = null

    const fetchStatus = async () => {
      try {
        const res = await api.getStatus()
        status.value = res.data
        statusText.value = `运行 ${res.data.uptime} | 任务 ${res.data.active_tasks}`
        statusType.value = res.data.active_tasks > 0 ? 'success' : 'info'
      } catch {
        statusText.value = '连接失败'
        statusType.value = 'danger'
      }
    }

    onMounted(() => {
      fetchStatus()
      timer = setInterval(fetchStatus, 10000)
    })

    onUnmounted(() => {
      if (timer) clearInterval(timer)
    })

    return { status, statusText, statusType }
  },
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #f5f7fa;
}

.app-container {
  min-height: 100vh;
}

.app-header {
  background: #1a1a2e;
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  height: 60px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.logo {
  font-size: 22px;
  font-weight: bold;
  color: #00d4ff;
}

.subtitle {
  font-size: 13px;
  color: #aaa;
}

.app-aside {
  background: white;
  border-right: 1px solid #e4e7ed;
  padding-top: 8px;
}

.app-aside .el-menu {
  border-right: none;
}

.app-main {
  padding: 20px;
  background: #f5f7fa;
}
</style>
