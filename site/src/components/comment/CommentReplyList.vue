<template>
  <div class="mt-4 pl-4 border-l-4 border-purple-500/20 space-y-4">
    <div v-for="reply in replies.items" :key="reply.id" class="flex items-start gap-4">
      <div class="relative group">
        <Avatar :src="reply.user.avatar" :username="reply.user.username" :size="48"
          class="rounded border-2 border-purple-300" />
      </div>

      <div class="flex-1 min-w-0">
        <!-- Reply Meta -->
        <div class="flex items-center gap-2 flex-wrap mb-1">
          <nuxt-link :to="`/users/${reply.user.username}`"
            class="font-bold text-purple-300 hover:text-purple-400 transition-colors">
            {{ reply.user.username }}
          </nuxt-link>

          <template v-if="reply.quote">
            <span class="text-xs text-gray-500">{{ $t('feed.replied_to') }}</span>
            <nuxt-link :to="`/users/${reply.quote.user.username}`"
              class="font-bold text-purple-300 hover:text-purple-400 transition-colors">
              {{ reply.quote.user.username }}
            </nuxt-link>
          </template>

          <span class="text-xs text-gray-500">{{ usePrettyDate(reply.createTime) }}</span>
        </div>

        <!-- Content -->
        <div class="text-gray-300 mb-2 prose prose-invert max-w-none text-sm">
          <div v-text="reply.content" />
        </div>

        <!-- Images -->
        <div v-if="reply.imageList && reply.imageList.length" class="flex flex-wrap gap-2 mb-2">
          <img v-for="(image, imageIndex) in reply.imageList" :key="imageIndex" :src="image.url"
            class="w-16 h-16 object-cover rounded border border-purple-500/30 hover:scale-105 transition-transform cursor-pointer">
        </div>

        <!-- Quoted Content -->
        <div v-if="reply.quote"
          class="mt-2 mb-2 p-3 bg-gray-800/50 rounded-lg border-l-4 border-purple-500/40 text-sm text-gray-400 italic">
          <div class="line-clamp-3" v-html="reply.quote.content" />
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4 text-sm text-gray-400">
          <button class="group flex items-center gap-1 transition-colors"
            :class="reply.liked ? 'text-blue-400' : 'hover:text-blue-400'" @click="like(reply)">
            <Icon name="MaterialSymbolsThumbUp" class="w-4 h-4" />
            <span v-if="reply.likeCount > 0">{{ reply.likeCount }}</span>
            <span v-else>Like</span>
          </button>

          <button class="group flex items-center gap-1 transition-colors"
            :class="myReply.commentId === reply.id ? 'text-purple-400' : 'hover:text-purple-400'"
            @click="switchShowReply(reply)">
            <Icon name="TablerMessageCircle" class="w-4 h-4" />
            <span>{{ myReply.commentId === reply.id ? $t('feed.actions.hide_reply') : $t('feed.actions.reply') }}</span>
          </button>
        </div>

        <!-- Reply Form -->
        <div v-if="myReply.commentId === reply.id"
          class="mt-3 p-3 rounded-lg bg-gray-800/50 border border-purple-500/20">
          <div class="flex items-center gap-2 mb-4">
            <span class="text-gray-400 font-medium text-xs tracking-wider">
              {{ $t('feed.replying_to') }}
            </span>
            <nuxt-link :to="`/users/${reply.user.username}`"
              class="font-bold text-purple-300 hover:text-purple-400 transition-colors">
              {{ reply.user.username }}
            </nuxt-link>
          </div>
          <CommentTextEditor v-model="myReply.input" :disabled="sending" :height="80" @submit="submitReply(reply)" />
        </div>
      </div>
    </div>

    <div v-if="replies.hasMore" class="pt-2">
      <button class="text-xs font-bold text-purple-400 hover:text-purple-300 transition-colors flex items-center gap-1"
        @click="loadMore">
        <Icon name="TablerChevronDown" /> {{ $t('feed.actions.view_more_replies') }}
      </button>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
const api = useApi()
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

const sending = ref(false)

const emit = defineEmits(['update:modelValue', 'reply'])
const myReply = reactive(props.modelValue)
const replies = ref(props.data)

const { user } = storeToRefs(userStore)

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
      await api.removeCommentReaction(comment.id)
      comment.liked = false
      comment.likeCount = comment.likeCount > 0 ? comment.likeCount - 1 : 0
      useMsgSuccess(i18n.t('message.unliked_success'))
    }
    else {
      await api.addCommentReaction(comment.id, 'like')
      comment.liked = true
      comment.likeCount = comment.likeCount + 1
      useMsgSuccess(i18n.t('message.liked_success'))
    }
  }
  catch (e) {
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
  }
  else {
    myReply.commentId = comment.id
  }
}

function hideReply() {
  myReply.commentId = 0
  myReply.input.content = ''
  myReply.input.imageList = []
}

async function submitReply(parent) {
  if (!myReply.input.content) {
    useMsgError(i18n.t('message.please_enter_comment'))
    return
  }
  if (sending.value) {
    return
  }

  sending.value = true
  try {
    const ret = await api.addCommentReply(myReply.commentId, {
      content: myReply.input.content || '',
      imageList: myReply.input.imageList,
      quoteId: myReply.commentId,
    })
    hideReply()
    emit('reply', ret)
    useMsgSuccess(i18n.t('message.comment_success'))
  }
  catch (e) {
    useCatchError(e)
  }
  finally {
    sending.value = false
  }
}
</script>
