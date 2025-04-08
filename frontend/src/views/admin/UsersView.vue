<script setup>
import moment from 'moment'
import {inject, onMounted, ref} from 'vue'
import useUsersService from '@/services/usersService.js'
import {useToast} from 'primevue/usetoast'

import DatePicker from 'primevue/datepicker'
import DaySelect from '@/components/DaySelect.vue'

const usersService = useUsersService()
const toast = useToast()

const appConfig = inject('config')

const entries = ref([])
const users = ref([])

const loading = ref(true)
const editMode = ref(false)

const platesNumbers = ref({});

const startDate = ref(moment().startOf('week'))
const endDate = ref(moment().endOf('week'))
const dateRange = ref(moment.range(startDate.value, endDate.value))

const loadTable = async () => {
  loading.value = true

  entries.value = await usersService.fetchUserEntries(
      startDate.value,
      endDate.value,
  )

  loading.value = false
}

const updateDateRange = async (start, end) => {
  startDate.value = start.startOf('week')
  endDate.value = end.endOf('week')
  dateRange.value = moment.range(startDate.value, endDate.value)
  editMode.value = false

  const weekdays = Array.from(moment.range(startDate.value, endDate.value).by('days')).filter(day => {
    return day.isoWeekday() >= 1 && day.isoWeekday() <= 5;
  });

  for (let day of weekdays) {
    const statsDay = day.format(appConfig.DATE_FORMAT);
    platesNumbers.value[statsDay] = await usersService.fetchNumberOfPlates(statsDay);
  }
  await loadTable()
}

const dayEntry = (day, userId) => {
  const entry = entries.value.get(day.format(appConfig.DATE_FORMAT) + '_' + userId)
  if (entry !== undefined) {
    return entry
  }

  return null
}

const changeUserDayStatus = async (data) => {
  const changed = await usersService.changeEntryStatus(data.id, data.status)
  if (changed) {
    await loadTable()
    toast.add({severity: 'success', summary: 'Status geändert...', life: 2000})
  }
}

const dates = ref([
  startDate.value.toDate(),
  endDate.value.toDate()
])

const changeDates = () => {
  const hasBothDays = dates.value.indexOf(null) === -1
  if (hasBothDays) {
    updateDateRange(moment(dates.value[0]), moment(dates.value[1]))
  }
}

const countApprovedPerDay = (day) => {
  return users.value.reduce((count, user) => {
    const entry = dayEntry(day, user.id);
    if (entry && entry.status === 'approved') {
      count++;
    }
    return count;
  }, 0);
};


const savePlatesNumber = async (day) => {
  try {
    const statsDay = day.format(appConfig.DATETIME_FORMAT);
    const numberOfPlates = parseInt(platesNumbers.value[day.format(appConfig.DATE_FORMAT)], 10);

    const success = await usersService.savePlatesNumber(statsDay, numberOfPlates);

    if (success) {
      toast.add({severity: 'success', summary: 'Anzahl der Teller wurde gespeichert!', life: 2000});
    } else {
      toast.add({severity: 'error', summary: 'Anzahl der Teller konnte nicht gespeichert werden.', life: 2000});
    }
  } catch (error) {
    console.error('Fehler beim Speichern der Telleranzahl:', error);
    toast.add({severity: 'error', summary: 'Beim Speichern der Telleranzahl ist ein Fehler aufgetreten', life: 2000});
  }
}

onMounted(async () => {
  try {
    const usersMap = await usersService.fetchUsers();
    users.value = Array.from(usersMap.entries()).map(([id, user]) => ({id, ...user}));

    users.value.forEach((user) => {
        const [firstName, lastName] = user.username.split('@')[0].split('.');
        user.username = firstName.toUpperCase() + '. ' + lastName.charAt(0).toUpperCase() + lastName.slice(1);
    })

    const weekdays = Array.from(moment.range(startDate.value, endDate.value).by('days')).filter(day => {
      return day.isoWeekday() >= 1 && day.isoWeekday() <= 5;
    });

    for (let day of weekdays) {
      const statsDay = day.format(appConfig.DATE_FORMAT);
      platesNumbers.value[statsDay] = await usersService.fetchNumberOfPlates(statsDay);
    }

    await loadTable();
  } catch (error) {
    console.error('Error during component mounting:', error);
  }
})
</script>

