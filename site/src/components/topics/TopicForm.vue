<template>
  <form class="space-y-6 w-full">
    <!-- Topic Title - Full Width -->
    <div class=" full-width-section space-y-2">
      <label for="topic-title" class="text-sm font-medium text-gray-200 flex items-center">
        <Icon name="Fa7SolidHeading" class="mr-2" />
        {{ $t('form.label.topic_title') }}
      </label>
      <div class="relative">
        <input id="topic-title" v-model="postForm.title" type="text" autocomplete="off"
          class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-600/50 transition-colors transition-shadow"
          :placeholder="$t('form.placeholder.enter_topic_title')" />
        <div class="absolute right-3 top-3">
          <Icon name="Fa7SolidEdit" class="text-gray-400" />
        </div>
      </div>
      <p class="text-xs text-gray-300">{{ $t('form.help.title_must_be_10_to_200_chars') }}</p>
    </div>

    <!-- Forum Category Selection - Full Width -->
    <div class="full-width-section space-y-2">
      <label for="forum-category" class="text-sm font-medium text-gray-200 flex items-center">
        <Icon name="Fa7SolidFolder" class="mr-2" />
        {{ $t('form.label.select_category') }}
      </label>
      <div class="relative">
        <select id="forum-category" v-model="postForm.forumId"
          class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-200 focus:outline-none focus:ring-2 focus:ring-purple-600/50 transition-colors transition-shadow appearance-none">
          <option value="" disabled>Select a category...</option>
          <option v-for="forum in forums" :key="forum.id" :value="forum.id">
            {{ forum.name }}
          </option>
        </select>
        <div class="absolute right-3 top-3 pointer-events-none">
          <Icon name="Fa7SolidChevronDown" class="text-gray-400" />
        </div>
      </div>

      <!-- Quick Category Selection -->
      <div v-if="forums.length" class="mt-3">
        <p class="text-xs text-gray-300 mb-2">{{ $t('form.label.quick_select_category') }}</p>
        <div class="flex flex-wrap gap-2">
          <button v-for="forum in forums.slice(0, 4)" :key="forum.id" type="button" :class="[
            'text-xs px-3 py-1 rounded-full border transition-all duration-200 ease-in-out',
            postForm.forumId === forum.id
              ? 'bg-purple-900/50 border-purple-700/50 text-purple-300'
              : 'bg-purple-900/30 hover:bg-purple-900/50 hover:-translate-y-0.5 border-purple-700/50 text-gray-300',
          ]" @click="postForm.forumId = forum.id">
            {{ forum.name }}
          </button>
        </div>
      </div>
    </div>

    <!-- Content Editor - Full Width -->
    <div class="full-width-section space-y-2">
      <div class="flex justify-between items-center">
        <label class="text-sm font-medium text-gray-200 flex items-center">
          <Icon name="Fa7SolidAlignLeft" class="mr-2" />
          {{ $t('form.label.topic_content') }}
        </label>
        <div class="text-xs text-gray-300">
          <span class="tabular-nums">{{ postForm.content.length }}/10000</span>
        </div>
      </div>

      <MarkdownEditor v-model="postForm.content" @onHtmlChanged="$emit('onHtmlChanged', $event)"
        class="md-editor-dark-custom" />

      <!-- Content Help -->
      <div
        class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-2 text-xs text-gray-400">
        <div>
          <Icon name="Fa7SolidCircleInfo" class="mr-1" />
          <span>{{ $t('form.label.topic_content_help') }}</span>
        </div>
      </div>
    </div>

    <!-- Tags Section -->
    <div class="full-width-section space-y-2">
      <label class="text-sm font-medium text-gray-300 flex items-center">
        <Icon name="Fa7SolidTag" class="mr-2" />
        Tags
      </label>
      <TagInput v-model="postForm.tags" />
      <p class="text-xs text-gray-400">{{ $t('form.label.add_relevant_tags') }}</p>
    </div>

    <!-- Poll/Voting Section -->
    <div class="space-y-4">
      <div class="flex items-center gap-4">
        <label class="relative inline-flex items-center cursor-pointer">
          <input v-model="enablePoll" type="checkbox" class="sr-only peer"
            @change="handlePollToggle(enablePoll)" />
          <div
            class="w-5 h-5 bg-gray-700 border border-gray-600 rounded flex items-center justify-center peer-checked:bg-purple-600 peer-checked:border-purple-600 transition-colors">
            <Icon v-show="enablePoll" name="Fa7SolidCheck" class="text-white text-xs transition-opacity" />
          </div>
        </label>
        <label class="text-sm font-medium text-gray-300 flex items-center">
          {{ $t('form.label.create_poll') }}
        </label>
      </div>

      <div v-if="enablePoll" class="space-y-4">
        <!-- Poll Question -->
        <div>
          <label for="poll-question" class="text-xs font-medium text-gray-300 block mb-1">
            <Icon name="Fa7SolidCircleQuestion" class="mr-1" />
            {{ $t('form.label.poll_question') }}
          </label>
          <input id="poll-question" v-model="pollOpts.question" type="text"
            class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-600/50 transition-colors transition-shadow"
            :placeholder="$t('form.placeholder.enter_poll_question')" />
        </div>

        <!-- Poll Options -->
        <div>
          <div class="flex justify-between items-center mb-2">
            <label class="text-xs font-medium text-gray-300">
              <Icon name="Fa7SolidListOl" class="mr-1" />
              {{ $t('form.label.poll_options') }}
            </label>
            <button type="button"
              class="text-xs px-2 py-1 bg-purple-600/30 hover:bg-purple-600/50 text-purple-300 rounded transition-colors"
              @click="addPollOption">
              <Icon name="Fa7SolidPlus" />
              {{ $t('form.label.add_poll_option') }}
            </button>
          </div>

          <div class="space-y-2 max-h-[300px] overflow-y-auto">
            <div v-for="(option, index) in pollOpts.options" :key="index"
              class="flex items-center gap-2 p-2 bg-gray-800/30 rounded-lg border border-gray-700/50 transition-all duration-200 ease-in-out hover:border-purple-500/50">
              <div class="flex-grow">
                <input v-model="pollOpts.options[index]" type="text"
                  class="w-full px-2 py-1 bg-transparent border-none text-white placeholder-gray-500 focus:outline-none focus:ring-0 text-sm"
                  :placeholder="`${$t('form.placeholder.enter_poll_option')} ${Number(index) + 1}`" />
              </div>
              <button type="button" :disabled="pollOpts.options.length <= 2"
                class="text-red-400 hover:text-red-300 p-1 disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                @click="removePollOption(Number(index))">
                <div>
                  <Icon name="Fa7SolidTimes" />
                </div>
              </button>
            </div>
          </div>
        </div>

        <!-- Poll Settings Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="sm:col-span-2">
            <label for="poll-duration" class="text-xs font-medium text-gray-300 block mb-1">
              <Icon name="Fa7SolidClock" class="mr-1" />
              {{ $t('form.label.poll_duration') }}
            </label>
            <div class="relative">
              <select id="poll-duration" v-model="pollOpts.durationHours"
                class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-200 focus:outline-none focus:ring-2 focus:ring-purple-600/50 transition-colors transition-shadow appearance-none">
                <option v-for="option in pollDurationOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
                </option>
              </select>
              <div class="absolute right-3 top-3 pointer-events-none">
                <Icon name="Fa7SolidChevronDown" class="text-gray-400" />
              </div>
            </div>
          </div>
        </div>

        <!-- Multiple Choice - New Row -->
        <div
          class="flex items-center justify-between p-3 bg-gray-800/30 rounded-lg border border-gray-700/50">
          <div>
            <h4 class="font-medium text-white text-sm">{{ $t('form.label.multiple_choices') }}</h4>
            <p class="text-xs text-gray-300 mt-1">{{ $t('form.label.multiple_choices_desc') }}</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="pollOpts.multiSelect" type="checkbox" class="sr-only peer" />
            <div
              class="w-10 h-5 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-5 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
            </div>
          </label>
        </div>

        <!-- Poll Privacy -->
        <div
          class="flex items-center justify-between p-3 bg-gray-800/30 rounded-lg border border-gray-700/50">
          <div>
            <h4 class="font-medium text-white text-sm">{{ $t('form.label.public_results') }}</h4>
            <p class="text-xs text-gray-300 mt-1">{{ $t('form.label.public_results_desc') }}</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="pollOpts.publicResults" type="checkbox" class="sr-only peer" checked />
            <div
              class="w-10 h-5 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-5 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
            </div>
          </label>
        </div>

        <!-- Poll Change Vote -->
        <div
          class="flex items-center justify-between p-3 bg-gray-800/30 rounded-lg border border-gray-700/50">
          <div>
            <h4 class="font-medium text-white text-sm">{{ $t('form.label.allow_vote_changes') }}</h4>
            <p class="text-xs text-gray-300 mt-1">{{ $t('form.label.allow_vote_changes_desc') }}</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input v-model="pollOpts.allowVoteChange" type="checkbox" class="sr-only peer" checked />
            <div
              class="w-10 h-5 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-5 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
            </div>
          </label>
        </div>
      </div>
    </div>

    <!-- Hidden Content Section - Full Width Below Poll -->
    <div class="full-width-section space-y-4">
      <div class="flex items-center gap-4">
        <label class="relative inline-flex items-center cursor-pointer">
          <input v-model="enableHiddenContent" type="checkbox" class="sr-only peer" />
          <div
            class="w-5 h-5 bg-gray-700 border border-gray-600 rounded flex items-center justify-center peer-checked:bg-purple-600 peer-checked:border-purple-600 transition-colors">
            <Icon v-show="enableHiddenContent" name="Fa7SolidCheck"
              class="text-white text-xs transition-opacity" />
          </div>
        </label>
        <span class="text-sm font-medium text-gray-300 flex items-center">
          {{ $t('form.label.enable_hidden_content') }}
        </span>
      </div>

      <div v-if="enableHiddenContent" class="space-y-3">
        <div class="relative">
          <textarea v-model="postForm.hiddenContent" rows="6"
            class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-600/50 transition-colors transition-shadow resize-none font-mono"
            :placeholder="$t('form.placeholder.enter_hidden_content')"></textarea>
        </div>
        <div class="flex items-center text-xs text-gray-400">
          <Icon name="Fa7SolidCircleInfo" class="mr-1" />
          <span>{{ $t('form.label.hidden_content_desc') }}</span>
        </div>
      </div>
    </div>

    <!-- Divider -->
    <div class="full-width-section border-t border-gray-700/50"></div>

    <!-- Action Buttons - Full Width -->
    <div class="full-width-section flex flex-col sm:flex-row justify-end items-center gap-4 pt-4">
      <nuxt-link to="/"
        class="px-8 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg font-medium border border-gray-700 flex items-center justify-center transition-colors gaming-title">
        <Icon name="Fa7SolidTimes" class="mr-2" />
        {{ $t('form.button.cancel') }}
      </nuxt-link>
      <button type="button" :disabled="publishing"
        @click="publishTopic"
        class="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center justify-center shadow-[0_0_15px_rgba(139,92,246,0.3),0_0_25px_rgba(139,92,246,0.2)] tracking-[0.5px] transition-all hover:scale-[1.02] gaming-title ">
        <Icon v-if="!publishing" name="Fa7SolidPaperPlane" class="mr-2" />
        <Icon v-else name="TablerLoader2" class="mr-2 animate-spin" />
        {{ publishing ? ('Publishing...') : (editMode ? $t('form.button.save_changes') : $t('form.button.create_topic'))
        }}
      </button>
    </div>
  </form>
