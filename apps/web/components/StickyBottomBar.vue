<template>
  <div
    :class="[
      visible ? 'md:hidden' : 'hidden',
      'fixed bottom-0 left-0 right-0 z-40 bg-white border-t border-gray-200 shadow-[0_-2px_10px_rgba(0,0,0,0.08)] px-4 py-3'
    ]"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  visible?: boolean
}>(), {
  visible: true,
})
</script>
