<template>
  <div class="min-h-screen flex items-center justify-center px-4">
    <div class="text-center max-w-md">
      <div class="mb-6">
        <svg v-if="is404" class="w-24 h-24 text-gray-200 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg v-else class="w-24 h-24 text-red-200 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
      </div>

      <h1 class="text-6xl font-bold text-gray-900 mb-2">{{ is404 ? '404' : '500' }}</h1>
      <p class="text-lg text-gray-500 mb-6">
        {{ is404 ? 'Halaman tidak ditemukan' : 'Terjadi kesalahan pada server' }}
      </p>
      <p class="text-sm text-gray-400 mb-8">
        {{ error?.statusMessage || (is404 ? 'URL yang kamu akses tidak tersedia.' : 'Silakan coba lagi nanti.') }}
      </p>

      <div class="flex items-center justify-center gap-3">
        <NuxtLink to="/" class="btn-primary">
          Kembali ke Beranda
        </NuxtLink>
        <button class="btn-outline" @click="handleClearError">
          Coba Lagi
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{ error?: { statusCode?: number; statusMessage?: string } }>()

const is404 = computed(() => props.error?.statusCode === 404)

function handleClearError() {
  clearError({ redirect: '/' })
}
</script>
