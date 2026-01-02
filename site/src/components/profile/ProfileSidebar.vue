<template>
  <div class="lg:col-span-1">
    <div
      class="bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))] border border-[rgba(139,92,246,0.2)] rounded-2xl p-6 lg:sticky lg:top-24 space-y-6">
      <div v-if="user.bio">
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="Fa7SolidUserLarge" class="mr-2" /> {{ $t('profile.info.bio') }}
        </h3>
        <p class="mt-3 text-gray-300 text-sm leading-relaxed whitespace-pre-wrap break-words">
          {{ user.bio }}
        </p>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="Fa7SolidCircleInfo" class="mr-2" /> {{ $t('profile.info.basic_info') }}
        </h3>
        <div class="mt-3 space-y-2">
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerCalendarPlus" class="text-purple-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">{{ $t('profile.info.join_date') }}</span>
            </div>
            <span class="text-sm font-bold gaming-title">{{ usePrettyDate(user.joinTime) }}</span>
          </div>
          <!-- Mock Data for design completeness -->
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerClockHour4" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">{{ $t('profile.info.play_time') }}</span>
            </div>
            <span class="text-sm font-bold gaming-title">{{ prettyPlayTime(600000) }}</span>
          </div>
          <!-- <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerMail" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Email</span>
            </div>
            <span class="text-sm font-bold gaming-title truncate max-w-[150px]">{{ user.email || 'Hidden' }}</span>
          </div> -->
          <div v-if="user.location"
            class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerMapPinFilled" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">{{ $t('profile.info.location') }}</span>
            </div>
            <span class="text-sm font-bold gaming-title">{{ user.location || 'Unknown' }}</span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="MdiChartLine" class="mr-2" /> Level
        </h3>
        <div
          class="mt-3 p-4 rounded-xl bg-gradient-to-br from-purple-500/10 to-pink-500/10 border border-purple-500/30">
          <div class="flex items-center justify-between mb-2">
            <div class="flex items-center gap-2">
              <span class="text-2xl font-bold text-purple-300 gaming-title">{{ getCurrentLevel(user.score) }}</span>
              <Icon name="TablerArrowRight" class="text-purple-400 text-xs" />
              <span class="text-xl font-bold text-gray-400 gaming-title">{{ getCurrentLevel(user.score) + 1 }}</span>
            </div>
            <span class="text-xs font-bold text-purple-300">{{ getLevelProgress(user.score) }}% to next level</span>
          </div>
          <div class="relative w-full h-3 bg-gray-700/50 rounded-full overflow-hidden">
            <div
              class="absolute inset-0 bg-gradient-to-r from-purple-600 via-pink-500 to-purple-600 rounded-full"
              :style="{ width: getLevelProgress(user.score) + '%' }">
            </div>
          </div>
          <div class="mt-2 flex items-center justify-between text-xs text-gray-400">
            <span>
              {{ currentExpInLevel }} / {{ expNeededForNextLevel }} XP
            </span>
            <span v-if="daysStreak > 0" class="flex items-center gap-1">
              <Icon name="RiFireLine" class="text-orange-400" />
              {{ daysStreak }} day streak
            </span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="TablerDeviceGamepad2" class=" mr-2" /> Main Games
        </h3>
        <div class="flex flex-wrap gap-2 mt-3">
          <span class="text-gray-500 text-sm italic">No game info provided</span>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="MaterialSymbolsTrophy" class="mr-2" /> Achievements
        </h3>
        <div class="mt-3 grid grid-cols-3 gap-2">
          <span v-if="achievements.length === 0"
            class="col-span-3 text-gray-500 text-sm italic whitespace-nowrap">
            No achievements
          </span>

          <AchievementCard
            v-for="name in achievements"
            :key="name"
            :name="name" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { getCurrentLevel, getLevelProgress, getTotalExpForLevel } from '@/composables/exp'
import AchievementCard from './AchievementCard.vue'

const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
})

const currentExpInLevel = computed(() => {
  const score = props.user?.score || 0
  const currentLevel = getCurrentLevel(score)
  const prevTotal = getTotalExpForLevel(currentLevel - 1)
  return Math.max(0, score - prevTotal)
})

const achievements = []

const expNeededForNextLevel = computed(() => {
  const score = props.user?.score || 0
  const currentLevel = getCurrentLevel(score)
  const nextTotal = getTotalExpForLevel(currentLevel)
  const prevTotal = getTotalExpForLevel(currentLevel - 1)
  return Math.max(1, nextTotal - prevTotal)
})

const daysStreak = computed(() => {
  return 0
})

function prettyPlayTime(seconds) {
  const hours = seconds / 3600
  return `${hours.toFixed(2)}h`
}
</script>
