<template>
  <div class="min-h-[80vh] flex items-center justify-center p-4 relative overflow-hidden">
    <!-- Ambient Background Effects -->
    <div class="absolute inset-0 pointer-events-none">
      <div class="absolute top-1/4 left-1/4 w-96 h-96 bg-purple-600/20 rounded-full blur-[100px] animate-pulse"></div>
      <div
        class="absolute bottom-1/4 right-1/4 w-96 h-96 bg-pink-600/20 rounded-full blur-[100px] animate-pulse delay-1000">
      </div>
    </div>

    <div class="relative w-full max-w-md z-10">
      <!-- Card -->
      <div
        class="relative overflow-hidden backdrop-blur-xl bg-[#1a1b26]/80 border border-white/10 rounded-2xl shadow-2xl p-8 sm:p-10">
        <!-- Gradient Top Border -->
        <div
          class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-transparent via-purple-500 to-transparent opacity-50">
        </div>

        <!-- Content -->
        <div class="flex flex-col items-center text-center">
          <!-- Logo Area -->
          <div class="mb-8 relative group">
            <div
              class="absolute -inset-4 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full opacity-20 group-hover:opacity-40 blur-xl transition-opacity duration-500">
            </div>
            <img src="/images/logo.png" alt="MineViet"
              class="relative w-24 h-24 object-contain drop-shadow-2xl transform group-hover:scale-105 transition-transform duration-500" />
          </div>

          <!-- Status Text -->
          <h2 class="text-2xl font-bold mb-3 text-white">
            <span v-if="!errorMsg" class="inline-flex items-center gap-2">
              {{ $t('page.authenticating') }}
            </span>
            <span v-else class="text-red-400">{{ $t('message.login_failed') }}</span>
          </h2>
          <div v-if="errorMsg">
            <p class="text-gray-400 mb-3 font-medium">
              {{ $t('message.failed_to_authenticate') }}
            </p>
            <p class="text-gray-400 mb-8 font-medium">
              {{ $t('message.failed_reason') }} {{ errorMsg }}
            </p>
          </div>

          <p v-else class="text-gray-400 mb-8 font-medium">
            {{ $t('message.waiting_for_verification') }}
          </p>

          <!-- Loader or Action -->
          <div v-if="!errorMsg" class="w-full max-w-[200px]">
            <div class="h-1.5 w-full bg-gray-700/50 rounded-full overflow-hidden relative">
              <div class="absolute inset-0 bg-gradient-to-r from-purple-500 to-pink-500 animate-progress w-1/2"></div>
            </div>
          </div>

          <div v-else>
            <button
              class="group relative inline-flex items-center justify-center px-8 py-3 font-bold text-white transition-all duration-200 bg-gradient-to-r from-purple-600 to-pink-600 rounded-lg hover:from-purple-500 hover:to-pink-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-600 focus:ring-offset-[#1a1b26] shadow-lg shadow-purple-900/50"
              @click="navigateTo('/')">
              <Icon name="Fa7SolidArrowLeft" class="mr-1 group-hover:-translate-x-1 transition-transform" />
              {{ $t('button.back_to_home') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const errorMsg = ref('')

definePageMeta({
  layout: 'default',
})

onMounted(async () => {
  const ticket = route.query.ticket as string
  const state = route.query.state as string
  if (!ticket) {
    errorMsg.value = 'Missing login ticket'
    return
  }
  try {
    await userStore.signin(ticket, state)
    router.replace('/')
  } catch (e: any) {
    errorMsg.value = e.data?.error?.message || e.message || 'Unknown error'
    console.error('Authentication error:', e)
  }
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
