<template>
  <div class="border-b border-purple-500/20">
    <div class="flex overflow-x-auto custom-scrollbar" role="tablist">
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="current === 'topics' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="select('topics')">
        <FontAwesome :icon="['fas', 'fire']" class="mr-2" /> POSTS
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="current === 'activity' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="select('activity')">
        <FontAwesome :icon="['fas', 'comments']" class="mr-2" /> ACTIVITIES
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="current === 'following' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="select('following')">
        <FontAwesome :icon="['fas', 'users']" class="mr-2" /> FOLLOWING
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="current === 'fans' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="select('fans')">
        <FontAwesome :icon="['fas', 'star']" class="mr-2" /> FANS
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="current === 'favorites' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="select('favorites')">
        <FontAwesome :icon="['fas', 'heart']" class="mr-2" /> FAVORITES
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps({
  modelValue: String,
})

const emit = defineEmits(['update:modelValue'])

const current = computed({
  get: () => props.modelValue,
  set: (v: string) => emit('update:modelValue', v),
})

function select(tab: string) {
  current.value = tab
}
</script>

<style scoped>
/* No local styles; component uses global utility classes */
</style>
