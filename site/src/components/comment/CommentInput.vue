<template>
  <div class="mb-6">
    <div
      class="bg-gray-800/40 rounded-xl border border-purple-500/20 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/40 focus-within:border-purple-500/60 focus-within:bg-gray-800/60 focus-within:shadow-[0_0_15px_rgba(168,85,247,0.1)] overflow-hidden group">
      <div ref="commentEditor" class="relative">
        <div v-if="quote"
          class="flex items-center justify-between gap-2 text-sm mx-3 mt-3 mb-1 bg-purple-500/10 p-2 px-3 rounded-lg border border-purple-500/20 animate-in fade-in slide-in-from-top-1">
          <div class="flex items-center gap-2">
            <span class="text-purple-400 font-medium text-xs uppercase tracking-wider">
              {{ $t('feed.replied_to') }}
            </span>
            <span class="text-gray-200 font-bold">{{ quote.user.nickname }}</span>
          </div>
          <button class="text-gray-500 hover:text-red-400 transition-colors p-1 rounded-full hover:bg-white/5"
            @click="cancelReply">
            <FontAwesome :icon="['fas', 'times']" />
          </button>
        </div>

        <CommentTextEditor ref="simpleEditor" v-model="value"
          :placeholder="$t('message.type_comment_placeholder')" @submit="create" />

        <!-- Loading Overlay -->
        <div v-if="sending"
          class="absolute inset-0 bg-gray-900/50 backdrop-blur-[1px] flex items-center justify-center z-10 rounded-xl">
          <div class="flex items-center gap-2 text-purple-400 font-bold gaming-title animate-pulse">
            <FontAwesome :icon="['fas', 'circle-notch']" spin /> SENDING...
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()

const props = defineProps({
  entityType: {
    type: String,
    required: true,
  },
  entityId: {
    type: Number,
    required: true,
  },
})

const emits = defineEmits(['created'])
const value = ref({
  content: '', // 内容
  imageList: [],
})
const sending = ref(false) // 发送中
const quote = ref(null) // 引用的对象
const commentEditor = ref(null) // 编辑器组件
const simpleEditor = ref(null) // 编辑器组件

async function create() {
  if (!value.value.content) {
    useMsgError(i18n.t('message.please_enter_comment'))
    return
  }
  if (sending.value) {
    return
  }
  if (simpleEditor.value && simpleEditor.value.isOnUpload()) {
    return
  }
  sending.value = true
  try {
    const data = await useHttpPostForm(`/api/topics/${route.params.slug}/comments`, {
      body: {
        content: value.value.content,
        imageList: value.value.imageList && value.value.imageList.length
          ? JSON.stringify(value.value.imageList)
          : '',
        quoteId: quote.value ? quote.value.id : '',
      },
    })
    emits('created', data)

    value.value.content = ''
    value.value.imageList = []
    quote.value = null
    simpleEditor.value.clear()
    useMsgSuccess(i18n.t('message.comment_success'))
  } catch (e) {
    console.error(e)
    useMsgError(e.message || e)
  } finally {
    sending.value = false
  }
}
function cancelReply() {
  quote.value = null
}
</script>

<style scoped lang="scss"></style>
