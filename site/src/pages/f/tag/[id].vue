<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <create-topic-btn v-if="user" />
            <forum-sidebar />
          </div>
          <div class="topics-main">
            <load-more-async
              v-slot="{ items }"
              url="/api/topic/tag/topics"
              :params="{ tagId: tagId }">
              <topic-list :topics="items" />
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
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const userStore = useUserStore()

const tagId = route.params.id
const { user } = storeToRefs(userStore)
const { data: tag } = await useAsyncData(() => useHttpGet(`/api/tag/${tagId}`))

useHead({
  title: useSiteTitle(i18n.t('page.by_tag', { tag: tag.value.name })),
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
