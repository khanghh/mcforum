<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 py-8">
    <div
      class="gaming-card sm:rounded-xl p-6 border border-purple-500/20 bg-[linear-gradient(145deg,rgba(30,30,60,0.8),rgba(20,20,40,0.9))]">
      <h1 class="text-xl sm:text-3xl font-bold text-white gaming-title">
        CHÍNH SÁCH COOKIE
      </h1>
      <p class="mt-2 text-sm text-gray-400">
        Cập nhật lần cuối: 09/01/2026
      </p>

      <p class="mt-6 text-gray-300 leading-relaxed">
        Chính sách Cookie này giải thích cách MineViet (“chúng tôi”) sử dụng cookie và các công nghệ tương tự trên
        website.
      </p>

      <div class="mt-8 space-y-8 text-gray-300">
        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">1. Cookie là gì?</h2>
          <p class="mt-3 leading-relaxed">
            Cookie là các tệp văn bản nhỏ được lưu trên trình duyệt của bạn khi truy cập website. Cookie giúp website
            hoạt
            động đúng và cải thiện trải nghiệm người dùng.
          </p>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">2. Cookie chúng tôi sử dụng</h2>
          <p class="mt-3 leading-relaxed">Website của chúng tôi có thể sử dụng các loại cookie sau:</p>

          <div class="mt-4 space-y-4">
            <div>
              <h3 class="font-bold text-white">a) Cookie cần thiết</h3>
              <ul class="mt-2 list-disc pl-5 space-y-1">
                <li>Dùng để website hoạt động bình thường</li>
                <li>Ghi nhớ phiên đăng nhập, cài đặt cơ bản</li>
                <li>Không thể tắt nếu không ảnh hưởng đến chức năng website</li>
              </ul>
            </div>

            <div>
              <h3 class="font-bold text-white">b) Cookie phân tích</h3>
              <ul class="mt-2 list-disc pl-5 space-y-1">
                <li>Thu thập dữ liệu ẩn danh về cách người dùng sử dụng website</li>
                <li>Ví dụ: số lượt truy cập, trang được xem</li>
                <li>Dùng các công cụ như Google Analytics (nếu có)</li>
              </ul>
            </div>

            <div>
              <h3 class="font-bold text-white">c) Cookie quảng cáo</h3>
              <ul class="mt-2 list-disc pl-5 space-y-1">
                <li>Dùng để hiển thị quảng cáo phù hợp</li>
                <li>Có thể được đặt bởi Google AdSense hoặc đối tác quảng cáo bên thứ ba</li>
              </ul>
            </div>
          </div>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">3. Cookie của bên thứ ba</h2>
          <p class="mt-3 leading-relaxed">
            Một số cookie có thể được đặt bởi dịch vụ bên thứ ba, bao gồm nhưng không giới hạn:
          </p>
          <ul class="mt-3 list-disc pl-5 space-y-1">
            <li>Google Analytics</li>
            <li>Google AdSense</li>
            <li>Cloudflare hoặc dịch vụ bảo mật tương tự</li>
          </ul>
          <p class="mt-3 leading-relaxed">
            Các bên này có chính sách cookie và quyền riêng tư riêng mà chúng tôi không kiểm soát.
          </p>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">4. Quản lý cookie</h2>
          <p class="mt-3 leading-relaxed">Bạn có thể:</p>
          <ul class="mt-3 list-disc pl-5 space-y-1">
            <li>Chấp nhận hoặc từ chối cookie thông qua cài đặt trình duyệt</li>
            <li>Xóa cookie đã lưu bất cứ lúc nào</li>
          </ul>
          <p class="mt-3 leading-relaxed">
            Lưu ý: Việc vô hiệu hóa cookie có thể làm một số chức năng website không hoạt động đúng.
          </p>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">5. Sự đồng ý của bạn</h2>
          <p class="mt-3 leading-relaxed">
            Bằng việc tiếp tục sử dụng website, bạn đồng ý với việc sử dụng cookie theo Chính sách này. Trong trường hợp
            pháp luật yêu cầu, website có thể hiển thị thông báo chấp thuận cookie.
          </p>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">6. Thay đổi chính sách</h2>
          <p class="mt-3 leading-relaxed">
            Chúng tôi có thể cập nhật Chính sách Cookie này theo thời gian. Mọi thay đổi sẽ được công bố trên website.
          </p>
        </section>

        <section>
          <h2 class="text-lg sm:text-xl font-bold text-purple-300 gaming-title">7. Liên hệ</h2>
          <p class="mt-3 leading-relaxed">
            Nếu bạn có câu hỏi liên quan đến Chính sách Cookie, vui lòng liên hệ:
          </p>
          <div class="mt-3 space-y-1">
            <div>
              <span class="text-gray-400">Email:</span>
              <a
                href="mailto:admin@mineviet.com"
                class="ml-1 text-purple-300 hover:text-purple-200 underline underline-offset-2 transition-colors">
                admin@mineviet.com
              </a>
            </div>

            <div>
              <span class="text-gray-400">Fanpage:</span>
              <a :href="socialLinks.facebook || '#'" target="_blank" rel="noopener"
                class="ml-1 text-purple-300 hover:text-purple-200 underline underline-offset-2 transition-colors">
                {{ socialLinks.facebook || '#' }}
              </a>
            </div>
            <div>
              <span class="text-gray-400">Discord:</span>
              <a :href="socialLinks.discord || '#'" target="_blank" rel="noopener"
                class="ml-1 text-purple-300 hover:text-purple-200 underline underline-offset-2 transition-colors">
                {{ socialLinks.discord || '#' }}
              </a>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useConfigStore } from '~/stores/config'

const configStore = useConfigStore()
const socialLinks = computed(() => configStore.config?.socialLinks || {})

useHead({
  title: useSiteTitle('Chính sách Cookie'),
})
</script>
