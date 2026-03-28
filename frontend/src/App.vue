<template>
  <div id="app">
    <!-- Mobile header -->
    <header class="mobile-header">
      <button class="menu-btn" @click="sidebarOpen = !sidebarOpen" aria-label="菜单">
        <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 10a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zM3 15a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd"/></svg>
      </button>
      <div class="mobile-logo">
        <div class="logo-mark">S</div>
        <span class="logo-text">Scarecrow</span>
      </div>
      <div style="width:36px"></div>
    </header>

    <!-- Sidebar overlay for mobile -->
    <div v-if="sidebarOpen" class="sidebar-overlay" @click="sidebarOpen = false"></div>

    <aside :class="['sidebar', { open: sidebarOpen }]">
      <div class="logo">
        <div class="logo-mark">S</div>
        <span class="logo-text">Scarecrow</span>
      </div>
      <nav>
        <router-link to="/" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg>
          <span>录播管理</span>
        </router-link>
        <router-link to="/templates" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M12 4v12M8 8l4-4 4 4"/></svg>
          <span>上传模板</span>
        </router-link>
        <router-link to="/videos" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg>
          <span>视频列表</span>
        </router-link>
        <router-link to="/files" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
          <span>文件管理</span>
        </router-link>
        <router-link to="/users" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14c-4.418 0-8 1.79-8 4v2h16v-2c0-2.21-3.582-4-8-4z"/></svg>
          <span>账号管理</span>
        </router-link>
        <router-link to="/config" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
          <span>系统配置</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <div class="version">v1.1.28</div>
      </div>
    </aside>
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue'
const sidebarOpen = ref(false)
</script>

<style>
/* ===== Design Tokens ===== */
:root {
  --c-bg: #09090b;
  --c-bg-elevated: #18181b;
  --c-bg-card: rgba(24, 24, 27, 0.7);
  --c-bg-input: rgba(9, 9, 11, 0.8);
  --c-border: rgba(255,255,255,0.06);
  --c-border-hover: rgba(255,255,255,0.12);
  --c-border-focus: #6366f1;
  --c-accent: #6366f1;
  --c-accent-hover: #818cf8;
  --c-accent-glow: rgba(99,102,241,0.25);
  --c-danger: #ef4444;
  --c-success: #22c55e;
  --c-text-1: #fafafa;
  --c-text-2: #a1a1aa;
  --c-text-3: #71717a;
  --radius: 12px;
  --radius-lg: 16px;
  --ease: cubic-bezier(.4,0,.2,1);
  --font: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  --font-mono: 'SF Mono', 'Fira Code', 'Cascadia Code', monospace;
}

/* ===== Reset ===== */
*,*::before,*::after { margin:0; padding:0; box-sizing:border-box; }
body {
  font-family: var(--font);
  background: var(--c-bg);
  color: var(--c-text-2);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  overflow: hidden;
}

