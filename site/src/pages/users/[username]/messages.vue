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
          <!-- Animated Background -->
          <div class="absolute inset-0 bg-gradient-to-br from-purple-500/5 via-transparent to-pink-500/5"></div>

          <div class="relative flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div
                class="w-14 h-14 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center ring-2 ring-purple-500/30 shadow-lg shadow-purple-500/20">
                <Icon name="Fa7SolidEnvelope" class="text-purple-400 text-2xl animate-pulse" />
              </div>
              <div>
                <h2
                  class="gaming-title text-3xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent font-['Saira Semi Condensed']">
                  {{ $t('message.messages') }}
                </h2>
                <p class="text-sm text-gray-400 mt-1">
                  {{ $t('message.stay_updated') || 'Stay updated with your activity' }}
                </p>
              </div>
            </div>

            <!-- Mark all as read button -->
            <button
              class="px-4 py-2 rounded-lg bg-purple-500/10 hover:bg-purple-500/20 text-purple-400 text-sm font-semibold transition-all duration-300 border border-purple-500/20 hover:border-purple-500/40">
              <Icon name="Fa7SolidCheckDouble" class="inline mr-2" />
              {{ $t('message.mark_all_read') || 'Mark all read' }}
            </button>
          </div>
        </div>

        <!-- Messages List -->
        <div class="space-y-3">
          <LoadMoreAsync v-slot="{ items }" :cursor="messagesCursor">
            <div v-if="items && items.length" class="space-y-3">
              <div
                v-for="message in items"
                :key="message.id"
                :class="[
                  'gaming-card rounded-xl overflow-hidden transition-all duration-200 hover:brightness-110',
                  message.status === 0 ? 'ring-2 ring-purple-500/30' : '',
                ]">
                <div class="relative">
                  <!-- Unread indicator -->
                  <div
                    v-if="message.status === 0"
                    class="absolute top-0 left-0 w-1 h-full bg-gradient-to-b from-purple-500 to-pink-500">
                  </div>

                  <div class="flex gap-4 p-5 pl-6">
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
                              {{ message.from.nickname }}
                            </nuxt-link>
                            <span v-else class="font-bold text-white text-base">
                              {{ message.from.nickname }}
                            </span>

                            <!-- Dynamic Title with clickable content link -->
                            <span class="text-gray-400 text-sm font-medium">
                              <i18n-t v-if="message.type === 0" keypath="message.commented_on_topic"
                                tag="span">
                                <template #topic>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.topic') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <i18n-t v-else-if="message.type === 1" keypath="message.replied_to_comment" tag="span">
                                <template #comment>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.comment') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <i18n-t v-else-if="message.type === 2" keypath="message.liked_topic" tag="span">
                                <template #topic>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.topic') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <i18n-t v-else-if="message.type === 3" keypath="message.favorited_topic" tag="span">
                                <template #topic>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.topic') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <i18n-t v-else-if="message.type === 4" keypath="message.recommended_topic" tag="span">
                                <template #topic>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.topic') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <i18n-t v-else-if="message.type === 5" keypath="message.deleted_topic" tag="span">
                                <template #topic>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.topic') }}
                                  </a>
                                </template>
                              </i18n-t>
                              <span v-else-if="message.type === 6">
                                {{ $t('message.follow_you') }}
                              </span>
                              <i18n-t v-else-if="message.type === 7" keypath="message.liked_comment" tag="span">
                                <template #comment>
                                  <a :href="message.detailUrl" target="_blank"
                                    class="text-purple-400 hover:text-purple-300 transition-colors">
                                    {{ $t('message.comment') }}
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

                      <!-- Actions -->
                      <div class="flex items-center gap-3 mt-3">
                        <!-- Mark as read -->
                        <button v-if="message.status === 0"
                          class="inline-flex items-center gap-2 px-3 py-2 rounded-lg bg-gray-700/30 hover:bg-gray-700/50 text-gray-400 hover:text-gray-300 text-xs font-medium transition-all duration-200">
                          <Icon name="Fa7SolidCheck" class="text-[10px]" />
                          {{ $t('message.mark_read') }}
                        </button>
                      </div>
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
import type { UserProfile } from '@/types'

