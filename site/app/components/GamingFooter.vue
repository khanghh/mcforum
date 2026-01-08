<template>
  <footer class="bg-gradient-to-b from-transparent to-purple-900/30 border-t border-purple-500/20 mt-6">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 py-12">
      <!-- Main Footer Content -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8">
        <!-- Brand Section -->
        <div class="space-y-4">
          <div class="flex items-center space-x-2">
            <Icon name="TablerDeviceGamepad2" class=" text-3xl text-purple-400" />
            <div>
              <div
                class="text-xl font-bold text-transparent bg-clip-text bg-gradient-to-r from-purple-400 to-pink-400 server-title">
                MineViet Forum
              </div>
            </div>
          </div>
          <p class="text-gray-400 text-sm">
            Chào mừng đến <strong>MineViet</strong> – cộng đồng Minecraft Việt Nam nơi sáng tạo không giới hạn.
            Xây dựng, phiêu lưu và chiến đấu cùng bạn bè!
          </p>
          <div class="flex space-x-3">
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <Icon name="TablerBrandDiscordFilled" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <Icon name="TablerBrandFacebookFilled" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <Icon name="TablerBrandTwitterFilled" />
            </a>
            <a href="#"
              class="w-10 h-10 bg-purple-500/20 rounded-lg flex items-center justify-center text-purple-400 hover:bg-purple-500/30 transition-colors">
              <Icon name="TablerBrandYoutubeFilled" />
            </a>
          </div>
        </div>

        <!-- Quick Links -->
        <div>
          <div class="text-lg font-bold text-purple-300 gaming-title mb-4">QUICK LINKS</div>
          <ul class="space-y-2 text-sm">
            <li>
              <a href="https://mineviet.com"
                class="text-gray-400 hover:text-purple-400 transition-colors"
                title="Trang chủ MineViet">
                Home
              </a>
            </li>
            <li>
              <a href="https://skin.mineviet.com"
                class="text-gray-400 hover:text-purple-400 transition-colors"
                title="Đổi skin Minecraft trên MineViet">
                Đổi skin
              </a>
            </li>
            <li>
              <a href="/staffs"
                class="text-gray-400 hover:text-purple-400 transition-colors"
                title="Xem danh sách Staff MineViet">
                Staff Team
              </a>
            </li>
            <li>
              <a href="https://play.mineviet.com"
                class="text-gray-400 hover:text-purple-400 transition-colors"
                title="Xem trạng thái và tham gia máy chủ MineViet">
                Server Status
              </a>
            </li>
            <li>
              <a href="https://ban.mineviet.com"
                class="text-gray-400 hover:text-purple-400 transition-colors"
                title="Xem danh sách vi phạm">
                Danh sách vi phạm
              </a>
            </li>
          </ul>
        </div>

        <!-- Resources -->
        <div>
          <div class="text-lg font-bold text-purple-300 gaming-title mb-4">RESOURCES</div>
          <ul class="space-y-2 text-sm">
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Wiki</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Tutorials</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Download Packs</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Rules & Guidelines</a></li>
            <li><a href="#" class="text-gray-400 hover:text-purple-400 transition-colors">Help Center</a></li>
          </ul>
        </div>

        <!-- Forum Stats (moved from sidebar) -->
        <div>
          <FooterForumStats />
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="pt-8 border-t border-purple-500/20">
        <div class="flex flex-col md:flex-row justify-between items-center gap-4">
          <p class="text-gray-400 text-sm gaming-title">
            © 2025 MineViet Network. All rights reserved.
          </p>
          <div class="flex flex-wrap justify-center gap-4 text-gray-400 text-sm gaming-title">
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
