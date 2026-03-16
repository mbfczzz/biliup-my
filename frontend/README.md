# Twitch 录播工具

## 项目结构

```
biliup/
├── frontend/          # Vue 3 前端
├── crates/           # Rust 后端
└── scarecrow/        # Python CLI 工具
```

## 前端 (Vue 3)

### 安装依赖
```bash
cd frontend
npm install
```

### 开发
```bash
npm run dev
```
访问 http://localhost:3000

### 构建
```bash
npm run build
```

## 后端 (Rust)

后端 API 运行在 `http://localhost:19159`

### 主要 API 端点
- `GET /v1/streamers` - 获取 Twitch 频道列表
- `POST /v1/streamers` - 添加 Twitch 频道
- `PUT /v1/streamers` - 更新频道
- `DELETE /v1/streamers/{id}` - 删除频道
- `PUT /v1/streamers/{id}/pause` - 暂停/恢复录播
- `GET /v1/configuration` - 获取配置
- `PUT /v1/configuration` - 更新配置

## Python CLI (scarecrow)

```bash
pip install -e .
scarecrow --help
```
