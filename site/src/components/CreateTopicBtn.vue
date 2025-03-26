<template>
  <el-button type="primary" class="publish-btn" size="large" @click="createTopic">
    <icon name="Plus" />&nbsp;
    {{ $t('navbar.create_topic') }}
  </el-button>
</template>

<script setup>
const i18n = useI18n()
const router = useRouter()

function createTopic() {
  router.push('/t/create')
}
</script>

<style lang="scss" scoped>
.publish-btn {
  width: 100%;
  margin-bottom: 10px;
}
</style>
