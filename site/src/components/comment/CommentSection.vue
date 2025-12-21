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
      <div v-if="isNeedEmailVerify" class="p-4 rounded-lg bg-yellow-500/10 border border-yellow-500/30 text-yellow-200">
        <i18n-t keypath="publish.verify_email_action" tag="p">
          <nuxt-link to="/user/profile/account" class="font-bold underline hover:text-yellow-100">
            {{ $t('navbar.profile') }} &gt; {{ $t('navbar.settings') }}
          </nuxt-link>
        </i18n-t>
      </div>
      <template v-else>
        <!-- Comment Input redesign wrapper -->
        <div class="mb-8 p-4 rounded-lg bg-purple-500/10 border border-purple-500/30">
          <h3 class="text-lg font-bold text-purple-300 mb-3 font-['Saira Semi Condensed']">Post a Comment</h3>
          <CommentInput ref="input" :entity-id="entityId" :entity-type="entityType" @created="commentCreated" />
        </div>
      </template>
    </template>

    <div v-else class="p-6 rounded-xl bg-gray-800/50 border border-gray-700 text-center">
      <div class="text-gray-400 mb-3">{{ $t('message.please_login_to_comment') }}</div>
      <button
        class="px-6 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold hover:scale-105 transition-transform"
        @click="useToSignIn()">
        {{ $t('page.signin') }}
      </button>
    </div>

    <!-- Comments List -->
    <CommentList ref="list" :entity-id="entityId" :entity-type="entityType" @reply="reply" />
  </div>
</template>

<script setup>
const props = defineProps({
  entityType: {
    type: String,
    required: true,
  },
  entityId: {
    type: Number,
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
const user = computed(() => {
  return userStore.user
})

const isLogin = computed(() => {
  return userStore.user !== null
})

const config = computed(() => {
  return configStore.config
})

const input = ref(null)
const list = ref(null)

// 是否需要先邮箱认证
const isNeedEmailVerify = computed(() => {
  return config.value.createCommentEmailVerified && user.value && !user.value.emailVerified
})

function commentCreated(data) {
  list.value.append(data)
  emits('created', data)
}
function reply(quote) {
  // this.$refs.input.reply(quote)
}
</script>

<style lang="scss" scoped>
/* Scoped styles replaced by Tailwind classes */
</style>
