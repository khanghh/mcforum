<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
    <UserCard v-for="u in users" :key="u.id" :user="u" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { UserProfile } from '@/types'
import UserCard from './UserCard.vue'

// Dummy following users (3 columns layout)
const users = ref<UserProfile[]>([
  { id: 101, username: 'alice', nickname: 'Alice', avatar: '/images/avatars/1.jpg', bio: 'Builder' } as any,
  { id: 102, username: 'bob', nickname: 'Bob', avatar: '/images/avatars/2.jpg', bio: 'Redstone wizard' } as any,
  { id: 103, username: 'carol', nickname: 'Carol', avatar: '/images/avatars/3.jpg', bio: 'Explorer' } as any,
  { id: 104, username: 'dave', nickname: 'Dave', avatar: '/images/avatars/4.jpg', bio: 'PvP master' } as any,
  { id: 105, username: 'eve', nickname: 'Eve', avatar: '/images/avatars/5.jpg', bio: 'Map maker' } as any,
  { id: 106, username: 'frank', nickname: 'Frank', avatar: '/images/avatars/6.jpg', bio: 'Admin' } as any,
])

</script>