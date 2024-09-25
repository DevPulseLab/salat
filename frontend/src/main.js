import { createApp } from 'vue'
import { createWebHistory, createRouter } from 'vue-router'
import '@/style.css'

import App from '@/App.vue'
import routes from '@/routes.js'

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')
