<template>
  <div class="border-b border-purple-500/20">
    <div class="flex overflow-x-auto custom-scrollbar" role="tablist">
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="activeTab === 'topics' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="changeTab('topics')">
        <Icon name="Fa7SolidComments" class="mr-2" /> POSTS
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="activeTab === 'activity' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="changeTab('activity')">
        <Icon name="RiFireLine" class="mr-2" /> ACTIVITIES
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="activeTab === 'followers' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="changeTab('followers')">
        <Icon name="TablerStar" class="mr-2" /> FOLLOWERS
      </button>
      <button
        class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
        :class="activeTab === 'following' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
        @click="changeTab('following')">
        <Icon name="PhUsersThreeFill" class="mr-2" /> FOLLOWING
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'

const props = defineProps({
  activeTab: {
    type: String as () => string | undefined,
  },
})

const route = useRoute()
const router = useRouter()

function changeTab(tab: string) {
  router.push({
    name: route.name as string, // keep the same named route
    params: {
      ...route.params, // keep existing params like id
      tab: tab, // change only the tab
    },
  })
}
</script>

<style scoped>
/* No local styles; component uses global utility classes */
</style>
