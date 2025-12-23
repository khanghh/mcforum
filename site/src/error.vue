<template>
  <div
    class="min-h-screen bg-gray-900 text-gray-100 font-sans selection:bg-purple-500 selection:text-white flex flex-col"
    style="background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);">
    <GamingNavbar />

    <main class="flex-grow flex items-center justify-center py-12 px-4 sm:px-6">
      <div class="max-w-md w-full text-center">
        <div class="mb-8 flex flex-col sm:flex-row items-center justify-center gap-4">
          <div class="relative inline-block">
            <div class="absolute -inset-4 bg-purple-500/20 rounded-full blur-xl animate-pulse"></div>
            <img src="/images/logo.png"
              class="relative z-10 w-20 h-20 object-contain drop-shadow-[0_0_15px_rgba(168,85,247,0.5)]"
              alt="Logo" />
          </div>
          <h1
            class="text-3xl font-bold tracking-wider text-white drop-shadow-[0_0_10px_rgba(168,85,247,0.8)]">
            MineViet Network
          </h1>
        </div>

        <!-- Error Card -->
        <div
          class="bg-[linear-gradient(145deg,rgba(30,30,60,0.6),rgba(20,20,40,0.8))] rounded-2xl p-8 border border-purple-500/20 relative overflow-hidden backdrop-blur-sm">
          <!-- Animated border overlay -->
          <div
            class="absolute inset-0 bg-[linear-gradient(45deg,#8b5cf6,#ec4899,#06b6d4,#8b5cf6)] bg-[length:300%_300%] animate-[gradientBorder_3s_ease_infinite] opacity-10">
          </div>

          <div class="relative z-10">
            <h1
              class="text-6xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 mb-2">
              {{ error ? error.statusCode : 'Error' }}
            </h1>

            <div class="text-xl text-purple-200 font-medium mb-6">
              <template v-if="error">
                <span v-if="error.statusCode === 404">PAGE NOT FOUND</span>
                <span v-else-if="error.statusCode === 403">ACCESS DENIED</span>
                <span v-else>INTERNAL SERVER ERROR</span>
              </template>
              <template v-else>
                Unknown Error
              </template>
            </div>

            <p class="text-gray-400 mb-8 leading-relaxed">
              <span v-if="error && error.statusCode === 404">
                {{ $t('page.errors.page_not_found') }}
              </span>
              <span v-else-if="error && error.statusCode === 403">
                {{ $t('page.errors.no_permission') }}
              </span>
              <span v-else>
                {{ $t('page.errors.internal_server_error') }}
              </span>
            </p>

            <button
              class="w-full sm:w-auto px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center justify-center mx-auto shadow-[0_0_10px_rgba(139,92,246,0.5),0_0_20px_rgba(139,92,246,0.3)] font-['Saira Semi Condensed',_sans-serif] tracking-[0.5px] transition-transform hover:scale-105"
              @click="handleError">
              <FontAwesome :icon="['fas', 'home']" class="mr-2" />
              {{ $t('links.return_home') }}
            </button>
          </div>
        </div>
      </div>
    </main>

    <GamingFooter />
  </div>
</template>

<script setup lang="ts">

defineProps({
  error: {
    type: Object,
    default: null,
  },
})

definePageMeta({
  layout: false,
})

const handleError = () => {
  clearError({ redirect: '/' })
}
</script>

<style scoped>
/* Animated gradient border (reused from gaming themes) */
@keyframes gradientBorder {
  0% {
    background-position: 0% 50%;
  }

  50% {
    background-position: 100% 50%;
  }

  100% {
    background-position: 0% 50%;
  }
}
</style>
