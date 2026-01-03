<template>
  <div v-if="images && images.length" class="mt-6">
    <div class="grid grid-cols-6 sm:grid-cols-8 md:grid-cols-10 gap-2">
      <div
        v-for="(image, index) in images"
        :key="index"
        class="aspect-square cursor-pointer overflow-hidden rounded-lg border border-purple-500/20 bg-gray-800/50 transition-all duration-300 hover:border-purple-500/50 hover:scale-105"
        @click="openImageViewer(index)">
        <img
          :src="image.preview || image.url"
          class="w-full h-full object-cover"
          loading="lazy"
          alt="Topic image" />
      </div>
    </div>

    <ImageViewer
      v-model="showImageViewer"
      :images="images"
      v-model:index="previewIndex" />
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  images: {
    type: Array as PropType<any[]>,
    default: () => [],
  },
})

const showImageViewer = ref(false)
const previewIndex = ref(0)

const openImageViewer = (index: number) => {
  previewIndex.value = index
  showImageViewer.value = true
}
</script>
