<template>
  <div v-if="user" class="widget no-margin">
    <div class="widget-header">
      <div class="account">
        <icon name="Settings" />
        <span>{{ $t('page.account_settings') }}</span>
      </div>
      <nuxt-link :to="'/user/' + user.id">
        <icon name="Undo2" />
        <span>{{ $t('links.back_to_profile') }}</span>
      </nuxt-link>
    </div>
    <div class="widget-content">
      <div class="settings">
        <div class="settings-item">
          <div class="settings-item-title">{{ $t('form.label.username') }}</div>
          <div class="settings-item-input">
            <div class="input-value">{{ user.username }}</div>
            <div class="action-box">
              <a v-if="!user.username" @click="showUsernameDialog">{{ $t('form.button.set') }}</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">{{ $t('form.label.email') }}</div>
          <div class="settings-item-input">
            <div class="input-value">
              <span>{{ user.email }}</span>
              <span
                v-if="user.emailVerified"
                style="margin-left: 4px; font-size: 80%">({{ $t('settings.email_verified') }})</span>
            </div>
            <div class="action-box">
              <a @click="showEmailDialog"> {{ $t('form.button.change') }} </a>
              <a v-if="user.email && !user.emailVerified"
                @click="requestEmailVerify">{{ $t('form.button.verify') }}</a>
            </div>
          </div>
        </div>

        <div class="settings-item">
          <div class="settings-item-title">{{ $t('form.label.password') }}</div>
          <div class="settings-item-input">
            <div class="input-value">
              {{ user.passwordSet ? $t('settings.password_set') : $t('profile.settings.password_not_set') }}
            </div>
            <div class="action-box">
              <a v-if="user.passwordSet" @click="showUpdatePasswordDialog">{{ $t('form.button.change') }}</a>
              <a v-else @click="showSetPasswordDialog">{{ $t('form.button.set') }}</a>
            </div>
          </div>
        </div>
      </div>
    </div>

    <my-dialog
      ref="usernameDialog"
      v-model:visible="usernameDialogVisible"
      :title="$t('dialog.title.setup_username')"
      :width="320"
      @ok="setUsername">
      <div style="padding: 30px 0">
        <input
          v-model="usernameForm.username"
          class="input is-small"
          type="text"
          :placeholder="$t('form.placeholder.enter_username')" />
      </div>
    </my-dialog>

    <my-dialog
      ref="emailDialog"
      v-model:visible="emailDialogVisible"
      :title="$t('dialog.title.setup_email')"
      :width="320"
      @ok="setEmail">
      <div style="padding: 30px 0">
        <input
          v-model="emailForm.email"
          class="input is-small"
          type="text"
          :placeholder="$t('form.placeholder.enter_email')" />
      </div>
    </my-dialog>

    <my-dialog
      ref="updatePasswordDialog"
      v-model:visible="updatePasswordDialogVisible"
      :title="$t('dialog.title.change_password')"
      :width="320"
      @ok="updatePassword">
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="updatePasswordForm.oldPassword"
            class="input is-small"
            type="password"
            :placeholder="$t('form.placeholder.enter_password')"
            @keydown.enter="updatePassword" />
          <span class="icon is-small is-left">
            <icon name="Lock" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="updatePasswordForm.password"
            class="input is-small"
            type="password"
            :placeholder="$t('form.placeholder.enter_new_password')"
            @keydown.enter="updatePassword" />
          <span class="icon is-small is-left">
            <icon name="Lock" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="updatePasswordForm.rePassword"
            class="input is-small"
            type="password"
            :placeholder="$t('form.placeholder.confirm_password')"
            @keydown.enter="updatePassword" />
          <span class="icon is-small is-left">
            <icon name="Lock" />
          </span>
        </div>
      </div>
    </my-dialog>

    <my-dialog
      ref="setPasswordDialog"
      v-model:visible="setPasswordDialogVisible"
      :title="$t('dialog.title.set_new_password')"
      :width="320"
      @ok="setPassword">
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="setPasswordForm.password"
            class="input is-small"
            type="password"
            :placeholder="$t('form.placeholder.enter_new_password')"
            @keydown.enter="setPassword" />
          <span class="icon is-small is-left">
            <icon name="Lock" />
          </span>
        </div>
      </div>
      <div class="field">
        <div class="control has-icons-left">
          <input
            v-model="setPasswordForm.rePassword"
            class="input is-small"
            type="password"
            :placeholder="$t('form.placeholder.confirm_password')"
            @keydown.enter="setPassword" />
          <span class="icon is-small is-left">
            <icon name="Lock" />
          </span>
        </div>
      </div>
    </my-dialog>
  </div>
