<template>
  <div class="mb-6">
    <nav class="flex text-sm text-gray-400">
      <template v-for="(item, index) in items" :key="index">
        <NuxtLink
          v-if="item.to"
          :to="item.to"
          class="hover:text-purple-400 transition-colors whitespace-nowrap">
          {{ item.label }}
        </NuxtLink>
        <span v-else class="text-purple-400 truncate max-w-md whitespace-nowrap">
          {{ item.label }}
        </span>
        <span v-if="index < items.length - 1" class="mx-2">›</span>
      </template>
    </nav>
  </div>
</template>

<script setup lang="ts">
interface BreadCrumbItem {
  label: string
  to?: string
}

withDefaults(defineProps<{
  items: BreadCrumbItem[]
}>(), {})
</script>