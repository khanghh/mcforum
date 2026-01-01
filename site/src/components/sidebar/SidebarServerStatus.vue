<template>
  <div class="gaming-card rounded-xl p-4">
    <h3 class="text-lg font-bold text-green-400 mb-4 flex items-center gaming-title">
      <Icon name="TablerServer" class="mr-2" /> SERVER STATUS
    </h3>
    <div class="space-y-3">
      <div class="flex justify-between items-center p-2 rounded bg-green-500/10">
        <span class="text-sm text-gray-400">Status</span>
        <span class="px-2 py-1 bg-green-500 text-white text-xs font-bold rounded gaming-title">
          ONLINE
        </span>
      </div>
      <div class="flex justify-between items-center p-2 rounded bg-blue-500/10">
        <span class="text-sm text-gray-400">Online</span>
        <span class="font-bold text-green-400 gaming-title">2,847</span>
      </div>
      <div class="flex justify-between items-center p-2 rounded bg-purple-500/10">
        <span class="text-sm text-gray-400">Version</span>
        <span class="font-bold text-purple-400 gaming-title">1.20.4</span>
      </div>
    </div>
    <a href="https://play.mineviet.com"
      class="w-full mt-4 px-4 py-3
             bg-gradient-to-r from-green-600 to-emerald-600
             text-white rounded-lg font-bold gaming-title
             hover:scale-105 transition-transform
             flex items-center justify-center">
      <Icon name="TablerPlayerPlayFilled" class="mr-2" /> JOIN SERVER
    </a>
  </div>
</template>

<script setup lang="ts">
const dialog = useConfirmDialog()
const toast = useToast()
</script>
