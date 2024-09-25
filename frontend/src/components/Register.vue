<script setup>
import {ref} from 'vue'
import http from '@/http'
import {useRouter} from 'vue-router'
import useAppStore from '@/stores/appStore'

const appStore = useAppStore()
const router = useRouter()

const defaultRegisterForm = () => {
    return {
        username: '',
        password: '',
        confirmPassword: ''
    }
}

const registerForm = ref(defaultRegisterForm())

const register = async () => {
    if (registerForm.value.username === '') {
        appStore.setAppMessage(400, 'Invalid username')
        return
    }

    if (registerForm.value.password === '' || registerForm.value.password !== registerForm.value.confirmPassword) {
        appStore.setAppMessage(400, 'Invalid password')
        return
    }

    const response = await http.post('/api/register', {
        username: registerForm.value.username,
        password: registerForm.value.password
    })

    if (response.status !== 200) {
        return
    }

    appStore.setAppMessage(200, response.data.message, 1500)
    registerForm.value = defaultRegisterForm()

    window.setTimeout(() => {
        router.replace({ name: 'home' })
    }, 2000)
}

</script>

<template>
    <h1>New account</h1>
    <div>
      <input type="text" v-model="registerForm.username">
    </div>
    <div>
      <input type="password" v-model="registerForm.password">
    </div>
    <div>
      <input type="password" v-model="registerForm.confirmPassword">
    </div>
    <div>
      <button @click="register">Register</button>
    </div>
</template>