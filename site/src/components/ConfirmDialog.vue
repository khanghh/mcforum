<template>
  <div v-if="wrapperVisible" class="fixed inset-0 z-100000 flex items-center justify-center">
    <div class="fixed inset-0 bg-black/50" @click="onCancel"></div>
    <transition
      appear
      @after-leave="onAfterLeave"
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 -translate-y-1 scale-95"
      enter-to-class="opacity-100 translate-y-0 scale-100"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0 scale-100"
      leave-to-class="opacity-0 -translate-y-1 scale-95"
    >
      <div v-if="visible" class="z-10 w-full max-w-lg p-6">
        <div class="rounded-2xl bg-white p-6 shadow-lg dark:bg-gray-900">
          <div class="mb-4">
            <div v-if="title" class="flex items-start space-x-3">
              <span :class="['flex-shrink-0 inline-flex items-center justify-center rounded-full p-2', iconBgClass]">
                <component :is="IconComponent" class="w-6 h-6" aria-hidden="true" />
              </span>
              <div>
                <h3 class="text-lg font-medium text-gray-900 dark:text-white">{{ title }}</h3>
                <p class="mt-2 text-sm text-gray-600 dark:text-gray-300" v-if="message">{{ message }}</p>
              </div>
            </div>
          </div>

          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              class="inline-flex items-center rounded-md border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 dark:border-gray-800 dark:bg-transparent dark:text-gray-200"
              @click="onCancel">
              {{ cancelText }}
            </button>

            <button
              type="button"
              :class="['inline-flex items-center rounded-md px-4 py-2 text-sm font-medium text-white', confirmBtnClass]"
              @click="onConfirm">
              {{ confirmText }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, computed } from 'vue'
import { useConfirmDialog } from '@/composables/useConfirmDialog'
import { PhInfoIcon, PhWarningIcon } from '@/icons'

const dialog = useConfirmDialog()

const { visible, title, message, confirmText, cancelText, variant } = dialog.state

const iconBgClass = computed(() => {
  const v = variant?.value ?? 'info'
  switch (v) {
    case 'warning':
      return 'bg-red-100 text-red-600 dark:bg-red-900 dark:text-red-300'
    default:
      return 'bg-blue-100 text-blue-600 dark:bg-gray-800 dark:text-blue-300'
  }
})

const IconComponent = computed(() => {
  const v = variant?.value ?? 'info'
  return v === 'warning' ? PhWarningIcon : PhInfoIcon
})

const confirmBtnClass = computed(() => {
  const v = variant?.value ?? 'info'
  return v === 'warning' ? 'bg-red-600 hover:bg-red-700' : 'bg-blue-600 hover:bg-blue-700'
})

const wrapperVisible = ref(false)

watch(visible, (v) => {
  if (v) wrapperVisible.value = true
})

function onAfterLeave() {
  wrapperVisible.value = false
}

function onConfirm() {
  dialog.confirm()
}

function onCancel() {
  dialog.cancel()
}

function handleKeydown(e: KeyboardEvent) {
  if ((e.key === 'Escape' || e.key === 'Esc') && wrapperVisible.value) {
    onCancel()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>
