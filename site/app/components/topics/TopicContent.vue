<template>
  <article class="mb-6 text-gray-300">
    <h1 class="text-xl sm:text-3xl font-bold mb-6 text-white gaming-title" itemprop="headline">
      {{ topic.title }}
    </h1>

    <div
      class="prose prose-invert max-w-none topic-content"
      v-html="topic.content"></div>

    <HiddenContent :topic="topic" />

    <ImageList :images="topic.images" />
  </article>
</template>

<script setup lang="ts">
import { type Topic } from '~/types/topic'
import ImageList from '~/components/topics/ImageList.vue'
import HiddenContent from '~/components/topics/HiddenContent.vue'

const props = defineProps({
  topic: {
    type: Object as PropType<Topic>,
    required: true,
  },
})
</script>

<style scoped lang="scss">
.topic-content :deep(a) {
  color: theme('colors.blue.400');

  :hover {
    color: theme('colors.blue.300');
  }
}
</style>