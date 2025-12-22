<template>
  <img
    :src="currentSrc"
    :alt="altText"
    :width="size"
    :height="size"
    :class="classes"
    @error="onError" />
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue'

const props = defineProps({
  username: {
    type: String,
    default: '',
  },
  src: {
    type: String,
    default: '',
  },
  size: {
    type: [Number, String],
    default: 48,
  },
  rounded: {
    type: Boolean,
    default: false,
  },
})


// internal defaults
const fallback = '/images/steve.png'
const altText = computed(() => (props.username ? `${props.username}` : 'avatar'))

const sizeParam = computed(() => (typeof props.size === 'number' ? props.size : parseInt(String(props.size)) || 48))

const avatarSrc = computed(() => {
  if (props.src) return props.src
  if (props.username) {
    return `https://minotar.net/avatar/${encodeURIComponent(props.username)}/${sizeParam.value}`
  }
  return fallback
})

// use a local mutable src so that when the image errors we can keep the fallback
const currentSrc = ref(avatarSrc.value)

watch(avatarSrc, (v) => {
  currentSrc.value = v
})

onMounted(() => {
  currentSrc.value = avatarSrc.value
})

const classes = computed(() => {
  const cls = ['object-cover', 'select-none', 'bg-gray-200']
  if (props.rounded) cls.push('rounded-full')
  else cls.push('rounded-none')
  return cls.join(' ')
})

function onError(e) {
  const t = e.target
  if (t && t.src && t.src !== fallback) {
    currentSrc.value = fallback
    t.src = fallback
  }
}
</script>

<style scoped></style>
