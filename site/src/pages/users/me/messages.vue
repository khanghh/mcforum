<template>
  <div v-if="user">
    <ProfileHeader v-model:user="user" />
    <div class="flex flex-col lg:flex-row gap-6 mt-6">
      <div class="w-full lg:w-80 lg:flex-shrink-0">
        <ProfileSidebar :user="user" />
      </div>

      <div class="flex-1 min-w-0">
        <!-- Messages Header -->
        <div class="gaming-card rounded-2xl p-6 mb-6 relative overflow-hidden">
          <div class="relative flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div
                class="w-14 h-14 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center ring-2 ring-purple-500/30 shadow-lg shadow-purple-500/20">
                <Icon name="Fa7SolidBell" class="text-purple-400 animate-pulse" size="25" />
              </div>
              <div>
                <h2
                  class="gaming-title text-3xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent">
                  {{ $t('page.messages') }}
                </h2>
                <p class="text-sm text-gray-400 mt-1">
                  {{ $t('page.messages_desc') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Messages List -->
        <div class="space-y-3">
          <LoadMoreAsync :cursor="messagesCursor">
            <template #default="{ items }">
              <div v-if="items && items.length" class="space-y-4">
                <MessageCard
                  v-for="message in items"
                  :key="message.id"
                  :message="message" />
              </div>
            </template>

            <template #empty>
              <div class="gaming-card rounded-2xl p-12 text-center">
                <div
                  class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                  <Icon name="TablerBellOff" class="text-purple-400" size="40" />
                </div>
                <h3 class="text-xl font-bold text-gray-300 mb-2">
                  {{ $t('message.no_notifications') }}
                </h3>
                <p class="text-gray-500">
                  {{ $t('message.no_notifications_desc') }}
                </p>
              </div>
            </template>
            <!-- Empty state -->
          </LoadMoreAsync>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'
import LoadMoreAsync from '~/components/LoadMoreAsync.vue'
import MessageCard from '~/components/notifications/MessageCard.vue'

definePageMeta({
  middleware: ['auth'],
  layout: 'default',
})

const i18n = useI18n()
const api = useApi()
const userStore = useUserStore()

const { user } = storeToRefs(userStore)

// Create cursor for loading messages
const messagesCursor = api.getMessages()

useHead({
  title: useSiteTitle(i18n.t('page.messages', { nickname: user.value?.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