</template>

<script setup lang="ts">

// import { MdEditor } from 'md-editor-v3'
import TagInput from '@/components/topics/TagInput.vue'
import type { TopicPollPayload, CreateTopicPayload } from '@/composables/api'

type Forum = {
  id: number
  name: string
  slug: string
}

const props = defineProps({
  modelValue: {
    type: Object as PropType<CreateTopicPayload>,
    required: true,
  },
  editMode: {
    type: Boolean,
    default: false,
  },
  submitCallback: {
    type: Function as PropType<(v: CreateTopicPayload) => Promise<any>>,
    default: async () => { },
  },
})

const enablePoll = ref(false)
const enableHiddenContent = ref(!!props.modelValue.hiddenContent)
const publishing = ref(false)
const postForm = props.modelValue
const forums = ref<Forum[]>([])

forums.value = await useApi().getForumList().catch(() => [])

const emits = defineEmits([
  'update:modelValue',
  'onHtmlChanged',
  'onSubmit'
])

// Watch for hidden content toggle changes
watch(enableHiddenContent, (newValue) => {
  if (!newValue) {
    postForm.hiddenContent = ''
  }
})

// Function to handle poll checkbox toggle
function handlePollToggle(enabled: boolean) {
  if (enabled) {
    postForm.poll = { ...pollOpts.value }
  } else {
    postForm.poll = undefined
  }
}

// Poll form state
const pollOpts = ref<TopicPollPayload>({
  question: '',
  options: ['', ''],
  durationHours: 7,
  multiSelect: false,
  publicResults: true,
  allowVoteChange: true,
})

const pollDurationOptions = [
  { value: 24, label: '1 day' },
  { value: 72, label: '3 days' },
  { value: 168, label: '7 days' },
  { value: 336, label: '14 days' },
  { value: 720, label: '30 days' },
  { value: 0, label: 'No limit' },
]

function publishTopic() {
  publishing.value = true
  return props.submitCallback(props.modelValue).finally(() => {
    publishing.value = false
  })
}

function addPollOption() {
  if (pollOpts.value.options.length < 10) {
    pollOpts.value.options.push('')
  }
}

function removePollOption(index: number) {
  if (pollOpts.value.options.length > 2) {
    pollOpts.value.options.splice(index, 1)
  }
}

</script>

<style scoped>
select option {
  background-color: #1a1a2e;
}
</style>