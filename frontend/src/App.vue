<template>
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
        <div class="nav-group-label">录制</div>
        <router-link to="/" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg>
          <span>直播录制</span>
        </router-link>
        <router-link to="/files" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/></svg>
          <span>本地文件</span>
        </router-link>
        <div class="nav-group-label">投稿</div>
        <router-link to="/templates" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M4 16v2a2 2 0 002 2h12a2 2 0 002-2v-2"/><path d="M12 4v12M8 8l4-4 4 4"/></svg>
          <span>投稿模板</span>
        </router-link>
        <router-link to="/videos" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14"/><rect x="3" y="6" width="12" height="12" rx="2"/></svg>
          <span>投稿记录</span>
        </router-link>
        <div class="nav-group-label">设置</div>
        <router-link to="/users" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14c-4.418 0-8 1.79-8 4v2h16v-2c0-2.21-3.582-4-8-4z"/></svg>
          <span>B 站账号</span>
        </router-link>
        <router-link to="/config" class="nav-item" exact-active-class="active" @click="sidebarOpen = false">
          <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 01-2.83 2.83l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/></svg>
          <span>系统设置</span>
        </router-link>
      </nav>
      <div class="sidebar-footer">
        <div class="version">Scarecrow v1.1.28</div>
      </div>
    </aside>
    <main class="main-content">
      <router-view v-slot="{ Component }">
        <transition name="page" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
</template>

<script setup>
import { ref } from 'vue'
const sidebarOpen = ref(false)
</script>

