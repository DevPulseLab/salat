import { createApp } from 'vue'
import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'

import '@/style.css'
import 'primeicons/primeicons.css'

import App from '@/App.vue'
import { router } from '@/routes.js'

createApp(App)
  .use(PrimeVue, {
    theme: {
        preset: Aura
    }
  })
  .use(router)
  .mount('#app')
