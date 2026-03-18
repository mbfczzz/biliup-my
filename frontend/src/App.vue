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
:root {
  --color-primary: #6366f1;
  --color-primary-dark: #4f46e5;
  --color-secondary: #8b5cf6;
  --color-success: #34d399;
  --color-danger: #dc2626;
  --color-bg-primary: #0a0a0a;
  --color-bg-secondary: #1a1a1a;
  --color-text-primary: #fff;
  --color-text-secondary: #e5e7eb;
  --color-text-tertiary: #9ca3af;
  --color-border: rgba(255, 255, 255, 0.06);
  --radius-sm: 10px;
  --radius-md: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
  --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.3);
  --shadow-md: 0 4px 24px rgba(0, 0, 0, 0.4);
  --shadow-lg: 0 8px 32px rgba(0, 0, 0, 0.5);
  --transition-fast: 0.15s cubic-bezier(0.4, 0, 0.2, 1);
  --transition-base: 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  --transition-slow: 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', sans-serif;
  background: var(--color-bg-primary);
  color: var(--color-text-secondary);
  overflow: hidden;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-rendering: optimizeLegibility;
}
#app {
  display: flex;
  height: 100vh;
  background: linear-gradient(135deg, var(--color-bg-primary) 0%, var(--color-bg-secondary) 100%);
}

.sidebar {
  width: 260px;
  background: linear-gradient(180deg, var(--color-bg-secondary) 0%, #151515 100%);
  border-right: 1px solid rgba(99, 102, 241, 0.1);
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 24px rgba(0, 0, 0, 0.5);
  position: relative;
  z-index: 10;
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
  user-select: none;
}

.logo-text {
  font-size: 22px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--color-primary) 0%, var(--color-secondary) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.5px;
  user-select: none;
}
nav {
  padding: 20px 12px;
  flex: 1;
  overflow-y: auto;
  scroll-behavior: smooth;
}

nav::-webkit-scrollbar { width: 4px; }
nav::-webkit-scrollbar-thumb {
  background: #2a2a2a;
  border-radius: 2px;
  transition: background var(--transition-base);
}
nav::-webkit-scrollbar-thumb:hover { background: #3a3a3a; }

.nav-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 16px;
  margin-bottom: 6px;
  color: var(--color-text-tertiary);
  text-decoration: none;
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  position: relative;
  overflow: hidden;
  cursor: pointer;
  user-select: none;
}

.nav-item::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  width: 3px;
  height: 100%;
  background: linear-gradient(180deg, var(--color-primary) 0%, var(--color-secondary) 100%);
  transform: scaleY(0);
  transition: transform var(--transition-slow);
}

.nav-item:hover {
  background: rgba(99, 102, 241, 0.08);
  color: var(--color-text-secondary);
  transform: translateX(4px);
}

.nav-item:active {
  transform: translateX(2px);
  transition: transform var(--transition-fast);
}

.nav-item:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

.nav-item.router-link-active {
  background: linear-gradient(90deg, rgba(99, 102, 241, 0.15) 0%, rgba(99, 102, 241, 0.05) 100%);
  color: var(--color-primary);
  font-weight: 600;
}

.nav-item.router-link-active::before {
  transform: scaleY(1);
}

.icon {
  font-size: 20px;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
  transition: transform var(--transition-base);
}

.nav-item:hover .icon {
  transform: scale(1.1);
}
.sidebar-footer {
  padding: 20px 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(0, 0, 0, 0.2);
}

.version {
  font-size: 12px;
  color: #6b7280;
  text-align: center;
  font-weight: 500;
  letter-spacing: 0.5px;
}

.main-content {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  background: var(--color-bg-primary);
  scroll-behavior: smooth;
}

.main-content::-webkit-scrollbar {
  width: 8px;
}

.main-content::-webkit-scrollbar-track {
  background: var(--color-bg-primary);
}

.main-content::-webkit-scrollbar-thumb {
  background: #2a2a2a;
  border-radius: 4px;
  transition: background var(--transition-base);
}

.main-content::-webkit-scrollbar-thumb:hover {
  background: #3a3a3a;
}

.main-content::-webkit-scrollbar-thumb:active {
  background: #4a4a4a;
}
</style>
