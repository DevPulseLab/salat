<script setup>
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'
import fullcalendarDe from '@fullcalendar/core/locales/de'
import {ref, onMounted, inject} from 'vue'

import {useConfirm} from 'primevue/useconfirm'
import useCalendarService from '@/services/calendarService.js'
import moment from 'moment'
import useAppStore from '@/stores/appStore.js'

import Legend from '@/components/Legend.vue'

import TrashIcon from '@/components/icon/Trash.vue'
import salatIcon from '@/assets/salat.svg';

import unicorn1 from '@/assets/unicorn.png';
import unicorn3 from '@/assets/unicorn2.png';
import unicorn2 from '@/assets/unicorn1.png';
import unicorn4 from '@/assets/unicorn3.png';
import unicorn5 from '@/assets/unicorn4.png';

const unicornImages = [unicorn1, unicorn2, unicorn3, unicorn4, unicorn5];
const randomUnicornImage = unicornImages[Math.floor(Math.random() * unicornImages.length)];

const appConfig = inject('config')

const appStore = useAppStore()
const calendarService = useCalendarService()
const confirm = useConfirm()

const today = moment()
const currentWeek = today.isoWeek()
const disableNextWeek = today.isoWeekday() >= 5 && today.hour() >= 12
const absenceDates = ref([])
const calendarContainer = ref(null)
const isLoading = ref(true)

const addEvent = (calendarApi, id, startDate, endDate, status) => {
  calendarApi.addEvent({
    id: id,
    title: "Salat",
    start: startDate.format(appConfig.DATE_FORMAT),
    end: endDate.format(appConfig.DATE_FORMAT),
    classNames: ['event-' + status],
    allDay: false
  })
}

const removeEvent = async (eventId) => {
  const eventIdAsNumber = Number(eventId);
  return await calendarService.removeEvent(eventIdAsNumber);
}

const getTooltipMessage = (classNames) => {
  if (classNames.includes('event-approved')) {
    return 'Dein Eintrag wurde genehmigt. Guten Appetit!:)';
  } else if (classNames.includes('event-rejected')) {
    return 'Dein Eintrag wurde abgelehnt. Trag dich nächstes Mal rechtzeitig ein.';
  } else if (classNames.includes('event-reserved')) {
    return 'Dein Eintrag wurde reserviert. Wir schauen, ob es genug Salat für alle gibt. Trag dich nächstes Mal rechtzeitig ein.';
  } else if (classNames.includes('event-cancelled')) {
    return 'Dein Eintrag wurde abgesagt.';
  }
  return 'Keine Information verfügbar.';
}

const isFullWindow = () => {
  return window.innerWidth >= 820
}

const updateCalendarSize = (calendar) => {
  if (isFullWindow()) {
    calendar.changeView('dayGridMonth')
    calendar.setOption('height', 'auto')
  } else {
    calendar.changeView('dayGridWeek')
    calendar.setOption('height', 390)
  }
}

const canDeleteEvent = (event) => {
  const weekNumber = moment(event.start).isoWeek()

  return !(
      weekNumber <= currentWeek ||
      event.classNames.includes('event-rejected') ||
      (event.classNames.includes('event-approved') && weekNumber === currentWeek) ||
      (event.classNames.includes('event-approved') && disableNextWeek && weekNumber === currentWeek + 1)
  )
}

