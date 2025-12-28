<template>
  <div class="space-y-6">
    <LoadMoreAsync ref="loadMoreComment" v-slot="{ items }" :cursor="commentsCursor">
      <TransitionGroup name="comment" tag="div" class="space-y-6">
        <div v-for="comment in items" :key="comment.id"
          class="border-b border-purple-500/20 pb-6 last:border-0 last:pb-0 mt-6 group">
          <div class="flex items-start gap-4">
            <div class="relative group">
              <Avatar :user="comment.user" :size="48" class="rounded border-2 border-purple-500/30" />
            </div>

            <div class="flex-1 min-w-0">
              <!-- Comment Meta -->
              <div class="flex items-center justify-between gap-2 flex-wrap mb-2">
                <div class="flex items-center gap-2 flex-wrap">
                  <nuxt-link :to="`/users/${comment.user.username}`"
                    class="font-bold text-purple-300 hover:text-purple-400 transition-colors">
                    {{ comment.user.username }}
                  </nuxt-link>
                  <span class="text-xs text-gray-500">{{ usePrettyDate(comment.createTime) }}</span>
                  <span v-if="comment.ipLocation" class="text-xs text-gray-600 bg-gray-800/50 px-2 py-0.5 rounded">
                    {{ comment.ipLocation }}
                  </span>
                </div>

                <button v-if="isCommentOwner(comment)"
                  class="group flex items-center gap-1 transition-colors text-gray-400 hover:text-purple-400 ml-auto opacity-0 group-hover:opacity-100 transition-opacity"
                  @click="deleteComment(comment)">
                  <Icon name="Fa7SolidTrashCan" class="w-4 h-4" />
                </button>
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
              <div class="flex items-center gap-4 text-sm text-gray-400 w-full">
                <button class="group flex items-center gap-1 transition-colors"
                  :class="comment.liked ? 'text-blue-400' : 'hover:text-blue-400'" @click="like(comment)">
                  <Icon name="MaterialSymbolsThumbUp" class="w-4 h-4" />
                  <span v-if="comment.likeCount > 0">{{ comment.likeCount }}</span>
                  <span v-else>Like</span>
                </button>

                <button class="group flex items-center gap-1 transition-colors"
                  :class="myReply.commentId === comment.id ? 'text-purple-400' : 'hover:text-gray-400'"
                  @click="switchShowReply(comment)">
                  <Icon name="TablerMessageCircle" class="w-4 h-4" />
                  <span>
                    {{ myReply.commentId === comment.id ? $t('feed.actions.hide_reply') : $t('feed.actions.reply') }}
                  </span>
                </button>
              </div>

              <!-- Reply Form -->
              <div v-if="myReply.commentId === comment.id"
                class="mt-4 p-4 rounded-lg bg-gray-800/50 border border-purple-500/20">
                <div v-if="replyTo" class="flex items-center gap-2 mb-4">
                  <span class="text-gray-400 font-medium text-xs tracking-wider">
                    {{ $t('feed.replying_to') }}
                  </span>
                  <nuxt-link :to="`/users/${replyTo.username}`"
                    class="font-bold text-sm text-purple-300 hover:text-purple-400 transition-colors">
                    {{ replyTo.username }}
                  </nuxt-link>
                </div>
                <CommentTextEditor v-model="myReply.input" :height="100" @submit="submitReply(comment)" />
              </div>

              <!-- Nested Replies -->
              <CommentReplyList v-if="comment.replies && comment.replies.items" v-model="myReply"
                :comment-id="comment.id"
                :data="comment.replies" @reply="onReply(comment, $event)" />
            </div>
          </div>
        </div>
      </TransitionGroup>
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, type Ref } from 'vue'
import type { Comment, CommentUser } from '@/types'

const userStore = useUserStore()
const api = useApi()
const dialog = useConfirmDialog()
const { t } = useI18n()

const props = defineProps({
  topicSlug: {
    type: String,
    required: true,
  },
})

const myReply = reactive({
  commentId: 0,
  input: {
    content: '',
    imageList: [],
  },
})

const replyTo = ref<CommentUser | null>(null)

const commentsCursor = api.getTopicComments(props.topicSlug)

const isLogin = computed(() => {
  return !!userStore.user
})

const isCommentOwner = (comment: Comment) => {
  return userStore.user && userStore.user.id === comment.user.id
}

const loadMoreComment = ref<any>(null)

const append = (comment: Comment) => {
  if (!loadMoreComment.value || !comment) return
  loadMoreComment.value.unshiftResults(comment)
}

const like = async (comment: Comment) => {
  if (!isLogin.value) return
  try {
    if (comment.liked) {
      await api.removeCommentReaction(comment.id)
      comment.liked = false
      comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0
    }
    else {
      await api.addCommentReaction(comment.id, 'like')
      comment.liked = true
      comment.likeCount = (comment.likeCount || 0) + 1
    }
  }
  catch (e) {
    useCatchError(e)
  }
}

const switchShowReply = (comment: Comment) => {
  if (!isLogin.value) return

  if (myReply.commentId === comment.id) {
    hideReply()
  }
  else {
    myReply.commentId = comment.id
    replyTo.value = comment.user
  }
}

const hideReply = () => {
  myReply.commentId = 0
  myReply.input.content = ''
  myReply.input.imageList = []
  replyTo.value = null
}

const submitReply = async (parent: Comment) => {
  if (!isLogin.value) return
  try {
    const ret = await api.addCommentReply(parent.id, {
      content: myReply.input.content || '',
      imageList: myReply.input.imageList,
    })
    hideReply()
    prependReply(parent, ret)
  }
  catch (e) {
    useCatchError(e)
  }
}

const onReply = (parent: Comment, comment: Comment) => {
  prependReply(parent, comment)
}

const prependReply = (parent: Comment, comment: Comment) => {
  if (parent.replies && parent.replies.items) {
    parent.replies.items.unshift(comment)
  }
  else {
    parent.replies = {
      items: [comment],
    }
  }
}

const deleteComment = async (comment: Comment) => {
  dialog.show({
    title: $t('dialog.title.confirm_delete'),
    message: $t('dialog.message.confirm_delete_comment'),
    confirmText: $t('dialog.button.confirm'),
    cancelText: $t('dialog.button.cancel'),
    variant: 'warning',
    icon: 'Fa7SolidTrashCan',
    onConfirm: async () => {
      await api.deleteComment(comment.id)
      if (loadMoreComment.value) {
        loadMoreComment.value.removeItem(comment)
      }
    },
  })
}

defineExpose({
  append,
})
</script>

<style scoped>
.comment-enter-active,
.comment-leave-active,
.comment-move {
  transition: all 220ms ease;
}

.comment-enter-from,
.comment-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.comment-enter-to,
.comment-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
