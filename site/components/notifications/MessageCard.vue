<template>
  <div
    :class="['gaming-card rounded-2xl p-6 relative overflow-hidden', message.status === MessageStatus.StatusUnread ? 'ring-2 ring-purple-500/30' : '']">
    <!-- bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] border border-[rgba(139,92,246,0.2)] rounded-2xl -->
    <div class="relative flex gap-4">
      <!-- Type Icon Badge -->
      <div class="flex-shrink-0 relative">
        <div class="relative">
          <!-- <Avatar :src="comment.user.avatar" :username="comment.user.username" :size="48"
                class="rounded border-2 border-purple-300" /> -->
          <Avatar :src="message.from.avatar" :username="message.from.username" :size="56"
            class="rounded border-2 border-purple-300" />

          <!-- Type badge overlay -->
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
              <span v-else class="font-bold text-white text-base">
                {{ message.from.nickname }}
              </span>

              <!-- Dynamic Title with clickable content link -->
              <span class="text-gray-400 text-sm font-medium">
                <i18n-t v-if="message.type === MessageType.TopicComment"
                  keypath="message.commented_on_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.CommentReply" keypath="message.replied_your_comment"
                  tag="span">
                  <template #comment>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.comment') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicLike" keypath="message.liked_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicFavorite" keypath="message.favorited_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicRecommended"
                  keypath="message.recommended_your_topic" tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicPinned" keypath="message.pinned_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicDelete" keypath="message.deleted_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <span v-else-if="message.type === MessageType.UserFollow">
                  {{ $t('message.follow_you') }}
                </span>
                <i18n-t v-else-if="message.type === MessageType.CommentLike" keypath="message.liked_your_comment"
                  tag="span">
                  <template #comment>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.comment') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.FollowingUserCreateTopic"
                  keypath="message.following_user_create_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicApproved" keypath="message.approved_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
                  </template>
                </i18n-t>
                <i18n-t v-else-if="message.type === MessageType.TopicRejected" keypath="message.rejected_your_topic"
                  tag="span">
                  <template #topic>
                    <nuxt-link :to="message.detailUrl"
                      class="text-purple-400 hover:text-purple-300 transition-colors">
                      {{ $t('message.topic') }}
                    </nuxt-link>
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
          <div v-if="message.type === MessageType.TopicComment || message.type === MessageType.CommentReply"
            class="relative mb-3 pl-4 py-2.5 pr-3 rounded-lg bg-gradient-to-r from-gray-700/30 to-gray-700/10 border-l-2 border-purple-500/50">
            <Icon name="Fa7SolidQuoteLeft" class="text-xs opacity-30 mr-1" />
            <span class="text-sm text-gray-400 leading-relaxed">
              {{ message.quoteContent }}
            </span>
          </div>

          <!-- Show as link for other types -->
          <div v-else class="mb-3">
            <a :href="message.detailUrl"
              class="inline-flex items-center gap-2 text-base font-medium text-purple-400 hover:text-purple-300 transition-colors group/quote">
              <span class="gaming-title text-lg decoration-purple-500/30 group-hover/quote:decoration-purple-500/60">
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
</template>

<script setup lang="ts">
import { type Message, MessageType, MessageStatus } from '@/types'

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
    [MessageType.TopicDelete]: 'Fa7SolidTrashCan',
    [MessageType.UserFollow]: 'Fa7SolidUserPlus',
    [MessageType.CommentLike]: 'MaterialSymbolsThumbUp',
    [MessageType.FollowingUserCreateTopic]: 'Fa7SolidNewspaper',
    [MessageType.TopicApproved]: 'Fa7SolidCheck',
    [MessageType.TopicRejected]: 'Fa7SolidTimesCircle',
  }
  return icons[type] || 'Fa7SolidBell'
}
</script>
