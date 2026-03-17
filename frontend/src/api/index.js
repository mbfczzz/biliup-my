import axios from 'axios'

const api = axios.create({
  baseURL: '/v1',
  timeout: 10000
})

export default {
  // 录播管理
  getStreamers: () => api.get('/streamers'),
  addStreamer: (data) => api.post('/streamers', data),
  updateStreamer: (data) => api.put('/streamers', data),
  deleteStreamer: (id) => api.delete(`/streamers/${id}`),
  pauseStreamer: (id) => api.put(`/streamers/${id}/pause`),

  // 配置
  getConfig: () => api.get('/configuration'),
  updateConfig: (data) => api.put('/configuration', data),

  // 上传模板
  getTemplates: () => api.get('/upload/streamers'),
  addTemplate: (data) => api.post('/upload/streamers', data),
  getTemplate: (id) => api.get(`/upload/streamers/${id}`),
  deleteTemplate: (id) => api.delete(`/upload/streamers/${id}`),

  // 用户管理
  getUsers: () => api.get('/users'),
  addUser: (data) => api.post('/users', data),
  deleteUser: (id) => api.delete(`/users/${id}`),

  // 登录
  getQrcode: () => api.get('/get_qrcode'),
  loginByQrcode: (data) => api.post('/login_by_qrcode', data),

  // 视频
  getVideos: () => api.get('/videos'),

  // 文件
  getFiles: (id) => api.get(`/streamer-info/files/${id}`),

  // 上传
  postUpload: (data) => api.post('/uploads', data),
  getStatus: () => api.get('/status')
}
