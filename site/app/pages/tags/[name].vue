<template>
  <div class="flex gap-6">
    <GamingSidebar />

    <!-- Middle Content Area (no sidebar) -->
    <div class="flex-1 min-w-0">
      <!-- Hero Banner -->
      <div class="gaming-card sm:rounded-xl p-6 sm:p-8 mb-6 relative overflow-hidden">
        <div class="absolute inset-0 animated-border opacity-10"></div>
        <div class="relative z-10">
          <h1
            class="text-3xl sm:text-4xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 gaming-title py-3">
            #{{ tagTitle }}
          </h1>
          <p class="text-gray-300 mb-4">
            {{ i18n.t('feed.tag_desc') }}
          </p>
        </div>
      </div>

      <!-- Threads List -->
      <div id="threads-container" class="space-y-4">
        <LoadMoreAsync
          v-slot="{ items }"
          :cursor="topicsCursor">
          <GamingTopicList :topics="items" show-pinned />
        </LoadMoreAsync>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Topic } from '@/types'
import type { CursorResult } from '@/composables/api'

definePageMeta({ layout: 'default' })

const i18n = useI18n()
const route = useRoute()
const router = useRouter()

const name = (route.params.name as string) || ''

const api = useApi()

const topicsCursor: CursorResult<Topic[]> = await api.getTopicsByTag(name)

const tagTitle = computed(() => name)

useHead({
  title: useSiteTitle(tagTitle.value),
})

</script>
