<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">{{ $t('page.signup') }}</div>
          <div class="widget-content">
            <div class="field">
              <label class="label">{{ $t('form.label.username') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.nickname" class="input is-success" type="text"
                  :placeholder="$t('form.placeholder.enter_username')" @keyup.enter="signup" />
                <span class="icon is-small is-left">
                  <icon name="UserRound" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">{{ $t('form.label.email') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.email" class="input is-success" type="text"
                  :placeholder="$t('form.placeholder.enter_email')"
                  @keyup.enter="signup" />
                <span class="icon is-small is-left">
                  <icon name="Mail" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">{{ $t('form.label.password') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.password" class="input" type="password"
                  :placeholder="$t('form.placeholder.enter_password')"
                  @keyup.enter="signup" />
                <span class="icon is-small is-left">
                  <icon name="Lock" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">{{ $t('form.label.confirm_password') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.rePassword" class="input" type="password"
                  :placeholder="$t('form.placeholder.confirm_password')" @keyup.enter="signup" />
                <span class="icon is-small is-left">
                  <icon name="Lock" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">{{ $t('form.label.captcha') }}</label>
              <div class="control has-icons-left">
                <div class="field is-horizontal">
                  <div class="field login-captcha-input">
                    <input v-model="form.captchaCode" class="input" type="text"
                      :placeholder="$t('form.placeholder.enter_captcha')"
                      @keyup.enter="signup" />
                    <span class="icon is-small is-left">
                      <icon name="ShieldCheck" />
                    </span>
                  </div>
                  <div v-if="form.captchaUrl" class="field login-captcha-img">
                    <a @click="refreshCaptcha"><img :src="form.captchaUrl" /></a>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">

              <div class="control">
                <button class="button is-link" @click="signup">{{ $t('form.button.signup') }}</button>
                <a class="button is-text" @click="toSignin">
                  {{ $t('form.link.already_have_account') }}&gt;&gt;
                </a>
              </div>
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
  title: useSiteTitle(i18n.t('page.signup')),
});

const route = useRoute();
const form = reactive({
  nickname: "",
  email: "",
  password: "",
  rePassword: "",
  captchaId: "",
  captchaUrl: "",
  captchaCode: "",
  redirect: route.query.redirect || "",
});

refreshCaptcha();

async function refreshCaptcha() {
  try {
    const { data: captcha } = await useAsyncData(() => {
      return useMyFetch("/api/captcha/request", {
        params: {
          captchaId: form.captchaId,
        },
      });
    });

    form.captchaId = captcha.value.captchaId;
    form.captchaUrl = captcha.value.captchaUrl;
    form.captchaCode = "";
  } catch (e) {
    useCatchError(e);
  }
}

async function signup() {
  try {
    const userStore = useUserStore();
    const { user, redirect } = await userStore.signup(form);
    if (redirect) {
      useLinkTo(redirect);
    } else {
      useLinkTo(`/user/${user.id}`);
    }
  } catch (err) {
    useCatchError(err);
    await refreshCaptcha();
  }
}

function toSignin() {
  if (form.redirect) {
    useLinkTo(`/user/signin?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signin");
  }
}
</script>

<style lang="scss" scoped></style>
