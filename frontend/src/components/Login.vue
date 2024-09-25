<script setup>
import http from '@/http'
import { ref } from 'vue'
import useUserStore from '@/stores/userStore'
import useAppStore from '@/stores/appStore'

const appStore = useAppStore()
const authStore = useUserStore()

const username = ref('')
const password = ref('')

const login = async () => {
    const response = await http.post('/api/login', {
        username: username.value,
        password: password.value
    })

    if (response.status !== 200) {
        return
    }

    username.value = ''
    password.value = ''

    authStore.setUserToken(response.data.token)
    appStore.setAppMessage(200, 'Login success')
}
</script>

<template>
    <div>
      <input type="text" v-model="username">
    </div>
    <div>
      <input type="password" v-model="password">
    </div>
    <div>
      <button @click="login">Login</button>
    </div>
</template>