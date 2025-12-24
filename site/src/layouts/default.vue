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