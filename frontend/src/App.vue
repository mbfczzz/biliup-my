<template>
  <div id="app">
    <aside class="sidebar">
      <div class="logo">
        <span class="logo-icon">🎬</span>
        <span class="logo-text">Scarecrow</span>
      </div>
      <nav>
        <router-link to="/" class="nav-item">
          <span class="icon">📺</span>
          <span>录播管理</span>
        </router-link>
        <router-link to="/templates" class="nav-item">
          <span class="icon">📤</span>
          <span>上传模板</span>
        </router-link>
        <router-link to="/videos" class="nav-item">
          <span class="icon">📹</span>
          <span>视频列表</span>
        </router-link>
        <router-link to="/files" class="nav-item">
          <span class="icon">📁</span>
          <span>文件管理</span>
        </router-link>
        <router-link to="/users" class="nav-item">
          <span class="icon">👥</span>
          <span>账号管理</span>
        </router-link>
        <router-link to="/config" class="nav-item">
          <span class="icon">⚙️</span>
          <span>系统配置</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <div class="version">v1.0.0</div>
      </div>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
  background: #0a0a0a;
  color: #e0e0e0;
  overflow: hidden;
}
#app {
  display: flex;
  height: 100vh;
  background: linear-gradient(135deg, #0a0a0a 0%, #1a1a1a 100%);
}
.sidebar {
  width: 260px;
  background: linear-gradient(180deg, #1a1a1a 0%, #151515 100%);
  border-right: 1px solid rgba(99, 102, 241, 0.1);
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.5);
}
.logo {
  padding: 28px 24px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, transparent 100%);
}
.logo-icon {
  font-size: 28px;
  filter: drop-shadow(0 2px 8px rgba(99, 102, 241, 0.5));
}
.logo-text {
  font-size: 22px;
  font-weight: 700;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
nav {
  padding: 20px 12px;
  flex: 1;
  overflow-y: auto;
}
nav::-webkit-scrollbar { width: 4px; }
nav::-webkit-scrollbar-thumb { background: #2a2a2a; border-radius: 2px; }
.nav-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  margin-bottom: 6px;
  color: #9ca3af;
  text-decoration: none;
  border-radius: 12px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}
.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: linear-gradient(180deg, #6366f1 0%, #8b5cf6 100%);
  transform: scaleY(0);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.nav-item:hover {
  background: rgba(99, 102, 241, 0.08);
  color: #e0e0e0;
  transform: translateX(4px);
}
.nav-item.router-link-active {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.15) 0%, rgba(99, 102, 241, 0.05) 100%);
  color: #6366f1;
  font-weight: 600;
}
.nav-item.router-link-active::before {
  transform: scaleY(1);
}
.icon {
  font-size: 20px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}
.sidebar-footer {
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}
.version {
  font-size: 12px;
  color: #6b7280;
  text-align: center;
}
.main-content {
  flex: 1;
  overflow-y: auto;
  background: #0a0a0a;
}
.main-content::-webkit-scrollbar { width: 8px; }
.main-content::-webkit-scrollbar-track { background: #0a0a0a; }
.main-content::-webkit-scrollbar-thumb {
  background: #2a2a2a;
  border-radius: 4px;
}
.main-content::-webkit-scrollbar-thumb:hover { background: #3a3a3a; }
</style>
