<template>
  <NuxtLink
    :to="`/events/${event.id}`"
    class="card p-0 overflow-hidden hover:shadow-md transition-shadow group"
  >
    <div class="h-40 bg-gradient-to-br from-primary-500 to-primary-700 flex items-center justify-center relative">
      <span class="text-white text-4xl font-bold">{{ event.title.charAt(0) }}</span>
      <div class="absolute top-3 right-3 flex gap-1.5">
        <span
          class="badge text-xs"
          :class="statusClass"
        >
          {{ statusLabel }}
        </span>
      </div>
    </div>
    <div class="p-4">
      <h3 class="font-semibold text-gray-900 mb-1 group-hover:text-primary-600 transition-colors">
        {{ event.title }}
      </h3>
      <div class="flex items-center gap-1 text-sm text-gray-500 mb-1">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
        </svg>
        {{ event.venue }}
      </div>
      <div class="flex items-center gap-1 text-sm text-gray-400">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        {{ formatDate(event.start_date) }}
      </div>
    </div>
  </NuxtLink>
</template>

<script setup lang="ts">
import type { Event } from '~/stores/event'

const props = defineProps<{
  event: Event
}>()

const statusClass = computed(() => {
  switch (props.event.status) {
    case 'published': return 'bg-green-100 text-green-800'
    case 'draft': return 'bg-yellow-100 text-yellow-800'
    case 'cancelled': return 'bg-red-100 text-red-800'
    case 'finished': return 'bg-gray-100 text-gray-800'
    default: return 'bg-gray-100 text-gray-800'
  }
})

const statusLabel = computed(() => {
  switch (props.event.status) {
    case 'published': return 'Published'
    case 'draft': return 'Draft'
    case 'cancelled': return 'Cancelled'
    case 'finished': return 'Finished'
    default: return props.event.status
  }
})

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}
</script>
