<script setup>
import GamingNavbar from '~/components/GamingNavbar.vue'
import GamingFooter from '~/components/GamingFooter.vue'
</script>

<template>
  <ConfirmDialog />
  <Toast />
  <div class="text-gray-100 flex flex-col min-h-screen font-sans">
    <GamingNavbar />

    <main class="flex-grow">
      <slot />
    </main>

    <GamingFooter />
  </div>
</template>

<script>
import ConfirmDialog from '~/components/ui/ConfirmDialog.vue'
import Toast from '~/components/ui/Toast.vue'
</script>

<style scoped>
.bg-pattern {
  background-color: #1a1a2e;
  background-image:
    radial-gradient(at 47% 33%, hsl(240, 33%, 20%) 0, transparent 59%),
    radial-gradient(at 82% 65%, hsl(267, 65%, 15%) 0, transparent 55%);
}
</style>
