<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>请先验证邮箱</p>
        </div>
        <div class="message-body">
          发表话题前，请先前往
          <strong><nuxt-link
              to="/user/profile/account"
              style="color: var(--text-link-color)">个人中心 &gt; 账号设置</nuxt-link></strong>
          页面设置邮箱，并完成邮箱认证。
        </div>
      </article>
      <div v-else class="article-create-form">
        <div class="article-form-title">{{ $t('page.create_article') }}</div>
        <div class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input"
              type="text"
              :placeholder="$t('form.placeholder.enter_post_content')" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              :placeholder="$t('form.placeholder.enter_post_content')" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <image-upload v-model="postForm.cover" :limit="1" size="120px">
              <template #add-image-button>
                <div class="cover-add-btn">
                  <icon name="Plus" color="#1c71d8" size="50px" />
                  <span>{{ $t('publish.add_image') }}</span>
                </div>
              </template>
            </image-upload>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a v-if="publishing" :class="{ 'is-loading': publishing }" disabled class="button is-success">
              {{ $t('form.button.post_article') }}
            </a>
            <a v-else :class="{ 'is-loading': publishing }" class="button is-success" @click="submitCreate">
              {{ $t('form.button.post_article') }}
            </a>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const i18n = useI18n();
const publishing = ref(false); // 当前是否正处于发布中...
const postForm = ref({
  title: "",
  tags: [],
  content: "",
});

const userStore = useUserStore();
const configStore = useConfigStore();
const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createArticleEmailVerified &&
    !userStore.user.emailVerified
  );
});

useHead({
  title: useSiteTitle(i18n.t('page.create_article')),
});

definePageMeta({
  middleware: "auth",
});

if (!configStore.isEnabledArticle) {
  throw createError({
    statusCode: 403,
    message: "已关闭文章功能",
  });
}

async function submitCreate() {
  if (publishing.value) {
    return;
  }
  publishing.value = true;
  try {
    const article = await useHttpPostForm("/api/article/create", {
      body: {
        title: postForm.value.title,
        content: postForm.value.content,
        tags: postForm.value.tags ? postForm.value.tags.join(",") : "",
        cover:
          postForm.value.cover && postForm.value.cover.length
            ? JSON.stringify(postForm.value.cover[0])
            : null,
      },
    });
    useMsg({
      message: "提交成功",
      onClose() {
        useLinkTo(`/article/${article.id}`);
      },
    });
  } catch (e) {
    useMsgError(e.message || e);
    publishing.value = false;
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
