<template>
  <nav
    class="bg-gradient-to-r from-purple-900/30 to-indigo-900/30 backdrop-blur-lg border-b border-purple-500/20 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <img src="/images/logo.png" class="w-12 h-12" />
        <nuxt-link href="/">
          <div
            class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 server-title ">
            MineViet Forum
          </div>
          <div class="hidden md:block text-xs text-purple-300">Minecraft Server Vietnam</div>
        </nuxt-link>
      </div>
      <ul class="flex items-center space-x-4 sm:space-x-6 text-gray-300 font-medium">
        <li>
          <a href="https://play.mineviet.com"
            class="text-gray-300 hover:text-purple-400 transition-colors"
            target="_blank">
            <Icon name="TablerServer" class="mr-1" size="25" />
            <span class="hidden sm:inline">Máy chủ</span>
          </a>
        </li>
        <li>
          <a href="https://skin.mineviet.com"
            class="text-gray-300 hover:text-purple-400 transition-colors">
            <Icon name="TablerPaletteFilled" class="mr-1" size="25" />
            <span class="hidden sm:inline">Đổi skin</span>
          </a>
        </li>
        <li>
          <a href="https://ban.mineviet.com"
            class="text-gray-300 hover:text-purple-400 transition-colors">
            <Icon name="TablerBan" class="mr-1" size="25" />
            <span class="hidden sm:inline">Vi phạm</span>
          </a>
        </li>
        <li v-if="user" class="flex-shrink-0">
          <nuxt-link href="/users/me/messages"
            class="text-gray-300 hover:text-purple-400 transition-colors flex items-center">
            <span class="relative mr-2 inline-flex items-center">
              <Icon name="Fa7SolidBell" size="25" />
              <span v-if="notifCount && notifCount > 0"
                class="absolute -top-1 -right-2 inline-flex items-center justify-center px-1.5 py-0.5 text-[10px] font-semibold leading-none text-white bg-red-500 rounded-full shadow">
                {{ notifCount }}
              </span>
            </span>
            <span class="hidden sm:inline">{{ $t('message.notifications') }}</span>
          </nuxt-link>
        </li>

        <li class="flex-shrink-0">
          <AvatarMenu v-if="user" class="flex-shrink-0" :user="user" />
          <a v-else href="/signin"
            class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg transition-colors font-semibold text-sm">
            {{ $t('page.signin') }}
          </a>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '~/stores/user'
import AvatarMenu from '~/components/navbar/AvatarMenu.vue'
const api = useApi()

const userStore = useUserStore()
const user = userStore.user

const notifCount = ref(0)

const resp = await api.getRecentMessages().catch(() => null)
notifCount.value = resp?.count || 0

// Mock notification count for UI preview
</script>
