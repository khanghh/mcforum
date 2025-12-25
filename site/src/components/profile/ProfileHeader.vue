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
    </div>

    <div class="px-6 -mt-12 sm:-mt-16 relative z-10 flex items-end justify-between gap-4">
      <div class="flex items-end gap-4">
        <div class="relative">
          <!-- Pro Frame Inner -->
          <div class="relative p-1">
            <AvatarDecorated :user="user" :is-self="isSelf" @uploaded="onUploaded" />
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
              <Icon name="shield-alt" /> ADMIN
            </span>
          </div>
          <div class="flex items-center gap-2 mt-1">
            <div
              :class="['gaming-card px-3 py-1 rounded-lg backdrop-blur-sm inline-flex items-center', { invisible: !user.statusMessage }]">
              <Icon name="TablerQuoteFilled" class="text-purple-400 text-sm mr-2" />
              <span class="text-sm text-purple-200 font-medium line-clamp-1">{{ user.statusMessage || ' ' }}</span>
            </div>
          </div>
        </div>
      </div>
      <div v-if="currentUser" class="hidden md:flex items-center gap-2 mb-3">
        <button v-if="!isSelf"
          :class="isFollowing
            ? 'px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors'
            : 'px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center neon-border gaming-title text-sm hover:scale-105 transition-transform'"
          :aria-label="isFollowing ? 'Unfollow user' : 'Follow user'"
          @click="isFollowing ? handleUnfollow() : handleFollow()">
          <Icon :name="isFollowing ? 'user-minus' : 'user-plus'" class="mr-2" />
          {{ isFollowing ? 'UNFOLLOW' : 'FOLLOW' }}
        </button>
        <button v-if="!isSelf"
          class="px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors"
          aria-label="Send message"
          @click="handleMessage">
          <Icon name="TablerMessageCircles" class="mr-2" /> MESSAGE
        </button>
        <nuxt-link v-if="isSelf" to="/users/me/profile"
          class="px-4 py-2 border-2 border-purple-500/50 text-purple-300 rounded-lg font-bold flex items-center gaming-title text-sm hover:bg-purple-500/10 transition-colors">
          <Icon name="Fa7SolidGear" class="mr-2" /> EDIT PROFILE
        </nuxt-link>
      </div>
    </div>

    <!-- Mobile Settings Button -->
    <div v-if="isSelf" class="absolute top-4 right-4 md:hidden">
      <nuxt-link to="/users/me/profile"
        class="gaming-card text-purple-300 px-4 py-2 rounded-lg text-sm font-medium backdrop-blur-sm block">
        <Icon name="Fa7SolidGear" />
      </nuxt-link>
    </div>
  </div>
</template>

<script setup>
import AvatarDecorated from '../AvatarDecorated.vue'

const userStore = useUserStore()
const api = useApi()

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

const isFollowing = ref(false)

const currentUser = computed(() => userStore.user)
const isSelf = computed(() => currentUser?.value && currentUser?.value?.id === props.user.id)

// Resolve before render (SSR-friendly). If not logged in / any error -> false.
if (typeof userStore.getCurrent === 'function') {
  await userStore.getCurrent().catch(() => undefined)
}

if (!isSelf.value) {
  isFollowing.value = await api.isFollowing(props.user.id).catch(() => false)
}

const handleFollow = async () => {
  try {
    await api.followUser(props.user.username)
    isFollowing.value = true
    useMsgSuccess('Followed successfully')
  }
  catch (e) {
    console.error(e)
    useMsgError('Follow failed')
  }
}

const handleUnfollow = async () => {
  try {
    await api.unfollowUser(props.user.username)
    isFollowing.value = false
    useMsgSuccess('Unfollowed successfully')
  }
  catch (e) {
    console.error(e)
    useMsgError('Unfollow failed')
  }
}

const handleMessage = () => {
  // Navigate to messages page or open chat
  navigateTo(`/messages/${props.user.id}`)
}

async function onUploaded(file) {
  try {
    const formData = new FormData()
    formData.append('image', file, file.name)
    const ret = await useHttpPostMultipart('/api/upload', formData)

    await useHttpPostForm('/api/users/update/avatar', {
      body: { avatar: ret.url },
    })

    await userStore.getCurrent()
    useMsgSuccess('头像更新成功')
  }
  catch (e) {
    console.error(e)
    useMsgError('头像更新失败')
  }
}

function calculateLevel(score) {
  return Math.floor(Math.sqrt(score || 0)) + 1
}
</script>

<style scoped>
.neon-border {
  box-shadow: 0 0 10px rgba(139, 92, 246, 0.5), 0 0 20px rgba(139, 92, 246, 0.3);
}
</style>
