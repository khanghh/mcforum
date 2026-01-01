<template>
  <div class="gaming-card rounded-xl p-4">
    <h3 class="text-lg font-bold text-amber-400 gaming-title mb-4 flex items-center">
      <Icon name="TablerMedal2" class="mr-2" /> {{ $t('widgets.top_contributors') }}
    </h3>
    <div class="space-y-3">
      <div v-for="(user, idx) in contributors" :key="user.id" class="flex items-center gap-3">
        <span class="text-2xl">{{ getMedalEmoji(idx) }}</span>
        <Avatar :src="user.avatar" :username="user.username" size="40"
          class="rounded border border-purple-500/50 flex-shrink-0" />
        <nuxt-link :to="`/users/${user.username}`" class="flex-1">
          <p :class="['font-bold text-sm hover:text-purple-400 transition-colors', getColorClass(idx)]">{{ user.username
          }}</p>
          <p class="text-xs text-gray-400">{{ user.score }} points</p>
        </nuxt-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type UserInfo } from '@/types'
const api = useApi()

const contributors = ref<UserInfo[]>()

contributors.value = await api.getTopContributors().catch(() => [])

function getColorClass(index: number): string {
  switch (index) {
    case 0:
      return 'text-amber-300 gaming-title'
    case 1:
      return 'text-orange-300 gaming-title'
    default:
      return 'text-red-300 gaming-title'
  }
}

function getMedalEmoji(index: number): string {
  switch (index) {
    case 0:
      return '🥇'
    case 1:
      return '🥈'
    case 2:
      return '🥉'
    default:
      return '🏅'
  }
}

</script>