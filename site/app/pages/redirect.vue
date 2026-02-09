<template>
  <div class="min-h-[80vh] flex items-center justify-center p-4 relative overflow-hidden">
    <!-- Ambient Background Effects -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-600/20 rounded-full blur-[100px] animate-pulse"></div>
      <div
        class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-pink-600/20 rounded-full blur-[100px] animate-pulse delay-1000">
      </div>
    </div>

    <div class="relative w-full max-w-3xl z-10">
      <div class="flex flex-col items-center text-center">
        <!-- Logo -->
        <div class="mb-8 relative">
          <div class="absolute -inset-5 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full opacity-20 blur-xl">
          </div>
          <img
            src="/images/logo.png"
            alt="Logo"
            class="relative w-24 h-24 object-contain drop-shadow-2xl" />
        </div>

        <!-- Text -->
        <h2 class="text-xl sm:text-3xl font-bold mb-3 text-white">
          {{ $t('message.redirecting_message') }}
        </h2>

        <p class="text-gray-400 font-medium mb-3">
          <i18n-t keypath="message.external_redirect" tag="span">
            <a
              :href="url"
              class="text-purple-300 hover:text-purple-200 underline underline-offset-4"
              rel="noopener noreferrer">
              {{ $t('message.click_to_redirect') }}
            </a>
          </i18n-t>
        </p>

        <!-- Loader -->
        <div class="mt-10 w-full max-w-[260px]">
          <div class="h-1.5 w-full bg-gray-700/50 rounded-full overflow-hidden relative">
            <div class="absolute inset-0 bg-gradient-to-r from-purple-500 to-pink-500 animate-progress w-1/2"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const i18n = useI18n()
const route = useRoute()

definePageMeta({
  layout: 'default',
})

const autoRedirectTime = 5000

const url = computed(() => {
  const value = route.query.url
  if (Array.isArray(value)) return value[0] ?? ''
  return (value as string | undefined) ?? ''
})

const normalized = computed(() => url.value.trim().toLowerCase())

if (!normalized.value.startsWith('http://') && !normalized.value.startsWith('https://')) {
  throw createError({
    statusCode: 500,
    message: i18n.t('message.internal_server_error'),
    fatal: true,
  })
}

onMounted(() => {
  setTimeout(() => {
    window.location.href = url.value
  }, autoRedirectTime)
})
</script>

<style scoped>
@keyframes progress {
  0% {
    left: -50%;
  }

  100% {
    left: 100%;
  }
}

.animate-progress {
  animation: progress 1.5s ease-in-out infinite;
}
</style>
