<script setup>
import useUsersService from '@/services/usersService.js'
import {inject, onMounted, ref} from 'vue'
import moment from 'moment'

import DatePicker from 'primevue/datepicker'

const appConfig = inject('config')

const usersService = useUsersService()

const loading = ref(true)
const editMode = ref(false)

const month = ref(moment().startOf('month'))
const stats = ref({})
const statsGuests = ref({})
const totalPlates = ref(0)
const totalGuests = ref(0)

const fetchStats = async () => {
  loading.value = true

  stats.value = {}
  statsGuests.value = {}
  totalPlates.value = 0
  totalGuests.value = 0
  for (let day of Array.from(moment.range(month.value, month.value.clone().endOf('month')).by('days'))) {
    const weekDay = day.isoWeekday()
    if (weekDay >= 6) {
      continue
    }

    const isoWeek = day.isoWeek()
    if (stats.value[isoWeek] === undefined) {
      stats.value[isoWeek] = {}
    }

    const plates = await usersService.fetchNumberOfPlates(day.format(appConfig.DATE_FORMAT))

    stats.value[isoWeek][day.format(appConfig.DATE_FORMAT)] = plates
    totalPlates.value += plates

    if (statsGuests.value[isoWeek] === undefined) {
      statsGuests.value[isoWeek] = {}
    }

    const guests = await usersService.fetchNumberOfGuests(day.format(appConfig.DATE_FORMAT))
    statsGuests.value[isoWeek][day.format(appConfig.DATE_FORMAT)] = guests
    totalGuests.value += guests
  }

  loading.value = false
}

const changeMonth = async selectedMonth => {
  month.value = moment(selectedMonth).startOf('month')
  editMode.value = false

  await fetchStats()
}

const getPlatesFromDayIndex = (week, index) => {
  const daysArray = Object.keys(stats.value[week])
  const platesArray = Object.values(stats.value[week])

  let formatedDate = null
  if (daysArray[index - 1]) {
    formatedDate = moment(daysArray[index - 1]).format(appConfig.VIEW_DATE_WITH_WEEKDAY_FORMAT)
  }

  return {
    day: formatedDate,
    plates: platesArray[index - 1] ?? null
  }
}

const getGuestsFromDayIndex = (week, index) => {
  const daysArray = Object.keys(statsGuests.value[week])
  const guestsArray = Object.values(statsGuests.value[week])

  let formatedDate = null
  if (daysArray[index - 1]) {
    formatedDate = moment(daysArray[index - 1]).format(appConfig.VIEW_DATE_WITH_WEEKDAY_FORMAT)
  }

  return {
    day: formatedDate,
    guests: guestsArray[index - 1] ?? null
  }
}

onMounted(async () => {
  await fetchStats()
})
</script>

<template>
  <h1 v-if="!editMode" @click="editMode = !editMode" class="text-xl font-bold text-center m-4">
    {{ month.format(appConfig.VIEW_MONTH_FORMAT) }}
    <i class="pi pi-pencil" style="font-size: 1rem"></i>
  </h1>
  <div v-else class="text-center">
    <DatePicker
        view="month"
        dateFormat="MM yy"
        :manualInput="false"
        @date-select="changeMonth" />
  </div>

  <p class="text-center mb-8">
    Die Statistiken arbeiten mit den tatsächlichen Werten, welche vom Admin eingetragen werden.
  </p>

  <div v-if="!loading" class="text-center">
    <p class="mb-2">
      <strong>Salate gesamt:</strong> {{ totalPlates + totalGuests }}<br />
      <span class="text-sm">{{ totalPlates }} Teller und {{ totalGuests }} Gäste</span>
    </p>

    <table class="mx-auto table-auto mb-4" v-for="(days, week) in stats">
      <thead>
      <tr class="border">
        <th
            class="bg-gray-300 px-2 py-1 w-[200px] text-center"
            v-for="dayIndex in 5"
        >
          {{ getPlatesFromDayIndex(week, dayIndex).day }}
        </th>
      </tr>
      </thead>
      <tbody>
      <tr class="border-b">
        <td
            class="px-2 py-1 w-[200px] text-center"
            v-for="dayIndex in 5"
        >
          {{ getPlatesFromDayIndex(week, dayIndex).plates }} ({{ getGuestsFromDayIndex(week, dayIndex).guests }})
        </td>
      </tr>
      </tbody>
    </table>
  </div>
  <div v-else>
    Loading ...
  </div>
</template>