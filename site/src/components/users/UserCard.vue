<template>
  <div class="gaming-card p-4 rounded-lg flex items-center gap-4">
    <Avatar :src="user.avatar" :username="user.username" :size="48" class="w-12 h-12 object-cover" />
    <div class="flex-1">
      <div class="flex items-center gap-2">
        <div class="font-bold text-sm text-white">{{ user.nickname || user.username }}</div>
        <span class="px-1 py-0.5 bg-blue-500 text-white text-[8px] font-bold rounded-full">LVL {{ level }}</span>
      </div>
      <div class="text-xs text-gray-400 mt-1 line-clamp-1">{{ user.statusMessage || '' }}</div>
    </div>
    <div>
      <button
        @click="toggleFollow"
        :class="isFollowing
          ? 'px-2 py-0.5 border-2 border-purple-500/50 text-purple-300 rounded-md text-xs font-bold hover:bg-purple-500/5 transition-colors'
          : 'px-2 py-0.5 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-md text-xs font-bold transition-transform'">
        {{ isFollowing ? 'UNFOLLOW' : 'FOLLOW' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import Avatar from '@/components/Avatar.vue'
import { useApi } from '@/composables/api'
import type { UserProfile } from '@/types'
const api = useApi()

type Props = {
  user: UserProfile
  isFollowing: boolean
}

const props = defineProps<Props>()
const user = props.user

const isFollowing = ref(props.isFollowing)
const level = computed(() => Math.floor(Math.sqrt(((user as any)?.score) || 0)) + 1)

async function toggleFollow() {
  if (isFollowing.value) {
    try {
      await api.unfollowUser(user.username)
      isFollowing.value = false
      // useMsgSuccess('Unfollowed')
    } catch (e) {
      console.error(e)
      // useMsgError('Unfollow failed')
    }
  } else {
    try {
      await api.followUser(user.username)
      isFollowing.value = true
      // useMsgSuccess('Followed')
    } catch (e) {
      console.error(e)
      // useMsgError('Follow failed')
    }
  }
}
</script>

<style scoped>
.gaming-card {
  background: rgba(17, 24, 39, 0.6);
}
</style>
