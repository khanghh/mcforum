<template>
  <aside class="hidden lg:block w-80 flex-shrink-0">
    <div class="space-y-6 sticky top-24">
      <SidebarCategories />
      <SidebarServerStatus />
      <SidebarTopContributors />
    </div>
  </aside>
</template>

<script setup>
import SidebarCategories from './sidebar/SidebarCategories.vue'
import SidebarServerStatus from './sidebar/SidebarServerStatus.vue'
import SidebarTopContributors from './sidebar/SidebarTopContributors.vue'
</script>
