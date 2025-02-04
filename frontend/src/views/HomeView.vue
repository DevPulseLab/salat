<script setup>
import { ref } from 'vue'
import useAppStore from '@/stores/appStore'
import { useRouter } from 'vue-router'
import useUserService from '@/services/userService'
import useUserStore from '@/stores/userStore'

import InputGroup from 'primevue/inputgroup'
import InputGroupAddon from 'primevue/inputgroupaddon'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import Button from 'primevue/button'
import Card from 'primevue/card'

const appStore = useAppStore()
const router = useRouter()
const userService = useUserService()
const { getUser } = useUserStore()

const username = ref('')
const password = ref('')

const login = async () => {
  const responseSuccess = await userService.login(username.value, password.value)
  if (!responseSuccess) {
    password.value = ''
    return
  }

  username.value = ''
  password.value = ''

  appStore.setAppMessage(200, 'Login success')
  router.replace({name: getUser().startRoute})
}
</script>

<template>
  <Card class="max-w-[400px] w-full mx-auto">
    <template #title>Login</template>
    <template #content>
      <form @submit.prevent="login">
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
          <Button label="Login" type="submit" raised />
        </div>
      </form>
    </template>
  </Card>
</template>