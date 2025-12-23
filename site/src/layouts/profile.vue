<script setup>
import GamingNavbar from '~/components/GamingNavbar.vue'
import GamingFooter from '~/components/GamingFooter.vue'
</script>

<template>
  <div
    class="min-h-screen bg-gray-900 text-gray-100 font-sans selection:bg-purple-500 selection:text-white"
    style="background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);">
    <GamingNavbar />

    <main class="flex-grow py-6">
      <div class="max-w-7xl mx-auto px-4 sm:px-6">
        <slot />
      </div>
    </main>

    <GamingFooter />
  </div>
</template>

<style scoped>
/* Layout uses Tailwind; no extra styles here. */
</style>
