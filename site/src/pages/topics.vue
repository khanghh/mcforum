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
            <load-more-async v-slot="{ items }" url="/api/topics">
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
const userStore = useUserStore()
const { user } = storeToRefs(userStore)

useHead({
  title: useSiteTitle('All Topics'),
})
</script>

<style scoped>
.topics-nav {}
</style>
