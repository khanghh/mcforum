<template>
  <div class="space-y-6">
    <!-- Preview Header -->
    <div class="border-b border-gray-700/50 pb-4">
      <h2 class="gaming-title text-2xl font-bold text-white mb-2">
        {{ postForm.title || 'Untitled Topic' }}
      </h2>
      <div class="flex items-center gap-4 text-sm text-gray-400">
        <div class="flex items-center gap-2">
          <Avatar :user="userStore.user" :size="32" class="rounded" />
          <span>{{ userStore.user?.username || 'Anonymous' }}</span>
        </div>
        <span>•</span>
        <span>{{ new Date().toLocaleDateString() }}</span>
        <span v-if="selectedForum">•</span>
        <span v-if="selectedForum" class="text-purple-400">{{ selectedForum.name }}</span>
      </div>
      <div v-if="postForm.tags && postForm.tags.length" class="flex flex-wrap gap-2 mt-3">
        <span v-for="tag in postForm.tags" :key="tag"
          class="px-3 py-1 bg-purple-900/30 text-purple-300 rounded-full text-xs border border-purple-700/50">
          {{ tag }}
        </span>
      </div>
    </div>

    <!-- Preview Content -->
    <div class="prose prose-invert max-w-none">
      <div v-if="postForm.content" class="text-gray-200" v-html="renderMarkdown(postForm.content)"></div>
      <div v-else class="text-gray-500 italic">No content to preview</div>
      <div class="italic">enableHideContent: {{ enabledHideContent }}</div>
      <div class="italic">enabledPool: {{ enabledPool }}</div>
    </div>

    <!-- Preview Hidden Content -->
    <div v-if="hiddenContent" class="mt-6">
      <div class="p-4 bg-yellow-900/20 border border-yellow-700/50 rounded-lg">
        <div class="flex items-center gap-2 mb-3">
          <Icon name="Fa7SolidEyeSlash" class="text-yellow-400" />
          <span class="font-medium text-yellow-300">Hidden Content</span>
        </div>
        <div class="prose prose-invert max-w-none">
          <div class="text-gray-300" v-html="renderMarkdown(postForm.hideContent)"></div>
        </div>
      </div>
    </div>

    <!-- Back to Edit Button -->
    <div class="pt-4 border-t border-gray-700/50">
      <button
        class="px-6 py-2.5 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition-colors flex items-center"
        @click="$emit('switchToCreate')">
        <Icon name="Fa7SolidArrowLeft" class="mr-2" />
        Back to Edit
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Forum } from '~/types'
const props = defineProps<{
  postForm: any
  selectedForum: Forum
}>()

const hiddenContent = computed(() => props.postForm.hideContent)

const userStore = useUserStore()

const enabledHideContent = computed(() => {
  return props.postForm.hideContent && props.postForm.hideContent.length > 0
})

const enabledPool = computed(() => {
  return props.postForm.poll
})

function renderMarkdown(content: string): string {
  return ''
}
</script>