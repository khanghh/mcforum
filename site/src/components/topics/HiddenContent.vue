<template>
  <div v-if="topic.hasHiddenContent" class="my-6">
    <div
      class="rounded-lg border border-dashed border-gray-600/30 bg-gray-800/50 overflow-hidden">
      <div class="flex items-center justify-between p-4">
        <div class="flex items-center gap-2 text-sm font-medium text-gray-300">
          <Icon name="Fa7SolidUnlock" />
          <span>
            {{ $t('feed.hidden_content') }}
          </span>
        </div>
        <button
          class="text-xs font-bold bg-purple-600 hover:bg-purple-700 text-white px-3 py-1.5 rounded transition-colors"
          @click="isExpanded = !isExpanded">
          {{ isExpanded ? $t('actions.hide_reply') : $t('actions.view') }}
        </button>
      </div>
      <div v-if="isExpanded && topic.hiddenContent"
        class="px-4 pb-4 prose prose-invert max-w-none border-t border-dashed border-gray-600/20 pt-4 font-mono">
        {{ topic.hiddenContent }}
      </div>
      <div v-else
        v-show="isExpanded"
        class="flex flex-col items-center justify-center border-t border-dashed border-gray-600/30 py-8 text-center">
        <Icon name="Fa7SolidLock" class="mb-3 text-4xl text-gray-500/50" />
        <template v-if="isLogin">
          <p class="font-medium text-gray-300">{{ $t('feed.hidden_content_locked') }}</p>
        </template>
        <template v-else>
          <p class="font-medium text-gray-300 mb-3">
            <i18n-t keypath="feed.hidden_content_login_required" scope="global">
              <template #login>
                <a :href="loginUrl" class="text-purple-400 hover:text-purple-300 underline">{{ $t('common.login') }}</a>
              </template>
            </i18n-t>
          </p>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { type Topic } from '~/types/topic'

const userStore = useUserStore()
const runtimeConfig = useRuntimeConfig()

const loginUrl = runtimeConfig.public.loginUrl
const isExpanded = ref(false)
const isLogin = computed(() => {
  return !!userStore.user
})

const props = defineProps({
  topic: {
    type: Object as PropType<Topic>,
    required: true,
  },
})

</script>