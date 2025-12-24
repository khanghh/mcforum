<template>
  <div v-if="user" class="max-w-3xl mx-auto px-4 py-8">
    <div class="bg-gray-800 border border-gray-700 rounded-2xl p-6 text-gray-100">
      <div class="flex items-center justify-between mb-6">
        <div class="flex items-center gap-3">
          <icon name="Settings" class="text-gray-200" />
          <h2 class="text-lg font-semibold">{{ $t('page.personal_info') }}</h2>
        </div>

        <nuxt-link :to="`/users/${user.username}`"
          class="inline-flex items-center gap-2 text-sm text-gray-300 hover:text-white">
          <icon name="Undo2" />
          <span>{{ $t('links.back_to_profile') }}</span>
        </nuxt-link>
      </div>

      <form @submit.prevent="submitForm" class="space-y-6">
        <!-- Avatar -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:gap-4">
          <label class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.avatar') }}</label>
          <div class="flex-1">
            <avatar-edit :value="user.avatar" />
          </div>
        </div>

        <!-- Nickname -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:gap-4">
          <label class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.nickname') }}</label>
          <div class="flex-1">
            <input
              v-model="form.nickname"
              class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400"
              type="text"
              autocomplete="off"
              :placeholder="$t('form.placeholder.enter_nickname')" />
          </div>
        </div>

        <!-- Bio -->
        <div class="flex flex-col sm:flex-row sm:items-start sm:gap-4">
          <label class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.bio') }}</label>
          <div class="flex-1">
            <textarea
              v-model="form.description"
              class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400"
              rows="2"
              :placeholder="$t('form.placeholder.write_bio')"></textarea>
          </div>
        </div>

        <!-- Website -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:gap-4">
          <label class="w-full sm:w-40 text-sm text-gray-300">{{ $t('form.label.website') }}</label>
          <div class="flex-1">
            <input
              v-model="form.homePage"
              class="w-full px-3 py-2 border rounded-md bg-gray-700 text-white placeholder-gray-400"
              type="text"
              autocomplete="off"
              :placeholder="$t('form.placeholder.enter_website')" />
          </div>
        </div>

        <!-- Save -->
        <div class="flex justify-end">
          <button type="submit"
            class="inline-flex items-center px-4 py-2 bg-green-600 hover:bg-green-500 text-white rounded-md">{{
              $t('form.button.save') }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
const i18n = useI18n()
definePageMeta({
  middleware: ['auth'],
  layout: 'profile',
})

useHead({
  title: useSiteTitle(i18n.t('page.personal_info')),
})

const userStore = useUserStore()
const user = computed(() => {
  return userStore.user
})

const form = ref({
  nickname: '',
  avatar: '',
  homePage: '',
  description: '',
})

if (user.value != null) {
  form.value.nickname = user.value.nickname
  form.value.avatar = user.value.avatar
  form.value.homePage = user.value.homePage
  form.value.description = user.value.description
}

async function submitForm() {
  try {
    await useHttpPostForm(`/api/users/edit/${user.value.id}`, {
      body: form.value,
    })
    await reload()
    useMsgSuccess(i18n.t('message.profile_update_success'))
  } catch (e) {
    console.error(e)
    useMsgError(i18n.t('message.profile_update_failure', { error: (e.message || e) }))
  }
}
async function reload() {
  user.value = await useHttpGet('/api/users/me')
  form.value = { ...user.value }
}
</script>

<style scoped>
/* Tailwind handles visual styling; no additional CSS required. */
</style>
