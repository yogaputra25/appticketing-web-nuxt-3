<template>
  <div class="card p-4 flex items-center justify-between">
    <div class="flex-1">
      <h4 class="font-semibold text-gray-900">{{ category.name }}</h4>
      <p v-if="category.description" class="text-sm text-gray-500 mt-0.5">{{ category.description }}</p>
      <div class="flex items-center gap-3 mt-2">
        <span v-if="isSoldOut" class="badge bg-red-100 text-red-700">Sold Out</span>
        <span v-else class="badge bg-green-100 text-green-700">
          {{ category.available_stock }} tersisa
        </span>
        <span class="text-xs text-gray-400">Max {{ category.max_per_user }} tiket</span>
      </div>
    </div>
    <div class="text-right ml-4">
      <p class="text-xl font-bold text-primary-600">
        {{ formatPrice(category.price) }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { TicketCategory } from '~/stores/event'

const props = defineProps<{
  category: TicketCategory
}>()

const isSoldOut = computed(() => props.category.available_stock <= 0)

function formatPrice(price: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}
</script>
