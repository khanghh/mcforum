<template>
  <footer class="bg-gradient-to-b from-transparent to-purple-900/30 border-t border-purple-500/20 mt-16">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-12">
      <!-- Main Footer Content -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
        <!-- Brand Section -->
        <div class="space-y-4">
          <div class="flex items-center space-x-2">
            <FontAwesome :icon="['fas', 'gamepad']" class="text-3xl text-purple-400" />
            <div>
              <h3
                class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 server-title">
                MineViet Forum
              </h3>
            </div>
          </div>
          <p class="text-gray-400 text-sm">
            The ultimate community for Minecraft players. Share builds, discover mods, and connect with fellow gamers.
          </p>
          <div class="flex space-x-3">
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <FontAwesome :icon="['fab', 'discord']" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <FontAwesome :icon="['fab', 'twitter']" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <FontAwesome :icon="['fab', 'youtube']" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <FontAwesome :icon="['fab', 'twitch']" />
            </a>
          </div>
        </div>

        <!-- Quick Links -->
        <div>
          <h4 class="text-lg font-bold text-purple-300 gaming-title mb-4">QUICK LINKS</h4>
          <ul class="space-y-2 text-sm">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Home</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Forums</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Leaderboards</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Staff Team</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Server Status</a></li>
          </ul>
        </div>

        <!-- Resources -->
        <div>
          <h4 class="text-lg font-bold text-purple-300 gaming-title mb-4">RESOURCES</h4>
          <ul class="space-y-2 text-sm">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Wiki</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Tutorials</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Download Packs</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Rules & Guidelines</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Help Center</a></li>
          </ul>
        </div>

        <!-- Community Stats -->
        <div>
          <h4 class="text-lg font-bold text-amber-400 gaming-title mb-4">COMMUNITY</h4>
          <div class="space-y-3">
            <div class="gaming-card p-3 rounded-lg">
              <div class="text-xs text-gray-400 mb-1">Active Players</div>
              <div class="text-2xl font-bold text-green-400 gaming-title">2,847</div>
            </div>
            <div class="gaming-card p-3 rounded-lg">
              <div class="text-xs text-gray-400 mb-1">Total Members</div>
              <div class="text-2xl font-bold text-purple-400 gaming-title">458K</div>
            </div>
            <div class="gaming-card p-3 rounded-lg">
              <div class="text-xs text-gray-400 mb-1">Online Now</div>
              <div class="text-2xl font-bold text-blue-400 gaming-title">2.4K</div>
            </div>
          </div>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="pt-8 border-t border-purple-500/20">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-gray-400 text-sm gaming-title">
            © 2025 MineViet Network. All rights reserved.
          </p>
          <div class="flex flex-wrap justify-center gap-4 text-sm text-gray-400">
            <a href="#" class="hover:text-purple-400 transition-colors">Privacy Policy</a>
            <span class="text-gray-600">•</span>
            <a href="#" class="hover:text-purple-400 transition-colors">Terms of Service</a>
            <span class="text-gray-600">•</span>
            <a href="#" class="hover:text-purple-400 transition-colors">Cookie Policy</a>
            <span class="text-gray-600">•</span>
            <a href="#" class="hover:text-purple-400 transition-colors">Contact Us</a>
          </div>
        </div>
      </div>
    </div>
  </footer>
</template>
