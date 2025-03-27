<template>
  <div class="replies">
    <div v-for="reply in replies.items" :key="reply.id" class="comment">
      <div class="comment-item-left">
        <my-avatar :user="reply.user" :size="30" has-border />
      </div>
      <div class="comment-item-main">
        <div class="comment-meta">
          <div>
            <nuxt-link :to="`/user/${reply.user.id}`" class="comment-nickname">
              {{ reply.user.nickname }}
            </nuxt-link>
            <template v-if="reply.quote">
              <span>&nbsp;{{ $t('feed.replied_to') }}&nbsp;</span>
              <nuxt-link :to="`/user/${reply.quote.user.id}`" class="comment-nickname">
                {{ reply.quote.user.nickname }}
              </nuxt-link>
            </template>
          </div>
          <time class="comment-time">{{ usePrettyDate(reply.createTime) }}</time>
        </div>
        <div class="comment-content-wrapper">
          <div v-if="reply.content" class="comment-content content">
            <div v-text="reply.content" />
          </div>
          <div v-if="reply.imageList && reply.imageList.length" class="comment-image-list">
            <img v-for="(image, imageIndex) in reply.imageList" :key="imageIndex" :src="image.url">
          </div>

          <div v-if="reply.quote" class="comment-quote">
            <div class="comment-quote-content content" v-html="reply.quote.content" />
            <div v-if="reply.quote.imageList && reply.quote.imageList.length" class="comment-quote-image-list">
              <img v-for="(image, imageIndex) in reply.imageList" :key="imageIndex" :src="image.url">
            </div>
          </div>
        </div>
        <div class="comment-actions">
          <div class="comment-action-item" :class="{ active: reply.liked }" @click="like(reply)">
            <icon name="ThumbsUp" :filled="reply.liked" />
            <span>&nbsp;{{ reply.liked ? $t('feed.actions.liked') : $t('feed.actions.like') }}&nbsp;</span>
            <span v-if="reply.likeCount > 0">{{ reply.likeCount }}</span>
          </div>
          <div class="comment-action-item" :class="{ active: myReply.commentId === reply.id }"
            @click="switchShowReply(reply)">
            <icon name="MessageSquareMore" />
            <span>
              &nbsp;
              {{ myReply.commentId === reply.id ? $t('feed.actions.hide_reply') : $t('feed.actions.reply') }}
            </span>
          </div>
        </div>
        <div v-if="myReply.commentId === reply.id" class="comment-reply-form">
          <TextEditor v-model="myReply.input" :height="80" @submit="submitReply()" />
        </div>
      </div>
    </div>
    <div v-if="replies.hasMore === true" class="comment-more">
      <a @click="loadMore">{{ $t('feed.actions.view_more_replies') }}</a>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
const userStore = useUserStore()

const props = defineProps({
  commentId: {
    type: Number,
    required: true,
  },
  modelValue: {
    type: Object,
    required: true,
  },
  data: {
    type: Object,
    required: true,
  },
})

const emit = defineEmits(['update:modelValue', 'reply'])
const myReply = reactive(props.modelValue)
const replies = ref(props.data)

const user = computed(() => {
  const userStore = useUserStore()
  return userStore.user
})

async function loadMore() {
  const ret = await useHttpGet(`/api/comments/${props.commentId}/replies`, {
    params: {
      cursor: replies.value.cursor,
    },
  })
  replies.value.cursor = ret.cursor
  replies.value.hasMore = ret.hasMore
  replies.value.items.push(...ret.items)
}

async function like(comment) {
  try {
    if (comment.liked) {
      await useHttpDelete(`/api/comments/${comment.id}/reactions/${userStore.user.id}`)
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

function switchShowReply(comment) {
  if (!user.value) {
    useMsgSignIn()
    return
  }

  if (myReply.commentId === comment.id) {
    hideReply(comment)
  } else {
    myReply.commentId = comment.id
  }
}

function hideReply() {
  myReply.commentId = 0
  myReply.input.content = ''
  myReply.input.imageList = []
}

async function submitReply(/* parent */) {
  try {
    const ret = await useHttpPostForm(`/api/comments/${myReply.commentId}/replies`, {
      body: {
        content: myReply.input.content || '',
        imageList: myReply.input.imageList ? JSON.stringify(myReply.input.imageList) : '',
      },
    })
    hideReply()
    emit('reply', ret)
    useMsgSuccess(i18n.t('message.comment_success'))
  } catch (e) {
    useCatchError(e)
  }
}
</script>

<style lang="scss" scoped>
.replies {
  margin-top: 10px;
  padding: 1px 10px;
  font-size: 12px;
  background-color: var(--bg-color2);

  .comment {
    display: flex;
    padding: 8px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .comment-item-main {
      flex: 1 1 auto;
      margin-left: 8px;

      .comment-meta {
        display: flex;
        justify-content: space-between;

        .comment-nickname {
          font-size: 12px;
          font-weight: 600;
          color: var(--text-color);

          &:hover {
            color: var(--text-link-color);
          }
        }

        .comment-time {
          font-size: 11px;
          color: var(--text-color3);
        }
      }

      .comment-content-wrapper {
        .comment-content {
          margin-top: 5px;
          margin-bottom: 0;
          color: var(--text-color2);
          white-space: pre-wrap;
        }

        .comment-image-list {
          margin-top: 5px;

          img {
            width: 62px;
            height: 62px;
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

        .comment-quote {
          position: relative;
          background-color: var(--bg-color3);
          border: 1px solid var(--border-color2);
          color: var(--text-color3);
          padding: 0 12px;
          margin: 5px 0;
          box-sizing: border-box;
          border-radius: 4px;

          &::after {
            position: absolute;
            content: "\201D";
            font-family: Georgia, serif;
            font-size: 36px;
            font-weight: bold;
            color: var(--text-color3);
            right: 2px;
            top: -8px;
          }

          .comment-quote-content {
            margin: 5px 0;
            color: var(--text-color3);
          }

          .comment-quote-image-list {
            margin-top: 5px;

            img {
              width: 50px;
              height: 50px;
              line-height: 50px;
              cursor: pointer;

              &:not(:last-child) {
                margin-right: 4px;
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
      }

      .comment-actions {
        margin-top: 5px;
        display: flex;
        align-items: center;

        .comment-action-item {
          display: inline-flex;
          line-height: 22px;
          font-size: 11px;
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
    }
  }

  .reply {
    display: flex;
  }

  .comment-more {
    margin: 10px;
    font-size: 13px;
    font-weight: 500;
    color: var(--text-link-color);
  }
}
</style>
