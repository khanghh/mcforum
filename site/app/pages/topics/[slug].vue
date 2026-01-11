<template>
  <div itemscope itemtype="http://schema.org/BlogPosting">
    <!-- Breadcrumb -->
    <BreadCrumb :items="breadcrumbItems" />

    <!-- Notification for pending posts -->
    <div v-if="topic.status == TopicStatus.PendingReview"
      class="mb-6 p-4 rounded-lg bg-yellow-500/20 border border-yellow-500/50 text-yellow-200">
      <Icon name="TablerAlertTriangle" class="mr-2" /> {{ $t('publish.topic_under_review') }}
    </div>
    <div v-if="topic.status == TopicStatus.Hidden"
      class="mb-6 p-4 rounded-lg bg-red-500/20 border border-red-500/50 text-red-200">
      <Icon name="TablerAlertTriangle" class="mr-2" /> {{ $t('publish.topic_waiting_for_deletion') }}
    </div>

    <!-- Original Post -->
    <div
      class="gaming-card rounded-xl p-6 mb-6 border border-purple-500/20 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] relative overflow-hidden">
      <!-- Animated border overlay -->
      <!-- <div
        class="absolute inset-0 bg-[linear-gradient(45deg,#8b5cf6,#ec4899,#06b6d4,#8b5cf6)] bg-[length:300%_300%] opacity-10 pointer-events-none">
      </div> -->

      <div class="relative z-10">
        <div class="flex items-center justify-between gap-3 mb-4">
          <div class="flex items-center gap-3">
            <div class="relative group">
              <Avatar :src="topic.user.avatar"
                :username="topic.user.username"
                class="w-12 h-12 rounded border-2 border-purple-300 flex-shrink-0 object-cover" />
            </div>
            <div>
              <div class="flex items-center gap-2 flex-wrap">
                <nuxt-link :to="`/users/${topic.user.username}`"
                  class="font-bold text-purple-300 gaming-title text-lg hover:text-purple-400 transition-colors"
                  itemprop="author">
                  {{ topic.user.username }}
                </nuxt-link>
                <span v-if="topic.user.type === 1 || topic.user.id === 1"
                  class="px-2 py-0.5 bg-red-500/20 text-red-400 text-xs font-bold rounded">
                  ADMIN
                </span>
                <span class="px-2 py-0.5 bg-purple-500/20 text-purple-400 text-xs font-bold rounded">
                  LVL
                  {{ getCurrentLevel(topic.user.score) }}
                </span>
              </div>

              <div class="text-sm text-gray-500 mt-1">
                {{ usePrettyDate(topic.createTime) }}
              </div>

              <div class="flex items-center gap-2 text-sm text-gray-400 mt-1">
                <span v-if="topic.ipLocation" class="flex items-center gap-1">
                  <Icon name="TablerMapPinFilled" class="text-gray-500" /> {{ topic.ipLocation }}
                </span>
              </div>
            </div>
          </div>

          <!-- Manage Menu -->
          <TopicManageMenu v-if="canManage" v-model="topic" class="relative z-20" />
        </div>

        <TopicContent :topic="topic" />

        <div v-if="topic.tags && topic.tags.length" class="flex flex-wrap gap-2 mb-6">
          <nuxt-link v-for="tag in topic.tags" :key="tag" :to="`/tags/${tag}`"
            class="px-3 py-1 bg-purple-500/10 text-purple-300 text-xs rounded-full border border-purple-500/20 hover:bg-purple-500/20 transition-colors">
            #{{ tag }}
          </nuxt-link>
        </div>

        <div class="flex items-center justify-between pt-4 border-t border-purple-500/20">
          <div class="flex items-center gap-4">
            <LikeButton :liked="topic.liked" :count="topic.likeCount" @click="toggleLike" />

            <div class="flex items-center gap-1 text-gray-400">
              <Icon name="IcRoundRemoveRedEye" />
              <span>{{ topic.viewCount }}</span>
              <span class="hidden sm:inline">{{ $t('feed.view_count') }}</span>
            </div>

            <div class="flex items-center gap-1 text-gray-400">
              <Icon name="TablerMessageCircle" />
              <span>{{ topic.commentCount }}</span>
              <span class="hidden sm:inline">{{ $t('feed.comment_count') }}</span>
            </div>

            <FavoriteButton v-if="user" :favorited="topic.favorited" @click="toggleFavorite">
              <span>{{ $t('actions.favorite') }}</span>
            </FavoriteButton>
          </div>

          <div class="flex items-center gap-2">
            <!-- <button class="text-gray-400 hover:text-purple-400 transition-colors">
              <Icon name="Fa7SolidShare" />
            </button> -->
            <button class="text-gray-400 hover:text-red-400 transition-colors">
              <Icon name="Fa7SolidFlag" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Comments Section -->
    <div
      class="gaming-card rounded-xl p-6 border border-purple-500/20 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))]">
      <h2 class="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 flex items-center text-white">
        <Icon name="Fa7SolidComments" class="mr-3 text-blue-400" />
        {{ $t('feed.comments') }}
      </h2>

      <!-- Wrapped Comment Component -->
      <CommentSection
        :topic-slug="topic.slug"
        :comment-count="topic.commentCount"
        @created="commentCreated" />
    </div>
  </div>
