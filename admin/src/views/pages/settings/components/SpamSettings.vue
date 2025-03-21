<template>
  <a-form :model="config" auto-label-width>
    <a-form-item label="Post verification code">
      <a-tooltip content="Whether to enable verification code verification when posting" placement="top">
        <a-switch v-model="config.topicCaptcha" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Post after email verification">
      <a-tooltip content="You need to verify your email address before you can post" placement="top">
        <a-switch v-model="config.createTopicEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Publish the article after email verification">
      <a-tooltip content="You need to verify your email address before you can publish an article" placement="top">
        <a-switch v-model="config.createArticleEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Comment after email verification">
      <a-tooltip content="You need to verify your email address before you can post a comment" placement="top">
        <a-switch v-model="config.createCommentEmailVerified" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Article review">
      <a-tooltip content="Whether to enable review after publishing the article" placement="top">
        <a-switch v-model="config.articlePending" />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="User observation period (seconds)">
      <a-tooltip
        content="During the observation period, users cannot post topics, updates, etc. Setting this value to 0 means there is no observation period."
        placement="top"
      >
        <a-input-number
          v-model="config.userObserveSeconds"
          mode="button"
          :min="0"
          :max="720"
        />
      </a-tooltip>
    </a-form-item>

    <a-form-item label="Email Whitelist">
      <a-select
        v-model="config.emailWhitelist"
        style="width: 100%"
        multiple
        filterable
        allow-create
        default-first-option
        placeholder="Email Whitelist"
      />
    </a-form-item>

    <a-form-item>
      <a-button type="primary" :loading="loading" @click="submit"
        >Save</a-button
      >
    </a-form-item>
  </a-form>
</template>

<script setup lang="ts">
  const loading = ref(false);
  const config = reactive({
    topicCaptcha: undefined,
    createTopicEmailVerified: undefined,
    createArticleEmailVerified: undefined,
    createCommentEmailVerified: undefined,
    articlePending: undefined,
    userObserveSeconds: undefined,
    emailWhitelist: undefined,
  });
  const loadConfig = async () => {
    const ret = await axios.get<any, any>('/api/admin/sys-config/all');
    config.topicCaptcha = ret.topicCaptcha;
    config.createTopicEmailVerified = ret.createTopicEmailVerified;
    config.createArticleEmailVerified = ret.createArticleEmailVerified;
    config.createCommentEmailVerified = ret.createCommentEmailVerified;
    config.articlePending = ret.articlePending;
    config.userObserveSeconds = ret.userObserveSeconds;
    config.emailWhitelist = ret.emailWhitelist;
  };

  loadConfig();

  const submit = async () => {
    loading.value = true;
    try {
      await axios.post('/api/admin/sys-config/save', config);
      await loadConfig();
      useNotificationSuccess('Submit successfully');
    } catch (e) {
      useHandleError(e);
    } finally {
      loading.value = false;
    }
  };
</script>

<style scoped lang="less"></style>
