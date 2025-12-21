<template>
  <nav
    class="bg-gradient-to-r from-purple-900/30 to-indigo-900/30 backdrop-blur-lg border-b border-purple-500/20 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-3 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <img src="/images/logo.png" class="w-12 h-12" />
        <a href="/">
          <div
            class="text-2xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 server-title ">
            MineViet Forum
          </div>
          <div class="hidden md:block text-xs text-purple-300">Minecraft Server Vietnam</div>
        </a>
      </div>
      <ul class="flex space-x-4 sm:space-x-6 text-gray-300 font-medium">
        <li>
          <a href="https://play.mineviet.com" class="text-gray-300 hover:text-purple-400 transition-colors"
            target="_blank">
            <FontAwesome :icon="['fas', 'server']" class="mr-1" />
            <span class="hidden sm:inline">Máy chủ</span>
          </a>
        </li>
        <li>
          <a href="https://skin.mineviet.com" class="text-gray-300 hover:text-purple-400 transition-colors">
            <FontAwesome :icon="['fas', 'palette']" class="mr-1" />
            <span class="hidden sm:inline">Đổi skin</span>
          </a>
        </li>
        <li>
          <a href="https://ban.mineviet.com" class="text-gray-300 hover:text-purple-400 transition-colors">
            <FontAwesome :icon="['fas', 'ban']" class="mr-1" />
            <span class="hidden sm:inline">Vi phạm</span>
          </a>
        </li>
        <li>
          <a href="#" class="text-gray-300 hover:text-purple-400 transition-colors">
            <FontAwesome :icon="['fas', 'envelope']" class="mr-1" />
            <span class="hidden sm:inline">Tin nhắn</span>
          </a>
        </li>
        <li>
          <a href="#" class="text-gray-300 hover:text-purple-400 transition-colors">
            <FontAwesome :icon="['fas', 'user-circle']" class="mr-1" />
            <span class="hidden sm:inline">Profile</span>
          </a>
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
// Logic for active link etc. can be added here
</script>
