<template>
  <div ref="dropdownRef" class="relative inline-block">
    <button class="flex items-center space-x-2 focus:outline-none" @click="toggleDropdown">
      <Avatar :src="avatarSrc" :username="user?.username"
        class="w-10 h-10 object-cover border border-purple-600 rounded-lg" />
    </button>

    <!-- Dropdown Menu (light style) -->
    <div v-if="isDropdownOpen"
      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-sm py-1 ring-1 ring-black ring-opacity-5 focus:outline-none z-50 border border-gray-200">
      <div class="px-4 py-2 flex items-center gap-3">
        <Avatar :src="avatarSrc" :username="user?.username"
          class="w-10 h-10 object-cover border border-purple-600 rounded-lg" />
        <div class="min-w-0">
          <div class="text-sm font-semibold text-gray-800 truncate">{{ user?.username }}</div>
          <div v-if="user?.email" class="text-xs text-gray-500 truncate">{{ user?.email }}</div>
        </div>
      </div>
      <div class="border-t border-gray-100 my-1"></div>
      <a :href="`/users/${user?.username}`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidUserLarge" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.profile') }}
      </a>
      <a :href="`/users/${user?.username}/favorites`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidStar" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.favorites') }}
      </a>
      <a href="/users/me/profile"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <Icon name="Fa7SolidGear" class="mr-2 w-4 text-gray-500" /> {{ $t('navbar.account_settings') }}
      </a>
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
import { ref, computed } from 'vue'
import type { UserProfile } from '~/types'
import { onClickOutside } from '@vueuse/core'
import { useUserStore } from '~/stores/user'

const props = defineProps<{ user: UserProfile }>()


const userStore = useUserStore()
const isDropdownOpen = ref(false)
const dropdownRef = ref(null)

const avatarSrc = computed(() => userStore.user?.avatar || '/images/default-avatar.png')

onClickOutside(dropdownRef, () => {
  isDropdownOpen.value = false
})

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
}

const logout = async () => {
  await userStore.signout()
  isDropdownOpen.value = false
  // Optional: Redirect to home or refresh
  window.location.href = '/'
}
</script>
