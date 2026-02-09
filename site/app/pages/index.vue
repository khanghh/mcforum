<template>
  <div class="flex gap-6">
    <GamingSidebar />
    <!-- Middle Content Area -->
    <div class="flex-1 min-w-0">
      <div class="gaming-card sm:rounded-xl p-6 sm:p-8 mb-6 relative overflow-hidden">
        <div class="absolute inset-0 animated-border opacity-10"></div>
        <div class="relative z-10">
          <h1
            class="text-3xl sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 gaming-title py-3">
            {{ $t('feed.whats_new') }}
          </h1>
          <p class="text-gray-300 mb-4">
            {{ $t('feed.whats_new_desc') }}
          </p>
        </div>
      </div>

      <!-- Threads List -->
      <div id="threads-container" class="space-y-4">
        <LoadMoreAsync v-slot="{ items }" :cursor="newTopicsCursor">
          <GamingTopicList :topics="items" show-pinned />
        </LoadMoreAsync>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const i18n = useI18n()

const api = useApi()

definePageMeta({
  layout: 'default',
})

const newTopicsCursor = api.getTopicFeeds(FeedType.WhatsNew)

</script>
