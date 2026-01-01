<template>
  <div>
    <h4 class="text-lg font-bold text-amber-400 gaming-title mb-4">FORUM STATS</h4>
    <div class="gaming-card rounded-xl p-4 gaming-title">
      <div class="space-y-2 text-sm">
        <div class="flex justify-between">
          <span class="text-gray-400">{{ $t('stats.total_topics') }}</span>
          <span class="font-bold text-gray-300">{{ prettyNumber(stats.topics) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">{{ $t('stats.total_posts') }}</span>
          <span class="font-bold text-gray-300">{{ prettyNumber(stats.posts) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">{{ $t('stats.total_members') }}</span>
          <span class="font-bold text-gray-300">{{ prettyNumber(stats.members) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">{{ $t('stats.total_visits') }}</span>
          <span class="font-bold text-gray-300">{{ prettyNumber(stats.visits) }}</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-400">{{ $t('stats.newest_member') }}</span>
          <nuxt-link :to="`/users/${stats.newestMember}`"
            class="font-bold text-blue-300 hover:text-purple-400 transition-colors">
            {{ stats.newestMember }}
          </nuxt-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const route = useRoute()
const api = useApi()

const { data: stats, refresh } = await useAsyncData(
  'forum-stats',
  () => api.getForumStats(),
  {
    default: () => ({
      topics: 0,
      posts: 0,
      members: 0,
      visits: 0,
      newestMember: '',
    }),
    watch: [() => route.fullPath],
  }
)

const prettyNumber = (num: number): string => {
  if (num >= 1_000_000) {
    return (num / 1_000_000).toFixed(1) + 'M'
  } else if (num >= 1_000) {
    return (num / 1_000).toFixed(1) + 'K'
  } else {
    return num.toString()
  }
}

</script>