const calendarOptions = {
  plugins: [dayGridPlugin, interactionPlugin],
  initialView: isFullWindow() ? 'dayGridMonth' : 'dayGridWeek',
  height: isFullWindow() ? 590 : 390,
  selectable: true,
  validRange: {
    start: moment().startOf('week').toISOString(),
    end: moment().startOf('week').add(3, 'weeks').endOf('week').toISOString()
  },
  headerToolbar: {
    left: '',
    center: 'title',
    right: 'today prev,next'
  },
  initialDate: moment().toISOString(),
  eventClick: (info) => {
    if (!canDeleteEvent(info.event)) {
      return
    }

    const eventId = info.event.id

    confirm.require({
      message: 'Bist du dir sicher, dass du diesen Eintrag löschen möchtest?',
      header: 'Eintrag löschen',
      acceptLabel: 'Löschen',
      rejectLabel: 'Abbrechen',
      accept: async () => {
        const response = await removeEvent(eventId);

        if (response.status === 200) {
          info.event.remove();
          appStore.setAppMessage(200, 'Dieser Eintrag wurde erfolgreich gelöscht.');
        } else {
          appStore.setAppMessage(400, 'Es können keine Einträge aus der Vergangenheit oder dieser Woche gelöscht werden');
        }
      },

      reject: () => appStore.setAppMessage(400, 'Löschung abgebrochen'),
      rejectClass: 'p-button-secondary'
    });
  },
  hiddenDays: [0, 6],
  firstDay: 1,
  locale: fullcalendarDe,
  weekNumbers: true,
  defaultAllDay: true,
  dayHeaderFormat: {
    weekday: 'long'
  },

  windowResize: (info) => {
    updateCalendarSize(info.view.calendar)
  },

  dayCellClassNames: function (info) {
    const weekNumber = moment(info.date).isoWeek(),
        isAbsenceWeek = absenceDates.value.some(absence => {
          const start = moment(absence.startDate),
              end = moment(absence.endDate);
          return moment(info.date).isBetween(start, end, 'day', '[]');
        });

    if (weekNumber <= currentWeek) {
      return ['disallow-week']
    }

    if (disableNextWeek && weekNumber === currentWeek + 1) {
      return ['disallow-week']
    }

    if (isAbsenceWeek) {
      return ['absence-week'];

    }

    return ['allow-week']
  },

  viewDidMount: (info) => {
    info.view.calendar.removeAllEvents()
  },

  datesSet: async (info) => {
    const calenderApi = info.view.calendar

    calenderApi.removeAllEvents()
    const userEvents = await calendarService.getEvents(
        moment(info.start),
        moment(info.end)
    )

    userEvents.forEach(entry => {
      addEvent(calenderApi, entry.id, entry.startDate, entry.endDate, entry.status)
    })
  },

  select: async (selectInfo) => {
    const calendarApi = selectInfo.view.calendar,
        weekNumber = moment(selectInfo.start).isoWeek(),
        currentDayOfWeek = today.isoWeekday(),
        absence = absenceDates.value.find(absence => {
          const start = moment(absence.startDate),
              end = moment(absence.endDate);
          return moment(selectInfo.start).isBetween(start, end, 'day', '[]');
        });
    calendarApi.unselect()

    if (absence) {
      const startDateFormatted = moment(absence.startDate).format('DD.MM.YYYY'),
          endDateFormatted = moment(absence.endDate).format('DD.MM.YYYY');

      confirm.require({
        message: `Lars ist ab dem ${startDateFormatted} bis zum ${endDateFormatted} nicht da`,
        header: 'Info',
        acceptLabel: 'Ok',
        rejectClass: '!hidden'
      });
      return;
    }

    if (weekNumber < currentWeek || (weekNumber === currentWeek && currentDayOfWeek >= 5 && currentDayOfWeek <= 7)) {
      confirm.require({
        message: 'Eintrag für diese Woche nicht mehr möglich',
        header: 'Not possible',
        acceptLabel: 'Ok',
        rejectClass: '!hidden'
      })
      return
    }

    const startDate = moment(selectInfo.startStr),
        endDate = moment(selectInfo.endStr).subtract(1, 'days')

    if (calendarApi.getEventById(startDate.format(appConfig.DATE_FORMAT)) instanceof Object) {
      confirm.require({
        message: 'Du hast an diesem Tag bereits einen Eintrag gemacht',
        header: 'Duplicate entry',
        acceptLabel: 'Ok',
        rejectClass: '!hidden'
      })
      return
    }

    const response = await calendarService.addEvent(startDate, endDate)

    if (response.status === 200) {
      response.data.calendarEntries.forEach(entry => {
        addEvent(calendarApi, entry.id, moment(entry.date), moment(entry.date), entry.status)
      })

      appStore.setAppMessage(200, 'Kalenderdaten wurden gespeichert')
    } else {
      appStore.setAppMessage(400, 'Kalenderdaten konnten nicht gespeichert werden')
    }
  }
}

onMounted(async () => {
  absenceDates.value = await calendarService.fetchAbsences(moment().startOf('year'), moment().endOf('year'));

  if (calendarContainer.value) {
    calendarContainer.value.getApi().refetchEvents();
  }

  isLoading.value = false;

  const styleElement = document.createElement('style')
  styleElement.textContent = `
    .fc .fc-daygrid-day.fc-day-today::before {
      content: '';
      background-image: url('${randomUnicornImage}');
      position: absolute;
      width: 64px;
      height: 64px;
      z-index: 10;
      background-repeat: no-repeat;
      right: 5px;
      bottom: 0;
    }
  `;
  document.head.appendChild(styleElement);
});
</script>

<template>
  <div class="relative">
    <p v-if="isLoading">Loading...</p>
    <FullCalendar v-else ref="calendarContainer" :options="calendarOptions">
      <template #eventContent="arg">
        <div class="calendar-entry text-white" v-tooltip.bottom="getTooltipMessage(arg.event.classNames)">
          <img class="salat" :src="salatIcon" alt="salat icon">
          {{ arg.event.title }}
          <TrashIcon v-if="canDeleteEvent(arg.event)" />
        </div>
      </template>
    </FullCalendar>

    <Legend />
  </div>
</template>

<style>
.p-dialog-footer {
  justify-content: flex-start !important;
  flex-direction: row-reverse !important;
}

.p-tooltip {
  font-size: 12px;
  min-width: 295px;
  padding: 8px;
}

.fc-daygrid-dot-event:hover {
  cursor: pointer;
}

.fc .fc-day-disabled {
  background: #ffff;
}

.fc .fc-daygrid-day.fc-day-today {
  background-color: #d2c288;
  position: relative;
}

.fc .fc-daygrid-week-number {
  background-color: transparent;
}

.fc .fc-daygrid-day-number {
  color: #474242;
}

.calendar-entry {
  display: flex;
  align-items: baseline;
  padding: 5px 5px 5px 15px;
  font-size: 18px;
  color: #fff;
  width: 100%;
}

.calendar-entry .salat {
  margin-right: 7px
}

.calendar-entry .icon-trash {
  margin-left: auto;
}

.disallow-week {
  background-color: #f6e0aa;
}

.allow-week {
  background-color: #b9dea7;
}

.absence-week {
  background-color: #e3e3ec;
}

.event-approved {
  background-color: #4faa2f;
}

.event-rejected {
  background-color: #aa2f2f;
}

.event-cancelled {
  background-color: #4a5565;
}

.event-reserved {
  background-color: #eec318;
}
</style>