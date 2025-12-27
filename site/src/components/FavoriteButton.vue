<template>
  <div class="relative inline-block" ref="wrapper">
    <button
      class="flex items-center gap-1 transition-colors duration-200 relative z-10"
      :class="favorited ? 'text-yellow-400' : 'text-gray-400 hover:text-yellow-400'"
      @click="handleClick">
      <Icon :name="favorited ? 'TablerStarFilled' : 'TablerStar'" />
      <span>{{ $t('feed.actions.favorite') }}</span>
    </button>

    <div ref="particles" class="absolute inset-0 pointer-events-none overflow-visible"></div>

  </div>
</template>

<script setup lang="ts">

type Props = {
  favorited: boolean
  onClick?: () => Promise<void>
}

const props = defineProps<Props>()

import { ref } from 'vue'

const particles = ref<HTMLElement | null>(null)

async function handleClick(event: MouseEvent) {
  if (typeof props.onClick === 'function') {
    await props.onClick()
  }
  if (props.favorited) {
    spawnParticles(12)
  }
}

function spawnParticles(count = 10) {
  const container = particles.value
  if (!container) return

  const colors = ['#f59e0b', '#ec4899', '#8b5cf6', '#06b6d4', '#ef4444']

  for (let i = 0; i < count; i++) {
    const p = document.createElement('div')
    p.className = 'absolute pointer-events-none'
    const size = Math.floor(Math.random() * 6) + 6 // 6-11px
    p.style.width = `${size}px`
    p.style.height = `${size}px`
    p.style.backgroundColor = colors[Math.floor(Math.random() * colors.length)]!
    p.style.borderRadius = '50%'
    p.style.left = '50%'
    p.style.top = '50%'
    p.style.transform = 'translate(-50%, -50%) scale(1)'
    p.style.opacity = '1'
    p.style.zIndex = '11'
    p.style.transition = 'transform 700ms cubic-bezier(0.2,0.8,0.2,1), opacity 700ms ease-out'

    container.appendChild(p)

    // calculate a random direction and distance
    const angle = Math.random() * Math.PI * 2
    const distance = 24 + Math.random() * 40
    const dx = Math.cos(angle) * distance
    const dy = Math.sin(angle) * distance - (Math.random() * 8)

    p.getBoundingClientRect()
    // animate on next frame
    requestAnimationFrame(() => {
      p.style.transform = `translate(calc(-50% + ${dx}px), calc(-50% + ${dy}px)) scale(0.6)`
      p.style.opacity = '0'
    })

    // remove after animation
    setTimeout(() => {
      if (p.parentNode === container) container.removeChild(p)
    }, 750)
  }
}

</script>