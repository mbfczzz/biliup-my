<template>
  <div class="page">
    <div class="page-header">
      <h2>上传模板</h2>
      <el-button type="primary" @click="showAddDialog = true">新建模板</el-button>
    </div>

    <el-table :data="templates" border stripe style="width: 100%">
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="template_name" label="模板名称" width="160" />
      <el-table-column prop="title" label="视频标题模板" min-width="200" />
      <el-table-column prop="tid" label="分区ID" width="80" />
      <el-table-column label="版权" width="80">
        <template #default="{ row }">
          {{ row.copyright === 1 ? '原创' : '转载' }}
        </template>
      </el-table-column>
      <el-table-column prop="tags" label="标签" min-width="150" />
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="editTemplate(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteTemplate(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="showAddDialog"
      :title="editing ? '编辑模板' : '新建模板'"
      width="600px"
    >
      <el-form :model="form" label-width="100px">
        <el-form-item label="模板名称" required>
          <el-input v-model="form.template_name" placeholder="给模板起个名字" />
        </el-form-item>
        <el-form-item label="视频标题">
          <el-input v-model="form.title" placeholder="{streamer} {title} %Y-%m-%d" />
          <div class="form-hint">可用变量: {streamer} 主播名, {title} 直播标题</div>
        </el-form-item>
        <el-form-item label="分区ID">
          <el-input-number v-model="form.tid" :min="1" />
          <div class="form-hint">21=日常, 171=电子竞技, 236=其他游戏</div>
        </el-form-item>
        <el-form-item label="版权类型">
          <el-radio-group v-model="form.copyright">
            <el-radio :label="1">原创</el-radio>
            <el-radio :label="2">转载</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="转载来源" v-if="form.copyright === 2">
          <el-input v-model="form.copyright_source" placeholder="https://www.twitch.tv/..." />
        </el-form-item>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="标签1,标签2,标签3" />
          <div class="form-hint">用英文逗号分隔，至少1个，最多10个</div>
        </el-form-item>
        <el-form-item label="简介">
          <el-input v-model="form.description" type="textarea" :rows="3" placeholder="视频简介" />
        </el-form-item>
        <el-form-item label="封面路径">
          <el-input v-model="form.cover_path" placeholder="留空自动截取" />
        </el-form-item>
        <el-form-item label="定时发布">
          <el-input-number v-model="form.dtime" :min="0" />
          <div class="form-hint">0=立即发布，其他为Unix时间戳</div>
        </el-form-item>
        <el-form-item label="杜比音效">
          <el-switch v-model="form.dolby" :active-value="1" :inactive-value="0" />
        </el-form-item>
        <el-form-item label="禁止转载">
          <el-switch v-model="form.no_reprint" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveTemplate" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '../api'

export default {
  name: 'Templates',
  setup() {
    const templates = ref([])
    const showAddDialog = ref(false)
    const editing = ref(null)
    const saving = ref(false)

    const defaultForm = {
      template_name: '',
      title: '{streamer} {title} %Y-%m-%d',
      tid: 21,
      copyright: 2,
      copyright_source: '',
      cover_path: '',
      description: '',
      dynamic: '',
      dtime: 0,
      dolby: 0,
      hires: 0,
      no_reprint: 0,
      tags: '',
      user_cookie: '',
    }

    const form = reactive({ ...defaultForm })

    const fetchTemplates = async () => {
      try {
        const res = await api.getUploadTemplates()
        templates.value = res.data || []
      } catch {
        ElMessage.error('加载模板失败')
      }
    }

    const editTemplate = (row) => {
      editing.value = row
      Object.assign(form, row)
      showAddDialog.value = true
    }

    const saveTemplate = async () => {
      if (!form.template_name) {
        ElMessage.warning('请输入模板名称')
        return
      }
      saving.value = true
      try {
        await api.addUploadTemplate(form)
        ElMessage.success(editing.value ? '更新成功' : '创建成功')
        showAddDialog.value = false
        editing.value = null
        Object.assign(form, defaultForm)
        fetchTemplates()
      } catch (e) {
        ElMessage.error(e.response?.data?.error || '操作失败')
      } finally {
        saving.value = false
      }
    }

    const deleteTemplate = async (row) => {
      try {
        await ElMessageBox.confirm(`确认删除模板 "${row.template_name}"？`, '确认', { type: 'warning' })
        await api.deleteUploadTemplate(row.id)
        ElMessage.success('删除成功')
        fetchTemplates()
      } catch {}
    }

    onMounted(fetchTemplates)

    return {
      templates, showAddDialog, editing, form, saving,
      editTemplate, saveTemplate, deleteTemplate,
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
.form-hint { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
