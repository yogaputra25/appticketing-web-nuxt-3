<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 mb-6">Overview</h2>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="card p-5">
        <p class="text-sm text-gray-500 mb-1">Total Events</p>
        <p class="text-3xl font-bold text-gray-900">-</p>
      </div>
      <div class="card p-5">
        <p class="text-sm text-gray-500 mb-1">Total Bookings</p>
        <p class="text-3xl font-bold text-gray-900">-</p>
      </div>
      <div class="card p-5">
        <p class="text-sm text-gray-500 mb-1">Revenue</p>
        <p class="text-3xl font-bold text-gray-900">-</p>
      </div>
      <div class="card p-5">
        <p class="text-sm text-gray-500 mb-1">Total Users</p>
        <p class="text-3xl font-bold text-gray-900">-</p>
      </div>
    </div>

    <p class="text-gray-500 text-sm">Data akan muncul setelah backend admin endpoints diimplementasikan.</p>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'admin',
  layout: 'admin',
})
</script>
