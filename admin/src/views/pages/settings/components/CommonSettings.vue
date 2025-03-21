<template>
  <a-form :model="config" auto-label-width>
    <a-form-item label="Website Title">
      <a-input 
        v-model="config.siteTitle" 
        type="text" 
        placeholder="Title displayed in the browser tab and metadata." 
      />
    </a-form-item>

    <a-form-item label="Website description">
      <a-textarea
        v-model="config.siteDescription"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        placeholder="Brief site description for SEO and metadata."
      />
    </a-form-item>

    <a-form-item label="Site keywords">
      <a-select
        v-model="config.siteKeywords"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="List of keywords for search engine optimization."
      />
    </a-form-item>

    <a-form-item label="System notice">
      <a-textarea
        v-model="config.siteNotification"
        :auto-size="{
          minRows: 2,
          maxRows: 5,
        }"
        placeholder=""
      />
    </a-form-item>

    <a-form-item label="Recommended tags">
      <a-select
        v-model="config.recommendTags"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="Predefined tags for content categorization."
      />
    </a-form-item>

    <a-form-item label="Default topic category">
      <a-select v-model="config.defaultNodeId" placeholder="Default category for new topics.">
        <a-option
          v-for="node in nodes"
          :key="node.id"
          :label="node.name"
          :value="node.id"
        />
      </a-select>
    </a-form-item>

    <a-form-item label="Enabled Modules">
      <a-checkbox v-model="config.modules.tweet" border>Tweet</a-checkbox>
      <a-checkbox v-model="config.modules.topic" border>Topic</a-checkbox>
      <a-checkbox v-model="config.modules.article" border>Article</a-checkbox>
    </a-form-item>

    <a-form-item label="External link jump page">
      <a-tooltip content="You need to manually confirm whether to go to the external link before jumping" placement="top">
        <a-switch v-model="config.urlRedirect" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Allow hidden content">
      <a-tooltip content="When posting, you can set the visible content after commenting" placement="top">
        <a-switch v-model="config.enableHideContent" />
      </a-tooltip>
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit">Save</a-button>
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  import { NodeDTO } from '@/composables/types';

  const loading = ref(false);
  const config = reactive({
    siteTitle: '',
    siteDescription: '',
    siteKeywords: [],
    siteNotification: '',
    recommendTags: [],
    defaultNodeId: undefined,
    urlRedirect: false,
    enableHideContent: false,
    modules: {
      tweet: false,
      topic: false,
      article: false,
    },
  });
  const nodes = ref<NodeDTO[]>([]);

  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/all');
    config.siteTitle = ret.siteTitle;
    config.siteDescription = ret.siteDescription;
    config.siteKeywords = ret.siteKeywords;
    config.siteNotification = ret.siteNotification;
    config.recommendTags = ret.recommendTags;
    config.defaultNodeId = ret.defaultNodeId;
    config.urlRedirect = ret.urlRedirect;
    config.enableHideContent = ret.enableHideContent;
    config.modules = ret.modules;
    nodes.value = await axios.get<any, NodeDTO[]>(
      '/api/admin/topic-node/nodes'
    );
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', config);
      await loadConfig();
      useNotificationSuccess('提交成功');
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style lang="scss" scoped></style>
