<template>
  <article class="mb-6 text-gray-300">
    <h1 class="text-2xl sm:text-3xl font-bold mb-6 text-white gaming-title" itemprop="headline">
      <span v-if="topic.pinned" class="text-red-500 mr-2" title="Sticky Post">
        <Icon name="TablerPin" />
      </span>
      <span v-if="topic.recommended" class="text-yellow-500 mr-2" title="Recommended">
        <Icon name="TablerStar" />
      </span>
      {{ topic.title }}
    </h1>

    <div class="prose prose-invert max-w-none " v-html="topic.content"></div>
    <ul v-if="topic.imageList && topic.imageList.length" class="topic-image-list">
      <li v-for="(image, index) in topic.imageList" :key="index">
        <div class="image-item">
          <!-- <el-image
            :src="image.preview"
            :preview-src-list="imageUrls"
            :initial-index="index" /> -->
        </div>
      </li>
    </ul>
  </article>
</template>

<script setup lang="ts">
import { type Topic } from '~/types/topic'

const props = defineProps({
  topic: {
    type: Object as PropType<Topic>,
    required: true,
  },
})

</script>