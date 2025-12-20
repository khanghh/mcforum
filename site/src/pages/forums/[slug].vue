<template>
  <section v-if="forum" class="main">
    <div class="container main-container left-main size-300">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <create-topic-btn v-if="user" />
            <forum-sidebar />
          </div>
          <div class="topics-main">
            <load-more-async
              ref="loader"
              v-slot="{ items }"
              :url="`/api/forums/${slug}`">
              <topic-list :topics="items" show-pinned />
            </load-more-async>
          </div>
        </div>
      </div>
      <div class="right-container">
        <check-in />
        <site-notice />
        <score-rank />
        <friend-links />
      </div>
    </div>
  </section>
  <NotFound v-else />
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

const slug = route.params.slug || 'whats-new'
const loader = ref(null)

const { data: forum } = await useAsyncData(() => useHttpGet(`/api/forums/${slug}/info`))

useHead({
  title: useSiteTitle(forum.value?.name || i18n.t('page.not_found')),
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

<style lang="scss" scoped></style>
