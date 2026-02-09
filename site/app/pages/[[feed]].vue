<template>
  <div class="flex gap-6">
    <GamingSidebar />

    <!-- Middle Content Area -->
    <div class="flex-1 min-w-0">
      <!-- Hero Banner -->
      <div class="gaming-card sm:rounded-xl p-6 sm:p-8 mb-6 relative overflow-hidden">
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
</template>

<script setup lang="ts">
import { useApi, FeedType } from '@/composables/api'
import type { CursorResult } from '@/composables/api'
import type { Topic } from '~/types'

const route = useRoute()
const userStore = useUserStore()
const i18n = useI18n()
const api = useApi()

definePageMeta({
  layout: 'default',
})

interface FeedInfo {
  title: string
  description: string
}


const feedType: FeedType = route.params.feed as FeedType

const { user } = storeToRefs(userStore)
if (feedType === FeedType.Followed && !user.value) {
  navigateTo('/')
}

const allowFeeds = [FeedType.WhatsNew, FeedType.Recommended, FeedType.Followed]
if (!feedType || !allowFeeds.includes(feedType)) {
  throw createError({ statusCode: 404, statusMessage: i18n.t('page.not_found'), fatal: true })
}

const feedInfo = computed<FeedInfo>(() => {
  switch (feedType) {
    case FeedType.Recommended:
      return {
        title: i18n.t('feed.recommended'),
        description: i18n.t('feed.recommended_desc'),
      }
    case FeedType.Followed:
      return {
        title: i18n.t('feed.followed'),
        description: i18n.t('feed.followed_desc'),
      }
    case FeedType.WhatsNew:
    default:
      return {
        title: i18n.t('feed.whats_new'),
        description: i18n.t('feed.whats_new_desc'),
      }
  }
})

const feedCursor: CursorResult<Topic[]> = api.getTopicFeeds(feedType)

useHead({
  title: useSiteTitle(feedInfo.value.title),
})
</script>