</template>

<script setup>
const i18n = useI18n();
definePageMeta({
  middleware: ["auth"],
  layout: "profile",
});

useHead({
  title: useSiteTitle(i18n.t('page.account_settings')),
});

const { data: user, refresh: userRefresh } = await useAsyncData("user", () =>
  useMyFetch("/api/user/current")
);

const usernameDialog = ref(null);
const usernameDialogVisible = ref(false);
const usernameForm = reactive({
  username: user.value ? user.value.username : "",
});
function showUsernameDialog() {
  usernameDialog.value.show();
}
async function setUsername() {
  try {
    await useHttpPostForm("/api/user/set/username", {
      body: {
        username: usernameForm.username,
      },
    });
    await userRefresh();
    useMsgSuccess(i18n.t('alert.set_username_success'));
    usernameDialog.value.close();
  } catch (err) {
    useMsgError(i18n.t('alert.set_username_failure', { error: (err.message || err) }));
  }
}

const emailDialog = ref(null);
const emailDialogVisible = ref(false);
const emailForm = reactive({
  email: user.value ? user.value.email : "",
});
function showEmailDialog() {
  emailDialog.value.show();
}
async function setEmail() {
  try {
    await useHttpPostForm("/api/user/set/email", {
      body: {
        email: emailForm.email,
      },
    });
    await userRefresh();
    useMsgSuccess(i18n.t('alert.email_update_success'));
    emailDialog.value.close();
  } catch (err) {
    useMsgError(i18n.t('alert.email_update_failure', { error: (err.message || err) }));
  }
}

async function requestEmailVerify() {
  const loading = useLoading();
  try {
    await useHttpPost("/api/user/send_verify_email");
    useMsgSuccess(
      i18n.t('alert.verify_email_sent')
    );
  } catch (err) {
    useMsgError(i18n.t('alert.verify_email_failure', { error: (err.message || err) }));
  } finally {
    loading.close();
  }
}

const updatePasswordDialog = ref(null);
const updatePasswordDialogVisible = ref(false);
const updatePasswordForm = reactive({
  password: "",
  rePassword: "",
});
function showUpdatePasswordDialog() {
  updatePasswordDialog.value.show();
}
async function updatePassword() {
  try {
    await useHttpPostForm("/api/user/update/password", {
      body: updatePasswordForm,
    });
    await userRefresh();
    useMsgSuccess(i18n.t('alert.password_update_success'));
    updatePasswordDialog.value.close();
  } catch (err) {
    useMsgError(i18n.t('alert.password_update_failure', { error: (err.message || err) }));
  }
}

const setPasswordDialog = ref(null);
const setPasswordDialogVisible = ref(false);
const setPasswordForm = reactive({
  password: "",
  rePassword: "",
});
function showSetPasswordDialog() {
  setPasswordDialog.value.show();
}
async function setPassword() {
  try {
    await useHttpPostForm("/api/user/set/password", {
      body: setPasswordForm,
    });
    await userRefresh();
    useMsgSuccess(i18n.t('alert.password_update_success'));
    setPasswordDialog.value.close();
  } catch (err) {
    useMsgError(i18n.t('alert.password_update_failure', { error: (err.message || err) }));
  }
}
</script>
<style lang="scss" scoped>
.field {
  margin-bottom: 10px;

  input {
    &:focus-visible {
      outline-width: 0;
    }
  }
}

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

.settings {
  .settings-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 0;

    @media (max-width: 600px) {
      flex-direction: column;
      align-items: flex-start;

      .settings-item-title {
        margin-bottom: 6px;
      }
    }

    &:not(:last-child) {
      border-bottom: 1px solid var(--border-color);
    }

    .settings-item-title {
      width: 100px;
      color: var(--text-color2);
      font-size: 15px;
    }

    .settings-item-input {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: space-between;
      width: 100%;
      font-size: 14px;

      .input-value {
        flex: 1;
        color: var(--text-color3);
      }

      .action-box {
        display: flex;
        align-items: center;
        column-gap: 10px;

        a {
          color: var(--text-link-color);
          font-size: 12px;

          &:hover {
            color: var(--text-link-hover-color);
          }
        }
      }
    }
  }
}
</style>
