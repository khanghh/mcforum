<template>
  <div class="relative inline-block" ref="dropdownRef">
    <button class="flex items-center space-x-2 focus:outline-none" @click="toggleDropdown">
      <Avatar :src="avatarSrc" :username="username"
        class='w-10 h-10 object-cover border border-purple-600 rounded-lg' />
    </button>

    <!-- Dropdown Menu (light style) -->
    <div v-if="isDropdownOpen"
      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-sm py-1 ring-1 ring-black ring-opacity-5 focus:outline-none z-50 border border-gray-200">
      <div class="px-4 py-2 text-sm font-semibold text-gray-800 truncate">
        {{ userStore.user?.nickname }}
      </div>
      <div class="border-t border-gray-100 my-1"></div>
      <a :href="`/users/${userStore.user?.username}`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <FontAwesome :icon="['fas', 'user']" class="mr-2 w-4 text-gray-500" /> Profile
      </a>
      <a :href="`/users/${userStore.user?.username}/favorites`"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <FontAwesome :icon="['fas', 'heart']" class="mr-2 w-4 text-gray-500" /> Favorite
      </a>
      <a href="/users/me/profile"
        class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900 transition-colors">
        <FontAwesome :icon="['fas', 'cog']" class="mr-2 w-4 text-gray-500" /> Account Settings
      </a>
      <div class="border-t border-gray-100 my-1"></div>
      <button
        class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-gray-100 hover:text-red-700 transition-colors"
        @click="logout">
        <FontAwesome :icon="['fas', 'sign-out-alt']" class="mr-2 w-4 text-red-600" /> Logout
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useUserStore } from '~/stores/user'

const userStore = useUserStore()
const isDropdownOpen = ref(false)
const dropdownRef = ref(null)
const avatarSrc = computed(() => userStore.user?.avatar || '/images/default-avatar.png')
const username = computed(() => {
  return userStore.user?.nickname
})

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