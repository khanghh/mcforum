<template>
  <div class="space-y-6">
    <!-- Header -->
    <div v-if="commentCount > 0" class="flex items-center justify-between">
      <h3 class="text-lg font-bold text-gray-300">{{ commentCount }} Comments</h3>
      <div class="text-sm text-gray-500">
        Sort by: <span class="text-purple-400 font-semibold cursor-pointer">Newest</span>
      </div>
    </div>

    <!-- Login Prompt or Input -->
    <template v-if="isLogin">
      <div class="mb-8 p-4 rounded-lg bg-purple-500/10 border border-purple-500/30">
        <h3 class="text-lg font-bold text-purple-300 mb-3 font-['Saira Semi Condensed']">Post a Comment</h3>
        <CommentInput ref="input" :topicSlug="topicSlug" @created="commentCreated" />
      </div>
    </template>

    <!-- Comments List -->
    <CommentList ref="list" :topicSlug="topicSlug" />
  </div>
</template>

<script setup>
const props = defineProps({
  topicSlug: {
    type: String,
    required: true,
  },
  commentCount: {
    type: Number,
    default: 0,
  },
})

const emits = defineEmits(['created'])
const userStore = useUserStore()
const configStore = useConfigStore()
const { user } = storeToRefs(userStore)

const isLogin = computed(() => {
  return !!user.value && !!user.value.id
})

const config = computed(() => {
  return configStore.config
})

const list = ref(null)

function commentCreated(data) {
  list.value.append(data)
  emits('created', data)
}

</script>
