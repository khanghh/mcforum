<template>
  <div>
    <section class="main">
      <div v-if="isPending" class="container main-container">
        <div
          class="notification is-warning"
          style="width: 100%; margin: 20px 0">
          {{ $t('message.post_under_review') }}
        </div>
      </div>
      <div class="container main-container left-main size-360">
        <div class="left-container">
          <div class="main-content no-padding no-bg">
            <article
              class="topic-detail"
              itemscope
              itemtype="http://schema.org/BlogPosting">
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
                  <topic-manage-menu v-model="topic" />
                </div>
              </div>

              <!-- 内容 -->
              <div
                class="topic-content content"
                :class="{
                  'topic-tweet': topic.type === 1,
                }"
                itemprop="articleBody">
                <h1 v-if="topic.title" class="topic-title" itemprop="headline">
                  {{ topic.title }}
                </h1>
                <div class="topic-content-detail line-numbers" v-html="topic.content" />
                <ul v-if="topic.imageList && topic.imageList.length" class="topic-image-list">
                  <li v-for="(image, index) in topic.imageList" :key="index">
                    <div class="image-item">
                      <el-image
                        :src="image.preview"
                        :preview-src-list="imageUrls"
                        :initial-index="index" />
                    </div>
                  </li>
                </ul>

                <div v-if="hideContent && hideContent.exists" class="topic-content-detail hide-content">
                  <div v-if="hideContent.show" class="widget has-border">
                    <div class="widget-header">
                      <span>
                        <icon name="LockOpen" />
                        <span>{{ $t('publish.hidden_content_unlocked') }}</span>
                      </span>
                    </div>
                    <div class="widget-content" v-html="hideContent.content" />
                  </div>
                  <div v-else class="hide-content-tip">
                    <icon name="Lock" />
                    <span>{{ $t('publish.hidden_content_locked') }}</span>
                  </div>
                </div>
              </div>

              <!-- 节点、标签 -->
              <div class="topic-tags">
                <nuxt-link
                  v-if="topic.node"
                  :to="`/topics/node/${topic.node.id}`"
                  class="topic-tag">
                  {{ topic.node.name }}
                </nuxt-link>
                <nuxt-link
                  v-for="tag in topic.tags"
                  :key="tag.id"
                  :to="`/topics/tag/${tag.id}`"
                  class="topic-tag">
                  #{{ tag.name }}
                </nuxt-link>
              </div>

              <!-- 点赞用户列表 -->
              <div
                v-if="likeUsers && likeUsers.length"
                class="topic-like-users">
                <my-avatar
                  v-for="likeUser in likeUsers"
                  :key="likeUser.id"
                  :user="likeUser"
                  :size="24"
                  has-border />
                <span class="like-count">{{ topic.likeCount }}</span>
              </div>

              <!-- 功能按钮 -->
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
                <div class="action" @click="like(topic)">
                  <icon name="ThumbsUp" color="#1c71d8" :filled="liked" />
                  <div class="action-text">
                    <span>{{ $t('feed.like_count') }}</span>
                    <span v-if="topic.likeCount > 0">
                      ({{ topic.likeCount }})
                    </span>
                  </div>
                </div>
                <div class="action" @click="addFavorite(topic.id)">
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
        </div>
        <div class="right-container">
          <user-info :user="topic.user" />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const hideContent = ref(null)

const { data: topic } = await useAsyncData('topic', () =>
  useMyFetch(`/api/topic/${route.params.id}`),
)

const { data: liked } = await useAsyncData('liked', () => {
  return useMyFetch('/api/like/liked', {
    params: {
      entityType: 'topic',
      entityId: route.params.id,
    },
  })
})

const { data: likeUsers, refresh: refreshLikeUsers } = await useAsyncData(
  () => {
    return useMyFetch(`/api/topic/recentlikes/${route.params.id}`)
  },
)

const imageUrls = computed(() => {
  if (!topic.value.imageList || !topic.value.imageList.length) {
    return []
  }
  const ret = []
  for (let i = 0; i < topic.value.imageList.length; i++) {
    ret.push(topic.value.imageList[i].url)
  }
  return ret
})

useHead({
  title: useTopicSiteTitle(topic.value),
})

const isPending = computed(() => {
  return topic.value.status === 2
})

async function like() {
  try {
    if (liked.value) {
      await useHttpPostForm('/api/like/unlike', {
        body: {
          entityType: 'topic',
          entityId: topic.value.id,
        },
      })
      liked.value = false
      topic.value.likeCount
        = topic.value.likeCount > 0 ? topic.value.likeCount - 1 : 0

      useMsgSuccess(i18n.t('message.unliked_success'))
      await refreshLikeUsers()
    }
    else {
      await useHttpPostForm('/api/like/like', {
        body: {
          entityType: 'topic',
          entityId: topic.value.id,
        },
      })
      liked.value = true
      topic.value.likeCount++

      useMsgSuccess(i18n.t('message.liked_success'))
      await refreshLikeUsers()
    }
  }
  catch (e) {
    useCatchError(e)
  }
}

async function addFavorite(topicId) {
  try {
    if (topic.value.favorited) {
      await useHttpPostForm('/api/favorite/delete', {
        body: {
          entityType: 'topic',
          entityId: topicId,
        },
      })
      topic.value.favorited = false
      useMsgSuccess(i18n.t('message.removed_from_favorite'))
    }
    else {
      await useHttpPostForm('/api/favorite/add', {
        body: {
          entityType: 'topic',
          entityId: topicId,
        },
      })
      topic.value.favorited = true
      useMsgSuccess(i18n.t('message.added_to_favorite'))
    }
  }
  catch (e) {
    useCatchError(e)
  }
}

async function commentCreated() {
  console.log('commentCreated...')
}
</script>

<style lang="scss" scoped></style>
