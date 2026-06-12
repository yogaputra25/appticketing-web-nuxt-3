<template>
  <div class="card p-5 md:p-6 text-center">
    <div class="mb-6">
      <div class="relative w-36 h-36 md:w-40 md:h-40 mx-auto mb-4">
        <svg class="w-36 h-36 md:w-40 md:h-40 -rotate-90" viewBox="0 0 120 120">
          <circle cx="60" cy="60" r="52" fill="none" stroke="#e5e7eb" stroke-width="8" />
          <circle
            cx="60" cy="60" r="52" fill="none" stroke="#3358ff"
            stroke-width="8" stroke-linecap="round"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="dashOffset"
            class="transition-all duration-700 ease-out"
          />
        </svg>
        <div class="absolute inset-0 flex items-center justify-center flex-col">
          <span class="text-5xl md:text-6xl font-bold text-gray-900">{{ position }}</span>
          <span class="text-sm text-gray-500">antrian</span>
        </div>
      </div>

      <p class="text-gray-600 text-sm md:text-base">
        <template v-if="position === 0">
          Giliranmu sudah tiba! Kamu memiliki waktu terbatas untuk melanjutkan ke booking.
        </template>
        <template v-else>
          Masih ada <strong class="text-gray-900">{{ position }}</strong> orang di depanmu.
        </template>
      </p>

      <p v-if="totalInQueue > 0 && position > 0" class="text-sm text-gray-400 mt-1">
        Total antrian: {{ totalInQueue }} orang
      </p>
    </div>

    <div v-if="position > 0" class="w-full bg-gray-200 rounded-full h-3 md:h-3.5 mb-4">
      <div
        class="bg-primary-500 h-3 md:h-3.5 rounded-full transition-all duration-700 ease-out"
        :style="{ width: `${progressPercent}%` }"
      />
    </div>

    <div v-if="position > 0" class="flex items-center justify-center gap-2 text-sm text-gray-500">
      <span class="relative flex h-2 w-2">
        <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary-400 opacity-75" />
        <span class="relative inline-flex rounded-full h-2 w-2 bg-primary-500" />
      </span>
      Memperbarui posisi...
    </div>

    <div v-if="error" class="mt-4 text-sm text-red-600 bg-red-50 rounded-lg p-3">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  position: number
  totalInQueue: number
  error?: string
}>()

const circumference = 2 * Math.PI * 52

const progressPercent = computed(() => {
  if (props.totalInQueue === 0) return 100
  const behind = props.totalInQueue - props.position
  return Math.max(0, Math.min(100, (behind / props.totalInQueue) * 100))
})

const dashOffset = computed(() => {
  return circumference - (progressPercent.value / 100) * circumference
})
</script>
