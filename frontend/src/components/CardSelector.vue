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
  }
})

const emit = defineEmits(['card-updated'])

const selectedCardColor = ref(props.initialColor || '')
const showCardPopup = ref(false)

const usersService = useUsersService()
const toast = useToast()

const beforeEnter = (el) => {
  el.style.opacity = 0
  el.style.transform = 'scaleX(0.2)'
  el.style.transformOrigin = 'left center'
}

const enter = (el, done) => {
  el.style.transition = 'all 200ms ease'
  requestAnimationFrame(() => {
    el.style.opacity = 1
    el.style.transform = 'scaleX(1)'
  })
  setTimeout(done, 200)
}

const leave = (el, done) => {
  el.style.transition = 'all 200ms ease'
  el.style.opacity = 0
  el.style.transform = 'scaleX(0.2)'
  setTimeout(done, 200)
}

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
    <transition
        name="popup-fade"
        @before-enter="beforeEnter"
        @enter="enter"
        @leave="leave"
    >
    <div
        v-if="showCardPopup"
        class="absolute z-10 bg-white shadow-md rounded p-1 pop-up"
    >
      <div class="flex space-x-4 justify-center">
        <button class="white-card" @click="selectCardColor('')"></button>
        <button class="yellow-card" @click="selectCardColor('yellow')"></button>
        <button class="red-card" @click="selectCardColor('red')"></button>
        <button @click="closePopup" class="text-gray-500 hover:text-black text-sm font-bold">×</button>
      </div>
    </div>
    </transition>
  </div>
</template>

<style scoped>
.pop-up {
  top: 0;
}

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