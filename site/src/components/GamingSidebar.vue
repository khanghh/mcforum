<template>
  <aside class="hidden lg:block w-80 flex-shrink-0">
    <div class="space-y-6 sticky top-24">
      <CategoriesWidget />
      <ServerStatusWidget />
      <TopContributorsWidget />
    </div>
  </aside>
</template>

<script setup>
import CategoriesWidget from './sidebar/CategoriesWidget.vue'
import ServerStatusWidget from './sidebar/ServerStatusWidget.vue'
import TopContributorsWidget from './sidebar/TopContributorsWidget.vue'
</script>
