<template>
  <div class="comments">
    <LoadMoreAsync ref="loadMore" v-slot="{ items }" :url="commentListUrl">
      <div v-for="comment in items" :key="comment.id" class="comment">
        <div class="comment-item-left">
          <MyAvatar :user="comment.user" :size="40" has-border />
        </div>
        <div class="comment-item-main">
          <div class="comment-meta">
            <nuxt-link :to="`/user/${comment.user.id}`" class="comment-nickname">
              {{ comment.user.nickname }}
            </nuxt-link>
            <div class="comment-meta-right">
              <time class="comment-time">
                {{ usePrettyDate(comment.createTime) }}
              </time>
              <span v-if="comment.ipLocation" class="comment-ip-area">{{ comment.ipLocation }}</span>
            </div>
          </div>
          <div class="comment-content-wrapper">
            <template v-if="comment.content">
              <div v-if="comment.contentType === 'text'" class="comment-content content" v-text="comment.content" />
              <div v-else class="comment-content content" v-html="comment.content" />
            </template>
            <div v-if="comment.imageList && comment.imageList.length" class="comment-image-list">
              <img v-for="(image, imageIndex) in comment.imageList" :key="imageIndex" :src="image.url">
            </div>
          </div>
          <div class="comment-actions">
            <div class="comment-action-item" :class="{ active: comment.liked }" @click="like(comment)">
              <icon name="ThumbsUp" :filled="comment.liked" />
              <span>&nbsp;{{ comment.liked ? $t('feed.actions.liked') : $t('feed.actions.like') }}&nbsp;</span>
              <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
            </div>
            <div class="comment-action-item" :class="{ active: myReply.commentId === comment.id }"
              @click="switchShowReply(comment)">
              <icon name="MessageSquareMore" />
              <span>
                &nbsp;
                {{ myReply.commentId === comment.id ? $t('feed.actions.hide_reply') : $t('feed.actions.reply') }}
              </span>
            </div>
          </div>
          <div v-if="myReply.commentId === comment.id" class="comment-reply-form">
            <TextEditor v-model="myReply.input" :height="100" @submit="submitReply(comment)" />
          </div>
          <CommentReplyList v-if="comment.replies && comment.replies.items" v-model="myReply" :comment-id="comment.id"
            :data="comment.replies" @reply="onReply(comment, $event)" />
        </div>
      </div>
    </LoadMoreAsync>
  </div>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()

defineProps({
  entityType: {
    type: String,
    required: true,
  },
  entityId: {
    type: Number,
    required: true,
  },
})

const myReply = reactive({
  commentId: 0,
  quoteId: 0,
  input: {
    content: '',
    imageList: [],
  },
})

const commentListUrl = computed(() => `/api/topics/${route.params.slug}/comments`)

const userStore = useUserStore()
const loadMore = ref(null)

const append = () => {
  if (loadMore.value) {
    // console.log(loadMore.value);
    // console.log(loadMore.value.unshiftResults);
    // loadMore.value.unshiftResults(data);
    loadMore.value.refresh()
  }
}

const like = async (comment) => {
  try {
    if (comment.liked) {
      await useHttpDelete(`/api/comments/${comment.id}/reactions/${userStore.user.id}`, {
      })
      comment.liked = false
      comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    } else {
      await useHttpPostForm(`/api/comments/${comment.id}/reactions`, {
        body: { type: 'like' },
      })
      comment.liked = true
      comment.likeCount = comment.likeCount + 1
      useMsgSuccess(i18n.t('message.liked_success'))
    }
  } catch (e) {
    useCatchError(e)
  }
}

const switchShowReply = (comment) => {
  if (!userStore.user) {
    useMsgSignIn()
    return
  }

  if (myReply.commentId === comment.id) {
    hideReply()
  } else {
    myReply.commentId = comment.id
  }
}

const hideReply = () => {
  myReply.commentId = 0
  myReply.input.content = ''
  myReply.input.imageList = []
}

const submitReply = async (parent) => {
  try {
    const ret = await useHttpPostForm(`/api/comments/${parent.id}/replies`, {
      body: {
        content: myReply.input.content,
        imageList: myReply.input.imageList ? JSON.stringify(myReply.input.imageList) : '',
      },
    })
    hideReply()
    prependReply(parent, ret)
    useMsgSuccess(i18n.t('message.comment_success'))
  } catch (e) {
    useCatchError(e)
  }
}

const onReply = (parent, comment) => {
  prependReply(parent, comment)
}

const prependReply = (parent, comment) => {
  if (parent.replies && parent.replies.items) {
    parent.replies.items.unshift(comment)
  } else {
    parent.replies = {
      items: [comment],
    }
  }
}

defineExpose({
  append,
})
</script>

<style scoped lang="scss">
.comments {
  padding: 10px;
  font-size: 14px;

  .comment {
    display: flex;
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 16px;

      .comment-meta {
        display: flex;
        justify-content: space-between;

        .comment-nickname {
          font-size: 14px;
          font-weight: 600;
          color: var(--text-color);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-meta-right {
          .comment-time {
            font-size: 13px;
            color: var(--text-color3);
          }

          .comment-ip-area {
            font-size: 13px;
            color: var(--text-color3);
            margin-left: 10px;
          }
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 10px;
          margin-bottom: 0;
          color: var(--text-color);
          white-space: pre-wrap;
        }

        .comment-image-list {
          margin-top: 10px;

          img {
            width: 72px;
            height: 72px;
            line-height: 72px;
            cursor: pointer;

            &:not(:last-child) {
              margin-right: 8px;
            }

            object-fit: cover;
            transition: all 0.5s ease-out 0.1s;

            &:hover {
              transform: matrix(1.04, 0, 0, 1.04, 0, 0);
              backface-visibility: hidden;
            }
          }
        }
      }

      .comment-actions {
        margin-top: 10px;
        display: flex;
        align-items: center;

        .comment-action-item {
          display: inline-flex;
          line-height: 22px;
          font-size: 13px;
          cursor: pointer;
          color: var(--text-color3);
          user-select: none;

          &:hover {
            color: var(--text-link-color);
          }

          &.active {
            color: var(--text-link-color);
            font-weight: 500;
          }

          &:not(:last-child) {
            margin-right: 16px;
          }
        }
      }

      .comment-reply-form {
        margin-top: 10px;
      }

      .comment-replies {
        margin-top: 10px;
        // padding: 10px;
        background-color: var(--bg-color2);
      }
    }
  }

  .reply {
    display: flex;
  }
}
</style>
