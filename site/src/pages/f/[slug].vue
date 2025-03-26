<template>
  <section class="main">
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
const userStore = useUserStore()

const slug = route.params.slug || 'whats-new'
const { user } = storeToRefs(userStore)
const loader = ref(null)

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
