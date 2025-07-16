<script setup>
import {ref, onMounted} from 'vue'
import useUserService from '@/services/userService'
import Message from 'primevue/message'

const userService = useUserService()
const penaltyCardMessage = ref(null)

onMounted(async () => {
  const userInfo = await userService.currentUserInfo()

  if (userInfo && userInfo?.penaltyCard) {
    let message = ''
    let type = ''
    switch (userInfo?.penaltyCard) {
        case 'yellow':
            message = 'Sie haben eine gelbe Karte erhalten.'
            type = 'warn'
            break;
        case 'red':
            message = 'Sie haben eine rote Karte erhalten.'
            type = 'error'
            break;
    }

    penaltyCardMessage.value = { message, type }
  }
})
</script>

<template>
  <Message v-if="penaltyCardMessage !== null" class="my-4" :severity="penaltyCardMessage.type">
    {{ penaltyCardMessage.message }}
  </Message>
</template>
