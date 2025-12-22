<template>
  <div class="relative inline-block" ref="dropdownRef">
    <button class="flex items-center space-x-2 focus:outline-none" @click="toggleDropdown">
      <Avatar :src="avatarSrc" :username="username" border class="w-10 h-10 rounded-full object-cover" />
    </button>

    <!-- Dropdown Menu -->
    <div v-if="isDropdownOpen"
      class="absolute right-0 mt-2 w-48 bg-gray-900 rounded-md shadow-lg py-1 ring-1 ring-black ring-opacity-5 focus:outline-none z-50 border border-purple-500/30 backdrop-blur-md">
      <div class="px-4 py-2 text-sm font-bold text-white truncate">
        {{ userStore.user?.nickname }}
      </div>
      <div class="border-t border-gray-700 my-1"></div>
      <a :href="`/user/${userStore.user?.id}`"
        class="block px-4 py-2 text-sm text-gray-300 hover:bg-purple-800/50 hover:text-white transition-colors">
        <FontAwesome :icon="['fas', 'user']" class="mr-2 w-4" /> Profile
      </a>
      <a href="/user/favorites"
        class="block px-4 py-2 text-sm text-gray-300 hover:bg-purple-800/50 hover:text-white transition-colors">
        <FontAwesome :icon="['fas', 'heart']" class="mr-2 w-4" /> Favorite
      </a>
      <a href="/user/profile"
        class="block px-4 py-2 text-sm text-gray-300 hover:bg-purple-800/50 hover:text-white transition-colors">
        <FontAwesome :icon="['fas', 'cog']" class="mr-2 w-4" /> Account Settings
      </a>
      <div class="border-t border-gray-700 my-1"></div>
      <button
        class="block w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-purple-800/50 hover:text-red-300 transition-colors"
        @click="logout">
        <FontAwesome :icon="['fas', 'sign-out-alt']" class="mr-2 w-4" /> Logout
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useUserStore } from '~/stores/user'

const userStore = useUserStore()
const isDropdownOpen = ref(false)
const dropdownRef = ref(null)
const avatarSrc = computed(() => {
  return userStore.user?.avatar || '/images/default-avatar.png'
})
console.log('User Store in AvatarMenu:', avatarSrc.value)
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