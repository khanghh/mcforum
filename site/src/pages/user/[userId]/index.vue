<template>
  <section class="main">
    <div class="container">
      <user-profile :user="user" />

      <div class="container main-container right-main size-300">
        <user-center-sidebar :user="user" />
        <div class="right-container">
          <div class="tabs-warp">
            <div class="tabs">
              <ul>
                <li :class="{ 'is-active': activeTab === 'topics' }">
                  <nuxt-link :to="'/user/' + user.id">
                    <span class="icon is-small">
                      <icon name="MessageSquareText" />
                    </span>
                    <span>{{ $t('page.topics') }}</span>
                  </nuxt-link>
                </li>
                <li :class="{ 'is-active': activeTab === 'articles' }">
                  <nuxt-link :to="'/user/' + user.id + '/articles'">
                    <span class="icon is-small">
                      <icon name="FileText" />
                    </span>
                    <span>{{ $t('page.articles') }}</span>
                  </nuxt-link>
                </li>
              </ul>
            </div>

            <div>
              <load-more-async
                v-slot="{ items }"
                url="/api/topic/user/topics"
                :params="{ userId: user.id }">
                <topic-list :topics="items" :show-avatar="false" />
              </load-more-async>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const user = await useHttpGet(`/api/user/${route.params.userId}`)
const activeTab = ref('topics')
useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.nickname })),
})
</script>

<style lang="scss" scoped>
.tabs-warp {
  background-color: var(--bg-color);
  padding: 0 10px 10px;

  .tabs {
    margin-bottom: 5px;
  }

  .more {
    text-align: right;
  }
}
</style>
