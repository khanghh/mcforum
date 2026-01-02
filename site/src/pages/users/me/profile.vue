<template>
  <div class="max-w-3xl mx-auto">
    <!-- Header with back button -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center gap-3">
        <div
          class="w-10 h-10 rounded-lg bg-gradient-to-r from-purple-600/20 to-pink-600/20 flex items-center justify-center">
          <Icon name="Fa7SolidUserEdit" class="text-purple-300" />
        </div>
        <h1
          class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 gaming-title">
          {{ $t('page.edit_profile') }}
        </h1>
      </div>

      <nuxt-link
        :to="`/users/${user.username}`"
        class="flex items-center text-sm text-gray-300 hover:text-white transition-colors px-4 py-2 rounded-lg bg-gray-800/50 hover:bg-gray-800/80">
        <Icon name="Fa7SolidArrowLeft" class="mr-2" />
        <span>{{ $t('links.back_to_profile') }}</span>
      </nuxt-link>
    </div>

    <!-- Edit Form Container (gradient border like profile.html, Tailwind-only) -->
    <div
      class="gradient-border rounded-2xl p-6 sm:p-8 relative overflow-hidden backdrop-blur-sm shadow-xl shadow-purple-900/10">
      <div class="absolute inset-0 bg-gradient-to-br from-gray-800/20 via-gray-900/30 to-indigo-900/20"
        data-v-21511f91="">
      </div>
      <form class="relative z-10 space-y-8" @submit.prevent="submitForm">
        <!-- Avatar -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-center">
            <Icon name="Fa7SolidImage" class="mr-2" />
            {{ $t('form.label.avatar') }}
          </label>
          <div>
            <AvatarDecorated :user="user">
              <AvatarEdit v-model="avatarUrl" :username="user.username" class="w-full h-full object-cover"
                size="120" />
            </AvatarDecorated>
            <div class="mt-2">
              <button type="button" class="text-sm text-red-400 hover:text-red-300 flex items-center gap-1"
                @click="removeAvatar">
                <Icon name="Fa7SolidTrashCan" />
                {{ $t('profile.remove_avatar') }}
              </button>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-600/50"></div>

        <!-- Nickname -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-center">
            <Icon name="Fa7SolidUserTag" class="mr-2" />
            {{ $t('form.label.nickname') }}
          </label>
          <div class="flex-1">
            <div class="relative">
              <input
                v-model="form.nickname"
                class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-purple-500 focus:ring-2 transition-colors transition-shadow"
                type="text"
                autocomplete="off"
                :placeholder="$t('form.placeholder.enter_nickname')" />
              <div class="absolute right-3 top-3">
                <Icon name="Fa7SolidPen" class="text-gray-500" />
              </div>
            </div>
            <p class="text-xs text-gray-400 mt-2">{{ $t('form.help.nickname_public') }}</p>
          </div>
        </div>

        <!-- Username (read-only) -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-center">
            <Icon name="Fa7SolidAt" class="mr-2" />
            {{ $t('form.label.username') }}
          </label>
          <div class="flex-1">
            <div
              class="w-full px-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-gray-300 flex items-center justify-between ">
              <div class="flex items-center min-w-0">
                <span class="ml-1 font-medium truncate text-gray-500">{{ user.username }}</span>
              </div>
              <span class="text-xs px-2 py-1 bg-gray-700 rounded text-gray-400">
                {{ $t('form.help.cannot_change') }}
              </span>
            </div>
            <p class="text-xs text-gray-400 mt-2">{{ $t('form.help.username_immutable') }}</p>
          </div>
        </div>

        <!-- Email (read-only) -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-center">
            <Icon name="TablerMail" class="mr-2" />
            {{ $t('form.label.email') }}
          </label>
          <div class="flex-1">
            <div
              class="w-full px-4 py-3 bg-gray-800/50 border border-gray-700 rounded-lg text-gray-400 flex items-center justify-between ">
              <span class="truncate text-gray-500">{{ user.email }}</span>
              <a
                href="https://auth.mineviet.com/profile"
                target="_blank"
                rel="noopener"
                class="ml-4 text-xs text-purple-400 hover:text-purple-300 transition-colors flex items-center gap-1">
                <Icon name="Fa7SolidExternalLinkAlt" />
                <span>{{ $t('form.help.change_in_account') }}</span>
              </a>
            </div>
            <div class="flex justify-between items-center mt-2 gap-3">
              <p class="text-xs text-gray-400">{{ $t('form.help.email_used_for_notifications') }}</p>
            </div>
          </div>
        </div>

        <!-- Bio -->
        <div class="flex flex-col sm:flex-row sm:items-start gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-start pt-3">
            <Icon name="Fa7SolidCircleInfo" class="mr-2" />
            {{ $t('form.label.bio') }}
          </label>
          <div class="flex-1">
            <div class="relative">
              <textarea
                v-model="form.bio"
                :maxlength="bioMax"
                class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:ring-purple-500 focus:ring-2 transition-colors transition-shadow"
                rows="4"
                :placeholder="$t('form.placeholder.write_bio')"
                @input="handleBioInput"></textarea>
              <div class="absolute right-3 top-3">
                <Icon name="Fa7SolidAlignLeft" class="text-gray-500" />
              </div>
            </div>
            <div class="flex justify-between items-center mt-2">
              <p class="text-xs text-gray-400">{{ $t('form.help.bio_limit') }}</p>
              <span class="text-xs" :class="bioCounterClass">{{ bioCount }}/{{ bioMax }}</span>
            </div>
          </div>
        </div>

        <!-- Location -->
        <div class="flex flex-col sm:flex-row sm:items-center gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-center">
            <Icon name="Fa7SolidGlobeAmericas" class="mr-2" />
            {{ $t('form.label.location') }}
          </label>
          <div class="flex-1">
            <div class="relative">
              <input
                v-model="form.location"
                class="w-full px-4 py-3 bg-gray-800/70 border border-gray-700 rounded-lg text-white placeholder-gray-400 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/30 transition-colors transition-shadow"
                type="text"
                autocomplete="off"
                :placeholder="$t('form.placeholder.enter_location')" />
              <div class="absolute right-3 top-3">
                <Icon name="Fa7SolidGlobeAmericas" class="text-gray-500" />
              </div>
            </div>
            <p class="text-xs text-gray-400 mt-2">{{ $t('form.help.location_optional') }}</p>
          </div>
        </div>

        <!-- Other settings -->
        <div class="flex flex-col sm:flex-row sm:items-start gap-3">
          <label class="w-full sm:w-40 text-sm font-medium text-gray-300 flex items-start pt-3">
            <Icon name="Fa7SolidGear" class="mr-2" />
            {{ $t('form.label.other_settings') }}
          </label>
          <div class="flex-1 space-y-4">
            <div
              class="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg border border-gray-700/50 ">
              <div>
                <h4 class="font-medium text-white">{{ $t('profile.settings.locked_profile') }}</h4>
                <p class="text-sm text-gray-400 mt-1">{{ $t('profile.settings.locked_profile_desc') }}</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input v-model="form.lockedProfile" type="checkbox" class="sr-only peer" />
                <div
                  class="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
                </div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg border border-gray-700/50 ">
              <div>
                <h4 class="font-medium text-white">{{ $t('profile.settings.show_location_title') }}</h4>
                <p class="text-sm text-gray-400 mt-1">{{ $t('profile.settings.show_location_desc') }}</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input v-model="form.showLocation" type="checkbox" class="sr-only peer" />
                <div
                  class="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
                </div>
              </label>
            </div>

            <div
              class="flex items-center justify-between p-4 bg-gray-800/50 rounded-lg border border-gray-700/50 ">
              <div>
                <h4 class="font-medium text-white">{{ $t('profile.settings.email_notifications_title') }}</h4>
                <p class="text-sm text-gray-400 mt-1">{{ $t('profile.settings.email_notifications_desc') }}</p>
              </div>
              <label class="relative inline-flex items-center cursor-pointer">
                <input v-model="form.emailNotify" type="checkbox" class="sr-only peer" />
                <div
                  class="w-12 h-6 bg-gray-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-6 after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-gradient-to-r from-purple-600 to-pink-600">
                </div>
              </label>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-600/50"></div>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row justify-end gap-4 pt-4">
          <button
            type="button"
            class="px-6 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg font-medium transition-colors border border-gray-700 gaming-title"
            @click="navigateTo(`/users/${user.username}`)">
            {{ $t('form.button.cancel') }}
          </button>
          <button
            type="submit"
            :disabled="saving"
            class="px-8 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg font-bold flex items-center justify-center shadow-[0_0_10px_rgba(139,92,246,0.5),0_0_20px_rgba(139,92,246,0.3)] gaming-title tracking-[0.5px] transition-transform hover:scale-[1.02] active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed">
            <span v-if="saving" class="flex items-center">
              <svg class="animate-spin mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none"
                viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"></path>
              </svg>
              {{ $t('form.button.save_changes') }}
            </span>
            <span v-else class="flex items-center">
              <Icon name="Fa7SolidSave" class="mr-2" />
              {{ $t('form.button.save_changes') }}
            </span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { storeToRefs } from 'pinia'

