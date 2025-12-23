<script setup>
import GamingNavbar from '~/components/GamingNavbar.vue'
import GamingFooter from '~/components/GamingFooter.vue'
import GamingSidebar from '~/components/GamingSidebar.vue'
</script>

<template>
  <div class="text-gray-100 flex flex-col min-h-screen custom-scrollbar font-sans">
    <GamingNavbar />

    <main class="flex-grow">
      <slot />
    </main>

    <GamingFooter />
  </div>
</template>

<style scoped>
/* Layout uses Tailwind; no extra styles here. */
</style>
