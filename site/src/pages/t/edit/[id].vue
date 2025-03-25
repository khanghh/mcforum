<template>
  <section class="main">
    <div class="container">
      <div v-if="postForm" class="topic-create-form">
        <div class="topic-form-title">
          {{ $t('page.edit_topic') }}
        </div>

        <div class="field">
          <div class="control">
            <div
              v-for="forum in forums"
              :key="forum.id"
              class="topic-tag"
              :class="{ selected: postForm.forumId === forum.id }"
              @click="postForm.forumId = forum.id">
              <span>{{ forum.name }}</span>
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
const i18n = useI18n()
definePageMeta({
  middleware: 'auth',
})

useHead({
  title: useSiteTitle(i18n.t('page.edit_topic')),
})

const route = useRoute()
const publishing = ref(false)
const configStore = useConfigStore()

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent
})

const { data: forums } = useAsyncData('forums', () =>
  useMyFetch('/api/forums/list'),
)

const { data: postForm } = useAsyncData(() => {
  publishing.value = false
  return useMyFetch(`/api/topics/edit/${route.params.id}`)
})

async function submitCreate() {
  if (publishing.value) {
    return
  }
  publishing.value = true

  try {
    useHttpPostForm(`/api/topic/edit/${postForm.value.id}`, {
      body: {
        forumId: postForm.value.forumId,
        title: postForm.value.title,
        content: postForm.value.content,
        hideContent: postForm.value.hideContent,
        tags: postForm.value.tags ? postForm.value.tags.join(',') : '',
      },
    })
    useMsg({
      message: '修改成功',
      onClose() {
        useLinkTo(`/topic/${postForm.value.id}`)
      },
    })
  }
  catch (e) {
    publishing.value = false
    useMsgError('提交失败：' + (e.message || e))
  }
}
</script>

<style lang="scss" scoped></style>
