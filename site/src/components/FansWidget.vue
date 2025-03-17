<template>
  <div class="widget">
    <div class="widget-header">
      <div>
        <span>{{ $t('profile.followers.title') }}</span>
        <span class="count">{{ user.fansCount }}</span>
      </div>
      <div class="slot">
        <nuxt-link :to="`/user/${user.id}/fans`">{{ $t('widget.view_more') }}</nuxt-link>
      </div>
    </div>
    <div class="widget-content">
      <div v-if="fansList && fansList.length">
        <user-follow-list :users="fansList" @onFollowed="onFollowed" />
      </div>
      <div v-else class="widget-tips">{{ $t('profile.followers.no_data') }}</div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  user: {
    type: Object,
    required: true,
  },
});

const fansList = ref([]);

onMounted(() => {
  loadData();
});

async function loadData() {
  const data = await useHttpGet(
    `/api/fans/recent/fans?userId=${props.user.id}`
  );
  fansList.value = data.results;
}

async function onFollowed(userId, followed) {
  await loadData();
}
</script>

<style lang="scss"></style>
