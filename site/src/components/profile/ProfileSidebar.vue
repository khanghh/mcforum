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
            <span class="text-sm font-bold gaming-title">{{ formatDate(user.createTime) }}</span>
          </div>
          <!-- Mock Data for design completeness -->
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerClockHour4" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">{{ $t('profile.info.play_time') }}</span>
            </div>
            <span class="text-sm font-bold gaming-title">2,847 hrs</span>
          </div>
          <!-- <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
            <div class="flex items-center gap-2">
              <Icon name="TablerMail" class="text-blue-400 text-sm" />
              <span class="text-sm font-bold text-gray-400">Email</span>
            </div>
            <span class="text-sm font-bold gaming-title truncate max-w-[150px]">{{ user.email || 'Hidden' }}</span>
          </div> -->
          <div class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5 transition-colors">
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
              <span class="text-2xl font-bold text-purple-300 gaming-title">{{ calculateLevel(user.score) }}</span>
              <Icon name="TablerArrowRight" class="text-purple-400 text-xs" />
              <span class="text-xl font-bold text-gray-400 gaming-title">{{ calculateLevel(user.score) + 1 }}</span>
            </div>
            <span class="text-xs font-bold text-purple-300">{{ calculateProgress(user.score) }}% to next level</span>
          </div>
          <div class="relative w-full h-3 bg-gray-700/50 rounded-full overflow-hidden">
            <div
              class="absolute inset-0 bg-gradient-to-r from-purple-600 via-pink-500 to-purple-600 rounded-full"
              :style="{ width: calculateProgress(user.score) + '%' }">
            </div>
          </div>
          <div class="mt-2 flex items-center justify-between text-xs text-gray-400">
            <span>{{ (user.score - totalExp(calculateLevel(user.score) - 1)) }} / {{
              totalExp(calculateLevel(user.score)) - totalExp(calculateLevel(user.score) - 1) }} XP</span>
            <span class="flex items-center gap-1">
              <Icon name="RiFireLine" class="text-orange-400" />
              5 day streak
            </span>
          </div>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="TablerDeviceGamepad2" class=" mr-2" /> Main Games
        </h3>
        <div class="flex flex-wrap gap-2 mt-3">
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-blue-500/20 text-blue-300 border border-blue-500/30">
            <Icon name="MdiCrosshairsGps" class="mr-1.5" /> SkyBlock
          </span>
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-purple-500/20 text-purple-300 border border-purple-500/30">
            <Icon name="Fa7SolidChess" class="mr-1.5" /> Survival
          </span>
          <span
            class="inline-flex items-center px-3 py-1.5 rounded-lg text-xs font-bold bg-orange-500/20 text-orange-300 border border-orange-500/30">
            <Icon name="TablerBombFilled" class="mr-1.5" /> Bedwars
          </span>
        </div>
      </div>

      <div>
        <h3 class="text-lg font-bold text-purple-300 gaming-title flex items-center">
          <Icon name="MaterialSymbolsTrophy" class="mr-2" /> Achievements
        </h3>
        <div class="mt-3 grid grid-cols-3 gap-2">
          <!-- Static Achievements List -->
          <div class="flex flex-col items-center gap-2 p-3 rounded-lg bg-amber-500/10 border border-amber-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center flex-shrink-0">
              <Icon name="TablerCrown" class="text-white text-sm" />
            </div>
            <div class="text-center">
              <p class="text-[10px] font-bold text-amber-300 gaming-title">CHAMPION</p>
            </div>
          </div>
          <div
            class="flex flex-col items-center gap-2 p-3 rounded-lg bg-purple-500/10 border border-purple-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center flex-shrink-0">
              <Icon name="TablerStar" class="text-white text-sm" />
            </div>
            <div class="text-center">
              <p class="text-[10px] font-bold text-purple-300 gaming-title">LEGEND</p>
            </div>
          </div>
          <div class="flex flex-col items-center gap-2 p-3 rounded-lg bg-blue-500/10 border border-blue-500/30">
            <div
              class="w-10 h-10 rounded-lg bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center flex-shrink-0">
              <Icon name="RiFireLine" class="text-white text-sm" />
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
  let level = 0;
  while (totalExp(level + 1) <= (score || 0)) {
    level++;
  }
  return level + 1;
}

function calculateProgress(score) {
  const currentLevel = calculateLevel(score);
  const currentLevelExp = totalExp(currentLevel - 1);
  const nextLevelExp = totalExp(currentLevel);
  const expInCurrentLevel = score - currentLevelExp;
  const expNeededForNextLevel = nextLevelExp - currentLevelExp;

  return Math.min(100, Math.round((expInCurrentLevel / expNeededForNextLevel) * 100));
}

function expToNext(level) {
  const base = 10;
  const growth = 5;

  let exp = base + level * growth;

  if (level > 500) {
    exp = exp * Math.pow(1.01, level - 500);
  }

  return Math.floor(exp); // round down to integer
}

function totalExp(level) {
  let total = 0;
  for (let i = 1; i <= level; i++) {
    total += expToNext(i);
  }
  return total;
}

</script>
