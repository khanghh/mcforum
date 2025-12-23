<template>
  <div class="lg:col-span-1">
    <div
      class="bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] border border-[rgba(139,92,246,0.2)] rounded-2xl p-6 lg:sticky lg:top-24 space-y-6">
      <div v-if="user.bio">
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <FontAwesome :icon="['fas', 'user']" class="mr-2" /> Bio
        </h3>
        <p class="mt-3 text-gray-300 text-sm leading-relaxed">
          {{ user.bio }}
        </p>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <FontAwesome :icon="['fas', 'info-circle']" class="mr-2" /> Basic Info
        </h3>
        <div class="mt-3 space-y-2">
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <FontAwesome :icon="['fas', 'calendar-plus']" class="text-purple-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Joined</span>
            </div>
            <span class="text-sm font-bold gaming-title">{{ formatDate(user.createTime) }}</span>
          </div>
          <!-- Mock Data for design completeness -->
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <FontAwesome :icon="['fas', 'clock']" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Play Time</span>
            </div>
            <span class="text-sm font-bold gaming-title">2,847 hrs</span>
          </div>
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <FontAwesome :icon="['fas', 'envelope']" class="text-green-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Email</span>
            </div>
            <span class="text-sm font-bold gaming-title truncate max-w-[150px]">{{ user.email || 'Hidden' }}</span>
          </div>
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <FontAwesome :icon="['fas', 'map-marker-alt']" class="text-cyan-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Location</span>
            </div>
            <span class="text-sm font-bold gaming-title">{{ user.location || 'Unknown' }}</span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <FontAwesome :icon="['fas', 'chart-line']" class="mr-2" /> Level Progress
        </h3>
        <div
          class="mt-3 p-4 rounded-xl bg-gradient-to-br from-purple-500/10 to-pink-500/10 border border-purple-500/30">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="text-2xl font-bold text-purple-300 gaming-title">{{ calculateLevel(user.score) }}</span>
              <FontAwesome :icon="['fas', 'arrow-right']" class="text-purple-400 text-xs" />
              <span class="text-xl font-bold text-gray-400 gaming-title">{{ calculateLevel(user.score) + 1 }}</span>
            </div>
            <span class="text-xs font-bold text-purple-300">75% to next level</span>
          </div>
          <div class="relative w-full h-3 bg-gray-700/50 rounded-full overflow-hidden">
            <div
              class="absolute inset-0 bg-gradient-to-r from-purple-600 via-pink-500 to-purple-600 rounded-full"
              style="width: 75%;">
            </div>
          </div>
          <div class="mt-2 flex items-center justify-between text-xs text-gray-400">
            <span>{{ user.score }} / 5,000 XP</span>
            <span class="flex items-center gap-1">
              <FontAwesome :icon="['fas', 'fire']" class="text-orange-400" />
              5 day streak
            </span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <FontAwesome :icon="['fas', 'gamepad']" class="mr-2" /> Main Games
        </h3>
        <div class="flex flex-wrap gap-2 mt-3">
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-blue-500/20 text-blue-300 border border-blue-500/30">
            <FontAwesome :icon="['fas', 'crosshairs']" class="mr-1.5" /> SkyBlock
          </span>
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-purple-500/20 text-purple-300 border border-purple-500/30">
            <FontAwesome :icon="['fas', 'chess']" class="mr-1.5" /> Survival
          </span>
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-orange-500/20 text-orange-300 border border-orange-500/30">
            <FontAwesome :icon="['fas', 'bomb']" class="mr-1.5" /> Bedwars
          </span>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <FontAwesome :icon="['fas', 'trophy']" class="mr-2" /> Achievements
        </h3>
        <div class="mt-3 grid grid-cols-3 gap-2">
          <!-- Static Achievements List -->
          <div class="flex flex-col items-center gap-2 p-3 rounded-lg bg-amber-500/10 border border-amber-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center flex-shrink-0">
              <FontAwesome :icon="['fas', 'crown']" class="text-white text-sm" />
            </div>
            <div class="text-center">
              <p class="text-[10px] font-bold text-amber-300 gaming-title">CHAMPION</p>
            </div>
          </div>
          <div
            class="flex flex-col items-center gap-2 p-3 rounded-lg bg-purple-500/10 border border-purple-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center flex-shrink-0">
              <FontAwesome :icon="['fas', 'star']" class="text-white text-sm" />
            </div>
            <div class="text-center">
              <p class="text-[10px] font-bold text-purple-300 gaming-title">LEGEND</p>
            </div>
          </div>
          <div class="flex flex-col items-center gap-2 p-3 rounded-lg bg-blue-500/10 border border-blue-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center flex-shrink-0">
              <FontAwesome :icon="['fas', 'fire']" class="text-white text-sm" />
            </div>
            <div class="text-center">
              <p class="text-[10px] font-bold text-blue-300 gaming-title">STREAK</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

function formatDate(timestamp) {
  if (!timestamp) return 'Unknown'
  return new Date(timestamp).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
}

function calculateLevel(score) {
  return Math.floor(Math.sqrt(score || 0)) + 1
}
</script>
