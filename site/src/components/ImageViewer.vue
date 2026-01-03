<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-[9999] flex items-center justify-center bg-black/50 backdrop-blur-sm"
        @click="close">
        <!-- Close Button -->
        <button
          class="absolute top-4 right-4 z-50 p-2 text-white/70 hover:text-white bg-black/20 hover:bg-white/10 rounded-full transition-all"
          @click.stop="close">
          <Icon name="Fa7SolidTimes" class="w-6 h-6" />
        </button>

        <!-- Navigation Buttons -->
        <button
          v-if="hasMultiple"
          class="absolute left-4 z-50 p-3 text-white/70 hover:text-white bg-black/20 hover:bg-white/10 rounded-full transition-all disabled:opacity-30 disabled:cursor-not-allowed"
          :disabled="currentIndex === 0"
          @click.stop="prev">
          <Icon name="Fa7SolidChevronLeft" class="w-6 h-6" />
        </button>

        <button
          v-if="hasMultiple"
          class="absolute right-4 z-50 p-3 text-white/70 hover:text-white bg-black/20 hover:bg-white/10 rounded-full transition-all disabled:opacity-30 disabled:cursor-not-allowed"
          :disabled="currentIndex === images.length - 1"
          @click.stop="next">
          <Icon name="Fa7SolidChevronRight" class="w-6 h-6" />
        </button>

        <!-- Image -->
        <div class="relative w-full h-full flex items-center justify-center p-4">
          <img
            :src="currentImage"
            class="max-w-full max-h-full object-contain rounded-lg shadow-2xl select-none"
            alt="Full screen preview"
            draggable="false"
            @click.stop />

          <!-- Counter -->
          <div v-if="hasMultiple"
            class="absolute bottom-6 left-1/2 -translate-x-1/2 px-3 py-1 bg-black/50 rounded-full text-white text-sm font-medium"
            @click.stop>
            {{ currentIndex + 1 }} / {{ images.length }}
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  images: {
    type: Array,
    default: () => [],
  },
  index: {
    type: Number,
    default: 0,
  },
})

const emit = defineEmits(['update:modelValue', 'update:index'])

const currentIndex = ref(props.index)

watch(() => props.index, (val) => {
  currentIndex.value = val
})

watch(() => props.modelValue, (val) => {
  if (val) {
    document.body.style.overflow = 'hidden'
    currentIndex.value = props.index
  } else {
    document.body.style.overflow = ''
  }
})

const currentImage = computed(() => {
  if (!props.images || props.images.length === 0) return ''
  const img = props.images[currentIndex.value]
  if (typeof img === 'string') return img
  return img.url || img.preview
})

const hasMultiple = computed(() => props.images && props.images.length > 1)

function close() {
  emit('update:modelValue', false)
}

function next() {
  if (currentIndex.value < props.images.length - 1) {
    currentIndex.value++
    emit('update:index', currentIndex.value)
  }
}

function prev() {
  if (currentIndex.value > 0) {
    currentIndex.value--
    emit('update:index', currentIndex.value)
  }
}

function onKeydown(e) {
  if (!props.modelValue) return

  if (e.key === 'Escape') {
    close()
  } else if (e.key === 'ArrowLeft') {
    prev()
  } else if (e.key === 'ArrowRight') {
    next()
  }
}

onMounted(() => {
  document.addEventListener('keydown', onKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', onKeydown)
  document.body.style.overflow = ''
})
</script>
