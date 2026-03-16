<template>
  <div id="app">
    <nav>
      <router-link to="/">Twitch 录播</router-link>
      <router-link to="/config">配置</router-link>
    </nav>
    <router-view />
  </div>
</template>

<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: sans-serif; }
nav { display: flex; gap: 20px; padding: 20px; background: #f5f5f5; }
nav a { text-decoration: none; color: #333; }
nav a.router-link-active { color: #00a1d6; font-weight: bold; }
</style>
