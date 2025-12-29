<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6">
    <ProfileHeader :user="user" />

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
                <Icon name="Fa7SolidEnvelope" class="text-purple-400 text-2xl animate-pulse" />
              </div>
              <div>
                <h2
                  class="gaming-title text-3xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent">
                  {{ $t('message.messages') }}
                </h2>
                <p class="text-sm text-gray-400 mt-1">
                  {{ $t('message.stay_updated') || 'Stay updated with your activity' }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Messages List -->
        <div class="space-y-3">
          <LoadMoreAsync v-slot="{ items }" :cursor="messagesCursor">
            <div v-if="items && items.length" class="space-y-3">
              <div
                v-for="message in items"
                :key="message.id"
                class="gaming-card rounded-2xl p-6 relative overflow-hidden">
                <!--bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] border border-[rgba(139,92,246,0.2)] rounded-2xl-->
                <div class="relative flex gap-4">
                  <!-- Type Icon Badge -->
                  <div class="flex-shrink-0 relative">
                    <div class="relative">
                      <Avatar :user="message.from" :size="56"
                        class="rounded ring-1 ring-offset-1 ring-offset-gray-800"
                        :class="getMessageTypeRingColor(message.type)" />

                      <!-- Type badge overlay -->
                      <div
                        :class="[
                          'absolute -bottom-1 -right-1 w-7 h-7 rounded-lg flex items-center justify-center shadow-lg',
                          getMessageTypeGradient(message.type),
                        ]">
                        <Icon :name="getMessageTypeIcon(message.type)" class="text-white text-xs" />
                      </div>
                    </div>
                  </div>

                  <!-- Message Content -->
                  <div class="flex-1 min-w-0">
                    <!-- Header -->
                    <div class="flex items-start justify-between gap-3 mb-2">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 flex-wrap">
                          <nuxt-link
                            v-if="message.from.id > 0"
                            :to="'/users/' + message.from.username"
                            class="gaming-title font-bold text-white hover:text-purple-400 transition-colors text-base">
                            {{ message.from.username }}
                          </nuxt-link>
                          <span v-else class="font-bold text-white text-base">
                            {{ message.from.nickname }}
                          </span>

                          <!-- Dynamic Title with clickable content link -->
                          <span class="text-gray-400 text-sm font-medium">
                            <i18n-t v-if="message.type === MessageType.TopicComment"
                              keypath="message.commented_on_topic"
                              tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.CommentReply"
                              keypath="message.replied_your_comment" tag="span">
                              <template #comment>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.comment') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.TopicLike"
                              keypath="message.liked_your_topic" tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.TopicFavorite"
                              keypath="message.favorited_your_topic" tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.TopicRecommended"
                              keypath="message.recommended_your_topic" tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.TopicPinned"
                              keypath="message.pinned_your_topic" tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.TopicDelete"
                              keypath="message.deleted_your_topic" tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                            <span v-else-if="message.type === MessageType.UserFollow">
                              {{ $t('message.follow_you') }}
                            </span>
                            <i18n-t v-else-if="message.type === MessageType.CommentLike"
                              keypath="message.liked_your_comment" tag="span">
                              <template #comment>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.comment') }}
                                </a>
                              </template>
                            </i18n-t>
                            <i18n-t v-else-if="message.type === MessageType.FollowingUserCreateTopic"
                              keypath="message.following_user_create_topic"
                              tag="span">
                              <template #topic>
                                <a :href="message.detailUrl" target="_blank"
                                  class="text-purple-400 hover:text-purple-300 transition-colors">
                                  {{ $t('message.topic') }}
                                </a>
                              </template>
                            </i18n-t>
                          </span>
                        </div>
                      </div>

                      <!-- Time -->
                      <span class="text-xs text-gray-500 whitespace-nowrap flex items-center gap-1.5">
                        <Icon name="Fa7SolidClock" class="text-[10px]" />
                        {{ usePrettyDate(message.createTime) }}
                      </span>
                    </div>

                    <!-- Quote Content (for comments) or Link (for other types) -->
                    <div v-if="message.quoteContent">
                      <!-- Show as quote for comments and replies -->
                      <div v-if="message.type === 0 || message.type === 1"
                        class="relative mb-3 pl-4 py-2.5 pr-3 rounded-lg bg-gradient-to-r from-gray-700/30 to-gray-700/10 border-l-2"
                        :class="getMessageTypeBorderColor(message.type)">
                        <Icon name="Fa7SolidQuoteLeft" class="text-xs opacity-30 mr-1" />
                        <span class="text-sm text-gray-400 leading-relaxed">
                          {{ message.quoteContent }}
                        </span>
                      </div>

                      <!-- Show as link for other types -->
                      <div v-else class="mb-3">
                        <a :href="message.detailUrl" target="_blank"
                          class="inline-flex items-center gap-2 text-base font-medium text-purple-400 hover:text-purple-300 transition-colors group/quote">
                          <span
                            class="gaming-title text-lg decoration-purple-500/30 group-hover/quote:decoration-purple-500/60">
                            {{ message.quoteContent }}
                          </span>
                        </a>
                      </div>
                    </div>

                    <!-- Main Content -->
                    <div v-if="message.content" class="text-gray-200 text-sm leading-relaxed mb-3">
                      {{ message.content }}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Empty state -->
            <div v-else class="gaming-card rounded-2xl p-12 text-center">
              <div
                class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center">
                <Icon name="Fa7SolidBellSlash" class="text-purple-400 text-3xl" />
              </div>
              <h3 class="text-xl font-bold text-gray-300 mb-2">
                {{ $t('message.no_notifications') }}
              </h3>
              <p class="text-gray-500">
                {{ $t('message.no_notifications_desc') }}
              </p>
            </div>
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
import { CursorResult } from '@/composables/api'
import { type UserProfile, MessageType } from '@/types'

