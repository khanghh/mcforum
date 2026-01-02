<template>
  <div v-if="user">
    <ProfileHeader v-model:user="user" />

    <div class="flex flex-col lg:flex-row gap-6 mt-6">
      <div class="w-full lg:w-80 lg:flex-shrink-0">
        <ProfileSidebar :user="user" />
      </div>

      <div class="flex-1 min-w-0 gaming-card rounded-2xl overflow-hidden min-h-[500px] p-6">
        <div class="space-y-4">
          <LoadMoreAsync v-slot="{ items }" :cursor="favoritesCursor">
            <GamingTopicList :topics="items" :show-pinned="false" />
          </LoadMoreAsync>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'

definePageMeta({
  middleware: ['auth'],
  layout: 'default',
})

const i18n = useI18n()
const api = useApi()
const userStore = useUserStore()

const { user } = storeToRefs(userStore)

const favoritesCursor = api.getMyFavorites()

useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.value?.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
