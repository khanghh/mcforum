<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6">
    <ProfileHeader :user="user" />

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-6">
      <ProfileSidebar :user="user" />

      <div class="lg:col-span-2 gaming-card rounded-2xl overflow-hidden min-h-[500px]">
        <TabsNavigation v-model="activeTab" />
        <div class="p-6">
          <UserTopics v-show="activeTab === 'topics'" :user="user" />
          <UserActivity v-show="activeTab === 'activity'" :user="user" />
          <UserFollowing v-show="activeTab === 'following'" :user="user" />
          <UserFans v-show="activeTab === 'fans'" :user="user" />
          <UserFavorites v-show="activeTab === 'favorites'" :user="user" />
        </div>
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
import UserFans from '~/components/users/UserFans.vue'
import UserFavorites from '~/components/users/UserFavorites.vue'
import type { UserProfile } from '@/types'
import { ref, watch } from 'vue'

definePageMeta({
  layout: 'profile',
})

const i18n = useI18n()
const route = useRoute()
const router = useRouter()
const api = useApi()
const username = route.params.username as string

const user: UserProfile = await api.getUser(username).catch(() => {
  throw createError({ statusCode: 404, statusMessage: 'User not found' })
})

const allowedTabs = ['topics', 'activity', 'following', 'fans', 'favorites']
const tabParam: string = route.params.tab as string
if (tabParam && !allowedTabs.includes(tabParam)) {
  throw createError({ statusCode: 404, statusMessage: 'Tab not found' })
}

const initialTab = allowedTabs.includes(tabParam) ? tabParam : 'topics'
const activeTab = ref(initialTab)

// keep URL in sync when tab changes
watch(activeTab, (tab) => {
  const params = { ...route.params }
  if (tab === 'topics') delete params.tab
  else params.tab = tab
  router.push({ params }).catch(() => { })
})

useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.nickname })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>

<style scoped>
/* Scrollbar styling */
.custom-scrollbar::-webkit-scrollbar {
  height: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #4b5563;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>