<template>
  <h1 v-if="!editMode" @click="editMode = !editMode" class="text-xl font-bold text-center m-4">
    {{ startDate.format(appConfig.VIEW_DATE_FORMAT) }} bis {{ endDate.format(appConfig.VIEW_DATE_FORMAT) }}
    <i class="pi pi-pencil" style="font-size: 1rem"></i>
  </h1>
  <div v-else class="text-center">
    <DatePicker
        v-model="dates"
        selectionMode="range"
        dateFormat="dd.mm.yy"
        :showButtonBar="true"
        :selectOtherMonths="true"
        :disabledDays="[0, 6]"
        :showWeek="true"
        :manualInput="false"
        @date-select="changeDates"
    />

    <div class="text-gray-500 mt-2 mb-6 text-xs">
      Der ausgesuchte Bereich aktualisiert sich automatisch auf die vollen Kalenderwochen!
    </div>
  </div>

  <table class="mx-auto table-auto mb-4" v-if="!loading" v-for="week in dateRange.by('weeks')">
    <thead>
    <tr>
      <th colspan="6">KW{{ week.isoWeek() }}</th>
    </tr>
    <tr class="border">
      <th class="bg-gray-300 px-2 py-1 w-[200px]">User</th>
      <th
          class="bg-gray-300 px-2 py-1 w-[200px]"
          v-for="day in moment.range(week.clone().startOf('week'), week.clone().endOf('week').subtract(2, 'day')).by('day')"
      >
        {{ day.format(appConfig.VIEW_DATE_FORMAT) }}
      </th>
    </tr>
    </thead>
    <tbody>
    <tr class="border-b" v-for="user in users" :key="user.id">
      <td class="px-2 py-1 w-[200px] border-l">{{ user.username }}</td>
      <td
          class="px-2 py-1 w-[200px] border-l border-r text-center"
          v-for="day in moment.range(week.clone().startOf('week'), week.clone().endOf('week').subtract(2, 'day')).by('day')"
      >
        <day-select v-bind="dayEntry(day, user.id)" @change="changeUserDayStatus"/>
      </td>
    </tr>
    <tr class="statistics border-t font-bold">
      <td class="px-2 py-1 w-[200px] border-l" style="background: #fdfdfd8a">Mit Status "Genehmigt" :</td>
      <td
          class="px-2 py-1 w-[200px] border-l border-r text-center" style="background: #fdfdfd8a"
          v-for="day in moment.range(week.clone().startOf('week'), week.clone().endOf('week').subtract(2, 'day')).by('day')"
      >
        {{ countApprovedPerDay(day) }}
      </td>
    </tr>
    <tr class="statistics border-t font-bold">
      <td class="px-2 py-1 w-[200px] border-l" style="background: #fdfdfd8a">Tatsächliche Anzahl der Leute:</td>
      <td
          class="px-2 py-1 w-[200px] border-l border-r text-center" style="background: #fdfdfd8a"
          v-for="day in moment.range(week.clone().startOf('week'), week.clone().endOf('week').subtract(2, 'day')).by('day')"
      >
        <div class="total-people-group">
          <input
              class="total-people-input"
              :placeholder="platesNumbers[day.format(appConfig.DATE_FORMAT)] || 0"
              v-model="platesNumbers[day.format(appConfig.DATE_FORMAT)]">
          <button class="total-people-button" @click="savePlatesNumber(day)">Ok</button>
        </div>
      </td>
    </tr>
    </tbody>
  </table>
  <div v-else>
    Loading ...
  </div>
</template>

<style>
.total-people-group {
  display: flex;
  justify-content: space-around;
}

.total-people-input {
  max-width: 90px;
  background: #fdfdfd8a;
  border-radius: 5px;
  padding: 3px 15px;
  outline: none;
}

.total-people-button {
  background-color: #7272758f;
  padding: 7px 15px;
  border-radius: 5px;
}

.total-people-button:hover {
  background-color: #727275;
}
</style>