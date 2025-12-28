<template>
  <div v-if="editing" class="field">
    <div class="control">
      <input v-model="topic.title"
        class="input topic-title"
        type="text"
        :placeholder="$t('form.placeholder.enter_post_title')" />
    </div>
    <div class="control">
      <MarkdownEditor
        v-model="topic.content"
        :placeholder="$t('form.placeholder.enter_post_content')" />
    </div>
    <div class="control">
      <tag-input v-model="topic.tags" />
    </div>
  </div>
  <div v-else>
    <h1 class="text-3xl font-bold mb-6 text-white gaming-title" itemprop="headline">
      <span v-if="topic.sticky" class="text-red-500 mr-2" title="Sticky Post">
        <Icon name="TablerPin" />
      </span>
      <span v-if="topic.recommend" class="text-yellow-500 mr-2" title="Recommended">
        <Icon name="TablerStar" />
      </span>
      {{ topic.title }}
    </h1>

    <div class="topic-content-detail line-numbers" v-html="topic.content" />
    <ul v-if="topic.imageList && topic.imageList.length" class="topic-image-list">
      <li v-for="(image, index) in topic.imageList" :key="index">
        <div class="image-item">
          <el-image
            :src="image.preview"
            :preview-src-list="imageUrls"
            :initial-index="index" />
        </div>
      </li>
    </ul>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  editing: {
    type: Boolean,
  },
})
const topic = ref(props.modelValue)
</script>
