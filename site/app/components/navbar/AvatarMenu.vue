<template>
  <div ref="dropdownRef" class="relative inline-block">
    <button class="flex items-center space-x-2 focus:outline-none" @click="toggleDropdown">
      <Avatar :src="avatarSrc" :username="user?.username"
        class="w-10 h-10 rounded border border-purple-300" />
    </button>

    <!-- Mobile Drawer -->
    <Teleport to="body">
      <div v-if="isMobileDrawerVisible" class="fixed inset-0 z-[9999] md:hidden" role="dialog" aria-modal="true">
        <div class="absolute inset-0 bg-black/40" @click="closeDropdown"></div>
        <Transition name="drawer-slide" @after-leave="onDrawerAfterLeave">
          <div v-if="isDropdownOpen"
            class="absolute right-0 top-0 h-full w-80 max-w-[85vw] bg-white border-l border-gray-200 shadow-xl flex flex-col">
            <!-- User Header -->
            <div class="px-4 py-4 flex items-center gap-3 border-b border-gray-100">
              <Avatar :src="avatarSrc" :username="user?.username" class="w-10 h-10 object-cover rounded" />
              <div class="min-w-0">
                <div class="text-sm font-semibold text-gray-800 truncate">{{ user?.username }}</div>
                <div v-if="user?.email" class="text-xs text-gray-500 truncate">{{ user?.email }}</div>
              </div>
            </div>

            <!-- User Menu (first) -->
            <div class="py-1">
              <nuxt-link :to="`/users/${user?.username}`"
                class="block px-4 py-3 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
                @click="closeDropdown">
                <Icon name="Fa7SolidUserLarge" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.profile') }}
              </nuxt-link>
              <nuxt-link :to="`/users/me/favorites`"
                class="block px-4 py-3 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
                @click="closeDropdown">
                <Icon name="Fa7SolidStar" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.favorites') }}
              </nuxt-link>
              <nuxt-link :to="`/users/me/profile`"
                class="block px-4 py-3 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
                @click="closeDropdown">
                <Icon name="Fa7SolidGear" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.account_settings') }}
              </nuxt-link>
            </div>

            <div class="border-t border-gray-100"></div>

            <!-- Categories List -->
            <div class="flex-1 overflow-y-auto py-2">
              <div class="px-4 pt-2 pb-2 text-xs font-semibold text-gray-500 uppercase tracking-wide">
                {{ $t('widgets.categories') }}
              </div>
              <nav class="space-y-1 px-2">
                <nuxt-link v-for="item in itemMenus" :key="item.urlPath" :to="item.urlPath"
                  class="flex items-center gap-2 px-3 py-2 rounded text-sm font-semibold text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors"
                  :class="{ 'bg-gray-100 text-gray-900': item.urlPath === menuPath }" @click="closeDropdown">
                  <Icon :name="getIcon(item)" :class="getIconColor(item)" />
                  <span class="truncate">{{ item.name }}</span>
                </nuxt-link>
              </nav>
            </div>

            <!-- Logout (bottom) -->
            <div class="border-t border-gray-100"></div>
            <div class="p-3">
              <button
                class="w-full text-left px-4 py-3 text-sm text-red-600 hover:bg-gray-100 hover:text-red-700 transition-colors rounded"
                @click="logout">
                <Icon name="TablerLogout" class="mr-2 w-4 text-red-600" /> {{ $t('navbar.logout') }}
              </button>
            </div>
          </div>
        </Transition>
      </div>
    </Teleport>

    <!-- Desktop Dropdown Menu (light style) -->
    <div v-if="isDropdownOpen"
      class="hidden md:block absolute right-0 mt-2 w-48 bg-white rounded-md shadow-sm py-1 ring-1 ring-black ring-opacity-5 focus:outline-none z-50 border border-gray-200">
      <div class="px-4 py-2 flex items-center gap-3">
        <Avatar :src="avatarSrc" :username="user?.username"
          class="w-10 h-10 object-cover rounded" />
        <div class="min-w-0">
          <div class="text-sm font-semibold text-gray-800 truncate">{{ user?.username }}</div>
          <div v-if="user?.email" class="text-xs text-gray-500 truncate">{{ user?.email }}</div>
        </div>
      </div>
      <div class="border-t border-gray-100 my-1"></div>
      <nuxt-link :to="`/users/${user?.username}`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidUserLarge" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.profile') }}
      </nuxt-link>
      <nuxt-link :to="`/users/me/favorites`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidStar" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.favorites') }}
      </nuxt-link>
      <nuxt-link :to="`/users/me/profile`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidGear" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.account_settings') }}
      </nuxt-link>
      <div class="border-t border-gray-100 my-1"></div>
      <button
        class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 hover:text-red-700 transition-colors"
        @click="logout">
        <Icon name="TablerLogout" class="mr-2 w-4 text-red-600" /> {{ $t('navbar.logout') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { UserProfile } from '~/types'
import { onClickOutside, useMediaQuery } from '@vueuse/core'
import { useUserStore } from '~/stores/user'
import { useConfigStore } from '~/stores/config'

const props = defineProps<{ user: UserProfile }>()

const userStore = useUserStore()
const configStore = useConfigStore()
const isDropdownOpen = ref(false)
const isMobileDrawerVisible = ref(false)
const dropdownRef = ref(null)
const route = useRoute()
const isDesktop = useMediaQuery('(min-width: 768px)')

const { user } = storeToRefs(userStore)
const avatarSrc = computed(() => user.value?.avatar)
const menuPath = computed(() => (route.path === '/' ? '/whats-new' : route.path))
const itemMenus = computed(() => configStore.config.menuItems || [])

onClickOutside(dropdownRef, () => {
  if (!isDesktop.value) return
  isDropdownOpen.value = false
})

const toggleDropdown = () => {
  if (isDesktop.value) {
    isDropdownOpen.value = !isDropdownOpen.value
    return
  }

  if (isDropdownOpen.value) {
    closeDropdown()
    return
  }

  // Mount the drawer first, then insert the panel
  isMobileDrawerVisible.value = true
  nextTick(() => {
    isDropdownOpen.value = true
  })
}

const closeDropdown = () => {
  isDropdownOpen.value = false
}

const onDrawerAfterLeave = () => {
  if (isDesktop.value) return
  if (!isDropdownOpen.value) isMobileDrawerVisible.value = false
}

watch(isDropdownOpen, (open) => {
  if (!process.client) return
  document.body.style.overflow = !isDesktop.value && open ? 'hidden' : ''
})

onBeforeUnmount(() => {
  if (!process.client) return
  document.body.style.overflow = ''
})

function getIcon(forum: { name: string }) {
  const name = forum.name.toLowerCase()
  if (name.includes('server')) return 'server'
  if (name.includes('build') || name.includes('creative')) return 'TablerCube'
  if (name.includes('mod') || name.includes('plugin')) return 'TablerPuzzle'
  if (name.includes('help') || name.includes('support')) return 'TablerQuestionCircle'
  if (name.includes('market') || name.includes('shop')) return 'TablerShoppingCart'
  if (name.includes('all') || name.includes('home')) return 'TablerHomeFilled'
  return 'UilCommentsAlt'
}

function getIconColor(forum: { name: string }) {
  const name = forum.name.toLowerCase()
  if (name.includes('server')) return 'text-blue-400'
  if (name.includes('build') || name.includes('creative')) return 'text-green-400'
  if (name.includes('mod') || name.includes('plugin')) return 'text-pink-400'
  if (name.includes('help') || name.includes('support')) return 'text-yellow-400'
  if (name.includes('market') || name.includes('shop')) return 'text-emerald-400'
  return 'text-purple-400'
}

const logout = async () => {
  await userStore.signout()
  isDropdownOpen.value = false
  // Optional: Redirect to home or refresh
  window.location.href = '/'
}
</script>

<style scoped>
.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: transform 220ms ease, opacity 220ms ease;
  will-change: transform, opacity;
}

/* Right -> Left when opening */
.drawer-slide-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

/* Left -> Right when closing */
.drawer-slide-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

@media (prefers-reduced-motion: reduce) {

  .drawer-slide-enter-active,
  .drawer-slide-leave-active {
    transition: none;
  }
}
</style>
