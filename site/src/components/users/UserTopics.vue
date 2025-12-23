<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :url="fetchTopicsUrl">
      <GamingTopicList :topics="items" :show-pinned="false" />
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserProfile } from '@/types'

const userStore = useUserStore()

interface UserTopicsProps {
  user: UserProfile
}

const props = defineProps<UserTopicsProps>()
const loggedInUser = computed(() => userStore.user)
const isSelf = loggedInUser.value?.username === props.user.username

const fetchTopicsUrl = computed(() => {
  const user = props.user
  if (isSelf) return '/api/users/me/topics'
  return `/api/users/${user.id}/topics`
})

</script>
