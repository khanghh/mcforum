<template>
  <div
    :class="[
      'relative hover:bg-white/5 rounded-lg p-5 transition-all duration-200 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] shadow-[0_0_0_1px_rgb(139,92,246,0.2)]',
      topic.pinned
        ? `border-l-4 ${borderColors[colorIndex]} hover:shadow-md hover:shadow-purple-500/20 hover:-translate-y-0.5`
        : 'shadow-[0_0_0_1px_rgb(139,92,246,0.2)]'
    ]">
    <div class="flex-1 min-w-0">
      <h2 class="font-bold mb-2 flex items-center gap-4">
        <nuxt-link :to="`/topics/${topic?.slug}`"
          class="text-white text-lg md:text-xl hover:text-purple-400 transition-colors gaming-title max-w-full truncate inline-block">
          {{ topic.title }}
        </nuxt-link>
        <div v-if="showPinned && topic.pinned"
          class="absolute -top-3 -right-3 flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 text-white text-xs font-bold shadow-lg shadow-purple-500/30 z-20">
          <Icon name="TablerPinnedFilled" class="w-3.5 h-3.5" />
          <span class="tracking-wide">{{ $t('feed.pinned') }}</span>
        </div>
      </h2>
      <div class="flex items-center gap-2 mb-2 flex-wrap">
        <!-- Avatar + user/info -->
        <div v-if="author" class="flex items-center gap-3">
          <Avatar :src="author.avatar" :username="author.username" size="40"
            class="rounded border border-purple-300 flex-shrink-0" />
          <div class="flex flex-col min-w-0">
            <nuxt-link :to="`/users/${author.username}`"
              class="font-bold text-gray-300 gaming-title text-md hover:text-purple-400 transition-colors truncate">
              {{ author.username }}
            </nuxt-link>
            <span class="text-xs text-gray-500 mt-0.5">{{ usePrettyDate(topic.createTime) }}</span>
          </div>
        </div>
      </div>

      <p class="text-gray-300 text-sm mb-3 line-clamp-2">
        {{ topic.summary || topic.content }}
      </p>

      <div v-if="topic.tags && topic.tags.length" class="flex flex-wrap gap-2 mb-3">
        <nuxt-link v-for="tag in topic.tags" :key="tag" :to="`/tags/${tag}`"
          class="px-3 py-1 bg-purple-500/10 text-purple-300 text-xs rounded-full border border-purple-500/20 hover:bg-purple-500/20 transition-colors">
          #{{ tag }}
        </nuxt-link>
      </div>

      <div class="flex items-center gap-4 text-sm text-gray-500 flex-wrap">
        <!-- Like -->
        <LikeButton :liked="topic.liked" :count="topic.likeCount" :onClick="likeTopic" :disabled="!user" />

        <!-- Comment -->
        <nuxt-link :to="`/topics/${topic?.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <Icon name="TablerMessageCircle" />
          {{ topic.commentCount }}
        </nuxt-link>

        <!-- View -->
        <nuxt-link :to="`/topics/${topic?.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <Icon name="IcRoundRemoveRedEye" />
          {{ topic.viewCount }}
        </nuxt-link>

        <!-- Forum Name -->
        <span v-if="topic.forum" :class="[
          'ml-auto px-3 py-1.5 text-xs font-bold rounded-lg bg-gradient-to-r border',
          labelColors[colorIndex]
        ]">
          {{ topic.forum.name.toUpperCase() }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import type { Topic } from '@/types'

const i18n = useI18n()
const userStore = useUserStore()
const api = useApi()

const { user } = storeToRefs(userStore)

type Props = {
  topic: Topic
  showPinned?: boolean
}
const props = defineProps<Props>()
const topic = props.topic

const author = computed(() => topic.user)

function slugHash(str: string): number {
  if (!str) return 0
  let h = 0
  for (let i = 0; i < str.length; i++) {
    h = (h << 5) - h + str.charCodeAt(i)
    h |= 0
  }
  return h >>> 0
}

const labelColors = [
  'from-red-600/20 to-pink-600/20 text-red-300 border-red-500/30',
  'from-green-600/20 to-emerald-600/20 text-green-300 border-green-500/30',
  'from-blue-600/20 to-cyan-600/20 text-blue-300 border-blue-500/30',
  'from-yellow-600/20 to-amber-600/20 text-yellow-300 border-yellow-500/30',
  'from-pink-600/20 to-rose-600/20 text-pink-300 border-pink-500/30',
  'from-purple-600/20 to-pink-600/20 text-purple-300 border-purple-500/30',
  'from-indigo-600/20 to-violet-600/20 text-indigo-300 border-indigo-500/30',
  'from-teal-600/20 to-emerald-600/20 text-teal-300 border-teal-500/30',
  'from-orange-600/20 to-amber-600/20 text-orange-300 border-orange-500/30',
  'from-cyan-600/20 to-blue-600/20 text-cyan-300 border-cyan-500/30',
]

const borderColors = [
  'border-red-500',
  'border-green-500',
  'border-blue-500',
  'border-yellow-500',
  'border-pink-500',
  'border-purple-500',
  'border-indigo-500',
  'border-teal-500',
  'border-orange-500',
  'border-cyan-500',
]

const colorIndex = slugHash(topic?.forum?.slug) % labelColors.length || 0

async function likeTopic() {
  try {
    if (!user.value || !user.value.id) {
      return
    }

    const slug = topic?.slug
    if (topic.liked) {
      await api.removeTopicReaction(slug)
      topic.liked = false
      topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0
    }
    else {
      await api.addTopicReaction(slug, 'like')
      topic.liked = true
      topic.likeCount = (topic.likeCount || 0) + 1
    }
  }
  catch (e) {
    useCatchError(e)
  }
}

</script>
