<template>
  <el-button type="primary" class="publish-btn" size="large" @click="createTopic">
    <icon name="SquarePen" />&nbsp;
    {{ $t('actions.create_topic') }}
  </el-button>
</template>

<script setup>
const i18n = useI18n()
const router = useRouter()

function createTopic() {
  router.push('/topics/create')
}
</script>

<style lang="scss" scoped>
.publish-btn {
  width: 100%;
  margin-bottom: 10px;
}
</style>
