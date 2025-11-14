<script setup>
import {ref, inject} from 'vue'
import Select from 'primevue/select'
import Tag from 'primevue/tag'
import moment from 'moment'

const appConfig = inject('config')

const emit = defineEmits(['change'])

const props = defineProps({
  id: {
    type: Number
  },
  date: {
    type: String
  },
  userId: {
    type: Number
  },
  status: {
    type: String,
    default: 'noentry'
  }
})

const status = ref(props.status)
const editMode = ref(false)

const changeDayStatus = () => {
  editMode.value = false
  emit('change', {
    id: props.id,
    userId: props.userId,
    date: moment(props.date),
    status: status.value
  })
}

const options = [
  {"status": appConfig.calendar.statusText.rejected, "value": appConfig.calendar.status.rejected},
  {"status": appConfig.calendar.statusText.reserved, "value": appConfig.calendar.status.reserved},
  {"status": appConfig.calendar.statusText.approved, "value": appConfig.calendar.status.approved},
  {"status": appConfig.calendar.statusText.cancelled, "value": appConfig.calendar.status.cancelled}
]

const classForStatus = () => {
  return {
    'text-green-700': status.value === appConfig.calendar.status.approved,
    'text-red-700': status.value === appConfig.calendar.status.rejected,
    'text-slate-600': status.value === appConfig.calendar.status.cancelled,
    'text-amber-700': status.value === appConfig.calendar.status.reserved,
  }
}
</script>

<template>
  <div v-if="props.status !== appConfig.calendar.status.noentry">
    <Tag v-if="!editMode" @click="editMode=!editMode" severity="secondary" class="w-full">
      <span :class="classForStatus()">{{ appConfig.calendar.statusText[status] }}</span>
      <i class="pi pi-pencil" style="font-size: 1rem"></i>
    </Tag>
    <Select
        v-else
        v-model="status"
        :options="options"
        option-label="status"
        option-value="value"
        size="small"
        @change="changeDayStatus"
        class="w-full"
    />
  </div>
</template>