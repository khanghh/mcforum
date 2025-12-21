<template>
  <div class="w-full">
    <div v-if="empty" class="text-center py-8">
      <div class="text-gray-500 text-sm font-medium py-4">
        {{ $t('message.no_data_available') }}
      </div>
    </div>
    <template v-else>
      <slot name="default" :items="pageData.items" />

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
          v-if="pageData.hasMore"
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

<script setup>
defineExpose({
  refresh,
  unshiftResults,
})

const props = defineProps({
  url: {
    type: String,
    required: true,
  },
  params: {
    type: Object,
    default() {
      return {}
    },
  },
})

const loading = ref(false)
const pageData = reactive({
  cursor: '',
  items: [],
  hasMore: true,
  firstLoaded: false,
})

const disabled = computed(() => {
  return loading.value || !pageData.hasMore
})

const empty = computed(() => {
  return pageData.firstLoaded && pageData.hasMore === false && pageData.items.length === 0
})

const { data: first } = await useAsyncData(
  `load-more:${props.url}:${JSON.stringify(props.params)}`,
  async () => {
    try {
      const data = await useHttpGet(props.url, {
        params: props.params || {},
      })
      return data
    } catch (e) {
      console.error(e)
      return null
    }
  },
  {
    watch: [() => props.url, () => props.params],
  },
)

if (first.value) {
  renderData(first.value)
  pageData.firstLoaded = true
}

async function loadMore() {
  loading.value = true
  try {
    const filters = Object.assign(props.params || {}, {
      cursor: pageData.cursor || '',
    })
    const data = await useHttpGet(props.url, {
      params: filters,
    })
    renderData(data)
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function refresh() {
  pageData.cursor = ''
  pageData.items = []
  pageData.hasMore = true
  return loadMore()
}

function renderData(data) {
  data = data || {}
  pageData.cursor = data.cursor
  pageData.hasMore = data.hasMore

  if (data.items && data.items.length) {
    data.items.forEach((item) => {
      pageData.items.push(item)
    })
  }
}

function unshiftResults(item) {
  if (item && pageData && pageData.items) {
    pageData.items.unshift(item)
  }
}
</script>
