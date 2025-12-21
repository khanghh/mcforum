<template>
  <div class="space-y-4">
    <div v-for="topic in topics" :key="topic.id" class="gaming-card rounded-xl p-5 hover:bg-white/5 transition-colors">
      <div class="flex-1 min-w-0">
        <h3 class="text-lg font-bold mb-2">
          <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
            class="text-white hover:text-purple-400 transition-colors">
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
            <my-avatar :user="topic.user" :size="24"
              class="w-6 h-6 rounded border border-purple-500/50 flex-shrink-0" />
            <nuxt-link :to="`/user/${topic.user.id}`"
              class="font-bold text-purple-300 gaming-title text-sm hover:text-purple-200 transition-colors">
              {{ topic.user.nickname }}
            </nuxt-link>
          </div>

          <!-- Time -->
          <span class="text-xs text-gray-500">{{ usePrettyDate(topic.createTime) }}</span>
        </div>

        <p class="text-gray-400 text-sm mb-3 line-clamp-2">
          {{ topic.summary || topic.content }}
        </p>

        <div class="flex items-center gap-4 text-sm text-gray-400 flex-wrap">
          <!-- Like -->
          <button class="flex items-center gap-1 hover:text-red-400 transition-colors"
            :class="{ 'text-red-400': topic.liked }"
            @click.prevent="like(topic)">
            <FontAwesome :icon="['fas', 'heart']" :class="{ 'text-red-400': topic.liked }" />
            {{ topic.likeCount }}
          </button>

          <!-- Comment -->
          <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
            class="flex items-center gap-1 hover:text-blue-400 transition-colors">
            <FontAwesome :icon="['fas', 'comment']" class="text-blue-400" />
            {{ topic.commentCount }}
          </nuxt-link>

          <!-- View -->
          <nuxt-link :to="`/topics/${topic.slug || topic.id}`"
            class="flex items-center gap-1 hover:text-green-400 transition-colors">
            <FontAwesome :icon="['fas', 'eye']" class="text-green-400" />
            {{ topic.viewCount }}
          </nuxt-link>

          <!-- Category/Node -->
          <nuxt-link
            v-if="topic.node"
            :to="`/forums/${topic.node.slug}`"
            class="ml-auto px-2 py-1 bg-purple-500/20 text-purple-400 text-xs font-bold rounded hover:bg-purple-500/30 transition-colors">
            {{ topic.node.name }}
          </nuxt-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
const userStore = useUserStore()
const user = userStore.user

defineProps({
  topics: {
    type: Array,
    default: () => [],
    required: false,
  },
  showAvatar: {
    type: Boolean,
    default: true,
  },
  showPinned: {
    type: Boolean,
    default: false,
  },
})

async function like(topic) {
  try {
    if (topic.liked) {
      await useHttpDelete(`/api/topics/${topic.slug}/reactions/${user.id}`)
      topic.liked = false
      topic.likeCount = topic.likeCount > 0 ? topic.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    } else {
      await useHttpPostForm(`/api/topics/${topic.slug}/reactions`, {
        body: { type: 'like' },
      })
      topic.liked = true
      topic.likeCount++
      useMsgSuccess(i18n.t('message.liked_success'))
    }
  } catch (e) {
    useCatchError(e)
  }
}
</script>

<style scoped>
/* Scoped styles are handled by tailwind classes mostly */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
