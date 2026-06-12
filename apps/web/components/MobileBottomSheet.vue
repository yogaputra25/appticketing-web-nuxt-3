<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50"
        @click.self="close"
      >
        <div class="absolute inset-0 bg-black/50" />

        <div
          ref="panelRef"
          role="dialog"
          aria-modal="true"
          class="absolute bottom-0 left-0 right-0 md:top-1/2 md:left-1/2 md:-translate-x-1/2 md:-translate-y-1/2 md:bottom-auto md:max-w-lg md:w-full"
        >
          <div
            class="bg-white rounded-t-2xl md:rounded-2xl shadow-xl max-h-[85vh] overflow-y-auto"
            @click.stop
          >
            <div class="sticky top-0 bg-white pt-3 pb-2 px-4 flex justify-center md:hidden">
              <div class="w-10 h-1.5 bg-gray-300 rounded-full" />
            </div>

            <div class="px-4 pb-4 pt-2 md:pt-4">
              <slot />
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const panelRef = ref<HTMLElement | null>(null)

function close() {
  emit('update:modelValue', false)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.modelValue) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
})

watch(() => props.modelValue, (val) => {
  if (val) {
    nextTick(() => {
      panelRef.value?.focus()
    })
  }
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
