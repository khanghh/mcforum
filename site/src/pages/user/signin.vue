<template>
  <section class="main">
    <div class="container">
      <div class="main-body no-bg">
        <div class="widget signin">
          <div class="widget-header">{{ $t('page.signin') }}</div>
          <div class="widget-content">
            <div class="field">
              <label class="label">{{ $t('form.label.username_or_email') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.username" class="input is-success" type="text"
                  :placeholder="$t('form.placeholder.enter_email')" @keyup.enter="signin" />
                <span class="icon is-small is-left">
                  <icon name="UserRound" />
                </span>
              </div>
            </div>

            <div class="field">
              <label class="label">{{ $t('form.label.password') }}</label>
              <div class="control has-icons-left">
                <input v-model="form.password" class="input" type="password"
                  :placeholder="$t('form.placeholder.enter_password')"
                  @keyup.enter="signin" />
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
                      @keyup.enter="signin" />
                    <span class="icon is-small is-left">
                      <icon name="ShieldCheck" />
                    </span>
                  </div>
                  <div v-if="form.captchaUrl" class="field login-captcha-img" @click="refreshCaptcha">
                    <img :src="form.captchaUrl" data-not-lazy />
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <button class="button is-link" @click="signin">{{ $t('form.button.signin') }}</button>
              <a class="button is-text" @click="toSignup">
                {{ $t('links.dont_have_account') }}
              </a>
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
});

const route = useRoute();
const form = reactive({
  username: "",
  password: "",
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

async function signin() {
  try {
    if (!form.username) {
      useMsgError(i18n.t('message.username_email_required'));
      return;
    }
    if (!form.password) {
      useMsgError(i18n.t('message.password_required'));
      return;
    }
    if (!form.captchaCode) {
      useMsgError(i18n.t('message.captcha_required'));
      return;
    }

    const userStore = useUserStore();
    const { user, redirect } = await userStore.signin(form);
    if (redirect) {
      useLinkTo(redirect);
    } else {
      useLinkTo(`/user/${user.id}`);
    }
  } catch (e) {
    useCatchError(e);
    await refreshCaptcha();
  }
}

function toSignup() {
  if (form.redirect) {
    useLinkTo(`/user/signup?redirect=${encodeURIComponent(form.redirect)}`);
  } else {
    useLinkTo("/user/signup");
  }
}
</script>

<style lang="scss" scoped></style>