// // Message types enum
// const MessageType = {
//   TopicComment: 0,
//   CommentReply: 1,
//   TopicLike: 2,
//   TopicFavorite: 3,
//   TopicRecommend: 4,
//   TopicDelete: 5,
//   UserFollow: 6,
//   CommentLike: 7,
//   FollowingUserCreateTopic: 8,
// } as const

definePageMeta({
  layout: 'profile',
})

const i18n = useI18n()
const route = useRoute()
const api = useApi()
const userStore = useUserStore()
const username = route.params.username as string

const user: UserProfile | null = await userStore.getCurrent()
if (user == null) {
  throw createError({ statusCode: 404, statusMessage: 'User not found' })
}

// Create cursor for loading messages
const messagesCursor = api.getMessages()

// Helper functions for message type styling
const getMessageTypeIcon = (type: number): string => {
  const icons: Record<number, string> = {
    [MessageType.TopicComment]: 'Fa7SolidCommentDots',
    [MessageType.CommentReply]: 'Fa7SolidMailReply',
    [MessageType.TopicLike]: 'MaterialSymbolsThumbUp',
    [MessageType.TopicFavorite]: 'Fa7SolidStar',
    [MessageType.TopicRecommended]: 'PhSealCheckFill',
    [MessageType.TopicPinned]: 'Fa7SolidThumbTack',
    [MessageType.TopicDelete]: 'Fa7SolidTrashCan',
    [MessageType.UserFollow]: 'Fa7SolidUserPlus',
    [MessageType.CommentLike]: 'MaterialSymbolsThumbUp',
    [MessageType.FollowingUserCreateTopic]: 'Fa7SolidPenNib',
  }
  return icons[type] || 'Fa7SolidBell'
}

const getMessageTypeGradient = (type: number): string => {
  // Unified purple/pink gradient for all types
  return 'bg-gradient-to-br from-purple-500 to-pink-600'
}

const getMessageTypeRingColor = (type: number): string => {
  // Unified purple ring for all types
  return 'ring-purple-500/50'
}

const getMessageTypeBorderColor = (type: number): string => {
  // Unified purple border for all types
  return 'border-purple-500/50'
}

useHead({
  title: useSiteTitle(i18n.t('page.messages', { nickname: user.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
