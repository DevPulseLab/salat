<script setup>
import http from '@/http'
import { onMounted, ref } from 'vue'

import AppMessage from '@/components/AppMessage.vue'

const pong = ref('...')

onMounted(async () => {
  const response = await http.get('/api/ping')
  if (response.status !== 200) {
    pong.value = ''
    return
  }
  pong.value = response.data.ping
})
</script>

<template>
  <AppMessage />

  <RouterView />

  <hr />

  <div>Server: {{ pong }}</div>
</template>
