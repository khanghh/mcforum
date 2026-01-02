<template>
  <div>
    <!-- Header with back button -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 rounded-lg bg-gradient-to-r from-purple-600/30 to-pink-600/30 flex items-center justify-center">
          <Icon name="Fa7SolidPen" class="text-purple-400 text-lg" />
        </div>
        <h1
          class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 gaming-title">
          Create New Topic
        </h1>
      </div>
      <nuxt-link to="/topics"
        class="flex items-center text-sm text-gray-300 hover:text-white transition-colors px-4 py-2 rounded-lg bg-gray-800/40 hover:bg-gray-700/60 border border-gray-700/50">
        <Icon name="Fa7SolidArrowLeft" class="mr-2" />
        <span>Back to Forum</span>
      </nuxt-link>
    </div>

    <!-- Tabs -->
    <div class="mb-6">
      <div class="flex gap-2 border-b border-gray-700/50">
        <button :class="[
          'px-6 py-3 font-medium transition-all duration-200',
          activeTab === 'form'
            ? 'text-purple-400 border-b-2 border-purple-500'
            : 'text-gray-400 hover:text-gray-300',
        ]" @click="activeTab = 'form'">
          <Icon name="Fa7SolidEdit" class="mr-2" />
          Create
        </button>
        <button :class="[
          'px-6 py-3 font-medium transition-all duration-200',
          activeTab === 'preview'
            ? 'text-purple-400 border-b-2 border-purple-500'
            : 'text-gray-400 hover:text-gray-300',
        ]" @click="activeTab = 'preview'">
          <Icon name="Fa7SolidEye" class="mr-2" />
          Preview
        </button>
      </div>
    </div>

    <!-- Create Tab -->
    <div v-show="activeTab === 'form'"
      class="gradient-border rounded-2xl p-6 sm:p-8 relative overflow-hidden backdrop-blur-sm shadow-xl shadow-purple-900/10">
      <!-- Subtle gradient overlay -->
      <div class="absolute inset-0 bg-gradient-to-br from-gray-800/20 via-gray-900/30 to-indigo-900/20"></div>

      <div class="relative z-10">
        <TopicForm v-model="postForm" :forums="forums" :submitCallback="createTopic"
          @onHtmlChanged="handleHtmlChange" />
      </div>
    </div>

    <!-- Preview Tab -->
    <div v-show="activeTab === 'preview'"
      class="gradient-border rounded-2xl p-6 sm:p-8 relative overflow-hidden backdrop-blur-sm shadow-xl shadow-purple-900/10">
      <!-- Subtle gradient overlay -->
      <div class="absolute inset-0 bg-gradient-to-br from-gray-800/20 via-gray-900/30 to-indigo-900/20"></div>

      <div class="relative z-10">
        <TopicPreview :payload="postForm" :author="author!!" :previewHTML="previewHTML" />
      </div>
    </div>

    <!-- Posting Guidelines - Full Width -->
    <div
      class="full-width-section mt-6 p-4 bg-gradient-to-r from-purple-900/20 to-indigo-900/20 rounded-lg border border-purple-500/30 backdrop-blur-sm">
      <div class="flex items-start">
        <div class="flex-shrink-0 mt-1">
          <Icon name="Fa7SolidScroll" class="text-purple-400" />
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-white">Posting Guidelines</h3>
          <ul class="text-xs text-gray-300 mt-1 space-y-1">
            <li>
              <Icon name="Fa7SolidCheckCircle" class="text-green-400 mr-1" />
              Do not post content that violates laws or forum rules
            </li>
            <li>
              <Icon name="Fa7SolidCheckCircle" class="text-green-400 mr-1" />
              Use proper Vietnamese with diacritics, avoid excessive abbreviations
            </li>
            <li>
              <Icon name="Fa7SolidCheckCircle" class="text-green-400 mr-1" />
              Use clear titles that accurately reflect the content
            </li>
            <li>
              <Icon name="Fa7SolidCheckCircle" class="text-green-400 mr-1" />
              Choose the correct category for your post
            </li>
            <li>
              <Icon name="Fa7SolidCheckCircle" class="text-green-400 mr-1" />
              No spam, unauthorized advertising, or sensitive content
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import TopicForm from '@/components/topics/TopicForm.vue'
import TopicPreview from '@/components/topics/TopicPreview.vue'
import type { UserInfo, Forum } from '@/types'
import type { CreateTopicPayload } from '@/composables/api'

import 'md-editor-v3/lib/style.css'

const i18n = useI18n()
const userStore = useUserStore()
const dialog = useConfirmDialog()
const api = useApi()

definePageMeta({
  middleware: ['auth'],
  layout: 'default',
})

const activeTab = ref('form')
const forums = ref<Forum[]>([])
const previewHTML = ref('')
const { user: author } = storeToRefs(userStore)

const postForm = ref<CreateTopicPayload>({
  forumId: 1,
  title: '',
  tags: [] as string[],
  content: '',
  hiddenContent: '',
  imageList: [] as string[],
})

function handleHtmlChange(html: string) {
  previewHTML.value = html
}

function createTopic(payload: CreateTopicPayload) {
  return api.createTopic(payload).then((topic) => {
    navigateTo(`/topics/${topic.slug}`)
  }).catch(async (error: any) => {
    const errMsg = error?.data?.error?.message || error.message || 'Unknown error'
    console.error('Error creating topic:', error)
    return dialog.show({
      title: 'Error Creating Topic',
      message: `${errMsg}`,
      variant: 'warning',
    })
  })
}

useHead({
  title: useSiteTitle(i18n.t('page.create_topic')),
})
</script>

<style scoped>
/* Static gradient border - cannot be done with Tailwind */
.gradient-border {
  position: relative;
  border: 2px solid transparent;
  background: linear-gradient(145deg, rgba(45, 45, 70, 0.9), rgba(35, 35, 55, 0.95)) padding-box,
    linear-gradient(45deg, #7c3aed, #8b5cf6, #a78bfa) border-box;
  background-size: 300% 300%;
}
</style>
