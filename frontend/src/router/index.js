import { createRouter, createWebHistory } from 'vue-router'
import Twitch from '../views/Twitch.vue'
import Config from '../views/Config.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Twitch },
    { path: '/config', component: Config }
  ]
})