// Message types enum
const MessageType = {
  TopicComment: 0,
  CommentReply: 1,
  TopicLike: 2,
  TopicFavorite: 3,
  TopicRecommend: 4,
  TopicDelete: 5,
  UserFollow: 6,
  CommentLike: 7,
} as const

definePageMeta({
  layout: 'profile',
})

const i18n = useI18n()
const route = useRoute()
const api = useApi()
const username = route.params.username as string

const user: UserProfile = await api.getUser(username).catch(() => {
  throw createError({ statusCode: 404, statusMessage: 'User not found' })
})

// Create cursor for loading messages
const messagesCursor = new CursorResult('/api/users/me/messages')

// Helper functions for message type styling
const getMessageTypeIcon = (type: number): string => {
  const icons: Record<number, string> = {
    [MessageType.TopicComment]: 'Fa7SolidComment',
    [MessageType.CommentReply]: 'Fa7SolidReply',
    [MessageType.TopicLike]: 'Fa7SolidHeart',
    [MessageType.TopicFavorite]: 'Fa7SolidStar',
    [MessageType.TopicRecommend]: 'Fa7SolidThumbsUp',
    [MessageType.TopicDelete]: 'Fa7SolidTrash',
    [MessageType.UserFollow]: 'Fa7SolidUserPlus',
    [MessageType.CommentLike]: 'Fa7SolidHeart',
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

const getMessageTypeBadgeClass = (type: number): string => {
  // Unified purple/pink badge for all types
  return 'bg-gradient-to-r from-purple-500/20 to-pink-500/20 border border-purple-500/30 text-purple-300'
}

const getMessageTypeBorderColor = (type: number): string => {
  // Unified purple border for all types
  return 'border-purple-500/50'
}

const getMessageTypeButtonClass = (type: number): string => {
  // Unified purple/pink button for all types
  return 'bg-gradient-to-r from-purple-600/20 to-pink-600/20 hover:from-purple-600/30 hover:to-pink-600/30 border border-purple-500/30 hover:border-purple-500/50 text-purple-300'
}

const getMessageTypeLabel = (type: number): string => {
  const labels: Record<number, string> = {
    [MessageType.TopicComment]: i18n.t('message.type.comment') || 'Comment',
    [MessageType.CommentReply]: i18n.t('message.type.reply') || 'Reply',
    [MessageType.TopicLike]: i18n.t('message.type.like') || 'Like',
    [MessageType.TopicFavorite]: i18n.t('message.type.favorite') || 'Favorite',
    [MessageType.TopicRecommend]: i18n.t('message.type.recommend') || 'Recommend',
    [MessageType.TopicDelete]: i18n.t('message.type.delete') || 'Delete',
    [MessageType.UserFollow]: i18n.t('message.type.follow') || 'Follow',
    [MessageType.CommentLike]: i18n.t('message.type.like') || 'Like',
  }
  return labels[type] || i18n.t('message.type.notification') || 'Notification'
}

const getMessageTypeActionLabel = (type: number): string => {
  const labels: Record<number, string> = {
    [MessageType.TopicComment]: i18n.t('message.action.view_comment') || 'View Comment',
    [MessageType.CommentReply]: i18n.t('message.action.view_reply') || 'View Reply',
    [MessageType.TopicLike]: i18n.t('message.action.view_topic') || 'View Topic',
    [MessageType.TopicFavorite]: i18n.t('message.action.view_topic') || 'View Topic',
    [MessageType.TopicRecommend]: i18n.t('message.action.view_topic') || 'View Topic',
    [MessageType.TopicDelete]: i18n.t('message.action.view_details') || 'View Details',
    [MessageType.UserFollow]: i18n.t('message.action.view_profile') || 'View Profile',
    [MessageType.CommentLike]: i18n.t('message.action.view_comment') || 'View Comment',
  }
  return labels[type] || i18n.t('message.action.view_details') || 'View Details'
}

useHead({
  title: useSiteTitle(i18n.t('page.messages', { nickname: user.nickname })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
