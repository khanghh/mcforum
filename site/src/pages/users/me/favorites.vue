<template>
  <ProfileHeader :user="user" />

  <div class="flex flex-col lg:flex-row gap-6 mt-6">
    <div class="w-full lg:w-80 lg:flex-shrink-0">
      <ProfileSidebar :user="user" />
    </div>

    <div class="flex-1 min-w-0 gaming-card rounded-2xl overflow-hidden min-h-[500px]">
    </div>
  </div>
</template>

<script setup lang="ts">
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'
import type { UserProfile } from '@/types/user'

definePageMeta({
  layout: 'default',
})

const i18n = useI18n()
const route = useRoute()
const api = useApi()

const user: UserProfile = await api.getCurrentUser().catch(() => {
  throw createError({ statusCode: 404, statusMessage: 'User not found' })
})

useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
