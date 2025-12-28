<template>
  <client-only>
    <MdEditor
      v-model="value"
      language="en-US"
      theme="dark"
      :toolbars="toolbars"
      :style="{ height: height }"
      :placeholder="placeholder"
      :class="class"
      :preview="false"
      :maxLength="contentLength"
      @on-change="change"
      @on-upload-img="uploadImg" />
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
  'update:content',
  'update:imageList',
])

const value = ref(props.modelValue)

const toolbars = ref([
  'bold',
  'underline',
  'italic',
  '-',
  'title',
  'strikeThrough',
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

<style lang="scss" scoped></style>
