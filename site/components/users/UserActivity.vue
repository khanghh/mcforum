<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="activityCursor">
      <div v-for="item in items" :key="item.id" class="gaming-card rounded-xl p-4">
        <div class="flex items-start gap-4">
          <div class="flex-shrink-0 mt-1 text-2xl text-purple-400">
            <Icon :name="iconFor(item)" />
          </div>
          <div class="flex-1">
            <div class="text-sm text-gray-400 mb-1">
              <span class="font-medium text-gray-200">{{ titleFor(item) }}</span>
              <span class="ml-2 text-xs text-gray-500">{{ usePrettyDate(item.createTime) }}</span>
            </div>
            <div class="text-gray-300 leading-relaxed">
              <template v-if="item.content">{{ item.content }}</template>
              <template v-else-if="item.data && item.data.content">{{ item.data.content }}</template>
              <template v-else>{{ JSON.stringify(item) }}</template>
            </div>
            <div v-if="item.entityId" class="mt-2 text-xs">
              <nuxt-link :to="`/topics/${item.entityId}`" class="text-blue-300 hover:underline">View related
                topic</nuxt-link>
            </div>
          </div>
        </div>
      </div>
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { CursorResult } from '@/composables/api'

const props = defineProps({ user: { type: Object, required: true } })

const activityCursor = new CursorResult('/api/activity/list', { userId: props.user.id })

function iconFor(item) {
  const t = (item.type || item.action || item.kind || '').toLowerCase()
  if (t.includes('comment')) return 'PhChatCircle'
  if (t.includes('post') || t.includes('create')) return 'PhFire'
  if (t.includes('like') || t.includes('favorite')) return 'PhHeart'
  return 'PhLightning'
}

function titleFor(item) {
  if (item.type) return item.type
  if (item.action) return item.action
  if (item.kind) return item.kind
  return 'Activity'
}
</script>

<style scoped>
/* No local styles */
</style>
