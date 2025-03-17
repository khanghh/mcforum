<template>
  <div class="widget">
    <div class="widget-header">
      <div>
        <span>{{ $t('profile.following.title') }}</span>
        <span class="count">{{ user.followCount }}</span>
      </div>
      <div class="slot">
        <nuxt-link :to="`/user/${user.id}/followed`">{{ $t('widget.view_more') }}</nuxt-link>
      </div>
    </div>
    <div class="widget-content">
      <div v-if="followList && followList.length">
        <user-follow-list :users="followList" @onFollowed="onFollowed" />
      </div>
      <div v-else class="widget-tips">{{ $t('profile.following.no_data') }}</div>
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

const followList = ref([]);

onMounted(() => {
  loadData();
});

async function loadData() {
  const data = await useHttpGet(
    `/api/fans/recent/follow?userId=${props.user.id}`
  );
  followList.value = data.results;
}

async function onFollowed(userId, followed) {
  await loadData();
}
</script>

<style lang="scss"></style>
