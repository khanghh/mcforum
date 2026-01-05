<template>
  <div class="space-y-4">
    <LoadMoreAsync v-slot="{ items }" :cursor="activitiesCursor">
      <ActivityCard v-for="item in items" :key="item.id" :message="item" />
    </LoadMoreAsync>
  </div>
</template>

<script setup lang="ts">
import ActivityCard from '~/components/profile/ActivityCard.vue'

const api = useApi()

const props = defineProps({
  user: {
    type: Object,
    required: true
  }
})

const activitiesCursor = api.getUserActivities(props.user.username) 
</script>

<style scoped>
/* No local styles */
</style>
