<template>
  <div
    :class="[{ 'opacity-60 pointer-events-none': disabled }, 'group relative w-full overflow-hidden rounded-xl bg-gray-800/40 backdrop-blur-sm border border-purple-500/20 transition-all duration-300 hover:border-purple-500/40 focus-within:border-purple-500/60 focus-within:bg-gray-800/60 focus-within:shadow-[0_0_15px_rgba(168,85,247,0.1)]']">
    <!-- Textarea -->
    <textarea
      ref="textarea"
      v-model="post.content"
      :placeholder="$t('form.placeholder.enter_post_content')"
      :style="{ 'min-height': `${height}px`, 'height': `${height}px` }"
      :class="['text-input w-full bg-transparent text-gray-200 border-none rounded-t-xl p-5 focus:outline-none focus:ring-0 placeholder-gray-500 resize-y text-base font-medium leading-relaxed', disabled ? 'cursor-not-allowed' : '']"
      :disabled="disabled"
      @input="onInput"
      @paste="handleParse"
      @drop="handleDrag"
      @keydown.ctrl.enter="doSubmit"
      @keydown.meta.enter="doSubmit" />

    <!-- Image Uploader Area -->
    <div v-show="showImageUpload" class="px-4 pb-4">
      <div class="p-4 rounded-lg bg-gray-900/40 border border-dashed border-gray-700/50">
        <image-upload
          ref="imageUploader"
          v-model="post.imageList"
          v-model:on-upload="imageUploading"
          @input="onInput" />
      </div>
    </div>

    <!-- Toolbar -->
    <div
      class="editor-toolbar border-t border-purple-500/20 px-4 py-3 flex items-center justify-between bg-gray-900/30 transition-colors duration-300 group-hover:bg-gray-800/40">
      <!-- Left Actions -->
      <div class="flex items-center gap-4">
        <button
          class="flex items-center gap-2 text-sm transition-colors duration-200"
          :class="showImageUpload ? 'text-purple-400 font-semibold' : 'text-gray-400 hover:text-purple-400'"
          @click="switchImageUpload">
          <FontAwesome :icon="['fas', 'image']" class="w-5 h-5" />
          <span>{{ $t('publish.add_image') }}</span>
        </button>
      </div>

      <!-- Right Actions -->
      <div class="flex items-center gap-4">
        <span class="hidden sm:inline-block text-xs text-gray-500 font-medium">Ctrl/⌘ + Enter</span>
        <button
          class="editor-submit-btn px-6 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white text-sm font-bold rounded-lg shadow-[0_0_10px_rgba(168,85,247,0.4)] hover:shadow-[0_0_20px_rgba(168,85,247,0.6)] transform hover:scale-105 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="imageUploading"
          @click="doSubmit">
          <span v-if="imageUploading" class="flex items-center gap-2">
            <FontAwesome :icon="['fas', 'spinner']" spin /> {{ $t('file.uploading') }}
          </span>
          <span v-else>{{ $t('form.button.send') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  height: {
    type: Number,
    default: 120,
  },
  value: {
    type: Object,
    default: () => ({
      content: '',
      imageList: [],
    }),
  },
})

const emit = defineEmits(['submit', 'update:modelValue'])

const post = ref({ ...props.value })
const showImageUpload = ref(false)
const imageUploading = ref(false)
const imageUploaderRef = ref(null)
const textareaRef = ref(null)

function doSubmit() {
  if (imageUploading.value) return
  emit('submit')
}

const onInput = () => {
  emit('update:modelValue', post.value)
}

const isOnUpload = () => imageUploading.value

const handleParse = (e) => {
  const items = e.clipboardData?.items
  if (!items?.length) return

  let file = null
  for (const item of items) {
    if (item.type.includes('image')) {
      file = item.getAsFile()
    }
  }

  if (file) {
    e.preventDefault() // Prevent default clipboard behavior
    showImageUpload.value = true
    nextTick(() => {
      imageUploaderRef.value?.addFiles([file])
    })
  }
}

const handleDrag = (e) => {
  e.stopPropagation()
  e.preventDefault()

  const items = e.dataTransfer?.items
  if (!items?.length) return

  const files = [...items]
    .filter(item => item.type.includes('image'))
    .map(item => item.getAsFile())

  if (files.length) {
    showImageUpload.value = true
    nextTick(() => {
      imageUploaderRef.value?.addFiles(files)
    })
  }
}

const switchImageUpload = () => {
  showImageUpload.value = !showImageUpload.value
}

const clear = () => {
  post.value.content = ''
  post.value.imageList = []
  showImageUpload.value = false
  imageUploaderRef.value?.clear()
  onInput()
}

const focus = () => {
  textareaRef.value?.focus()
}

watch(
  () => props.value,
  (newValue) => {
    post.value = { ...newValue }
  },
  { deep: true },
)

defineExpose({
  isOnUpload,
  focus,
  clear,
})
</script>

<style lang="scss" scoped>
/* Scoped styles replaced by Tailwind classes */
</style>
