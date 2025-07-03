<script setup>
import {ref, watch} from 'vue'
import {useToast} from 'primevue/usetoast'
import useUsersService from '@/services/usersService.js'

const props = defineProps({
  userId: {
    type: Number,
    required: true
  },
  initialColor: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['card-updated'])

const selectedCardColor = ref(props.initialColor || '')
const showCardPopup = ref(false)

const usersService = useUsersService()
const toast = useToast()

const togglePopup = () => {
  showCardPopup.value = !showCardPopup.value
}

const closePopup = () => {
  showCardPopup.value = false
}

const selectCardColor = async (color) => {
  try {
    const success = await usersService.setPenaltyCard(color, props.userId)
    if (success) {
      selectedCardColor.value = color
      emit('card-updated', color)
      toast.add({severity: 'success', summary: 'Karte gespeichert', life: 1500})
    } else {
      toast.add({severity: 'error', summary: 'Fehler beim Speichern', life: 1500})
    }
  } catch (error) {
    console.error('Fehler beim Speichern der Karte:', error)
    toast.add({severity: 'error', summary: 'Serverfehler', life: 1500})
  } finally {
    closePopup()
  }
}
</script>

<template>
  <div class="relative inline-block">
    <!--penalty card-->
    <button
        :class="`${selectedCardColor || 'white'}-card`"
        @click="togglePopup"
    ></button>

    <!--pop-up-->
    <div
        v-if="showCardPopup"
        class="absolute z-10 mt-2 bg-white shadow-md rounded p-4 w-max min-w-[180px]"
    >
      <div class="flex justify-end mb-2">
        <button @click="closePopup" class="text-gray-500 hover:text-black text-sm font-bold">×</button>
      </div>
      <div class="flex space-x-4 justify-center">
        <button class="white-card" @click="selectCardColor('')"></button>
        <button class="yellow-card" @click="selectCardColor('yellow')"></button>
        <button class="red-card" @click="selectCardColor('red')"></button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.white-card,
.yellow-card,
.red-card {
  width: 25px;
  height: 30px;
  border: 1px solid #64748b;
  border-radius: 4px;
  display: inline-block;
  transition: transform 0.1s ease;
}

.white-card {
  background-color: white;
}

.yellow-card {
  background-color: #facc15;
}

.red-card {
  background-color: #f87171;
}

.white-card:hover,
.yellow-card:hover,
.red-card:hover {
  transform: scale(1.1);
  cursor: pointer;
}
</style>