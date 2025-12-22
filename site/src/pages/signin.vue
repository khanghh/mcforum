<template>
  <section class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="bg-white py-8 px-6 shadow-lg rounded-2xl">
        <div class="text-center mb-6">
          <h2 class="mt-2 text-2xl font-extrabold text-gray-900">{{ $t('page.signin') }}</h2>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700">{{ $t('form.label.username_or_email') }}</label>
            <div class="mt-1">
              <input v-model="form.username" type="text" @keyup.enter="signin"
                :placeholder="$t('form.placeholder.enter_email')"
                class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">{{ $t('form.label.password') }}</label>
            <div class="mt-1">
              <input v-model="form.password" type="password" @keyup.enter="signin"
                :placeholder="$t('form.placeholder.enter_password')"
                class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700">{{ $t('form.label.captcha') }}</label>
            <div class="mt-1 flex gap-3 items-center">
              <input v-model="form.captchaCode" type="text" @keyup.enter="signin"
                :placeholder="$t('form.placeholder.enter_captcha')"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" />

              <button v-if="form.captchaUrl" @click="refreshCaptcha" type="button" class="flex-shrink-0">
                <img :src="form.captchaUrl" alt="captcha" class="h-10 w-28 object-cover rounded-md border" />
              </button>
            </div>
          </div>

          <div class="pt-4">
            <button @click="signin"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none">
              {{ $t('form.button.signin') }}
            </button>
            <div class="mt-3 text-center">
              <button @click="toSignup" class="text-sm text-indigo-600 hover:underline">{{ $t('links.dont_have_account')
              }}</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
const i18n = useI18n()

useHead({
  title: useSiteTitle(i18n.t('page.signin')),
})

const route = useRoute()
const form = reactive({
  username: '',
  password: '',
  captchaId: '',
  captchaUrl: '',
  captchaCode: '',
  redirect: route.query.redirect || '',
})

refreshCaptcha()

async function refreshCaptcha() {
  try {
    const { data: captcha } = await useAsyncData(() => {
      return useHttpGet('/api/captcha/request', {
        params: {
          captchaId: form.captchaId,
        },
      })
    })
    form.captchaId = captcha.value.captchaId
    form.captchaUrl = captcha.value.captchaUrl
    form.captchaCode = ''
  } catch (e) {
    useCatchError(e)
  }
}

async function signin() {
  try {
    if (!form.username) {
      useMsgError(i18n.t('message.username_email_required'))
      return
    }
    if (!form.password) {
      useMsgError(i18n.t('message.password_required'))
      return
    }
    if (!form.captchaCode) {
      useMsgError(i18n.t('message.captcha_required'))
      return
    }

    const userStore = useUserStore()
    const { user, redirect } = await userStore.signin(form)
    if (redirect) {
      useLinkTo(redirect)
    } else {
      useLinkTo(`/user/${user.id}`)
    }
  } catch (e) {
    useCatchError(e)
    await refreshCaptcha()
  }
}

function toSignup() {
  if (form.redirect) {
    useLinkTo(`/signup?redirect=${encodeURIComponent(form.redirect)}`)
  } else {
    useLinkTo('/signup')
  }
}
</script>

<style scoped>
/* no custom styles — layout handled with Tailwind utilities */
</style>
