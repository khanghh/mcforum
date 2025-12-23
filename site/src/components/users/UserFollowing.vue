<template>
  <div class="py-10 text-center text-gray-400">Following list implementation coming soon.</div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { UserProfile } from '@/types'

const userStore = useUserStore()
const api = useApi()

interface UserFollowingProps {
  user: UserProfile
}
const props = defineProps<UserFollowingProps>()



</script>