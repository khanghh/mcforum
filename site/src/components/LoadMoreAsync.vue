<template>
  <div class="w-full">
    <div v-if="empty" class="text-center py-8">
      <div class="text-gray-500 text-sm font-medium py-4">
        {{ $t('message.no_data_available') }}
      </div>
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

      <div class="text-center py-8">
        <button
          v-if="pageHasMore"
          class="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold font-['Saira Semi Condensed'] tracking-wider shadow-[0_0_10px_rgba(139,92,246,0.5),0_0_20px_rgba(139,92,246,0.3)] hover:scale-105 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:scale-100 flex items-center justify-center gap-2 mx-auto min-w-[200px]"
          :disabled="disabled"
          @click="loadMore">
          <FontAwesome v-if="loading" :icon="['fas', 'spinner']" spin />
          <span>{{ loading ? 'LOADING...' : $t('feed.actions.view_more') }}</span>
          <FontAwesome v-if="!loading" :icon="['fas', 'chevron-down']" />
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
import { ref, computed } from 'vue'


type Props = {
  cursor: CursorResult<any[]>
}

const props = defineProps<Props>()

const loading = ref(false)
const pageCursor = ref<any>('')
const pageItems = ref<any[]>([])
const pageHasMore = ref<boolean>(false)

// initialize page data from the passed cursor (do not auto-load)
pageItems.value = props.cursor.items ?? []
pageHasMore.value = props.cursor.hasMore ?? false
pageCursor.value = props.cursor.cursor

// expose API after functions are declared

const disabled = computed(() => loading.value || !pageHasMore.value)

const empty = computed(() => pageHasMore.value === false && pageItems.value.length === 0)

// Do not auto-load on mounted — render current `cursor.items` only.

async function loadMore() {
  loading.value = true
  try {
    const newItems = await props.cursor.loadNext()
    if (newItems && newItems.length) {
      newItems.forEach((it) => pageItems.value.push(it))
    }
    pageHasMore.value = props.cursor.hasMore
    pageCursor.value = props.cursor.cursor
  } catch (err: any) {
    throw createError({
      statusCode: err.statusCode || 500,
      statusMessage: err.message || 'Failed to load more items',
    })
  } finally {
    loading.value = false
  }
}

function unshiftResults(item: any) {
  if (item && pageItems.value) {
    pageItems.value.unshift(item)
  }
}

// expose helpers
defineExpose({
  unshiftResults,
})
</script>
