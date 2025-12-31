<template>
  <client-only>
    <MdEditor
      v-model="value"
      language="en-US"
      theme="dark"
      :toolbars="toolbars"
      :style="{ height: height }"
      :placeholder="placeholder"
      :class="['md-editor-dark-custom', props.class]"
      :preview="false"
      :maxLength="contentLength"
      @onHtmlChanged="handleHtmlChange"
      @onChange="change"
      @onUploadImg="uploadImg" />
  </client-only>
</template>

<script setup>
import { MdEditor } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  height: {
    type: String,
    default: '500px',
  },
  placeholder: {
    type: String,
    default: '',
  },
  class: {
    type: [Array, String],
    default: '',
  },
})

const contentLength = 10000

const emits = defineEmits([
  'update:modelValue',
  'htmlChanged',
])

const value = ref(props.modelValue)

const toolbars = ref([
  'bold',
  'underline',
  'italic',
  'title',
  'strikeThrough',
  '-',
  'sub',
  'sup',
  'quote',
  'unorderedList',
  'orderedList',
  'task',
  '-',
  'codeRow',
  'code',
  'link',
  'image',
  'table',
  // "mermaid",
  // "katex",
  '-',
  'revoke',
  'next',
  'save',
  '=',
  'pageFullscreen',
  'fullscreen',
  'preview',
  'htmlPreview',
  'catalog',
])

function change(v) {
  emits('update:modelValue', v)
}

function handleHtmlChange(html) {
  emits('htmlChanged', html)
}

async function uploadImg(files, callback) {
  const res = await Promise.all(
    files.map((file) => {
      return new Promise((rev, rej) => {
        const formData = new FormData()
        formData.append('image', file, file.name)
        useHttp('/api/upload', {
          method: 'POST',
          body: formData,
        })
          .then(res => rev(res))
          .catch(error => rej(error))
      })
    }),
  )
  callback(res.map(item => item.url))
}
</script>

<style lang="scss" scoped>
.md-editor-dark-custom {
  --md-color: #cbd5e1;
  --md-hover-color: #e2e8f0;
  --md-bk-color: #0f141e;
  --md-bk-color-outstand: #111827;
  --md-bk-hover-color: #111827;
  --md-border-color: #374151;
  --md-border-hover-color: #4b5563;
  --md-border-active-color: #8b5cf6;
  --md-modal-mask: rgba(189, 40, 40, 0.34);
  --md-modal-shadow: 0px 6px 24px 2px rgba(17, 24, 39, 0.6);
  --md-scrollbar-bg-color: #0b1220;
  --md-scrollbar-thumb-color: #2f3a4a;
  --md-scrollbar-thumb-hover-color: #465569;
  --md-scrollbar-thumb-active-color: #465569;
}
</style>
