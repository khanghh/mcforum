<template>
  <div class="text-gray-100 flex flex-col min-h-screen custom-scrollbar font-sans">
    <GamingNavbar />

    <main class="flex-grow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6">
        <div class="flex gap-6">
          <GamingSidebar />

          <!-- Middle Content Area -->
          <div class="flex-1 min-w-0">
            <!-- Hero Banner -->
            <div class="gaming-card rounded-xl p-6 sm:p-8 mb-6 relative overflow-hidden">
              <div class="absolute inset-0 animated-border opacity-10"></div>
              <div class="relative z-10">
                <h1
                  class="text-3xl sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 gaming-title py-3">
                  {{ feed.title }}
                </h1>
                <p class="text-gray-300 mb-4">
                  {{ feed.description }}
                </p>
              </div>
            </div>

            <!-- Threads List -->
            <div id="threads-container" class="space-y-4">
              <LoadMoreAsync v-slot="{ items }" :url="feed.fetchTopicsUrl">
                <GamingTopicList :topics="items" show-pinned />
              </LoadMoreAsync>
            </div>
          </div>
        </div>
      </div>
    </main>

    <GamingFooter />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const { t } = useI18n()

definePageMeta({
  layout: false,
})

type FeedInfo = {
  title: string
  fetchTopicsUrl: string
  description: string
}

const feed = ref < FeedInfo > ({
  title: t('feed.whats_new.title'),
  fetchTopicsUrl: '/api/feeds/whats-new',
  description: t('feed.whats_new.description'),
})

useHead({
  title: useSiteTitle(feed.value.title),
  bodyAttrs: {
    class: 'bg-[#0f0f23]',
  },
})
</script>

<style lang="scss">
@import "~/assets/css/gaming-design.scss";
</style>
