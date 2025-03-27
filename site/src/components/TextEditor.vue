<template>
  <div class="text-editor">
    <textarea
      ref="textarea"
      v-model="post.content"
      :placeholder="$t('form.placeholder.enter_post_content')"
      :style="{ 'min-height': `${height}px`, 'height': `${height}px` }"
      @input="onInput"
      @paste="handleParse"
      @drop="handleDrag"
      @keydown.ctrl.enter="doSubmit"
      @keydown.meta.enter="doSubmit" />
    <div v-show="showImageUpload" class="text-editor-image-uploader">
      <image-upload
        ref="imageUploader"
        v-model="post.imageList"
        v-model:on-upload="imageUploading"
        @input="onInput" />
    </div>
    <div class="text-editor-bar">
      <div class="text-editor-actions">
        <div
          class="text-editor-action-item"
          :class="{ active: showImageUpload }"
          @click="switchImageUpload">
          <icon name="ImagePlus" />
          <span>&nbsp;{{ $t('publish.add_image') }}</span>
        </div>
      </div>
      <div class="text-editor-btn">
        <span>Ctrl/⌘ + Enter</span>
        <button class="button is-success is-small" @click="doSubmit">
          {{ $t('form.button.send') }}
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

function doSubmit() {
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
    imageUploaderRef.value?.addFiles([file])
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
    imageUploaderRef.value?.addFiles(files)
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

const imageUploaderRef = ref(null)
const textareaRef = ref(null)

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
.text-editor {
  border: 1px solid var(--border-color);

  textarea {
    width: 100%;
    font-family: inherit;
    background: var(--bg-color2);
    border: 0;
    outline: 0;
    display: block;
    position: relative;
    resize: none;
    line-height: 1.8;
    padding: 15px 15px 20px;
    overflow: auto;
    overscroll-behavior: contain;
    transition: all 100ms linear;
    color: var(--text-color);
  }

  .text-editor-image-uploader {
    padding: 10px;
  }

  .text-editor-bar {
    background-color: var(--bg-color);
    border-top: 1px solid var(--border-color);
    padding: 5px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .text-editor-actions {
      .text-editor-action-item {
        cursor: pointer;
        color: var(--text-color3);
        user-select: none;

        i,
        span {
          font-size: 16px;
        }

        &:hover {
          color: var(--text-link-color);
        }

        &.active {
          color: var(--text-link-color);
          font-weight: 500;
        }
      }
    }

    .text-editor-btn {
      display: flex;
      align-items: center;

      span {
        font-size: 12px;
        color: var(--text-color3);
        margin-right: 5px;
      }
    }
  }
}
</style>
