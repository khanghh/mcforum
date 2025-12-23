<template>
  <section class="main">
    <div class="container">
      <article v-if="isNeedEmailVerify" class="message is-warning">
        <div class="message-header">
          <p>{{ $t('publish.verify_email_title') }}</p>
        </div>
        <div class="message-body">
          <i18n-t keypath="publish.verify_email_action" tag="p">
            <strong>
              <nuxt-link
                to="/users/me/account"
                style="color: var(--text-link-color)">
                {{ $t('navbar.profile') }} &gt; {{ $t('navbar.settings')
                }}
              </nuxt-link>
            </strong>
          </i18n-t>
        </div>
      </article>
      <div v-else class="topic-create-form">
        <div class="topic-form-title">
          {{ postForm.type === 0 ? $t('page.create_topic') : $t('page.post_status') }}
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

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <input
              v-model="postForm.title"
              class="input topic-title"
              type="text"
              :placeholder="$t('form.placeholder.enter_post_content')" />
          </div>
        </div>

        <div v-if="postForm.type === 0" class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.content"
              :placeholder="$t('form.placeholder.enter_post_content')" />
          </div>
        </div>

        <div v-if="postForm.type === 0 && isEnableHideContent" class="field">
          <div class="control">
            <markdown-editor
              v-model="postForm.hideContent"
              height="200px"
              :placeholder="$t('form.placeholder.enter_hidden_content')" />
          </div>
        </div>

        <div v-if="postForm.type === 1" class="field">
          <div class="control">
            <simple-editor
              ref="simpleEditorComponent"
              v-model:content="postForm.content"
              v-model:image-list="postForm.imageList" />
          </div>
        </div>

        <div class="field">
          <div class="control">
            <tag-input v-model="postForm.tags" />
          </div>
        </div>

        <div v-if="postForm.captchaUrl" class="field is-horizontal">
          <div class="field control has-icons-left">
            <input
              v-model="postForm.captchaCode"
              class="input"
              type="text"
              :placeholder="$t('form.placeholder.enter_captcha')"
              style="max-width: 150px; margin-right: 20px" />
            <span class="icon is-small is-left">
              <icon name="ShieldCheck" />
            </span>
          </div>
          <div class="field">
            <a @click="showCaptcha">
              <img :src="postForm.captchaUrl" style="height: 40px" />
            </a>
          </div>
        </div>

        <div class="field is-grouped">
          <div class="control">
            <a
              :class="{ 'is-loading': publishing }"
              class="button is-success"
              @click="createTopic">
              {{ postForm.type === 1 ? $t('form.button.post_status')
                : $t('form.button.post_topic') }}
            </a>
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

const userStore = useUserStore()
const configStore = useConfigStore()
const route = useRoute()
const router = useRouter()

const forumId = parseInt(route.query.forumId) || configStore.config.defaultForumId || 0

const postForm = ref({
  type: Number.parseInt(route.query.type) || 0,
  forumId: forumId,
  title: '',
  tags: [],
  content: '',
  hideContent: '',
  imageList: [],

  captchaId: '',
  captchaUrl: '',
  captchaCode: '',
})
const publishing = ref(false)
const simpleEditorComponent = ref(null)

const isNeedEmailVerify = computed(() => {
  return (
    configStore.config.createTopicEmailVerified && !userStore.user.emailVerified
  )
})

const isEnableHideContent = computed(() => {
  return configStore.config.enableHideContent
})

const topicCaptchaEnabled = computed(() => {
  return configStore.config.topicCaptcha
})

const { data: forums } = useAsyncData('forums', () =>
  useHttpGet('/api/forums'),
)

init()

watch(
  () => route.query,
  (newQuery, oldQuery) => {
    // console.log(newQuery, oldQuery);

    init()
  },
  { deep: true },
)

onMounted(() => {
  showCaptcha()
})

function init() {
  postForm.value.type = Number.parseInt(route.query.type) || 0
  useHead({
    title: postForm.value.type === 0 ? i18n.t('page.create_topic') : i18n.t('page.post_status'),
  })
}

async function createTopic() {
  if (publishing.value) {
    return
  }

  if (postForm.value.type === 1) {
    if (simpleEditorComponent.value.loading) {
      useMsgWarning('图片上传中,请稍后重试...')
      return
    }
  }

  publishing.value = true
  try {
    const topic = await useHttpPost('/api/topics', {
      body: postForm.value,
    })
    router.push(`/topics/${topic.slug}`)
  } catch (e) {
    showCaptcha()
    useMsgError(e.message || e)
    publishing.value = false
  }
}

async function showCaptcha() {
  if (topicCaptchaEnabled.value) {
    try {
      const ret = await useHttpGet('/api/captcha/request', {
        params: {
          captchaId: postForm.value.captchaId || '',
        },
      })
      postForm.value.captchaId = ret.captchaId
      postForm.value.captchaUrl = ret.captchaUrl
    } catch (e) {
      useMsgError(e.message || e)
    }
  }
}
</script>

<style lang="scss" scoped></style>
