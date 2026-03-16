<template>
  <div class="page">
    <div class="page-header">
      <h2>录制管理</h2>
      <el-button type="primary" @click="showAddDialog = true">添加频道</el-button>
    </div>

    <el-table :data="streamers" border stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column label="频道" min-width="200">
        <template #default="{ row }">
          <div>
            <a :href="row.url" target="_blank" class="channel-link">
              {{ getChannelName(row.url) }}
            </a>
            <div class="remark" v-if="row.remark">{{ row.remark }}</div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row)" effect="dark" size="small">
            {{ getStatusLabel(row) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="format" label="格式" width="80" />
      <el-table-column label="分段" width="100">
        <template #default="{ row }">
          {{ row.segment_time > 0 ? (row.segment_time / 60) + '分钟' : '不分段' }}
        </template>
      </el-table-column>
      <el-table-column label="上传模板" width="120">
        <template #default="{ row }">
          {{ getTemplateName(row.upload_template_id) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="togglePause(row)">
            {{ row.paused ? '恢复' : '暂停' }}
          </el-button>
          <el-button size="small" @click="editStreamer(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteStreamer(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editingStreamer ? '编辑频道' : '添加频道'"
      width="560px"
    >
      <el-form :model="form" label-width="100px">
        <el-form-item label="Twitch URL" required>
          <el-input v-model="form.url" placeholder="https://www.twitch.tv/频道名" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" placeholder="可选备注" />
        </el-form-item>
        <el-form-item label="文件名模板">
          <el-input v-model="form.filename_prefix" placeholder="{streamer}_%Y-%m-%d_%H_%M_%S" />
          <div class="form-hint">可用变量: {streamer} 频道名, %Y年 %m月 %d日 %H时 %M分 %S秒</div>
        </el-form-item>
        <el-form-item label="录制格式">
          <el-select v-model="form.format" style="width: 100%">
            <el-option label="FLV" value="flv" />
            <el-option label="MP4" value="mp4" />
            <el-option label="TS" value="ts" />
          </el-select>
        </el-form-item>
        <el-form-item label="分段时间(秒)">
          <el-input-number v-model="form.segment_time" :min="0" :step="1800" />
          <div class="form-hint">0 = 不分段, 3600 = 1小时</div>
        </el-form-item>
        <el-form-item label="上传模板">
          <el-select v-model="form.upload_template_id" style="width: 100%" clearable placeholder="不自动上传">
            <el-option
              v-for="t in templates"
              :key="t.id"
              :label="t.template_name"
              :value="t.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="录制时段">
          <el-input v-model="form.time_range" placeholder="如: 00:00-06:00 (留空=全天)" />
        </el-form-item>
        <el-form-item label="排除关键词">
          <el-input v-model="form.excluded_keywords" placeholder="标题包含这些词时不录制，逗号分隔" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveStreamer" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

export default {
  name: 'Streamers',
  setup() {
    const streamers = ref([])
    const templates = ref([])
    const showAddDialog = ref(false)
    const editingStreamer = ref(null)
    const saving = ref(false)

    const defaultForm = {
      url: '',
      remark: '',
      filename_prefix: '{streamer}_%Y-%m-%d_%H_%M_%S',
      format: 'flv',
      segment_time: 3600,
      file_size: 0,
      upload_template_id: 0,
      time_range: '',
      excluded_keywords: '',
      opt_args: '',
    }

    const form = reactive({ ...defaultForm })

    const fetchData = async () => {
      try {
        const [s, t] = await Promise.all([
          api.getStreamers(),
          api.getUploadTemplates(),
        ])
        streamers.value = s.data || []
        templates.value = t.data || []
      } catch (e) {
        ElMessage.error('加载数据失败')
      }
    }

    const getChannelName = (url) => {
      return url.replace(/\/$/, '').split('/').pop()
    }

    const getStatusType = (row) => {
      if (row.paused) return 'warning'
      switch (row.status) {
        case 'recording': return 'success'
        case 'uploading': return 'primary'
        case 'error': return 'danger'
        case 'checking': return 'info'
        default: return 'info'
      }
    }

    const getStatusLabel = (row) => {
      if (row.paused) return '已暂停'
      switch (row.status) {
        case 'recording': return '录制中'
        case 'uploading': return '上传中'
        case 'error': return '错误'
        case 'checking': return '检测中'
        default: return '待机'
      }
    }

    const getTemplateName = (id) => {
      if (!id) return '无'
      const t = templates.value.find((t) => t.id === id)
      return t ? t.template_name : '未知'
    }

    const editStreamer = (row) => {
      editingStreamer.value = row
      Object.assign(form, row)
      showAddDialog.value = true
    }

    const saveStreamer = async () => {
      if (!form.url) {
        ElMessage.warning('请输入 Twitch URL')
        return
      }
      saving.value = true
      try {
        if (editingStreamer.value) {
          await api.updateStreamer({ ...form, id: editingStreamer.value.id })
          ElMessage.success('更新成功')
        } else {
          await api.addStreamer(form)
          ElMessage.success('添加成功')
        }
        showAddDialog.value = false
        editingStreamer.value = null
        Object.assign(form, defaultForm)
        fetchData()
      } catch (e) {
        ElMessage.error(e.response?.data?.error || '操作失败')
      } finally {
        saving.value = false
      }
    }

    const deleteStreamer = async (row) => {
      try {
        await ElMessageBox.confirm(`确认删除 ${getChannelName(row.url)}？`, '确认', {
          type: 'warning',
        })
        await api.deleteStreamer(row.id)
        ElMessage.success('删除成功')
        fetchData()
      } catch {}
    }

    const togglePause = async (row) => {
      try {
        await api.pauseStreamer(row.id)
        ElMessage.success(row.paused ? '已恢复' : '已暂停')
        fetchData()
      } catch {
        ElMessage.error('操作失败')
      }
    }

    onMounted(() => {
      fetchData()
      // 定时刷新状态
      setInterval(fetchData, 15000)
    })

    return {
      streamers, templates, showAddDialog, editingStreamer, form, saving,
      getChannelName, getStatusType, getStatusLabel, getTemplateName,
      editStreamer, saveStreamer, deleteStreamer, togglePause,
    }
  },
}
</script>

<style scoped>
.page { max-width: 1200px; }
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.page-header h2 { font-size: 20px; color: #303133; }
.channel-link {
  color: #6441a5;
  font-weight: 500;
  text-decoration: none;
}
.channel-link:hover { text-decoration: underline; }
.remark { font-size: 12px; color: #909399; margin-top: 2px; }
.form-hint { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
