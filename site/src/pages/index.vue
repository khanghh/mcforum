<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6">
    <div class="flex gap-6">
      <GamingSidebar />
      <!-- Middle Content Area -->
      <div class="flex-1 min-w-0">
        <div class="gaming-card rounded-xl p-6 sm:p-8 mb-6 relative overflow-hidden">
          <div class="absolute inset-0 animated-border opacity-10"></div>
          <div class="relative z-10">
            <h1
              class="text-3xl sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 gaming-title py-3">
              {{ $t('feed.whats_new.title') }}
            </h1>
            <p class="text-gray-300 mb-4">
              {{ $t('feed.whats_new.description') }}
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
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()

const api = useApi()

definePageMeta({
  layout: 'default',
})

const newTopicsCursor = await api.getTopicFeeds(FeedType.WhatsNew).catch(err => {
  const status = err.statusCode || 500
  const message = err.message || t('page.server_error')
  throw createError({ statusCode: status, statusMessage: message, fatal: true })
})

useHead({
  title: useSiteTitle(t('feed.whats_new.title')),
  bodyAttrs: {
    class: 'bg-[#0f0f23]',
  },
})
</script>
