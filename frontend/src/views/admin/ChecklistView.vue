<script setup>
import {inject, onMounted, ref, computed} from 'vue'
import moment from 'moment'
import useUsersService from '@/services/usersService.js'
import CardSelector from '@/components/CardSelector.vue'
import Checkbox from 'primevue/checkbox'
import {useToast} from 'primevue/usetoast'
import InputText from 'primevue/inputtext'

const toast = useToast()
const usersService = useUsersService()
const today = moment()

const appConfig = inject('config')

const loading = ref(true)
const users = ref([])
const entries = ref([])
const checkedUsers = ref([])


const searchQuery = ref('');
const filteredUsers = computed(() => {
  if (!searchQuery.value) {
    return users.value;
  }
  return users.value.filter(user =>
      usersService.getNameFromEmail(user).toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const checkUser = async (userId) => {
  const visitDate = today.startOf('day').format(appConfig.DATETIME_FORMAT)

  try {
    const success = await usersService.setCheckboxValue(visitDate, userId)
    if (success) {
      toast.add({severity: 'success', summary: 'Nutzer wurde eingecheckt!', life: 2000})
    } else {
      toast.add({severity: 'error', summary: 'Nutzer konnte nicht eingecheckt werden.', life: 2000})
    }
  } catch (error) {
    console.error('Error checking user:', error);
    toast.add({severity: 'error', summary: 'Error. Nutzer konnte nicht eingecheckt werden.', life: 2000})
  }
};

const loadCheckboxValues = async () => {
  try {
    const checkboxValues = await usersService.getCheckboxValue(today, today)
    checkboxValues.forEach(item => {
      if (!item.isVisit || checkedUsers.value.includes(item.userId)) {
        return
      }

      checkedUsers.value.push(item.userId)
    });
  } catch (error) {
    console.error('Error during load checkbox values:', error);
  }
}

onMounted(async () => {
  const usersMap = await usersService.fetchUsers();
  users.value = Array.from(usersMap.entries()).map(([id, user]) => ({id, ...user}));
  entries.value = await usersService.fetchUserEntries(today, today)

  await loadCheckboxValues()

  loading.value = false
})
</script>

<template>
  <div class="flex justify-center my-4">
    <InputText v-model="searchQuery" placeholder="Nutzer suchen..." class="w-1/3" />
  </div>
  <div>
    <table class="mx-auto table-auto mb-4" v-if="!loading">
      <thead>
      <tr>
        <th colspan="3">{{ today.format(appConfig.VIEW_DATE_FORMAT) }}</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="user in filteredUsers">
        <td class="px-2 py-1 w-[200px] border">
          {{ usersService.getNameFromEmail(user) }}
          <CardSelector
              :initialColor="user.penaltyCard"
              :userId="user.id"
              @card-updated="(color) => user.penaltyCard = color"
          />
        </td>
        <td class="px-2 py-1 w-[150px] border text-center">{{ appConfig.calendar.statusText[entries.get(today.format(appConfig.DATE_FORMAT) + '_' + user.id)?.status] ?? '---' }}</td>
        <td class="px-2 py-1 w-[100px] border text-center">
          <Checkbox v-model="checkedUsers" :value="user.id" @click="checkUser(user.id)" size="large" />
        </td>
      </tr>
      </tbody>
    </table>
  </div>
</template>