<template>
  <section class="main">
    <div class="container" v-if="postForm">
      <div class="article-create-form">
        <h1 class="title">{{ $t('page.edit_article') }}</h1>

        <div class="field">
          <div class="control">
            <div
              v-for="node in nodes"
              :key="node.id"
              class="article-tag"
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
              class="input article-title"
              type="text"
              :placeholder="$t('form.placeholder.enter_post_title')" />
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
              placeholder="隐藏内容，评论后可见" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <button
              class="button is-success"
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
  title: useSiteTitle(i18n.t('page.edit_article')),
});

const route = useRoute();
const configStore = useConfigStore();

const publishing = ref(false);

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent;
});

const { data: nodes } = useAsyncData("nodes", () =>
  useMyFetch("/api/article/nodes")
);

const { data: postForm } = useAsyncData(() => {
  publishing.value = false;
  return useMyFetch(`/api/article/edit/${route.params.id}`)
});

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;

  try {
    useHttpPostForm(`/api/article/edit/${postForm.value.id}`, {
      body: {
        title: postForm.value.title,
        content: postForm.value.content,
        cover: postForm.value.cover,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
      },
    });
    useMsg({
      message: i18n.t('alert.edit_post_success'),
      onClose() {
        useLinkTo(`/article/${postForm.value.id}`);
      },
    });
  } catch (e) {
    publishing.value = false;
    useMsgError(i18n.t('alert.edit_post_success', { error: (e.message || e) }));
  }
}
</script>

<style lang="scss" scoped>
.article-create-form {
  background-color: var(--bg-color);
  padding: 30px;

  .article-form-title {
    font-size: 36px;
    font-weight: 500;
    margin-bottom: 10px;
  }

  .field {
    margin-bottom: 10px;

    input {
      &:focus-visible {
        outline-width: 0;
      }
    }
  }
}

.cover-add-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  i {
    font-size: 24px;
    color: #1878f3;
  }

  span {
    font-size: 14px;
    color: var(--text-color3);
  }
}
</style>
