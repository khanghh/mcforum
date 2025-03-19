<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div>
        <icon name="Settings" />
        <span>{{ $t('page.personal_info') }}</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <icon name="Undo2" />
        <span>{{ $t('links.back_to_profile') }}</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <!-- 头像 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t('form.label.avatar') }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <avatar-edit :value="user.avatar" />
            </div>
          </div>
        </div>
      </div>

      <!-- 昵称 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t('form.label.nickname') }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.nickname"
                class="input"
                type="text"
                autocomplete="off"
                :placeholder="$t('form.placeholder.enter_nickname')" />
            </div>
          </div>
        </div>
      </div>

      <!-- 个性签名 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t('form.label.bio') }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <textarea
                v-model="form.description"
                class="textarea"
                rows="2"
                :placeholder="$t('form.placeholder.write_bio')" />
            </div>
          </div>
        </div>
      </div>

      <!-- 个人主页 -->
      <div class="field is-horizontal">
        <div class="field-label is-normal">
          <label class="label">{{ $t('form.label.website') }}</label>
        </div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <input
                v-model="form.homePage"
                class="input"
                type="text"
                autocomplete="off"
                :placeholder="$t('form.placeholder.enter_website')" />
            </div>
          </div>
        </div>
      </div>

      <div class="field is-horizontal">
        <div class="field-label is-normal"></div>
        <div class="field-body">
          <div class="field">
            <div class="control">
              <a class="button is-success" @click="submitForm">{{ $t('form.button.save') }}</a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n();
definePageMeta({
  middleware: ["auth"],
  layout: "profile",
});

useHead({
  title: useSiteTitle(i18n.t('page.personal_info')),
});

const userStore = useUserStore();
const user = computed(() => {
  return userStore.user;
});

const form = ref({
  nickname: "",
  avatar: "",
  homePage: "",
  description: "",
});

if (user.value != null) {
  form.value.nickname = user.value.nickname;
  form.value.avatar = user.value.avatar;
  form.value.homePage = user.value.homePage;
  form.value.description = user.value.description;
}

async function submitForm() {
  try {
    await useHttpPostForm(`/api/user/edit/${user.value.id}`, {
      body: form.value,
    });
    await reload();
    useMsgSuccess(i18n.t('message.profile_update_success'));
  } catch (e) {
    console.error(e);
    useMsgError(i18n.t('message.profile_update_failure', { error: (e.message || e) }));
  }
}
async function reload() {
  user.value = await useHttpGet("/api/user/current");
  form.value = { ...user.value };
}
</script>
<style lang="scss" scoped>
.widget-header {
  padding: 18px 0;

  &>div {
    display: flex;
    align-items: center;

    span {
      margin-left: 6px;
    }
  }

  &>a {
    display: flex;
    align-items: center;
    font-weight: 500;
    font-size: 12px;

    span {
      margin-left: 6px;
    }
  }
}

.field {
  margin-bottom: 10px;

  input,
  textarea {
    &:focus-visible {
      outline-width: 0;
    }
  }
}
</style>
