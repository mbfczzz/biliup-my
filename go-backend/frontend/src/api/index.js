import axios from 'axios'

const http = axios.create({
  baseURL: '/v1',
  timeout: 30000,
})

export default {
  // 主播管理
  getStreamers: () => http.get('/streamers'),
  addStreamer: (data) => http.post('/streamers', data),
  updateStreamer: (data) => http.put('/streamers', data),
  deleteStreamer: (id) => http.delete(`/streamers/${id}`),
  pauseStreamer: (id) => http.put(`/streamers/${id}/pause`),

  // 全局配置
  getConfiguration: () => http.get('/configuration'),
  updateConfiguration: (data) => http.put('/configuration', data),

  // 上传模板
  getUploadTemplates: () => http.get('/upload/streamers'),
  addUploadTemplate: (data) => http.post('/upload/streamers', data),
  getUploadTemplate: (id) => http.get(`/upload/streamers/${id}`),
  deleteUploadTemplate: (id) => http.delete(`/upload/streamers/${id}`),

  // 用户管理
  getUsers: () => http.get('/users'),
  addUser: (data) => http.post('/users', data),
  deleteUser: (id) => http.delete(`/users/${id}`),

  // B站登录
  getQRCode: () => http.get('/get_qrcode'),
  loginByQRCode: (qrcodeKey) => http.post('/login_by_qrcode', { qrcode_key: qrcodeKey }),

  // 系统
  getStatus: () => http.get('/status'),
  getVideos: () => http.get('/videos'),

  // B站信息
  getMyInfo: () => axios.get('/bili/space/myinfo'),
  getArchivePre: () => axios.get('/bili/archive/pre'),
}