</template>

<script setup>
import TopicContent from '~/components/topics/TopicContent.vue'
import BreadCrumb from '~/components/ui/BreadCrumb.vue'
import { TopicStatus } from '~/types'

const i18n = useI18n()
const route = useRoute()
const toast = useToast()
const api = useApi()
const userStore = useUserStore()

const slug = route.params.slug
if (!slug) {
  throw createError({ statusCode: 404, statusMessage: 'Topic not found', fatal: true })
}

const { data } = await useAsyncData(`topic`, () => api.getTopic(slug))
if (!data.value) {
  throw createError({ statusCode: 404, statusMessage: 'Topic not found', fatal: true })
}
const topic = ref(data.value)

const user = computed(() => userStore.user)
const canManage = computed(() => {
  if (!user.value) return false
  return userIsManager(user.value) || (topic.value && topic.value.user.id === user.value.id)
})

const breadcrumbItems = computed(() => {
  const items = [
    { label: 'Forums', to: '/' }
  ]
  if (topic.value.forum) {
    items.push({ label: topic.value.forum.name, to: `/forums/${topic.value.forum.slug}` })
  }
  items.push({ label: topic.value.title })
  return items
})


async function toggleLike() {
  if (topic.value.status !== TopicStatus.Active) return
  if (!user.value) return
  if (topic.value.liked) {
    await api.removeTopicReaction(slug)
    topic.value.liked = false
    topic.value.likeCount = topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0
  } else {
    await api.addTopicReaction(slug, 'like')
    topic.value.liked = true
    topic.value.likeCount++
  }
}

async function toggleFavorite() {
  if (topic.value.status !== TopicStatus.Active) return
  try {
    if (topic.value.favorited) {
      await api.setTopicFavorite(topic.value.id, false)
      topic.value.favorited = false
    } else {
      await api.setTopicFavorite(topic.value.id, true)
      topic.value.favorited = true
    }
  } catch (e) {
    const errMsg = e.data?.error?.message || e.message || e
    toast.error(i18n.t('message.opration_failed', { error: errMsg }))
  }
}

function commentCreated() {
  topic.value.commentCount++
}

useSeoMeta({
  title: `${topic.value.title} – MineViet Network`,
  ogTitle: `${topic.value.title} – MineViet Network`,
  description: 'Tham gia thảo luận về chủ đề này trên diễn đàn MineViet!',
  ogDescription: 'Tham gia thảo luận về chủ đề này trên diễn đàn MineViet!',
  twitterTitle: `${topic.value.title} – MineViet Network`,
  twitterDescription: 'Tham gia thảo luận về chủ đề này trên diễn đàn MineViet!',
})
</script>
