<template>
  <nav class="sticky top-20 w-full gaming-card rounded-xl p-4 border border-white/5 bg-gray-800/50 backdrop-blur-sm">
    <ul class="space-y-2">
      <li v-for="item in menuItems" :key="item.urlPath">
        <nuxt-link
          :to="forumUrl(item.urlPath)"
          class="flex items-center px-3 py-2.5 rounded-lg text-sm font-bold transition-all duration-200 group"
          :class="[
            isActive(item.urlPath)
              ? 'bg-purple-600 text-white shadow-lg shadow-purple-500/20'
              : 'text-gray-400 hover:bg-white/5 hover:text-gray-200'
          ]">
          <img
            :src="forumLogo(item)"
            class="w-6 h-6 rounded object-cover mr-3 transition-transform group-hover:scale-110"
            alt="">
          <span>{{ item.name }}</span>
        </nuxt-link>
      </li>
    </ul>
  </nav>
</template>

<script setup>
const route = useRoute()
const configStore = useConfigStore()

const menuItems = configStore.config.menuItems

function forumLogo(forum) {
  return forum.logo || '/images/node.png'
}

function isActive(urlPath) {
  return route.path === urlPath
}
</script>
