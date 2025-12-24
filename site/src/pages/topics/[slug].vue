<template>
  <div v-if="topic"
    class="min-h-screen bg-gray-900 text-gray-100 font-sans selection:bg-purple-500 selection:text-white flex flex-col"
    style="background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);">
    <GamingNavbar />

    <!-- Main Content -->
    <main class="flex-grow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 py-6" itemscope itemtype="http://schema.org/BlogPosting">
        <!-- Notification for pending posts -->
        <div v-if="isPending" class="mb-6 p-4 rounded-lg bg-yellow-500/20 border border-yellow-500/50 text-yellow-200">
          <FontAwesome :icon="['fas', 'exclamation-triangle']" class="mr-2" /> {{ $t('message.post_under_review') }}
        </div>

        <!-- Breadcrumb -->
        <div class="mb-6">
          <nav class="flex text-sm text-gray-400">
            <nuxt-link to="/forums" class="hover:text-purple-400 transition-colors">Forums</nuxt-link>
            <span class="mx-2">›</span>
            <nuxt-link v-if="topic.forum" :to="`/forums/${topic.forum.slug}`"
              class="hover:text-purple-400 transition-colors">
              {{ topic.forum.name }}
            </nuxt-link>
            <span class="mx-2">›</span>
            <span class="text-purple-400 truncate max-w-md">{{ topic.title }}</span>
          </nav>
        </div>

        <!-- Original Post -->
        <div
          class="gaming-card rounded-xl p-6 mb-6 border border-purple-500/20 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] relative overflow-hidden">
          <!-- Animated border overlay -->
          <div
            class="absolute inset-0 bg-[linear-gradient(45deg,#8b5cf6,#ec4899,#06b6d4,#8b5cf6)] bg-[length:300%_300%] animate-[gradientBorder_3s_ease_infinite] opacity-10 pointer-events-none">
          </div>

          <div class="relative z-10">
            <div class="flex items-center justify-between gap-3 mb-4">
              <div class="flex items-center gap-3">
                <div class="relative group">
                  <Avatar :src="topic.user.avatar || '/images/default-avatar.png'"
                    :username="topic.user.nickname"
                    class="w-12 h-12 rounded border-2 border-purple-500/50 flex-shrink-0 object-cover" />
                </div>
                <div>
                  <div class="flex items-center gap-2 flex-wrap">
                    <nuxt-link :to="`/users/${topic.user.username}`"
                      class="font-bold text-purple-300 gaming-title text-lg hover:text-purple-400 transition-colors"
                      itemprop="author">
                      {{ topic.user.nickname }}
                    </nuxt-link>
                    <span v-if="topic.user.type === 1 || topic.user.id === 1"
                      class="px-2 py-0.5 bg-red-500/20 text-red-400 text-xs font-bold rounded">
                      ADMIN
                    </span>
                    <span class="px-2 py-0.5 bg-purple-500/20 text-purple-400 text-xs font-bold rounded">
                      LVL
                      {{ calculateLevel(topic.user.score) }}
                    </span>
                  </div>

                  <div class="text-sm text-gray-500 mt-1">
                    {{ usePrettyDate(topic.createTime) }}
                  </div>

                  <div class="flex items-center gap-2 text-sm text-gray-400 mt-1">
                    <span v-if="topic.ipLocation" class="flex items-center gap-1">
                      <FontAwesome :icon="['fas', 'map-marker-alt']" class="text-gray-500" /> {{ topic.ipLocation }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Manage Menu -->
              <TopicManageMenu v-model="topic" class="relative z-20" @onSwitchEditMode="onSwitchEditMode" />
            </div>

            <h1 class="text-3xl font-bold mb-6 text-white" itemprop="headline">
              <span v-if="topic.sticky" class="text-red-500 mr-2" title="Sticky Post">
                <FontAwesome :icon="['fas', 'thumbtack']" />
              </span>
              <span v-if="topic.recommend" class="text-yellow-500 mr-2" title="Recommended">
                <FontAwesome :icon="['fas', 'star']" />
              </span>
              {{ topic.title }}
            </h1>

            <div class="prose prose-invert max-w-none mb-6 text-gray-300" itemprop="articleBody">
              <topic-content v-model="topic" :editing="topic.editing" />
            </div>

            <!-- Tags -->
            <div v-if="topic.tags && topic.tags.length" class="flex flex-wrap gap-2 mb-6">
              <nuxt-link v-for="tag in topic.tags" :key="tag" :to="`/tags/${tag}`"
                class="px-3 py-1 bg-purple-500/10 text-purple-300 text-xs rounded-full border border-purple-500/20 hover:bg-purple-500/20 transition-colors">
                #{{ tag }}
              </nuxt-link>
            </div>

            <div class="flex items-center justify-between pt-4 border-t border-purple-500/20">
              <div class="flex items-center gap-4">
                <button
                  class="group flex items-center gap-1 transition-all duration-200"
                  :class="liked ? 'text-green-400' : 'text-gray-400 hover:text-green-400'"
                  @click="toggleLike">
                  <FontAwesome :icon="['fas', 'arrow-up']"
                    class="transform group-hover:-translate-y-0.5 transition-transform" />
                  <span class="font-bold">{{ topic.likeCount || 0 }}</span>
                </button>

                <!-- Downvote (Visual Only/Disabled as API doesn't support yet, generic placeholder logic) -->
                <!-- <button class="group flex items-center gap-1 text-gray-400 hover:text-red-400 transition-all duration-200">
                  <FontAwesome :icon="['fas', 'arrow-down']" class="transform group-hover:translate-y-0.5 transition-transform" />
                </button> -->

                <div class="flex items-center gap-1 text-gray-400">
                  <FontAwesome :icon="['far', 'eye']" />
                  <span>{{ topic.viewCount }} Views</span>
                </div>

                <div class="flex items-center gap-1 text-gray-400">
                  <FontAwesome :icon="['far', 'comment']" />
                  <span>{{ topic.commentCount }} Comments</span>
                </div>

                <button
                  class="flex items-center gap-1 transition-colors duration-200"
                  :class="topic.favorited ? 'text-yellow-400' : 'text-gray-400 hover:text-yellow-400'"
                  @click="toggleFavorite">
                  <FontAwesome :icon="topic.favorited ? ['fas', 'bookmark'] : ['far', 'bookmark']" />
                  <span>{{ topic.favorited ? 'Saved' : 'Save' }}</span>
                </button>
              </div>

              <div class="flex items-center gap-2">
                <button class="text-gray-400 hover:text-purple-400 transition-colors">
                  <FontAwesome :icon="['fas', 'share']" />
                </button>
                <button class="text-gray-400 hover:text-red-400 transition-colors">
                  <FontAwesome :icon="['fas', 'flag']" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Comments Section -->
        <div
          class="gaming-card rounded-xl p-6 border border-purple-500/20 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))]">
          <h2 class="text-2xl font-bold mb-6 flex items-center text-white">
            <FontAwesome :icon="['fas', 'comments']" class="mr-3 text-blue-400" />
            Comments ({{ topic.commentCount }})
          </h2>

          <!-- Wrapped Comment Component -->
          <CommentSection
            :topicSlug="topic.slug"
            :comment-count="topic.commentCount"
            entity-type="topic"
            @created="commentCreated" />
        </div>
      </div>
    </main>

    <GamingFooter />
  </div>
  <NotFound v-else />
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const userStore = useUserStore()
const api = useApi()

const slug = route.params.slug
const user = userStore.user

const { data: topic } = await useAsyncData('topic', () => api.getTopic(slug))

// Handle 404 manually if topic is null (handled by v-if/v-else logic currently, but strict check is better)
if (!topic.value) {
  // Falls through to NotFound component via v-else
}

definePageMeta({
  layout: false,
})

useHead({
  title: topic.value ? useTopicSiteTitle(topic.value) : i18n.t('page.not_found'),
  bodyAttrs: {
    class: 'bg-[#0f0f23]',
  },
})

const isPending = computed(() => {
  return topic.value?.status === 2
})

const liked = ref(topic.value?.liked || false)

function calculateLevel(score) {
  return Math.floor(Math.sqrt(score || 0)) + 1
}

async function toggleLike() {
  try {
    if (liked.value) {
      await api.removeTopicReaction(slug)
      liked.value = false
      topic.value.likeCount = topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    } else {
      await api.addTopicReaction(slug, 'like')
      liked.value = true
      topic.value.likeCount++
      useMsgSuccess(i18n.t('message.liked_success'))
    }
  } catch (e) {
    useCatchError(e)
  }
}

async function toggleFavorite() {
  try {
    if (topic.value.favorited) {
      await useHttpDelete(`/api/me/favorites/${topic.value.id}`)
      topic.value.favorited = false
      useMsgSuccess(i18n.t('message.removed_from_favorite'))
    } else {
      await useHttpPutForm(`/api/me/favorites`, {
        body: { topicId: topic.value.id },
      })
      topic.value.favorited = true
      useMsgSuccess(i18n.t('message.added_to_favorite'))
    }
  } catch (e) {
    useCatchError(e)
  }
}

async function onSwitchEditMode() {
  if (topic.value.editing) {
    useHttpPutForm(`/api/topics/${topic.value.slug}`, {
      body: {
        forumId: topic.value.forumId,
        title: topic.value.title,
        content: topic.value.content,
        tags: topic.value.tags,
        imageList: [],
      },
    }).then((tmp) => {
      topic.value.content = tmp.content
      topic.value.editing = false
      useLinkTo(`/topics/${tmp.slug}`)
    }).catch((err) => {
      alert(err)
    })
  } else {
    useHttpGet(`/api/topics/${topic.value.slug}/edit`).then((tmp) => {
      topic.value.content = tmp.content
      topic.value.editing = true
    }).catch((err) => {
      alert(err)
    })
  }
}

function commentCreated() {
  topic.value.commentCount++
}
</script>
