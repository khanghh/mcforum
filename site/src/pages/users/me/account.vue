<template>
  <div v-if="user" class="max-w-3xl mx-auto px-4 py-8">
    <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 text-gray-100">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-3">
          <icon name="Settings" class="text-gray-200" />
          <h2 class="text-lg font-semibold">{{ $t('page.account_settings') }}</h2>
        </div>

        <nuxt-link :to="'/users/' + user.id"
          class="inline-flex items-center gap-2 text-sm text-gray-300 hover:text-white">
          <icon name="Undo2" />
          <span>{{ $t('links.back_to_profile') }}</span>
        </nuxt-link>
      </div>

      <div class="space-y-4">
        <!-- Username -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-4 border-t border-gray-700">
          <div class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.username') }}</div>
          <div class="flex-1 flex items-center justify-between gap-4">
            <div class="text-sm text-gray-100">{{ user.username || '-' }}</div>
            <div class="flex items-center gap-3 text-sm">
              <button v-if="!user.username" class="text-indigo-400 hover:underline" @click="showUsernameDialog">
                {{
                  $t('form.button.set') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Email -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-4 border-t border-gray-700">
          <div class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.email') }}</div>
          <div class="flex-1 flex items-center justify-between gap-4">
            <div class="text-sm text-gray-100">
              <span>{{ user.email || '-' }}</span>
              <span v-if="user.emailVerified" class="ml-2 text-xs text-gray-400">
                ({{
                  $t('profile.settings.email_verified') }})
              </span>
            </div>
            <div class="flex items-center gap-3 text-sm">
              <button class="text-indigo-400 hover:underline" @click="showEmailDialog">
                {{ $t('form.button.change')
                }}
              </button>
              <button v-if="user.email && !user.emailVerified" class="text-indigo-400 hover:underline"
                @click="requestEmailVerify">
                {{ $t('form.button.verify') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Password -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 py-4 border-t border-gray-700">
          <div class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.password') }}</div>
          <div class="flex-1 flex items-center justify-between gap-4">
            <div class="text-sm text-gray-100">
              {{ user.passwordSet ? $t('profile.settings.password_set')
                : $t('profile.settings.password_not_set') }}
            </div>
            <div class="flex items-center gap-3 text-sm">
              <button v-if="user.passwordSet" class="text-indigo-400 hover:underline"
                @click="showUpdatePasswordDialog">
                {{ $t('form.button.change') }}
              </button>
              <button v-else class="text-indigo-400 hover:underline" @click="showSetPasswordDialog">
                {{
                  $t('form.button.set') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Dialogs -->
    <my-dialog
      ref="usernameDialog"
      v-model:visible="usernameDialogVisible"
      :title="$t('dialog.title.setup_username')"
      :width="320"
      @ok="setUsername">
      <div class="py-6">
        <input
          v-model="usernameForm.username"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400"
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
      <div class="py-6">
        <input
          v-model="emailForm.email"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400"
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
      <div class="space-y-3">
        <input v-model="updatePasswordForm.oldPassword"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400" type="password"
          :placeholder="$t('form.placeholder.enter_password')" @keydown.enter="updatePassword" />
        <input v-model="updatePasswordForm.password"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400" type="password"
          :placeholder="$t('form.placeholder.enter_new_password')" @keydown.enter="updatePassword" />
        <input v-model="updatePasswordForm.rePassword"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400" type="password"
          :placeholder="$t('form.placeholder.confirm_password')" @keydown.enter="updatePassword" />
      </div>
    </my-dialog>

    <my-dialog
      ref="setPasswordDialog"
      v-model:visible="setPasswordDialogVisible"
      :title="$t('dialog.title.set_new_password')"
      :width="320"
      @ok="setPassword">
      <div class="space-y-3">
        <input v-model="setPasswordForm.password"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400" type="password"
          :placeholder="$t('form.placeholder.enter_new_password')" @keydown.enter="setPassword" />
        <input v-model="setPasswordForm.rePassword"
          class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400" type="password"
          :placeholder="$t('form.placeholder.confirm_password')" @keydown.enter="setPassword" />
      </div>
    </my-dialog>
  </div>
</template>

<script setup>
const i18n = useI18n()
definePageMeta({
  middleware: ['auth'],
  layout: 'profile',
})

useHead({
  title: useSiteTitle(i18n.t('page.account_settings')),
})

const { data: user, refresh: userRefresh } = await useAsyncData('user', () =>
  useHttpGet('/api/users/me'),
)

const usernameDialog = ref(null)
const usernameDialogVisible = ref(false)
const usernameForm = reactive({
  username: user.value ? user.value.username : '',
})
function showUsernameDialog() {
  usernameDialog.value.show()
}
async function setUsername() {
  try {
    await useHttpPostForm('/api/users/set/username', {
      body: {
        username: usernameForm.username,
      },
    })
    await userRefresh()
    useMsgSuccess(i18n.t('message.set_username_success'))
    usernameDialog.value.close()
  }
  catch (err) {
    useMsgError(i18n.t('message.set_username_failure', { error: (err.message || err) }))
  }
}

const emailDialog = ref(null)
const emailDialogVisible = ref(false)
const emailForm = reactive({
  email: user.value ? user.value.email : '',
})
function showEmailDialog() {
  emailDialog.value.show()
}
async function setEmail() {
  try {
    await useHttpPostForm('/api/users/set/email', {
      body: {
        email: emailForm.email,
      },
    })
    await userRefresh()
    useMsgSuccess(i18n.t('message.email_update_success'))
    emailDialog.value.close()
  }
  catch (err) {
    useMsgError(i18n.t('message.email_update_failure', { error: (err.message || err) }))
  }
}

async function requestEmailVerify() {
  const loading = useLoading()
  try {
    await useHttpPost('/api/users/send_verify_email')
    useMsgSuccess(
      i18n.t('message.verify_email_sent'),
    )
  }
  catch (err) {
    useMsgError(i18n.t('message.verify_email_failure', { error: (err.message || err) }))
  }
  finally {
    loading.close()
  }
}

const updatePasswordDialog = ref(null)
const updatePasswordDialogVisible = ref(false)
const updatePasswordForm = reactive({
  password: '',
  rePassword: '',
})
function showUpdatePasswordDialog() {
  updatePasswordDialog.value.show()
}
async function updatePassword() {
  try {
    await useHttpPostForm('/api/users/update/password', {
      body: updatePasswordForm,
    })
    await userRefresh()
    useMsgSuccess(i18n.t('message.password_update_success'))
    updatePasswordDialog.value.close()
  }
  catch (err) {
    useMsgError(i18n.t('message.password_update_failure', { error: (err.message || err) }))
  }
}

const setPasswordDialog = ref(null)
const setPasswordDialogVisible = ref(false)
const setPasswordForm = reactive({
  password: '',
  rePassword: '',
})
function showSetPasswordDialog() {
  setPasswordDialog.value.show()
}
async function setPassword() {
  try {
    await useHttpPostForm('/api/users/set/password', {
      body: setPasswordForm,
    })
    await userRefresh()
    useMsgSuccess(i18n.t('message.password_update_success'))
    setPasswordDialog.value.close()
  }
  catch (err) {
    useMsgError(i18n.t('message.password_update_failure', { error: (err.message || err) }))
  }
}
</script>

<style scoped>
/* Tailwind handles visual styling; no additional CSS required. */
</style>
