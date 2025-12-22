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
                  {{ forumTitle }}
                </h1>
                <p class="text-gray-300 mb-4">
                  {{ forumDescription }}
                </p>
                <div class="flex flex-wrap gap-3">
                  <button
                    class="px-4 py-2 bg-purple-600 text-white rounded-lg font-semibold text-sm hover:bg-purple-700 transition-colors"
                    @click="navigateTo('/topics/create')">
                    <FontAwesome :icon="['fas', 'plus']" class="mr-2" /> New Thread
                  </button>
                  <button
                    class="px-4 py-2 border border-purple-500/50 text-purple-300 rounded-lg font-semibold text-sm hover:bg-purple-500/20 transition-colors">
                    <FontAwesome :icon="['fas', 'filter']" class="mr-2" /> Filter
                  </button>
                </div>
              </div>
            </div>

            <!-- Threads List -->
            <div id="threads-container" class="space-y-4">
              <LoadMoreAsync
                ref="loader"
                v-slot="{ items }"
                :url="`/api/forums/${slug}`">
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

<script setup>
definePageMeta({
  layout: false,
})

const i18n = useI18n()
const route = useRoute()

const slug = route.params.slug || 'whats-new'

const { data, error } = await useAsyncData(() => useHttpGet('/api/forums'))

if (error?.value) {
  const status = error.value.statusCode || error.value.status || 500
  const message = error.value.message || i18n.t('page.server_error') || 'Server Error'
  throw createError({ statusCode: status, statusMessage: message, fatal: true })
}

const forum = computed(() => {
  return data.value?.find(f => f.slug === slug)
})

if (!forum.value) {
  throw createError({ statusCode: 404, statusMessage: i18n.t('page.not_found'), fatal: true })
}

const forumTitle = computed(() => {
  return forum.value ? forum.value.name : i18n.t('page.not_found')
})

const forumDescription = computed(() => {
  return forum.value ? forum.value.description : ''
})


useHead({
  title: useSiteTitle(forumTitle?.value),
  bodyAttrs: {
    class: 'bg-[#0f0f23]',
  },
  meta: [
    {
      hid: 'description',
      name: 'description',
      content: useSiteDescription(),
    },
    { hid: 'keywords', name: 'keywords', content: useSiteKeywords() },
  ],
})
</script>

<style lang="scss">
@import "~/assets/css/gaming-design.scss";
</style>
