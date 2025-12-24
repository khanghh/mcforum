<template>
  <LoadMoreAsync v-slot="{ items }" :cursor="followingCursor">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
      <UserCard v-for="u in items" :key="u.id" :user="u" :isFollowing="true" />
    </div>
  </LoadMoreAsync>

</template>

<script setup lang="ts">
import type { UserProfile } from '@/types/user'
import UserCard from './UserCard.vue'
import { useApi } from '@/composables/api'
const api = useApi()
const userStore = useUserStore()

interface Props {
  user: UserProfile
}
const props = defineProps<Props>()

const loggedInUser = userStore.user
const isSelf = loggedInUser?.username === props.user.username

const getUserFollowing = async () => {
  if (isSelf) {
    return await api.getMyFollowers()
  }
  return await api.getUserFollowers(props.user.username)
}

const followingCursor = await getUserFollowing().catch(err => {
  const status = err.statusCode || 500
  const message = err.message || 'Server error'
  throw createError({ statusCode: status, statusMessage: message, fatal: true })
})

</script>