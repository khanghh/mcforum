<template>
  <div v-if="name === 'CHAMPION'"
    class="flex flex-col items-center gap-2 p-3 rounded-lg border bg-amber-500/10 border-amber-500/30">
    <div
      class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-gradient-to-br from-amber-500 to-orange-500">
      <Icon name="TablerCrown" class="text-white text-sm" />
    </div>
    <div class="text-center">
      <p class="text-[10px] font-bold gaming-title text-amber-300">
        CHAMPION
      </p>
    </div>
  </div>

  <div v-else-if="name === 'LEGEND'"
    class="flex flex-col items-center gap-2 p-3 rounded-lg border bg-purple-500/10 border-purple-500/30">
    <div
      class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-gradient-to-br from-purple-500 to-pink-500">
      <Icon name="TablerStar" class="text-white text-sm" />
    </div>
    <div class="text-center">
      <p class="text-[10px] font-bold gaming-title text-purple-300">
        LEGEND
      </p>
    </div>
  </div>

  <div v-else-if="name === 'STREAK'"
    class="flex flex-col items-center gap-2 p-3 rounded-lg border bg-blue-500/10 border-blue-500/30">
    <div
      class="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-gradient-to-br from-blue-500 to-cyan-500">
      <Icon name="RiFireLine" class="text-white text-sm" />
    </div>
    <div class="text-center">
      <p class="text-[10px] font-bold gaming-title text-blue-300">
        STREAK
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  name: string
}>()
</script>
