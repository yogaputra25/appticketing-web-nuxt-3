<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <h2 class="text-xl font-semibold text-gray-900">Events</h2>
      <NuxtLink to="/admin/events/new" class="btn-primary text-sm !py-1.5 !px-3 md:!px-4 touch-target inline-flex items-center gap-1">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span class="hidden sm:inline">Tambah Event</span>
      </NuxtLink>
    </div>

    <div class="overflow-x-auto -mx-4 md:-mx-0">
      <div class="inline-block min-w-full align-middle px-4 md:px-0">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
              <th class="pb-3 font-medium">Nama</th>
              <th class="pb-3 font-medium">Status</th>
              <th class="pb-3 font-medium">Tanggal</th>
              <th class="pb-3 font-medium">Tiket Terjual</th>
              <th class="pb-3 font-medium" />
            </tr>
          </thead>
          <tbody>
            <tr class="border-b border-gray-100">
              <td class="py-3 text-gray-500 text-sm" colspan="5">Data akan muncul setelah backend diimplementasikan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})
</script>
