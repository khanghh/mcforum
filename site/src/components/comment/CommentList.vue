<template>
  <div class="space-y-6">
    <LoadMoreAsync ref="loadMore" v-slot="{ items }" :url="commentListUrl">
      <div v-for="comment in items" :key="comment.id"
        class="border-b border-purple-500/20 pb-6 last:border-0 last:pb-0 mt-6">
        <div class="flex items-start gap-4">
          <div class="relative group">
            <MyAvatar :user="comment.user" :size="30" class="rounded-lg border-2 border-purple-500/30" />
          </div>

          <div class="flex-1 min-w-0">
            <!-- Comment Meta -->
            <div class="flex items-center gap-2 flex-wrap mb-2">
              <nuxt-link :to="`/user/${comment.user.id}`"
                class="font-bold text-purple-300 hover:text-purple-400 transition-colors">
                {{ comment.user.nickname }}
              </nuxt-link>
              <span class="text-xs text-gray-500">{{ usePrettyDate(comment.createTime) }}</span>
              <span v-if="comment.ipLocation" class="text-xs text-gray-600 bg-gray-800/50 px-2 py-0.5 rounded">
                {{ comment.ipLocation }}
              </span>
            </div>

            <!-- Content -->
            <div class="text-gray-300 mb-3 prose prose-invert max-w-none text-sm leading-relaxed">
              <template v-if="comment.content">
                <div v-if="comment.contentType === 'text'" v-text="comment.content" />
                <div v-else v-html="comment.content" />
              </template>
            </div>

            <!-- Images -->
            <div v-if="comment.imageList && comment.imageList.length" class="flex flex-wrap gap-2 mb-3">
              <img v-for="(image, imageIndex) in comment.imageList" :key="imageIndex" :src="image.url"
                class="w-24 h-24 object-cover rounded-lg border border-purple-500/30 hover:scale-105 transition-transform cursor-pointer">
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-4 text-sm text-gray-400">
              <button class="group flex items-center gap-1 transition-colors"
                :class="comment.liked ? 'text-green-400' : 'hover:text-green-400'" @click="like(comment)">
                <icon name="ThumbsUp" class="w-4 h-4" :filled="comment.liked" />
                <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
                <span v-else>Like</span>
              </button>

              <button class="group flex items-center gap-1 transition-colors"
                :class="myReply.commentId === comment.id ? 'text-purple-400' : 'hover:text-purple-400'"
                @click="switchShowReply(comment)">
                <icon name="MessageSquareMore" class="w-4 h-4" />
                <span>
                  {{ myReply.commentId === comment.id ? $t('feed.actions.hide_reply') : $t('feed.actions.reply') }}
                </span>
              </button>
            </div>

            <!-- Reply Form -->
            <div v-if="myReply.commentId === comment.id"
              class="mt-4 p-4 rounded-lg bg-gray-800/50 border border-purple-500/20">
              <CommentTextEditor v-model="myReply.input" :height="100" @submit="submitReply(comment)" />
            </div>

            <!-- Nested Replies -->
            <CommentReplyList v-if="comment.replies && comment.replies.items" v-model="myReply" :comment-id="comment.id"
              :data="comment.replies" @reply="onReply(comment, $event)" />
          </div>
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

<style scoped>
/* Scoped styles replaced by Tailwind classes */
</style>
