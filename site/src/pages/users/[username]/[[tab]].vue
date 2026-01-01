<template>
  <ProfileHeader :user="user" />

  <div class="flex flex-col lg:flex-row gap-6 mt-6">
    <div class="w-full lg:w-80 lg:flex-shrink-0">
      <ProfileSidebar :user="user" />
    </div>

    <div class="flex-1 min-w-0 gaming-card rounded-2xl overflow-hidden min-h-[500px]">
      <TabsNavigation :user="user" :activeTab="activeTab" />
      <div class="p-6">
        <UserTopics v-show="activeTab === 'topics'" :user="user" />
        <UserActivity v-show="activeTab === 'activity'" :user="user" />
        <UserFollowing v-show="activeTab === 'following'" :user="user" />
        <UserFollowers v-show="activeTab === 'followers'" :user="user" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'
import TabsNavigation from '~/components/users/TabsNavigation.vue'
import UserTopics from '~/components/users/UserTopics.vue'
import UserActivity from '~/components/users/UserActivity.vue'
import UserFollowing from '~/components/users/UserFollowing.vue'
import UserFollowers from '~/components/users/UserFollowers.vue'
import type { UserProfile } from '@/types/user'

definePageMeta({
  layout: 'default',
})

const i18n = useI18n()
const route = useRoute()
const api = useApi()

const username = route.params.username as string

const user: UserProfile = await api.getUser(username).catch(() => {
  throw createError({ statusCode: 404, statusMessage: 'User not found' })
})

const allowedTabs = ['topics', 'activity', 'followers', 'following']
const tabParam: string = route.params.tab as string
if (tabParam && !allowedTabs.includes(tabParam)) {
  throw createError({ statusCode: 404, statusMessage: 'Tab not found' })
}

const activeTab = computed(() => tabParam || 'topics')

useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
