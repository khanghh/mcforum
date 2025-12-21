<template>
  <div class="relative">
    <!-- Cover Photo -->
    <div
      class="h-48 sm:h-64 bg-gradient-to-r from-purple-900 via-indigo-900 to-blue-900 rounded-2xl overflow-hidden relative">
      <img
        :src="user.backgroundImage || 'https://images.unsplash.com/photo-1511512578047-dfb367046420?auto=format&fit=crop&w=1600&q=80'"
        class="w-full h-full object-cover opacity-40" alt="Gaming Cover">
      <!-- Animated overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-purple-900/80 via-transparent to-transparent"></div>
      <div class="absolute inset-0 animated-border opacity-20"></div>
    </div>

    <div class="px-6 -mt-12 sm:-mt-16 relative z-10 flex items-end justify-between gap-4">
      <div class="flex items-end gap-4">
        <div class="relative">
          <!-- Pro Frame Outer Ring -->
          <div
            class="absolute -inset-2 bg-gradient-to-r from-amber-500 via-purple-500 to-pink-500 rounded-full opacity-75 blur-sm animate-pulse">
          </div>
          <!-- Pro Frame Inner -->
          <div class="relative animated-border p-1 rounded-full">
            <!-- Frame Decorations -->
            <div class="absolute -top-1 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10">
              <div
                class="bg-gradient-to-r from-amber-500 to-orange-500 px-2 py-0.5 rounded-full border-2 border-gray-900 flex items-center gap-1">
                <FontAwesome :icon="['fas', 'crown']" class="text-white text-xs" />
                <span class="text-white text-xs font-bold gaming-title">PRO</span>
              </div>
            </div>
            <!-- Corner Stars -->
            <div class="absolute top-0 right-0 text-amber-400 text-xl transform translate-x-1 -translate-y-1">
              <FontAwesome :icon="['fas', 'star']" class="drop-shadow-[0_0_10px_rgba(251,191,36,0.8)]" />
            </div>
            <div class="absolute bottom-0 left-0 text-purple-400 text-xl transform -translate-x-1 translate-y-1">
              <FontAwesome :icon="['fas', 'star']" class="drop-shadow-[0_0_10px_rgba(168,85,247,0.8)]" />
            </div>

            <my-avatar :user="user" :size="128" :round="true"
              class="w-24 h-24 sm:w-32 sm:h-32 rounded-full border-4 border-gray-900 shadow-2xl relative z-10 block" />

            <div
              class="absolute bottom-0 right-0 bg-green-500 w-6 h-6 rounded-full border-3 border-gray-900 flex items-center justify-center z-20">
              <span class="w-2 h-2 bg-white rounded-full animate-pulse"></span>
            </div>
          </div>
        </div>
        <div class="pb-2">
          <div class="flex items-center gap-2 flex-wrap">
            <h1
              class="text-2xl sm:text-3xl font-bold text-white gaming-title drop-shadow-[0_0_15px_rgba(139,92,246,0.8)]">
              {{ user.nickname }}
            </h1>
            <span v-if="user.id === 1 || user.type === 1"
              class="px-3 py-1 bg-gradient-to-r from-red-600 to-red-500 text-white text-xs font-bold rounded-full shadow-lg flex items-center gap-1">
              <FontAwesome :icon="['fas', 'shield-alt']" /> ADMIN
            </span>
            <span
              class="px-3 py-1 bg-gradient-to-r from-amber-500 to-orange-500 text-white text-xs font-bold rounded-full shadow-lg">
              LVL
              {{ calculateLevel(user.score) }}
            </span>
          </div>
          <div class="flex items-center gap-2 mt-1">
            <div v-if="user.description"
              class="gaming-card px-3 py-1 rounded-lg backdrop-blur-sm inline-flex items-center">
              <FontAwesome :icon="['fas', 'quote-left']" class="text-purple-400 text-sm mr-2" />
              <span class="text-sm text-purple-200 font-medium line-clamp-1">{{ user.description }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="hidden md:flex items-center gap-2 mb-3">
        <button v-if="!isSelf"
          class="px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center neon-border gaming-title text-sm hover:scale-105 transition-transform"
          aria-label="Follow user">
          <FontAwesome :icon="['fas', 'user-plus']" class="mr-2" /> FOLLOW
        </button>
        <button v-if="!isSelf"
          class="px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors"
          aria-label="Send message">
          <FontAwesome :icon="['fas', 'comments']" class="mr-2" /> MESSAGE
        </button>
        <nuxt-link v-if="isSelf" to="/user/profile"
          class="px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors">
          <FontAwesome :icon="['fas', 'cog']" class="mr-2" /> EDIT PROFILE
        </nuxt-link>
      </div>
    </div>

    <!-- Mobile Settings Button -->
    <div v-if="isSelf" class="absolute top-4 right-4 md:hidden">
      <nuxt-link to="/user/profile"
        class="gaming-card text-purple-300 px-4 py-2 rounded-lg text-sm font-medium backdrop-blur-sm block">
        <FontAwesome :icon="['fas', 'cog']" />
      </nuxt-link>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

const userStore = useUserStore()
const currentUser = computed(() => userStore.user)
const isSelf = computed(() => currentUser.value && currentUser.value.id === props.user.id)

function calculateLevel(score) {
  return Math.floor(Math.sqrt(score || 0)) + 1
}
</script>

<style scoped>
/* Animated gradient border */
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

.animated-border {
  background: linear-gradient(45deg, #8b5cf6, #ec4899, #06b6d4, #8b5cf6);
  background-size: 300% 300%;
  animation: gradientBorder 3s ease infinite;
}

.neon-border {
  box-shadow: 0 0 10px rgba(139, 92, 246, 0.5), 0 0 20px rgba(139, 92, 246, 0.3);
}
</style>
