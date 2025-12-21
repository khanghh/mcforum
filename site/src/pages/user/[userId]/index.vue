<template>
  <div class="min-h-screen bg-gray-900 text-gray-100 font-sans selection:bg-purple-500 selection:text-white"
    style="background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);">
    <GamingNavbar />

    <main class="flex-grow py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6">
        <ProfileHeader :user="user" />

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mt-6">
          <ProfileSidebar :user="user" />

          <div class="lg:col-span-2">
            <div class="gaming-card rounded-2xl overflow-hidden min-h-[500px]">
              <!-- Tabs Navigation -->
              <div class="border-b border-purple-500/20">
                <div class="flex overflow-x-auto custom-scrollbar" role="tablist">
                  <button
                    class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
                    :class="activeTab === 'topics' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
                    @click="activeTab = 'topics'">
                    <FontAwesome :icon="['fas', 'fire']" class="mr-2" /> POSTS
                  </button>
                  <button
                    class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
                    :class="activeTab === 'comments' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
                    @click="activeTab = 'comments'">
                    <FontAwesome :icon="['fas', 'comments']" class="mr-2" /> ACTIVITY
                  </button>
                  <button
                    class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
                    :class="activeTab === 'following' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
                    @click="activeTab = 'following'">
                    <FontAwesome :icon="['fas', 'users']" class="mr-2" /> FOLLOWING
                  </button>
                  <button
                    class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
                    :class="activeTab === 'fans' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
                    @click="activeTab = 'fans'">
                    <FontAwesome :icon="['fas', 'star']" class="mr-2" /> FANS
                  </button>
                  <button
                    class="tab-btn flex-1 flex items-center justify-center px-6 py-4 border-b-2 font-bold whitespace-nowrap gaming-title transition-colors"
                    :class="activeTab === 'favorites' ? 'border-purple-500 text-purple-400' : 'border-transparent text-gray-400 hover:text-gray-200'"
                    @click="activeTab = 'favorites'">
                    <FontAwesome :icon="['fas', 'heart']" class="mr-2" /> FAVORITES
                  </button>
                </div>
              </div>

              <!-- Tab Content -->
              <div class="p-6">
                <!-- Topics Content -->
                <div v-show="activeTab === 'topics'" class="space-y-4">
                  <LoadMoreAsync v-slot="{ items }" url="/api/topic/user/topics" :params="{ userId: user.id }">
                    <GamingTopicList :topics="items" :show-pinned="false" />
                  </LoadMoreAsync>
                </div>

                <!-- Activity (Comments) Content -->
                <div v-show="activeTab === 'comments'" class="space-y-4">
                  <LoadMoreAsync v-slot="{ items }" url="/api/comment/list"
                    :params="{ userId: user.id, entityType: 'topic' }">
                    <!-- Placeholder for comments list, using a simplified card for now or custom component -->
                    <div v-for="item in items" :key="item.id"
                      class="gaming-card rounded-xl p-4 border-l-4 border-blue-500/50">
                      <div class="flex items-start gap-4">
                        <FontAwesome :icon="['fas', 'comment-dots']" class="text-2xl text-blue-400 mt-1" />
                        <div class="flex-1">
                          <div class="text-gray-400 text-sm mb-1">Commented on <nuxt-link
                              :to="`/topics/${item.entityId}`" class="text-blue-300 hover:underline">Topic #{{
                                item.entityId }}</nuxt-link></div>
                          <p class="text-gray-300 leading-relaxed">{{ item.content }}</p>
                          <div class="mt-2 text-xs text-gray-500">{{ usePrettyDate(item.createTime) }}</div>
                        </div>
                      </div>
                    </div>
                  </LoadMoreAsync>
                </div>

                <!-- Following Content -->
                <div v-show="activeTab === 'following'" class="space-y-4">
                  <!-- Placeholder content for now -->
                  <div class="text-gray-400 text-center py-10">Following list implementation coming soon.</div>
                </div>

                <!-- Fans Content -->
                <div v-show="activeTab === 'fans'" class="space-y-4">
                  <!-- Placeholder content for now -->
                  <div class="text-gray-400 text-center py-10">Fans list implementation coming soon.</div>
                </div>

                <!-- Favorites Content -->
                <div v-show="activeTab === 'favorites'" class="space-y-4">
                  <LoadMoreAsync v-slot="{ items }" url="/api/topic/user/favorites" :params="{ userId: user.id }">
                    <GamingTopicList :topics="items" :show-pinned="false" />
                  </LoadMoreAsync>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </main>

    <GamingFooter />
  </div>
</template>

<script setup>
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'

definePageMeta({
  layout: false,
})

const i18n = useI18n()
const route = useRoute()
const user = await useHttpGet(`/api/user/${route.params.userId}`)
const activeTab = ref('topics')

useHead({
  title: useSiteTitle(i18n.t('page.profile', { nickname: user.nickname })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>

<style scoped>
/* Scrollbar styling */
.custom-scrollbar::-webkit-scrollbar {
  height: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #4b5563;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #6b7280;
}
</style>
