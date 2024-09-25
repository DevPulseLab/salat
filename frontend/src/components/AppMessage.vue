<script setup>
import useAppStore from '@/stores/appStore'

const { appMessage } = useAppStore()
</script>

<template>
    <div class="app-message" :class="appMessage.type" v-if="appMessage !== null" v-html="appMessage.message"></div>
</template>

<style scoped>
.app-message {
    padding: 10px;
}

.error {
  background-color: #fac9c9;
}

.warning {
  background-color: #fafabc;
}

.success {
  background-color: #c1fabc;
}
</style>