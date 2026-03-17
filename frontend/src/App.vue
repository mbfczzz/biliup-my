<template>
  <div id="app">
    <aside class="sidebar">
      <div class="logo">Scarecrow</div>
      <nav>
        <router-link to="/" class="nav-item">
          <span class="icon">📺</span>
          <span>录播管理</span>
        </router-link>
        <router-link to="/config" class="nav-item">
          <span class="icon">⚙️</span>
          <span>系统配置</span>
        </router-link>
      </nav>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: #0f0f0f;
  color: #e0e0e0;
}
#app { display: flex; height: 100vh; }
.sidebar {
  width: 240px;
  background: #1a1a1a;
  border-right: 1px solid #2a2a2a;
  display: flex;
  flex-direction: column;
}
.logo {
  padding: 24px 20px;
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  border-bottom: 1px solid #2a2a2a;
}
nav { padding: 16px 0; }
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: #a0a0a0;
  text-decoration: none;
  transition: all 0.2s;
}
.nav-item:hover { background: #252525; color: #fff; }
.nav-item.router-link-active {
  background: #2a2a2a;
  color: #6366f1;
  border-left: 3px solid #6366f1;
}
.icon { font-size: 18px; }
.main-content {
  flex: 1;
  overflow-y: auto;
  background: #0f0f0f;
}
</style>
