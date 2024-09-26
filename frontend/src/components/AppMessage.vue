<script setup>
import useAppStore from '@/stores/appStore'
import Message from 'primevue/message'
const { appMessage } = useAppStore()
</script>

<template>
  <Message v-if="appMessage !== null" class="my-4" :severity="appMessage.type">{{ appMessage.message }}</Message>
</template>