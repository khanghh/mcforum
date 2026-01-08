<template>
  <nav
    class="bg-gradient-to-r from-purple-900/30 to-indigo-900/30 backdrop-blur-lg border-b border-purple-500/20 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <img src="/images/logo.png" class="w-12 h-12" alt="MineViet Logo" />
        <nuxt-link href="/">
          <div
            class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 server-title ">
            MineViet Forum
          </div>
          <div class="hidden md:block text-xs text-purple-300">Minecraft Server Vietnam</div>
        </nuxt-link>
      </div>
      <ul class="flex items-center space-x-4 sm:space-x-6 text-gray-300 font-medium">
        <li v-for="(nav, idx) in siteNavs" :key="idx">
          <a :href="nav.url"
            class="text-gray-300 hover:text-purple-400 transition-colors"
            :target="nav.url && nav.url.startsWith('http') ? '_blank' : null">
            <span class="hidden sm:inline">{{ nav.title }}</span>
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
          </nuxt-link>
        </li>

        <li class="flex-shrink-0">
          <AvatarMenu v-if="user" class="flex-shrink-0" :user="user" />
          <a v-else :href="loginUrl || '#'"
            class="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded-lg transition-colors font-semibold text-sm">
            {{ $t('page.signin') }}
          </a>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
import { useUserStore } from '~/stores/user'
import { useConfigStore } from '~/stores/config'
import AvatarMenu from '~/components/navbar/AvatarMenu.vue'

const api = useApi()
const route = useRoute()
const configStore = useConfigStore()
const userStore = useUserStore()
const runtimeConfig = useRuntimeConfig()

const { user } = storeToRefs(userStore)
const { loginUrl } = runtimeConfig.public

const siteNavs = computed(() => configStore.config?.siteNavs || [])

const { data: recentMsgResp } = await useAsyncData(
  'recent-messages',
  () => api.getRecentMessages().catch(() => null),
  {
    default: () => ({
      topics: 0,
      posts: 0,
      members: 0,
      visits: 0,
      newestMember: '',
    }),
    watch: [() => route.fullPath],
  }
)

const notifCount = computed(() => {
  return recentMsgResp.value?.count || 0
})

// Mock notification count for UI preview
</script>
