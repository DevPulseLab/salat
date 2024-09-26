<script setup>
import http from '@/http'
import { ref } from 'vue'
import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Card from 'primevue/card'
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
    password.value = ''
    return
  }

  username.value = ''
  password.value = ''

  authStore.setUserToken(response.data.token)
  appStore.setAppMessage(200, 'Login success')
}
</script>

<template>
  <Card class="max-w-[400px] w-full mx-auto">
    <template #title>Login</template>
    <template #content>
      <InputGroup class="mb-4">
        <InputGroupAddon>
          <i class="pi pi-user"></i>
        </InputGroupAddon>
        <InputText placeholder="Username" v-model="username" />
      </InputGroup>

      <InputGroup class="mb-4">
        <InputGroupAddon>
          <i class="pi pi-key"></i>
        </InputGroupAddon>
        <Password placeholder="Password" v-model="password" :feedback="false" />
      </InputGroup>

      <div class="flex justify-center">
        <Button label="Login" @click="login" raised />
      </div>
    </template>
  </Card>
</template>