<template>
  <ConfirmDialog />
  <Toast />

  <div class="min-h-screen flex flex-col text-gray-100 font-sans">
    <GamingNavbar />

    <main class="flex-1 flex flex-col mx-auto w-full max-w-7xl px-4 sm:px-6 py-6 py-6">
      <slot />
    </main>

    <GamingFooter />
  </div>
</template>


<script setup>
import GamingNavbar from '~/components/GamingNavbar.vue'
import GamingFooter from '~/components/GamingFooter.vue'
import ConfirmDialog from '~/components/ui/ConfirmDialog.vue'
import Toast from '~/components/ui/Toast.vue'
</script>
