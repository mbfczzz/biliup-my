<template>
  <div class="page">
    <div class="page-header">
      <h2>B站账号</h2>
      <el-button type="primary" @click="startLogin">扫码登录</el-button>
    </div>

    <!-- 已登录的账号列表 -->
    <el-table :data="users" border stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="username" label="用户名" width="150" />
      <el-table-column prop="nickname" label="昵称" width="150" />
      <el-table-column prop="uid" label="UID" width="120" />
      <el-table-column prop="created_at" label="添加时间" width="180" />
      <el-table-column label="操作" width="100">
        <template #default="{ row }">
          <el-button size="small" type="danger" @click="deleteUser(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 扫码登录对话框 -->
    <el-dialog v-model="showLoginDialog" title="扫码登录B站" width="400px" :close-on-click-modal="false">
      <div class="qr-container" v-loading="qrLoading">
        <div v-if="qrImage" class="qr-wrapper">
          <img :src="qrImage" alt="登录二维码" class="qr-image" />
          <p class="qr-hint">请使用B站手机App扫描二维码</p>
          <p class="qr-status" :class="loginStatusClass">{{ loginStatusText }}</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="showLoginDialog = false">取消</el-button>
        <el-button @click="startLogin">刷新二维码</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

export default {
  name: 'Account',
  setup() {
    const users = ref([])
    const showLoginDialog = ref(false)
    const qrImage = ref('')
    const qrCodeKey = ref('')
    const qrLoading = ref(false)
    const loginStatusText = ref('等待扫码...')
    const loginStatusClass = ref('')
    let pollTimer = null

    const fetchUsers = async () => {
      try {
        const res = await api.getUsers()
        users.value = res.data || []
      } catch {
        ElMessage.error('加载账号列表失败')
      }
    }

    const startLogin = async () => {
      showLoginDialog.value = true
      qrLoading.value = true
      loginStatusText.value = '等待扫码...'
      loginStatusClass.value = ''

      try {
        const res = await api.getQRCode()
        qrImage.value = res.data.image
        qrCodeKey.value = res.data.qrcode_key

        // 开始轮询
        if (pollTimer) clearInterval(pollTimer)
        pollTimer = setInterval(pollLoginStatus, 3000)
      } catch (e) {
        ElMessage.error('获取二维码失败')
      } finally {
        qrLoading.value = false
      }
    }

    const pollLoginStatus = async () => {
      try {
        const res = await api.loginByQRCode(qrCodeKey.value)
        const data = res.data

        switch (data.code) {
          case 0:
            // 登录成功
            loginStatusText.value = '登录成功!'
            loginStatusClass.value = 'success'
            clearInterval(pollTimer)

            // 保存 cookie
            await api.addUser({
              username: 'bilibili_user',
              cookie_data: data.cookie,
            })

            ElMessage.success('登录成功')
            setTimeout(() => {
              showLoginDialog.value = false
              fetchUsers()
            }, 1000)
            break
          case 86038:
            // 二维码过期
            loginStatusText.value = '二维码已过期，请刷新'
            loginStatusClass.value = 'expired'
            clearInterval(pollTimer)
            break
          case 86090:
            // 已扫码未确认
            loginStatusText.value = '已扫码，请在手机上确认'
            loginStatusClass.value = 'scanned'
            break
          case 86101:
            // 未扫码
            loginStatusText.value = '等待扫码...'
            break
        }
      } catch {}
    }

    const deleteUser = async (row) => {
      try {
        await ElMessageBox.confirm(`确认删除账号 ${row.username}？`, '确认', { type: 'warning' })
        await api.deleteUser(row.id)
        ElMessage.success('删除成功')
        fetchUsers()
      } catch {}
    }

    onMounted(fetchUsers)
    onUnmounted(() => {
      if (pollTimer) clearInterval(pollTimer)
    })

    return {
      users, showLoginDialog, qrImage, qrLoading,
      loginStatusText, loginStatusClass,
      startLogin, deleteUser,
    }
  },
}
</script>

<style scoped>
.page { max-width: 900px; }
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.page-header h2 { font-size: 20px; color: #303133; }
.qr-container { text-align: center; min-height: 300px; }
.qr-image { width: 256px; height: 256px; border: 1px solid #eee; border-radius: 8px; }
.qr-hint { margin-top: 12px; color: #606266; font-size: 14px; }
.qr-status { margin-top: 8px; font-size: 14px; font-weight: 500; }
.qr-status.success { color: #67c23a; }
.qr-status.scanned { color: #e6a23c; }
.qr-status.expired { color: #f56c6c; }
</style>