/* ===== Layout ===== */
#app { display:flex; height:100vh; padding-top:0; }
@media (max-width:768px) { #app { padding-top:56px; } }

/* ===== Page transition ===== */
.page-enter-active { animation: fadeSlideIn .25s var(--ease); }
.page-leave-active { animation: fadeSlideOut .15s var(--ease); }
@keyframes fadeSlideIn { from { opacity:0; transform:translateY(8px); } to { opacity:1; transform:translateY(0); } }
@keyframes fadeSlideOut { from { opacity:1; transform:translateY(0); } to { opacity:0; transform:translateY(-4px); } }

/* ===== Shared Components ===== */

/* -- Page -- */
.page { padding:40px 48px; max-width:1200px; margin:0 auto; min-height:100vh; width:100%; }
.page-header { margin-bottom:32px; }
.page-header-row { display:flex; justify-content:space-between; align-items:center; margin-bottom:32px; gap:16px; flex-wrap:wrap; }
.page-title { font-size:24px; font-weight:700; color:var(--c-text-1); letter-spacing:-.3px; line-height:1.3; }
.page-subtitle { font-size:13px; color:var(--c-text-3); margin-top:4px; letter-spacing:.1px; }
@media (max-width:768px) {
  .page { padding:24px 16px; min-height:calc(100vh - 56px); }
  .page-header { margin-bottom:24px; }
  .page-header-row { margin-bottom:24px; }
  .page-title { font-size:20px; }
}
@media (min-width:769px) and (max-width:1024px) {
  .page { padding:32px 28px; }
}

/* -- Card -- */
.card { background:var(--c-bg-card); backdrop-filter:blur(16px); border:1px solid var(--c-border); border-radius:var(--radius); overflow:hidden; transition:border-color .2s var(--ease), box-shadow .2s var(--ease); }
.card:hover { border-color:var(--c-border-hover); box-shadow:0 2px 12px rgba(0,0,0,.15); }

/* -- Table -- */
.table-wrap { overflow-x:auto; -webkit-overflow-scrolling:touch; }
table { width:100%; border-collapse:collapse; min-width:500px; }
th { background:rgba(255,255,255,.02); padding:12px 24px; text-align:left; font-weight:600; color:var(--c-text-3); font-size:11px; text-transform:uppercase; letter-spacing:.8px; white-space:nowrap; }
td { padding:16px 24px; border-top:1px solid var(--c-border); color:var(--c-text-2); font-size:14px; }
tbody tr { transition:background .15s var(--ease); }
tbody tr:hover { background:rgba(99,102,241,.04); }
@media (max-width:768px) {
  th { padding:10px 14px; font-size:10px; }
  td { padding:12px 14px; font-size:13px; }
}

/* -- Buttons -- */
.btn { display:inline-flex; align-items:center; justify-content:center; gap:6px; font-family:var(--font); font-weight:600; font-size:13px; border:none; border-radius:var(--radius); cursor:pointer; transition:all .2s var(--ease); white-space:nowrap; }
.btn:active { transform:scale(.97); }
.btn:disabled { opacity:.4; pointer-events:none; }

.btn-primary { background:var(--c-accent); color:#fff; padding:10px 20px; box-shadow:0 1px 3px rgba(0,0,0,.2), 0 0 0 1px rgba(99,102,241,.3); }
.btn-primary:hover { background:var(--c-accent-hover); box-shadow:0 4px 16px var(--c-accent-glow), 0 0 0 1px rgba(129,140,248,.4); }

.btn-ghost { background:transparent; color:var(--c-text-2); padding:10px 20px; border:1px solid var(--c-border); }
.btn-ghost:hover { background:rgba(255,255,255,.04); border-color:var(--c-border-hover); color:var(--c-text-1); }

.btn-sm { padding:6px 12px; font-size:12px; border-radius:8px; }
.btn-danger { background:rgba(239,68,68,.12); color:#f87171; border:1px solid rgba(239,68,68,.15); }
.btn-danger:hover { background:rgba(239,68,68,.2); border-color:rgba(239,68,68,.3); }

/* -- Badge -- */
.badge { display:inline-flex; align-items:center; gap:6px; padding:4px 10px; border-radius:20px; font-size:11px; font-weight:600; letter-spacing:.3px; }
.badge-live { background:rgba(34,197,94,.12); color:#4ade80; }
.badge-live::before { content:''; width:6px; height:6px; border-radius:50%; background:#4ade80; animation:pulse-dot 1.5s infinite; }
@keyframes pulse-dot { 0%,100%{opacity:1} 50%{opacity:.4} }
.badge-offline { background:rgba(113,113,122,.12); color:var(--c-text-3); }
.badge-success { background:rgba(34,197,94,.12); color:#4ade80; }

/* -- Form -- */
.form-group { margin-bottom:20px; }
.form-group:last-child { margin-bottom:0; }
.form-label { display:block; margin-bottom:8px; color:var(--c-text-2); font-size:13px; font-weight:600; }
.form-input, .form-select {
  width:100%; padding:10px 14px; background:var(--c-bg-input); border:1.5px solid var(--c-border);
  border-radius:10px; color:var(--c-text-1); font-size:14px; font-family:var(--font);
  transition:border-color .15s var(--ease), box-shadow .15s var(--ease);
}
.form-input::placeholder { color:var(--c-text-3); }
.form-input:hover, .form-select:hover { border-color:var(--c-border-hover); }
.form-input:focus, .form-select:focus { outline:none; border-color:var(--c-accent); box-shadow:0 0 0 3px var(--c-accent-glow); }
.form-hint { margin-top:6px; font-size:12px; color:var(--c-text-3); }

/* -- Modal -- */
.modal-overlay { position:fixed; inset:0; background:rgba(0,0,0,.7); backdrop-filter:blur(4px); display:flex; align-items:center; justify-content:center; z-index:1000; animation:fadeIn .15s var(--ease); padding:16px; }
@keyframes fadeIn { from{opacity:0} to{opacity:1} }
.modal-box { background:var(--c-bg-elevated); border:1px solid var(--c-border); border-radius:var(--radius-lg); width:480px; max-width:100%; box-shadow:0 24px 48px rgba(0,0,0,.5); animation:modalIn .2s var(--ease); }
@media (max-width:768px) {
  .modal-overlay { align-items:flex-end; padding:0; }
  .modal-box { width:100%; max-width:100%; border-radius:var(--radius-lg) var(--radius-lg) 0 0; animation:modalSlideUp .25s var(--ease); }
  @keyframes modalSlideUp { from{opacity:0;transform:translateY(100%)} to{opacity:1;transform:translateY(0)} }
}
@keyframes modalIn { from{opacity:0;transform:scale(.96) translateY(8px)} to{opacity:1;transform:scale(1) translateY(0)} }
.modal-head { display:flex; justify-content:space-between; align-items:center; padding:20px 24px; border-bottom:1px solid var(--c-border); }
.modal-title { font-size:17px; font-weight:700; color:var(--c-text-1); }
.modal-close { width:32px; height:32px; display:flex; align-items:center; justify-content:center; background:none; border:none; color:var(--c-text-3); font-size:18px; cursor:pointer; border-radius:8px; transition:all .15s; }
.modal-close:hover { background:rgba(255,255,255,.06); color:var(--c-text-1); }
.modal-body { padding:24px; }
.modal-foot { display:flex; gap:10px; padding:16px 24px; border-top:1px solid var(--c-border); }
.modal-foot .btn { flex:1; }

/* -- Loading / Empty -- */
.state-box { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:100px 20px; }
.spinner { width:36px; height:36px; border:3px solid var(--c-border); border-top-color:var(--c-accent); border-radius:50%; animation:spin .7s linear infinite; margin-bottom:16px; }
@keyframes spin { to{transform:rotate(360deg)} }
.state-text { color:var(--c-text-3); font-size:14px; }
.empty-box { text-align:center; padding:72px 40px; border:1px solid var(--c-border); border-radius:var(--radius-lg); background:var(--c-bg-card); }
.empty-icon { width:56px; height:56px; margin:0 auto 20px; border-radius:14px; background:rgba(99,102,241,.08); color:var(--c-accent); display:flex; align-items:center; justify-content:center; }
.empty-icon svg { width:26px; height:26px; }
.empty-title { font-size:16px; font-weight:700; color:var(--c-text-1); margin-bottom:6px; }
.empty-desc { color:var(--c-text-3); font-size:13px; margin-bottom:24px; line-height:1.6; }
@media (max-width:768px) {
  .empty-box { padding:48px 20px; }
  .empty-title { font-size:15px; }
}

/* -- Scrollbar -- */
::-webkit-scrollbar { width:6px; }
::-webkit-scrollbar-track { background:transparent; }
::-webkit-scrollbar-thumb { background:rgba(255,255,255,.08); border-radius:3px; }
::-webkit-scrollbar-thumb:hover { background:rgba(255,255,255,.14); }
</style>

<style scoped>
/* Mobile header */
.mobile-header {
  display:none;
  position:fixed; top:0; left:0; right:0; height:56px;
  background:var(--c-bg-elevated);
  border-bottom:1px solid var(--c-border);
  align-items:center; justify-content:space-between;
  padding:0 16px; z-index:20;
}
.menu-btn {
  width:36px; height:36px; display:flex; align-items:center; justify-content:center;
  background:none; border:none; color:var(--c-text-2); cursor:pointer;
  border-radius:8px; transition:background .15s;
}
.menu-btn:hover { background:rgba(255,255,255,.06); }
.mobile-logo { display:flex; align-items:center; gap:8px; }

/* Sidebar overlay */
.sidebar-overlay {
  display:none; position:fixed; inset:0; background:rgba(0,0,0,.5);
  backdrop-filter:blur(2px); z-index:25;
}

.sidebar {
  width:240px; min-width:240px;
  background:var(--c-bg-elevated);
  border-right:1px solid var(--c-border);
  display:flex; flex-direction:column;
  z-index:30;
}

.logo {
  padding:24px 20px;
  display:flex; align-items:center; gap:12px;
  border-bottom:1px solid var(--c-border);
}
.logo-mark {
  width:32px; height:32px;
  background:linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a78bfa 100%);
  border-radius:9px;
  display:flex; align-items:center; justify-content:center;
  font-size:15px; font-weight:800; color:#fff;
  box-shadow:0 2px 10px rgba(99,102,241,.3);
}
.logo-text {
  font-size:18px; font-weight:800; color:var(--c-text-1);
  letter-spacing:-.3px;
}

nav { padding:12px 8px; flex:1; overflow-y:auto; }

.nav-item {
  display:flex; align-items:center; gap:10px;
  padding:10px 12px; margin-bottom:2px;
  color:var(--c-text-3); text-decoration:none;
  border-radius:10px; font-size:14px; font-weight:500;
  transition:all .15s var(--ease);
  position:relative;
}
.nav-icon { width:20px; height:20px; flex-shrink:0; opacity:.8; }
.nav-item:hover { color:var(--c-text-2); background:rgba(255,255,255,.04); }
.nav-item.active {
  color:var(--c-text-1); background:rgba(99,102,241,.1);
  font-weight:600;
}
.nav-item.active .nav-icon { color:var(--c-accent); opacity:1; }

.sidebar-footer {
  padding:16px 20px;
  border-top:1px solid var(--c-border);
}
.version { font-size:11px; color:var(--c-text-3); text-align:center; letter-spacing:.5px; }

.main-content {
  flex:1; overflow-y:auto; overflow-x:hidden;
  background:var(--c-bg);
  display:flex; flex-direction:column;
}

/* ===== Responsive ===== */
@media (max-width:768px) {
  .mobile-header { display:flex; }
  .sidebar-overlay { display:block; }
  .sidebar {
    position:fixed; top:0; left:0; bottom:0;
    transform:translateX(-100%);
    transition:transform .25s var(--ease);
    box-shadow:4px 0 24px rgba(0,0,0,.3);
  }
  .sidebar.open { transform:translateX(0); }
  .nav-item { padding:12px 14px; font-size:15px; }
}
@media (min-width:769px) {
  .mobile-header { display:none !important; }
  .sidebar-overlay { display:none !important; }
}
</style>
