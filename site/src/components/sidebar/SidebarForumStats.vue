<template>
  <div class="gaming-card rounded-xl p-4 gaming-title">
    <h3 class="text-lg font-bold text-gray-400 mb-4">FORUM STATS</h3>
    <div class="space-y-2 text-sm">
      <div class="flex justify-between">
        <span class="text-gray-400">Threads</span>
        <span class="font-bold text-purple-300">{{ usePrettyNumber(stats.threads) }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-gray-400">Posts</span>
        <span class="font-bold text-pink-300">{{ usePrettyNumber(stats.posts) }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-gray-400">Members</span>
        <span class="font-bold text-blue-300">{{ usePrettyNumber(stats.members) }}</span>
      </div>
      <div class="flex justify-between">
        <span class="text-gray-400">Online</span>
        <span class="font-bold text-green-300">{{ usePrettyNumber(stats.online) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const stats = ref({
  threads: 0,
  posts: 0,
  members: 0,
  online: 0,
})

</script>