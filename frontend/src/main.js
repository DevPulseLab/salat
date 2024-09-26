import { createApp } from 'vue'
import { createWebHistory, createRouter } from 'vue-router'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'

import '@/style.css'
import 'primeicons/primeicons.css'

import App from '@/App.vue'
import routes from '@/routes.js'

const router = createRouter({
  history: createWebHistory(),
  routes,
})

createApp(App)
  .use(PrimeVue, {
    theme: {
        preset: Aura
    }
  })
  .use(router)
  .mount('#app')
