<template>
  <div :class="['relative inline-block', customClass]" :style="{ width: sizePx, height: sizePx }">
    <div class="absolute inset-0 flex items-center justify-center">
      <div :style="{ width: innerSizePx, height: innerSizePx }" class="overflow-hidden">
        <Avatar
          :username="username"
          :src="src"
          :size="innerSize"
          class="w-full h-full" />
      </div>
    </div>

    <img
      :src="frame"
      alt="avatar-frame"
      class="pointer-events-none absolute inset-0 w-full h-full" />
  </div>
</template>

<script setup>
import Avatar from './Avatar.vue'

const props = defineProps({
  src: { type: String, default: '' },
  username: { type: String, default: '' },
  size: { type: [Number, String], default: 48 },
  rounded: { type: Boolean, default: true },
  frame: { type: String, default: '/frames/member.png' },
  customClass: { type: String, default: '' },
})

const sizePx = computed(() => (typeof props.size === 'number' ? props.size + 'px' : (parseInt(String(props.size)) || 48) + 'px'))

// Avatar inner size slightly smaller than the frame (86%)
const innerSize = computed(() => {
  const n = typeof props.size === 'number' ? props.size : (parseInt(String(props.size)) || 48)
  return Math.max(24, Math.round(n * 0.7))
})

const innerSizePx = computed(() => `${innerSize.value}px`)
</script>

<style scoped>
.pointer-events-none {
  pointer-events: none;
}

img[alt="avatar-frame"] {
  object-fit: contain;
}
</style>
