import { createRouter, createWebHistory } from 'vue-router'
import Twitch from '../views/Twitch.vue'
import Config from '../views/Config.vue'
import Templates from '../views/Templates.vue'
import Videos from '../views/Videos.vue'
import Files from '../views/Files.vue'
import Users from '../views/Users.vue'

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Twitch },
    { path: '/templates', component: Templates },
    { path: '/videos', component: Videos },
    { path: '/files', component: Files },
    { path: '/users', component: Users },
    { path: '/config', component: Config }
  ]
})
