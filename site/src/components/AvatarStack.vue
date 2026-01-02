<template>
  <div class="flex">
    <Avatar
      v-for="(user, index) in visibleUsers"
      :key="user.id || index"
      :src="user.avatar"
      :username="user.username"
      :size="size"
      class="border border-purple-300 rounded"
      :class="index > 0 ? '-ml-5' : ''" />
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  users: {
    type: Array,
    default: () => [],
  },
  size: {
    type: Number,
    default: 32,
  },
})

const visibleUsers = computed(() => props.users.slice(0, 4))
</script>
