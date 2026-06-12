<template>
  <div class="flex gap-2 overflow-x-auto no-scrollbar md:justify-center md:flex-wrap" role="tablist">
    <button
      v-for="tab in tabs"
      :key="tab.value"
      role="tab"
      :aria-selected="modelValue === tab.value"
      class="touch-target whitespace-nowrap px-4 py-2 rounded-lg text-sm font-medium transition-colors shrink-0"
      :class="modelValue === tab.value
        ? 'bg-primary-600 text-white shadow-sm'
        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'"
      @click="emit('update:modelValue', tab.value)"
    >
      {{ tab.label }}
    </button>
  </div>
</template>

<script setup lang="ts">
interface Tab {
  label: string
  value: string
}

defineProps<{
  tabs: Tab[]
  modelValue: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
</script>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
