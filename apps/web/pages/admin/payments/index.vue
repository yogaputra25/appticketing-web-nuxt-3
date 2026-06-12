<template>
  <div>
    <h2 class="text-xl font-semibold text-gray-900 mb-6">Payments</h2>

    <div class="overflow-x-auto -mx-4 md:-mx-0">
      <div class="inline-block min-w-full align-middle px-4 md:px-0">
        <table class="w-full min-w-[600px]">
          <thead>
            <tr class="border-b border-gray-200 text-left text-sm text-gray-500">
              <th class="pb-3 font-medium">ID</th>
              <th class="pb-3 font-medium">Booking</th>
              <th class="pb-3 font-medium">Metode</th>
              <th class="pb-3 font-medium">Status</th>
              <th class="pb-3 font-medium">Jumlah</th>
              <th class="pb-3 font-medium">Tanggal</th>
            </tr>
          </thead>
          <tbody>
            <tr class="border-b border-gray-100">
              <td class="py-3 text-gray-500 text-sm" colspan="6">Data akan muncul setelah backend diimplementasikan.</td>
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
