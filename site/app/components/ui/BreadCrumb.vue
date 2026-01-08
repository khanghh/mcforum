<template>
  <div class="mb-6">
    <nav class="flex text-sm text-gray-400">
      <template v-for="(item, index) in items" :key="index">
        <component :is="item.to ? 'nuxt-link' : 'span'"
          :to="item.to"
          :class="item.to ? 'hover:text-purple-400 transition-colors' : 'text-purple-400 truncate max-w-md'">
          {{ item.label }}
        </component>
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