export function useViewport() {
  const isMobile = ref(false)
  const isTablet = ref(false)
  const isDesktop = ref(false)

  function update() {
    if (import.meta.client) {
      isMobile.value = window.matchMedia('(max-width: 639px)').matches
      isTablet.value = window.matchMedia('(min-width: 640px) and (max-width: 1023px)').matches
      isDesktop.value = window.matchMedia('(min-width: 1024px)').matches
    }
  }

  onMounted(() => {
    update()
    const mqlMobile = window.matchMedia('(max-width: 639px)')
    const mqlTablet = window.matchMedia('(min-width: 640px) and (max-width: 1023px)')
    const mqlDesktop = window.matchMedia('(min-width: 1024px)')

    mqlMobile.addEventListener('change', update)
    mqlTablet.addEventListener('change', update)
    mqlDesktop.addEventListener('change', update)

    onUnmounted(() => {
      mqlMobile.removeEventListener('change', update)
      mqlTablet.removeEventListener('change', update)
      mqlDesktop.removeEventListener('change', update)
    })
  })

  return { isMobile, isTablet, isDesktop }
}
