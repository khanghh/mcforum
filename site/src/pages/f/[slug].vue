<template>
  <section class="main">
    <div class="container main-container left-main size-320">
      <div class="left-container">
        <div class="main-content no-padding no-bg topics-wrapper">
          <div class="topics-nav">
            <forum-sidebar />
          </div>
          <div class="topics-main">
            <load-more-async
              ref="loader"
              v-slot="{ items }"
              :url="apiUrl">
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
</template>

<script setup>
const route = useRoute()

const loader = ref(null)
const slug = route.params.slug || 'whats-new'

const apiUrl = computed(() => `/api/forums/${slug}`)

useHead({
  title: useSiteTitle('test'),
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
