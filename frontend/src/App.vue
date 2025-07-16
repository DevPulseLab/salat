<script setup>
import {onMounted, provide, ref, watch} from 'vue'
import Message from 'primevue/message'
import Menubar from 'primevue/menubar'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'

import AppMessage from '@/components/AppMessage.vue'
import PenaltyCard from '@/components/PenaltyCard.vue'
import {getRoutes} from '@/routes.js'

import {useRouter, useRoute} from 'vue-router'
import useAppStore from '@/stores/appStore.js'
import useUserStore from '@/stores/userStore.js'

const appStore = useAppStore()
const userStore = useUserStore()

const route = useRoute()
const router = useRouter()

provide('config', {
  VIEW_DATE_FORMAT: 'DD.MM.YYYY',
  VIEW_DATE_WITH_WEEKDAY_FORMAT: 'ddd DD.MM.YYYY',
  VIEW_MONTH_FORMAT: 'MMMM YYYY',
  DATE_FORMAT: 'YYYY-MM-DD',
  DATETIME_FORMAT: 'YYYY-MM-DD\\THH:mm:ss\\Z',
  calendar: {
    status: {
      noentry: 'noentry',
      rejected: 'rejected',
      approved: 'approved',
      reserved: 'reserved',
    },
    statusText: {
      noentry: '---',
      rejected: 'Abgelehnt',
      approved: 'Genehmigt',
      reserved: 'Erwartet',
    }
  }
})

const routes = ref({})
const showTokenExpired = ref(false)

onMounted(() => {
  routes.value = getRoutes()
  if (appStore.isDarkModeEnabled.value) {
    appStore.enableDarkMode()
  }
})

watch(userStore.isAuthenticated, isAuthenticated => {
  if (!isAuthenticated && route.meta.requiresAuth) {
    showTokenExpired.value = true
    window.setTimeout(() => {
      showTokenExpired.value = false
      router.replace({name: 'home'})
    }, 2000)
  }
})
</script>

<template>
  <div class="container p-2 xl:mt-10 md:mx-auto">
    <button label="Toggle Dark Mode" class="float-right" @click="appStore.toggleDarkMode()">
      <i class="pi" :class="{'pi-moon': !appStore.isDarkModeEnabled.value, 'pi-sun': appStore.isDarkModeEnabled.value}"></i>
    </button>

    <h1 class="text-3xl font-bold mb-4 text-center">
      SalatBar App
    </h1>

    <Menubar :model="routes" class="mb-4" v-if="Object.values(routes).length > 1">
      <template #item="{ item, props }">
        <router-link v-slot="{ href, navigate }" :to="{ name: item.name }">
          <a :href="href" v-bind="props.action" @click="navigate">
            <span class="ml-2">{{ item.meta.label }}</span>
          </a>
        </router-link>
      </template>
    </Menubar>

    <Toast />
    <ConfirmDialog />
    <AppMessage />
    <PenaltyCard />

    <div v-if="!showTokenExpired">
      <RouterView />
    </div>
    <div v-else class="text-center">
      <Message severity="info">Dein Token ist abgelaufen.<br />Du wirst automatisch neu eingeloggt ...</Message>
    </div>
  </div>
</template>