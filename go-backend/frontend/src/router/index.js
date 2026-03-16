import { createRouter, createWebHistory } from 'vue-router'
import Streamers from '../views/Streamers.vue'
import Templates from '../views/Templates.vue'
import Config from '../views/Config.vue'
import Account from '../views/Account.vue'
import Logs from '../views/Logs.vue'

const routes = [
  { path: '/', name: 'Streamers', component: Streamers },
  { path: '/templates', name: 'Templates', component: Templates },
  { path: '/config', name: 'Config', component: Config },
  { path: '/account', name: 'Account', component: Account },
  { path: '/logs', name: 'Logs', component: Logs },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
