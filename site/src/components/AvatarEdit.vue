<template>
  <label
    class="relative inline-block group cursor-pointer"
    :style="{ width: size + 'px', height: size + 'px' }">
    <Avatar
      :username="username"
      :src="previewSrc || src"
      :size="size"
      :rounded="rounded"
      class="w-full h-full" />

    <div
      class="absolute inset-0 bg-black bg-opacity-40 flex items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-opacity"
      aria-hidden="true">
      <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M4 17v2a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2v-2M12 12v8M7 8l5-5 5 5M12 3v10" />
      </svg>
      <span class="ml-2 text-sm">Change</span>
    </div>

    <input
      type="file"
      accept="image/*"
      class="sr-only"
      @change="onFileChange" />
  </label>
</template>

<script setup>
import { ref, onUnmounted } from 'vue'
import Avatar from './Avatar.vue'

const props = defineProps({
  username: {
    type: String,
    default: '',
  },
  src: {
    type: String,
    default: '',
  },
  size: {
    type: [Number, String],
    default: 48,
  },
  rounded: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['uploaded'])

const previewSrc = ref('')
let objectUrl = null

function onFileChange(e) {
  const file = e.target.files && e.target.files[0]
  if (!file) return
  if (objectUrl) URL.revokeObjectURL(objectUrl)
  objectUrl = URL.createObjectURL(file)
  previewSrc.value = objectUrl
  emit('uploaded', file)
  e.target.value = ''
}

onUnmounted(() => {
  if (objectUrl) URL.revokeObjectURL(objectUrl)
})
</script>

<style scoped></style>
