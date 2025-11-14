<template>
  <div class="mt-4 section mb-4 bg-white">
    <div class="p-4">
      <h2 class="text-xl font-bold mb-2">Regeln für die Eintragung</h2>
      <p class="mb-2">
        Freitag bis 12 Uhr kann ein Salat für die folgende Woche eingetragen werden.<br />
        Einträge die nach dieser Zeit getätigt werden, werden automatisch als "reserviert"
        <span class="inline-block text-center w-8 h-8 py-1 event-reserved"><img src="@/assets/salat.svg" class="inline-block" alt="salat icon"></span> markiert.
      </p>
      <p>
        Einträge mit einem Papierkorb können noch von Dir selber gelöscht werden.
      </p>
    </div>
  </div>

  <div class="mb-4">
    <div class="mb-4 section flex items-center space-x-4 bg-white">
      <div class="inline-block p-4 event-approved">
        <img src="@/assets/salat.svg" alt="salat icon">
      </div>
      <p>
        Der Salat wurde bestätigt.
      </p>
    </div>

    <div class="mb-4 section flex items-center space-x-4 bg-white">
      <div class="inline-block p-4 event-reserved">
        <img src="@/assets/salat.svg" alt="salat icon">
      </div>
      <p>
        Der Salat wurde reserviert, muss aber vom Admin noch bestätigt werden.
      </p>
    </div>

    <div class="mb-4 section flex items-center space-x-4 bg-white">
      <div class="inline-block p-4 event-rejected">
        <img src="@/assets/salat.svg" alt="salat icon">
      </div>
      <p>
        Der Salat wurde leider abgelehnt.
      </p>
    </div>

    <div class="mb-4 section flex items-center space-x-4 bg-white">
      <div class="inline-block p-4 event-cancelled">
        <img src="@/assets/salat.svg" alt="salat icon">
      </div>
      <p>
        Der Salat wurde leider abgesagt.
      </p>
    </div>
  </div>
</template>