<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="favoritesCursor">
      <GamingTopicList :topics="items" :show-pinned="false" />
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { CursorResult } from '@/composables/api'

const props = defineProps({ user: { type: Object, required: true } })

const favoritesCursor = new CursorResult('/api/topic/users/favorites', { userId: props.user.id })
</script>

<style scoped>
/* No local styles */
</style>
