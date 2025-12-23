<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="topicsCursor">
      <GamingTopicList :topics="items" :show-pinned="false" />
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserProfile } from '@/types'
import { useApi } from '@/composables/api'

const api = useApi()
const userStore = useUserStore()
const i18n = useI18n()

interface UserTopicsProps {
  user: UserProfile
}

const props = defineProps<UserTopicsProps>()
const loggedInUser = computed(() => userStore.user)
const isSelf = loggedInUser.value?.username === props.user.username

const fetchUserTopics = async () => {
  if (isSelf) {
    return await api.getMyTopics()
  }
  return await api.getUserTopics(props.user.username)
}

const topicsCursor = await fetchUserTopics().catch(err => {
  const status = err.statusCode || 500
  const message = err.message || i18n.t('page.server_error')
  throw createError({ statusCode: status, statusMessage: message, fatal: true })
})
console.log('topicsCursor', topicsCursor)

</script>
