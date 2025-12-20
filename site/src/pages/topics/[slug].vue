<template>
  <section v-if="topic != null" class="main">
    <div v-if="isPending" class="container main-container">
      <div class="notification is-warning"
        style="width: 100%; margin: 20px 0">
        {{ $t('message.post_under_review') }}
      </div>
    </div>
    <div class="container main-container left-main size-300">
      <div class="left-container">
        <article class="topic-detail" itemscope itemtype="http://schema.org/BlogPosting">
          <div class="topic-header">
            <div class="topic-header-left">
              <my-avatar :user="topic.user" :size="45" />
            </div>

            <div class="topic-header-center">
              <div class="topic-nickname" itemprop="headline">
                <nuxt-link
                  itemprop="author"
                  itemscope
                  itemtype="http://schema.org/Person"
                  :to="`/user/${topic.user.id}`">
                  {{ topic.user.nickname }}
                </nuxt-link>
              </div>
              <div class="topic-meta">
                <span class="meta-item">
                  {{ $t('feed.published_on') }}
                  <time
                    :datetime="usePrettyDate(topic.createTime)"
                    itemprop="datePublished">
                    {{ usePrettyDate(topic.createTime) }}
                  </time>
                </span>
                <span v-if="topic.ipLocation" class="meta-item">{{ topic.ipLocation }}</span>
              </div>
            </div>

            <div class="topic-header-right">
              <topic-manage-menu v-model="topic" @onSwitchEditMode="onSwitchEditMode" />
            </div>
          </div>

          <div class="topic-content content" itemprop="articleBody">
            <topic-content v-model="topic" :editing="topic.editing" />
          </div>

          <div class="topic-tags">
            <nuxt-link
              v-if="topic.forum"
              :to="`/forums/${topic.forum.slug}`"
              class="topic-tag">
              {{ topic.forum.name }}
            </nuxt-link>
            <nuxt-link
              v-for="tag in topic.tags"
              :key="tag.name"
              :to="`/tags/${tag}`"
              class="topic-tag">
              #{{ tag }}
            </nuxt-link>
          </div>

          <div class="topic-actions">
            <div class="action disabled">
              <icon name="BookOpenText" size="1em" />
              <div class="action-text">
                <span>{{ $t('feed.view_count') }}</span>
                <span v-if="topic.viewCount > 0" class="action-text">
                  ({{ topic.viewCount }})
                </span>
              </div>
            </div>
            <div class="action" @click="toggleLike(topic)">
              <icon name="ThumbsUp" color="#1c71d8" :filled="liked" />
              <div class="action-text">
                <span>{{ $t('feed.like_count') }}</span>
                <span v-if="topic.likeCount > 0">
                  ({{ topic.likeCount }})
                </span>
              </div>
            </div>
            <div class="action" @click="toggleFavorite(topic.id)">
              <icon name="Star" color="#f6d32d" :filled="topic.favorited" />
              <div class="action-text">
                <span>{{ $t('feed.actions.favorite') }}</span>
              </div>
            </div>
          </div>
        </article>

        <!-- 评论 -->
        <comment
          :entity-id="topic.id"
          :comment-count="topic.commentCount"
          entity-type="topic"
          @created="commentCreated" />
      </div>
      <div class="right-container">
        <user-info :user="topic.user" />
      </div>
    </div>
  </section>
  <not-found v-else />
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const userStore = useUserStore()

const slug = route.params.slug
const user = userStore.user

const { data: topic } = await useAsyncData('topic', () =>
  useHttpGet(`/api/topics/${slug}`),
)

useHead({
  title: topic.value ? useTopicSiteTitle(topic.value) : i18n.t('page.not_found'),
})
const isPending = computed(() => {
  return topic.value.status === 2
})

const liked = ref(topic.value?.liked || false)
async function toggleLike() {
  try {
    if (liked.value) {
      await useHttpDelete(`/api/topics/${slug}/reactions/${user.id}`)
      liked.value = false
      topic.value.likeCount = topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    } else {
      await useHttpPostForm(`/api/topics/${slug}/reactions`, {
        body: { type: 'like' },
      })
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

<style lang="scss" scoped></style>
