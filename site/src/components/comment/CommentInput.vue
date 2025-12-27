<template>
  <div class="mb-6">
    <div
      class="bg-gray-800/40 rounded-xl border border-purple-500/20 backdrop-blur-sm transition-all duration-300 hover:border-purple-500/40 focus-within:border-purple-500/60 focus-within:bg-gray-800/60 focus-within:shadow-[0_0_15px_rgba(168,85,247,0.1)] overflow-hidden group">
      <CommentTextEditor ref="commentEditor" v-model="value"
        :placeholder="$t('form.placeholder.enter_your_comment')" :disabled="sending" @submit="submitComment" />
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
const route = useRoute()
const api = useApi()

const props = defineProps({
  topicSlug: {
    type: String,
    required: true,
  }
})

const emits = defineEmits(['created'])
const value = ref({
  content: '', // 内容
  imageList: [],
})
const sending = ref(false) // 发送中
const commentEditor = ref(null) // 编辑器组件

function clearTextEditor() {
  value.value.content = ''
  value.value.imageList = []
  commentEditor.value.clear()
}

async function submitComment() {
  if (!value.value.content) {
    useMsgError(i18n.t('message.please_enter_comment'))
    return
  }
  if (sending.value) {
    return
  }
  if (commentEditor.value && commentEditor.value.isOnUpload()) {
    return
  }
  sending.value = true
  try {
    const data = await api.addTopicComment(props.topicSlug, {
      content: value.value.content,
      imageList: value.value.imageList,
    })
    emits('created', data)
    clearTextEditor()
    useMsgSuccess(i18n.t('message.comment_success'))
  } catch (e) {
    console.error(e)
    useMsgError(e.message || e)
  } finally {
    sending.value = false
  }
}

</script>
