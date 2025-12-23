<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" url="/api/topic/users/favorites" :params="{ userId: user.id }">
      <GamingTopicList :topics="items" :show-pinned="false" />
    </LoadMoreAsync>
  </div>
</template>

<script setup>
defineProps({ user: { type: Object, required: true } })
</script>

<style scoped>
/* No local styles */
</style>
