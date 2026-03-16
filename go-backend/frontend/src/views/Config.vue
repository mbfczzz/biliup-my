<template>
  <div class="page">
    <h2 style="margin-bottom: 20px;">系统配置</h2>

    <el-card>
      <el-form :model="config" label-width="140px" v-loading="loading">
        <el-divider content-position="left">录制设置</el-divider>

        <el-form-item label="下载目录">
          <el-input v-model="config.download_dir" placeholder="./downloads" />
        </el-form-item>

        <el-form-item label="检测间隔(秒)">
          <el-input-number v-model.number="config.check_interval" :min="10" :max="600" :step="10" />
          <div class="form-hint">检查主播是否在线的时间间隔</div>
        </el-form-item>

        <el-form-item label="分段时间(秒)">
          <el-input-number v-model.number="config.segment_time" :min="0" :step="1800" />
          <div class="form-hint">默认分段时间，0=不分段</div>
        </el-form-item>

        <el-form-item label="文件名模板">
          <el-input v-model="config.filename_prefix" placeholder="{streamer}_%Y-%m-%d_%H_%M_%S" />
        </el-form-item>

        <el-form-item label="下载器">
          <el-select v-model="config.downloader" style="width: 200px">
            <el-option label="FFmpeg" value="ffmpeg" />
            <el-option label="Streamlink" value="streamlink" />
          </el-select>
        </el-form-item>

        <el-divider content-position="left">上传设置</el-divider>

        <el-form-item label="上传线路">
          <el-select v-model="config.upload_line" style="width: 200px">
            <el-option label="自动选择" value="AUTO" />
            <el-option label="bda2" value="bda2" />
            <el-option label="alia" value="alia" />
            <el-option label="cos" value="cos" />
            <el-option label="kodo" value="kodo" />
          </el-select>
        </el-form-item>

        <el-form-item label="上传线程数">
          <el-input-number v-model.number="config.upload_threads" :min="1" :max="16" />
        </el-form-item>

        <el-form-item label="过滤阈值(MB)">
          <el-input-number v-model.number="config.filtering_threshold" :min="0" />
          <div class="form-hint">小于此大小的文件不上传</div>
        </el-form-item>

        <el-form-item label="上传延迟(秒)">
          <el-input-number v-model.number="config.delay" :min="0" />
          <div class="form-hint">录制完成后延迟多久开始上传</div>
        </el-form-item>

        <el-divider content-position="left">Twitch 设置</el-divider>

        <el-form-item label="Auth Token">
          <el-input v-model="config.twitch_cookie" placeholder="可选：Twitch OAuth Token" show-password />
          <div class="form-hint">用于绕过地区限制和广告，留空则匿名</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="saveConfig" :loading="saving">保存配置</el-button>
          <el-button @click="fetchConfig">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '../api'

export default {
  name: 'Config',
  setup() {
    const config = reactive({
      download_dir: './downloads',
      check_interval: '60',
      segment_time: '3600',
      file_size: '0',
      filename_prefix: '{streamer}_%Y-%m-%d_%H_%M_%S',
      downloader: 'ffmpeg',
      upload_line: 'AUTO',
      upload_threads: '3',
      filtering_threshold: '10',
      delay: '0',
      twitch_cookie: '',
    })

    const loading = ref(false)
    const saving = ref(false)

    const fetchConfig = async () => {
      loading.value = true
      try {
        const res = await api.getConfiguration()
        Object.assign(config, res.data)
      } catch {
        ElMessage.error('加载配置失败')
      } finally {
        loading.value = false
      }
    }

    const saveConfig = async () => {
      saving.value = true
      try {
        // 确保数字类型是字符串
        const data = {}
        for (const [k, v] of Object.entries(config)) {
          data[k] = String(v)
        }
        await api.updateConfiguration(data)
        ElMessage.success('配置已保存')
      } catch {
        ElMessage.error('保存失败')
      } finally {
        saving.value = false
      }
    }

    onMounted(fetchConfig)

    return { config, loading, saving, fetchConfig, saveConfig }
  },
}
</script>

<style scoped>
.page { max-width: 800px; }
.form-hint { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
