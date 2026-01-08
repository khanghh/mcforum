<template>
  <div class="relative flex gap-4 mb-8">
    <!-- Type Icon Badge -->
    <div class="flex-shrink-0 relative">
      <div class="relative">
        <Avatar :src="message.from.avatar" :username="message.from.username" :size="70"
          class="rounded border-2 border-purple-300" />
        <div
          class="absolute -bottom-1 -right-1 w-6 h-6 rounded-lg flex items-center justify-center shadow-lg bg-gradient-to-br from-purple-500 to-pink-600">
          <Icon :name="getMessageTypeIcon(message.type)" class="text-white text-xs" size="16" />
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

            <span class="text-gray-400 text-sm font-medium">
              {{ getActivityTitle(message.type) }}
            </span>
          </div>
        </div>

        <!-- Time -->
        <span class="text-xs text-gray-500 whitespace-nowrap flex items-center gap-1.5">
          <Icon name="Fa7SolidClock" class="text-[10px]" />
          {{ usePrettyDate(message.createTime) }}
        </span>
      </div>

      <div
        class="relative mb-3 pl-2 py-2 pr-3 rounded-lg bg-gradient-to-r from-gray-700/30 to-gray-700/10 border-l-2 border-purple-500/50 flex items-center justify-between">
        <div class="flex items-center">
          <Icon name="Fa7SolidQuoteLeft" class="text-xs opacity-30 mr-1" />
          <span class="text-sm text-gray-400 leading-relaxed font-medium gaming-title">
            {{ message.title }}
          </span>
        </div>
        <nuxt-link :to="message.detailUrl" class="text-purple-400 hover:underline text-sm">
          {{ $t('actions.view') }} ›
        </nuxt-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type Message, MessageType, MessageStatus } from '@/types'
const i18n = useI18n()

defineProps({
  message: {
    type: Object as () => Message,
    required: true,
  },
})
// Helper functions for message type styling
const getMessageTypeIcon = (type: number): string => {
  const icons: Record<number, string> = {
    [MessageType.TopicComment]: 'Fa7SolidCommentDots',
    [MessageType.CommentReply]: 'Fa7SolidMailReply',
    [MessageType.TopicLike]: 'MaterialSymbolsThumbUp',
    [MessageType.TopicFavorite]: 'Fa7SolidStar',
    [MessageType.TopicRecommended]: 'PhSealCheckFill',
    [MessageType.TopicPinned]: 'Fa7SolidThumbTack',
    [MessageType.CommentLike]: 'MaterialSymbolsThumbUp',
  }
  return icons[type] || 'Fa7SolidBell'
}

function getActivityTitle(type: number) {
  const titles: Record<number, string> = {
    [MessageType.TopicComment]: i18n.t('activity.commentd_on_a_topic'),
    [MessageType.CommentReply]: i18n.t('activity.replied_a_comment'),
    [MessageType.TopicLike]: i18n.t('activity.liked_a_topic'),
    [MessageType.TopicFavorite]: i18n.t('activity.favorited_a_topic'),
    [MessageType.TopicRecommended]: i18n.t('activity.recommended_a_topic'),
    [MessageType.TopicPinned]: i18n.t('activity.pinned_a_topic'),
    [MessageType.CommentLike]: i18n.t('activity.liked_a_comment'),
  }
  return titles[type] || ''
}
</script>