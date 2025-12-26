<template>
  <div class="relative inline-block">
    <button
      class="group flex items-center gap-1 transition-all duration-200"
      :class="liked ? 'text-blue-400' : 'text-gray-400 hover:text-blue-400'"
      @click="handleClick">
      <Icon name="MaterialSymbolsThumbUp" />
      <span>{{ count }}</span>
    </button>
    <div ref="floatNumbers" class="absolute inset-0 pointer-events-none overflow-visible"></div>
  </div>
</template>

<script setup lang="ts">

type Props = {
  liked: boolean
  count: number
  onClick?: () => Promise<void>
}

const props = defineProps<Props>()

import { ref } from 'vue'

const floatNumbers = ref<HTMLElement | null>(null)

async function handleClick(event: MouseEvent) {
  if (typeof props.onClick === 'function') {
    await props.onClick()
  }
  if (props.liked) {
    createFloatingNumber()
  }
}

function createFloatingNumber() {
  const container = floatNumbers.value
  if (!container) return

  const floatNumber = document.createElement('div')
  floatNumber.textContent = '+1'
  floatNumber.className = 'absolute text-blue-400 font-bold text-2xl pointer-events-none'
  floatNumber.style.left = '50%'
  floatNumber.style.top = '0'
  floatNumber.style.transform = 'translateX(-50%)'
  floatNumber.style.opacity = '1'
  floatNumber.style.transition = 'transform 800ms ease-out, opacity 800ms ease-out'

  container.appendChild(floatNumber)

  // animate
  requestAnimationFrame(() => {
    floatNumber.style.transform = 'translate(-50%, -30px) scale(0.5)'
    floatNumber.style.opacity = '0'
  })

  // remove after animation
  setTimeout(() => {
    if (floatNumber.parentNode === container) container.removeChild(floatNumber)
  }, 800)
}

</script>