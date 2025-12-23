<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="activitiesCursor">
      <div v-for="item in items" :key="item.id" class="gaming-card rounded-xl p-4 border-l-4 border-blue-500/50">
        <div class="flex items-start gap-4">
          <FontAwesome :icon="['fas', 'comment-dots']" class="text-2xl text-blue-400 mt-1" />
          <div class="flex-1">
            <div class="text-gray-400 text-sm mb-1">Commented on <nuxt-link :to="`/topics/${item.entityId}`"
                class="text-blue-300 hover:underline">Topic #{{ item.entityId }}</nuxt-link></div>
            <p class="text-gray-300 leading-relaxed">{{ item.content }}</p>
            <div class="mt-2 text-xs text-gray-500">{{ usePrettyDate(item.createTime) }}</div>
          </div>
        </div>
      </div>
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { CursorResult } from '@/composables/api'

const props = defineProps({ user: { type: Object, required: true } })

const activitiesCursor = new CursorResult('/api/comment/list', { userId: props.user.id, entityType: 'topic' })
</script>

<style scoped>
/* No local styles */
</style>
