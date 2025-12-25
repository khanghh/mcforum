<template>
  <div class="gaming-card rounded-xl p-5 hover:bg-white/5 transition-colors">
    <div class="flex-1 min-w-0">
      <h3 class="text-lg font-bold mb-2">
        <nuxt-link :to="`/topics/${topic.slug}`"
          class="text-white hover:text-purple-400 transition-colors gaming-title">
          <span v-if="showPinned && topic.pinned"
            class="mr-2 px-2 py-0.5 bg-red-500/20 text-red-400 text-xs font-bold rounded border border-red-500/50">
            PINNED
          </span>
          {{ topic.title }}
        </nuxt-link>
      </h3>
      <div class="flex items-center gap-2 mb-2 flex-wrap">
        <!-- Avatar -->
        <div v-if="author" class="flex items-center gap-2">
          <Avatar :src="author.avatar" :username="author.nickname" :size="24"
            class="w-6 h-6 rounded border border-purple-500/50 flex-shrink-0" />
          <nuxt-link :to="`/users/${author.username}`"
            class="font-bold text-purple-300 gaming-title text-md hover:text-purple-200 transition-colors">
            {{ author.nickname }}
          </nuxt-link>
        </div>

        <!-- Time -->
        <span class="text-xs text-gray-500">{{ usePrettyDate(topic.createTime) }}</span>
      </div>

      <p class="text-gray-300 text-sm mb-3 line-clamp-2">
        {{ topic.summary || topic.content }}
      </p>

      <div v-if="topic.tags && topic.tags.length" class="flex items-center gap-2 flex-wrap mb-3">
        <nuxt-link
          v-for="tag in topic.tags"
          :key="tag"
          :to="`/tags/${tag}`"
          class="px-2 py-0.5 bg-white/3 text-sm text-gray-400 rounded hover:text-gray-300 border border-white/10 hover:bg-white/5 transition-colors">
          #{{ tag }}
        </nuxt-link>
      </div>

      <div class="flex items-center gap-4 text-sm text-gray-500 flex-wrap">
        <!-- Like -->
        <button class="flex items-center gap-1 hover:text-blue-400 transition-colors"
          :class="{ 'text-blue-400': topic.liked }"
          @click.prevent="likeTopic">
          <Icon name="MaterialSymbolsThumbUp" :class="{ 'text-blue-400': topic.liked }" />
          {{ topic.likeCount }}
        </button>

        <!-- Comment -->
        <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <Icon name="TablerMessageCircle" />
          {{ topic.commentCount }}
        </nuxt-link>

        <!-- View -->
        <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <Icon name="IcRoundRemoveRedEye" />
          {{ topic.viewCount }}
        </nuxt-link>

        <!-- Forum Name -->
        <span v-if="topic.forum" :class="['ml-auto px-2 py-1 text-xs font-bold rounded', forumLabelColor]">
          {{ topic.forum.name.toUpperCase() }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import type { Topic } from '@/types'

const i18n = useI18n()
const userStore = useUserStore()
const user = userStore.user
const api = useApi()

type Props = {
  topic: Topic
  showPinned?: boolean
}
const props = defineProps<Props>()
const topic = props.topic

const author = computed(() => topic.user)

const forumLabelColor = computed(() => {
  if (!topic.forum) return ''
  const hash = slugHash(topic.forum.slug)
  return labelColors[hash % labelColors.length]
})

async function likeTopic() {
  try {
    if (!user || !user.id) {
      useCatchError(new Error(i18n.t('auth.login_required')))
      return
    }

    const slug = topic.slug
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


const labelColors = [
  'bg-red-500/20 text-red-400 border-red-500/50',
  'bg-green-500/20 text-green-400 border-green-500/50',
  'bg-blue-500/20 text-blue-400 border-blue-500/50',
  'bg-yellow-500/20 text-yellow-400 border-yellow-500/50',
  'bg-purple-500/20 text-purple-400 border-purple-500/50',
  'bg-pink-500/20 text-pink-400 border-pink-500/50',
  'bg-indigo-500/20 text-indigo-400 border-indigo-500/50',
  'bg-teal-500/20 text-teal-400 border-teal-500/50',
  'bg-orange-500/20 text-orange-400 border-orange-500/50',
  'bg-cyan-500/20 text-cyan-400 border-cyan-500/50',
]

function slugHash(str: string): number {
  let h = 0
  for (let i = 0; i < str.length; i++) {
    h = (h << 5) - h + str.charCodeAt(i)
    h |= 0
  }
  return h >>> 0
}
</script>
