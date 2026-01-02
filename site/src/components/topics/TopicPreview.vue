<template>
  <client-only>
    <div class="space-y-6 w-full">
      <!-- Preview Header -->
      <div class="flex items-center justify-between gap-3 mb-4">
        <div class="flex items-center gap-3">
          <div class="relative group">
            <Avatar :src="author.avatar"
              :username="author.username"
              class="w-12 h-12 rounded border-2 border-purple-300 flex-shrink-0 object-cover" />
          </div>
          <div>
            <div class="flex items-center gap-2 flex-wrap">
              <nuxt-link :to="`/users/${author.username}`"
                class="font-bold text-purple-300 gaming-title text-lg hover:text-purple-400 transition-colors"
                itemprop="author">
                {{ author.username }}
              </nuxt-link>
              <span v-if="author.role"
                class="px-2 py-0.5 bg-red-500/20 text-red-400 text-xs font-bold rounded uppercase">
                {{ author.role }}
              </span>
            </div>

            <div class="text-sm text-gray-500 mt-1">
              {{ usePrettyDate(Date.now()) }}
            </div>
          </div>
        </div>
      </div>

      <article>
        <h1 class="text-2xl sm:text-3xl font-bold mb-6 text-white gaming-title" itemprop="headline">
          {{ payload.title }}
        </h1>
        <!-- Preview Content -->
        <div class="prose prose-invert max-w-none">
          <div v-if="previewHTML" class="text-gray-200" v-html="previewHTML"></div>
          <div v-else class="text-gray-500 italic">No content to preview</div>
        </div>

        <!-- Preview Hidden Content -->
        <div v-if="hiddenContent" class="mt-6">
          <div class="p-4 bg-yellow-900/20 border border-yellow-700/50 rounded-lg">
            <div class="flex items-center gap-2 mb-3">
              <Icon name="Fa7SolidEyeSlash" class="text-yellow-400" />
              <span class="font-medium text-yellow-300">Hidden Content</span>
            </div>
            <div class="prose prose-invert max-w-none">
              <div class="text-gray-300" v-html="renderMarkdown(payload.hiddenContent || '')"></div>
            </div>
          </div>
        </div>


        <div v-if="payload.tags" class="flex flex-wrap gap-2 mb-6">
          <nuxt-link v-for="tag in payload.tags" :key="tag" :to="`/tags/${tag}`"
            class="px-3 py-1 bg-purple-500/10 text-purple-300 text-xs rounded-full border border-purple-500/20 hover:bg-purple-500/20 transition-colors">
            #{{ tag }}
          </nuxt-link>
        </div>
      </article>

      <!-- Back to Edit Button -->
      <div class="pt-4 border-t border-gray-700/50">
        <button
          class="px-6 py-2.5 bg-purple-600 hover:bg-purple-700 text-white rounded-lg font-medium transition-colors flex items-center"
          @click="$emit('onBackToEdit')">
          <Icon name="Fa7SolidArrowLeft" class="mr-2" />
          Back to Edit
        </button>
      </div>
    </div>
  </client-only>
</template>

<script setup lang="ts">
import type { UserInfo } from '~/types/user'
const props = defineProps({
  payload: {
    type: Object as PropType<CreateTopicPayload>,
    required: true,
  },
  author: {
    type: Object as PropType<UserInfo>,
    required: true,
  },
  previewHTML: {
    type: String,
    required: false,
  },
})
defineEmits(['onBackToEdit'])

const hiddenContent = computed(() => props.payload.hiddenContent)

const enabledHideContent = computed(() => {
  return props.payload.hiddenContent && props.payload.hiddenContent.length > 0
})

const enabledPool = computed(() => {
  return props.payload.poll
})

function renderMarkdown(content: string): string {
  return ''
}
</script>