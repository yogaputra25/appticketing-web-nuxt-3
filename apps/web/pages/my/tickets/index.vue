<template>
  <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-6 md:py-8">
    <div class="mb-6 md:mb-8">
      <h1 class="text-2xl md:text-3xl font-bold text-gray-900">Tiket Saya</h1>
      <p class="text-gray-500 mt-1 text-sm md:text-base">Semua tiket yang sudah kamu beli</p>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-8 h-8 border-4 border-primary-600 border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="tickets.length === 0" class="text-center py-16">
      <svg class="w-16 h-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 5H9a2 2 0 00-2 2v10a2 2 0 002 2h6a2 2 0 002-2V7a2 2 0 00-2-2z" />
      </svg>
      <p class="text-gray-500 text-lg">Belum ada tiket</p>
      <NuxtLink to="/events" class="text-primary-600 hover:text-primary-700 mt-2 inline-block font-medium">
        Cari Event &rarr;
      </NuxtLink>
    </div>

    <div v-else>
      <div class="space-y-3 md:hidden">
        <NuxtLink
          v-for="ticket in tickets"
          :key="ticket.id"
          :to="`/my/tickets/${ticket.id}`"
          class="card p-4 block hover:shadow-md transition-shadow"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <h3 class="font-semibold text-gray-900 truncate">{{ ticket.event_title }}</h3>
              <p class="text-sm text-gray-500 mt-0.5">{{ ticket.category_name }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(ticket.created_at) }}</p>
            </div>
            <span class="badge shrink-0" :class="statusClass(ticket.status)">{{ statusLabel(ticket.status) }}</span>
          </div>
        </NuxtLink>
      </div>

      <div class="hidden md:block overflow-x-auto">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
              <th class="pb-3 font-medium">Event</th>
              <th class="pb-3 font-medium">Kategori</th>
              <th class="pb-3 font-medium">Status</th>
              <th class="pb-3 font-medium">Tanggal</th>
              <th class="pb-3 font-medium" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="ticket in tickets" :key="ticket.id" class="border-b border-gray-100">
              <td class="py-3 font-medium text-gray-900">{{ ticket.event_title }}</td>
              <td class="py-3 text-sm text-gray-600">{{ ticket.category_name }}</td>
              <td class="py-3">
                <span class="badge" :class="statusClass(ticket.status)">{{ statusLabel(ticket.status) }}</span>
              </td>
              <td class="py-3 text-sm text-gray-500">{{ formatDate(ticket.created_at) }}</td>
              <td class="py-3">
                <NuxtLink :to="`/my/tickets/${ticket.id}`" class="text-primary-600 hover:text-primary-700 text-sm font-medium">
                  Detail
                </NuxtLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth',
})

const ticketStore = useTicketStore()
const tickets = computed(() => ticketStore.tickets)
const loading = computed(() => ticketStore.loading)
const total = computed(() => ticketStore.total)

function statusClass(status: string) {
  const classes: Record<string, string> = {
    active: 'bg-green-100 text-green-800',
    used: 'bg-gray-100 text-gray-800',
    refunded: 'bg-red-100 text-red-800',
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

function statusLabel(status: string) {
  const labels: Record<string, string> = {
    active: 'Aktif',
    used: 'Terpakai',
    refunded: 'Refund',
  }
  return labels[status] || status
}

function formatDate(dateStr: string) {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

onMounted(async () => {
  await ticketStore.fetchMyTickets()
})
</script>
