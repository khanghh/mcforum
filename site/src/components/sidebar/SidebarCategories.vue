<template>
  <div class="gaming-card rounded-xl p-4">
    <h3 class="text-lg font-bold text-purple-300 mb-4 flex items-center gaming-title">
      <Icon name="IcRoundViewList" class="mr-2" /> CATEGORIES
    </h3>
    <nav class="space-y-1">
      <nuxt-link v-for="item in itemMenus" :key="item.urlPath" :to="item.urlPath"
        class="category-item block px-3 py-2 rounded text-md font-semibold transition-all hover:bg-purple-500/10 gaming-title"
        :class="{ active: item.urlPath === menuPath }">
        <Icon :name="getIcon(item)" class="mr-2" :class="getIconColor(item)" /> {{ item.name }}
      </nuxt-link>
    </nav>
  </div>
</template>

<script setup>
const route = useRoute()
const configStore = useConfigStore()
const menuPath = route.path === '/' ? '/whats-new' : route.path

const itemMenus = computed(() => configStore.config.menuItems || [])

function getIcon(forum) {
  const name = forum.name.toLowerCase()
  if (name.includes('server')) return 'server'
  if (name.includes('build') || name.includes('creative')) return 'TablerCube'
  if (name.includes('mod') || name.includes('plugin')) return 'TablerPuzzle'
  if (name.includes('help') || name.includes('support')) return 'TablerQuestionCircle'
  if (name.includes('market') || name.includes('shop')) return 'TablerShoppingCart'
  if (name.includes('all') || name.includes('home')) return 'TablerHomeFilled'
  return 'comments'
}

function getIconColor(forum) {
  const name = forum.name.toLowerCase()
  if (name.includes('server')) return 'text-blue-400'
  if (name.includes('build') || name.includes('creative')) return 'text-green-400'
  if (name.includes('mod') || name.includes('plugin')) return 'text-pink-400'
  if (name.includes('help') || name.includes('support')) return 'text-yellow-400'
  if (name.includes('market') || name.includes('shop')) return 'text-emerald-400'
  return 'text-purple-400'
}
</script>
