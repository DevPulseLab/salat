<script setup>
import { inject, onMounted, ref, computed } from 'vue'
import moment from 'moment'
import useUsersService from '@/services/usersService.js'
import CardSelector from '@/components/CardSelector.vue'
import Checkbox from 'primevue/checkbox'
import { useToast } from 'primevue/usetoast'
import InputText from 'primevue/inputtext'

const toast = useToast()
const usersService = useUsersService()
const today = moment()

const appConfig = inject('config')
const todayFormatted = today.format(appConfig.DATE_FORMAT)

const loading = ref(true)
const users = ref([])
const entries = ref([])
const checkedUsers = ref([])

const searchQuery = ref('');
const filteredUsers = computed(() => {
  const userList = users.value.filter(user => entries.value.get(todayFormatted + '_' + user.id))

  if (!searchQuery.value) {
    return userList;
  }

  return userList.filter(user =>
    usersService.getNameFromEmail(user).toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

const checkUser = async (userId) => {
  const visitDate = today.startOf('day').format(appConfig.DATETIME_FORMAT)

  try {
    const success = await usersService.setCheckboxValue(visitDate, userId)
    if (success) {
      toast.add({ severity: 'success', summary: 'Nutzer wurde eingecheckt!', life: 2000 })
    } else {
      toast.add({ severity: 'error', summary: 'Nutzer konnte nicht eingecheckt werden.', life: 2000 })
    }
  } catch (error) {
    console.error('Error checking user:', error);
    toast.add({ severity: 'error', summary: 'Error. Nutzer konnte nicht eingecheckt werden.', life: 2000 })
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
  users.value = Array.from(usersMap.entries()).map(([id, user]) => ({ id, ...user }));
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
    <div class="mx-auto mb-4" v-if="!loading">
      <!-- Header -->
      <div class="font-bold text-center py-1 bg-white border border-gray-200 mb-2" style="border-radius: 8px;">
        {{ today.format(appConfig.VIEW_DATE_FORMAT) }}
      </div>

      <!-- Body -->
      <div class="grid grid-cols-3 gap-3" style="top: 10px;">
        <!-- No entries case -->
        <div v-if="filteredUsers.length === 0" class="px-24 py-4 w-full text-center">
          Keine Einträge für heute
        </div>

        <!-- User entries -->
        <div v-else v-for="user in filteredUsers" class="flex bg-white border border-gray-200" style="border-radius: 8px" :key="user.id">
          <!-- User name and card selector -->
          <div class="px-2 py-1 w-[200px] flex items-center">
            <CardSelector :initialColor="user.penaltyCard" :userId="user.id"
              @card-updated="(color) => user.penaltyCard = color" class="mr-2" />
            {{ usersService.getNameFromEmail(user) }}
          </div>

          <!-- Status -->
          <div class="px-2 py-1 w-[150px] text-center flex items-center justify-center">
            {{ appConfig.calendar.statusText[entries.get(todayFormatted + '_' + user.id)?.status] ?? '---' }}
          </div>

          <!-- Checkbox -->
          <div class="px-2 py-1 w-[100px] text-center flex items-center justify-center">
            <Checkbox v-model="checkedUsers" :value="user.id" @click="checkUser(user.id)" size="large" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>