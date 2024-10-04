<script setup>
import {ref} from 'vue'
import {useRouter} from 'vue-router'
import useAppStore from '@/stores/appStore'
import useUserService from '@/services/userService'

import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Card from 'primevue/card'

const appStore = useAppStore()
const userService = useUserService()
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

  const response = await userService.register(registerForm.value.username, registerForm.value.password)
  if (!response.success) {
    return
  }

  registerForm.value = defaultRegisterForm()
  appStore.setAppMessage(200, response.message, 1500)

  window.setTimeout(() => {
    router.replace({ name: 'home' })
  }, 2000)
}
</script>

<template>
  <Card class="max-w-[600px] mx-auto">
    <template #title>Register</template>
    <template #content>
      <InputGroup class="mb-4">
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText placeholder="Username" v-model="registerForm.username" />
      </InputGroup>

      <InputGroup class="mb-4">
        <InputGroupAddon>
          <i class="pi pi-key"></i>
        </InputGroupAddon>
        <Password placeholder="Password" v-model="registerForm.password" :feedback="false" />
      </InputGroup>

      <InputGroup class="mb-4">
        <InputGroupAddon>
          <i class="pi pi-key"></i>
        </InputGroupAddon>
        <Password placeholder="Confirm password " v-model="registerForm.confirmPassword" :feedback="false" />
      </InputGroup>

      <div class="flex justify-center">
        <Button label="Register" @click="register" raised />
      </div>
    </template>
  </Card>
</template>