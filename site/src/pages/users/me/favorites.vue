<template>
  <div v-if="user">
    <ProfileHeader v-model:user="user" />

    <div class="flex flex-col lg:flex-row gap-6 mt-6">
      <div class="w-full lg:w-80 lg:flex-shrink-0">
        <ProfileSidebar :user="user" />
      </div>

      <div class="flex-1 min-w-0">
        <!-- Favorites Header -->
        <div class="gaming-card rounded-2xl p-6 mb-6 relative overflow-hidden">
          <div class="relative flex items-center justify-between">
            <div class="flex items-center gap-4">
              <div
                class="w-14 h-14 rounded-xl bg-gradient-to-br from-purple-500/20 to-pink-500/20 flex items-center justify-center ring-2 ring-purple-500/30 shadow-lg shadow-purple-500/20">
                <Icon name="Fa7SolidStar" class="text-purple-400 animate-pulse" size="25" />
              </div>
              <div>
                <h2
                  class="gaming-title text-3xl font-bold bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 bg-clip-text text-transparent">
                  {{ $t('page.favorites') }}
                </h2>
                <p class="text-sm text-gray-400 mt-1">
                  {{ $t('page.favorites_desc') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Favorites List -->
        <div class="space-y-4">
          <LoadMoreAsync :cursor="favoritesCursor">
            <template #default="{ items }">
              <div v-if="items && items.length" class="space-y-4">
                <div v-for="item in items" :key="item.id"
                  class="relative hover:bg-white/5 rounded-lg p-5 transition-all duration-200 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] shadow-[0_0_0_1px_rgb(139,92,246,0.2)] group">
                  <div class="flex-1 min-w-0">
                    <!-- Title -->
                    <h2 class="font-bold mb-3 flex items-center justify-between gap-4">
                      <a :href="item.url"
                        class="text-white text-lg md:text-xl hover:text-purple-400 transition-colors gaming-title truncate inline-block">
                        {{ item.title }}
                      </a>

                      <!-- Entity Type Badge -->
                      <div
                        class="hidden sm:flex items-center gap-1.5 px-2 py-0.5 rounded-md bg-white/5 border border-white/10 text-xs text-gray-400 flex-shrink-0">
                        <Icon name="Fa7SolidBookmark" size="10" />
                        <span class="capitalize">{{ item.entityType }}</span>
                      </div>
                    </h2>

                    <!-- User Info & Meta -->
                    <div class="flex items-center gap-3 mb-3">
                      <nuxt-link :to="'/users/' + item.user.username" class="flex-shrink-0">
                        <Avatar :src="item.user.avatar" :username="item.user.username" size="32"
                          class="rounded border border-purple-300 flex-shrink-0" />
                      </nuxt-link>
                      <div class="flex flex-col min-w-0">
                        <nuxt-link :to="'/users/' + item.user.username"
                          class="font-bold text-gray-300 gaming-title text-sm hover:text-purple-400 transition-colors truncate">
                          {{ item.user.username }}
                        </nuxt-link>
                        <span class="text-xs text-gray-500 mt-0.5">{{ usePrettyDate(item.createTime) }}</span>
                      </div>
                    </div>

                    <!-- Snippet -->
                    <p v-if="item.content" class="text-gray-300 text-sm mb-0 line-clamp-2 leading-relaxed">
                      {{ item.content }}
                    </p>
                  </div>
                </div>
              </div>
            </template>

            <template #empty>
              <div class="gaming-card rounded-2xl p-12 text-center border dashed border-white/10 bg-transparent">
                <div
                  class="w-20 h-20 mx-auto mb-4 rounded-full bg-gradient-to-br from-purple-500/10 to-pink-500/10 flex items-center justify-center ring-1 ring-white/10">
                  <Icon name="Fa7SolidStar" class="text-gray-600" size="32" />
                </div>
                <h3 class="text-xl font-bold text-gray-300 mb-2">
                  No favorites yet
                </h3>
                <p class="text-gray-500 text-sm max-w-xs mx-auto">
                  Mark topics as favorite to quickly access them here later.
                </p>
              </div>
            </template>
          </LoadMoreAsync>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import ProfileHeader from '~/components/profile/ProfileHeader.vue'
import ProfileSidebar from '~/components/profile/ProfileSidebar.vue'
import LoadMoreAsync from '~/components/LoadMoreAsync.vue'

definePageMeta({
  middleware: ['auth'],
  layout: 'default',
})

const i18n = useI18n()
const api = useApi()
const userStore = useUserStore()

const { user } = storeToRefs(userStore)

const favoritesCursor = api.getMyFavorites()

useHead({
  title: useSiteTitle(i18n.t('page.favorites', { nickname: user.value?.username })),
  bodyAttrs: {
    class: 'bg-gaming-bg',
  },
})
</script>
