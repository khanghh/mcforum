<template>
  <Teleport to="body">
    <div class="fixed top-40 right-0 z-[9999] flex flex-col gap-3 pointer-events-none max-w-md w-full px-4 sm:px-0">
      <TransitionGroup
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
        move-class="transition duration-200 ease-in-out">
        <ToastItem
          v-for="toast in toasts"
          :key="toast.id"
          :toast="toast" />
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import ToastItem from './ToastItem.vue'

const { toasts } = useToast()
</script>
