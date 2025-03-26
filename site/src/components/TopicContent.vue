<template>
  <div v-if="editing" class="field">
    <div class="control">
      <input v-model="topic.title"
        class="input topic-title"
        type="text"
        :placeholder="$t('form.placeholder.enter_post_title')" />
    </div>
    <div class="control">
      <markdown-editor
        v-model="topic.content"
        :placeholder="$t('form.placeholder.enter_post_content')" />
    </div>
    <div class="control">
      <tag-input v-model="topic.tags" />
    </div>
  </div>
  <div v-else>
    <h1 v-if="topic.title" class="topic-title" itemprop="headline">
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
