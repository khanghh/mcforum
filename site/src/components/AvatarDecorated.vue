<template>
  <div class="relative items-center p-1">
    <!-- Frame Decorations -->
    <div class="absolute -top-1 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-10">
      <div v-if="roleName"
        :class="`bg-gradient-to-r ${roleColor} px-2 py-0.5 border-2 border-gray-900 flex items-center gap-1`">
        <Icon v-if="roleIcon" :name="roleIcon" class="text-white text-xs" />
        <span class="text-white text-xs font-bold gaming-title capitalize">{{ roleName }}</span>
      </div>
    </div>

    <!-- Square animated gradient frame -->
    <div
      class="animated-border avatar-frame w-24 h-24 sm:w-32 sm:h-32 rounded-xl relative z-10 shadow-2xl flex items-center justify-center">
      <div class="avatar-inner w-full h-full bg-white rounded-lg overflow-hidden flex items-center justify-center">
        <slot name="default">
          <img src="/images/avatars/steve.png" class="w-full h-full object-cover" />
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
// Avatar removed: consumers should provide avatar via default slot

const props = defineProps({
  user: { type: Object, required: true },
})

const roleName = computed(() => props.user.role || '')

const roleIcon = computed(() => {
  switch (props.user.role) {
    case 'owner':
      return 'TablerCrown'
    default:
      return null
  }
})

const roleColor = computed(() => {
  switch (props.user.role) {
    case 'owner':
      return 'from-amber-500 to-orange-500'
    case 'admin':
      return 'from-red-500 to-pink-500'
    case 'moderator':
      return 'from-blue-500 to-cyan-500'
    default:
      return ''
  }
})


</script>

<style scoped>
/* Animated gradient border */
@keyframes gradientBorder {
  0% {
    background-position: 0% 50%;
  }

  50% {
    background-position: 100% 50%;
  }

  100% {
    background-position: 0% 50%;
  }
}

.animated-border {
  background: linear-gradient(45deg, #8b5cf6, #ec4899, #06b6d4, #8b5cf6);
  background-size: 300% 300%;
  animation: gradientBorder 3s ease infinite;
}

/* Square avatar frame styles */
.avatar-frame {
  display: inline-block;
  padding: 0.25rem;
  /* small inner gap for the animated border */
  border-radius: 0.75rem;
}

.avatar-inner {
  border-radius: 0.5rem;
  overflow: hidden;
}
</style>