const i18n = useI18n()
const api = useApi()
const toast = useToast()
const dialog = useConfirmDialog()
const userStore = useUserStore()

definePageMeta({
  middleware: ['auth'],
  layout: 'default',
})

const { user } = storeToRefs(userStore)
const avatarUrl = ref(user.value.avatar)

const form = ref({
  nickname: user.value.nickname,
  bio: user.value.bio,
  avatar: avatarUrl.value,
  location: user.value.location,
  lockedProfile: user.value.settings.lockedProfile,
  showLocation: user.value.settings.showLocation,
  emailNotify: user.value.settings.emailNotify,
})

const bioMax = 200
const bioCount = computed(() => (form.value.bio || '').length)
const bioCounterClass = computed(() => {
  if (bioCount.value > bioMax - 20) return 'text-red-400'
  if (bioCount.value > bioMax - 50) return 'text-yellow-400'
  return 'text-gray-400'
})
const saving = ref(false)

async function submitForm() {
  saving.value = true
  try {
    await api.updateProfile({
      nickname: form.value.nickname,
      bio: form.value.bio,
      avatar: avatarUrl.value,
      location: form.value.location,
      lockedProfile: form.value.lockedProfile,
      showLocation: form.value.showLocation,
      emailNotify: form.value.emailNotify,
    })
    await reload()
    toast.success(i18n.t('message.profile_update_success'))
  } catch (e) {
    console.error(e?.data?.error.message || e)
    if (e?.data?.error.message) {
      toast.error(e.data.error.message)
      return
    }
    toast.error(i18n.t('message.profile_update_failure'))
  } finally {
    saving.value = false
  }
}

