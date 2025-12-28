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
                  {{ feedInfo?.title.toUpperCase() }}
                </h1>
                <p class="text-gray-300 mb-4">
                  {{ feedInfo?.description }}
                </p>
              </div>
            </div>

            <!-- Threads List -->
            <div id="threads-container" class="space-y-4">
              <LoadMoreAsync v-slot="{ items }" :cursor="feedCursor">
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
import { useApi, FeedType } from '@/composables/api'
import type { CursorResult } from '@/composables/api'
import type { Topic } from '~/types'

const route = useRoute()
const { t } = useI18n()
const api = useApi()


definePageMeta({
  layout: false,
})

interface FeedInfo {
  title: string
  description: string
}

const feedType: FeedType = route.params.feed as FeedType
const allowFeeds = ['recommended', 'followed', 'whats-new']
if (!feedType || !allowFeeds.includes(feedType)) {
  throw createError({ statusCode: 404, statusMessage: 'Page Not Found', fatal: true })
}

const feedInfo = computed<FeedInfo>(() => {
  switch (feedType) {
    case FeedType.Recommended:
      return {
        title: t('feed.recommended.title'),
        description: t('feed.recommended.description'),
      }
    case FeedType.Followed:
      return {
        title: t('feed.followed.title'),
        description: t('feed.followed.description'),
      }
    case FeedType.WhatsNew:
    default:
      return {
        title: t('feed.whats_new.title'),
        description: t('feed.whats_new.description'),
      }
  }
})


const feedCursor: CursorResult<Topic[]> = api.getTopicFeeds(feedType)

useHead({
  title: useSiteTitle(feedInfo.value.title),
  bodyAttrs: {
    class: 'bg-[#0f0f23]',
  },
})
</script>
