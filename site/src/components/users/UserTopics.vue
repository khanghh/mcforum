<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="topicsCursor">
      <GamingTopicList :topics="items" :show-pinned="false" />
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import type { UserProfile, Topic } from '@/types'
import { useApi } from '@/composables/api'

const api = useApi()
const userStore = useUserStore()

interface Props {
  user: UserProfile
}

const props = defineProps<Props>()
const { user: loggedInUser } = storeToRefs(userStore)
const isSelf = loggedInUser.value?.username === props.user.username

const fetchUserTopics = (): CursorResult<Topic[]> => {
  if (isSelf) {
    return api.getMyTopics()
  }
  return api.getUserTopics(props.user.username)
}

const topicsCursor = fetchUserTopics()
</script>
