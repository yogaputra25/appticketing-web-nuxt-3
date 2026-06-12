<template>
  <span class="tabular-nums" :class="{ 'text-red-600': isUrgent }">
    {{ display }}
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{
  targetDate: string
}>()

const emit = defineEmits<{
  expired: []
}>()

const display = ref('')
const isUrgent = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

function update() {
  const now = Date.now()
  const target = new Date(props.targetDate).getTime()
  const diff = target - now

  if (diff <= 0) {
    display.value = '00:00:00'
    isUrgent.value = true
    if (timer) clearInterval(timer)
    emit('expired')
    return
  }

  const hours = Math.floor(diff / 3600000)
  const minutes = Math.floor((diff % 3600000) / 60000)
  const seconds = Math.floor((diff % 60000) / 1000)

  display.value = `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`
  isUrgent.value = diff < 60000
}

onMounted(() => {
  update()
  timer = setInterval(update, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
