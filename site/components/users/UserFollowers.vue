<template>
  <LoadMoreAsync v-slot="{ items }" :cursor="followingCursor">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
      <UserCard v-for="user in items" :key="user.id" :user="user" />
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

const { user: loggedInUser } = storeToRefs(userStore)
const isSelf = loggedInUser.value?.username === props.user.username

const getUserFollowing = async () => {
  if (isSelf) {
    return await api.getMyFollowers()
  }
  return await api.getUserFollowers(props.user.username)
}

const followingCursor = await getUserFollowing()

</script>