function handleBioInput(e) {
  const val = e?.target?.value ?? ''
  if (val.length > bioMax) {
    form.value.bio = val.slice(0, bioMax)
    e.target.value = form.value.bio
  } else {
    form.value.bio = val
  }
}

async function reload() {
  if (typeof userStore.getCurrent === 'function') {
    await userStore.getCurrent().catch(() => undefined)
  }
  const u = user.value
  console.log('Reloaded user:', u)
  if (!u) return
  avatarUrl.value = u.avatar
  form.value.nickname = u.nickname
  form.value.bio = u.bio
  form.value.location = u.location
  form.value.lockedProfile = u.settings.lockedProfile
  form.value.showLocation = u.settings.showLocation
  form.value.emailNotify = u.settings.emailNotify
}

async function removeAvatar() {
  await dialog.show({
    title: i18n.t('dialog.title.confirm_delete'),
    message: i18n.t('dialog.message.confirm_remove_avatar'),
    confirmText: i18n.t('dialog.button.confirm'),
    cancelText: i18n.t('dialog.button.cancel'),
    icon: 'Fa7SolidTrashCan',
    variant: 'warning',
    onConfirm: async () => {
      try {
        await api.removeAvatar()
        await reload()
      } catch (e) {
        const errMsg = e?.data?.error?.message || e.message || 'Unknown error'
        toast.error(i18n.t('message.action_failure', { erroro: errMsg }))
      }
    },
  })
}

useHead({
  title: useSiteTitle(i18n.t('page.edit_profile')),
})
</script>

<style scoped>
/* Static gradient border without animation */
.gradient-border {
  position: relative;
  border: 2px solid transparent;
  background: linear-gradient(145deg, rgba(45, 45, 70, 0.9), rgba(35, 35, 55, 0.95)) padding-box,
    linear-gradient(45deg, #7c3aed, #8b5cf6, #a78bfa) border-box;
  background-size: 300% 300%;
}
</style>