<style>
/* ===== Design Tokens ===== */
:root {
  --c-bg: #f6f5f4;
  --c-bg-elevated: #ffffff;
  --c-bg-card: #ffffff;
  --c-bg-input: #f5f3f2;
  --c-bg-subtle: #faf9f8;
  --c-bg-2: #f0efee;
  --c-border: rgba(0,0,0,0.05);
  --c-border-hover: rgba(0,0,0,0.1);
  --c-border-focus: #c48a9c;
  --c-accent: #c48a9c;
  --c-accent-hover: #b47688;
  --c-accent-soft: rgba(196,138,156,0.08);
  --c-accent-glow: rgba(196,138,156,0.15);
  --c-danger: #d14d4d;
  --c-danger-soft: rgba(209,77,77,0.06);
  --c-success: #3d9a5f;
  --c-success-soft: rgba(61,154,95,0.08);
  --c-warning: #d4940a;
  --c-text-1: #1a1a1a;
  --c-text-2: #555555;
  --c-text-3: #8c8c8c;
  --c-text-4: #b0b0b0;
  --radius: 12px;
  --radius-lg: 16px;
  --radius-xl: 20px;
  --shadow-sm: 0 1px 2px rgba(0,0,0,.04), 0 1px 4px rgba(0,0,0,.04);
  --shadow-md: 0 2px 8px rgba(0,0,0,.06), 0 4px 16px rgba(0,0,0,.04);
  --shadow-lg: 0 4px 16px rgba(0,0,0,.08), 0 8px 32px rgba(0,0,0,.06);
  --ease: cubic-bezier(.4,0,.2,1);
  --ease-out: cubic-bezier(0,.5,.2,1);
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
.page-enter-active { animation: fadeSlideIn .3s var(--ease-out); }
.page-leave-active { animation: fadeSlideOut .15s var(--ease); }
@keyframes fadeSlideIn { from { opacity:0; transform:translateY(12px); } to { opacity:1; transform:translateY(0); } }
@keyframes fadeSlideOut { from { opacity:1; transform:translateY(0); } to { opacity:0; transform:translateY(-6px); } }

/* ===== Shared Components ===== */

/* -- Page -- */
.page { min-height:100vh; padding:48px 56px; max-width:1000px; margin:0 auto; }
.page-header { margin-bottom:36px; }
.page-header-row { display:flex; justify-content:space-between; align-items:center; margin-bottom:36px; gap:16px; flex-wrap:wrap; }
.page-title {
  font-size:28px; font-weight:800; color:var(--c-text-1);
  letter-spacing:-.5px; line-height:1.2;
}
.page-subtitle { font-size:14px; color:var(--c-text-3); margin-top:8px; letter-spacing:.1px; }
@media (max-width:768px) {
  .page { padding:24px 20px; min-height:calc(100vh - 56px); }
  .page-header { margin-bottom:28px; }
  .page-header-row { margin-bottom:28px; }
  .page-title { font-size:22px; }
}
@media (min-width:769px) and (max-width:1024px) {
  .page { padding:36px 36px; }
}

/* -- Card -- */
.card {
  background:#fff;
  border:1px solid var(--c-border);
  border-radius:var(--radius-lg);
  overflow:hidden;
  transition:all .25s var(--ease);
  box-shadow:var(--shadow-sm);
}
.card:hover { box-shadow:var(--shadow-md); border-color:var(--c-border-hover); }

/* -- Table -- */
.table-wrap { overflow-x:auto; -webkit-overflow-scrolling:touch; }
table { width:100%; border-collapse:collapse; min-width:500px; }
th {
  background:var(--c-bg-subtle);
  padding:12px 24px; text-align:left;
  font-weight:600; color:var(--c-text-3);
  font-size:11px; text-transform:uppercase; letter-spacing:1px;
  white-space:nowrap; border-bottom:1px solid var(--c-border);
}
td { padding:16px 24px; border-top:1px solid rgba(0,0,0,.03); color:var(--c-text-2); font-size:14px; }
tbody tr { transition:all .15s var(--ease); }
tbody tr:hover { background:var(--c-accent-soft); }
@media (max-width:768px) {
  th { padding:10px 14px; font-size:10px; }
  td { padding:12px 14px; font-size:13px; }
}

/* -- Buttons -- */
.btn {
  display:inline-flex; align-items:center; justify-content:center; gap:6px;
  font-family:var(--font); font-weight:600; font-size:13px;
  border:none; border-radius:10px; cursor:pointer;
  transition:all .2s var(--ease); white-space:nowrap;
  position:relative; overflow:hidden;
}
.btn:active { transform:scale(.97); }
.btn:disabled { opacity:.4; pointer-events:none; }

.btn-primary {
  background:linear-gradient(135deg, var(--c-accent), var(--c-accent-hover));
  color:#fff; padding:10px 22px;
  box-shadow:0 2px 8px rgba(196,138,156,.25), 0 1px 2px rgba(196,138,156,.15);
}
.btn-primary:hover {
  box-shadow:0 4px 16px rgba(196,138,156,.35), 0 2px 4px rgba(196,138,156,.2);
  transform:translateY(-1px);
  filter:brightness(1.05);
}
.btn-primary:active { transform:translateY(0) scale(.97); filter:brightness(.98); }

.btn-ghost {
  background:transparent; color:var(--c-text-2);
  padding:10px 22px; border:1.5px solid rgba(0,0,0,.08);
}
.btn-ghost:hover { background:rgba(0,0,0,.03); border-color:rgba(0,0,0,.12); color:var(--c-text-1); }

.btn-sm { padding:6px 14px; font-size:12px; border-radius:8px; }
.btn-danger {
  background:var(--c-danger-soft); color:#d14d4d;
  border:1px solid rgba(209,77,77,.12);
}
.btn-danger:hover { background:rgba(209,77,77,.12); }

/* -- Badge -- */
.badge {
  display:inline-flex; align-items:center; gap:6px;
  padding:5px 12px; border-radius:20px;
  font-size:11px; font-weight:600; letter-spacing:.3px;
}
.badge-live { background:var(--c-success-soft); color:var(--c-success); }
.badge-live::before {
  content:''; width:7px; height:7px; border-radius:50%;
  background:var(--c-success);
  box-shadow:0 0 0 3px rgba(61,154,95,.2);
  animation:pulse-dot 2s ease-in-out infinite;
}
@keyframes pulse-dot { 0%,100%{opacity:1;box-shadow:0 0 0 3px rgba(61,154,95,.2)} 50%{opacity:.6;box-shadow:0 0 0 6px rgba(61,154,95,.1)} }
.badge-offline { background:rgba(0,0,0,.04); color:var(--c-text-4); }
.badge-success { background:var(--c-success-soft); color:var(--c-success); }

/* -- Form -- */
.form-group { margin-bottom:20px; }
.form-group:last-child { margin-bottom:0; }
.form-label {
  display:flex; align-items:center; gap:4px;
  margin-bottom:8px; color:var(--c-text-2); font-size:13px; font-weight:600; letter-spacing:.1px;
}
.form-input, .form-select {
  width:100%; padding:11px 16px; background:var(--c-bg-input); border:1.5px solid rgba(0,0,0,.05);
  border-radius:12px; color:var(--c-text-1); font-size:14px; font-family:var(--font);
  transition:all .2s var(--ease);
}
.form-input::placeholder { color:var(--c-text-4); }
.form-input:hover, .form-select:hover { border-color:var(--c-border-hover); background:rgba(255,255,255,.8); }
.form-input:focus, .form-select:focus {
  outline:none; border-color:var(--c-accent);
  box-shadow:0 0 0 4px var(--c-accent-glow);
  background:#fff;
}
select.form-input {
  appearance:none;
  background-image:url("data:image/svg+xml,%3Csvg width='10' height='6' viewBox='0 0 10 6' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1l4 4 4-4' stroke='%23b0b0b0' stroke-width='1.5' stroke-linecap='round' stroke-linejoin='round'/%3E%3C/svg%3E");
  background-repeat:no-repeat; background-position:right 14px center; padding-right:36px;
}
.form-hint { margin-top:6px; font-size:12px; color:var(--c-text-4); line-height:1.5; }

/* -- Modal -- */
.modal-overlay {
  position:fixed; inset:0;
  background:rgba(10,10,10,.45); backdrop-filter:blur(12px) saturate(180%);
  display:flex; align-items:center; justify-content:center;
  z-index:1000; animation:overlayIn .25s var(--ease); padding:20px;
}
@keyframes overlayIn { from{opacity:0} to{opacity:1} }
.modal-box {
  background:#fff;
  border:1px solid rgba(0,0,0,.06);
  border-radius:22px; width:480px; max-width:100%;
  box-shadow: 0 8px 40px rgba(0,0,0,.12), 0 2px 12px rgba(0,0,0,.06), 0 0 0 1px rgba(255,255,255,.7) inset;
  animation:modalIn .3s var(--ease-out);
  display:flex; flex-direction:column;
}
.modal-wide { width:640px; }
@media (max-width:768px) {
  .modal-overlay { align-items:flex-end; padding:0; }
  .modal-box { width:100%; max-width:100%; border-radius:22px 22px 0 0; animation:modalSlideUp .35s var(--ease-out); max-height:92vh; }
  .modal-wide { width:100%; }
  @keyframes modalSlideUp { from{opacity:0;transform:translateY(100%)} to{opacity:1;transform:translateY(0)} }
}
@keyframes modalIn { from{opacity:0;transform:scale(.96) translateY(8px)} to{opacity:1;transform:scale(1) translateY(0)} }
.modal-head {
  display:flex; justify-content:space-between; align-items:center;
  padding:28px 32px 0; flex-shrink:0;
}
.modal-title { font-size:19px; font-weight:800; color:var(--c-text-1); letter-spacing:-.3px; }
.modal-close {
  width:36px; height:36px; display:flex; align-items:center; justify-content:center;
  background:rgba(0,0,0,.03); border:none; color:var(--c-text-3); font-size:18px;
  cursor:pointer; border-radius:50%; transition:all .2s var(--ease);
}
.modal-close:hover { background:rgba(0,0,0,.07); color:var(--c-text-1); transform:rotate(90deg); }
.modal-body { padding:28px 32px; flex:1; min-height:0; }
.modal-body-scroll { overflow-y:auto; max-height:calc(80vh - 140px); }
.modal-body-scroll::-webkit-scrollbar { width:5px; }
.modal-body-scroll::-webkit-scrollbar-thumb { background:rgba(0,0,0,.08); border-radius:3px; }
.modal-foot {
  display:flex; gap:10px; padding:20px 32px 28px;
  border-top:1px solid rgba(0,0,0,.04); flex-shrink:0;
}
.modal-foot .btn { flex:1; padding:12px 22px; font-size:14px; }

/* -- Loading / Empty -- */
.state-box { display:flex; flex-direction:column; align-items:center; justify-content:center; padding:120px 20px; }
.spinner {
  width:40px; height:40px;
  border:3px solid rgba(0,0,0,.06);
  border-top-color:var(--c-accent);
  border-radius:50%; animation:spin .8s cubic-bezier(.5,.1,.5,.9) infinite;
  margin-bottom:20px;
}
@keyframes spin { to{transform:rotate(360deg)} }
.state-text { color:var(--c-text-4); font-size:14px; font-weight:500; }
.empty-box {
  text-align:center; padding:80px 48px;
  border:1px solid var(--c-border); border-radius:var(--radius-xl);
  background:#fff; box-shadow:var(--shadow-sm);
}
.empty-icon {
  width:64px; height:64px; margin:0 auto 24px;
  border-radius:16px;
  background:linear-gradient(135deg, rgba(196,138,156,.1), rgba(196,138,156,.04));
  color:var(--c-accent); display:flex; align-items:center; justify-content:center;
}
.empty-icon svg { width:28px; height:28px; }
.empty-title { font-size:17px; font-weight:700; color:var(--c-text-1); margin-bottom:8px; letter-spacing:-.2px; }
.empty-desc { color:var(--c-text-3); font-size:13px; margin-bottom:28px; line-height:1.7; }
@media (max-width:768px) {
  .empty-box { padding:48px 20px; }
  .empty-title { font-size:15px; }
}

/* -- Toast / Notification -- */
.toast {
  position:fixed; bottom:24px; right:24px;
  background:#fff; border:1px solid var(--c-border);
  border-radius:var(--radius); padding:14px 20px;
  box-shadow:var(--shadow-lg); font-size:14px; font-weight:500;
  color:var(--c-text-1); z-index:2000;
  animation:toastIn .3s var(--ease-out);
}
@keyframes toastIn { from{opacity:0;transform:translateY(12px)} to{opacity:1;transform:translateY(0)} }

/* -- Scrollbar -- */
::-webkit-scrollbar { width:6px; }
::-webkit-scrollbar-track { background:transparent; }
::-webkit-scrollbar-thumb { background:rgba(0,0,0,.08); border-radius:3px; }
::-webkit-scrollbar-thumb:hover { background:rgba(0,0,0,.15); }
</style>

<style scoped>
/* Mobile header */
.mobile-header {
  display:none;
  position:fixed; top:0; left:0; right:0; height:56px;
  background:rgba(255,255,255,.85);
  backdrop-filter:blur(12px) saturate(180%);
  border-bottom:1px solid var(--c-border);
  align-items:center; justify-content:space-between;
  padding:0 16px; z-index:20;
}
.menu-btn {
  width:36px; height:36px; display:flex; align-items:center; justify-content:center;
  background:none; border:none; color:var(--c-text-2); cursor:pointer;
  border-radius:8px; transition:background .15s;
}
.menu-btn:hover { background:rgba(0,0,0,.04); }
.mobile-logo { display:flex; align-items:center; gap:8px; }

/* Sidebar overlay */
.sidebar-overlay {
  display:none; position:fixed; inset:0; background:rgba(0,0,0,.5);
  backdrop-filter:blur(4px); z-index:25;
}

.sidebar {
  width:240px; min-width:240px;
  background:rgba(255,255,255,.92);
  backdrop-filter:blur(20px) saturate(180%);
  border-right:1px solid var(--c-border);
  display:flex; flex-direction:column;
  z-index:30;
}

.logo {
  padding:28px 20px 20px;
  display:flex; align-items:center; gap:12px;
}
.logo-mark {
  width:34px; height:34px;
  background:linear-gradient(135deg, #d4a0b0, #b47688);
  border-radius:10px;
  display:flex; align-items:center; justify-content:center;
  font-size:15px; font-weight:800; color:#fff;
  box-shadow:0 2px 8px rgba(196,138,156,.35);
}
.logo-text {
  font-size:18px; font-weight:800; color:var(--c-text-1);
  letter-spacing:-.4px;
}

nav { padding:8px 10px; flex:1; overflow-y:auto; }

.nav-group-label {
  padding:16px 12px 6px;
  font-size:10px; font-weight:700;
  text-transform:uppercase; letter-spacing:1.2px;
  color:var(--c-text-4);
}
.nav-group-label:first-child { padding-top:4px; }

.nav-item {
  display:flex; align-items:center; gap:10px;
  padding:9px 12px; margin-bottom:2px;
  color:var(--c-text-3); text-decoration:none;
  border-radius:10px; font-size:13.5px; font-weight:500;
  transition:all .15s var(--ease);
  position:relative;
}
.nav-icon { width:19px; height:19px; flex-shrink:0; opacity:.6; }
.nav-item:hover { color:var(--c-text-2); background:rgba(0,0,0,.03); }
.nav-item:hover .nav-icon { opacity:.8; }
.nav-item.active {
  color:var(--c-accent); background:var(--c-accent-soft);
  font-weight:600;
}
.nav-item.active::before {
  content:''; position:absolute; left:0; top:50%; transform:translateY(-50%);
  width:3px; height:16px; border-radius:0 3px 3px 0;
  background:var(--c-accent);
}
.nav-item.active .nav-icon { color:var(--c-accent); opacity:1; }

.sidebar-footer {
  padding:16px 20px;
  border-top:1px solid var(--c-border);
}
.version {
  font-size:11px; color:var(--c-text-4); text-align:center;
  letter-spacing:.3px; font-weight:500;
}

.main-content {
  flex:1; overflow-y:auto; overflow-x:hidden;
  background:var(--c-bg);
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
