<script setup>
import { useToast } from '@/composables/useToast'
import { useConfirmDialog } from '@/composables/useConfirmDialog'

const i18n = useI18n()
const props = defineProps({
  modelValue: {
    type: Array,
    default() {
      return []
    },
  },
  accept: {
    type: String,
    default: 'image/*',
  },
  limit: {
    type: Number,
    default: 9,
  },
  sizeLimit: {
    type: Number,
    default: 1024 * 1024 * 20,
  },
  size: {
    type: String,
    default: '94px',
  },
})
const emits = defineEmits(['update:modelValue'])
const fileList = ref(props.modelValue)
const previewFiles = ref([])
const currentInput = ref(null)
const loading = ref(false)
const toast = useToast()
const confirmDialog = useConfirmDialog()

function onClick() {
  if (currentInput.value) {
    currentInput.value.dispatchEvent(new MouseEvent('click'))
  }
}

function onInput(e) {
  const files = e.target.files
  addFiles(files)
}

function addFiles(files) {
  if (!files || !files.length) return // 没有文件
  if (!checkSizeLimit(files)) return // 文件大小检查
  if (!checkLengthLimit(files)) return // 文件数量检查

  const fileArray = []
  for (let i = 0; i < files.length; i++) {
    const url = getObjectURL(files[i])
    previewFiles.value.push({
      name: files[i].name,
      url,
      progress: 0,
      deleted: false,
      size: files[i].size,
    })
    fileArray.push(files[i])
  }
  const promiseList = fileArray.reduce((result, file, index, array) => {
    result.push(uploadFile(file, index, array.length))
    return result
  }, [])
  uploadFiles(promiseList)
}

function uploadFile(file, index, length) {
  //   const config = {
  //     onUploadProgress(progressEvent) {
  //       if (!progressEvent.lengthComputable) {
  //         // 当进度不可估量,直接等于 100
  //         previewFiles.value[
  //           previewFiles.value.length - length + index
  //         ].progress = 100;
  //         return;
  //       }
  //       previewFiles.value[previewFiles.value.length - length + index].progress
  //         = Number.parseInt(
  //           Math.round(
  //             (progressEvent.loaded / progressEvent.total) * 100,
  //           ).toString(),
  //         ) * 0.9;
  //     },
  //   };
  const formData = new FormData()
  formData.append('image', file, file.name)
  return useHttp('/api/upload', {
    method: 'POST',
    body: formData,
  })
}

function uploadFiles(promiseList) {
  loading.value = true

  Promise.all(promiseList).then(
    (resList) => {
      // 请求响应后，更新到 100%
      previewFiles.value.forEach((item) => {
        item.progress = 100
      })
      resList.forEach((item) => {
        fileList.value.push(item)
      })
      if (currentInput.value) {
        currentInput.value.value = ''
      }
      loading.value = false
      emits('update:modelValue', fileList.value)
    },
    (e) => {
      if (currentInput.value) {
        currentInput.value.value = ''
      }

      // 失败的时候取消对应的预览照片
      const length = promiseList.length
      previewFiles.value.splice(previewFiles.value.length - length, length)
      console.error(e)

      loading.value = false
    },
  )
}

function removeItem(index) {
  confirmDialog.show({
    title: i18n.t('dialog.title.prompt'),
    message: i18n.t('dialog.message.confirm_action_post'),
    confirmText: i18n.t('dialog.button.confirm'),
    cancelText: i18n.t('dialog.button.cancel'),
    variant: 'warning',
    onConfirm: () => {
      previewFiles.value[index].deleted = true // 删除动画
      fileList.value.splice(index, 1)
      emits('update:modelValue', fileList.value) // 避免和回显冲突，先修改 fileList
      setTimeout(() => {
        previewFiles.value.splice(index, 1)
        toast.success(i18n.t('message.delete_success'))
      }, 900)
    },
    onCancel: () => console.log('canceled delete'),
  })
}

function checkSizeLimit(files) {
  let pass = true
  for (let i = 0; i < files.length; i++) {
    if (files[i].size > props.sizeLimit) {
      pass = false
    }
  }
  if (!pass)
    toast.error(i18n.t('message.image_size_limit_error', { limit: `${props.sizeLimit / 1024 / 1024} MB` }))
  return pass
}
function checkLengthLimit(files) {
  if (previewFiles.value.length + files.length > props.limit) {
    toast.warning(i18n.t('message.image_upload_limit', { limit: props.limit }))
    return false
  }
  else {
    return true
  }
}
function getObjectURL(file) {
  let url = null
  if (window.createObjectURL) {
    // basic
    url = window.createObjectURL(file)
  }
  else if (window.URL) {
    // mozilla(firefox)
    url = window.URL.createObjectURL(file)
  }
  else if (window.webkitURL) {
    // webkit or chrome
    url = window.webkitURL.createObjectURL(file)
  }
  return url
}
function clear() {
  fileList.value = []
  previewFiles.value = []
}

defineExpose({
  onClick,
  clear,
  loading,
})
</script>

<template>
  <div class="flex flex-wrap gap-3">
    <div
      v-for="(image, index) in previewFiles"
      :key="index"
      class="relative group rounded-xl overflow-hidden border border-purple-500/30 bg-gradient-to-br from-gray-900/80 to-gray-800/80"
      :class="{ 'transition-all duration-1000 transform -translate-y-full opacity-0': image.deleted }"
      :style="{ width: size, height: size }">
      <img :src="image.url"
        class="w-full h-full object-cover opacity-80 group-hover:opacity-100 transition-opacity duration-300">

      <!-- Progress Bar -->
      <div v-show="image.progress < 100"
        class="absolute inset-0 flex flex-col justify-center items-center bg-black/60 z-10 p-2">
        <div class="w-full h-1.5 bg-gray-700 rounded-full overflow-hidden mb-2">
          <div
            class="h-full bg-gradient-to-r from-purple-500 to-pink-500 rounded-full transition-all duration-300"
            :style="{ width: image.progress + '%' }">
          </div>
        </div>
        <span class="text-xs text-purple-200 font-medium">Uploading...</span>
      </div>

      <!-- Delete Overlay -->
      <div
        class="absolute inset-0 flex items-center justify-center bg-black/60 opacity-0 transition-opacity duration-200"
        :class="{ 'opacity-100': image.progress === 100 && !image.deleted, 'group-hover:opacity-100': image.progress === 100 }">
        <div
          class="cursor-pointer p-2 rounded-full bg-red-500/20 text-red-400 hover:bg-red-500 hover:text-white transition-colors duration-200"
          @click="removeItem(index)">
          <icon name="trash" size="20px" />
        </div>
      </div>
    </div>

    <!-- Add Button -->
    <div
      v-show="previewFiles.length < limit"
      class="relative flex items-center justify-center rounded-xl border-2 border-dashed border-purple-500/30 bg-gradient-to-br from-purple-500/5 to-pink-500/5 hover:from-purple-500/10 hover:to-pink-500/10 hover:border-purple-400 transition-all duration-300 cursor-pointer group"
      :style="{ width: size, height: size }"
      @click="onClick($event)">
      <input ref="currentInput" :accept="accept" type="file" multiple class="hidden" @input="onInput">
      <div class="flex flex-col items-center justify-center">
        <slot name="add-image-button">
          <Icon name="TablerPlus"
            class="text-purple-400 group-hover:scale-110 transition-transform duration-300 drop-shadow-[0_0_8px_rgba(168,85,247,0.5)]"
            size="24px" />
        </slot>
      </div>
    </div>
  </div>
</template>
