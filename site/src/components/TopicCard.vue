<template>
  <div class="gaming-card rounded-xl p-5 hover:bg-white/5 transition-colors">
    <div class="flex-1 min-w-0">
      <h3 class="text-lg font-bold mb-2">
        <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
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
        <div class="flex items-center gap-2">
          <Avatar :src="topic.user.avatar" :username="topic.user.nickname" :size="24"
            class="w-6 h-6 rounded border border-purple-500/50 flex-shrink-0" />
          <nuxt-link :to="`/users/${topic.user.username}`"
            class="font-bold text-purple-300 gaming-title text-md hover:text-purple-200 transition-colors">
            {{ topic.user.nickname }}
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
        <button class="flex items-center gap-1 hover:text-red-400 transition-colors"
          :class="{ 'text-red-400': topic.liked }"
          @click.prevent="like(topic)">
          <FontAwesome :icon="['fas', 'heart']" :class="{ 'text-red-400': topic.liked }" />
          {{ topic.likeCount }}
        </button>

        <!-- Comment -->
        <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <FontAwesome :icon="['fas', 'comment']" />
          {{ topic.commentCount }}
        </nuxt-link>

        <!-- View -->
        <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
          class="flex items-center gap-1 hover:text-gray-400 transition-colors">
          <FontAwesome :icon="['fas', 'eye']" />
          {{ topic.viewCount }}
        </nuxt-link>

      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Topic } from '@/types'
const i18n = useI18n()
const userStore = useUserStore()
const user = userStore.user

type Props = {
  topic: Topic
  showPinned?: boolean
}
const props = defineProps<Props>()
  
const topic = props.topic

async function like(topic: Topic) {
  try {
    if (topic.liked) {
      await useHttpDelete(`/api/topics/${topic.slug}/reactions/${user.id}`)
      topic.liked = false
      topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0
      // useMsgSuccess(i18n.t('message.unliked_success'))
    } else {
      await useHttpPostForm(`/api/topics/${topic.slug}/reactions`, {
        body: { type: 'like' },
      })
      topic.liked = true
      topic.likeCount++
      // useMsgSuccess(i18n.t('message.liked_success'))
    }
  } catch (e) {
    useCatchError(e)
  }
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
