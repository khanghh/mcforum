<template>
  <div class="space-y-4">
    <TopicCard
      v-for="topic in topics"
      :key="topic.id"
      :topic="topic"
      :show-pinned="showPinned" />
  </div>
</template>

<script setup>
import TopicCard from './TopicCard.vue'

defineProps({
  topics: {
    type: Array,
    default: () => [],
    required: false,
  },
  showAvatar: {
    type: Boolean,
    default: true,
  },
  showPinned: {
    type: Boolean,
    default: false,
  },
})
</script>
