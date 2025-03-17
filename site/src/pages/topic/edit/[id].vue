<template>
  <section class="main">
    <div class="container">
      <div class="topic-create-form" v-if="postForm">
        <div class="topic-form-title">{{ $t('publish.edit_topic') }}</div>

        <div class="field">
          <div class="control">
            <div
              v-for="node in nodes"
              :key="node.id"
              class="topic-tag"
              :class="{ selected: postForm.nodeId === node.id }"
              @click="postForm.nodeId = node.id">
              <span>{{ node.name }}</span>
            </div>
          </div>
        </div>

        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input topic-title"
              type="text"
              :placeholder="$t('form.placeholder.enter_title')" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              :placeholder="$t('form.placeholder.enter_post_content')" />
          </div>
        </div>

        <div v-if="isEnableHideContent || postForm.hideContent" class="field">
          <div class="control">
            <markdown-editor
              ref="mdEditor"
              v-model="postForm.hideContent"
              height="200px"
              :placeholder="$t('form.placeholder.enter_hidden_content')" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <button class="button is-success"
              :class="{ 'is-loading': publishing }"
              :disabled="publishing"
              @click="submitCreate">
              {{ $t('form.button.submmit_changes') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const i18n = useI18n();
definePageMeta({
  middleware: "auth",
});

useHead({
  title: useSiteTitle(i18n.t('publish.edit_topic')),
});

const route = useRoute();
const publishing = ref(false);
const configStore = useConfigStore();

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const { data: nodes } = useAsyncData("nodes", () =>
  useMyFetch("/api/topic/nodes")
);

const { data: postForm } = useAsyncData(() => {
  publishing.value = false
  return useMyFetch(`/api/topic/edit/${route.params.id}`)
});


async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPostForm(`/api/topic/edit/${postForm.value.id}`, {
      body: {
        nodeId: postForm.value.nodeId,
        title: postForm.value.title,
        content: postForm.value.content,
        hideContent: postForm.value.hideContent,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      },
    });
    useMsg({
      message: "修改成功",
      onClose() {
        useLinkTo(`/topic/${postForm.value.id}`);
      },
    });
  } catch (e) {
    publishing.value = false;
    useMsgError("提交失败：" + (e.message || e));
  }
}

</script>

<style lang="scss" scoped></style>
