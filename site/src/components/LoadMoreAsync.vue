<template>
  <div class="load-more">
    <slot v-if="empty" name="empty">
      <my-empty />
    </slot>
    <template v-else>
      <slot name="default" :items="pageData.items" />
      <div v-if="loading" class="loading">
        <el-skeleton :rows="3" animated />
      </div>
      <div class="has-more">
        <button class="button is-primary is-small" :disabled="disabled" @click="loadMore">
          <span v-if="loading" class="icon">
            <icon name="LoaderCircle" class="icon-loading" />
          </span>
          <span>{{ pageData.hasMore ? $t('message.load_more') : $t('message.no_more_content') }}</span>
        </button>
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
  // 请求URL
  url: {
    type: String,
    required: true,
  },
  // 请求参数
  params: {
    type: Object,
    default() {
      return {}
    },
  },
})
// 是否正在加载中
const loading = ref(false)
const pageData = reactive({
  cursor: '',
  items: [],
  hasMore: true,
})

const disabled = computed(() => {
  return loading.value || !pageData.hasMore
})

const empty = computed(() => {
  return pageData.hasMore === false && pageData.items.length === 0
})

const { data: first } = await useAsyncData(() => {
  return useMyFetch(props.url, {
    params: props.params || {},
  })
})

renderData(first.value)

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
  }
  catch (err) {
    console.error(err)
  }
  finally {
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

<style lang="scss" scoped>
.load-more {
  .loading {
    background-color: var(--bg-color);
    padding: 10px;
  }

  .has-more {
    text-align: center;
    padding: 20px;

    button {
      width: 150px;
    }
  }

  .no-more {
    text-align: center;
    padding: 10px 0;
    color: var(--text-color3);
    font-size: 14px;
  }

  .icon-loading {
    animation: rotating 3s infinite linear;
  }
}
</style>
