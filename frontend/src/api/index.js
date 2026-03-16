import axios from 'axios'

const api = axios.create({
  baseURL: '/v1',
  timeout: 10000
})

export default {
  getStreamers: () => api.get('/streamers'),
  addStreamer: (data) => api.post('/streamers', data),
  updateStreamer: (data) => api.put('/streamers', data),
  deleteStreamer: (id) => api.delete(`/streamers/${id}`),
  pauseStreamer: (id) => api.put(`/streamers/${id}/pause`),
  getConfig: () => api.get('/configuration'),
  updateConfig: (data) => api.put('/configuration', data)
}
