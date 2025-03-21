<template>
  <div class="container">
    <div class="container-header">
      <a-form :model="filters" layout="inline" :size="appStore.table.size">
        <a-form-item>
          <a-input v-model="filters.id" placeholder="UserID" />
        </a-form-item>
        <a-form-item>
          <a-input v-model="filters.username" placeholder="Username" />
        </a-form-item>
        <a-form-item>
          <a-input v-model="filters.nickname" placeholder="Nickname" />
        </a-form-item>
        <a-form-item>
          <a-select v-model="filters.type" placeholder="User Type" allow-clear>
            <a-option :value="0" label="User" />
            <a-option :value="1" label="Staff" />
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" @click="list">
            <template #icon> <icon-search /> </template>
            Search
          </a-button>
        </a-form-item>
      </a-form>

      <div class="action-btns">
        <!-- <a-button type="primary" :size="appStore.table.size" @click="showAdd">
          <template #icon>
            <icon-plus />
          </template>
          Add
        </a-button> -->
      </div>
    </div>
    <div class="container-main">
      <a-table
        :loading="loading"
        :data="data.results"
        :size="appStore.table.size"
        :bordered="appStore.table.bordered"
        :pagination="pagination"
        :sticky-header="true"
        style="height: 100%"
        column-resizable
        @page-change="onPageChange"
        @page-size-change="onPageSizeChange"
      >
        <template #columns>
          <a-table-column title="ID" data-index="id"></a-table-column>
          <a-table-column title="Type" data-index="type">
            <template #cell="{ record }">
              <a-tag v-if="record.type === 1" color="blue">Staff</a-tag>
              <a-tag v-else>User</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="Avatar" data-index="avatar">
            <template #cell="{ record }">
              <a-avatar>
                <img v-if="record.avatar" :src="record.avatar" />
                <span v-else>{{ record.nickname }}</span>
              </a-avatar>
            </template>
          </a-table-column>
          <a-table-column title="Nickname" data-index="nickname"></a-table-column>
          <a-table-column title="Email" data-index="email"></a-table-column>
          <a-table-column title="Score" data-index="score"></a-table-column>
          <a-table-column title="Fobidden" data-index="forbidden">
            <template #cell="{ record }">
              {{ record.forbidden ? 'Mute' : '-' }}
            </template>
          </a-table-column>
          <a-table-column title="Registration Time" data-index="createTime">
            <template #cell="{ record }">
              {{ useFormatDate(record.createTime) }}
            </template>
          </a-table-column>
          <a-table-column title="Actions">
            <template #cell="{ record }">
              <a-button
                type="primary"
                :size="appStore.table.size"
                @click="showEdit(record.id)"
                >Edit</a-button
              >
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <Edit ref="edit" @ok="list" />
  </div>
</template>

<script setup lang="ts">
  import Edit from './components/Edit.vue';

  const appStore = useAppStore();
  const loading = ref(false);
  const edit = ref();
  const filters = reactive({
    limit: 20,
    page: 1,

    id: undefined,
    username: undefined,
    nickname: undefined,
    type: undefined,
  });

  const data = reactive({
    page: {
      page: 1,
      limit: 20,
      total: 0,
    },
    results: [],
  });

  const pagination = computed(() => {
    return {
      total: data.page.total,
      current: data.page.page,
      pageSize: data.page.limit,
      showTotal: true,
      showJumper: true,
      showPageSize: true,
      pageSizeOptions: [20, 50, 100, 200, 300, 500],
    };
  });

  onMounted(() => {
    useTableHeight();
  });

  const list = async () => {
    loading.value = true;
    try {
      const ret = await axios.postForm<any>(
        '/api/admin/user/list',
        jsonToFormData(filters)
      );
      data.page = ret.page;
      data.results = ret.results;
    } finally {
      loading.value = false;
    }
  };

  list();

  const showAdd = () => {
    edit.value.show();
  };

  const showEdit = (id: any) => {
    edit.value.showEdit(id);
  };

  const onPageChange = (page: number) => {
    filters.page = page;
    list();
  };

  const onPageSizeChange = (pageSize: number) => {
    filters.limit = pageSize;
    list();
  };
</script>

<style scoped lang="less"></style>
