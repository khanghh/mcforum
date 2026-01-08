<template>
  <div class="w-full">
    <div v-if="empty" class="text-center pt-8">
      <slot name="empty" :items="pageItems">
        <div class="text-gray-500 text-sm font-medium">
          {{ $t('message.no_data_available') }}
        </div>
      </slot>
    </div>
    <template v-else>
      <slot name="default" :items="pageItems" />

      <!-- Custom Loading Skeleton -->
      <div v-if="loading" class="py-4 space-y-4">
        <div class="w-full p-6 rounded-xl bg-gray-800/40 border border-purple-500/20 animate-pulse">
          <div class="flex items-center gap-4 mb-4">
            <div class="w-10 h-10 rounded-lg bg-purple-500/20"></div>
            <div class="flex-1">
              <div class="h-4 bg-purple-500/20 rounded w-1/4 mb-2"></div>
              <div class="h-3 bg-purple-500/20 rounded w-1/6"></div>
            </div>
          </div>
          <div class="space-y-2">
            <div class="h-4 bg-purple-500/20 rounded w-full"></div>
            <div class="h-4 bg-purple-500/20 rounded w-5/6"></div>
            <div class="h-4 bg-purple-500/20 rounded w-4/6"></div>
          </div>
        </div>
      </div>

      <div class="text-center pt-8">
        <button
          v-if="hasMore"
          class="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold font-['Saira Semi Condensed'] tracking-wider shadow-[0_0_10px_rgba(139,92,246,0.5),0_0_20px_rgba(139,92,246,0.3)] hover:scale-105 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center gap-2 mx-auto min-w-[200px]"
          :disabled="disabled"
          @click="loadMore">
          <Icon v-if="loading" name="spinner" class="animate-spin" />
          <span>{{ loading ? 'LOADING...' : $t('actions.view_more') }}</span>
          <Icon v-if="!loading" name="TablerChevronDown" />
        </button>
        <div v-else class="text-gray-500 text-sm font-medium py-4">
          {{ $t('message.no_more_content') }}
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { CursorResult } from '@/composables/api'

type Props = {
  cursor: CursorResult<any[]>
}
const props = defineProps<Props>()
defineExpose({
  unshiftResults,
  removeItem
})

const loading = ref(false)
const pageItems = ref<any[]>([])
const hasMore = ref<boolean>(true)

const preloadItems = await props.cursor.loadNext().catch(() => [])
pageItems.value = preloadItems
hasMore.value = props.cursor.hasMore

const disabled = computed(() => {
  return loading.value || !hasMore
})

const empty = computed(() => {
  return !hasMore.value && !pageItems.value?.length
})

function unshiftResults(item: any) {
  if (item && pageItems.value) {
    pageItems.value.unshift(item)
  }
}

function removeItem(item: any) {
  if (item && pageItems.value) {
    const index = pageItems.value.indexOf(item)
    if (index !== -1) {
      pageItems.value.splice(index, 1)
    }
  }
}

async function loadMore() {
  loading.value = true
  try {
    const items = await props.cursor.loadNext()
    pageItems.value.push(...items)
    hasMore.value = props.cursor.hasMore
  } catch (err) {
    console.error('Load more failed:', err)
  } finally {
    loading.value = false
  }
}



</script>
