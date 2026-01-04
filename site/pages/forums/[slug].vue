<template>
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
            {{ forumTitle.toUpperCase() }}
          </h1>
          <p class="text-gray-300 mb-4">
            {{ forumDescription }}
          </p>
          <div class="flex flex-wrap gap-3">
            <button
              class="flex items-center justify-center px-4 py-2 bg-purple-600 text-white rounded-lg font-semibold text-sm hover:bg-purple-700 transition-colors"
              @click="navigateTo('/topics/create')">
              <Icon name="TablerPlus" class="mr-2" /> {{ $t('actions.create_topic') }}
            </button>
            <button
              class="flex items-center px-4 py-2 border border-purple-500/50 text-purple-300 rounded-lg font-semibold text-sm hover:bg-purple-500/20 transition-colors">
              <Icon name="Fa7SolidFilter" class="mr-2" /> {{ $t('actions.filter_topics') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Threads List -->
      <div id="threads-container" class="space-y-4">
        <LoadMoreAsync v-slot="{ items }" :cursor="forumCursor">
          <GamingTopicList :topics="items" show-pinned />
        </LoadMoreAsync>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useApi } from '@/composables/api'

// Expose auto-imports to template for TypeScript

definePageMeta({
  layout: 'default',
})

const i18n = useI18n()
const route = useRoute()

const slug = route.params.slug as string
const api = useApi()

const forumList = await api.getForumList().catch((err) => {
  const status = err.statusCode || 500
  const message = err.message || i18n.t('page.server_error')
  throw createError({ statusCode: status, statusMessage: message, fatal: true })
})

const forumInfo = forumList.find(forum => forum.slug === slug)
if (!forumInfo || !forumInfo.slug) {
  throw createError({ statusCode: 404, statusMessage: i18n.t('page.not_found'), fatal: true })
}

const forumTitle = computed(() => forumInfo.name || slug)
const forumDescription = computed(() => forumInfo.description || '')
const forumCursor = api.getForumTopics(slug)

useHead({
  title: useSiteTitle(forumInfo.name),
})
</script>
